package algorithm

import (
	"container/heap"
	"fmt"
	"math"
	"sort"
)

// ============================================================
// 贪心算法
// 408考点：贪心策略、最优子结构、与DP的区别
// ============================================================

// --- 1. 活动选择问题 ---
// 贪心策略：每次选择结束时间最早且与已选活动不冲突的活动

// Activity 活动定义
type Activity struct {
	Name  string
	Start int
	End   int
}

// ActivitySelection 贪心活动选择
func ActivitySelection(activities []Activity) []Activity {
	// 按结束时间排序
	sorted := make([]Activity, len(activities))
	copy(sorted, activities)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].End < sorted[j].End
	})

	selected := []Activity{sorted[0]}
	lastEnd := sorted[0].End

	for i := 1; i < len(sorted); i++ {
		if sorted[i].Start >= lastEnd {
			selected = append(selected, sorted[i])
			lastEnd = sorted[i].End
		}
	}
	return selected
}

// --- 2. 哈夫曼编码 ---
// 贪心策略：每次取频率最小的两个节点合并

// HuffmanNode 哈夫曼树节点
type HuffmanNode struct {
	Char  rune
	Freq  int
	Left  *HuffmanNode
	Right *HuffmanNode
}

// huffmanHeap 用于构建最小堆
type huffmanHeap []*HuffmanNode

func (h huffmanHeap) Len() int            { return len(h) }
func (h huffmanHeap) Less(i, j int) bool  { return h[i].Freq < h[j].Freq }
func (h huffmanHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *huffmanHeap) Push(x interface{}) { *h = append(*h, x.(*HuffmanNode)) }
func (h *huffmanHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

// BuildHuffmanTree 构建哈夫曼树
func BuildHuffmanTree(charFreq map[rune]int) *HuffmanNode {
	h := &huffmanHeap{}
	heap.Init(h)

	for char, freq := range charFreq {
		heap.Push(h, &HuffmanNode{Char: char, Freq: freq})
	}

	for h.Len() > 1 {
		left := heap.Pop(h).(*HuffmanNode)
		right := heap.Pop(h).(*HuffmanNode)
		parent := &HuffmanNode{
			Freq:  left.Freq + right.Freq,
			Left:  left,
			Right: right,
		}
		heap.Push(h, parent)
	}

	return heap.Pop(h).(*HuffmanNode)
}

// GetHuffmanCodes 获取哈夫曼编码表
func GetHuffmanCodes(root *HuffmanNode) map[rune]string {
	codes := make(map[rune]string)
	var traverse func(node *HuffmanNode, code string)
	traverse = func(node *HuffmanNode, code string) {
		if node == nil {
			return
		}
		if node.Left == nil && node.Right == nil {
			if code == "" {
				code = "0" // 只有一个字符的特殊情况
			}
			codes[node.Char] = code
			return
		}
		traverse(node.Left, code+"0")
		traverse(node.Right, code+"1")
	}
	traverse(root, "")
	return codes
}

// --- 3. Prim 最小生成树 ---
// 贪心策略：每次选择连接已选顶点和未选顶点的最小权边

// Edge 边定义
type Edge struct {
	From, To int
	Weight   float64
}

// PrimMST Prim算法（邻接矩阵版）
func PrimMST(graph [][]float64) ([]Edge, float64) {
	n := len(graph)
	if n == 0 {
		return nil, 0
	}

	inMST := make([]bool, n)
	key := make([]float64, n) // 到MST的最小距离
	parent := make([]int, n)  // 父节点

	for i := range key {
		key[i] = math.Inf(1)
		parent[i] = -1
	}
	key[0] = 0

	var mstEdges []Edge
	totalWeight := 0.0

	for count := 0; count < n; count++ {
		// 找到未加入MST的最小key顶点
		u := -1
		minKey := math.Inf(1)
		for v := 0; v < n; v++ {
			if !inMST[v] && key[v] < minKey {
				minKey = key[v]
				u = v
			}
		}

		if u == -1 {
			break
		}

		inMST[u] = true
		if parent[u] != -1 {
			mstEdges = append(mstEdges, Edge{From: parent[u], To: u, Weight: graph[parent[u]][u]})
			totalWeight += graph[parent[u]][u]
		}

		// 更新相邻顶点的key
		for v := 0; v < n; v++ {
			if !inMST[v] && graph[u][v] < key[v] {
				key[v] = graph[u][v]
				parent[v] = u
			}
		}
	}
	return mstEdges, totalWeight
}

// --- 4. Kruskal 最小生成树 ---
// 贪心策略：按边权排序，每次选最小边且不构成环

// UnionFind 并查集
type UnionFind struct {
	parent []int
	rank   []int
}

// NewUnionFind 创建并查集
func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	rank := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	return &UnionFind{parent: parent, rank: rank}
}

// Find 查找根节点（带路径压缩）
func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x])
	}
	return uf.parent[x]
}

// Union 合并两个集合（按秩合并）
func (uf *UnionFind) Union(x, y int) bool {
	px, py := uf.Find(x), uf.Find(y)
	if px == py {
		return false // 已在同一集合
	}
	if uf.rank[px] < uf.rank[py] {
		px, py = py, px
	}
	uf.parent[py] = px
	if uf.rank[px] == uf.rank[py] {
		uf.rank[px]++
	}
	return true
}

// KruskalMST Kruskal算法
func KruskalMST(n int, edges []Edge) ([]Edge, float64) {
	// 按权重排序
	sorted := make([]Edge, len(edges))
	copy(sorted, edges)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Weight < sorted[j].Weight
	})

	uf := NewUnionFind(n)
	var mstEdges []Edge
	totalWeight := 0.0

	for _, e := range sorted {
		if uf.Union(e.From, e.To) {
			mstEdges = append(mstEdges, e)
			totalWeight += e.Weight
			if len(mstEdges) == n-1 {
				break
			}
		}
	}
	return mstEdges, totalWeight
}

// GreedyExample 贪心算法示例
func GreedyExample() {
	fmt.Println("\n--- 贪心算法 ---")

	// 1. 活动选择
	fmt.Println("\n【活动选择问题】")
	activities := []Activity{
		{"A1", 1, 4}, {"A2", 3, 5}, {"A3", 0, 6},
		{"A4", 5, 7}, {"A5", 3, 9}, {"A6", 5, 9},
		{"A7", 6, 10}, {"A8", 8, 11}, {"A9", 8, 12},
		{"A10", 2, 14}, {"A11", 12, 16},
	}
	fmt.Println("活动列表:")
	for _, a := range activities {
		fmt.Printf("  %s: [%d, %d)\n", a.Name, a.Start, a.End)
	}
	selected := ActivitySelection(activities)
	fmt.Printf("贪心选择（按结束时间）: ")
	for _, a := range selected {
		fmt.Printf("%s ", a.Name)
	}
	fmt.Printf("(共%d个活动)\n", len(selected))

	// 2. 哈夫曼编码
	fmt.Println("\n【哈夫曼编码】")
	charFreq := map[rune]int{
		'a': 45, 'b': 13, 'c': 12, 'd': 16, 'e': 9, 'f': 5,
	}
	fmt.Printf("字符频率: %v\n", charFreq)
	root := BuildHuffmanTree(charFreq)
	codes := GetHuffmanCodes(root)
	fmt.Println("哈夫曼编码:")
	totalBits := 0
	for char, code := range codes {
		freq := charFreq[char]
		fmt.Printf("  '%c' (频率=%d): %s\n", char, freq, code)
		totalBits += freq * len(code)
	}
	fmt.Printf("加权路径长度(WPL): %d\n", totalBits)

	// 3. Prim MST
	fmt.Println("\n【Prim最小生成树】")
	inf := math.Inf(1)
	graph := [][]float64{
		{0, 2, inf, 6, inf},
		{2, 0, 3, 8, 5},
		{inf, 3, 0, inf, 7},
		{6, 8, inf, 0, 9},
		{inf, 5, 7, 9, 0},
	}
	mstEdges, totalWeight := PrimMST(graph)
	fmt.Println("MST边:")
	for _, e := range mstEdges {
		fmt.Printf("  %d -- %d (权重=%.0f)\n", e.From, e.To, e.Weight)
	}
	fmt.Printf("总权重: %.0f\n", totalWeight)

	// 4. Kruskal MST
	fmt.Println("\n【Kruskal最小生成树】")
	edges := []Edge{
		{0, 1, 2}, {0, 3, 6}, {1, 2, 3},
		{1, 3, 8}, {1, 4, 5}, {2, 4, 7}, {3, 4, 9},
	}
	mstEdges2, totalWeight2 := KruskalMST(5, edges)
	fmt.Println("MST边:")
	for _, e := range mstEdges2 {
		fmt.Printf("  %d -- %d (权重=%.0f)\n", e.From, e.To, e.Weight)
	}
	fmt.Printf("总权重: %.0f\n", totalWeight2)

	fmt.Println("\n408考点总结:")
	fmt.Println("  - 贪心与DP的区别：贪心无回溯，DP保存子问题解")
	fmt.Println("  - Prim适合稠密图 O(V²)，Kruskal适合稀疏图 O(ElogE)")
	fmt.Println("  - 哈夫曼编码是前缀编码，WPL最小")
}
