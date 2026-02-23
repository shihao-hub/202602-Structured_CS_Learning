package memory

import (
	"fmt"
)

// PageTableEntry 页表项
// 408 考点：页表项的结构（页框号、有效位、修改位、访问位）
type PageTableEntry struct {
	FrameNumber    int  // 页框号（物理页号）
	Valid          bool // 有效位：该页是否在内存中
	Modified       bool // 修改位（脏位）：该页是否被修改过
	Referenced     bool // 访问位：该页最近是否被访问过
	LoadTime       int  // 装入时间：用于 FIFO
	LastAccessTime int  // 最后访问时间：用于 LRU
}

// TLBEntry TLB 表项
// 408 考点：TLB 是页表的高速缓存，采用全相联方式
type TLBEntry struct {
	Valid       bool // 有效位
	PageNumber  int  // 虚拟页号
	FrameNumber int  // 物理页框号
	AccessTime  int  // 访问时间（用于 LRU）
}

// PageReplacementPolicy 页面替换算法
type PageReplacementPolicy int

const (
	PageFIFO    PageReplacementPolicy = iota // 先进先出
	PageLRU                                  // 最近最少使用
	PageClock                                // 时钟算法（NRU）
	PageOptimal                              // 最佳置换算法
)

// VirtualMemorySimulator 虚拟存储器模拟器
type VirtualMemorySimulator struct {
	PageTable          []PageTableEntry      // 页表
	TLB                []TLBEntry            // TLB
	PhysicalMemory     []int                 // 物理内存（存储页号，-1 表示空闲）
	PageSize           int                   // 页面大小
	NumPages           int                   // 虚拟页数
	NumFrames          int                   // 物理页框数
	TLBSize            int                   // TLB 大小
	Policy             PageReplacementPolicy // 页面替换策略
	CurrentTime        int                   // 当前时间
	ClockPointer       int                   // 时钟指针（用于 Clock 算法）
	PageFaults         int                   // 缺页次数
	TLBHits            int                   // TLB 命中次数
	TLBMisses          int                   // TLB 缺失次数
	AccessCount        int                   // 总访问次数
	FutureAccesses     []int                 // 未来访问序列（用于 OPT 算法）
	CurrentAccessIndex int                   // 当前访问索引
}

// NewVirtualMemorySimulator 创建虚拟存储器模拟器
func NewVirtualMemorySimulator(numPages, numFrames, tlbSize int, policy PageReplacementPolicy) *VirtualMemorySimulator {
	pageTable := make([]PageTableEntry, numPages)
	for i := range pageTable {
		pageTable[i] = PageTableEntry{
			FrameNumber: -1,
			Valid:       false,
			Modified:    false,
			Referenced:  false,
		}
	}

	tlb := make([]TLBEntry, tlbSize)
	for i := range tlb {
		tlb[i] = TLBEntry{Valid: false}
	}

	physicalMemory := make([]int, numFrames)
	for i := range physicalMemory {
		physicalMemory[i] = -1 // -1 表示空闲
	}

	return &VirtualMemorySimulator{
		PageTable:      pageTable,
		TLB:            tlb,
		PhysicalMemory: physicalMemory,
		PageSize:       4096, // 4KB
		NumPages:       numPages,
		NumFrames:      numFrames,
		TLBSize:        tlbSize,
		Policy:         policy,
		CurrentTime:    0,
		ClockPointer:   0,
		PageFaults:     0,
		TLBHits:        0,
		TLBMisses:      0,
		AccessCount:    0,
	}
}

// SetFutureAccesses 设置未来访问序列（用于 OPT 算法）
func (vms *VirtualMemorySimulator) SetFutureAccesses(accesses []int) {
	vms.FutureAccesses = accesses
	vms.CurrentAccessIndex = 0
}

// Access 访问虚拟地址
// 408 考点：地址转换过程（TLB → 页表 → 缺页处理）
func (vms *VirtualMemorySimulator) Access(virtualAddress int) string {
	vms.AccessCount++
	vms.CurrentTime++
	if vms.Policy == PageOptimal {
		vms.CurrentAccessIndex++
	}

	// 1. 分解虚拟地址
	pageNumber, offset := vms.parseVirtualAddress(virtualAddress)

	var result string
	result += fmt.Sprintf("访问虚拟地址 0x%04X (页号: %d, 偏移: 0x%03X)\n",
		virtualAddress, pageNumber, offset)

	// 2. 查找 TLB
	frameNumber, tlbHit := vms.searchTLB(pageNumber)
	if tlbHit {
		vms.TLBHits++
		physicalAddress := frameNumber*vms.PageSize + offset
		result += fmt.Sprintf("  ✓ TLB 命中! 页框号: %d, 物理地址: 0x%04X\n",
			frameNumber, physicalAddress)
		vms.PageTable[pageNumber].Referenced = true
		vms.PageTable[pageNumber].LastAccessTime = vms.CurrentTime
		return result
	}

	// 3. TLB 缺失，查找页表
	vms.TLBMisses++
	result += "  ✗ TLB 缺失，查找页表...\n"

	if vms.PageTable[pageNumber].Valid {
		// 页表命中，页面在内存中
		frameNumber = vms.PageTable[pageNumber].FrameNumber
		physicalAddress := frameNumber*vms.PageSize + offset
		result += fmt.Sprintf("  ✓ 页表命中! 页框号: %d, 物理地址: 0x%04X\n",
			frameNumber, physicalAddress)

		// 更新 TLB
		vms.updateTLB(pageNumber, frameNumber)
		result += "  → 更新 TLB\n"

		// 更新访问信息
		vms.PageTable[pageNumber].Referenced = true
		vms.PageTable[pageNumber].LastAccessTime = vms.CurrentTime
	} else {
		// 4. 缺页中断
		vms.PageFaults++
		result += "  ✗ 缺页中断！需要从磁盘调入页面\n"

		// 选择被替换的页框
		frameNumber = vms.selectVictimFrame()

		// 检查该页框是否已被占用
		oldPage := -1
		for _, frame := range vms.PhysicalMemory {
			if frame != -1 && vms.PageTable[frame].FrameNumber == frameNumber {
				oldPage = frame
				break
			}
		}

		if oldPage != -1 {
			// 如果页面被修改过，需要写回磁盘
			if vms.PageTable[oldPage].Modified {
				result += fmt.Sprintf("  → 页面 %d 已修改，写回磁盘\n", oldPage)
			}
			// 更新旧页表项
			vms.PageTable[oldPage].Valid = false
			vms.PageTable[oldPage].FrameNumber = -1
		}

		// 调入新页面
		vms.PhysicalMemory[frameNumber] = pageNumber
		vms.PageTable[pageNumber].Valid = true
		vms.PageTable[pageNumber].FrameNumber = frameNumber
		vms.PageTable[pageNumber].Modified = false
		vms.PageTable[pageNumber].Referenced = true
		vms.PageTable[pageNumber].LoadTime = vms.CurrentTime
		vms.PageTable[pageNumber].LastAccessTime = vms.CurrentTime

		physicalAddress := frameNumber*vms.PageSize + offset
		result += fmt.Sprintf("  → 页面 %d 调入页框 %d, 物理地址: 0x%04X\n",
			pageNumber, frameNumber, physicalAddress)

		// 更新 TLB
		vms.updateTLB(pageNumber, frameNumber)
		result += "  → 更新 TLB\n"
	}

	return result
}

// parseVirtualAddress 解析虚拟地址
func (vms *VirtualMemorySimulator) parseVirtualAddress(address int) (pageNumber, offset int) {
	offset = address % vms.PageSize
	pageNumber = address / vms.PageSize
	return
}

// searchTLB 在 TLB 中查找页号
func (vms *VirtualMemorySimulator) searchTLB(pageNumber int) (frameNumber int, hit bool) {
	for i := range vms.TLB {
		if vms.TLB[i].Valid && vms.TLB[i].PageNumber == pageNumber {
			vms.TLB[i].AccessTime = vms.CurrentTime
			return vms.TLB[i].FrameNumber, true
		}
	}
	return -1, false
}

// updateTLB 更新 TLB（使用 LRU 替换）
func (vms *VirtualMemorySimulator) updateTLB(pageNumber, frameNumber int) {
	// 查找空闲或最久未使用的 TLB 项
	victim := 0
	minTime := vms.TLB[0].AccessTime

	for i := range vms.TLB {
		if !vms.TLB[i].Valid {
			victim = i
			break
		}
		if vms.TLB[i].AccessTime < minTime {
			minTime = vms.TLB[i].AccessTime
			victim = i
		}
	}

	vms.TLB[victim] = TLBEntry{
		Valid:       true,
		PageNumber:  pageNumber,
		FrameNumber: frameNumber,
		AccessTime:  vms.CurrentTime,
	}
}

// selectVictimFrame 选择被替换的页框
// 408 考点：四种页面替换算法
func (vms *VirtualMemorySimulator) selectVictimFrame() int {
	// 首先查找空闲页框
	for i := range vms.PhysicalMemory {
		if vms.PhysicalMemory[i] == -1 {
			return i
		}
	}

	// 根据替换策略选择牺牲页框
	switch vms.Policy {
	case PageFIFO:
		return vms.selectFIFO()
	case PageLRU:
		return vms.selectLRU()
	case PageClock:
		return vms.selectClock()
	case PageOptimal:
		return vms.selectOptimal()
	default:
		return 0
	}
}

// selectFIFO FIFO 算法：选择最早装入的页面
func (vms *VirtualMemorySimulator) selectFIFO() int {
	minLoadTime := vms.CurrentTime
	victim := 0

	for i, pageNum := range vms.PhysicalMemory {
		if pageNum != -1 && vms.PageTable[pageNum].LoadTime < minLoadTime {
			minLoadTime = vms.PageTable[pageNum].LoadTime
			victim = i
		}
	}

	return victim
}

// selectLRU LRU 算法：选择最久未使用的页面
func (vms *VirtualMemorySimulator) selectLRU() int {
	minAccessTime := vms.CurrentTime
	victim := 0

	for i, pageNum := range vms.PhysicalMemory {
		if pageNum != -1 && vms.PageTable[pageNum].LastAccessTime < minAccessTime {
			minAccessTime = vms.PageTable[pageNum].LastAccessTime
			victim = i
		}
	}

	return victim
}

// selectClock Clock 算法（NRU - Not Recently Used）
// 408 考点：时钟算法，循环检查访问位
func (vms *VirtualMemorySimulator) selectClock() int {
	for {
		pageNum := vms.PhysicalMemory[vms.ClockPointer]

		if pageNum != -1 {
			if !vms.PageTable[pageNum].Referenced {
				// 找到访问位为 0 的页面
				victim := vms.ClockPointer
				vms.ClockPointer = (vms.ClockPointer + 1) % vms.NumFrames
				return victim
			}
			// 清除访问位，继续查找
			vms.PageTable[pageNum].Referenced = false
		}

		vms.ClockPointer = (vms.ClockPointer + 1) % vms.NumFrames
	}
}

// selectOptimal 最佳置换算法（OPT）：选择未来最长时间不访问的页面
// 408 考点：理论最优算法，实际无法实现
func (vms *VirtualMemorySimulator) selectOptimal() int {
	maxFutureDistance := -1
	victim := 0

	for i, pageNum := range vms.PhysicalMemory {
		if pageNum == -1 {
			continue
		}

		// 查找该页面在未来何时被访问
		futureDistance := vms.findNextAccess(pageNum)

		if futureDistance > maxFutureDistance {
			maxFutureDistance = futureDistance
			victim = i
		}
	}

	return victim
}

// findNextAccess 查找页面在未来的访问距离
func (vms *VirtualMemorySimulator) findNextAccess(pageNumber int) int {
	for i := vms.CurrentAccessIndex; i < len(vms.FutureAccesses); i++ {
		virtualAddr := vms.FutureAccesses[i]
		pageNum, _ := vms.parseVirtualAddress(virtualAddr)
		if pageNum == pageNumber {
			return i - vms.CurrentAccessIndex
		}
	}
	// 如果未来不再访问，返回最大值
	return len(vms.FutureAccesses) + 1
}

// GetStatistics 获取统计信息
func (vms *VirtualMemorySimulator) GetStatistics() string {
	pageFaultRate := 0.0
	tlbHitRate := 0.0

	if vms.AccessCount > 0 {
		pageFaultRate = float64(vms.PageFaults) / float64(vms.AccessCount) * 100
		tlbHitRate = float64(vms.TLBHits) / float64(vms.AccessCount) * 100
	}

	return fmt.Sprintf(`
虚拟存储器统计信息：
  总访问次数:   %d
  缺页次数:     %d
  缺页率:       %.2f%%
  TLB 命中次数: %d
  TLB 缺失次数: %d
  TLB 命中率:   %.2f%%
  
配置信息：
  虚拟页数:     %d
  物理页框数:   %d
  TLB 大小:     %d
  页面大小:     %d 字节
  替换算法:     %s
`,
		vms.AccessCount, vms.PageFaults, pageFaultRate,
		vms.TLBHits, vms.TLBMisses, tlbHitRate,
		vms.NumPages, vms.NumFrames, vms.TLBSize, vms.PageSize,
		getPolicyNameVM(vms.Policy))
}

func getPolicyNameVM(policy PageReplacementPolicy) string {
	switch policy {
	case PageFIFO:
		return "FIFO (先进先出)"
	case PageLRU:
		return "LRU (最近最少使用)"
	case PageClock:
		return "Clock (时钟算法/NRU)"
	case PageOptimal:
		return "OPT (最佳置换)"
	default:
		return "未知"
	}
}

// VirtualMemoryExample 虚拟存储器示例
// 408 考点：用访问序列演示页面替换过程
func VirtualMemoryExample() {
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  虚拟存储器模拟")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 访问序列：典型的 408 考试序列
	// 页面引用串: 1, 2, 3, 4, 1, 2, 5, 1, 2, 3, 4, 5
	addresses := []int{
		0x1000, // 页 1
		0x2000, // 页 2
		0x3000, // 页 3
		0x4000, // 页 4
		0x1000, // 页 1
		0x2000, // 页 2
		0x5000, // 页 5
		0x1000, // 页 1
		0x2000, // 页 2
		0x3000, // 页 3
		0x4000, // 页 4
		0x5000, // 页 5
	}

	policies := []PageReplacementPolicy{PageFIFO, PageLRU, PageClock, PageOptimal}

	for _, policy := range policies {
		fmt.Printf("\n【%s 页面替换算法】\n", getPolicyNameVM(policy))
		fmt.Println("配置：虚拟页数=8, 物理页框=3, TLB大小=2")
		fmt.Println("────────────────────────────────────────")

		simulator := NewVirtualMemorySimulator(8, 3, 2, policy)

		// OPT 算法需要知道未来访问序列
		if policy == PageOptimal {
			simulator.SetFutureAccesses(addresses)
		}

		for i, addr := range addresses {
			fmt.Printf("\n第 %d 次访问：\n", i+1)
			result := simulator.Access(addr)
			fmt.Print(result)
		}

		fmt.Println(simulator.GetStatistics())
		fmt.Println("════════════════════════════════════════")
	}

	fmt.Println("\n" + virtualMemory408Summary())
}

// virtualMemory408Summary 408 考试总结
func virtualMemory408Summary() string {
	return `
╔════════════════════════════════════════════════════════════════╗
║                 408 考试要点总结 - 虚拟存储器                 ║
╠════════════════════════════════════════════════════════════════╣
║ 1. 地址转换过程：                                              ║
║    虚拟地址 → TLB 查找 → 页表查找 → 物理地址                  ║
║    • TLB 命中：直接得到物理页框号                             ║
║    • TLB 缺失但页表命中：从页表获取页框号，更新 TLB           ║
║    • 缺页中断：调入页面，可能需要页面替换                     ║
║                                                                ║
║ 2. 地址结构：                                                  ║
║    虚拟地址 = 虚拟页号 | 页内偏移                             ║
║    物理地址 = 物理页框号 | 页内偏移                           ║
║    • 页号位数 = log₂(虚拟页数)                                ║
║    • 页内偏移位数 = log₂(页面大小)                            ║
║                                                                ║
║ 3. 页表项结构：                                                ║
║    • 页框号：物理内存中的位置                                 ║
║    • 有效位(V)：页面是否在内存中                              ║
║    • 修改位(M)：页面是否被修改过（脏位）                      ║
║    • 访问位(A)：页面是否被访问过（用于 Clock）                ║
║                                                                ║
║ 4. 页面替换算法比较：                                          ║
║    • FIFO：实现简单，可能产生 Belady 异常                     ║
║    • LRU：性能好，硬件开销大，不会产生 Belady 异常            ║
║    • Clock：LRU 的近似，实现简单，性能接近 LRU                ║
║    • OPT：理论最优，实际无法实现，用于性能评估                ║
║                                                                ║
║ 5. 缺页率计算：                                                ║
║    缺页率 = 缺页次数 / 总访问次数 × 100%                      ║
║                                                                ║
║ 6. 有效访问时间（EAT）：                                       ║
║    EAT = (1-p) × 内存访问时间 + p × 缺页处理时间              ║
║    其中 p 为缺页率                                             ║
║                                                                ║
║ 7. TLB 对性能的影响：                                          ║
║    • TLB 命中：1 次内存访问                                   ║
║    • TLB 缺失但页在内存：2 次内存访问（页表+数据）            ║
║    • 缺页：需要磁盘 I/O，时间远大于内存访问                   ║
║                                                                ║
║ 8. 典型考题：                                                  ║
║    • 给定页面引用串，画出各算法的页面调度过程                 ║
║    • 计算不同算法的缺页次数和缺页率                           ║
║    • 分析 Belady 异常现象                                     ║
║    • 计算页表大小和地址转换                                   ║
╚════════════════════════════════════════════════════════════════╝
`
}
