package hash_map

type ProbeStrategy int

const (
	Linear ProbeStrategy = iota
	Quadratic
	Double
)

type HashFunc[K comparable] func(key K) uint64

type HashTable[K comparable, V any] struct {
	slots    []*Node[K, V]
	size     int
	count    int
	strategy ProbeStrategy
	h1       HashFunc[K] // основной хеш — определяет стартовый слот
	h2       HashFunc[K] // вспомогательный хеш — используется только для Double
}

type Node[K comparable, V any] struct {
	key     K
	value   V
	deleted bool
}

// h2 обязателен только для стратегии Double, но проще всегда требовать обе функции,
// чем городить nil-проверки в проде.
func NewHashMap[K comparable, V any](size int, strategy ProbeStrategy, h1, h2 HashFunc[K]) *HashTable[K, V] {
	return &HashTable[K, V]{
		slots:    make([]*Node[K, V], size),
		size:     size,
		strategy: strategy,
		h1:       h1,
		h2:       h2,
	}
}

// baseHash возвращает индекс слота для первой попытки.
func (ht *HashTable[K, V]) baseHash(key K) int {
	return int(ht.h1(key) % uint64(ht.size))
}

// stepHash — шаг для двойного хеширования.
// Возвращает значение в диапазоне [1, size-1], т.е. никогда 0.
func (ht *HashTable[K, V]) stepHash(key K) int {
	step := int(ht.h2(key) % uint64(ht.size-1))
	return step + 1
}

// probe возвращает индекс i-й попытки согласно выбранной стратегии.
func (ht *HashTable[K, V]) probe(base, step, i int) int {
	switch ht.strategy {
	case Linear:
		return (base + i) % ht.size
	case Quadratic:
		return (base + i*i) % ht.size
	case Double:
		return (base + i*step) % ht.size
	default:
		return (base + i) % ht.size
	}
}

func (ht *HashTable[K, V]) Insert(key K, value V) bool {
	if ht.count >= ht.size {
		return false
	}

	base := ht.baseHash(key)
	step := ht.stepHash(key)
	firstFree := -1

	for i := 0; i < ht.size; i++ {
		idx := ht.probe(base, step, i)
		node := ht.slots[idx]

		if node == nil {
			if firstFree == -1 {
				firstFree = idx
			}
			break
		}
		if node.deleted {
			if firstFree == -1 {
				firstFree = idx
			}
			continue
		}
		if node.key == key {
			node.value = value
			return true
		}
	}

	if firstFree == -1 {
		return false
	}

	ht.slots[firstFree] = &Node[K, V]{key: key, value: value}
	ht.count++

	return true
}

func (ht *HashTable[K, V]) Search(key K) (V, bool) {
	base := ht.baseHash(key)
	step := ht.stepHash(key)
	var zeroValue V

	for i := 0; i < ht.size; i++ {
		idx := ht.probe(base, step, i)
		node := ht.slots[idx]

		if node == nil {
			return zeroValue, false
		}
		if !node.deleted && node.key == key {
			return node.value, true
		}
	}

	return zeroValue, false
}

func (ht *HashTable[K, V]) Delete(key K) bool {
	base := ht.baseHash(key)
	step := ht.stepHash(key)

	for i := 0; i < ht.size; i++ {
		idx := ht.probe(base, step, i)
		node := ht.slots[idx]

		if node == nil {
			return false
		}
		if !node.deleted && node.key == key {
			node.deleted = true
			ht.count--
			return true
		}
	}

	return false
}
