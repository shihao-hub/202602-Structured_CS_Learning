package algorithm

import "fmt"

// ============================================================
// 八大排序算法实现
// 408考点：排序算法的时间/空间复杂度、稳定性、适用场景
// ============================================================

// SortResult 排序结果，包含统计信息
type SortResult struct {
	Data        []int
	Comparisons int // 比较次数
	Swaps       int // 交换次数
}

// --- 1. 冒泡排序 (Bubble Sort) ---
// 时间复杂度: 最好 O(n), 平均 O(n²), 最坏 O(n²)
// 空间复杂度: O(1)
// 稳定性: 稳定
func BubbleSort(arr []int) SortResult {
	n := len(arr)
	data := make([]int, n)
	copy(data, arr)
	comparisons, swaps := 0, 0

	for i := 0; i < n-1; i++ {
		swapped := false
		for j := 0; j < n-1-i; j++ {
			comparisons++
			if data[j] > data[j+1] {
				data[j], data[j+1] = data[j+1], data[j]
				swaps++
				swapped = true
			}
		}
		// 优化：如果一趟没有交换，说明已有序
		if !swapped {
			break
		}
	}
	return SortResult{Data: data, Comparisons: comparisons, Swaps: swaps}
}

// --- 2. 选择排序 (Selection Sort) ---
// 时间复杂度: 最好 O(n²), 平均 O(n²), 最坏 O(n²)
// 空间复杂度: O(1)
// 稳定性: 不稳定（例如 [5a, 5b, 2] → [2, 5b, 5a]）
func SelectionSort(arr []int) SortResult {
	n := len(arr)
	data := make([]int, n)
	copy(data, arr)
	comparisons, swaps := 0, 0

	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			comparisons++
			if data[j] < data[minIdx] {
				minIdx = j
			}
		}
		if minIdx != i {
			data[i], data[minIdx] = data[minIdx], data[i]
			swaps++
		}
	}
	return SortResult{Data: data, Comparisons: comparisons, Swaps: swaps}
}

// --- 3. 插入排序 (Insertion Sort) ---
// 时间复杂度: 最好 O(n), 平均 O(n²), 最坏 O(n²)
// 空间复杂度: O(1)
// 稳定性: 稳定
func InsertionSort(arr []int) SortResult {
	n := len(arr)
	data := make([]int, n)
	copy(data, arr)
	comparisons, swaps := 0, 0

	for i := 1; i < n; i++ {
		key := data[i]
		j := i - 1
		for j >= 0 {
			comparisons++
			if data[j] > key {
				data[j+1] = data[j]
				swaps++
				j--
			} else {
				break
			}
		}
		data[j+1] = key
	}
	return SortResult{Data: data, Comparisons: comparisons, Swaps: swaps}
}

// --- 4. 希尔排序 (Shell Sort) ---
// 时间复杂度: 取决于增量序列，约 O(n^1.3)
// 空间复杂度: O(1)
// 稳定性: 不稳定
func ShellSort(arr []int) SortResult {
	n := len(arr)
	data := make([]int, n)
	copy(data, arr)
	comparisons, swaps := 0, 0

	// 使用 Shell 增量序列: n/2, n/4, ..., 1
	for gap := n / 2; gap > 0; gap /= 2 {
		for i := gap; i < n; i++ {
			key := data[i]
			j := i - gap
			for j >= 0 {
				comparisons++
				if data[j] > key {
					data[j+gap] = data[j]
					swaps++
					j -= gap
				} else {
					break
				}
			}
			data[j+gap] = key
		}
	}
	return SortResult{Data: data, Comparisons: comparisons, Swaps: swaps}
}

// --- 5. 归并排序 (Merge Sort) ---
// 时间复杂度: 最好 O(nlogn), 平均 O(nlogn), 最坏 O(nlogn)
// 空间复杂度: O(n)
// 稳定性: 稳定
func MergeSort(arr []int) SortResult {
	n := len(arr)
	data := make([]int, n)
	copy(data, arr)
	comparisons := 0
	swaps := 0

	var mergeSort func(arr []int) []int
	mergeSort = func(arr []int) []int {
		if len(arr) <= 1 {
			return arr
		}
		mid := len(arr) / 2
		left := mergeSort(arr[:mid])
		right := mergeSort(arr[mid:])
		return merge(left, right, &comparisons, &swaps)
	}

	result := mergeSort(data)
	return SortResult{Data: result, Comparisons: comparisons, Swaps: swaps}
}

// merge 合并两个有序数组
func merge(left, right []int, comparisons, moves *int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		*comparisons++
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
		*moves++
	}
	for ; i < len(left); i++ {
		result = append(result, left[i])
		*moves++
	}
	for ; j < len(right); j++ {
		result = append(result, right[j])
		*moves++
	}
	return result
}

// --- 6. 快速排序 (Quick Sort) ---
// 时间复杂度: 最好 O(nlogn), 平均 O(nlogn), 最坏 O(n²)
// 空间复杂度: O(logn)（递归栈）
// 稳定性: 不稳定
func QuickSort(arr []int) SortResult {
	n := len(arr)
	data := make([]int, n)
	copy(data, arr)
	comparisons, swaps := 0, 0

	var quickSort func(low, high int)
	quickSort = func(low, high int) {
		if low < high {
			pivot := partition(data, low, high, &comparisons, &swaps)
			quickSort(low, pivot-1)
			quickSort(pivot+1, high)
		}
	}

	if n > 0 {
		quickSort(0, n-1)
	}
	return SortResult{Data: data, Comparisons: comparisons, Swaps: swaps}
}

// partition 分区操作，选取最后一个元素为基准
func partition(arr []int, low, high int, comparisons, swaps *int) int {
	pivot := arr[high]
	i := low - 1

	for j := low; j < high; j++ {
		*comparisons++
		if arr[j] <= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
			*swaps++
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	*swaps++
	return i + 1
}

// --- 7. 堆排序 (Heap Sort) ---
// 时间复杂度: 最好 O(nlogn), 平均 O(nlogn), 最坏 O(nlogn)
// 空间复杂度: O(1)
// 稳定性: 不稳定
func HeapSort(arr []int) SortResult {
	n := len(arr)
	data := make([]int, n)
	copy(data, arr)
	comparisons, swaps := 0, 0

	// 建立大顶堆（从最后一个非叶子节点开始下沉）
	for i := n/2 - 1; i >= 0; i-- {
		heapify(data, n, i, &comparisons, &swaps)
	}

	// 逐个取出堆顶元素
	for i := n - 1; i > 0; i-- {
		data[0], data[i] = data[i], data[0]
		swaps++
		heapify(data, i, 0, &comparisons, &swaps)
	}
	return SortResult{Data: data, Comparisons: comparisons, Swaps: swaps}
}

// heapify 维护大顶堆性质（下沉操作）
func heapify(arr []int, n, i int, comparisons, swaps *int) {
	largest := i
	left := 2*i + 1
	right := 2*i + 2

	if left < n {
		*comparisons++
		if arr[left] > arr[largest] {
			largest = left
		}
	}
	if right < n {
		*comparisons++
		if arr[right] > arr[largest] {
			largest = right
		}
	}
	if largest != i {
		arr[i], arr[largest] = arr[largest], arr[i]
		*swaps++
		heapify(arr, n, largest, comparisons, swaps)
	}
}

// --- 8. 基数排序 (Radix Sort) ---
// 时间复杂度: O(d*(n+k))，d为最大位数，k为基数(10)
// 空间复杂度: O(n+k)
// 稳定性: 稳定
// 注意：基数排序仅适用于非负整数
func RadixSort(arr []int) SortResult {
	n := len(arr)
	data := make([]int, n)
	copy(data, arr)

	if n <= 1 {
		return SortResult{Data: data}
	}

	// 找最大值确定位数
	maxVal := data[0]
	for _, v := range data {
		if v > maxVal {
			maxVal = v
		}
	}

	// 按每一位进行计数排序
	for exp := 1; maxVal/exp > 0; exp *= 10 {
		countingSortByDigit(data, n, exp)
	}
	return SortResult{Data: data}
}

// countingSortByDigit 按指定位进行计数排序
func countingSortByDigit(arr []int, n, exp int) {
	output := make([]int, n)
	count := make([]int, 10) // 基数为10

	// 统计每个数字出现次数
	for i := 0; i < n; i++ {
		digit := (arr[i] / exp) % 10
		count[digit]++
	}

	// 计算累计计数
	for i := 1; i < 10; i++ {
		count[i] += count[i-1]
	}

	// 从后往前放置元素（保证稳定性）
	for i := n - 1; i >= 0; i-- {
		digit := (arr[i] / exp) % 10
		output[count[digit]-1] = arr[i]
		count[digit]--
	}

	copy(arr, output)
}

// SortingExample 排序算法示例
func SortingExample() {
	fmt.Println("\n--- 八大排序算法对比 ---")

	arr := []int{64, 34, 25, 12, 22, 11, 90, 78, 45, 33}
	fmt.Printf("原始数组: %v\n\n", arr)

	// 排序算法列表
	type sortFunc struct {
		name   string
		fn     func([]int) SortResult
		stable bool
	}

	algorithms := []sortFunc{
		{"冒泡排序", BubbleSort, true},
		{"选择排序", SelectionSort, false},
		{"插入排序", InsertionSort, true},
		{"希尔排序", ShellSort, false},
		{"归并排序", MergeSort, true},
		{"快速排序", QuickSort, false},
		{"堆排序  ", HeapSort, false},
		{"基数排序", RadixSort, true},
	}

	fmt.Printf("%-10s %-30s 比较次数  交换/移动  稳定性\n", "算法", "排序结果")
	fmt.Println("----------------------------------------------------------------------")

	for _, alg := range algorithms {
		result := alg.fn(arr)
		stableStr := "不稳定"
		if alg.stable {
			stableStr = "稳定"
		}
		fmt.Printf("%-10s %-30v %-8d %-8d %s\n",
			alg.name, result.Data, result.Comparisons, result.Swaps, stableStr)
	}

	fmt.Println("\n408考点总结:")
	fmt.Println("  - 平均时间复杂度 O(nlogn): 归并、快速、堆排序")
	fmt.Println("  - 平均时间复杂度 O(n²): 冒泡、选择、插入")
	fmt.Println("  - 最坏情况快排退化为 O(n²)（已有序时）")
	fmt.Println("  - 空间复杂度：归并O(n)，快排O(logn)，其余O(1)")
	fmt.Println("  - 稳定排序：冒泡、插入、归并、基数")
}
