# Linked List

<p style="text-align: left">
  <a href="#русский">Русский</a> ・ <a href="#english">English</a>
</p>

---

## Русский

Связный список — линейная структура данных, в которой элементы (узлы) хранятся не в непрерывном блоке памяти, а по отдельности и связаны между собой через указатели. Каждый узел знает только адрес следующего узла — этим список отличается от массива или среза, где все элементы лежат подряд.

Реализация в этом пакете — односвязный список: переход возможен только в одну сторону, от начала к концу.

### Внутреннее устройство

Список хранит три поля: `head`, `tail` и `size`.

- `head` — указатель на первый узел списка.
- `tail` — указатель на последний узел списка. Он нужен, чтобы добавлять элемент в конец за O(1), не пробегая весь список каждый раз.
- `size` — текущее количество элементов.

Каждый узел (`Node[T]`) состоит из значения `Value` и указателя `Next` на следующий узел. У последнего узла `Next` равен `nil` — это признак конца списка.

Если список пуст, оба указателя, `head` и `tail`, равны `nil`.

### Добавление и удаление

- `Prepend` создаёт новый узел и ставит его перед текущим `head`, а сам становится новым `head`. Если список был пуст, этот же узел становится и `tail`.
- `Append` создаёт узел и подвешивает его к `tail.Next`, после чего `tail` сдвигается на него. Если список был пуст, новый узел становится и `head`, и `tail`.
- `InsertAt` вставляет значение на произвольную позицию. Для позиций `0` и `size` используются `Prepend` и `Append`, для остальных случаев список проходится до узла перед нужной позицией, и указатели переставляются.
- `RemoveFirst` возвращает значение из `head` и сдвигает `head` на следующий узел. Если после этого список опустел, `tail` тоже обнуляется.
- `RemoveLast` вынужден пройти список от начала до предпоследнего узла, поскольку у узлов нет указателя назад, на предыдущий элемент. Из-за этого удаление последнего элемента стоит O(n), в отличие от добавления.
- `Remove` ищет первый узел с заданным значением и вырезает его из цепочки, переставляя `Next` у предыдущего узла.

### Поиск и обход

`Search` и `Get` проходят список последовательно от `head`, пока не найдут нужный элемент или не дойдут до конца — прямого доступа по индексу, как в массиве, у связного списка нет.

`Traverse` принимает функцию-колбэк и вызывает её для значения каждого узла по порядку.

`Reverse` разворачивает список на месте: проходит по узлам, у каждого перекладывает `Next` так, чтобы он указывал на предыдущий узел, и в конце меняет местами `head` и `tail`.

`ToSlice` собирает значения всех узлов в обычный срез — удобно для отладки или передачи данных туда, где нужен произвольный доступ по индексу.

### Использование

```go
list := New[int]()
list.Append(1)
list.Append(2)
list.Prepend(0)

val, _ := list.Get(1) // val = 1
list.Reverse()
slice := list.ToSlice() // [2, 1, 0]
```

### Операции

| Операция | Сложность | Описание |
|---|---|---|
| `Prepend(x)` | O(1) | добавить элемент в начало списка |
| `Append(x)` | O(1) | добавить элемент в конец списка |
| `InsertAt(pos, x)` | O(n) | вставить элемент на произвольную позицию |
| `RemoveFirst()` | O(1) | удалить и вернуть первый элемент |
| `RemoveLast()` | O(n) | удалить и вернуть последний элемент |
| `Remove(x)` | O(n) | удалить первый узел с заданным значением |
| `Search(x)` | O(n) | проверить, есть ли значение в списке |
| `Get(pos)` | O(n) | получить значение по позиции |
| `Traverse(fn)` | O(n) | пройти список и вызвать функцию для каждого значения |
| `Reverse()` | O(n) | развернуть список на месте |
| `ToSlice()` | O(n) | получить список значений в виде среза |
| `Len()` | O(1) | текущее количество элементов |
| `IsEmpty()` | O(1) | проверка, пуст ли список |

### Сборка и тестирование

```bash
go test -v ./...
```

---

## English

A linked list is a linear data structure where elements (nodes) are not stored in one contiguous block of memory, but individually, connected to each other through pointers. Each node only knows the address of the next node — this is what sets it apart from an array or a slice, where all elements sit next to each other.

The implementation in this package is a singly linked list: traversal is only possible in one direction, from the start toward the end.

### Internal layout

The list keeps three fields: `head`, `tail`, and `size`.

- `head` points to the first node of the list.
- `tail` points to the last node of the list. It exists so that adding an element to the end costs O(1) instead of walking the whole list every time.
- `size` holds the current number of elements.

Each node (`Node[T]`) consists of a value `Value` and a pointer `Next` to the following node. The last node has `Next` set to `nil` — that is the marker for the end of the list.

When the list is empty, both `head` and `tail` are `nil`.

### Insertion and removal

- `Prepend` creates a new node and places it before the current `head`, becoming the new `head`. If the list was empty, this same node also becomes `tail`.
- `Append` creates a node and attaches it to `tail.Next`, then moves `tail` to it. If the list was empty, the new node becomes both `head` and `tail`.
- `InsertAt` inserts a value at an arbitrary position. Positions `0` and `size` fall back to `Prepend` and `Append`; for anything in between, the list is walked up to the node just before the target position, and the pointers are relinked.
- `RemoveFirst` returns the value at `head` and moves `head` to the next node. If the list becomes empty as a result, `tail` is cleared too.
- `RemoveLast` has to walk the list from the beginning up to the second-to-last node, since nodes have no pointer back to their predecessor. Because of this, removing the last element costs O(n), unlike appending.
- `Remove` looks for the first node holding the given value and cuts it out of the chain by relinking the previous node's `Next`.

### Search and traversal

`Search` and `Get` walk the list sequentially from `head` until they find the target element or reach the end — a linked list has no direct index-based access the way an array does.

`Traverse` takes a callback function and calls it with the value of every node in order.

`Reverse` reverses the list in place: it walks the nodes, flips each one's `Next` to point at the previous node, and finally swaps `head` and `tail`.

`ToSlice` collects the values of all nodes into a plain slice — useful for debugging or for passing the data to something that needs index-based access.

### Usage

```go
list := New[int]()
list.Append(1)
list.Append(2)
list.Prepend(0)

val, _ := list.Get(1) // val = 1
list.Reverse()
slice := list.ToSlice() // [2, 1, 0]
```

### Operations

| Operation | Complexity | Description |
|---|---|---|
| `Prepend(x)` | O(1) | add an element to the front of the list |
| `Append(x)` | O(1) | add an element to the back of the list |
| `InsertAt(pos, x)` | O(n) | insert an element at an arbitrary position |
| `RemoveFirst()` | O(1) | remove and return the first element |
| `RemoveLast()` | O(n) | remove and return the last element |
| `Remove(x)` | O(n) | remove the first node holding the given value |
| `Search(x)` | O(n) | check whether a value exists in the list |
| `Get(pos)` | O(n) | get the value at a given position |
| `Traverse(fn)` | O(n) | walk the list and call a function for each value |
| `Reverse()` | O(n) | reverse the list in place |
| `ToSlice()` | O(n) | get the list's values as a slice |
| `Len()` | O(1) | current number of elements |
| `IsEmpty()` | O(1) | check whether the list is empty |

### Build and test

```bash
go test -v ./...
```

---

<br>

> Связный список не хранит карту всего пути — только адрес следующего шага. Чтобы дойти до конца, нужно пройти через каждое звено цепи, ни одно не пропустить.
>
> *A linked list keeps no map of the whole path — only the address of the next step. To reach the end, you have to pass through every link in the chain, skipping none.*