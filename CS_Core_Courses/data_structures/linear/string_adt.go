package linear

import "fmt"

// ============================================================
// 串 (String ADT)
// 408考点：串的存储结构、基本操作、模式匹配（KMP见algorithm包）
// ============================================================

// MyString 串的顺序存储结构
type MyString struct {
	data   []byte
	length int
}

// NewMyString 创建串
func NewMyString(s string) *MyString {
	data := []byte(s)
	return &MyString{data: data, length: len(data)}
}

// Length 返回串长
func (s *MyString) Length() int {
	return s.length
}

// CharAt 返回指定位置的字符（从0开始）
func (s *MyString) CharAt(index int) (byte, bool) {
	if index < 0 || index >= s.length {
		return 0, false
	}
	return s.data[index], true
}

// SubString 截取子串 [start, start+length)
func (s *MyString) SubString(start, length int) *MyString {
	if start < 0 || start >= s.length || length <= 0 {
		return NewMyString("")
	}
	end := start + length
	if end > s.length {
		end = s.length
	}
	return NewMyString(string(s.data[start:end]))
}

// Concat 串连接
func (s *MyString) Concat(other *MyString) *MyString {
	newData := make([]byte, s.length+other.length)
	copy(newData, s.data)
	copy(newData[s.length:], other.data)
	return &MyString{data: newData, length: s.length + other.length}
}

// Index 定位子串（朴素匹配），返回首次出现位置，-1表示不存在
func (s *MyString) Index(sub *MyString) int {
	if sub.length == 0 || sub.length > s.length {
		return -1
	}
	for i := 0; i <= s.length-sub.length; i++ {
		match := true
		for j := 0; j < sub.length; j++ {
			if s.data[i+j] != sub.data[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

// Replace 替换子串（替换第一个匹配）
func (s *MyString) Replace(old, new_ *MyString) *MyString {
	idx := s.Index(old)
	if idx == -1 {
		return NewMyString(string(s.data))
	}
	result := make([]byte, 0, s.length-old.length+new_.length)
	result = append(result, s.data[:idx]...)
	result = append(result, new_.data...)
	result = append(result, s.data[idx+old.length:]...)
	return &MyString{data: result, length: len(result)}
}

// Insert 在指定位置插入串
func (s *MyString) Insert(pos int, sub *MyString) *MyString {
	if pos < 0 || pos > s.length {
		return NewMyString(string(s.data))
	}
	result := make([]byte, 0, s.length+sub.length)
	result = append(result, s.data[:pos]...)
	result = append(result, sub.data...)
	result = append(result, s.data[pos:]...)
	return &MyString{data: result, length: len(result)}
}

// Delete 删除子串 [start, start+length)
func (s *MyString) Delete(start, length int) *MyString {
	if start < 0 || start >= s.length || length <= 0 {
		return NewMyString(string(s.data))
	}
	end := start + length
	if end > s.length {
		end = s.length
	}
	result := make([]byte, 0, s.length-(end-start))
	result = append(result, s.data[:start]...)
	result = append(result, s.data[end:]...)
	return &MyString{data: result, length: len(result)}
}

// Compare 串比较: 返回 <0, 0, >0
func (s *MyString) Compare(other *MyString) int {
	minLen := s.length
	if other.length < minLen {
		minLen = other.length
	}
	for i := 0; i < minLen; i++ {
		if s.data[i] != other.data[i] {
			return int(s.data[i]) - int(other.data[i])
		}
	}
	return s.length - other.length
}

// String 转为Go字符串
func (s *MyString) String() string {
	return string(s.data)
}

// StringADTExample 串操作示例
func StringADTExample() {
	fmt.Println("\n--- 串 (String ADT) ---")

	s1 := NewMyString("Hello, World!")
	fmt.Printf("串 s1: \"%s\" (长度=%d)\n", s1, s1.Length())

	// 子串
	sub := s1.SubString(7, 5)
	fmt.Printf("子串 s1[7:12]: \"%s\"\n", sub)

	// 连接
	s2 := NewMyString(" Welcome!")
	s3 := s1.Concat(s2)
	fmt.Printf("连接: \"%s\"\n", s3)

	// 定位
	target := NewMyString("World")
	idx := s1.Index(target)
	fmt.Printf("定位 \"World\": 位置=%d\n", idx)

	// 替换
	old := NewMyString("World")
	new_ := NewMyString("Go")
	s4 := s1.Replace(old, new_)
	fmt.Printf("替换 World→Go: \"%s\"\n", s4)

	// 插入
	ins := NewMyString("Beautiful ")
	s5 := s1.Insert(7, ins)
	fmt.Printf("在位置7插入: \"%s\"\n", s5)

	// 删除
	s6 := s1.Delete(5, 7)
	fmt.Printf("删除 s1[5:12]: \"%s\"\n", s6)

	// 比较
	sa := NewMyString("abc")
	sb := NewMyString("abd")
	fmt.Printf("比较 \"%s\" 与 \"%s\": %d\n", sa, sb, sa.Compare(sb))
}
