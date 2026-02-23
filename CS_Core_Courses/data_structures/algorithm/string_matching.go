package algorithm

import "fmt"

// ============================================================
// 字符串匹配算法
// 408考点：KMP算法、next数组构造（重点）
// ============================================================

// --- 1. 朴素模式匹配 (Brute Force) ---
// 时间复杂度: 最坏 O(m*n)
func BruteForceMatch(text, pattern string) (int, int) {
	n, m := len(text), len(pattern)
	comparisons := 0

	for i := 0; i <= n-m; i++ {
		j := 0
		for j < m {
			comparisons++
			if text[i+j] != pattern[j] {
				break
			}
			j++
		}
		if j == m {
			return i, comparisons // 匹配成功
		}
	}
	return -1, comparisons // 未找到
}

// --- 2. KMP 算法 ---
// 时间复杂度: O(m+n)
// 核心思想：利用已匹配信息，避免主串回溯

// BuildNext 构造next数组
// next[j] 表示 pattern[0..j-1] 的最长相等前后缀长度
//
// 构造过程详解：
//
//	next[0] = -1 （特殊标记）
//	next[1] = 0  （长度为1的串，无真前后缀）
//	对于 j >= 2:
//	  比较 pattern[j-1] 和 pattern[next[j-1]]
//	  若相等: next[j] = next[j-1] + 1
//	  若不等: 回退到 next[next[j-1]]，继续比较
func BuildNext(pattern string) []int {
	m := len(pattern)
	next := make([]int, m)
	next[0] = -1
	if m == 1 {
		return next
	}
	next[1] = 0

	// i 指向当前要计算next值的位置
	// j 指向当前比较的前缀末尾位置
	i := 2
	j := 0
	for i < m {
		if pattern[i-1] == pattern[j] {
			j++
			next[i] = j
			i++
		} else if j > 0 {
			j = next[j] // 回退
		} else {
			next[i] = 0
			i++
		}
	}
	return next
}

// BuildNextImproved 构造改进的nextval数组
// 优化：当 pattern[j] == pattern[next[j]] 时，继续回退
// 避免不必要的比较
func BuildNextImproved(pattern string) []int {
	m := len(pattern)
	next := BuildNext(pattern)
	nextval := make([]int, m)
	nextval[0] = -1

	for j := 1; j < m; j++ {
		if next[j] != -1 && pattern[j] == pattern[next[j]] {
			nextval[j] = nextval[next[j]] // 继续回退
		} else {
			nextval[j] = next[j]
		}
	}
	return nextval
}

// KMPSearch KMP字符串匹配
func KMPSearch(text, pattern string) (int, int) {
	n, m := len(text), len(pattern)
	if m == 0 {
		return 0, 0
	}

	next := BuildNext(pattern)
	comparisons := 0
	i, j := 0, 0

	for i < n {
		comparisons++
		if text[i] == pattern[j] {
			i++
			j++
			if j == m {
				return i - m, comparisons // 匹配成功
			}
		} else if j > 0 {
			j = next[j] // 主串不回溯，模式串回退
		} else {
			i++ // 模式串无法回退，主串前进
		}
	}
	return -1, comparisons
}

// KMPSearchAll KMP查找所有匹配位置
func KMPSearchAll(text, pattern string) []int {
	n, m := len(text), len(pattern)
	if m == 0 {
		return nil
	}

	next := BuildNext(pattern)
	positions := make([]int, 0)
	i, j := 0, 0

	for i < n {
		if text[i] == pattern[j] {
			i++
			j++
			if j == m {
				positions = append(positions, i-m)
				j = next[j-1] // 继续查找下一个匹配
				if j < 0 {
					j = 0
				}
			}
		} else if j > 0 {
			j = next[j]
		} else {
			i++
		}
	}
	return positions
}

// StringMatchingExample 字符串匹配示例
func StringMatchingExample() {
	fmt.Println("\n--- 字符串匹配算法 ---")

	text := "ABABDABACDABABCABAB"
	pattern := "ABABCABAB"
	fmt.Printf("主串:   \"%s\"\n", text)
	fmt.Printf("模式串: \"%s\"\n\n", pattern)

	// 1. 朴素匹配
	pos, comp := BruteForceMatch(text, pattern)
	fmt.Printf("朴素匹配: 位置=%d, 比较次数=%d\n", pos, comp)

	// 2. KMP匹配
	pos, comp = KMPSearch(text, pattern)
	fmt.Printf("KMP匹配:  位置=%d, 比较次数=%d\n", pos, comp)

	// 3. next数组详解
	fmt.Println("\n--- next数组构造详解 ---")
	patterns := []string{"ABABCABAB", "ABAABCAC", "ABCDABD"}
	for _, p := range patterns {
		next := BuildNext(p)
		nextval := BuildNextImproved(p)
		fmt.Printf("\n模式串: \"%s\"\n", p)
		fmt.Printf("  下标:    ")
		for i := range p {
			fmt.Printf("%3d", i)
		}
		fmt.Println()
		fmt.Printf("  字符:    ")
		for _, c := range p {
			fmt.Printf("%3c", c)
		}
		fmt.Println()
		fmt.Printf("  next:    ")
		for _, v := range next {
			fmt.Printf("%3d", v)
		}
		fmt.Println()
		fmt.Printf("  nextval: ")
		for _, v := range nextval {
			fmt.Printf("%3d", v)
		}
		fmt.Println()
	}

	// 4. 查找所有匹配
	fmt.Println("\n--- 查找所有匹配位置 ---")
	text2 := "AABAABAABAAB"
	pattern2 := "AABAA"
	fmt.Printf("主串: \"%s\", 模式串: \"%s\"\n", text2, pattern2)
	positions := KMPSearchAll(text2, pattern2)
	fmt.Printf("所有匹配位置: %v\n", positions)

	fmt.Println("\n408考点总结:")
	fmt.Println("  - KMP的核心是next数组（最长相等前后缀）")
	fmt.Println("  - next数组手算是必考题型")
	fmt.Println("  - KMP主串指针不回溯，时间复杂度 O(m+n)")
	fmt.Println("  - nextval是next的优化版，避免多余比较")
}
