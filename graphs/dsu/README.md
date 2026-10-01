# Disjoint-Set Union (DSU)

<p style="text-align: left">
  <a href="#русский">Русский</a> ・ <a href="#english">English</a>
</p>

---

## Русский

DSU (Disjoint-Set Union), также известная как Union-Find, хранит разбиение элементов на непересекающиеся множества. Структура быстро отвечает на два вопроса: к какому множеству относится элемент и находятся ли два элемента в одной группе.

### Внутреннее устройство

Каждое множество представлено деревом. `parent[x]` указывает на родителя элемента, а корень указывает сам на себя и служит представителем множества.

`Find` поднимается к корню и выполняет **сжатие пути**: все пройденные элементы начинают ссылаться прямо на корень. `Union` использует **объединение по размеру**: корень меньшего дерева присоединяется к корню большего. Вместе эти оптимизации не дают деревьям становиться высокими.

`components` хранит текущее число независимых множеств, а `size` у корня — размер компоненты.

### Зачем нужен DSU

DSU применяют для динамической проверки связности: объединения пользователей или устройств в группы, кластеризации, обнаружения цикла в неориентированном графе, алгоритма Краскала, объединения аккаунтов и обработки островов на карте.

Структура хорошо поддерживает добавление связей, но не удаление: после `Union` разделить компоненту обратно обычная DSU не умеет.

### Использование

```go
d, _ := New(5) // {0} {1} {2} {3} {4}

d.Union(0, 1)
d.Union(1, 2)

connected, _ := d.Connected(0, 2) // true
size, _ := d.Size(1)              // 3
count := d.Components()           // 3
root, _ := d.Find(2)
```

### Операции и сложность

Пусть `n` — число элементов, `α(n)` — обратная функция Аккермана, растущая настолько медленно, что на практических размерах не превышает небольшой константы.

| Операция | Амортизированное время | Память |
|---|---:|---:|
| `New(n)` | O(n) | O(n) |
| `Find(x)` | O(α(n)) | O(1) |
| `Union(a, b)` | O(α(n)) | O(1) |
| `Connected(a, b)` | O(α(n)) | O(1) |
| `Size(x)` | O(α(n)) | O(1) |
| `Components()` | O(1) | O(1) |

### Тестирование

```bash
go test -v ./graphs/dsu
```

---

## English

DSU (Disjoint-Set Union), also known as Union-Find, maintains a partition of elements into non-overlapping sets. It quickly answers two questions: which set contains an element, and whether two elements belong to the same group.

### Internal layout

Each set is represented by a tree. `parent[x]` points to an element's parent; a root points to itself and acts as the set representative.

`Find` climbs to the root and performs **path compression**, making every visited element point directly to the root. `Union` applies **union by size**, attaching the smaller tree's root below the larger tree's root. Together, these optimizations keep trees shallow.

`components` stores the current number of independent sets, while a root's `size` stores its component size.

### Why DSU is useful

DSU supports dynamic connectivity: grouping users or devices, clustering, detecting cycles in an undirected graph, Kruskal's algorithm, account merging, and island processing on a grid.

It handles adding connections well, but not removing them: a conventional DSU cannot split a component after `Union`.

### Usage

```go
d, _ := New(5) // {0} {1} {2} {3} {4}

d.Union(0, 1)
d.Union(1, 2)

connected, _ := d.Connected(0, 2) // true
size, _ := d.Size(1)              // 3
count := d.Components()           // 3
root, _ := d.Find(2)
```

### Operations and complexity

Let `n` be the element count and `α(n)` the inverse Ackermann function, which grows so slowly that it stays below a small constant at practical sizes.

| Operation | Amortized time | Space |
|---|---:|---:|
| `New(n)` | O(n) | O(n) |
| `Find(x)` | O(α(n)) | O(1) |
| `Union(a, b)` | O(α(n)) | O(1) |
| `Connected(a, b)` | O(α(n)) | O(1) |
| `Size(x)` | O(α(n)) | O(1) |
| `Components()` | O(1) | O(1) |

### Testing

```bash
go test -v ./graphs/dsu
```

---

> DSU не хранит все связи. Она хранит достаточно, чтобы знать, кто уже принадлежит одному целому.
>
> *DSU does not store every connection. It stores just enough to know what already belongs to one whole.*
