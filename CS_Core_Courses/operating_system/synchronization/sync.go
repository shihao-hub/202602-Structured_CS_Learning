package synchronization

import (
	"fmt"
	"sync"
	"time"
)

// Semaphore 信号量实现
type Semaphore struct {
	permits int
	mu      sync.Mutex
	cond    *sync.Cond
}

// NewSemaphore 创建信号量
func NewSemaphore(permits int) *Semaphore {
	s := &Semaphore{
		permits: permits,
	}
	s.cond = sync.NewCond(&s.mu)
	return s
}

// Acquire 获取信号量（P操作）
func (s *Semaphore) Acquire() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for s.permits <= 0 {
		s.cond.Wait()
	}
	s.permits--
}

// Release 释放信号量（V操作）
func (s *Semaphore) Release() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.permits++
	s.cond.Signal()
}

// GetPermits 获取当前可用许可数
func (s *Semaphore) GetPermits() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.permits
}

// Mutex 互斥锁（简单封装演示）
type Mutex struct {
	locked bool
	owner  int
	mu     sync.Mutex
	cond   *sync.Cond
}

// NewMutex 创建互斥锁
func NewMutex() *Mutex {
	m := &Mutex{
		locked: false,
		owner:  -1,
	}
	m.cond = sync.NewCond(&m.mu)
	return m
}

// Lock 加锁
func (m *Mutex) Lock(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for m.locked {
		m.cond.Wait()
	}
	m.locked = true
	m.owner = id
}

// Unlock 解锁
func (m *Mutex) Unlock(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.owner == id {
		m.locked = false
		m.owner = -1
		m.cond.Signal()
	}
}

// IsLocked 检查是否被锁定
func (m *Mutex) IsLocked() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.locked
}

// BoundedBuffer 有界缓冲区（生产者-消费者模型）
type BoundedBuffer struct {
	buffer   []int
	capacity int
	count    int
	in       int
	out      int
	mu       sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
}

// NewBoundedBuffer 创建有界缓冲区
func NewBoundedBuffer(capacity int) *BoundedBuffer {
	bb := &BoundedBuffer{
		buffer:   make([]int, capacity),
		capacity: capacity,
		count:    0,
		in:       0,
		out:      0,
	}
	bb.notFull = sync.NewCond(&bb.mu)
	bb.notEmpty = sync.NewCond(&bb.mu)
	return bb
}

// Put 放入数据（生产者）
func (bb *BoundedBuffer) Put(item int) {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	for bb.count == bb.capacity {
		bb.notFull.Wait()
	}

	bb.buffer[bb.in] = item
	bb.in = (bb.in + 1) % bb.capacity
	bb.count++

	bb.notEmpty.Signal()
}

// Get 取出数据（消费者）
func (bb *BoundedBuffer) Get() int {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	for bb.count == 0 {
		bb.notEmpty.Wait()
	}

	item := bb.buffer[bb.out]
	bb.out = (bb.out + 1) % bb.capacity
	bb.count--

	bb.notFull.Signal()
	return item
}

// Size 获取当前缓冲区大小
func (bb *BoundedBuffer) Size() int {
	bb.mu.Lock()
	defer bb.mu.Unlock()
	return bb.count
}

// ReadWriteLock 读写锁
type ReadWriteLock struct {
	readers   int
	writers   int
	writeWait int
	mu        sync.Mutex
	readCond  *sync.Cond
	writeCond *sync.Cond
}

// NewReadWriteLock 创建读写锁
func NewReadWriteLock() *ReadWriteLock {
	rwl := &ReadWriteLock{
		readers:   0,
		writers:   0,
		writeWait: 0,
	}
	rwl.readCond = sync.NewCond(&rwl.mu)
	rwl.writeCond = sync.NewCond(&rwl.mu)
	return rwl
}

// ReadLock 获取读锁
func (rwl *ReadWriteLock) ReadLock() {
	rwl.mu.Lock()
	defer rwl.mu.Unlock()

	for rwl.writers > 0 || rwl.writeWait > 0 {
		rwl.readCond.Wait()
	}
	rwl.readers++
}

// ReadUnlock 释放读锁
func (rwl *ReadWriteLock) ReadUnlock() {
	rwl.mu.Lock()
	defer rwl.mu.Unlock()

	rwl.readers--
	if rwl.readers == 0 {
		rwl.writeCond.Signal()
	}
}

// WriteLock 获取写锁
func (rwl *ReadWriteLock) WriteLock() {
	rwl.mu.Lock()
	defer rwl.mu.Unlock()

	rwl.writeWait++
	for rwl.readers > 0 || rwl.writers > 0 {
		rwl.writeCond.Wait()
	}
	rwl.writeWait--
	rwl.writers++
}

// WriteUnlock 释放写锁
func (rwl *ReadWriteLock) WriteUnlock() {
	rwl.mu.Lock()
	defer rwl.mu.Unlock()

	rwl.writers--
	rwl.readCond.Broadcast()
	rwl.writeCond.Signal()
}

// SynchronizationExample 同步机制示例
func SynchronizationExample() {
	fmt.Println("=== 进程同步 (Synchronization) 示例 ===")

	// 信号量示例
	fmt.Println("\n1. 信号量示例:")
	semaphore := NewSemaphore(3)
	fmt.Printf("初始信号量: %d\n", semaphore.GetPermits())

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("  线程 %d 尝试获取信号量...\n", id)
			semaphore.Acquire()
			fmt.Printf("  线程 %d 获取信号量成功，当前剩余: %d\n", id, semaphore.GetPermits())
			time.Sleep(100 * time.Millisecond)
			semaphore.Release()
			fmt.Printf("  线程 %d 释放信号量，当前剩余: %d\n", id, semaphore.GetPermits())
		}(i)
	}
	wg.Wait()

	// 生产者-消费者示例
	fmt.Println("\n2. 生产者-消费者示例:")
	buffer := NewBoundedBuffer(5)
	done := make(chan bool)

	// 生产者
	go func() {
		for i := 1; i <= 10; i++ {
			buffer.Put(i)
			fmt.Printf("  生产者: 放入 %d, 缓冲区大小: %d\n", i, buffer.Size())
			time.Sleep(50 * time.Millisecond)
		}
		done <- true
	}()

	// 消费者
	go func() {
		for i := 1; i <= 10; i++ {
			item := buffer.Get()
			fmt.Printf("  消费者: 取出 %d, 缓冲区大小: %d\n", item, buffer.Size())
			time.Sleep(100 * time.Millisecond)
		}
		done <- true
	}()

	<-done
	<-done

	// 读写锁示例
	fmt.Println("\n3. 读写锁示例:")
	rwLock := NewReadWriteLock()
	sharedData := 0

	// 读者
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			rwLock.ReadLock()
			fmt.Printf("  读者 %d: 读取数据 = %d\n", id, sharedData)
			time.Sleep(50 * time.Millisecond)
			rwLock.ReadUnlock()
		}(i)
	}

	// 写者
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(25 * time.Millisecond)
		rwLock.WriteLock()
		sharedData = 100
		fmt.Printf("  写者: 更新数据 = %d\n", sharedData)
		rwLock.WriteUnlock()
	}()

	wg.Wait()

	fmt.Println("\n4. 同步原语说明:")
	fmt.Println("  - 信号量 (Semaphore): 控制同时访问资源的线程数")
	fmt.Println("  - 互斥锁 (Mutex): 保证同一时刻只有一个线程访问临界区")
	fmt.Println("  - 条件变量 (Condition): 线程间的等待/通知机制")
	fmt.Println("  - 读写锁 (RWLock): 允许多读单写的并发访问")
	fmt.Println("  - 有界缓冲区: 生产者-消费者模型的经典实现")

	fmt.Println()
}

// RunAllSynchronizationExamples 运行所有同步相关的示例
func RunAllSynchronizationExamples() {
	SynchronizationExample()
}
