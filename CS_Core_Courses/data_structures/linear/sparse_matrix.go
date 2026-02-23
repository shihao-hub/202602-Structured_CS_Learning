package linear

import "fmt"

// ============================================================
// 稀疏矩阵 (Sparse Matrix)
// 408考点：三元组表示法、矩阵转置（快速转置算法）
// ============================================================

// Triple 三元组
type Triple struct {
	Row int
	Col int
	Val int
}

// SparseMatrix 稀疏矩阵（三元组顺序表）
type SparseMatrix struct {
	Rows     int      // 行数
	Cols     int      // 列数
	NumTerms int      // 非零元素个数
	Data     []Triple // 三元组数组（按行优先排列）
}

// NewSparseMatrix 从二维数组创建稀疏矩阵
func NewSparseMatrix(matrix [][]int) *SparseMatrix {
	rows := len(matrix)
	if rows == 0 {
		return &SparseMatrix{}
	}
	cols := len(matrix[0])

	sm := &SparseMatrix{Rows: rows, Cols: cols}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if matrix[i][j] != 0 {
				sm.Data = append(sm.Data, Triple{Row: i, Col: j, Val: matrix[i][j]})
			}
		}
	}
	sm.NumTerms = len(sm.Data)
	return sm
}

// ToDense 转换回二维数组
func (sm *SparseMatrix) ToDense() [][]int {
	matrix := make([][]int, sm.Rows)
	for i := range matrix {
		matrix[i] = make([]int, sm.Cols)
	}
	for _, t := range sm.Data {
		matrix[t.Row][t.Col] = t.Val
	}
	return matrix
}

// TransposeSimple 简单转置法
// 时间复杂度: O(Cols * NumTerms)
func (sm *SparseMatrix) TransposeSimple() *SparseMatrix {
	result := &SparseMatrix{
		Rows:     sm.Cols,
		Cols:     sm.Rows,
		NumTerms: sm.NumTerms,
		Data:     make([]Triple, 0, sm.NumTerms),
	}

	// 按列号从小到大依次扫描
	for col := 0; col < sm.Cols; col++ {
		for _, t := range sm.Data {
			if t.Col == col {
				result.Data = append(result.Data, Triple{Row: t.Col, Col: t.Row, Val: t.Val})
			}
		}
	}
	return result
}

// TransposeFast 快速转置法
// 时间复杂度: O(Cols + NumTerms)
// 核心思想：预先计算每列非零元素个数和起始位置
func (sm *SparseMatrix) TransposeFast() *SparseMatrix {
	result := &SparseMatrix{
		Rows:     sm.Cols,
		Cols:     sm.Rows,
		NumTerms: sm.NumTerms,
		Data:     make([]Triple, sm.NumTerms),
	}

	if sm.NumTerms == 0 {
		return result
	}

	// num[col] = 原矩阵第col列的非零元素个数
	num := make([]int, sm.Cols)
	for _, t := range sm.Data {
		num[t.Col]++
	}

	// cpot[col] = 原矩阵第col列第一个非零元素在转置矩阵中的位置
	cpot := make([]int, sm.Cols)
	cpot[0] = 0
	for i := 1; i < sm.Cols; i++ {
		cpot[i] = cpot[i-1] + num[i-1]
	}

	// 依次放置每个非零元素
	for _, t := range sm.Data {
		pos := cpot[t.Col]
		result.Data[pos] = Triple{Row: t.Col, Col: t.Row, Val: t.Val}
		cpot[t.Col]++
	}
	return result
}

// Add 稀疏矩阵加法
func (sm *SparseMatrix) Add(other *SparseMatrix) *SparseMatrix {
	if sm.Rows != other.Rows || sm.Cols != other.Cols {
		return nil
	}

	result := &SparseMatrix{Rows: sm.Rows, Cols: sm.Cols}
	i, j := 0, 0

	for i < sm.NumTerms && j < other.NumTerms {
		posA := sm.Data[i].Row*sm.Cols + sm.Data[i].Col
		posB := other.Data[j].Row*other.Cols + other.Data[j].Col

		if posA < posB {
			result.Data = append(result.Data, sm.Data[i])
			i++
		} else if posA > posB {
			result.Data = append(result.Data, other.Data[j])
			j++
		} else {
			sum := sm.Data[i].Val + other.Data[j].Val
			if sum != 0 {
				result.Data = append(result.Data, Triple{Row: sm.Data[i].Row, Col: sm.Data[i].Col, Val: sum})
			}
			i++
			j++
		}
	}

	for ; i < sm.NumTerms; i++ {
		result.Data = append(result.Data, sm.Data[i])
	}
	for ; j < other.NumTerms; j++ {
		result.Data = append(result.Data, other.Data[j])
	}
	result.NumTerms = len(result.Data)
	return result
}

// Print 打印稀疏矩阵信息
func (sm *SparseMatrix) Print() {
	fmt.Printf("稀疏矩阵 (%d×%d, 非零元素=%d)\n", sm.Rows, sm.Cols, sm.NumTerms)
	fmt.Println("三元组表:")
	fmt.Printf("  %5s %5s %5s\n", "行", "列", "值")
	for _, t := range sm.Data {
		fmt.Printf("  %5d %5d %5d\n", t.Row, t.Col, t.Val)
	}
}

// PrintDense 以二维矩阵形式打印
func (sm *SparseMatrix) PrintDense() {
	matrix := sm.ToDense()
	for _, row := range matrix {
		for _, v := range row {
			fmt.Printf("%4d", v)
		}
		fmt.Println()
	}
}

// SparseMatrixExample 稀疏矩阵示例
func SparseMatrixExample() {
	fmt.Println("\n--- 稀疏矩阵 ---")

	// 创建稀疏矩阵
	matrix := [][]int{
		{0, 12, 9, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{-3, 0, 0, 0, 14, 0},
		{0, 0, 24, 0, 0, 0},
		{0, 18, 0, 0, 0, 0},
		{15, 0, 0, -7, 0, 0},
	}

	fmt.Println("原始矩阵:")
	for _, row := range matrix {
		for _, v := range row {
			fmt.Printf("%4d", v)
		}
		fmt.Println()
	}

	sm := NewSparseMatrix(matrix)
	fmt.Println()
	sm.Print()

	// 简单转置
	fmt.Println("\n简单转置结果:")
	transSimple := sm.TransposeSimple()
	transSimple.PrintDense()

	// 快速转置
	fmt.Println("\n快速转置结果:")
	transFast := sm.TransposeFast()
	transFast.PrintDense()

	fmt.Println("\n408考点总结:")
	fmt.Println("  - 三元组表按行优先存储非零元素")
	fmt.Println("  - 简单转置 O(Cols*NumTerms)，快速转置 O(Cols+NumTerms)")
	fmt.Println("  - 快速转置核心：预计算列偏移数组cpot[]")
	fmt.Println("  - 稀疏矩阵适用于非零元素远少于总元素的场景")
}
