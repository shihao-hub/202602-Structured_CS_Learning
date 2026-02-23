package scheduling

import (
	"fmt"
	"math"
	"sort"
)

// DiskRequest 磁盘请求
type DiskRequest struct {
	Track int // 磁道号
}

// DiskScheduler 磁盘调度器
type DiskScheduler struct {
	CurrentHead int            // 当前磁头位置
	DiskSize    int            // 磁盘大小（磁道数）
	Requests    []DiskRequest  // 请求队列
	Direction   int            // 磁头移动方向 (1: 向外, -1: 向内)
}

// DiskScheduleResult 磁盘调度结果
type DiskScheduleResult struct {
	Order         []int // 服务顺序
	TotalMovement int   // 总磁头移动距离
}

// NewDiskScheduler 创建磁盘调度器
func NewDiskScheduler(currentHead, diskSize int, requests []int) *DiskScheduler {
	ds := &DiskScheduler{
		CurrentHead: currentHead,
		DiskSize:    diskSize,
		Requests:    make([]DiskRequest, len(requests)),
		Direction:   1, // 默认向外
	}
	for i, track := range requests {
		ds.Requests[i] = DiskRequest{Track: track}
	}
	return ds
}

// FCFS 首先来先服务（First Come First Served）
// 408考点：最简单的磁盘调度算法，按请求到达顺序服务
// 优点：公平、简单、无饥饿
// 缺点：平均寻道时间长，磁头移动距离大
func (ds *DiskScheduler) FCFS() *DiskScheduleResult {
	result := &DiskScheduleResult{
		Order:         make([]int, 0, len(ds.Requests)),
		TotalMovement: 0,
	}

	current := ds.CurrentHead
	for _, req := range ds.Requests {
		result.Order = append(result.Order, req.Track)
		result.TotalMovement += int(math.Abs(float64(req.Track - current)))
		current = req.Track
	}

	return result
}

// SSTF 最短寻道时间优先（Shortest Seek Time First）
// 408考点：贪心算法，每次选择距离当前磁头最近的请求
// 优点：平均寻道时间短，吞吐量高
// 缺点：可能产生饥饿现象（远端请求长期得不到服务）
func (ds *DiskScheduler) SSTF() *DiskScheduleResult {
	result := &DiskScheduleResult{
		Order:         make([]int, 0, len(ds.Requests)),
		TotalMovement: 0,
	}

	// 复制请求列表
	remaining := make([]int, len(ds.Requests))
	for i, req := range ds.Requests {
		remaining[i] = req.Track
	}

	current := ds.CurrentHead
	for len(remaining) > 0 {
		// 找到距离当前位置最近的请求
		minDist := math.MaxInt32
		minIdx := -1
		for i, track := range remaining {
			dist := int(math.Abs(float64(track - current)))
			if dist < minDist {
				minDist = dist
				minIdx = i
			}
		}

		// 服务该请求
		result.Order = append(result.Order, remaining[minIdx])
		result.TotalMovement += minDist
		current = remaining[minIdx]

		// 从待服务列表中移除
		remaining = append(remaining[:minIdx], remaining[minIdx+1:]...)
	}

	return result
}

// SCAN 扫描算法（电梯算法）
// 408考点：磁头在一个方向上移动，服务路径上所有请求，到达边界后反向
// 优点：避免饥饿，寻道性能较好
// 缺点：对边缘磁道和中间磁道不公平（中间磁道等待时间更短）
func (ds *DiskScheduler) SCAN() *DiskScheduleResult {
	result := &DiskScheduleResult{
		Order:         make([]int, 0, len(ds.Requests)),
		TotalMovement: 0,
	}

	// 分组：当前位置以内和以外的请求
	inner := make([]int, 0)  // 内侧请求（小于当前位置）
	outer := make([]int, 0)  // 外侧请求（大于等于当前位置）

	for _, req := range ds.Requests {
		if req.Track < ds.CurrentHead {
			inner = append(inner, req.Track)
		} else {
			outer = append(outer, req.Track)
		}
	}

	// 排序
	sort.Ints(inner)
	sort.Ints(outer)

	current := ds.CurrentHead

	// 假设初始方向向外
	if len(outer) > 0 {
		// 先处理外侧请求（向外扫描）
		for _, track := range outer {
			result.Order = append(result.Order, track)
			result.TotalMovement += track - current
			current = track
		}

		// 到达边界后，如果还有内侧请求，需要反向
		if len(inner) > 0 {
			// 到达最外侧
			result.TotalMovement += (ds.DiskSize - 1) - current
			current = ds.DiskSize - 1

			// 反向处理内侧请求（从大到小）
			for i := len(inner) - 1; i >= 0; i-- {
				result.Order = append(result.Order, inner[i])
				result.TotalMovement += current - inner[i]
				current = inner[i]
			}
		}
	} else if len(inner) > 0 {
		// 只有内侧请求，向内扫描（从大到小）
		for i := len(inner) - 1; i >= 0; i-- {
			result.Order = append(result.Order, inner[i])
			result.TotalMovement += current - inner[i]
			current = inner[i]
		}
	}

	return result
}

// C-SCAN 循环扫描算法（Circular SCAN）
// 408考点：单向扫描，到达边界后直接返回起始端继续扫描（不服务返程请求）
// 优点：对所有磁道更加公平，等待时间方差小
// 缺点：返回起始端的空移动增加了总的寻道时间
func (ds *DiskScheduler) CSCAN() *DiskScheduleResult {
	result := &DiskScheduleResult{
		Order:         make([]int, 0, len(ds.Requests)),
		TotalMovement: 0,
	}

	// 分组并排序
	inner := make([]int, 0)
	outer := make([]int, 0)

	for _, req := range ds.Requests {
		if req.Track < ds.CurrentHead {
			inner = append(inner, req.Track)
		} else {
			outer = append(outer, req.Track)
		}
	}

	sort.Ints(inner)
	sort.Ints(outer)

	current := ds.CurrentHead

	// 先向外扫描
	for _, track := range outer {
		result.Order = append(result.Order, track)
		result.TotalMovement += track - current
		current = track
	}

	// 到达最外侧
	if len(outer) > 0 {
		result.TotalMovement += (ds.DiskSize - 1) - current
		current = ds.DiskSize - 1
	}

	// 返回到磁盘起始位置（0号磁道）
	if len(inner) > 0 {
		result.TotalMovement += current - 0
		current = 0

		// 从内侧最小的开始服务
		for _, track := range inner {
			result.Order = append(result.Order, track)
			result.TotalMovement += track - current
			current = track
		}
	}

	return result
}

// LOOK 改进的SCAN算法
// 408考点：与SCAN类似，但磁头不必到达边界，只需到达该方向最远的请求即可反向
// 优点：减少了不必要的边界移动，性能优于SCAN
// 缺点：实现稍复杂
func (ds *DiskScheduler) LOOK() *DiskScheduleResult {
	result := &DiskScheduleResult{
		Order:         make([]int, 0, len(ds.Requests)),
		TotalMovement: 0,
	}

	// 分组并排序
	inner := make([]int, 0)
	outer := make([]int, 0)

	for _, req := range ds.Requests {
		if req.Track < ds.CurrentHead {
			inner = append(inner, req.Track)
		} else {
			outer = append(outer, req.Track)
		}
	}

	sort.Ints(inner)
	sort.Ints(outer)

	current := ds.CurrentHead

	// 先向外扫描到最远请求
	if len(outer) > 0 {
		for _, track := range outer {
			result.Order = append(result.Order, track)
			result.TotalMovement += track - current
			current = track
		}
	}

	// 反向处理内侧请求（从大到小）
	if len(inner) > 0 {
		for i := len(inner) - 1; i >= 0; i-- {
			result.Order = append(result.Order, inner[i])
			result.TotalMovement += current - inner[i]
			current = inner[i]
		}
	}

	return result
}

// C-LOOK 循环LOOK算法
// 408考点：LOOK的循环版本，到达该方向最远请求后直接跳到另一端最远请求
// 优点：结合了C-SCAN和LOOK的优点，性能更好
func (ds *DiskScheduler) CLOOK() *DiskScheduleResult {
	result := &DiskScheduleResult{
		Order:         make([]int, 0, len(ds.Requests)),
		TotalMovement: 0,
	}

	// 分组并排序
	inner := make([]int, 0)
	outer := make([]int, 0)

	for _, req := range ds.Requests {
		if req.Track < ds.CurrentHead {
			inner = append(inner, req.Track)
		} else {
			outer = append(outer, req.Track)
		}
	}

	sort.Ints(inner)
	sort.Ints(outer)

	current := ds.CurrentHead

	// 先向外扫描到最远请求
	for _, track := range outer {
		result.Order = append(result.Order, track)
		result.TotalMovement += track - current
		current = track
	}

	// 跳到内侧最小的请求（单向循环）
	if len(inner) > 0 && len(outer) > 0 {
		result.TotalMovement += current - inner[0]
		current = inner[0]
		result.Order = append(result.Order, current)

		// 继续向外处理剩余内侧请求
		for i := 1; i < len(inner); i++ {
			result.Order = append(result.Order, inner[i])
			result.TotalMovement += inner[i] - current
			current = inner[i]
		}
	} else if len(inner) > 0 {
		// 如果没有外侧请求，直接处理内侧
		for _, track := range inner {
			result.Order = append(result.Order, track)
			result.TotalMovement += int(math.Abs(float64(track - current)))
			current = track
		}
	}

	return result
}

// DiskSchedulerExample 磁盘调度算法示例
func DiskSchedulerExample() {
	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║         操作系统 - 磁盘调度算法 (Disk Scheduling)         ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")

	// 测试数据
	currentHead := 53
	diskSize := 200
	requests := []int{98, 183, 37, 122, 14, 124, 65, 67}

	fmt.Printf("\n初始状态:\n")
	fmt.Printf("  当前磁头位置: %d\n", currentHead)
	fmt.Printf("  磁盘大小: %d 磁道\n", diskSize)
	fmt.Printf("  请求队列: %v\n", requests)
	fmt.Println()

	// 创建调度器
	ds := NewDiskScheduler(currentHead, diskSize, requests)

	// 1. FCFS
	fmt.Println("1. FCFS (First Come First Served) - 先来先服务")
	fcfsResult := ds.FCFS()
	fmt.Printf("   服务顺序: %v\n", fcfsResult.Order)
	fmt.Printf("   总磁头移动距离: %d\n", fcfsResult.TotalMovement)
	fmt.Printf("   平均寻道长度: %.2f\n\n", float64(fcfsResult.TotalMovement)/float64(len(requests)))

	// 2. SSTF
	fmt.Println("2. SSTF (Shortest Seek Time First) - 最短寻道时间优先")
	sstfResult := ds.SSTF()
	fmt.Printf("   服务顺序: %v\n", sstfResult.Order)
	fmt.Printf("   总磁头移动距离: %d\n", sstfResult.TotalMovement)
	fmt.Printf("   平均寻道长度: %.2f\n", float64(sstfResult.TotalMovement)/float64(len(requests)))
	fmt.Printf("   ⚠️  注意：可能产生饥饿现象\n\n")

	// 3. SCAN
	fmt.Println("3. SCAN (Elevator Algorithm) - 扫描算法/电梯算法")
	scanResult := ds.SCAN()
	fmt.Printf("   服务顺序: %v\n", scanResult.Order)
	fmt.Printf("   总磁头移动距离: %d\n", scanResult.TotalMovement)
	fmt.Printf("   平均寻道长度: %.2f\n\n", float64(scanResult.TotalMovement)/float64(len(requests)))

	// 4. C-SCAN
	fmt.Println("4. C-SCAN (Circular SCAN) - 循环扫描算法")
	cscanResult := ds.CSCAN()
	fmt.Printf("   服务顺序: %v\n", cscanResult.Order)
	fmt.Printf("   总磁头移动距离: %d\n", cscanResult.TotalMovement)
	fmt.Printf("   平均寻道长度: %.2f\n\n", float64(cscanResult.TotalMovement)/float64(len(requests)))

	// 5. LOOK
	fmt.Println("5. LOOK - 改进的SCAN算法")
	lookResult := ds.LOOK()
	fmt.Printf("   服务顺序: %v\n", lookResult.Order)
	fmt.Printf("   总磁头移动距离: %d\n", lookResult.TotalMovement)
	fmt.Printf("   平均寻道长度: %.2f\n", float64(lookResult.TotalMovement)/float64(len(requests)))
	fmt.Printf("   ✓ 相比SCAN，不需要到达边界\n\n")

	// 6. C-LOOK
	fmt.Println("6. C-LOOK - 循环LOOK算法")
	clookResult := ds.CLOOK()
	fmt.Printf("   服务顺序: %v\n", clookResult.Order)
	fmt.Printf("   总磁头移动距离: %d\n", clookResult.TotalMovement)
	fmt.Printf("   平均寻道长度: %.2f\n\n", float64(clookResult.TotalMovement)/float64(len(requests)))

	// 性能比较
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    算法性能比较                            ║")
	fmt.Println("╠═════════════╦════════════════╦═══════════════════════════╣")
	fmt.Println("║   算法      ║ 总移动距离     ║        特点               ║")
	fmt.Println("╠═════════════╬════════════════╬═══════════════════════════╣")
	fmt.Printf("║ FCFS        ║ %-14d ║ 简单公平，性能差       ║\n", fcfsResult.TotalMovement)
	fmt.Printf("║ SSTF        ║ %-14d ║ 性能好，可能饥饿       ║\n", sstfResult.TotalMovement)
	fmt.Printf("║ SCAN        ║ %-14d ║ 避免饥饿，往返扫描     ║\n", scanResult.TotalMovement)
	fmt.Printf("║ C-SCAN      ║ %-14d ║ 更公平，单向扫描       ║\n", cscanResult.TotalMovement)
	fmt.Printf("║ LOOK        ║ %-14d ║ SCAN改进，不到边界     ║\n", lookResult.TotalMovement)
	fmt.Printf("║ C-LOOK      ║ %-14d ║ C-SCAN改进，最优       ║\n", clookResult.TotalMovement)
	fmt.Println("╚═════════════╩════════════════╩═══════════════════════════╝")

	// 408考点总结
	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                  408 考点总结                              ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println("1. 磁盘调度算法分类:")
	fmt.Println("   • FCFS: 公平但效率低")
	fmt.Println("   • SSTF: 贪心算法，可能饥饿")
	fmt.Println("   • SCAN/C-SCAN: 电梯算法，避免饥饿")
	fmt.Println("   • LOOK/C-LOOK: 改进的电梯算法")
	fmt.Println()
	fmt.Println("2. 关键性能指标:")
	fmt.Println("   • 平均寻道时间/距离")
	fmt.Println("   • 响应时间方差（公平性）")
	fmt.Println("   • 是否存在饥饿现象")
	fmt.Println()
	fmt.Println("3. 算法选择原则:")
	fmt.Println("   • 负载轻: FCFS")
	fmt.Println("   • 负载重且请求密集: SSTF")
	fmt.Println("   • 需要避免饥饿: SCAN/LOOK")
	fmt.Println("   • 要求公平性: C-SCAN/C-LOOK")
	fmt.Println()
	fmt.Println("4. 考试常见题型:")
	fmt.Println("   • 给定请求序列，计算某算法的总寻道距离")
	fmt.Println("   • 比较不同算法的性能优劣")
	fmt.Println("   • 判断算法是否会产生饥饿")
	fmt.Println()
}
