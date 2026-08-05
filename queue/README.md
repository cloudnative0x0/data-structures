# Queue

<p style="text-align: left">
  <a href="#русский">Русский</a> ・ <a href="#english">English</a>
</p>

---

## Русский

Очередь — линейная структура данных, работающая по принципу **FIFO** (First In, First Out): элемент, добавленный первым, извлекается первым. Это полная противоположность стеку.

Пример: если добавить элементы в порядке `1, 2, 3`, то первым выйдет `1`, затем `2`, затем `3`. Порядок выхода в точности повторяет порядок входа.

### Внутреннее устройство

Реализация использует кольцевой буфер на основе среза `arr`, два индекса — `head` и `tail`, а также счётчик `count`.

- `head` указывает на ячейку, из которой будет читать следующая операция `Dequeue` или `Peek`.
- `tail` указывает на ячейку, в которую запишет следующая операция `Enqueue`.
- `count` хранит текущее число элементов в очереди. Оно необходимо, чтобы различать состояния «пуста» и «полна», так как при `head == tail` буфер может быть и пустым, и полностью заполненным — только счётчик даёт ответ.

При добавлении элемента (`Enqueue`) значение записывается в `arr[tail]`, после чего `tail` сдвигается циклически вперёд по формуле `(tail + 1) % len(arr)`, а `count` увеличивается на единицу.

При извлечении (`Dequeue`) читается `arr[head]`, затем `head` сдвигается вперёд так же циклически, а `count` уменьшается.

Фиксированная ёмкость очереди задаётся при создании и не меняется. Попытка добавить элемент в заполненную очередь приводит к ошибке, как и попытка извлечь из пустой.

Дополнительно доступны методы `Size()` (возвращает `count`) и `Cap()` (возвращает `len(arr)`), позволяющие узнать текущую заполненность и предельную вместимость.

### Использование

```go
q, err := NewQueue[int](3)
if err != nil {
    // обработать неверную ёмкость
}
q.Enqueue(10)
q.Enqueue(20)
val, _ := q.Dequeue() // val = 10
```

### Операции

| Операция | Сложность | Описание |
|---|---|---|
| `Enqueue(x)` | O(1) | добавить элемент в конец очереди |
| `Dequeue()` | O(1) | извлечь элемент из начала очереди |
| `Peek()` | O(1) | посмотреть элемент в начале, не удаляя |
| `IsEmpty()` | O(1) | проверка, пуста ли очередь |
| `IsFull()` | O(1) | проверка, заполнена ли очередь |
| `Size()` | O(1) | текущее количество элементов |
| `Cap()` | O(1) | максимальная ёмкость |

### Сборка и тестирование

```bash
go test -v ./...
```

Корректность проверяется с помощью набора модульных тестов, покрывающих граничные случаи, циклическое поведение индексов и стресс-теста, который сравнивает поведение очереди с эталонной реализацией на срезе.

---

## English

A queue is a linear data structure that follows the **FIFO** principle (First In, First Out): the element inserted first is the first one to be removed. It is the exact opposite of a stack.

Example: inserting elements in the order `1, 2, 3` will retrieve `1` first, then `2`, then `3`. The order of removal precisely matches the order of insertion.

### Internal layout

The implementation uses a circular buffer built on a slice `arr`, two indices `head` and `tail`, and a counter `count`.

- `head` points to the cell from which the next `Dequeue` or `Peek` operation will read.
- `tail` points to the cell where the next `Enqueue` operation will write.
- `count` holds the current number of elements. It is essential to tell apart the empty and full states, because when `head == tail` the buffer could be either empty or completely full — only the counter provides the answer.

When an element is added (`Enqueue`), the value is stored in `arr[tail]`, then `tail` advances circularly via `(tail + 1) % len(arr)`, and `count` increases by one.

When an element is removed (`Dequeue`), the value at `arr[head]` is read, then `head` advances circularly in the same way, and `count` decreases.

The fixed capacity of the queue is set at creation and never changes. Inserting into a full queue or removing from an empty one both result in errors.

Additionally, `Size()` (returns `count`) and `Cap()` (returns `len(arr)`) are provided to inspect how many elements are currently stored and what the maximum capacity is.

### Usage

```go
q, err := NewQueue[int](3)
if err != nil {
    // handle invalid capacity
}
q.Enqueue(10)
q.Enqueue(20)
val, _ := q.Dequeue() // val = 10
```

### Operations

| Operation | Complexity | Description |
|---|---|---|
| `Enqueue(x)` | O(1) | add an element to the back of the queue |
| `Dequeue()` | O(1) | remove and return the front element |
| `Peek()` | O(1) | return the front element without removing |
| `IsEmpty()` | O(1) | check whether the queue is empty |
| `IsFull()` | O(1) | check whether the queue is at capacity |
| `Size()` | O(1) | current number of elements |
| `Cap()` | O(1) | maximum capacity |

### Build and test

```bash
go test -v ./...
```

Correctness is verified through a suite of unit tests that cover edge cases, circular index wrapping, and a stress test that compares the queue's behavior against a reference slice-based implementation.

---

<br>

> Очередь не смотрит на важность, не слышит просьб пропустить вперёд. Она знает только одно: кто пришёл раньше, тот и уйдёт раньше. В этом её слепая справедливость.
>
> *A queue does not look at importance, nor does it hear pleas to skip ahead. It knows only one thing: whoever came first leaves first. In that lies its blind justice.*