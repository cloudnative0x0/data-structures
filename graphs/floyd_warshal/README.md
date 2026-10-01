# Floyd–Warshall Algorithm

<p style="text-align: left">
  <a href="#русский">Русский</a> ・ <a href="#english">English</a>
</p>

---

## Русский

Алгоритм Флойда–Уоршелла находит кратчайшие расстояния **между всеми парами вершин** взвешенного графа. Он поддерживает отрицательные рёбра, но кратчайшие пути не определены при наличии отрицательного цикла.

### Внутреннее устройство

Граф хранится в матрице весов. `weights[i][j]` равен стоимости прямого ребра `i → j`, `0` на диагонали и `Infinity`, если ребра нет. Повторное добавление ребра сохраняет минимальный вес.

Основная динамика последовательно разрешает каждую вершину `k` как промежуточную:

```text
dist[i][j] = min(dist[i][j], dist[i][k] + dist[k][j])
```

После шага `k` матрица содержит лучшие пути, внутренние вершины которых взяты только из `0..k`. Матрица `next` запоминает первый переход и позволяет восстановить маршрут, а отрицательное значение `dist[v][v]` обнаруживает отрицательный цикл.

### Зачем нужен Floyd–Warshall

Алгоритм полезен, когда граф относительно небольшой или плотный и ожидается много запросов между произвольными парами: маршруты между офисами, таблицы стоимости между сервисами, анализ сетевой достижимости, игровые карты и транзитивное замыкание.

Для одного источника на большом разреженном графе обычно выгоднее Дейкстра или Bellman–Ford. Матрицы Floyd–Warshall требуют O(V²) памяти.

### Использование

```go
g, _ := New(4)
g.AddEdge(0, 1, 4)
g.AddEdge(0, 2, 10)
g.AddEdge(1, 2, -2)
g.AddEdge(2, 3, 3)

result, err := g.ShortestPaths()
distance, found, _ := result.Distance(0, 3) // 5, true
path, found, _ := result.Path(0, 3)         // [0 1 2 3], true
```

`AddEdge` создаёт ориентированное ребро, `AddUndirectedEdge` — два симметричных. Недостижимый путь возвращает `found == false`. При отрицательном цикле `ShortestPaths` возвращает `ErrNegativeCycle`.

### Операции и сложность

Пусть `V` — число вершин.

| Операция | Время | Доп. память |
|---|---:|---:|
| `New(V)` | O(V²) | O(V²) |
| `AddEdge(a, b, weight)` | O(1) | O(1) |
| `ShortestPaths()` | O(V³) | O(V²) |
| `Distance(a, b)` | O(1) | O(1) |
| `Path(a, b)` | O(L) | O(L) |

Здесь `L` — число вершин в восстановленном пути.

### Тестирование

```bash
go test -v ./graphs/floyd_warshal
```

---

## English

The Floyd–Warshall algorithm finds shortest distances **between every pair of vertices** in a weighted graph. It supports negative edges, but shortest paths are undefined when a negative cycle exists.

### Internal layout

The graph is stored as a weight matrix. `weights[i][j]` is the cost of a direct `i → j` edge, `0` on the diagonal, and `Infinity` when no edge exists. Adding an edge repeatedly retains its smallest weight.

The dynamic program successively allows every vertex `k` as an intermediate point:

```text
dist[i][j] = min(dist[i][j], dist[i][k] + dist[k][j])
```

After step `k`, the matrix contains the best paths whose internal vertices come only from `0..k`. A `next` matrix records the first hop for path reconstruction, while a negative `dist[v][v]` reveals a negative cycle.

### Why Floyd–Warshall is useful

The algorithm is useful when a graph is relatively small or dense and many arbitrary-pair queries are expected: routes between offices, service-to-service cost tables, network reachability, game maps, and transitive closure.

For one source in a large sparse graph, Dijkstra or Bellman–Ford is usually more efficient. Floyd–Warshall matrices require O(V²) space.

### Usage

```go
g, _ := New(4)
g.AddEdge(0, 1, 4)
g.AddEdge(0, 2, 10)
g.AddEdge(1, 2, -2)
g.AddEdge(2, 3, 3)

result, err := g.ShortestPaths()
distance, found, _ := result.Distance(0, 3) // 5, true
path, found, _ := result.Path(0, 3)         // [0 1 2 3], true
```

`AddEdge` creates a directed edge; `AddUndirectedEdge` creates both directions. An unreachable route returns `found == false`. When a negative cycle exists, `ShortestPaths` returns `ErrNegativeCycle`.

### Operations and complexity

Let `V` be the number of vertices.

| Operation | Time | Extra space |
|---|---:|---:|
| `New(V)` | O(V²) | O(V²) |
| `AddEdge(a, b, weight)` | O(1) | O(1) |
| `ShortestPaths()` | O(V³) | O(V²) |
| `Distance(a, b)` | O(1) | O(1) |
| `Path(a, b)` | O(L) | O(L) |

Here `L` is the number of vertices in the reconstructed path.

### Testing

```bash
go test -v ./graphs/floyd_warshal
```

---

> Флойд–Уоршелл рассматривает каждую вершину как возможного посредника и постепенно превращает локальные рёбра в глобальную карту расстояний.
>
> *Floyd–Warshall treats every vertex as a possible intermediary, gradually turning local edges into a global map of distances.*
