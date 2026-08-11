# FastSwissMap

<p style="text-align: left">
  <a href="#русский">Русский</a> ・ <a href="#english">English</a>
</p>

---

## Русский

FastSwissMap — хеш-таблица общего назначения, построенная на двух идеях: swiss table (батчевый поиск по control-байтам внутри группы из 8 слотов) и extendible hashing (директория указателей на bucket'ы, растущая делением, а не полным рехэшем).

Обычная swiss table при переполнении рехэширует всю таблицу целиком. Здесь вместо этого каждый bucket — отдельная небольшая swiss table. Пока bucket не достиг `maxBucketCapacity` (64 слота), он просто удваивается на месте. Когда предел достигнут, bucket расщепляется на два по одному биту хэша, а рехэшу подвергаются только записи этого одного bucket'а, а не вся карта.

### Внутреннее устройство

**Control-байты и группы.** Каждая группа занимает один `uint64` (`groupSize = 8`, по байту на слот). Пустой слот кодируется как `0b10000000`, удалённый (tombstone) — как `0b11111110`. Поиск совпадающего байта в группе делает `matchByte` — трюк SWAR (широкое XOR + вычитание + маска знаковых битов), который находит все позиции сразу, без цикла по восьми байтам. У этого трюка есть цена: вычитание всего 64-битного слова разом может дать ложное совпадение из-за переноса между байтами (borrow chain), поэтому кандидаты из маски перепроверяются побайтово в `firstVerified`.

**Разбиение хэша.** `splitHash` делит 64-битный хэш на `h1` (верхние 57 бит — индекс в директории плюс затравка пробинга внутри bucket'а) и `h2` (младшие 7 бит — то, что кладётся в control-байт). Индекс в директории берётся как верхние `globalDepth` бит хэша (`dirIndex`), то есть один и тот же топ-срез бит определяет и позицию в директории, и по какую сторону разреза окажется запись при расщеплении bucket'а.

**Пробинг.** Внутри bucket'а группы перебираются квадратичным пробингом — смещения растут по треугольным числам (0, 1, 3, 6, 10, ...), а не на единицу за шаг, что снижает первичную кластеризацию. Последовательность биективна по модулю числа групп, поэтому цикл `probe < b.groups` гарантированно обходит каждую группу ровно один раз и завершается.

**growthLeft.** Счётчик реально свободных (`ctrlEmpty`) слотов. Удаление не восстанавливает его — только помечает слот как tombstone или, если группа никогда не была заполнена целиком, сразу как `ctrlEmpty`. Если бы tombstone восполнял growthLeft, bucket мог бы физически забиться удалёнными байтами при малом числе живых записей и никогда не вырасти: `Get`/`Put` крутились бы в поиске пустого слота, которого физически не осталось.

**Удаление.** Если в группе целевого слота уже есть хотя бы один `ctrlEmpty`, значит группа никогда не заполнялась полностью и не участвует в обрыве цепочки пробинга — слот сразу понижается до `ctrlEmpty`, а `growthLeft` растёт. Иначе ставится tombstone: без него цепочка пробинга для других ключей, проходящая через этот слот, оборвалась бы преждевременно.

**Рост bucket'а.** Пока `capacity() < maxBucketCapacity`, `Put` при нехватке места вызывает `rehashInPlace` — bucket того же `localDepth`, вдвое больше, без tombstone'ов. Когда предел достигнут, вызывается `splitBucket`: если `localDepth` bucket'а уже равен `globalDepth` директории, директория сперва удваивается (`growDirectory`), затем bucket делится на `lo`/`hi` по значению бита хэша на позиции `64 - newDepth` — первого бита, ещё не использованного директорией. Оба новых bucket'а получают `localDepth + 1`, а во все ячейки директории, ссылавшиеся на старый bucket, записываются `lo` или `hi` в зависимости от значения этого бита.

### Использование

```go
m := NewFastSwissMap[string, int](hashString)
m.Put("a", 1)
m.Put("b", 2)
val, ok := m.Get("a") // val = 1, ok = true
m.Delete("b")
n := m.Len() // n = 1
```

Функция хэширования передаётся снаружи (`hash func(K) uint64`) — таблица не привязана к конкретной хэш-функции.

### Операции

| Операция | Сложность | Описание |
|---|---|---|
| `Get(key)` | O(1) амортизированно | вернуть значение по ключу |
| `Put(key, val)` | O(1) амортизированно | вставить или обновить значение |
| `Delete(key)` | O(1) амортизированно | удалить запись по ключу |
| `Len()` | O(1) | текущее число элементов |

Максимальная средняя загрузка группы — `maxAvgGroupLoad = 7` из 8 слотов (~87.5%), после чего `growthLeft` обнуляется и bucket растёт.

### Сборка и тестирование

```bash
go test -v ./...
```

---

## English

FastSwissMap is a general-purpose hash map combining two ideas: swiss table (batched lookup over control bytes within groups of 8 slots) and extendible hashing (a directory of bucket pointers that grows by splitting, not by a full rehash).

A plain swiss table rehashes everything once it overflows. Here, each bucket is its own small swiss table instead. As long as a bucket is below `maxBucketCapacity` (64 slots), it just doubles in place. Once the limit is hit, the bucket splits into two along a single hash bit, and only that one bucket's entries get rehashed — not the whole map.

### Internal layout

**Control bytes and groups.** Each group occupies one `uint64` (`groupSize = 8`, one byte per slot). An empty slot is encoded as `0b10000000`, a deleted one (tombstone) as `0b11111110`. Matching a byte within a group is done by `matchByte` — a SWAR trick (broadcast XOR, subtract, mask the sign bits) that finds all matching positions at once instead of looping over eight bytes. The trick has a catch: subtracting the whole 64-bit word at once can produce a false positive from a borrow chain crossing byte boundaries, so candidates from the mask get re-verified byte by byte in `firstVerified`.

**Hash split.** `splitHash` divides the 64-bit hash into `h1` (upper 57 bits — directory index plus the probe seed inside the bucket) and `h2` (lower 7 bits — what gets stored in the control byte). The directory index is the top `globalDepth` bits of the hash (`dirIndex`), so the same top-bit slice both selects the directory slot and decides which side of the split an entry lands on when a bucket is divided.

**Probing.** Groups within a bucket are visited with quadratic probing — offsets grow by triangular numbers (0, 1, 3, 6, 10, ...) instead of incrementing by one each step, which avoids primary clustering. The sequence is a bijection modulo the group count, so the `probe < b.groups` loop is guaranteed to visit every group exactly once and terminate.

**growthLeft.** A counter of genuinely empty (`ctrlEmpty`) slots. Deletion never refills it — a slot is marked either as a tombstone or, if its group was never fully occupied, directly as `ctrlEmpty`. If tombstones did refill growthLeft, a bucket could fill up with deleted bytes while the live entry count stayed low and never trigger growth: `Get`/`Put` would spin looking for an empty slot that no longer physically exists.

**Deletion.** If the target slot's group already contains an empty byte, the group was never fully occupied and doesn't gate probe termination — the slot is downgraded straight to `ctrlEmpty` and `growthLeft` increases. Otherwise it becomes a tombstone: without it, probing for other keys whose chain passes through this slot would terminate too early.

**Bucket growth.** While `capacity() < maxBucketCapacity`, a `Put` that runs out of room triggers `rehashInPlace` — a bucket of the same `localDepth`, twice the size, with tombstones cleared. Once the limit is reached, `splitBucket` runs: if the bucket's `localDepth` already equals the directory's `globalDepth`, the directory doubles first (`growDirectory`), then the bucket splits into `lo`/`hi` based on the hash bit at position `64 - newDepth` — the first bit not yet consumed by the directory. Both new buckets get `localDepth + 1`, and every directory slot that pointed at the old bucket is rewritten to `lo` or `hi` depending on that bit.

### Usage

```go
m := NewFastSwissMap[string, int](hashString)
m.Put("a", 1)
m.Put("b", 2)
val, ok := m.Get("a") // val = 1, ok = true
m.Delete("b")
n := m.Len() // n = 1
```

The hash function is supplied by the caller (`hash func(K) uint64`) — the table isn't tied to any specific hashing scheme.

### Operations

| Operation | Complexity | Description |
|---|---|---|
| `Get(key)` | O(1) amortized | return the value for a key |
| `Put(key, val)` | O(1) amortized | insert or update a value |
| `Delete(key)` | O(1) amortized | remove an entry by key |
| `Len()` | O(1) | current number of elements |

Maximum average group load is `maxAvgGroupLoad = 7` out of 8 slots (~87.5%), past which `growthLeft` hits zero and the bucket grows.

### Build and test

```bash
go test -v ./...
```

---

<br>

> Директория не хранит записи — она хранит право спросить нужный bucket. Когда bucket переполняется, делится не карта, а один-единственный узел ответственности.
>
> *The directory holds no entries of its own — only the right to ask the correct bucket. When a bucket overflows, what splits is not the map, but one single unit of responsibility.*