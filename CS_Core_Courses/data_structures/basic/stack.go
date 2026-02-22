package basic

import "fmt"

// Stack 栈结构 (使用数组实现)
type Stack struct {
	data     []interface{}
	capacity int
	top      int
}

// NewStack 创建新栈
func NewStack(capacity int) *Stack {
	return &Stack{
		data:     make([]interface{}, capacity),
		capacity: capacity,
		top:      -1,
	}
}

// Push 入栈
func (s *Stack) Push(value interface{}) bool {
	if s.IsFull() {
		return false
	}
	s.top++
	s.data[s.top] = value
	return true
}

// Pop 出栈
func (s *Stack) Pop() interface{} {
	if s.IsEmpty() {
		return nil
	}
	value := s.data[s.top]
	s.top--
	return value
}

// Peek 查看栈顶元素
func (s *Stack) Peek() interface{} {
	if s.IsEmpty() {
		return nil
	}
	return s.data[s.top]
}

// IsEmpty 判断栈是否为空
func (s *Stack) IsEmpty() bool {
	return s.top == -1
}

// IsFull 判断栈是否已满
func (s *Stack) IsFull() bool {
	return s.top == s.capacity-1
}

// Size 获取栈大小
func (s *Stack) Size() int {
	return s.top + 1
}

// Capacity 获取栈容量
func (s *Stack) Capacity() int {
	return s.capacity
}

// Clear 清空栈
func (s *Stack) Clear() {
	s.top = -1
}

// Print 打印栈内容
func (s *Stack) Print() {
	fmt.Printf("Stack[%d/%d]: [", s.Size(), s.capacity)
	for i := 0; i <= s.top; i++ {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(s.data[i])
	}
	fmt.Println("]")
}

// ReverseStack 反转栈 (不使用额外空间)
func (s *Stack) ReverseStack() {
	if s.IsEmpty() || s.Size() == 1 {
		return
	}

	// 递归反转
	s.reverseStackRecursive(0, s.Size()-1)
}

func (s *Stack) reverseStackRecursive(start, end int) {
	if start >= end {
		return
	}
	// 交换元素
	s.data[start], s.data[end] = s.data[end], s.data[start]
	s.reverseStackRecursive(start+1, end-1)
}

// ToArray 转换为数组
func (s *Stack) ToArray() []interface{} {
	result := make([]interface{}, s.Size())
	for i := 0; i <= s.top; i++ {
		result[i] = s.data[i]
	}
	return result
}

// 示例函数
func StackExample() {
	fmt.Println("=== 栈 (Stack) 示例 ===")

	// 创建容量为5的栈
	stack := NewStack(5)
	fmt.Println("创建容量为5的空栈:")
	stack.Print()

	// 入栈操作
	fmt.Println("\n入栈操作:")
	fmt.Println("Push 10:", stack.Push(10))
	stack.Print()

	fmt.Println("Push 20:", stack.Push(20))
	stack.Print()

	fmt.Println("Push 30:", stack.Push(30))
	stack.Print()

	fmt.Println("Push 40:", stack.Push(40))
	stack.Print()

	fmt.Println("Push 50:", stack.Push(50))
	stack.Print()

	fmt.Println("Push 60 (栈满):", stack.Push(60))
	stack.Print()

	// 查看栈顶元素
	fmt.Println("\n栈顶元素 (Peek):", stack.Peek())

	// 出栈操作
	fmt.Println("\n出栈操作:")
	value := stack.Pop()
	fmt.Printf("Pop: %v, ", value)
	stack.Print()

	value = stack.Pop()
	fmt.Printf("Pop: %v, ", value)
	stack.Print()

	// 栈信息
	fmt.Printf("\n栈大小: %d, 容量: %d, 是否为空: %t, 是否已满: %t\n",
		stack.Size(), stack.Capacity(), stack.IsEmpty(), stack.IsFull())

	// 反转栈
	fmt.Println("\n反转栈:")
	stack.ReverseStack()
	stack.Print()

	// 转换为数组
	fmt.Println("\n转换为数组:")
	arr := stack.ToArray()
	fmt.Printf("Array: %v\n", arr)

	// 清空栈
	fmt.Println("\n清空栈:")
	stack.Clear()
	stack.Print()

	// 验证栈的LIFO特性
	fmt.Println("\n验证栈的LIFO特性:")
	lifoStack := NewStack(3)
	items := []interface{}{"First", "Second", "Third"}

	fmt.Println("入栈顺序:")
	for _, item := range items {
		lifoStack.Push(item)
		fmt.Printf("Push: %v\n", item)
	}

	fmt.Println("\n出栈顺序:")
	for !lifoStack.IsEmpty() {
		value := lifoStack.Pop()
		fmt.Printf("Pop: %v\n", value)
	}
	fmt.Println()
}