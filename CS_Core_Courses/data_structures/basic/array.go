package basic

import "fmt"

// Array 基础数组结构
type Array struct {
	data     []interface{}
	capacity int
	size     int
}

// NewArray 创建新数组
func NewArray(capacity int) *Array {
	return &Array{
		data:     make([]interface{}, capacity),
		capacity: capacity,
		size:     0,
	}
}

// Get 获取指定索引的元素
func (a *Array) Get(index int) interface{} {
	if index < 0 || index >= a.size {
		return nil
	}
	return a.data[index]
}

// Set 设置指定索引的元素
func (a *Array) Set(index int, value interface{}) bool {
	if index < 0 || index >= a.size {
		return false
	}
	a.data[index] = value
	return true
}

// Insert 在指定位置插入元素
func (a *Array) Insert(index int, value interface{}) bool {
	if index < 0 || index > a.size || a.size >= a.capacity {
		return false
	}

	// 将index及之后的元素后移
	for i := a.size; i > index; i-- {
		a.data[i] = a.data[i-1]
	}

	a.data[index] = value
	a.size++
	return true
}

// Delete 删除指定位置的元素
func (a *Array) Delete(index int) bool {
	if index < 0 || index >= a.size {
		return false
	}

	// 将index之后的元素前移
	for i := index; i < a.size-1; i++ {
		a.data[i] = a.data[i+1]
	}

	a.size--
	return true
}

// Append 在末尾添加元素
func (a *Array) Append(value interface{}) bool {
	if a.size >= a.capacity {
		return false
	}
	a.data[a.size] = value
	a.size++
	return true
}

// Size 获取数组大小
func (a *Array) Size() int {
	return a.size
}

// Capacity 获取数组容量
func (a *Array) Capacity() int {
	return a.capacity
}

// IsEmpty 判断数组是否为空
func (a *Array) IsEmpty() bool {
	return a.size == 0
}

// Print 打印数组内容
func (a *Array) Print() {
	fmt.Printf("Array[%d/%d]: [", a.size, a.capacity)
	for i := 0; i < a.size; i++ {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(a.data[i])
	}
	fmt.Println("]")
}

// Find 查找元素索引，返回-1表示未找到
func (a *Array) Find(value interface{}) int {
	for i := 0; i < a.size; i++ {
		if a.data[i] == value {
			return i
		}
	}
	return -1
}

// Reverse 反转数组
func (a *Array) Reverse() {
	for i, j := 0, a.size-1; i < j; i, j = i+1, j-1 {
		a.data[i], a.data[j] = a.data[j], a.data[i]
	}
}

// Clone 克隆数组
func (a *Array) Clone() *Array {
	newArray := NewArray(a.capacity)
	for i := 0; i < a.size; i++ {
		newArray.data[i] = a.data[i]
	}
	newArray.size = a.size
	return newArray
}

// 示例函数
func ArrayExample() {
	fmt.Println("=== 数组 (Array) 示例 ===")

	// 创建容量为5的数组
	arr := NewArray(5)
	fmt.Println("创建容量为5的空数组:")
	arr.Print()

	// 添加元素
	fmt.Println("\n添加元素 10, 20, 30:")
	arr.Append(10)
	arr.Append(20)
	arr.Append(30)
	arr.Print()

	// 插入元素
	fmt.Println("\n在索引1处插入元素 15:")
	arr.Insert(1, 15)
	arr.Print()

	// 查找元素
	fmt.Println("\n查找元素 20 的索引:", arr.Find(20))
	fmt.Println("查找元素 99 的索引:", arr.Find(99))

	// 删除元素
	fmt.Println("\n删除索引2的元素:")
	arr.Delete(2)
	arr.Print()

	// 修改元素
	fmt.Println("\n将索引1的元素修改为 25:")
	arr.Set(1, 25)
	arr.Print()

	// 反转数组
	fmt.Println("\n反转数组:")
	arr.Reverse()
	arr.Print()

	// 数组信息
	fmt.Printf("\n数组大小: %d, 容量: %d, 是否为空: %t\n",
		arr.Size(), arr.Capacity(), arr.IsEmpty())

	// 克隆数组
	fmt.Println("\n克隆数组:")
	clonedArr := arr.Clone()
	clonedArr.Print()
	fmt.Println()
}