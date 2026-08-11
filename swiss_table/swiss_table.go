package swiss_table

import "math/bits"

const (
	groupSize         = 8
	ctrlEmpty         = 0b10000000
	ctrlDeleted       = 0b11111110
	maxAvgGroupLoad   = 7
	maxBucketCapacity = 64
)

type slot[K comparable, V any] struct {
	key K
	val V
}

func broadcast(b byte) uint64 { return uint64(b) * 0x0101010101010101 }

func matchByte(word uint64, b byte) uint64 {
	x := word ^ broadcast(b)
	return (x - 0x0101010101010101) &^ x & 0x8080808080808080
}

func setByte(word uint64, pos int, b byte) uint64 {
	shift := uint(pos) * 8
	return (word &^ (0xFF << shift)) | (uint64(b) << shift)
}

func byteAt(word uint64, pos int) byte { return byte(word >> (uint(pos) * 8)) }

// firstVerified: matchByte даёт ложные срабатывания (borrow-цепочка при вычитании
// всего 64-битного слова разом), поэтому позиции из mask перепроверяются побайтово.
// firstVerified: matchByte can false-positive (borrow chain from subtracting the whole
// 64-bit word at once), so positions from mask are re-checked byte by byte.
func firstVerified(word uint64, mask uint64, want byte) (pos int, ok bool) {
	for mask != 0 {
		p := bits.TrailingZeros64(mask) / 8
		if byteAt(word, p) == want {
			return p, true
		}
		mask &= mask - 1
	}
	return 0, false
}

// splitHash: h1 — верхние 57 бит (индекс в директории + затравка пробинга),
// h2 — младшие 7 бит (кладутся в control-байт).
// splitHash: h1 — upper 57 bits (directory index + probe seed),
// h2 — lower 7 bits (stored as the control byte).
func splitHash(hash uint64) (h1 uint64, h2 byte) {
	return hash >> 7, byte(hash & 0x7F)
}

// dirIndex = hash >> (64-globalDepth): те же top globalDepth бит, что и в cockroachdb/swiss.
// dirIndex = hash >> (64-globalDepth): same top-globalDepth-bits rule as cockroachdb/swiss.
func dirIndex(hash uint64, globalDepth uint) uint64 {
	if globalDepth == 0 {
		return 0
	}
	return hash >> (64 - globalDepth)
}

// probeSeq — квадратичный пробинг между группами: смещения — треугольные числа
// (0,1,3,6,10,...) вместо "+1", убирает первичную кластеризацию.
// probeSeq — quadratic probing across groups: triangular-number offsets
// (0,1,3,6,10,...) instead of "+1", removes primary clustering.
type probeSeq struct {
	mask   uint64
	offset uint64
	index  uint64
}

func makeProbeSeq(h1 uint64, mask uint64) probeSeq {
	return probeSeq{mask: mask, offset: h1 & mask}
}

func (s probeSeq) next() probeSeq {
	s.index++
	s.offset = (s.offset + s.index) & s.mask
	return s
}

// bucket — отдельная swiss-таблица, адресуется через extendible hashing.
// Растёт удвоением до maxBucketCapacity, дальше — split на два bucket'а
// с localDepth+1 вместо рехэша всей карты.
// bucket — a standalone swiss table, addressed via extendible hashing.
// Doubles in place up to maxBucketCapacity, then splits into two buckets
// with localDepth+1 instead of rehashing the whole map.
type bucket[K comparable, V any] struct {
	ctrl   []uint64
	slots  []slot[K, V]
	groups int
	used   int

	// growthLeft — число ещё доступных НАСТОЯЩИХ empty-слотов. Tombstone его
	// не восполняет: иначе bucket может физически забиться deleted-байтами
	// при малом used и никогда не вырасти → Get/Put зациклятся, не найдя ни
	// одного ctrlEmpty.
	// growthLeft — count of remaining REAL empty slots. Tombstones don't refill it:
	// otherwise a bucket can fill up with deleted bytes while used stays low and
	// never grow → Get/Put would loop forever finding no ctrlEmpty.
	growthLeft int
	localDepth uint
}

func newBucket[K comparable, V any](groups int, localDepth uint) *bucket[K, V] {
	b := &bucket[K, V]{
		ctrl:       make([]uint64, groups),
		slots:      make([]slot[K, V], groups*groupSize),
		groups:     groups,
		localDepth: localDepth,
	}
	b.growthLeft = (b.capacity() * maxAvgGroupLoad) / groupSize
	empty := broadcast(ctrlEmpty)
	for i := range b.ctrl {
		b.ctrl[i] = empty
	}
	return b
}

func (b *bucket[K, V]) capacity() int { return b.groups * groupSize }

func (b *bucket[K, V]) get(h1 uint64, h2 byte, key K) (val V, ok bool) {
	seq := makeProbeSeq(h1, uint64(b.groups-1))
	// probe < b.groups — явная граница: квадратичная последовательность по модулю
	// степени двойки обходит каждую группу ровно раз за b.groups шагов, поэтому
	// цикл гарантированно конечен. Страховка на случай нарушения инварианта growthLeft.

	// probe < b.groups — explicit bound: quadratic probing modulo a power of two
	// visits every group exactly once in b.groups steps, so the loop always
	// terminates. Safety net in case the growthLeft invariant is ever violated.
	for probe := 0; probe < b.groups; probe++ {
		word := b.ctrl[seq.offset]

		for matches := matchByte(word, h2); matches != 0; matches &= matches - 1 {
			pos := bits.TrailingZeros64(matches) / 8
			idx := int(seq.offset)*groupSize + pos
			if b.slots[idx].key == key {
				return b.slots[idx].val, true
			}
		}
		if empty := matchByte(word, ctrlEmpty); empty != 0 {
			if _, ok := firstVerified(word, empty, ctrlEmpty); ok {
				var zero V
				return zero, false
			}
		}
		seq = seq.next()
	}
	var zero V
	return zero, false
}

// putIfPresent обновляет значение существующего ключа; false, если дошли до
// пустого слота (ключа в bucket'е нет).
// putIfPresent updates the value of an existing key; false once an empty slot
// is reached (key not present).
func (b *bucket[K, V]) putIfPresent(h1 uint64, h2 byte, key K, val V) bool {
	seq := makeProbeSeq(h1, uint64(b.groups-1))
	for probe := 0; probe < b.groups; probe++ {
		word := b.ctrl[seq.offset]

		for matches := matchByte(word, h2); matches != 0; matches &= matches - 1 {
			pos := bits.TrailingZeros64(matches) / 8
			idx := int(seq.offset)*groupSize + pos
			if b.slots[idx].key == key {
				b.slots[idx].val = val
				return true
			}
		}
		if empty := matchByte(word, ctrlEmpty); empty != 0 {
			if _, ok := firstVerified(word, empty, ctrlEmpty); ok {
				return false
			}
		}
		seq = seq.next()
	}
	return false
}

// insert вставляет новую пару, предполагая, что ключа ещё нет (проверка — на
// вызывающей стороне). false = свободного слота нет, bucket нужно расщепить/растить.
// insert adds a new pair assuming the key is absent (caller's responsibility to check).
// false = no free slot, bucket must be grown or split.
func (b *bucket[K, V]) insert(h1 uint64, h2 byte, key K, val V) bool {
	seq := makeProbeSeq(h1, uint64(b.groups-1))
	delGroup, delPos := -1, -1

	for probe := 0; probe < b.groups; probe++ {
		g := int(seq.offset)
		word := b.ctrl[g]

		if del := matchByte(word, ctrlDeleted); del != 0 && delGroup == -1 {
			if p, ok := firstVerified(word, del, ctrlDeleted); ok {
				delGroup, delPos = g, p
			}
		}

		if empty := matchByte(word, ctrlEmpty); empty != 0 {
			if p, ok := firstVerified(word, empty, ctrlEmpty); ok {
				tg, tp := g, p
				if delGroup != -1 {
					tg, tp = delGroup, delPos
				}
				idx := tg*groupSize + tp
				if delGroup == -1 {
					b.growthLeft--
				}
				b.ctrl[tg] = setByte(b.ctrl[tg], tp, h2)
				b.slots[idx] = slot[K, V]{key: key, val: val}
				b.used++
				return true
			}
		}
		seq = seq.next()
	}
	return false
}

// delete: если в группе уже есть empty-байт, группа никогда не была полностью
// заполненной и не участвует в обрыве цепочки пробинга — слот сразу помечается
// empty (downgrade), иначе ставится tombstone, чтобы не сломать пробинг для
// ключей, чья цепочка проходит через этот слот.
// delete: if the group already has an empty byte, it was never fully occupied
// and doesn't gate probe termination — the slot is downgraded straight to empty;
// otherwise it becomes a tombstone so probing for other keys through this slot
// still works.
func (b *bucket[K, V]) delete(h1 uint64, h2 byte, key K) bool {
	seq := makeProbeSeq(h1, uint64(b.groups-1))
	for probe := 0; probe < b.groups; probe++ {
		g := int(seq.offset)
		word := b.ctrl[g]

		for matches := matchByte(word, h2); matches != 0; matches &= matches - 1 {
			pos := bits.TrailingZeros64(matches) / 8
			idx := g*groupSize + pos
			if b.slots[idx].key == key {
				if empty := matchByte(word, ctrlEmpty); empty != 0 {
					if _, ok := firstVerified(word, empty, ctrlEmpty); ok {
						b.ctrl[g] = setByte(word, pos, ctrlEmpty)
						b.growthLeft++
					} else {
						b.ctrl[g] = setByte(word, pos, ctrlDeleted)
					}
				} else {
					b.ctrl[g] = setByte(word, pos, ctrlDeleted)
				}
				var zero slot[K, V]
				b.slots[idx] = zero
				b.used--
				return true
			}
		}
		if empty := matchByte(word, ctrlEmpty); empty != 0 {
			if _, ok := firstVerified(word, empty, ctrlEmpty); ok {
				return false
			}
		}
		seq = seq.next()
	}
	return false
}

// rehashInPlace переносит живые записи в новый bucket той же localDepth, вдвое
// больше — убирает tombstone'ы без расщепления по биту хэша.
// rehashInPlace copies live entries into a new bucket, same localDepth, double
// size — clears tombstones without splitting on a hash bit.
func (b *bucket[K, V]) rehashInPlace(hash func(K) uint64) *bucket[K, V] {
	grown := newBucket[K, V](b.groups*2, b.localDepth)
	for g := 0; g < b.groups; g++ {
		word := b.ctrl[g]
		for pos := 0; pos < groupSize; pos++ {
			c := byteAt(word, pos)
			if c != ctrlEmpty && c != ctrlDeleted {
				idx := g*groupSize + pos
				h1, h2 := splitHash(hash(b.slots[idx].key))
				grown.insert(h1, h2, b.slots[idx].key, b.slots[idx].val)
			}
		}
	}
	return grown
}

// split делит bucket на два с localDepth+1 по значению бита hash на позиции
// bitPos = 64-newDepth (первый ещё не использованный директорией бит).
// split partitions a bucket into two with localDepth+1, keyed on the hash bit
// at bitPos = 64-newDepth (the first bit not yet consumed by the directory).
func (b *bucket[K, V]) split(hash func(K) uint64) (lo, hi *bucket[K, V]) {
	newDepth := b.localDepth + 1
	groups := b.groups
	if groups < 1 {
		groups = 1
	}
	lo = newBucket[K, V](groups, newDepth)
	hi = newBucket[K, V](groups, newDepth)

	bitPos := 64 - newDepth
	for g := 0; g < b.groups; g++ {
		word := b.ctrl[g]
		for pos := 0; pos < groupSize; pos++ {
			c := byteAt(word, pos)
			if c != ctrlEmpty && c != ctrlDeleted {
				idx := g*groupSize + pos
				key, val := b.slots[idx].key, b.slots[idx].val
				full := hash(key)
				h1, h2 := splitHash(full)
				dst := lo
				if (full>>bitPos)&1 == 1 {
					dst = hi
				}
				if !dst.insert(h1, h2, key, val) {
					// dst сам переполнился при перераспределении — растим его на месте.
					// dst overflowed during redistribution — grow it in place.
					dst = dst.rehashInPlace(hash)
					dst.insert(h1, h2, key, val)
				}
				if (full>>bitPos)&1 == 1 {
					hi = dst
				} else {
					lo = dst
				}
			}
		}
	}
	return lo, hi
}

// FastSwissMap — swiss table с квадратичным пробингом и extendible hashing:
// рост = split одного переполненного bucket'а, а не рехэш всей карты.
// FastSwissMap — swiss table with quadratic probing and extendible hashing:
// growth means splitting one overflowing bucket, not rehashing the whole map.
type FastSwissMap[K comparable, V any] struct {
	dir         []*bucket[K, V]
	globalDepth uint
	count       int
	hash        func(K) uint64
}

func NewFastSwissMap[K comparable, V any](hash func(K) uint64) *FastSwissMap[K, V] {
	b := newBucket[K, V](1, 0)
	return &FastSwissMap[K, V]{
		dir:  []*bucket[K, V]{b},
		hash: hash,
	}
}

func (m *FastSwissMap[K, V]) bucketFor(hash uint64) (idx uint64, b *bucket[K, V]) {
	idx = dirIndex(hash, m.globalDepth)
	return idx, m.dir[idx]
}

func (m *FastSwissMap[K, V]) Get(key K) (V, bool) {
	full := m.hash(key)
	h1, h2 := splitHash(full)
	_, b := m.bucketFor(full)
	return b.get(h1, h2, key)
}

func (m *FastSwissMap[K, V]) Put(key K, val V) {
	full := m.hash(key)
	h1, h2 := splitHash(full)

	idx, b := m.bucketFor(full)
	if b.putIfPresent(h1, h2, key, val) {
		return
	}

	// growthLeft>0 гарантирует хотя бы один настоящий empty где-то в bucket'е —
	// insert() пройдёт по всем группам (probeSeq — биекция по modulo) и найдёт его.

	// growthLeft>0 guarantees at least one real empty somewhere in the bucket —
	// insert() sweeps all groups (probeSeq is a bijection mod groups) and finds it.
	if b.growthLeft > 0 {
		if b.insert(h1, h2, key, val) {
			m.count++
			return
		}
	}

	if b.capacity() < maxBucketCapacity {
		// Не достигли предела bucket'а — растим удвоением на месте (заодно чистим tombstone'ы).
		// Below the bucket size limit — grow in place by doubling (also clears tombstones).
		grown := b.rehashInPlace(m.hash)
		m.installBucket(idx, grown)
		grown.insert(h1, h2, key, val)
		m.count++
		return
	}

	// Bucket достиг maxBucketCapacity — расщепляем вместо рехэша всей карты.
	// Bucket hit maxBucketCapacity — split instead of rehashing the whole map.
	m.splitBucket(idx, b)
	m.Put(key, val) // директория и bucket изменились — вставляем заново / directory and bucket changed — retry insertion
}

func (m *FastSwissMap[K, V]) Delete(key K) bool {
	full := m.hash(key)
	h1, h2 := splitHash(full)
	_, b := m.bucketFor(full)
	if b.delete(h1, h2, key) {
		m.count--
		return true
	}
	return false
}

func (m *FastSwissMap[K, V]) Len() int { return m.count }

// installBucket пишет bucket во все 2^(globalDepth-localDepth) ячеек директории,
// которые на него указывают.
// installBucket writes the bucket into all 2^(globalDepth-localDepth) directory
// slots that reference it.
func (m *FastSwissMap[K, V]) installBucket(anyIdx uint64, b *bucket[K, V]) {
	step := uint64(1) << (m.globalDepth - b.localDepth)
	start := (anyIdx / step) * step
	for i := uint64(0); i < step; i++ {
		m.dir[start+i] = b
	}
}

// splitBucket расщепляет переполненный bucket на два. Если его localDepth уже
// равна globalDepth, директория сначала удваивается.
// splitBucket splits an overflowing bucket into two. If its localDepth already
// equals globalDepth, the directory is doubled first.
func (m *FastSwissMap[K, V]) splitBucket(idx uint64, b *bucket[K, V]) {
	// origLocalDepth — глубина ДО расщепления: задаёт суммарное число ячеек
	// директории для lo+hi. Глубина уже созданных lo/hi (localDepth+1) для этого
	// не годится — с ней часть ячеек останется не переписана и укажет на
	// выброшенный bucket (был баг именно на этом месте).

	// origLocalDepth — depth BEFORE the split: defines the total directory slot
	// count for lo+hi combined. Using the already-incremented lo/hi depth here is
	// wrong — some slots would stay unwritten and point at the discarded bucket
	// (this exact spot was the bug).
	origLocalDepth := b.localDepth
	if origLocalDepth == m.globalDepth {
		m.growDirectory()
		// Директория удвоилась: старый индекс i -> пара новых {2i, 2i+1}.
		// Directory doubled: old index i -> new pair {2i, 2i+1}.
		idx *= 2
	}
	lo, hi := b.split(m.hash)

	step := uint64(1) << (m.globalDepth - origLocalDepth)
	start := (idx / step) * step
	half := step / 2
	for i := uint64(0); i < half; i++ {
		m.dir[start+i] = lo
	}
	for i := half; i < step; i++ {
		m.dir[start+i] = hi
	}
}

func (m *FastSwissMap[K, V]) growDirectory() {
	newDir := make([]*bucket[K, V], len(m.dir)*2)
	for i, b := range m.dir {
		newDir[2*i] = b
		newDir[2*i+1] = b
	}
	m.dir = newDir
	m.globalDepth++
}
