# Graph Breadth-First Search (BFS)

<p style="text-align: left">
  <a href="#русский">Русский</a> ・ <a href="#english">English</a>
</p>

---

## Русский

BFS (поиск в ширину) обходит невзвешенный граф слоями: сначала стартовую вершину, затем соседей на расстоянии одного ребра, потом вершины на расстоянии двух рёбер и так далее. Очередь FIFO сохраняет этот порядок, поэтому первое достижение вершины одновременно задаёт путь с минимальным числом рёбер.

### Внутреннее устройство

`Graph` хранит списки смежности `adjacency`: `adjacency[v]` содержит вершины, в которые ведут рёбра из `v`. В ориентированном графе `AddEdge(a, b)` добавляет только `a → b`, в неориентированном — также `b → a`.

Массив `visited` не позволяет повторно добавлять вершины в очередь при циклах. В `Distances` значение `-1` означает недостижимую вершину. `ShortestPath` дополнительно запоминает родителя каждой обнаруженной вершины и восстанавливает маршрут от цели к старту.

### Зачем нужен BFS

BFS применяют для поиска кратчайшего пути в невзвешенной сети: минимального числа пересадок, социальных связей, переходов между состояниями, кликов в интерфейсе или шагов по клеточной карте. Для взвешенных рёбер нужен другой алгоритм, например Дейкстра или Floyd–Warshall.

### Использование

```go
g, _ := New(5, false) // неориентированный граф
g.AddEdge(0, 1)
g.AddEdge(0, 2)
g.AddEdge(1, 3)
g.AddEdge(2, 4)

order, _ := g.Traverse(0)             // [0 1 2 3 4]
distance, _ := g.Distances(0)          // [0 1 1 2 2]
path, found, _ := g.ShortestPath(0, 4) // [0 2 4], true
```

### Операции и сложность

Пусть `V` — число вершин, `E` — число рёбер.

| Операция | Время | Доп. память |
|---|---:|---:|
| `New(V, directed)` | O(V) | O(V) |
| `AddEdge(a, b)` | O(1) амортизированно | O(1) |
| `Traverse(start)` | O(V + E) | O(V) |
| `Distances(start)` | O(V + E) | O(V) |
| `ShortestPath(start, target)` | O(V + E) | O(V) |

### Тестирование

```bash
go test -v ./graphs/bfs
```

---

## English

BFS (breadth-first search) explores an unweighted graph in layers: the start vertex first, then vertices one edge away, then two edges away, and so on. A FIFO queue preserves this order, so the first discovery of a vertex also gives a path with the minimum number of edges.

### Internal layout

`Graph` stores adjacency lists: `adjacency[v]` contains vertices reached by outgoing edges from `v`. In a directed graph, `AddEdge(a, b)` adds only `a → b`; in an undirected graph it also adds `b → a`.

A `visited` array prevents cycles from placing vertices into the queue repeatedly. In `Distances`, `-1` marks an unreachable vertex. `ShortestPath` also remembers each discovered vertex's parent and reconstructs the route from target back to start.

### Why BFS is useful

BFS finds shortest paths in unweighted networks: minimum transfers, social connections, state transitions, UI clicks, or moves on a grid. Weighted edges require a different algorithm, such as Dijkstra or Floyd–Warshall.

### Usage

```go
g, _ := New(5, false) // undirected graph
g.AddEdge(0, 1)
g.AddEdge(0, 2)
g.AddEdge(1, 3)
g.AddEdge(2, 4)

order, _ := g.Traverse(0)             // [0 1 2 3 4]
distance, _ := g.Distances(0)          // [0 1 1 2 2]
path, found, _ := g.ShortestPath(0, 4) // [0 2 4], true
```

### Operations and complexity

Let `V` be the number of vertices and `E` the number of edges.

| Operation | Time | Extra space |
|---|---:|---:|
| `New(V, directed)` | O(V) | O(V) |
| `AddEdge(a, b)` | O(1) amortized | O(1) |
| `Traverse(start)` | O(V + E) | O(V) |
| `Distances(start)` | O(V + E) | O(V) |
| `ShortestPath(start, target)` | O(V + E) | O(V) |

### Testing

```bash
go test -v ./graphs/bfs
```

---

> BFS сначала проверяет ближайшее. В невзвешенном графе эта осторожность превращается в гарантию кратчайшего пути.
>
> *BFS checks what is nearest first. In an unweighted graph, that caution becomes a shortest-path guarantee.*
