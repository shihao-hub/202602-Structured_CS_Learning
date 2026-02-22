package basic

import "fmt"

// Node 链表节点
type Node struct {
	data interface{}
	next *Node
}

// LinkedList 链表结构
type LinkedList struct {
	head *Node
	size int
}

// NewLinkedList 创建新链表
func NewLinkedList() *LinkedList {
	return &LinkedList{
		head: nil,
		size: 0,
	}
}

// GetHead 获取头节点
func (ll *LinkedList) GetHead() *Node {
	return ll.head
}

// Size 获取链表大小
func (ll *LinkedList) Size() int {
	return ll.size
}

// IsEmpty 判断链表是否为空
func (ll *LinkedList) IsEmpty() bool {
	return ll.head == nil
}

// InsertAtHead 在头部插入节点
func (ll *LinkedList) InsertAtHead(data interface{}) {
	newNode := &Node{data: data, next: ll.head}
	ll.head = newNode
	ll.size++
}

// InsertAtTail 在尾部插入节点
func (ll *LinkedList) InsertAtTail(data interface{}) {
	newNode := &Node{data: data, next: nil}

	if ll.head == nil {
		ll.head = newNode
	} else {
		current := ll.head
		for current.next != nil {
			current = current.next
		}
		current.next = newNode
	}
	ll.size++
}

// InsertAtIndex 在指定位置插入节点
func (ll *LinkedList) InsertAtIndex(index int, data interface{}) bool {
	if index < 0 || index > ll.size {
		return false
	}

	if index == 0 {
		ll.InsertAtHead(data)
		return true
	}

	if index == ll.size {
		ll.InsertAtTail(data)
		return true
	}

	newNode := &Node{data: data}
	current := ll.head
	for i := 0; i < index-1; i++ {
		current = current.next
	}

	newNode.next = current.next
	current.next = newNode
	ll.size++
	return true
}

// DeleteAtHead 删除头节点
func (ll *LinkedList) DeleteAtHead() interface{} {
	if ll.head == nil {
		return nil
	}

	data := ll.head.data
	ll.head = ll.head.next
	ll.size--
	return data
}

// DeleteAtTail 删除尾节点
func (ll *LinkedList) DeleteAtTail() interface{} {
	if ll.head == nil {
		return nil
	}

	if ll.head.next == nil {
		data := ll.head.data
		ll.head = nil
		ll.size--
		return data
	}

	current := ll.head
	for current.next.next != nil {
		current = current.next
	}

	data := current.next.data
	current.next = nil
	ll.size--
	return data
}

// DeleteAtIndex 删除指定位置的节点
func (ll *LinkedList) DeleteAtIndex(index int) interface{} {
	if index < 0 || index >= ll.size {
		return nil
	}

	if index == 0 {
		return ll.DeleteAtHead()
	}

	if index == ll.size-1 {
		return ll.DeleteAtTail()
	}

	current := ll.head
	for i := 0; i < index-1; i++ {
		current = current.next
	}

	data := current.next.data
	current.next = current.next.next
	ll.size--
	return data
}

// Get 获取指定位置的节点数据
func (ll *LinkedList) Get(index int) interface{} {
	if index < 0 || index >= ll.size {
		return nil
	}

	current := ll.head
	for i := 0; i < index; i++ {
		current = current.next
	}
	return current.data
}

// Set 设置指定位置的节点数据
func (ll *LinkedList) Set(index int, data interface{}) bool {
	if index < 0 || index >= ll.size {
		return false
	}

	current := ll.head
	for i := 0; i < index; i++ {
		current = current.next
	}
	current.data = data
	return true
}

// Find 查找元素，返回索引
func (ll *LinkedList) Find(data interface{}) int {
	current := ll.head
	for i := 0; current != nil; i++ {
		if current.data == data {
			return i
		}
		current = current.next
	}
	return -1
}

// Contains 判断链表是否包含某个元素
func (ll *LinkedList) Contains(data interface{}) bool {
	return ll.Find(data) != -1
}

// Print 打印链表
func (ll *LinkedList) Print() {
	fmt.Printf("LinkedList[%d]: ", ll.size)
	current := ll.head
	for current != nil {
		fmt.Printf("%v", current.data)
		if current.next != nil {
			fmt.Print(" -> ")
		}
		current = current.next
	}
	fmt.Println()
}

// Reverse 反转链表
func (ll *LinkedList) Reverse() {
	if ll.head == nil || ll.head.next == nil {
		return
	}

	var prev *Node
	current := ll.head
	var next *Node

	for current != nil {
		next = current.next
		current.next = prev
		prev = current
		current = next
	}

	ll.head = prev
}

// ToArray 转换为数组
func (ll *LinkedList) ToArray() []interface{} {
	result := make([]interface{}, ll.size)
	current := ll.head
	for i := 0; current != nil; i++ {
		result[i] = current.data
		current = current.next
	}
	return result
}

// Clear 清空链表
func (ll *LinkedList) Clear() {
	ll.head = nil
	ll.size = 0
}

// 示例函数
func LinkedListExample() {
	fmt.Println("=== 链表 (Linked List) 示例 ===")

	// 创建链表
	ll := NewLinkedList()
	fmt.Println("创建空链表:")
	ll.Print()

	// 头部插入
	fmt.Println("\n在头部插入元素 30, 20, 10:")
	ll.InsertAtHead(30)
	ll.InsertAtHead(20)
	ll.InsertAtHead(10)
	ll.Print()

	// 尾部插入
	fmt.Println("\n在尾部插入元素 40, 50:")
	ll.InsertAtTail(40)
	ll.InsertAtTail(50)
	ll.Print()

	// 指定位置插入
	fmt.Println("\n在索引2处插入元素 25:")
	ll.InsertAtIndex(2, 25)
	ll.Print()

	// 查找元素
	fmt.Println("\n查找元素 40 的索引:", ll.Find(40))
	fmt.Println("链表是否包含元素 99:", ll.Contains(99))

	// 获取和设置元素
	fmt.Println("\n获取索引3的元素:", ll.Get(3))
	fmt.Println("将索引1的元素设置为 15:")
	ll.Set(1, 15)
	ll.Print()

	// 删除操作
	fmt.Println("\n删除头节点:")
	ll.DeleteAtHead()
	ll.Print()

	fmt.Println("\n删除尾节点:")
	ll.DeleteAtTail()
	ll.Print()

	fmt.Println("\n删除索引2的节点:")
	ll.DeleteAtIndex(2)
	ll.Print()

	// 反转链表
	fmt.Println("\n反转链表:")
	ll.Reverse()
	ll.Print()

	// 链表信息
	fmt.Printf("\n链表大小: %d, 是否为空: %t\n", ll.Size(), ll.IsEmpty())

	// 转换为数组
	fmt.Println("\n转换为数组:")
	arr := ll.ToArray()
	fmt.Printf("Array: %v\n", arr)

	// 清空链表
	fmt.Println("\n清空链表:")
	ll.Clear()
	ll.Print()
	fmt.Println()
}