package tree

import "fmt"

// Heap 堆结构（最小堆/最大堆）
type Heap struct {
	data  []int
	isMin bool // true为最小堆，false为最大堆
}

// NewMinHeap 创建最小堆
func NewMinHeap() *Heap {
	return &Heap{
		data:  make([]int, 0),
		isMin: true,
	}
}

// NewMaxHeap 创建最大堆
func NewMaxHeap() *Heap {
	return &Heap{
		data:  make([]int, 0),
		isMin: false,
	}
}

// Size 获取堆大小
func (h *Heap) Size() int {
	return len(h.data)
}

// IsEmpty 检查堆是否为空
func (h *Heap) IsEmpty() bool {
	return len(h.data) == 0
}

// Parent 获取父节点索引
func (h *Heap) parent(i int) int {
	return (i - 1) / 2
}

// LeftChild 获取左子节点索引
func (h *Heap) leftChild(i int) int {
	return 2*i + 1
}

// RightChild 获取右子节点索引
func (h *Heap) rightChild(i int) int {
	return 2*i + 2
}

// compare 比较两个元素（根据堆类型）
func (h *Heap) compare(i, j int) bool {
	if h.isMin {
		return h.data[i] < h.data[j]
	}
	return h.data[i] > h.data[j]
}

// swap 交换两个元素
func (h *Heap) swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

// Insert 插入元素
func (h *Heap) Insert(value int) {
	h.data = append(h.data, value)
	h.siftUp(len(h.data) - 1)
}

// siftUp 上浮操作
func (h *Heap) siftUp(i int) {
	for i > 0 && h.compare(i, h.parent(i)) {
		h.swap(i, h.parent(i))
		i = h.parent(i)
	}
}

// Peek 查看堆顶元素
func (h *Heap) Peek() (int, bool) {
	if h.IsEmpty() {
		return 0, false
	}
	return h.data[0], true
}

// Extract 提取堆顶元素
func (h *Heap) Extract() (int, bool) {
	if h.IsEmpty() {
		return 0, false
	}

	top := h.data[0]
	lastIdx := len(h.data) - 1
	h.data[0] = h.data[lastIdx]
	h.data = h.data[:lastIdx]

	if len(h.data) > 0 {
		h.siftDown(0)
	}

	return top, true
}

// siftDown 下沉操作
func (h *Heap) siftDown(i int) {
	for {
		smallest := i
		left := h.leftChild(i)
		right := h.rightChild(i)

		if left < len(h.data) && h.compare(left, smallest) {
			smallest = left
		}
		if right < len(h.data) && h.compare(right, smallest) {
			smallest = right
		}

		if smallest == i {
			break
		}

		h.swap(i, smallest)
		i = smallest
	}
}

// BuildHeap 从数组构建堆
func (h *Heap) BuildHeap(arr []int) {
	h.data = make([]int, len(arr))
	copy(h.data, arr)

	// 从最后一个非叶子节点开始下沉
	for i := len(h.data)/2 - 1; i >= 0; i-- {
		h.siftDown(i)
	}
}

// ToArray 转换为数组
func (h *Heap) ToArray() []int {
	result := make([]int, len(h.data))
	copy(result, h.data)
	return result
}

// Print 打印堆
func (h *Heap) Print() {
	heapType := "最小堆"
	if !h.isMin {
		heapType = "最大堆"
	}
	fmt.Printf("%s (size=%d): %v\n", heapType, len(h.data), h.data)
}

// HeapSort 堆排序（返回排序后的数组）
func HeapSort(arr []int, ascending bool) []int {
	var heap *Heap
	if ascending {
		heap = NewMinHeap()
	} else {
		heap = NewMaxHeap()
	}

	heap.BuildHeap(arr)

	result := make([]int, 0, len(arr))
	for !heap.IsEmpty() {
		val, _ := heap.Extract()
		result = append(result, val)
	}

	return result
}

// HeapExample 堆示例
func HeapExample() {
	fmt.Println("=== 堆 (Heap) 示例 ===")

	// 最小堆示例
	fmt.Println("\n1. 最小堆:")
	minHeap := NewMinHeap()

	values := []int{35, 10, 25, 5, 20, 15, 30}
	fmt.Printf("插入顺序: %v\n", values)
	for _, v := range values {
		minHeap.Insert(v)
		fmt.Printf("插入 %d: ", v)
		minHeap.Print()
	}

	fmt.Println("\n提取元素:")
	for !minHeap.IsEmpty() {
		val, _ := minHeap.Extract()
		fmt.Printf("提取 %d, 剩余: %v\n", val, minHeap.ToArray())
	}

	// 最大堆示例
	fmt.Println("\n2. 最大堆:")
	maxHeap := NewMaxHeap()

	for _, v := range values {
		maxHeap.Insert(v)
	}
	fmt.Printf("最大堆内容: ")
	maxHeap.Print()

	fmt.Println("\n提取元素:")
	for !maxHeap.IsEmpty() {
		val, _ := maxHeap.Extract()
		fmt.Printf("提取 %d, 剩余: %v\n", val, maxHeap.ToArray())
	}

	// 堆排序
	fmt.Println("\n3. 堆排序:")
	unsorted := []int{64, 34, 25, 12, 22, 11, 90}
	fmt.Printf("原数组: %v\n", unsorted)
	fmt.Printf("升序排序: %v\n", HeapSort(unsorted, true))
	fmt.Printf("降序排序: %v\n", HeapSort(unsorted, false))

	// 从数组构建堆
	fmt.Println("\n4. 从数组构建堆:")
	arr := []int{50, 30, 40, 10, 20, 35, 45}
	buildHeap := NewMinHeap()
	buildHeap.BuildHeap(arr)
	fmt.Printf("原数组: %v\n", arr)
	fmt.Printf("构建最小堆: %v\n", buildHeap.ToArray())

	fmt.Println()
}
