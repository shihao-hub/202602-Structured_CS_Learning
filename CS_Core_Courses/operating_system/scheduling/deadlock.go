package scheduling

import (
	"fmt"
)

// BankerSystem 银行家算法系统
type BankerSystem struct {
	NumProcesses int     // 进程数量
	NumResources int     // 资源类型数量
	Max          [][]int // 最大需求矩阵 [进程][资源类型]
	Allocation   [][]int // 已分配矩阵
	Need         [][]int // 需求矩阵 (Max - Allocation)
	Available    []int   // 可用资源向量
	Finish       []bool  // 进程是否完成
	SafeSequence []int   // 安全序列
}

// NewBankerSystem 创建银行家算法系统
func NewBankerSystem(numProcesses, numResources int, max, allocation [][]int, available []int) *BankerSystem {
	bs := &BankerSystem{
		NumProcesses: numProcesses,
		NumResources: numResources,
		Max:          copyMatrix(max),
		Allocation:   copyMatrix(allocation),
		Need:         make([][]int, numProcesses),
		Available:    make([]int, numResources),
		Finish:       make([]bool, numProcesses),
		SafeSequence: make([]int, 0),
	}

	// 复制可用资源
	copy(bs.Available, available)

	// 计算需求矩阵 Need = Max - Allocation
	for i := 0; i < numProcesses; i++ {
		bs.Need[i] = make([]int, numResources)
		for j := 0; j < numResources; j++ {
			bs.Need[i][j] = bs.Max[i][j] - bs.Allocation[i][j]
		}
	}

	return bs
}

// copyMatrix 复制二维矩阵
func copyMatrix(src [][]int) [][]int {
	dst := make([][]int, len(src))
	for i := range src {
		dst[i] = make([]int, len(src[i]))
		copy(dst[i], src[i])
	}
	return dst
}

// IsSafe 安全性检查（银行家算法核心）
// 408考点：判断系统是否处于安全状态
// 安全状态：存在一个进程执行序列，使得每个进程都能顺利执行完毕
func (bs *BankerSystem) IsSafe() bool {
	// 工作向量，初始等于可用资源
	work := make([]int, bs.NumResources)
	copy(work, bs.Available)

	// 完成标记，初始都为false
	finish := make([]bool, bs.NumProcesses)

	// 安全序列
	safeSeq := make([]int, 0, bs.NumProcesses)

	// 循环查找可以执行的进程
	count := 0
	for count < bs.NumProcesses {
		found := false

		for i := 0; i < bs.NumProcesses; i++ {
			if finish[i] {
				continue
			}

			// 检查该进程的需求是否能被满足
			canAllocate := true
			for j := 0; j < bs.NumResources; j++ {
				if bs.Need[i][j] > work[j] {
					canAllocate = false
					break
				}
			}

			if canAllocate {
				// 模拟分配和释放
				for j := 0; j < bs.NumResources; j++ {
					work[j] += bs.Allocation[i][j]
				}
				finish[i] = true
				safeSeq = append(safeSeq, i)
				found = true
				count++
				break
			}
		}

		if !found {
			// 找不到可执行的进程，系统不安全
			return false
		}
	}

	// 所有进程都能执行完，系统安全
	bs.SafeSequence = safeSeq
	return true
}

// RequestResources 资源请求（银行家算法）
// 408考点：进程请求资源时的处理流程
// 1. 检查请求是否超过需求
// 2. 检查请求是否超过可用资源
// 3. 试探性分配
// 4. 安全性检查
// 5. 若安全则分配，否则回滚
func (bs *BankerSystem) RequestResources(processID int, request []int) bool {
	// 1. 检查请求是否超过声明的最大需求
	for j := 0; j < bs.NumResources; j++ {
		if request[j] > bs.Need[processID][j] {
			fmt.Printf("   ❌ 错误：进程 P%d 请求的资源超过其声明的最大需求\n", processID)
			return false
		}
	}

	// 2. 检查请求是否超过系统可用资源
	for j := 0; j < bs.NumResources; j++ {
		if request[j] > bs.Available[j] {
			fmt.Printf("   ⏳ 进程 P%d 必须等待，资源不足\n", processID)
			return false
		}
	}

	// 3. 试探性分配资源
	oldAvailable := make([]int, bs.NumResources)
	copy(oldAvailable, bs.Available)

	for j := 0; j < bs.NumResources; j++ {
		bs.Available[j] -= request[j]
		bs.Allocation[processID][j] += request[j]
		bs.Need[processID][j] -= request[j]
	}

	// 4. 安全性检查
	if bs.IsSafe() {
		fmt.Printf("   ✓ 分配成功，系统处于安全状态\n")
		fmt.Printf("   安全序列: %v\n", formatProcessSequence(bs.SafeSequence))
		return true
	} else {
		// 5. 系统不安全，回滚
		for j := 0; j < bs.NumResources; j++ {
			bs.Available[j] = oldAvailable[j]
			bs.Allocation[processID][j] -= request[j]
			bs.Need[processID][j] += request[j]
		}
		fmt.Printf("   ❌ 分配失败，会导致系统不安全，请求被拒绝\n")
		return false
	}
}

// Print 打印系统状态
func (bs *BankerSystem) Print() {
	fmt.Println("\n系统当前状态:")
	fmt.Printf("  可用资源 Available: %v\n\n", bs.Available)

	fmt.Println("  进程资源分配表:")
	fmt.Printf("  %-8s", "进程")
	for j := 0; j < bs.NumResources; j++ {
		fmt.Printf("  Max[%d]", j)
	}
	for j := 0; j < bs.NumResources; j++ {
		fmt.Printf("  Alloc[%d]", j)
	}
	for j := 0; j < bs.NumResources; j++ {
		fmt.Printf("  Need[%d]", j)
	}
	fmt.Println()

	for i := 0; i < bs.NumProcesses; i++ {
		fmt.Printf("  P%-7d", i)
		for j := 0; j < bs.NumResources; j++ {
			fmt.Printf("  %-7d", bs.Max[i][j])
		}
		for j := 0; j < bs.NumResources; j++ {
			fmt.Printf("  %-8d", bs.Allocation[i][j])
		}
		for j := 0; j < bs.NumResources; j++ {
			fmt.Printf("  %-7d", bs.Need[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}

// formatProcessSequence 格式化进程序列
func formatProcessSequence(seq []int) string {
	if len(seq) == 0 {
		return "[]"
	}
	result := "["
	for i, p := range seq {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("P%d", p)
	}
	result += "]"
	return result
}

// DeadlockDetection 死锁检测
// 408考点：使用资源分配图检测死锁
// 简化版：检查是否存在循环等待
type DeadlockDetection struct {
	NumProcesses int     // 进程数量
	NumResources int     // 资源数量
	Allocation   [][]int // 已分配矩阵
	Request      [][]int // 请求矩阵
	Available    []int   // 可用资源
}

// NewDeadlockDetection 创建死锁检测系统
func NewDeadlockDetection(numProcesses, numResources int, allocation, request [][]int, available []int) *DeadlockDetection {
	return &DeadlockDetection{
		NumProcesses: numProcesses,
		NumResources: numResources,
		Allocation:   copyMatrix(allocation),
		Request:      copyMatrix(request),
		Available:    append([]int{}, available...),
	}
}

// DetectDeadlock 检测死锁
// 408考点：死锁检测算法
// 原理：模拟银行家算法的安全性检查，如果不能找到安全序列，则存在死锁
func (dd *DeadlockDetection) DetectDeadlock() []int {
	work := make([]int, dd.NumResources)
	copy(work, dd.Available)

	finish := make([]bool, dd.NumProcesses)
	deadlocked := make([]int, 0)

	// 找出所有能完成的进程
	changed := true
	for changed {
		changed = false
		for i := 0; i < dd.NumProcesses; i++ {
			if finish[i] {
				continue
			}

			// 检查进程i的请求是否能被满足
			canFinish := true
			for j := 0; j < dd.NumResources; j++ {
				if dd.Request[i][j] > work[j] {
					canFinish = false
					break
				}
			}

			if canFinish {
				// 进程i可以完成，释放其资源
				for j := 0; j < dd.NumResources; j++ {
					work[j] += dd.Allocation[i][j]
				}
				finish[i] = true
				changed = true
			}
		}
	}

	// 未完成的进程即为死锁进程
	for i := 0; i < dd.NumProcesses; i++ {
		if !finish[i] {
			deadlocked = append(deadlocked, i)
		}
	}

	return deadlocked
}

// DeadlockExample 死锁示例
func DeadlockExample() {
	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║            操作系统 - 死锁处理 (Deadlock Handling)        ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")

	// ============ 示例1: 银行家算法 - 安全状态 ============
	fmt.Println("\n【示例1】银行家算法 - 安全状态演示")
	fmt.Println("─────────────────────────────────────────────────")

	// 5个进程，3种资源 (A, B, C)
	max1 := [][]int{
		{7, 5, 3}, // P0
		{3, 2, 2}, // P1
		{9, 0, 2}, // P2
		{2, 2, 2}, // P3
		{4, 3, 3}, // P4
	}
	allocation1 := [][]int{
		{0, 1, 0}, // P0
		{2, 0, 0}, // P1
		{3, 0, 2}, // P2
		{2, 1, 1}, // P3
		{0, 0, 2}, // P4
	}
	available1 := []int{3, 3, 2} // 可用资源

	banker1 := NewBankerSystem(5, 3, max1, allocation1, available1)
	banker1.Print()

	fmt.Println("进行安全性检查:")
	if banker1.IsSafe() {
		fmt.Printf("  ✓ 系统处于安全状态\n")
		fmt.Printf("  安全序列: %v\n", formatProcessSequence(banker1.SafeSequence))
	} else {
		fmt.Printf("  ❌ 系统处于不安全状态\n")
	}

	// 进程P1请求资源
	fmt.Println("\n进程 P1 请求资源 [1, 0, 2]:")
	request := []int{1, 0, 2}
	banker1.RequestResources(1, request)

	// ============ 示例2: 银行家算法 - 不安全请求 ============
	fmt.Println("\n\n【示例2】银行家算法 - 拒绝导致不安全的请求")
	fmt.Println("─────────────────────────────────────────────────")

	// 进程P4请求过多资源
	fmt.Println("进程 P4 请求资源 [3, 3, 0]:")
	request2 := []int{3, 3, 0}
	banker1.RequestResources(4, request2)

	// ============ 示例3: 死锁检测 ============
	fmt.Println("\n\n【示例3】死锁检测")
	fmt.Println("─────────────────────────────────────────────────")

	// 构造一个死锁场景
	allocation3 := [][]int{
		{0, 1, 0}, // P0 持有资源B
		{2, 0, 0}, // P1 持有资源A
		{3, 0, 3}, // P2 持有资源A和C
		{2, 1, 1}, // P3 持有资源A、B、C
		{0, 0, 2}, // P4 持有资源C
	}
	request3 := [][]int{
		{0, 0, 1}, // P0 请求资源C
		{1, 0, 1}, // P1 请求资源A和C
		{3, 1, 0}, // P2 请求资源A和B
		{0, 0, 1}, // P3 请求资源C
		{0, 1, 0}, // P4 请求资源B
	}
	available3 := []int{0, 0, 0} // 无可用资源

	dd := NewDeadlockDetection(5, 3, allocation3, request3, available3)

	fmt.Println("当前资源分配情况:")
	fmt.Printf("  可用资源: %v\n", dd.Available)
	fmt.Println("\n  进程资源分配表:")
	fmt.Printf("  %-8s  Allocation    Request\n", "进程")
	for i := 0; i < dd.NumProcesses; i++ {
		fmt.Printf("  P%-7d  %v       %v\n", i, dd.Allocation[i], dd.Request[i])
	}

	fmt.Println("\n执行死锁检测:")
	deadlocked := dd.DetectDeadlock()
	if len(deadlocked) > 0 {
		fmt.Printf("  ❌ 检测到死锁！\n")
		fmt.Printf("  死锁进程: %v\n", formatProcessSequence(deadlocked))
	} else {
		fmt.Printf("  ✓ 系统无死锁\n")
	}

	// ============ 示例4: 正常无死锁情况 ============
	fmt.Println("\n\n【示例4】死锁检测 - 无死锁情况")
	fmt.Println("─────────────────────────────────────────────────")

	allocation4 := [][]int{
		{0, 1, 0}, // P0
		{2, 0, 0}, // P1
		{3, 0, 2}, // P2
		{2, 1, 1}, // P3
		{0, 0, 2}, // P4
	}
	request4 := [][]int{
		{0, 0, 0}, // P0 无请求
		{0, 0, 1}, // P1 请求C
		{0, 0, 0}, // P2 无请求
		{0, 0, 0}, // P3 无请求
		{0, 1, 0}, // P4 请求B
	}
	available4 := []int{3, 2, 1}

	dd2 := NewDeadlockDetection(5, 3, allocation4, request4, available4)

	fmt.Println("当前资源分配情况:")
	fmt.Printf("  可用资源: %v\n", dd2.Available)

	fmt.Println("\n执行死锁检测:")
	deadlocked2 := dd2.DetectDeadlock()
	if len(deadlocked2) > 0 {
		fmt.Printf("  ❌ 检测到死锁！\n")
		fmt.Printf("  死锁进程: %v\n", formatProcessSequence(deadlocked2))
	} else {
		fmt.Printf("  ✓ 系统无死锁\n")
	}

	// ============ 408考点总结 ============
	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                     408 考点总结                           ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")

	fmt.Println("\n1. 死锁的四个必要条件（缺一不可）:")
	fmt.Println("   • 互斥：资源不能共享，一次只能被一个进程使用")
	fmt.Println("   • 占有并等待：进程已获得资源，同时等待新资源")
	fmt.Println("   • 非抢占：资源不能被强制剥夺，只能主动释放")
	fmt.Println("   • 循环等待：存在进程资源的环形等待链")

	fmt.Println("\n2. 死锁处理策略:")
	fmt.Println("   • 预防（Prevention）：破坏四个必要条件之一")
	fmt.Println("   • 避免（Avoidance）：银行家算法，动态检查安全性")
	fmt.Println("   • 检测与恢复（Detection & Recovery）：允许死锁发生，检测后恢复")
	fmt.Println("   • 鸵鸟策略：忽略死锁问题（如UNIX）")

	fmt.Println("\n3. 银行家算法（Banker's Algorithm）:")
	fmt.Println("   • 核心思想：在分配资源前进行安全性检查")
	fmt.Println("   • 关键数据结构：")
	fmt.Println("     - Available: 可用资源向量")
	fmt.Println("     - Max: 最大需求矩阵")
	fmt.Println("     - Allocation: 已分配矩阵")
	fmt.Println("     - Need: 需求矩阵 (Need = Max - Allocation)")
	fmt.Println("   • 安全状态：存在安全序列，使所有进程都能顺利完成")
	fmt.Println("   • 算法步骤：")
	fmt.Println("     1. 检查请求是否超过需求")
	fmt.Println("     2. 检查请求是否超过可用资源")
	fmt.Println("     3. 试探性分配")
	fmt.Println("     4. 执行安全性算法")
	fmt.Println("     5. 若安全则分配，否则回滚")

	fmt.Println("\n4. 死锁检测算法:")
	fmt.Println("   • 类似银行家算法的安全性检查")
	fmt.Println("   • 不能完成的进程即为死锁进程")
	fmt.Println("   • 检测时机：定期检测或资源利用率下降时")

	fmt.Println("\n5. 死锁恢复方法:")
	fmt.Println("   • 资源剥夺：抢占死锁进程的资源")
	fmt.Println("   • 撤销进程：终止死锁进程（一个或全部）")
	fmt.Println("   • 进程回退：回退到安全状态")

	fmt.Println("\n6. 考试常见题型:")
	fmt.Println("   • 判断系统是否安全，给出安全序列")
	fmt.Println("   • 判断某个资源请求能否被满足")
	fmt.Println("   • 计算Need矩阵")
	fmt.Println("   • 识别死锁的四个必要条件")
	fmt.Println("   • 给出资源分配图，判断是否存在死锁")
	fmt.Println()
}
