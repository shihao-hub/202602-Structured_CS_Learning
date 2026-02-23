package algorithm

import (
	"fmt"
	"math"
)

// ============================================================
// 查找算法实现
// 408考点：顺序查找、二分查找、分块查找的原理与复杂度
// ============================================================

// SearchResult 查找结果
type SearchResult struct {
	Index       int  // 找到的位置，-1表示未找到
	Found       bool // 是否找到
	Comparisons int  // 比较次数
}

// --- 1. 顺序查找 (Sequential Search) ---
// 时间复杂度: 最好 O(1), 平均 O(n), 最坏 O(n)
// 适用场景: 无序表或链式存储
func SequentialSearch(arr []int, target int) SearchResult {
	comparisons := 0
	for i, v := range arr {
		comparisons++
		if v == target {
			return SearchResult{Index: i, Found: true, Comparisons: comparisons}
		}
	}
	return SearchResult{Index: -1, Found: false, Comparisons: comparisons}
}

// SequentialSearchWithSentinel 带哨兵的顺序查找
// 减少循环中的边界判断，效率略高
func SequentialSearchWithSentinel(arr []int, target int) SearchResult {
	n := len(arr)
	if n == 0 {
		return SearchResult{Index: -1, Found: false, Comparisons: 0}
	}

	// 复制数组并在末尾添加哨兵
	data := make([]int, n+1)
	copy(data, arr)
	data[n] = target // 哨兵

	comparisons := 0
	i := 0
	for data[i] != target {
		comparisons++
		i++
	}
	comparisons++ // 最后一次匹配的比较

	if i < n {
		return SearchResult{Index: i, Found: true, Comparisons: comparisons}
	}
	return SearchResult{Index: -1, Found: false, Comparisons: comparisons}
}

// --- 2. 二分查找 (Binary Search) ---
// 前提: 数组必须有序
// 时间复杂度: O(logn)
// 空间复杂度: O(1)（迭代版）

// BinarySearch 二分查找（迭代版）
func BinarySearch(arr []int, target int) SearchResult {
	low, high := 0, len(arr)-1
	comparisons := 0

	for low <= high {
		mid := low + (high-low)/2 // 防止溢出
		comparisons++

		if arr[mid] == target {
			return SearchResult{Index: mid, Found: true, Comparisons: comparisons}
		} else if arr[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return SearchResult{Index: -1, Found: false, Comparisons: comparisons}
}

// BinarySearchRecursive 二分查找（递归版）
// 空间复杂度: O(logn)（递归栈）
func BinarySearchRecursive(arr []int, target int) SearchResult {
	comparisons := 0
	idx := binarySearchHelper(arr, target, 0, len(arr)-1, &comparisons)
	return SearchResult{Index: idx, Found: idx != -1, Comparisons: comparisons}
}

func binarySearchHelper(arr []int, target, low, high int, comparisons *int) int {
	if low > high {
		return -1
	}
	mid := low + (high-low)/2
	*comparisons++

	if arr[mid] == target {
		return mid
	} else if arr[mid] < target {
		return binarySearchHelper(arr, target, mid+1, high, comparisons)
	}
	return binarySearchHelper(arr, target, low, mid-1, comparisons)
}

// --- 3. 插值查找 (Interpolation Search) ---
// 适用于数据分布均匀的有序表
// 时间复杂度: 平均 O(log(logn))，最坏 O(n)
func InterpolationSearch(arr []int, target int) SearchResult {
	low, high := 0, len(arr)-1
	comparisons := 0

	for low <= high && target >= arr[low] && target <= arr[high] {
		comparisons++

		// 插值公式：根据值的分布估算位置
		if arr[high] == arr[low] {
			if arr[low] == target {
				return SearchResult{Index: low, Found: true, Comparisons: comparisons}
			}
			break
		}

		mid := low + int(float64(high-low)*float64(target-arr[low])/float64(arr[high]-arr[low]))

		if mid < low || mid > high {
			break
		}

		if arr[mid] == target {
			return SearchResult{Index: mid, Found: true, Comparisons: comparisons}
		} else if arr[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return SearchResult{Index: -1, Found: false, Comparisons: comparisons}
}

// --- 4. 分块查找 (Block Search / Index Sequential Search) ---
// 块间有序，块内无序
// 时间复杂度: O(√n)（块间二分+块内顺序）

// Block 索引块
type Block struct {
	MaxKey   int   // 块中最大关键字
	StartIdx int   // 块的起始位置
	Elements []int // 块内元素
}

// BlockSearch 分块查找
func BlockSearch(blocks []Block, target int) SearchResult {
	comparisons := 0

	// 第一步：在索引表中确定目标所在的块（二分查找）
	blockIdx := -1
	low, high := 0, len(blocks)-1
	for low <= high {
		mid := low + (high-low)/2
		comparisons++
		if target <= blocks[mid].MaxKey {
			blockIdx = mid
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	if blockIdx == -1 {
		return SearchResult{Index: -1, Found: false, Comparisons: comparisons}
	}

	// 第二步：在确定的块内顺序查找
	block := blocks[blockIdx]
	for i, v := range block.Elements {
		comparisons++
		if v == target {
			return SearchResult{
				Index:       block.StartIdx + i,
				Found:       true,
				Comparisons: comparisons,
			}
		}
	}
	return SearchResult{Index: -1, Found: false, Comparisons: comparisons}
}

// SearchingExample 查找算法示例
func SearchingExample() {
	fmt.Println("\n--- 查找算法对比 ---")

	// 有序数组用于二分查找
	sortedArr := []int{2, 5, 8, 12, 16, 23, 38, 45, 56, 72, 91}
	target := 23
	fmt.Printf("有序数组: %v\n", sortedArr)
	fmt.Printf("查找目标: %d\n\n", target)

	// 1. 顺序查找
	result := SequentialSearch(sortedArr, target)
	fmt.Printf("顺序查找: 位置=%d, 比较次数=%d\n", result.Index, result.Comparisons)

	// 2. 二分查找（迭代）
	result = BinarySearch(sortedArr, target)
	fmt.Printf("二分查找(迭代): 位置=%d, 比较次数=%d\n", result.Index, result.Comparisons)

	// 3. 二分查找（递归）
	result = BinarySearchRecursive(sortedArr, target)
	fmt.Printf("二分查找(递归): 位置=%d, 比较次数=%d\n", result.Index, result.Comparisons)

	// 4. 插值查找
	result = InterpolationSearch(sortedArr, target)
	fmt.Printf("插值查找: 位置=%d, 比较次数=%d\n", result.Index, result.Comparisons)

	// 5. 分块查找
	fmt.Println("\n--- 分块查找示例 ---")
	// 将数组分为3块
	blockSize := int(math.Ceil(float64(len(sortedArr)) / 3.0))
	blocks := make([]Block, 0)
	for i := 0; i < len(sortedArr); i += blockSize {
		end := i + blockSize
		if end > len(sortedArr) {
			end = len(sortedArr)
		}
		elements := sortedArr[i:end]
		maxKey := elements[0]
		for _, v := range elements {
			if v > maxKey {
				maxKey = v
			}
		}
		blocks = append(blocks, Block{MaxKey: maxKey, StartIdx: i, Elements: elements})
	}

	fmt.Printf("索引表: ")
	for _, b := range blocks {
		fmt.Printf("[max=%d, start=%d] ", b.MaxKey, b.StartIdx)
	}
	fmt.Println()

	result = BlockSearch(blocks, target)
	fmt.Printf("分块查找: 位置=%d, 比较次数=%d\n", result.Index, result.Comparisons)

	fmt.Println("\n408考点总结:")
	fmt.Println("  - 顺序查找ASL = (n+1)/2")
	fmt.Println("  - 二分查找ASL = log₂(n+1) - 1，仅适用于有序顺序表")
	fmt.Println("  - 分块查找结合了索引查找和顺序查找的优点")
	fmt.Println("  - 二分查找的判定树是一棵平衡二叉排序树")
}
