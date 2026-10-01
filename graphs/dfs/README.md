# Graph Depth-First Search (DFS)

<p style="text-align: left">
  <a href="#русский">Русский</a> ・ <a href="#english">English</a>
</p>

---

## Русский

DFS (поиск в глубину) выбирает одно ребро и продолжает идти по ветке, пока это возможно. Затем алгоритм возвращается к последней развилке и исследует следующую ветку. Пакет содержит итеративную реализацию через стек и классическую рекурсивную реализацию.

### Внутреннее устройство

`Graph` хранит ориентированный или неориентированный граф в списках смежности. Итеративный `Traverse` использует стек LIFO. Соседи кладутся в обратном порядке, чтобы фактический обход соответствовал порядку добавления рёбер. `visited` защищает от бесконечного обхода циклов.

`TraverseRecursive` выражает алгоритм короче, но использует стек вызовов Go. Для очень длинной цепочки безопаснее итеративный вариант. `HasPath` завершает работу сразу после обнаружения цели.

DFS не гарантирует кратчайший путь: результат зависит от порядка соседей.

### Зачем нужен DFS

DFS подходит для проверки достижимости, обхода зависимостей, поиска циклов, топологической сортировки, компонент связности, лабиринтов и backtracking. Он особенно удобен, когда нужно полностью исследовать одну ветку перед переходом к следующей.

### Использование

```go
g, _ := New(5, true) // ориентированный граф
g.AddEdge(0, 1)
g.AddEdge(0, 2)
g.AddEdge(1, 3)
g.AddEdge(3, 4)

order, _ := g.Traverse(0)
recursive, _ := g.TraverseRecursive(0)
reachable, _ := g.HasPath(0, 4) // true
```

### Операции и сложность

Пусть `V` — число вершин, `E` — число рёбер.

| Операция | Время | Доп. память |
|---|---:|---:|
| `New(V, directed)` | O(V) | O(V) |
| `AddEdge(a, b)` | O(1) амортизированно | O(1) |
| `Traverse(start)` | O(V + E) | O(V) |
| `TraverseRecursive(start)` | O(V + E) | O(V) |
| `HasPath(start, target)` | O(V + E) | O(V) |

### Тестирование

```bash
go test -v ./graphs/dfs
```

Тесты включают циклы, ориентированную достижимость и итеративный обход цепочки из 100 000 вершин.

---

## English

DFS (depth-first search) chooses an edge and follows that branch as far as possible. It then backtracks to the latest fork and explores the next branch. The package provides both an iterative stack-based implementation and the classic recursive form.

### Internal layout

`Graph` stores a directed or undirected graph as adjacency lists. Iterative `Traverse` uses a LIFO stack. Neighbours are pushed in reverse so the actual visit order follows edge insertion order. A `visited` array prevents cycles from causing infinite traversal.

`TraverseRecursive` expresses the algorithm more compactly but consumes the Go call stack. The iterative form is safer for a very long chain. `HasPath` stops as soon as it discovers the target.

DFS does not guarantee a shortest path: its result depends on neighbour order.

### Why DFS is useful

DFS fits reachability checks, dependency traversal, cycle detection, topological sorting, connected components, mazes, and backtracking. It is especially useful when one branch should be explored completely before moving to the next.

### Usage

```go
g, _ := New(5, true) // directed graph
g.AddEdge(0, 1)
g.AddEdge(0, 2)
g.AddEdge(1, 3)
g.AddEdge(3, 4)

order, _ := g.Traverse(0)
recursive, _ := g.TraverseRecursive(0)
reachable, _ := g.HasPath(0, 4) // true
```

### Operations and complexity

Let `V` be the number of vertices and `E` the number of edges.

| Operation | Time | Extra space |
|---|---:|---:|
| `New(V, directed)` | O(V) | O(V) |
| `AddEdge(a, b)` | O(1) amortized | O(1) |
| `Traverse(start)` | O(V + E) | O(V) |
| `TraverseRecursive(start)` | O(V + E) | O(V) |
| `HasPath(start, target)` | O(V + E) | O(V) |

### Testing

```bash
go test -v ./graphs/dfs
```

Tests include cycles, directed reachability, and iterative traversal of a 100,000-vertex chain.

---

> DFS не обещает самый короткий путь. Он обещает не оставить выбранную ветку неисследованной.
>
> *DFS does not promise the shortest path. It promises not to leave the chosen branch unexplored.*
