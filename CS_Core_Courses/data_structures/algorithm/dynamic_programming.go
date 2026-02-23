package algorithm

import (
	"fmt"
	"math"
)

// ============================================================
// 动态规划经典问题
// 408考点：DP思想、状态定义、状态转移方程、最优子结构
// ============================================================

// --- 1. 0/1 背包问题 ---
// 状态定义: dp[i][w] = 前i个物品在容量w下的最大价值
// 状态转移: dp[i][w] = max(dp[i-1][w], dp[i-1][w-weight[i]] + value[i])
func Knapsack01(weights, values []int, capacity int) (int, [][]int) {
	n := len(weights)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}

	for i := 1; i <= n; i++ {
		for w := 0; w <= capacity; w++ {
			dp[i][w] = dp[i-1][w] // 不选第i个物品
			if w >= weights[i-1] && dp[i-1][w-weights[i-1]]+values[i-1] > dp[i][w] {
				dp[i][w] = dp[i-1][w-weights[i-1]] + values[i-1] // 选第i个物品
			}
		}
	}
	return dp[n][capacity], dp
}

// Knapsack01Optimized 空间优化版（一维数组）
func Knapsack01Optimized(weights, values []int, capacity int) int {
	dp := make([]int, capacity+1)

	for i := 0; i < len(weights); i++ {
		// 逆序遍历，确保每个物品只用一次
		for w := capacity; w >= weights[i]; w-- {
			if dp[w-weights[i]]+values[i] > dp[w] {
				dp[w] = dp[w-weights[i]] + values[i]
			}
		}
	}
	return dp[capacity]
}

// --- 2. 完全背包问题 ---
// 每个物品可以选无限次
// 状态转移: dp[i][w] = max(dp[i-1][w], dp[i][w-weight[i]] + value[i])
// 注意与0/1背包的区别：第二项用dp[i]而非dp[i-1]
func KnapsackComplete(weights, values []int, capacity int) int {
	dp := make([]int, capacity+1)

	for i := 0; i < len(weights); i++ {
		// 正序遍历（与0/1背包的逆序区分）
		for w := weights[i]; w <= capacity; w++ {
			if dp[w-weights[i]]+values[i] > dp[w] {
				dp[w] = dp[w-weights[i]] + values[i]
			}
		}
	}
	return dp[capacity]
}

// --- 3. 最长公共子序列 (LCS) ---
// 状态定义: dp[i][j] = X[0..i-1]和Y[0..j-1]的LCS长度
// 状态转移:
//
//	if X[i-1] == Y[j-1]: dp[i][j] = dp[i-1][j-1] + 1
//	else: dp[i][j] = max(dp[i-1][j], dp[i][j-1])
func LCS(x, y string) (int, string) {
	m, n := len(x), len(y)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if x[i-1] == y[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = dp[i-1][j]
				if dp[i][j-1] > dp[i][j] {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}

	// 回溯找出LCS内容
	lcs := make([]byte, 0)
	i, j := m, n
	for i > 0 && j > 0 {
		if x[i-1] == y[j-1] {
			lcs = append([]byte{x[i-1]}, lcs...)
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			i--
		} else {
			j--
		}
	}

	return dp[m][n], string(lcs)
}

// --- 4. 最长递增子序列 (LIS) ---
// 状态定义: dp[i] = 以arr[i]结尾的最长递增子序列长度
// 状态转移: dp[i] = max(dp[j] + 1)，其中 j < i 且 arr[j] < arr[i]
// 时间复杂度: O(n²)
func LIS(arr []int) (int, []int) {
	n := len(arr)
	if n == 0 {
		return 0, nil
	}

	dp := make([]int, n)
	prev := make([]int, n) // 记录前驱用于回溯

	for i := range dp {
		dp[i] = 1
		prev[i] = -1
	}

	maxLen := 1
	maxIdx := 0

	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if arr[j] < arr[i] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
				prev[i] = j
			}
		}
		if dp[i] > maxLen {
			maxLen = dp[i]
			maxIdx = i
		}
	}

	// 回溯找出LIS
	lis := make([]int, 0)
	for idx := maxIdx; idx != -1; idx = prev[idx] {
		lis = append([]int{arr[idx]}, lis...)
	}

	return maxLen, lis
}

// --- 5. Floyd 最短路径算法 ---
// 求所有顶点对之间的最短路径
// 状态定义: dist[i][j] = 从i到j的最短距离
// 状态转移: dist[i][j] = min(dist[i][j], dist[i][k] + dist[k][j])
// 时间复杂度: O(n³)
func Floyd(graph [][]float64) [][]float64 {
	n := len(graph)
	dist := make([][]float64, n)
	for i := range dist {
		dist[i] = make([]float64, n)
		copy(dist[i], graph[i])
	}

	// k为中间节点
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if dist[i][k]+dist[k][j] < dist[i][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}
	return dist
}

// --- 6. 矩阵链乘法 ---
// 确定矩阵连乘的最优加括号方式，使乘法次数最少
// 状态定义: dp[i][j] = 计算 A_i...A_j 所需的最少乘法次数
// 状态转移: dp[i][j] = min(dp[i][k] + dp[k+1][j] + p[i-1]*p[k]*p[j])
func MatrixChainMultiplication(dimensions []int) (int, [][]int) {
	n := len(dimensions) - 1 // 矩阵个数
	dp := make([][]int, n)
	split := make([][]int, n) // 记录最优分割点

	for i := range dp {
		dp[i] = make([]int, n)
		split[i] = make([]int, n)
	}

	// l 为链长
	for l := 2; l <= n; l++ {
		for i := 0; i <= n-l; i++ {
			j := i + l - 1
			dp[i][j] = math.MaxInt32
			for k := i; k < j; k++ {
				cost := dp[i][k] + dp[k+1][j] + dimensions[i]*dimensions[k+1]*dimensions[j+1]
				if cost < dp[i][j] {
					dp[i][j] = cost
					split[i][j] = k
				}
			}
		}
	}
	return dp[0][n-1], split
}

// printOptimalParens 打印最优括号方案
func printOptimalParens(split [][]int, i, j int) string {
	if i == j {
		return fmt.Sprintf("A%d", i+1)
	}
	return "(" + printOptimalParens(split, i, split[i][j]) + " x " +
		printOptimalParens(split, split[i][j]+1, j) + ")"
}

// DPExample 动态规划示例
func DPExample() {
	fmt.Println("\n--- 动态规划经典问题 ---")

	// 1. 0/1 背包
	fmt.Println("\n【0/1背包问题】")
	weights := []int{2, 3, 4, 5}
	values := []int{3, 4, 5, 6}
	capacity := 8
	fmt.Printf("物品: 重量=%v, 价值=%v, 背包容量=%d\n", weights, values, capacity)
	maxValue, _ := Knapsack01(weights, values, capacity)
	fmt.Printf("最大价值: %d\n", maxValue)
	maxValueOpt := Knapsack01Optimized(weights, values, capacity)
	fmt.Printf("空间优化版最大价值: %d\n", maxValueOpt)

	// 2. 完全背包
	fmt.Println("\n【完全背包问题】")
	fmt.Printf("物品: 重量=%v, 价值=%v, 背包容量=%d (每种物品无限)\n", weights, values, capacity)
	maxValueC := KnapsackComplete(weights, values, capacity)
	fmt.Printf("最大价值: %d\n", maxValueC)

	// 3. LCS
	fmt.Println("\n【最长公共子序列 LCS】")
	x, y := "ABCBDAB", "BDCABA"
	fmt.Printf("X = \"%s\", Y = \"%s\"\n", x, y)
	length, lcs := LCS(x, y)
	fmt.Printf("LCS长度: %d, LCS: \"%s\"\n", length, lcs)

	// 4. LIS
	fmt.Println("\n【最长递增子序列 LIS】")
	arr := []int{10, 9, 2, 5, 3, 7, 101, 18}
	fmt.Printf("数组: %v\n", arr)
	lisLen, lis := LIS(arr)
	fmt.Printf("LIS长度: %d, LIS: %v\n", lisLen, lis)

	// 5. Floyd
	fmt.Println("\n【Floyd最短路径】")
	inf := math.Inf(1)
	graph := [][]float64{
		{0, 3, inf, 7},
		{8, 0, 2, inf},
		{5, inf, 0, 1},
		{2, inf, inf, 0},
	}
	fmt.Println("邻接矩阵 (inf表示不可达):")
	for i, row := range graph {
		fmt.Printf("  %d: ", i)
		for _, v := range row {
			if math.IsInf(v, 1) {
				fmt.Printf("%5s ", "inf")
			} else {
				fmt.Printf("%5.0f ", v)
			}
		}
		fmt.Println()
	}
	dist := Floyd(graph)
	fmt.Println("最短路径矩阵:")
	for i, row := range dist {
		fmt.Printf("  %d: ", i)
		for _, v := range row {
			fmt.Printf("%5.0f ", v)
		}
		fmt.Println()
	}

	// 6. 矩阵链乘法
	fmt.Println("\n【矩阵链乘法】")
	dimensions := []int{30, 35, 15, 5, 10, 20, 25}
	fmt.Printf("矩阵维度: %v (6个矩阵)\n", dimensions)
	minCost, split := MatrixChainMultiplication(dimensions)
	fmt.Printf("最少乘法次数: %d\n", minCost)
	fmt.Printf("最优加括号方案: %s\n", printOptimalParens(split, 0, len(dimensions)-2))
}
