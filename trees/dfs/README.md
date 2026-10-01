# Depth-First Search (DFS) Binary Tree

<p style="text-align: left">
  <a href="#русский">Русский</a> ・ <a href="#english">English</a>
</p>

---

## Русский

Этот пакет демонстрирует **поиск в глубину** (Depth-First Search, DFS) на бинарном дереве. DFS выбирает ветку и идёт по ней настолько глубоко, насколько возможно. Когда продолжать некуда, алгоритм возвращается к ближайшей развилке и исследует следующую ветку.

Для дерева

```text
        1
      /   \
     2     3
    / \   / \
   4   5 6   7
```

основной обход `TraverseDFS` использует pre-order и выдаёт `1, 2, 4, 5, 3, 6, 7`.

### Внутреннее устройство

`Node` хранит `Value` и ссылки `Left`/`Right`. `New` принимает значения в порядке уровней и за O(n) строит полное бинарное дерево: для элемента с индексом `i` дети находятся по индексам `2*i+1` и `2*i+2`.

`Search` и `TraverseDFS` реализованы итеративно через стек LIFO. Сначала в стек кладётся правый ребёнок, затем левый: левый оказывается на вершине и обрабатывается первым. Итеративный вариант не расходует стек вызовов Go и поэтому подходит для очень глубоких деревьев лучше рекурсии.

Пакет также содержит три классических рекурсивных порядка:

- **pre-order** (`root → left → right`) — копирование дерева, сериализация, выполнение выражения в префиксной форме;
- **in-order** (`left → root → right`) — получение отсортированных значений у бинарного дерева поиска;
- **post-order** (`left → right → root`) — удаление дерева, подсчёт размера каталогов, вычисление выражений снизу вверх.

`New` создаёт обычное полное бинарное дерево, а не бинарное дерево поиска. Поэтому `Search` не может отбросить половину дерева по сравнению значений и в худшем случае проверяет каждый узел.

### Зачем нужен DFS

DFS применяют там, где нужно полностью исследовать одну ветку перед переходом к следующей: обход файлов и DOM, анализ зависимостей, поиск циклов, топологическая сортировка, backtracking-задачи, проверка существования пути и обработка вложенных категорий. В отличие от BFS, первый найденный DFS путь не обязан быть кратчайшим.

### Использование

```go
root := New([]int{1, 2, 3, 4, 5, 6, 7})

found := root.Search(6) // узел со значением 6

var values []int
root.TraverseDFS(func(v int) {
    values = append(values, v)
}) // [1 2 4 5 3 6 7]
```

Дерево можно собрать вручную, если нужна конкретная форма:

```go
root := &Node{
    Value: 10,
    Left:  &Node{Value: 5},
    Right: &Node{Value: 20},
}
```

### Операции и сложность

Пусть `n` — число узлов, `h` — высота дерева.

| Операция | Время | Доп. память | Описание |
|---|---:|---:|---|
| `New(values)` | O(n) | O(n) | создать узлы и связать их по индексам |
| `Search(value)` | O(n) | O(h) в среднем, O(n) в худшем | найти первое совпадение в pre-order DFS |
| `TraverseDFS(visit)` | O(n) | O(h) в среднем, O(n) в худшем | итеративный pre-order обход |
| `TraversePreOrder` | O(n) | O(h) | рекурсивный обход root-left-right |
| `TraverseInOrder` | O(n) | O(h) | рекурсивный обход left-root-right |
| `TraversePostOrder` | O(n) | O(h) | рекурсивный обход left-right-root |

Для сбалансированного дерева `h = O(log n)`. Для дерева-цепочки `h = O(n)`. У итеративного DFS фактический размер стека зависит от числа ещё не исследованных развилок: на цепочке он может оставаться O(1), но общая верхняя граница — O(n).

### Сборка и тестирование

```bash
go test -v ./trees/dfs
```

Тесты проверяют построение дерева, все порядки обхода, отсутствующие и повторяющиеся значения, nil-получатели и итеративный обход цепочки из 100 000 узлов.

---

## English

This package demonstrates **depth-first search** (DFS) on a binary tree. DFS selects one branch and follows it as deeply as possible. Once it cannot continue, it backtracks to the nearest fork and explores the next branch.

For this tree:

```text
        1
      /   \
     2     3
    / \   / \
   4   5 6   7
```

the primary `TraverseDFS` traversal uses pre-order and produces `1, 2, 4, 5, 3, 6, 7`.

### Internal layout

`Node` stores a `Value` and `Left`/`Right` links. `New` accepts values in level order and builds a complete binary tree in O(n): the children of item `i` are located at indices `2*i+1` and `2*i+2`.

`Search` and `TraverseDFS` are iterative and use a LIFO stack. The right child is pushed before the left one, leaving the left child on top to be processed first. The iterative form does not consume the Go call stack, so it is safer than recursion for very deep trees.

The package also provides three classic recursive orders:

- **pre-order** (`root → left → right`) — copying a tree, serialization, prefix-expression evaluation;
- **in-order** (`left → root → right`) — obtaining sorted values from a binary search tree;
- **post-order** (`left → right → root`) — deleting a tree, computing directory sizes, bottom-up expression evaluation.

`New` creates a plain complete binary tree, not a binary search tree. Therefore `Search` cannot discard half of the tree by comparing values and may inspect every node in the worst case.

### Why DFS is useful

DFS is useful when one branch must be fully explored before moving to the next: traversing file trees and the DOM, dependency analysis, cycle detection, topological sorting, backtracking, reachability checks, and processing nested categories. Unlike BFS, the first path found by DFS is not necessarily the shortest one.

### Usage

```go
root := New([]int{1, 2, 3, 4, 5, 6, 7})

found := root.Search(6) // node containing 6

var values []int
root.TraverseDFS(func(v int) {
    values = append(values, v)
}) // [1 2 4 5 3 6 7]
```

Build the tree manually when a specific shape is required:

```go
root := &Node{
    Value: 10,
    Left:  &Node{Value: 5},
    Right: &Node{Value: 20},
}
```

### Operations and complexity

Let `n` be the number of nodes and `h` their tree height.

| Operation | Time | Extra space | Description |
|---|---:|---:|---|
| `New(values)` | O(n) | O(n) | allocate nodes and link them by index |
| `Search(value)` | O(n) | O(h) average, O(n) worst case | find the first pre-order DFS match |
| `TraverseDFS(visit)` | O(n) | O(h) average, O(n) worst case | iterative pre-order traversal |
| `TraversePreOrder` | O(n) | O(h) | recursive root-left-right traversal |
| `TraverseInOrder` | O(n) | O(h) | recursive left-root-right traversal |
| `TraversePostOrder` | O(n) | O(h) | recursive left-right-root traversal |

For a balanced tree, `h = O(log n)`. For a chain-shaped tree, `h = O(n)`. The actual iterative DFS stack size depends on the number of unexplored forks: it can stay O(1) on a chain, while its general upper bound is O(n).

### Build and test

```bash
go test -v ./trees/dfs
```

Tests cover tree construction, every traversal order, missing and duplicate values, nil receivers, and iterative traversal of a 100,000-node chain.

---

<br>

> DFS выбирает путь и следует ему до конца. Глубина даёт ответ быстрее, когда нужная ветка выбрана верно, но не обещает, что найденный путь окажется кратчайшим.
>
> *DFS chooses a path and follows it to the end. Depth finds an answer quickly when the right branch is chosen, but it never promises that the discovered path is the shortest.*
