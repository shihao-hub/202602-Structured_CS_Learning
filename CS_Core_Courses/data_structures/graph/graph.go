package graph

import (
	"fmt"
	"math"
)

// Graph 图结构（邻接表表示）
type Graph struct {
	vertices int            // 顶点数
	edges    int            // 边数
	directed bool           // 是否有向
	adjList  map[int][]Edge // 邻接表
}

// Edge 边结构
type Edge struct {
	To     int     // 目标顶点
	Weight float64 // 权重
}

// NewGraph 创建图
func NewGraph(vertices int, directed bool) *Graph {
	return &Graph{
		vertices: vertices,
		edges:    0,
		directed: directed,
		adjList:  make(map[int][]Edge),
	}
}

// AddEdge 添加边
func (g *Graph) AddEdge(from, to int, weight float64) {
	g.adjList[from] = append(g.adjList[from], Edge{To: to, Weight: weight})
	g.edges++

	if !g.directed {
		g.adjList[to] = append(g.adjList[to], Edge{To: from, Weight: weight})
	}
}

// AddUnweightedEdge 添加无权边
func (g *Graph) AddUnweightedEdge(from, to int) {
	g.AddEdge(from, to, 1.0)
}

// GetNeighbors 获取邻居节点
func (g *Graph) GetNeighbors(vertex int) []Edge {
	return g.adjList[vertex]
}

// GetVertices 获取顶点数
func (g *Graph) GetVertices() int {
	return g.vertices
}

// GetEdges 获取边数
func (g *Graph) GetEdges() int {
	return g.edges
}

// BFS 广度优先搜索
func (g *Graph) BFS(start int) []int {
	visited := make(map[int]bool)
	result := make([]int, 0)
	queue := []int{start}
	visited[start] = true

	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		result = append(result, vertex)

		for _, edge := range g.adjList[vertex] {
			if !visited[edge.To] {
				visited[edge.To] = true
				queue = append(queue, edge.To)
			}
		}
	}

	return result
}

// DFS 深度优先搜索
func (g *Graph) DFS(start int) []int {
	visited := make(map[int]bool)
	result := make([]int, 0)
	g.dfsHelper(start, visited, &result)
	return result
}

func (g *Graph) dfsHelper(vertex int, visited map[int]bool, result *[]int) {
	visited[vertex] = true
	*result = append(*result, vertex)

	for _, edge := range g.adjList[vertex] {
		if !visited[edge.To] {
			g.dfsHelper(edge.To, visited, result)
		}
	}
}

// Dijkstra 最短路径算法
func (g *Graph) Dijkstra(start int) (map[int]float64, map[int]int) {
	dist := make(map[int]float64)
	prev := make(map[int]int)
	visited := make(map[int]bool)

	// 初始化距离
	for i := 0; i < g.vertices; i++ {
		dist[i] = math.Inf(1)
		prev[i] = -1
	}
	dist[start] = 0

	for i := 0; i < g.vertices; i++ {
		// 找到未访问的最小距离顶点
		u := -1
		minDist := math.Inf(1)
		for v := 0; v < g.vertices; v++ {
			if !visited[v] && dist[v] < minDist {
				u = v
				minDist = dist[v]
			}
		}

		if u == -1 {
			break
		}

		visited[u] = true

		// 更新邻居距离
		for _, edge := range g.adjList[u] {
			newDist := dist[u] + edge.Weight
			if newDist < dist[edge.To] {
				dist[edge.To] = newDist
				prev[edge.To] = u
			}
		}
	}

	return dist, prev
}

// GetShortestPath 获取从起点到终点的最短路径
func (g *Graph) GetShortestPath(start, end int) ([]int, float64) {
	dist, prev := g.Dijkstra(start)

	if math.IsInf(dist[end], 1) {
		return nil, math.Inf(1)
	}

	path := make([]int, 0)
	for v := end; v != -1; v = prev[v] {
		path = append([]int{v}, path...)
	}

	return path, dist[end]
}

// HasCycle 检测是否有环（使用DFS）
func (g *Graph) HasCycle() bool {
	if g.directed {
		return g.hasCycleDirected()
	}
	return g.hasCycleUndirected()
}

func (g *Graph) hasCycleDirected() bool {
	white := make(map[int]bool) // 未访问
	gray := make(map[int]bool)  // 正在访问
	black := make(map[int]bool) // 已访问

	for i := 0; i < g.vertices; i++ {
		white[i] = true
	}

	for i := 0; i < g.vertices; i++ {
		if white[i] {
			if g.hasCycleDirectedDFS(i, white, gray, black) {
				return true
			}
		}
	}
	return false
}

func (g *Graph) hasCycleDirectedDFS(v int, white, gray, black map[int]bool) bool {
	white[v] = false
	gray[v] = true

	for _, edge := range g.adjList[v] {
		if gray[edge.To] {
			return true
		}
		if white[edge.To] && g.hasCycleDirectedDFS(edge.To, white, gray, black) {
			return true
		}
	}

	gray[v] = false
	black[v] = true
	return false
}

func (g *Graph) hasCycleUndirected() bool {
	visited := make(map[int]bool)

	for i := 0; i < g.vertices; i++ {
		if !visited[i] {
			if g.hasCycleUndirectedDFS(i, visited, -1) {
				return true
			}
		}
	}
	return false
}

func (g *Graph) hasCycleUndirectedDFS(v int, visited map[int]bool, parent int) bool {
	visited[v] = true

	for _, edge := range g.adjList[v] {
		if !visited[edge.To] {
			if g.hasCycleUndirectedDFS(edge.To, visited, v) {
				return true
			}
		} else if edge.To != parent {
			return true
		}
	}
	return false
}

// TopologicalSort 拓扑排序（仅对有向无环图）
func (g *Graph) TopologicalSort() ([]int, bool) {
	if !g.directed {
		return nil, false
	}

	if g.HasCycle() {
		return nil, false
	}

	visited := make(map[int]bool)
	stack := make([]int, 0)

	for i := 0; i < g.vertices; i++ {
		if !visited[i] {
			g.topologicalSortDFS(i, visited, &stack)
		}
	}

	// 反转栈
	result := make([]int, len(stack))
	for i := 0; i < len(stack); i++ {
		result[i] = stack[len(stack)-1-i]
	}

	return result, true
}

func (g *Graph) topologicalSortDFS(v int, visited map[int]bool, stack *[]int) {
	visited[v] = true

	for _, edge := range g.adjList[v] {
		if !visited[edge.To] {
			g.topologicalSortDFS(edge.To, visited, stack)
		}
	}

	*stack = append(*stack, v)
}

// Print 打印图
func (g *Graph) Print() {
	graphType := "无向图"
	if g.directed {
		graphType = "有向图"
	}
	fmt.Printf("%s (V=%d, E=%d):\n", graphType, g.vertices, g.edges)
	for v := 0; v < g.vertices; v++ {
		if edges, ok := g.adjList[v]; ok && len(edges) > 0 {
			fmt.Printf("  %d -> ", v)
			for i, edge := range edges {
				if i > 0 {
					fmt.Print(", ")
				}
				if edge.Weight != 1.0 {
					fmt.Printf("%d(%.1f)", edge.To, edge.Weight)
				} else {
					fmt.Printf("%d", edge.To)
				}
			}
			fmt.Println()
		}
	}
}

// GraphExample 图示例
func GraphExample() {
	fmt.Println("=== 图 (Graph) 示例 ===")

	// 无向图
	fmt.Println("\n1. 无向图:")
	undirected := NewGraph(6, false)
	undirected.AddUnweightedEdge(0, 1)
	undirected.AddUnweightedEdge(0, 2)
	undirected.AddUnweightedEdge(1, 2)
	undirected.AddUnweightedEdge(1, 3)
	undirected.AddUnweightedEdge(2, 4)
	undirected.AddUnweightedEdge(3, 4)
	undirected.AddUnweightedEdge(3, 5)
	undirected.AddUnweightedEdge(4, 5)
	undirected.Print()

	// 图遍历
	fmt.Println("\n2. 图遍历:")
	fmt.Printf("BFS (从顶点 0 开始): %v\n", undirected.BFS(0))
	fmt.Printf("DFS (从顶点 0 开始): %v\n", undirected.DFS(0))

	// 环检测
	fmt.Println("\n3. 环检测:")
	fmt.Printf("无向图是否有环: %t\n", undirected.HasCycle())

	// 加权有向图
	fmt.Println("\n4. 加权有向图和最短路径:")
	weighted := NewGraph(5, true)
	weighted.AddEdge(0, 1, 4)
	weighted.AddEdge(0, 2, 2)
	weighted.AddEdge(1, 2, 3)
	weighted.AddEdge(1, 3, 2)
	weighted.AddEdge(1, 4, 3)
	weighted.AddEdge(2, 1, 1)
	weighted.AddEdge(2, 3, 4)
	weighted.AddEdge(2, 4, 5)
	weighted.AddEdge(4, 3, 1)
	weighted.Print()

	// 最短路径
	fmt.Println("\nDijkstra 最短路径 (从顶点 0):")
	dist, _ := weighted.Dijkstra(0)
	for v := 0; v < weighted.vertices; v++ {
		if math.IsInf(dist[v], 1) {
			fmt.Printf("  到顶点 %d: 不可达\n", v)
		} else {
			fmt.Printf("  到顶点 %d: %.1f\n", v, dist[v])
		}
	}

	// 具体路径
	path, cost := weighted.GetShortestPath(0, 3)
	fmt.Printf("\n从顶点 0 到顶点 3 的最短路径: %v, 距离: %.1f\n", path, cost)

	// 拓扑排序
	fmt.Println("\n5. 拓扑排序:")
	dag := NewGraph(6, true)
	dag.AddUnweightedEdge(5, 2)
	dag.AddUnweightedEdge(5, 0)
	dag.AddUnweightedEdge(4, 0)
	dag.AddUnweightedEdge(4, 1)
	dag.AddUnweightedEdge(2, 3)
	dag.AddUnweightedEdge(3, 1)
	dag.Print()

	if order, ok := dag.TopologicalSort(); ok {
		fmt.Printf("拓扑排序结果: %v\n", order)
	} else {
		fmt.Println("无法进行拓扑排序（图中有环或不是有向图）")
	}

	fmt.Println()
}
