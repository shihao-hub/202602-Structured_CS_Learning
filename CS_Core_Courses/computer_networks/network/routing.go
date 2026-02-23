package network

import (
	"fmt"
	"math"
	"strings"
)

// RoutingTable 路由表
type RoutingTable struct {
	Destination string  // 目的网络
	NextHop     string  // 下一跳
	Cost        float64 // 距离/代价
}

// String 格式化输出路由表项
func (r *RoutingTable) String() string {
	return fmt.Sprintf("目的: %-10s | 下一跳: %-8s | 代价: %.1f", r.Destination, r.NextHop, r.Cost)
}

// Node 网络节点
type Node struct {
	Name          string                    // 节点名称
	RoutingTable  map[string]*RoutingTable  // 路由表 (目的地 -> 路由项)
	Neighbors     map[string]float64        // 邻居节点及其链路代价
	DistanceTable map[string]map[string]float64 // 距离表 (用于距离向量算法)
}

// NewNode 创建新节点
func NewNode(name string) *Node {
	return &Node{
		Name:          name,
		RoutingTable:  make(map[string]*RoutingTable),
		Neighbors:     make(map[string]float64),
		DistanceTable: make(map[string]map[string]float64),
	}
}

// AddNeighbor 添加邻居节点
func (n *Node) AddNeighbor(neighbor string, cost float64) {
	n.Neighbors[neighbor] = cost
	// 初始化路由表: 直连邻居,下一跳就是自己
	n.RoutingTable[neighbor] = &RoutingTable{
		Destination: neighbor,
		NextHop:     neighbor,
		Cost:        cost,
	}
	// 初始化距离表
	if n.DistanceTable[neighbor] == nil {
		n.DistanceTable[neighbor] = make(map[string]float64)
	}
	n.DistanceTable[neighbor][neighbor] = cost
}

// PrintRoutingTable 打印路由表
func (n *Node) PrintRoutingTable() {
	fmt.Printf("\n【节点 %s 的路由表】\n", n.Name)
	fmt.Println("┌────────────┬──────────┬────────┐")
	fmt.Println("│ 目的网络   │ 下一跳   │  代价  │")
	fmt.Println("├────────────┼──────────┼────────┤")
	for dest, route := range n.RoutingTable {
		if dest == n.Name {
			continue // 跳过自己到自己的路由
		}
		fmt.Printf("│ %-10s │ %-8s │ %6.1f │\n", route.Destination, route.NextHop, route.Cost)
	}
	fmt.Println("└────────────┴──────────┴────────┘")
}

// DistanceVectorRouting 距离向量路由算法 (类 RIP - Bellman-Ford)
// 对应 408 考点: 距离向量算法,路由信息交换
type DistanceVectorRouting struct {
	Nodes map[string]*Node // 所有节点
}

// NewDistanceVectorRouting 创建距离向量路由
func NewDistanceVectorRouting() *DistanceVectorRouting {
	return &DistanceVectorRouting{
		Nodes: make(map[string]*Node),
	}
}

// AddNode 添加节点
func (dv *DistanceVectorRouting) AddNode(name string) {
	dv.Nodes[name] = NewNode(name)
	// 初始化到自己的距离为 0
	dv.Nodes[name].RoutingTable[name] = &RoutingTable{
		Destination: name,
		NextHop:     name,
		Cost:        0,
	}
}

// AddLink 添加链路 (双向)
func (dv *DistanceVectorRouting) AddLink(node1, node2 string, cost float64) {
	dv.Nodes[node1].AddNeighbor(node2, cost)
	dv.Nodes[node2].AddNeighbor(node1, cost)
}

// UpdateRoutingTables 更新路由表 (一轮交换)
func (dv *DistanceVectorRouting) UpdateRoutingTables() bool {
	updated := false

	// 每个节点从邻居接收路由信息
	for _, node := range dv.Nodes {
		// 遍历所有邻居
		for neighbor, linkCost := range node.Neighbors {
			neighborNode := dv.Nodes[neighbor]

			// 遍历邻居的路由表
			for dest, route := range neighborNode.RoutingTable {
				if dest == node.Name {
					continue // 跳过到自己的路由
				}

				// 计算通过该邻居到达目的地的新代价
				newCost := linkCost + route.Cost

				// 如果没有到该目的地的路由,或者新路由代价更小,则更新
				currentRoute, exists := node.RoutingTable[dest]
				if !exists || newCost < currentRoute.Cost {
					node.RoutingTable[dest] = &RoutingTable{
						Destination: dest,
						NextHop:     neighbor,
						Cost:        newCost,
					}
					updated = true
				}
			}
		}
	}

	return updated
}

// Run 运行距离向量算法直到收敛
func (dv *DistanceVectorRouting) Run(maxIterations int) int {
	fmt.Println("\n【距离向量路由算法 - 类 RIP/Bellman-Ford】")
	iterations := 0

	for i := 0; i < maxIterations; i++ {
		iterations++
		fmt.Printf("\n第 %d 轮路由信息交换...\n", iterations)
		updated := dv.UpdateRoutingTables()
		if !updated {
			fmt.Println("✓ 路由表已收敛,无需继续交换")
			break
		}
	}

	return iterations
}

// LinkStateRouting 链路状态路由算法 (类 OSPF - Dijkstra)
// 对应 408 考点: 链路状态算法,最短路径优先
type LinkStateRouting struct {
	Nodes map[string]*Node           // 所有节点
	Graph map[string]map[string]float64 // 全局链路状态数据库
}

// NewLinkStateRouting 创建链路状态路由
func NewLinkStateRouting() *LinkStateRouting {
	return &LinkStateRouting{
		Nodes: make(map[string]*Node),
		Graph: make(map[string]map[string]float64),
	}
}

// AddNode 添加节点
func (ls *LinkStateRouting) AddNode(name string) {
	ls.Nodes[name] = NewNode(name)
	ls.Graph[name] = make(map[string]float64)
}

// AddLink 添加链路 (双向)
func (ls *LinkStateRouting) AddLink(node1, node2 string, cost float64) {
	ls.Graph[node1][node2] = cost
	ls.Graph[node2][node1] = cost
	ls.Nodes[node1].AddNeighbor(node2, cost)
	ls.Nodes[node2].AddNeighbor(node1, cost)
}

// Dijkstra 从源节点运行 Dijkstra 算法
func (ls *LinkStateRouting) Dijkstra(source string) {
	node := ls.Nodes[source]

	// 初始化
	dist := make(map[string]float64)   // 最短距离
	prev := make(map[string]string)    // 前驱节点
	visited := make(map[string]bool)   // 已访问标记

	// 所有节点初始距离为无穷大
	for name := range ls.Nodes {
		dist[name] = math.Inf(1)
	}
	dist[source] = 0

	// Dijkstra 主循环
	for len(visited) < len(ls.Nodes) {
		// 找到未访问节点中距离最小的
		var current string
		minDist := math.Inf(1)
		for name := range ls.Nodes {
			if !visited[name] && dist[name] < minDist {
				current = name
				minDist = dist[name]
			}
		}

		if current == "" {
			break // 没有可达节点
		}

		visited[current] = true

		// 更新邻居节点的距离
		for neighbor, cost := range ls.Graph[current] {
			if !visited[neighbor] {
				newDist := dist[current] + cost
				if newDist < dist[neighbor] {
					dist[neighbor] = newDist
					prev[neighbor] = current
				}
			}
		}
	}

	// 构建路由表
	for dest := range ls.Nodes {
		if dest == source {
			continue
		}

		// 回溯找到下一跳
		nextHop := dest
		for prev[nextHop] != source && prev[nextHop] != "" {
			nextHop = prev[nextHop]
		}

		if dist[dest] != math.Inf(1) {
			node.RoutingTable[dest] = &RoutingTable{
				Destination: dest,
				NextHop:     nextHop,
				Cost:        dist[dest],
			}
		}
	}
}

// Run 运行链路状态算法
func (ls *LinkStateRouting) Run() {
	fmt.Println("\n【链路状态路由算法 - 类 OSPF/Dijkstra】")
	fmt.Println("每个节点独立计算最短路径树...")

	// 每个节点运行 Dijkstra 算法
	for name := range ls.Nodes {
		ls.Dijkstra(name)
	}

	fmt.Println("✓ 所有节点已完成最短路径计算")
}

// RoutingExample 路由算法示例
func RoutingExample() {
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("【网络层 - 路由算法示例】")
	fmt.Println(strings.Repeat("─", 50))

	// 构建网络拓扑
	//       A
	//      / \
	//     2   3
	//    /     \
	//   B---1---C
	//    \     /
	//     4   2
	//      \ /
	//       D

	fmt.Println("\n网络拓扑:")
	fmt.Println("         A")
	fmt.Println("        / \\")
	fmt.Println("      2/   \\3")
	fmt.Println("      /     \\")
	fmt.Println("     B---1---C")
	fmt.Println("      \\     /")
	fmt.Println("      4\\   /2")
	fmt.Println("        \\ /")
	fmt.Println("         D")

	// 1. 距离向量路由算法
	fmt.Println("\n" + strings.Repeat("═", 50))
	dv := NewDistanceVectorRouting()
	dv.AddNode("A")
	dv.AddNode("B")
	dv.AddNode("C")
	dv.AddNode("D")
	dv.AddLink("A", "B", 2)
	dv.AddLink("A", "C", 3)
	dv.AddLink("B", "C", 1)
	dv.AddLink("B", "D", 4)
	dv.AddLink("C", "D", 2)

	iterations := dv.Run(10)
	fmt.Printf("\n算法在 %d 轮后收敛\n", iterations)

	fmt.Println("\n最终路由表:")
	for _, nodeName := range []string{"A", "B", "C", "D"} {
		dv.Nodes[nodeName].PrintRoutingTable()
	}

	// 2. 链路状态路由算法
	fmt.Println("\n" + strings.Repeat("═", 50))
	ls := NewLinkStateRouting()
	ls.AddNode("A")
	ls.AddNode("B")
	ls.AddNode("C")
	ls.AddNode("D")
	ls.AddLink("A", "B", 2)
	ls.AddLink("A", "C", 3)
	ls.AddLink("B", "C", 1)
	ls.AddLink("B", "D", 4)
	ls.AddLink("C", "D", 2)

	ls.Run()

	fmt.Println("\n最终路由表:")
	for _, nodeName := range []string{"A", "B", "C", "D"} {
		ls.Nodes[nodeName].PrintRoutingTable()
	}

	// 408 考点提示
	fmt.Println("\n📚 408 考点总结:")
	fmt.Println("  ✓ 距离向量算法 (RIP): 与邻居交换路由信息,Bellman-Ford")
	fmt.Println("  ✓ 链路状态算法 (OSPF): 全局链路状态,Dijkstra 最短路径")
	fmt.Println("  ✓ RIP 特点: 周期性更新,跳数限制 15,慢收敛")
	fmt.Println("  ✓ OSPF 特点: 事件触发更新,层次化,快速收敛")
}
