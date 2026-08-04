# Stack

<p style="text-align: left">
  <a href="#русский">Русский</a> ・ <a href="#english">English</a>
</p>

---

## Русский

Стек — линейная структура данных, работающая по принципу **LIFO** (Last In, First Out): последний вошедший элемент выходит первым.

Пример: если добавить элементы в порядке `1, 2, 3`, то `3` окажется наверху и будет извлечён первым, следом `2`, затем `1`. Порядок входа и выхода полностью противоположный.

### Внутреннее устройство

Реализация хранит срез `arr` и число `top` — индекс последнего добавленного элемента. Активная часть стека — это `arr[1..top]`, где:

- `arr[1]` — нижний, самый первый добавленный элемент
- `arr[top]` — вершина стека, элемент, добавленный последним

Если `top == 0`, стек пуст. Если `top == n` (максимальная вместимость), попытка добавить новый элемент приводит к переполнению.

Срез создаётся размером `n+1`, а не `n` — нулевая ячейка `arr[0]` физически существует, но не используется, поскольку индексация ведётся с 1, как в оригинальном псевдокоде.

### Использование

```go
s := NewStack(5)

s.Push(1)
s.Push(2)
s.Push(3)

val, err := s.Pop() // val = 3, err = nil
```

### Операции

| Операция | Сложность | Описание |
|---|---|---|
| `Push(x)` | O(1) | добавить элемент на вершину |
| `Pop()` | O(1) | снять и вернуть элемент с вершины |
| `IsEmpty()` | O(1) | проверка на пустоту |
| `IsFull()` | O(1) | проверка на переполнение |
| `Size()` | O(1) | текущее количество элементов |

### Сборка и тестирование

```bash
go test -v ./...
```

Корректность проверяется stress-тестом — реализация сравнивается с обычным срезом, работающим как стек, на большом числе случайных последовательностей операций.

---

## English

A stack is a linear data structure that follows the **LIFO** principle (Last In, First Out): the element inserted last is the first one to leave.

Example: inserting elements in the order `1, 2, 3` places `3` on top, so it comes out first, followed by `2`, then `1`. The exit order is the exact reverse of the entry order.

### Internal layout

The implementation keeps a slice `arr` and a number `top` — the index of the most recently inserted element. The active part of the stack is `arr[1..top]`, where:

- `arr[1]` is the bottom, the very first element inserted
- `arr[top]` is the top of the stack, the most recently inserted element

If `top == 0`, the stack is empty. If `top == n` (the fixed capacity), inserting another element causes an overflow.

The slice is allocated with size `n+1`, not `n` — the zero cell `arr[0]` exists but is never used, since indexing starts at 1, matching the original pseudocode.

### Usage

```go
s := NewStack(5)

s.Push(1)
s.Push(2)
s.Push(3)

val, err := s.Pop() // val = 3, err = nil
```

### Operations

| Operation | Complexity | Description |
|---|---|---|
| `Push(x)` | O(1) | insert an element on top |
| `Pop()` | O(1) | remove and return the top element |
| `IsEmpty()` | O(1) | check whether the stack is empty |
| `IsFull()` | O(1) | check whether the stack is at capacity |
| `Size()` | O(1) | current number of elements |

### Build and test

```bash
go test -v ./...
```

Correctness is checked with a stress test — the implementation is compared against a plain slice used as a stack, over a large number of randomized operation sequences.

---

<br>

> Стек не помнит, кто был первым. Он помнит только, кто был последним — и отдаёт именно его. Так устроена память толпы, поколений, целых цивилизаций: наверху всегда то, что случилось недавно, а самое первое, изначальное, погребено глубже всего и достаётся последним, если вообще достаётся.
>
> *A stack has no memory of who came first. It only remembers who came last — and gives that one back. This is how the memory of a crowd works, of generations, of entire civilizations: what happened recently always sits on top, while the very first, the original, lies buried deepest and is reached last, if it is reached at all.*