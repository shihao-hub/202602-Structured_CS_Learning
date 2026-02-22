package basic

import (
	"fmt"
)

// HashEntry 哈希表条目
type HashEntry struct {
	key   interface{}
	value interface{}
	next  *HashEntry
}

// HashTable 哈希表结构
type HashTable struct {
	buckets []*HashEntry
	size    int
	capacity int
}

// NewHashTable 创建新哈希表
func NewHashTable(capacity int) *HashTable {
	return &HashTable{
		buckets: make([]*HashEntry, capacity),
		size:    0,
		capacity: capacity,
	}
}

// hash 哈希函数
func (ht *HashTable) hash(key interface{}) int {
	var hash int
	switch v := key.(type) {
	case int:
		hash = v
	case string:
		hash = 0
		for _, char := range v {
			hash = hash*31 + int(char)
		}
	default:
		// 使用默认的Go哈希函数
		hash = 0
		keyStr := fmt.Sprintf("%v", key)
		for _, char := range keyStr {
			hash = hash*31 + int(char)
		}
	}
	return (hash & 0x7fffffff) % ht.capacity
}

// Put 插入键值对
func (ht *HashTable) Put(key, value interface{}) bool {
	if key == nil {
		return false
	}

	index := ht.hash(key)

	// 检查key是否已存在
	entry := ht.buckets[index]
	for entry != nil {
		if entry.key == key {
			entry.value = value
			return true
		}
		entry = entry.next
	}

	// 创建新条目
	newEntry := &HashEntry{key: key, value: value}
	newEntry.next = ht.buckets[index]
	ht.buckets[index] = newEntry
	ht.size++
	return true
}

// Get 获取值
func (ht *HashTable) Get(key interface{}) interface{} {
	if key == nil {
		return nil
	}

	index := ht.hash(key)
	entry := ht.buckets[index]

	for entry != nil {
		if entry.key == key {
			return entry.value
		}
		entry = entry.next
	}

	return nil
}

// Remove 删除键值对
func (ht *HashTable) Remove(key interface{}) interface{} {
	if key == nil {
		return nil
	}

	index := ht.hash(key)
	entry := ht.buckets[index]
	var prev *HashEntry

	for entry != nil {
		if entry.key == key {
			if prev == nil {
				ht.buckets[index] = entry.next
			} else {
				prev.next = entry.next
			}
			ht.size--
			return entry.value
		}
		prev = entry
		entry = entry.next
	}

	return nil
}

// Contains 检查是否包含指定键
func (ht *HashTable) Contains(key interface{}) bool {
	return ht.Get(key) != nil
}

// Size 获取哈希表大小
func (ht *HashTable) Size() int {
	return ht.size
}

// Capacity 获取哈希表容量
func (ht *HashTable) Capacity() int {
	return ht.capacity
}

// IsEmpty 判断哈希表是否为空
func (ht *HashTable) IsEmpty() bool {
	return ht.size == 0
}

// Clear 清空哈希表
func (ht *HashTable) Clear() {
	ht.buckets = make([]*HashEntry, ht.capacity)
	ht.size = 0
}

// Keys 获取所有键
func (ht *HashTable) Keys() []interface{} {
	keys := make([]interface{}, 0, ht.size)
	for _, bucket := range ht.buckets {
		entry := bucket
		for entry != nil {
			keys = append(keys, entry.key)
			entry = entry.next
		}
	}
	return keys
}

// Values 获取所有值
func (ht *HashTable) Values() []interface{} {
	values := make([]interface{}, 0, ht.size)
	for _, bucket := range ht.buckets {
		entry := bucket
		for entry != nil {
			values = append(values, entry.value)
			entry = entry.next
		}
	}
	return values
}

// Print 打印哈希表
func (ht *HashTable) Print() {
	fmt.Printf("HashTable[%d/%d]:\n", ht.size, ht.capacity)
	for i, bucket := range ht.buckets {
		if bucket != nil {
			fmt.Printf("  [%d]: ", i)
			entry := bucket
			for entry != nil {
				fmt.Printf("%v=%v", entry.key, entry.value)
				if entry.next != nil {
					fmt.Print(" -> ")
				}
				entry = entry.next
			}
			fmt.Println()
		}
	}
}

// getLoadFactor 获取负载因子
func (ht *HashTable) getLoadFactor() float64 {
	return float64(ht.size) / float64(ht.capacity)
}

// 示例函数
func HashTableExample() {
	fmt.Println("=== 哈希表 (Hash Table) 示例 ===")

	// 创建容量为10的哈希表
	hashTable := NewHashTable(10)
	fmt.Println("创建容量为10的空哈希表:")
	hashTable.Print()

	// 插入键值对
	fmt.Println("\n插入键值对:")
	pairs := map[interface{}]interface{}{
		"name":     "Alice",
		"age":      25,
		"city":     "Beijing",
		1:          "One",
		2:          "Two",
		"hello":    "World",
		"email":    "alice@example.com",
	}

	for key, value := range pairs {
		success := hashTable.Put(key, value)
		fmt.Printf("Put(%v, %v): %t\n", key, value, success)
	}
	hashTable.Print()

	// 获取值
	fmt.Println("\n获取值:")
	keys := []interface{}{"name", "age", "nonexistent", 1, "hello"}
	for _, key := range keys {
		value := hashTable.Get(key)
		fmt.Printf("Get(%v): %v\n", key, value)
	}

	// 检查键是否存在
	fmt.Println("\n检查键是否存在:")
	for _, key := range keys {
		exists := hashTable.Contains(key)
		fmt.Printf("Contains(%v): %t\n", key, exists)
	}

	// 更新值
	fmt.Println("\n更新值:")
	fmt.Println("更新前 age =", hashTable.Get("age"))
	hashTable.Put("age", 26)
	fmt.Println("更新后 age =", hashTable.Get("age"))

	// 删除键值对
	fmt.Println("\n删除键值对:")
	keysToRemove := []interface{}{"city", "nonexistent"}
	for _, key := range keysToRemove {
		value := hashTable.Remove(key)
		fmt.Printf("Remove(%v): %v\n", key, value)
	}
	hashTable.Print()

	// 哈希表信息
	fmt.Printf("\n哈希表大小: %d, 容量: %d, 是否为空: %t, 负载因子: %.2f\n",
		hashTable.Size(), hashTable.Capacity(), hashTable.IsEmpty(), hashTable.getLoadFactor())

	// 获取所有键和值
	fmt.Println("\n获取所有键和值:")
	keys = hashTable.Keys()
	values := hashTable.Values()
	fmt.Printf("Keys: %v\n", keys)
	fmt.Printf("Values: %v\n", values)

	// 清空哈希表
	fmt.Println("\n清空哈希表:")
	hashTable.Clear()
	hashTable.Print()

	// 哈希冲突示例
	fmt.Println("\n=== 哈希冲突示例 ===")
	smallHashTable := NewHashTable(3) // 小容量容易产生冲突
	conflictKeys := []interface{}{"a", "d", "g", "j", "m"} // 这些键可能产生相同哈希值

	fmt.Println("在小容量哈希表中插入容易冲突的键:")
	for i, key := range conflictKeys {
		value := fmt.Sprintf("value%d", i)
		smallHashTable.Put(key, value)
		fmt.Printf("Put(%v, %v)\n", key, value)
	}
	smallHashTable.Print()
	fmt.Println()
}