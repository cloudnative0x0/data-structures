# HashMap

<p style="text-align: left">
  <a href="#русский">Русский</a> ・ <a href="#english">English</a>
</p>

---

## Русский

Хеш-таблица — структура данных, хранящая пары ключ-значение и обеспечивающая доступ к ним в среднем за O(1), без перебора всех элементов подряд. Вместо поиска по всей коллекции от ключа вычисляется число (хеш), которое сразу указывает на слот в массиве, где этот ключ должен находиться.

Пример: если положить пары `("a", 1)`, `("b", 2)`, `("c", 3)`, поиск значения по ключу `"b"` не требует просмотра `"a"` и `"c"` — хеш от `"b"` сразу даёт индекс нужного слота, и таблица заглядывает почти сразу в нужное место.

### Внутреннее устройство

Реализация использует open addressing (открытую адресацию): все элементы лежат в одном срезе `slots` фиксированного размера, без связных списков или отдельных бакетов на каждый слот. Если слот, на который указывает хеш ключа, уже занят другим ключом — это коллизия, и таблица последовательно проверяет следующие слоты по формуле пробирования, пока не найдёт либо искомый ключ, либо свободное место.

Поля структуры:

- `slots []*Node[K, V]` — массив слотов; `nil` в слоте означает "сюда ещё никогда не клали элемент"
- `size` — фиксированная вместимость таблицы, задаётся один раз при создании и не меняется — ресайза в этой реализации нет
- `count` — число реально хранимых, неудалённых элементов
- `strategy` — выбранная стратегия пробирования: `Linear`, `Quadratic` или `Double`
- `h1` — основная хеш-функция, определяет стартовый слот для ключа
- `h2` — вспомогательная хеш-функция; используется только стратегией `Double`, задаёт размер шага между попытками

`Node` хранит ключ, значение и флаг `deleted`. Удаление не освобождает слот физически — `deleted` помечает узел как "мёртвый", но занятое место (tombstone). Это принципиально: при поиске таблица понимает, что цепочка коллизий закончилась, только встретив `nil`. Если бы удаление ставило `nil` вместо пометки, все ключи, которые когда-то были вытеснены коллизией за этот слот, стали бы ненаходимыми — они физически остались бы в таблице, но поиск обрывался бы раньше, чем до них доходит.

`baseHash` берёт `h1(key)` и приводит его к диапазону `[0, size)` через остаток от деления — это стартовый слот. `stepHash` берёт `h2(key)` и приводит к диапазону `[1, size-1]`, специально исключая 0: нулевой шаг означал бы, что при `Double` таблица вечно проверяет один и тот же слот и никогда не продвигается дальше.

`probe(base, step, i)` — это и есть формула пробирования, единая точка, где выбор `strategy` превращается в конкретный индекс i-й попытки:

- `Linear`: `base + i` — следующий слот подряд
- `Quadratic`: `base + i²` — расстояние между попытками растёт квадратично, это разбивает длинные цепочки коллизий (кластеризацию), характерные для линейного пробирования
- `Double`: `base + i*step` — шаг зависит от самого ключа через `h2`, поэтому два ключа, столкнувшихся в одном стартовом слоте, с высокой вероятностью пойдут разными путями дальше

`Insert` идёт по формуле пробирования до `size` шагов. По дороге запоминает первый встреченный свободный или помеченный удалённым слот (`firstFree`) — именно туда в итоге ляжет новый элемент, если ключа в таблице ещё нет. Если попадается живой узел с тем же ключом — это не вставка, а обновление значения на месте. Встреченный `nil` останавливает поиск: дальше по цепочке пробирования этого ключа быть не может, потому что при вставке `nil` никогда не пропускается.

`Search` и `Delete` устроены зеркально — идут по той же формуле пробирования, останавливаются на первом `nil`, но, в отличие от `Insert`, не трогают tombstone-узлы: `deleted` для них просто означает "не то, что ищем, идём дальше".

Стоит понимать ограничение этой реализации: таблица не делает rehash и не растёт. `count` уменьшается при удалении, но сам слот остаётся физически занятым tombstone-ом до тех пор, пока `Insert` не перезапишет его новым ключом. При большом числе удалений без последующих вставок таблица может выглядеть почти пустой по `count`, но пробирование всё ещё будет идти через длинные цепочки старых tombstone-ов.

### Использование

```go
h1 := func(s string) uint64 {
	var h uint64 = 14695981039346656037
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= 1099511628211
	}
	return h
}

h2 := func(s string) uint64 {
	var h uint64 = 5381
	for i := 0; i < len(s); i++ {
		h = ((h << 5) + h) + uint64(s[i])
	}
	return h
}

ht := hash_map.NewHashMap[string, int](16, hash_map.Double, h1, h2)

ht.Insert("a", 1)
ht.Insert("b", 2)

val, found := ht.Search("a") // val = 1, found = true

ht.Delete("a")
_, found = ht.Search("a") // found = false
```

`h1` и `h2` должны быть алгоритмически независимыми функциями — если `h2` производна от `h1` предсказуемым образом, шаг пробирования при `Double` будет коррелировать со стартовым слотом, и вся выгода от двойного хеширования потеряется.

### Операции

| Операция | Сложность (средняя) | Сложность (худшая) | Описание |
|---|---|---|---|
| `Insert(key, value)` | O(1) | O(size) | вставить новую пару или обновить значение существующего ключа |
| `Search(key)` | O(1) | O(size) | найти значение по ключу |
| `Delete(key)` | O(1) | O(size) | пометить узел с ключом как удалённый |

Худший случай наступает при высокой заполненности таблицы или неудачном выборе `h1`/`h2`, когда пробирование вынуждено обойти почти все слоты, прежде чем найти нужный или свободный.

### Сборка и тестирование

```bash
go test -v ./...
```

Тесты покрывают все три стратегии пробирования, поведение с принудительными коллизиями (когда `h1` у разных ключей совпадает), поиск и вставку через цепочку tombstone-узлов, полностью заполненную таблицу и работу с ключами разных типов.

---

## English

A hash table is a data structure that stores key-value pairs and provides access to them in O(1) on average, without scanning every element in order. Instead of searching the whole collection, a number (a hash) is computed from the key, and that number points straight at the slot in the array where the key should live.

Example: after inserting `("a", 1)`, `("b", 2)`, `("c", 3)`, looking up `"b"` does not require touching `"a"` or `"c"` — the hash of `"b"` gives the index of the right slot directly, and the table looks almost straight there.

### Internal layout

The implementation uses open addressing: every element lives in a single fixed-size slice `slots`, with no linked lists or per-slot buckets. If the slot a key's hash points to is already taken by another key — a collision — the table walks forward through further slots following a probing formula, until it finds either the key it's looking for or a free spot.

Struct fields:

- `slots []*Node[K, V]` — the array of slots; `nil` in a slot means "an element has never been placed here"
- `size` — the table's fixed capacity, set once at creation and never changed — this implementation has no resizing
- `count` — the number of currently stored, non-deleted elements
- `strategy` — the chosen probing strategy: `Linear`, `Quadratic`, or `Double`
- `h1` — the primary hash function, determines the starting slot for a key
- `h2` — the secondary hash function, used only by the `Double` strategy to set the step size between attempts

`Node` holds a key, a value, and a `deleted` flag. Deletion does not free the slot physically — `deleted` marks the node as dead but still occupying its place (a tombstone). This matters: while probing, the table only knows a collision chain has ended once it hits `nil`. If deletion set the slot to `nil` instead, every key that had once been pushed past that slot by a collision would become unreachable — still physically present in the table, but the search would stop before reaching it.

`baseHash` takes `h1(key)` and reduces it to `[0, size)` via a modulo — that's the starting slot. `stepHash` takes `h2(key)` and reduces it to `[1, size-1]`, deliberately excluding 0: a zero step would mean `Double` keeps checking the same slot forever and never advances.

`probe(base, step, i)` is the probing formula itself, the single place where `strategy` turns into a concrete index for the i-th attempt:

- `Linear`: `base + i` — the next slot in sequence
- `Quadratic`: `base + i²` — the distance between attempts grows quadratically, breaking up the long collision runs (clustering) typical of linear probing
- `Double`: `base + i*step` — the step depends on the key itself through `h2`, so two keys that collide on the same starting slot are likely to take different paths from there

`Insert` follows the probing formula for up to `size` attempts. Along the way it remembers the first free or tombstoned slot it sees (`firstFree`) — that's where the new element ends up if the key isn't already in the table. If it runs into a live node with the same key, that's not an insertion but an update in place. Hitting `nil` stops the search: this key's chain cannot continue past it, since `Insert` never leaves a `nil` gap inside a chain.

`Search` and `Delete` mirror each other — they follow the same probing formula, stop at the first `nil`, but unlike `Insert` they don't touch tombstoned nodes: for them, `deleted` just means "not what we're looking for, keep going".

One limitation worth knowing: the table never rehashes or grows. `count` drops on deletion, but the slot itself stays physically occupied by a tombstone until `Insert` overwrites it with a new key. After many deletions with few new insertions, the table can look nearly empty by `count` while probing still has to walk through long runs of old tombstones.

### Usage

```go
h1 := func(s string) uint64 {
	var h uint64 = 14695981039346656037
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= 1099511628211
	}
	return h
}

h2 := func(s string) uint64 {
	var h uint64 = 5381
	for i := 0; i < len(s); i++ {
		h = ((h << 5) + h) + uint64(s[i])
	}
	return h
}

ht := hash_map.NewHashMap[string, int](16, hash_map.Double, h1, h2)

ht.Insert("a", 1)
ht.Insert("b", 2)

val, found := ht.Search("a") // val = 1, found = true

ht.Delete("a")
_, found = ht.Search("a") // found = false
```

`h1` and `h2` need to be algorithmically independent functions — if `h2` is predictably derived from `h1`, the probing step under `Double` will correlate with the starting slot, and the whole point of double hashing is lost.

### Operations

| Operation | Complexity (average) | Complexity (worst case) | Description |
|---|---|---|---|
| `Insert(key, value)` | O(1) | O(size) | insert a new pair or update an existing key's value |
| `Search(key)` | O(1) | O(size) | look up a value by key |
| `Delete(key)` | O(1) | O(size) | mark the node with this key as deleted |

The worst case shows up when the table is nearly full, or when `h1`/`h2` are poorly chosen, forcing probing to walk through almost every slot before finding the right one or a free one.

### Build and test

```bash
go test -v ./...
```

The tests cover all three probing strategies, behavior under forced collisions (when different keys share the same `h1`), search and insertion through chains of tombstoned nodes, a fully saturated table, and keys of different types.

---

<br>

> Функция хеширования не спрашивает, кто ты. Она сразу назначает тебе место — и если оно занято, тебя сдвигают дальше, по формуле, к которой ты не имеешь отношения. Сосед по адресу — не тот, кто на тебя похож, а тот, кто просто пришёл раньше и попал туда же. Удалённое место не становится пустым: оно остаётся отмеченным, чтобы путь к тем, кто прошёл дальше по цепочке, не оборвался — даже если тебя самого там больше нет.
>
> *A hash function doesn't ask who you are. It assigns you a place right away — and if that place is taken, you get pushed further along, by a formula you had no part in. Your neighbor at that address isn't someone like you, just someone who arrived earlier and landed on the same spot. A deleted place doesn't become empty: it stays marked, so the path to whoever moved further down the chain isn't cut off — even once you yourself are gone.*