package algorithm

import "fmt"

// ============================================================
// 回溯法
// 408考点：回溯思想、解空间树、剪枝策略
// ============================================================

// --- 1. N皇后问题 ---
// 在 n×n 棋盘上放置 n 个皇后，使任意两个皇后不在同一行、列、对角线

// SolveNQueens 求解N皇后问题，返回所有解
func SolveNQueens(n int) [][]int {
	solutions := make([][]int, 0)
	board := make([]int, n) // board[i] 表示第i行皇后放在第几列

	var backtrack func(row int)
	backtrack = func(row int) {
		if row == n {
			// 找到一个解
			solution := make([]int, n)
			copy(solution, board)
			solutions = append(solutions, solution)
			return
		}

		for col := 0; col < n; col++ {
			if isQueenSafe(board, row, col) {
				board[row] = col
				backtrack(row + 1)
				// 回溯：自动被下次赋值覆盖
			}
		}
	}

	backtrack(0)
	return solutions
}

// isQueenSafe 检查在(row, col)放置皇后是否安全
func isQueenSafe(board []int, row, col int) bool {
	for i := 0; i < row; i++ {
		// 同列检查
		if board[i] == col {
			return false
		}
		// 对角线检查
		if abs(board[i]-col) == abs(i-row) {
			return false
		}
	}
	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// PrintQueenBoard 可视化打印棋盘
func PrintQueenBoard(solution []int) {
	n := len(solution)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if solution[i] == j {
				fmt.Print("Q ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
}

// --- 2. 图着色问题 ---
// 用 m 种颜色给图的顶点着色，使相邻顶点颜色不同

// GraphColoring 图着色
func GraphColoring(adjMatrix [][]bool, numColors int) [][]int {
	n := len(adjMatrix)
	solutions := make([][]int, 0)
	colors := make([]int, n)

	var backtrack func(v int)
	backtrack = func(v int) {
		if v == n {
			solution := make([]int, n)
			copy(solution, colors)
			solutions = append(solutions, solution)
			return
		}

		for c := 1; c <= numColors; c++ {
			if isColorSafe(adjMatrix, colors, v, c) {
				colors[v] = c
				backtrack(v + 1)
				colors[v] = 0 // 回溯
			}
		}
	}

	backtrack(0)
	return solutions
}

// isColorSafe 检查给顶点v着色c是否安全
func isColorSafe(adjMatrix [][]bool, colors []int, v, c int) bool {
	for i := 0; i < len(adjMatrix); i++ {
		if adjMatrix[v][i] && colors[i] == c {
			return false
		}
	}
	return true
}

// --- 3. 子集和问题 ---
// 从集合中找出所有和为目标值的子集

// SubsetSum 求子集和
func SubsetSum(set []int, target int) [][]int {
	solutions := make([][]int, 0)
	current := make([]int, 0)

	var backtrack func(idx, currentSum int)
	backtrack = func(idx, currentSum int) {
		if currentSum == target {
			solution := make([]int, len(current))
			copy(solution, current)
			solutions = append(solutions, solution)
			return
		}

		if idx >= len(set) || currentSum > target {
			return // 剪枝
		}

		// 选择当前元素
		current = append(current, set[idx])
		backtrack(idx+1, currentSum+set[idx])

		// 不选当前元素（回溯）
		current = current[:len(current)-1]
		backtrack(idx+1, currentSum)
	}

	backtrack(0, 0)
	return solutions
}

// --- 4. 全排列生成 ---
// 生成数组的所有排列

// Permutations 生成全排列
func Permutations(arr []int) [][]int {
	results := make([][]int, 0)
	n := len(arr)
	perm := make([]int, n)
	copy(perm, arr)

	var backtrack func(first int)
	backtrack = func(first int) {
		if first == n {
			result := make([]int, n)
			copy(result, perm)
			results = append(results, result)
			return
		}

		for i := first; i < n; i++ {
			// 交换
			perm[first], perm[i] = perm[i], perm[first]
			backtrack(first + 1)
			// 回溯
			perm[first], perm[i] = perm[i], perm[first]
		}
	}

	backtrack(0)
	return results
}

// BacktrackingExample 回溯法示例
func BacktrackingExample() {
	fmt.Println("\n--- 回溯法 ---")

	// 1. N皇后
	fmt.Println("\n【4皇后问题】")
	solutions := SolveNQueens(4)
	fmt.Printf("4皇后共有 %d 个解:\n", len(solutions))
	for i, sol := range solutions {
		fmt.Printf("\n解 %d: 列位置=%v\n", i+1, sol)
		PrintQueenBoard(sol)
	}

	fmt.Println("\n【8皇后问题】")
	solutions8 := SolveNQueens(8)
	fmt.Printf("8皇后共有 %d 个解\n", len(solutions8))

	// 2. 图着色
	fmt.Println("\n【图着色问题】")
	adjMatrix := [][]bool{
		{false, true, true, true},
		{true, false, true, false},
		{true, true, false, true},
		{true, false, true, false},
	}
	fmt.Println("邻接矩阵 (4个顶点):")
	for _, row := range adjMatrix {
		fmt.Printf("  %v\n", row)
	}
	colorSolutions := GraphColoring(adjMatrix, 3)
	fmt.Printf("3种颜色的着色方案数: %d\n", len(colorSolutions))
	if len(colorSolutions) > 0 {
		fmt.Printf("第一个方案: %v\n", colorSolutions[0])
	}

	// 3. 子集和
	fmt.Println("\n【子集和问题】")
	set := []int{3, 7, 1, 8, 4}
	target := 11
	fmt.Printf("集合: %v, 目标和: %d\n", set, target)
	subsets := SubsetSum(set, target)
	fmt.Printf("满足条件的子集:\n")
	for _, s := range subsets {
		fmt.Printf("  %v\n", s)
	}

	// 4. 全排列
	fmt.Println("\n【全排列】")
	arr := []int{1, 2, 3}
	fmt.Printf("数组: %v\n", arr)
	perms := Permutations(arr)
	fmt.Printf("全排列 (%d个):\n", len(perms))
	for _, p := range perms {
		fmt.Printf("  %v\n", p)
	}

	fmt.Println("\n408考点总结:")
	fmt.Println("  - 回溯法本质是DFS+剪枝")
	fmt.Println("  - 解空间树：子集树(2^n)、排列树(n!)")
	fmt.Println("  - 约束条件剪枝可大幅减少搜索空间")
}
