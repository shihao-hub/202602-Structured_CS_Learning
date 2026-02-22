package basic

import "fmt"

// Queue 队列结构 (使用数组实现)
type Queue struct {
	data     []interface{}
	capacity int
	front    int
	rear     int
	size     int
}

// NewQueue 创建新队列
func NewQueue(capacity int) *Queue {
	return &Queue{
		data:     make([]interface{}, capacity),
		capacity: capacity,
		front:    0,
		rear:     -1,
		size:     0,
	}
}

// Enqueue 入队
func (q *Queue) Enqueue(value interface{}) bool {
	if q.IsFull() {
		return false
	}
	q.rear = (q.rear + 1) % q.capacity
	q.data[q.rear] = value
	q.size++
	return true
}

// Dequeue 出队
func (q *Queue) Dequeue() interface{} {
	if q.IsEmpty() {
		return nil
	}
	value := q.data[q.front]
	q.front = (q.front + 1) % q.capacity
	q.size--
	return value
}

// Front 查看队首元素
func (q *Queue) Front() interface{} {
	if q.IsEmpty() {
		return nil
	}
	return q.data[q.front]
}

// Rear 查看队尾元素
func (q *Queue) Rear() interface{} {
	if q.IsEmpty() {
		return nil
	}
	return q.data[q.rear]
}

// IsEmpty 判断队列是否为空
func (q *Queue) IsEmpty() bool {
	return q.size == 0
}

// IsFull 判断队列是否已满
func (q *Queue) IsFull() bool {
	return q.size == q.capacity
}

// Size 获取队列大小
func (q *Queue) Size() int {
	return q.size
}

// Capacity 获取队列容量
func (q *Queue) Capacity() int {
	return q.capacity
}

// Clear 清空队列
func (q *Queue) Clear() {
	q.front = 0
	q.rear = -1
	q.size = 0
}

// Print 打印队列内容
func (q *Queue) Print() {
	fmt.Printf("Queue[%d/%d]: [", q.size, q.capacity)
	if !q.IsEmpty() {
		count := 0
		i := q.front
		for count < q.size {
			if count > 0 {
				fmt.Print(", ")
			}
			fmt.Print(q.data[i])
			i = (i + 1) % q.capacity
			count++
		}
	}
	fmt.Println("]")
}

// Contains 检查队列是否包含某个元素
func (q *Queue) Contains(value interface{}) bool {
	if q.IsEmpty() {
		return false
	}

	count := 0
	i := q.front
	for count < q.size {
		if q.data[i] == value {
			return true
		}
		i = (i + 1) % q.capacity
		count++
	}
	return false
}

// ToArray 转换为数组
func (q *Queue) ToArray() []interface{} {
	result := make([]interface{}, q.size)
	if !q.IsEmpty() {
		count := 0
		i := q.front
		for count < q.size {
			result[count] = q.data[i]
			i = (i + 1) % q.capacity
			count++
		}
	}
	return result
}

// Clone 克隆队列
func (q *Queue) Clone() *Queue {
	newQueue := NewQueue(q.capacity)
	if !q.IsEmpty() {
		count := 0
		i := q.front
		for count < q.size {
			newQueue.Enqueue(q.data[i])
			i = (i + 1) % q.capacity
			count++
		}
	}
	return newQueue
}

// 示例函数
func QueueExample() {
	fmt.Println("=== 队列 (Queue) 示例 ===")

	// 创建容量为5的队列
	queue := NewQueue(5)
	fmt.Println("创建容量为5的空队列:")
	queue.Print()

	// 入队操作
	fmt.Println("\n入队操作:")
	fmt.Println("Enqueue 10:", queue.Enqueue(10))
	queue.Print()

	fmt.Println("Enqueue 20:", queue.Enqueue(20))
	queue.Print()

	fmt.Println("Enqueue 30:", queue.Enqueue(30))
	queue.Print()

	fmt.Println("Enqueue 40:", queue.Enqueue(40))
	queue.Print()

	fmt.Println("Enqueue 50:", queue.Enqueue(50))
	queue.Print()

	fmt.Println("Enqueue 60 (队列满):", queue.Enqueue(60))
	queue.Print()

	// 查看队首和队尾元素
	fmt.Printf("\n队首元素 (Front): %v\n", queue.Front())
	fmt.Printf("队尾元素 (Rear): %v\n", queue.Rear())

	// 出队操作
	fmt.Println("\n出队操作:")
	value := queue.Dequeue()
	fmt.Printf("Dequeue: %v, ", value)
	queue.Print()

	value = queue.Dequeue()
	fmt.Printf("Dequeue: %v, ", value)
	queue.Print()

	// 队列信息
	fmt.Printf("\n队列大小: %d, 容量: %d, 是否为空: %t, 是否已满: %t\n",
		queue.Size(), queue.Capacity(), queue.IsEmpty(), queue.IsFull())

	// 检查元素是否存在
	fmt.Println("\n检查元素是否存在:")
	fmt.Println("Contains 30:", queue.Contains(30))
	fmt.Println("Contains 99:", queue.Contains(99))

	// 转换为数组
	fmt.Println("\n转换为数组:")
	arr := queue.ToArray()
	fmt.Printf("Array: %v\n", arr)

	// 克隆队列
	fmt.Println("\n克隆队列:")
	clonedQueue := queue.Clone()
	clonedQueue.Print()

	// 验证队列的FIFO特性
	fmt.Println("\n验证队列的FIFO特性:")
	fifoQueue := NewQueue(3)
	items := []interface{}{"First", "Second", "Third"}

	fmt.Println("入队顺序:")
	for _, item := range items {
		fifoQueue.Enqueue(item)
		fmt.Printf("Enqueue: %v\n", item)
	}

	fmt.Println("\n出队顺序:")
	for !fifoQueue.IsEmpty() {
		value := fifoQueue.Dequeue()
		fmt.Printf("Dequeue: %v\n", value)
	}

	// 清空队列
	fmt.Println("\n清空队列:")
	queue.Clear()
	queue.Print()
	fmt.Println()
}