# Binary Heap

<p style="text-align: left">
  <a href="#русский">Русский</a> ・ <a href="#english">English</a>
</p>

---

## Русский

Бинарная куча (binary heap) — это почти полное бинарное дерево, в котором каждый родитель имеет приоритет не ниже своих детей. В **min-heap** наверху всегда минимальный элемент, в **max-heap** — максимальный.

Куча не хранит все значения в отсортированном порядке. Она гарантирует только, что лучший по приоритету элемент находится в корне и может быть получен за O(1). Благодаря этому вставка и извлечение приоритетного элемента выполняются за O(log n), а не требуют полной сортировки после каждого изменения.

### Внутреннее устройство

Дерево хранится в непрерывном срезе `data`, без указателей между узлами. Для элемента с индексом `i`:

- родитель находится по индексу `(i-1)/2`;
- левый ребёнок — `2*i+1`;
- правый ребёнок — `2*i+2`.

Например, массив `[1, 3, 2, 8, 5]` представляет min-heap:

```text
        1
      /   \
     3     2
    / \
   8   5
```

Функция `less(a, b)` определяет приоритет. Если она возвращает `true`, значение `a` должно находиться ближе к корню, чем `b`. Поэтому одна реализация поддерживает min-heap, max-heap и пользовательские типы.

**Вставка (`Push`).** Новый элемент добавляется в конец среза, затем поднимается вверх (`siftUp`), пока имеет больший приоритет, чем родитель.

**Извлечение (`Pop`).** Корень сохраняется как результат, последний элемент переносится на его место и опускается вниз (`siftDown`) через ребёнка с более высоким приоритетом.

**Построение (`New`).** Готовый набор значений копируется, после чего внутренние узлы просеиваются снизу вверх. Алгоритм Флойда строит кучу за O(n), что быстрее последовательных O(n log n) вставок.

### Зачем нужна бинарная куча

Главный бизнес-сценарий — **очередь с приоритетом**: сначала обрабатывается не самая старая задача, а самая важная. Примеры: диспетчеризация фоновых заданий, выбор ближайшего события по времени, планирование запросов, обработка инцидентов, алгоритмы Дейкстры и Прима, поиск top-K элементов и объединение отсортированных потоков.

Куча хорошо подходит, когда постоянно нужен один экстремальный элемент. Если требуется быстрый поиск произвольного значения или полностью отсортированный обход без разрушения структуры, нужны другие структуры.

### Использование

Min-heap для целых чисел:

```go
h := NewMinHeap(7, 2, 9)
h.Push(1)

top, _ := h.Peek() // 1; куча не меняется
value, _ := h.Pop() // 1; элемент удалён
size := h.Len()     // 3
```

Max-heap:

```go
h := NewMaxHeap(7, 2, 9)
value, _ := h.Pop() // 9
```

Пользовательский приоритет:

```go
type Job struct {
    Name     string
    Priority int
}

h := New(func(a, b Job) bool {
    return a.Priority > b.Priority
})

h.Push(Job{Name: "email", Priority: 1})
h.Push(Job{Name: "incident", Priority: 10})
next, _ := h.Pop() // incident
```

Пустые `Peek` и `Pop` возвращают `ErrEmpty`. Конструктор паникует при `nil` вместо компаратора, потому что без правила приоритета поддерживать инвариант невозможно.

### Операции и сложность

Пусть `n` — число элементов в куче.

| Операция | Время | Доп. память | Описание |
|---|---:|---:|---|
| `New(less, values...)` | O(n) | O(n) | скопировать значения и построить кучу снизу вверх |
| `NewMinHeap(values...)` | O(n) | O(n) | создать min-heap целых чисел |
| `NewMaxHeap(values...)` | O(n) | O(n) | создать max-heap целых чисел |
| `Push(value)` | O(log n) | O(1)* | добавить значение и поднять его |
| `Peek()` | O(1) | O(1) | вернуть корень без удаления |
| `Pop()` | O(log n) | O(1) | удалить корень и восстановить инвариант |
| `Len()` / `IsEmpty()` | O(1) | O(1) | получить размер или проверить пустоту |

\* Амортизированно вставка может расширить внутренний срез и выделить O(n) памяти, но само просеивание работает на месте.

### Сборка и тестирование

```bash
go test -v ./binary_heap
```

Тесты проверяют min/max-порядок, дубликаты, отрицательные числа, пользовательский тип, пустую кучу, независимость от исходного среза и 20 000 случайных операций со сравнением против эталонной модели.

---

## English

A binary heap is an almost complete binary tree in which every parent has at least as much priority as its children. A **min-heap** always keeps the smallest element on top; a **max-heap** keeps the largest one there.

A heap does not keep every value fully sorted. It only guarantees that the highest-priority value is at the root and can be read in O(1). This makes insertion and priority removal O(log n), without sorting the entire collection after every change.

### Internal layout

The tree is stored in one contiguous `data` slice with no pointers between nodes. For an item at index `i`:

- its parent is at `(i-1)/2`;
- its left child is at `2*i+1`;
- its right child is at `2*i+2`.

For example, `[1, 3, 2, 8, 5]` represents this min-heap:

```text
        1
      /   \
     3     2
    / \
   8   5
```

The `less(a, b)` function defines priority. If it returns `true`, `a` belongs closer to the root than `b`. This lets one implementation support min-heaps, max-heaps, and application-specific types.

**Insertion (`Push`).** A new item is appended to the slice and moves upward (`siftUp`) while it has more priority than its parent.

**Removal (`Pop`).** The root is saved as the result, the final item moves into its place, and then moves downward (`siftDown`) through the higher-priority child.

**Construction (`New`).** The input values are copied, then internal nodes are sifted from the bottom upward. Floyd's algorithm builds the heap in O(n), faster than O(n log n) repeated insertion.

### Why a binary heap is useful

The primary business use case is a **priority queue**: the most important task, rather than the oldest task, is processed first. Examples include background-job dispatch, selecting the next timed event, request scheduling, incident handling, Dijkstra's and Prim's algorithms, top-K queries, and merging sorted streams.

A heap fits workloads that repeatedly need one extreme value. When arbitrary-value lookup or a fully sorted non-destructive traversal is required, a different structure is usually a better choice.

### Usage

An integer min-heap:

```go
h := NewMinHeap(7, 2, 9)
h.Push(1)

top, _ := h.Peek() // 1; heap is unchanged
value, _ := h.Pop() // 1; item is removed
size := h.Len()     // 3
```

A max-heap:

```go
h := NewMaxHeap(7, 2, 9)
value, _ := h.Pop() // 9
```

Application-specific priority:

```go
type Job struct {
    Name     string
    Priority int
}

h := New(func(a, b Job) bool {
    return a.Priority > b.Priority
})

h.Push(Job{Name: "email", Priority: 1})
h.Push(Job{Name: "incident", Priority: 10})
next, _ := h.Pop() // incident
```

Empty `Peek` and `Pop` calls return `ErrEmpty`. The constructor panics for a `nil` comparator because the invariant cannot be maintained without a priority rule.

### Operations and complexity

Let `n` be the number of values in the heap.

| Operation | Time | Extra space | Description |
|---|---:|---:|---|
| `New(less, values...)` | O(n) | O(n) | copy the values and build the heap bottom-up |
| `NewMinHeap(values...)` | O(n) | O(n) | create an integer min-heap |
| `NewMaxHeap(values...)` | O(n) | O(n) | create an integer max-heap |
| `Push(value)` | O(log n) | O(1)* | add a value and move it upward |
| `Peek()` | O(1) | O(1) | return the root without removal |
| `Pop()` | O(log n) | O(1) | remove the root and restore the invariant |
| `Len()` / `IsEmpty()` | O(1) | O(1) | read the size or check for emptiness |

\* Amortized insertion may grow the backing slice and allocate O(n) memory, while the sifting operation itself works in place.

### Build and test

```bash
go test -v ./binary_heap
```

Tests cover min/max ordering, duplicates, negative numbers, a custom type, empty-heap behavior, independence from the input slice, and 20,000 random operations checked against a reference model.

---

<br>

> Куча не знает полного порядка. Она знает только, кто должен быть следующим, — и часто именно этого достаточно для быстрой системы.
>
> *A heap does not know the complete order. It only knows who should come next—and that is often all a fast system needs.*
