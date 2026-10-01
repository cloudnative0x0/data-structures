# Breadth-First Search (BFS) Binary Tree

<p style="text-align: left">
  <a href="#русский">Русский</a> ・ <a href="#english">English</a>
</p>

---

## Русский

Этот пакет содержит обычное бинарное дерево и показывает **поиск в ширину** (Breadth-First Search, BFS). BFS посещает дерево уровень за уровнем: сначала корень, затем всех его детей, потом внуков и так далее. Поэтому алгоритм использует очередь FIFO: первый обнаруженный узел обрабатывается первым.

Для дерева

```text
        1
      /   \
     2     3
    / \   / \
   4   5 6   7
```

порядок BFS будет `1, 2, 3, 4, 5, 6, 7`.

### Внутреннее устройство

`Node` хранит целое значение и две ссылки: `Left` и `Right`. `New` по очереди передаёт значения в `Insert`. Вставка выполняет BFS и помещает новый узел в первую свободную позицию слева направо. Благодаря этому дерево остаётся **полным**: все уровни, кроме последнего, заполнены, а последний заполняется слева направо.

`Search` тоже выполняет BFS и возвращает первый найденный узел. Если значения повторяются, это будет самое верхнее совпадение, а на одном уровне — самое левое.

`Delete` находит целевой узел, копирует в него значение самого глубокого правого узла, а затем удаляет последний узел. Такой способ сохраняет форму полного бинарного дерева, но не поддерживает порядок бинарного дерева поиска: здесь нет правила `left < root < right`.

Помимо BFS пакет предоставляет классические глубинные обходы:

- `TraversePreOrder`: корень → левое поддерево → правое поддерево;
- `TraverseInOrder`: левое поддерево → корень → правое поддерево;
- `TraversePostOrder`: левое поддерево → правое поддерево → корень;
- `TraverseBFS`: уровень за уровнем слева направо.

### Зачем нужен BFS

BFS полезен, когда бизнес-правило зависит от расстояния или уровня: поиск ближайшего сотрудника в иерархии, обход категорий по уровням, минимальное число переходов в невзвешенном графе, поиск ближайшего элемента интерфейса или заполнение дерева без разрывов. Для поиска глубокого решения при ограниченной памяти чаще подходит DFS.

### Использование

```go
root := New([]int{1, 2, 3, 4, 5, 6, 7})

found := root.Search(6) // узел со значением 6

var values []int
root.TraverseBFS(func(v int) {
    values = append(values, v)
}) // [1 2 3 4 5 6 7]

root = Delete(root, 3)
```

### Операции и сложность

Пусть `n` — число узлов, `w` — максимальная ширина дерева, `h` — его высота.

| Операция | Время | Доп. память | Описание |
|---|---:|---:|---|
| `New(values)` | O(n²) | O(w) | текущая реализация вызывает BFS-вставку для каждого значения |
| `Insert(root, value)` | O(n) | O(w) | найти первую свободную позицию по уровням |
| `Search(value)` | O(n) | O(w) | найти первое совпадение в порядке BFS |
| `Delete(root, value)` | O(n) | O(w) | поиск цели и последнего узла |
| `TraverseBFS(visit)` | O(n) | O(w) | обойти дерево по уровням |
| DFS-обходы | O(n) | O(h) | рекурсивный pre/in/post-order обход |

В худшем случае очередь BFS хранит целый уровень. Для полного бинарного дерева это O(n). Рекурсивные DFS-обходы расходуют память стека пропорционально высоте.

### Сборка и тестирование

```bash
go test -v ./trees/bfs
```

---

## English

This package contains a plain binary tree and demonstrates **breadth-first search** (BFS). BFS visits the tree level by level: the root first, then all of its children, then its grandchildren, and so on. It therefore uses a FIFO queue: the first discovered node is processed first.

For this tree:

```text
        1
      /   \
     2     3
    / \   / \
   4   5 6   7
```

the BFS order is `1, 2, 3, 4, 5, 6, 7`.

### Internal layout

`Node` stores an integer value and two links, `Left` and `Right`. `New` passes values to `Insert` one by one. Insertion runs BFS and puts the new node into the first free position from left to right. This keeps the tree **complete**: every level except possibly the last is full, and the final level is filled from the left.

`Search` also runs BFS and returns the first matching node. With duplicate values, this means the shallowest match, and the leftmost one among matches on the same level.

`Delete` finds the target, copies the value of the deepest rightmost node into it, and removes that final node. This preserves the shape of a complete binary tree, but it does not maintain binary-search-tree ordering: there is no `left < root < right` rule here.

The package also exposes the classic depth-first traversals:

- `TraversePreOrder`: root → left subtree → right subtree;
- `TraverseInOrder`: left subtree → root → right subtree;
- `TraversePostOrder`: left subtree → right subtree → root;
- `TraverseBFS`: level by level, left to right.

### Why BFS is useful

BFS fits business rules based on distance or level: finding the nearest employee in an organization hierarchy, processing category levels, finding the fewest hops in an unweighted graph, locating the nearest UI element, or filling a tree without gaps. DFS is often a better fit for deep exploration under a tighter memory budget.

### Usage

```go
root := New([]int{1, 2, 3, 4, 5, 6, 7})

found := root.Search(6) // node containing 6

var values []int
root.TraverseBFS(func(v int) {
    values = append(values, v)
}) // [1 2 3 4 5 6 7]

root = Delete(root, 3)
```

### Operations and complexity

Let `n` be the number of nodes, `w` the maximum tree width, and `h` its height.

| Operation | Time | Extra space | Description |
|---|---:|---:|---|
| `New(values)` | O(n²) | O(w) | the current implementation performs one BFS insertion per value |
| `Insert(root, value)` | O(n) | O(w) | find the first free level-order position |
| `Search(value)` | O(n) | O(w) | find the first match in BFS order |
| `Delete(root, value)` | O(n) | O(w) | locate the target and final node |
| `TraverseBFS(visit)` | O(n) | O(w) | visit the tree level by level |
| DFS traversals | O(n) | O(h) | recursive pre/in/post-order traversal |

At worst, the BFS queue holds an entire level. For a complete binary tree that is O(n). Recursive DFS traversals use call-stack space proportional to the height.

### Build and test

```bash
go test -v ./trees/bfs
```

---

<br>

> BFS не спешит в глубину: сначала он выслушивает всех, кто рядом. Поэтому первым найденным ответом часто оказывается ближайший.
>
> *BFS does not rush into the depths: it first listens to everyone nearby. That is why the first answer it finds is often the closest one.*
