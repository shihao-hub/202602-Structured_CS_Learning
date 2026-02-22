package memory

import (
	"fmt"
)

// MemoryBlock 内存块
type MemoryBlock struct {
	Start     int  // 起始地址
	Size      int  // 大小
	Free      bool // 是否空闲
	ProcessID int  // 占用进程ID（-1表示空闲）
}

// MemoryManager 内存管理器
type MemoryManager struct {
	TotalSize int            // 总内存大小
	Blocks    []*MemoryBlock // 内存块列表
}

// NewMemoryManager 创建内存管理器
func NewMemoryManager(totalSize int) *MemoryManager {
	mm := &MemoryManager{
		TotalSize: totalSize,
		Blocks:    make([]*MemoryBlock, 0),
	}
	// 初始化为一个大的空闲块
	mm.Blocks = append(mm.Blocks, &MemoryBlock{
		Start:     0,
		Size:      totalSize,
		Free:      true,
		ProcessID: -1,
	})
	return mm
}

// FirstFit 首次适应算法
func (mm *MemoryManager) FirstFit(processID, size int) bool {
	for i, block := range mm.Blocks {
		if block.Free && block.Size >= size {
			return mm.allocateBlock(i, processID, size)
		}
	}
	return false
}

// BestFit 最佳适应算法
func (mm *MemoryManager) BestFit(processID, size int) bool {
	bestIdx := -1
	bestSize := mm.TotalSize + 1

	for i, block := range mm.Blocks {
		if block.Free && block.Size >= size && block.Size < bestSize {
			bestIdx = i
			bestSize = block.Size
		}
	}

	if bestIdx != -1 {
		return mm.allocateBlock(bestIdx, processID, size)
	}
	return false
}

// WorstFit 最差适应算法
func (mm *MemoryManager) WorstFit(processID, size int) bool {
	worstIdx := -1
	worstSize := -1

	for i, block := range mm.Blocks {
		if block.Free && block.Size >= size && block.Size > worstSize {
			worstIdx = i
			worstSize = block.Size
		}
	}

	if worstIdx != -1 {
		return mm.allocateBlock(worstIdx, processID, size)
	}
	return false
}

// allocateBlock 分配内存块
func (mm *MemoryManager) allocateBlock(idx, processID, size int) bool {
	block := mm.Blocks[idx]

	if block.Size == size {
		// 刚好合适
		block.Free = false
		block.ProcessID = processID
	} else {
		// 分割块
		newBlock := &MemoryBlock{
			Start:     block.Start + size,
			Size:      block.Size - size,
			Free:      true,
			ProcessID: -1,
		}
		block.Size = size
		block.Free = false
		block.ProcessID = processID

		// 插入新块
		newBlocks := make([]*MemoryBlock, 0, len(mm.Blocks)+1)
		newBlocks = append(newBlocks, mm.Blocks[:idx+1]...)
		newBlocks = append(newBlocks, newBlock)
		newBlocks = append(newBlocks, mm.Blocks[idx+1:]...)
		mm.Blocks = newBlocks
	}

	return true
}

// Free 释放内存
func (mm *MemoryManager) Free(processID int) bool {
	freed := false
	for _, block := range mm.Blocks {
		if block.ProcessID == processID {
			block.Free = true
			block.ProcessID = -1
			freed = true
		}
	}

	if freed {
		mm.compact()
	}
	return freed
}

// compact 合并相邻的空闲块
func (mm *MemoryManager) compact() {
	if len(mm.Blocks) <= 1 {
		return
	}

	newBlocks := make([]*MemoryBlock, 0)
	current := mm.Blocks[0]

	for i := 1; i < len(mm.Blocks); i++ {
		next := mm.Blocks[i]
		if current.Free && next.Free {
			// 合并
			current.Size += next.Size
		} else {
			newBlocks = append(newBlocks, current)
			current = next
		}
	}
	newBlocks = append(newBlocks, current)
	mm.Blocks = newBlocks
}

// GetFreeMemory 获取空闲内存总量
func (mm *MemoryManager) GetFreeMemory() int {
	total := 0
	for _, block := range mm.Blocks {
		if block.Free {
			total += block.Size
		}
	}
	return total
}

// GetUsedMemory 获取已用内存总量
func (mm *MemoryManager) GetUsedMemory() int {
	return mm.TotalSize - mm.GetFreeMemory()
}

// GetFragmentation 获取碎片数（空闲块数量）
func (mm *MemoryManager) GetFragmentation() int {
	count := 0
	for _, block := range mm.Blocks {
		if block.Free {
			count++
		}
	}
	return count
}

// Print 打印内存状态
func (mm *MemoryManager) Print() {
	fmt.Printf("内存状态 (总大小: %d, 已用: %d, 空闲: %d, 碎片数: %d):\n",
		mm.TotalSize, mm.GetUsedMemory(), mm.GetFreeMemory(), mm.GetFragmentation())
	fmt.Println("地址    大小    状态      进程")
	fmt.Println("-" + "------" + "------" + "--------" + "------")
	for _, block := range mm.Blocks {
		status := "空闲"
		process := "-"
		if !block.Free {
			status = "已分配"
			process = fmt.Sprintf("P%d", block.ProcessID)
		}
		fmt.Printf("%-7d %-7d %-9s %s\n", block.Start, block.Size, status, process)
	}
	fmt.Println()
}

// Page 页面结构
type Page struct {
	PageNumber  int  // 页号
	FrameNumber int  // 帧号（物理内存）
	Valid       bool // 有效位
	Referenced  bool // 引用位
	Modified    bool // 修改位
}

// PageTable 页表
type PageTable struct {
	Pages    []*Page // 页表项
	PageSize int     // 页大小
}

// NewPageTable 创建页表
func NewPageTable(numPages, pageSize int) *PageTable {
	pt := &PageTable{
		Pages:    make([]*Page, numPages),
		PageSize: pageSize,
	}
	for i := 0; i < numPages; i++ {
		pt.Pages[i] = &Page{
			PageNumber:  i,
			FrameNumber: -1,
			Valid:       false,
			Referenced:  false,
			Modified:    false,
		}
	}
	return pt
}

// MapPage 映射页面到帧
func (pt *PageTable) MapPage(pageNum, frameNum int) {
	if pageNum >= 0 && pageNum < len(pt.Pages) {
		pt.Pages[pageNum].FrameNumber = frameNum
		pt.Pages[pageNum].Valid = true
	}
}

// TranslateAddress 地址转换（逻辑地址 -> 物理地址）
func (pt *PageTable) TranslateAddress(logicalAddr int) (int, bool) {
	pageNum := logicalAddr / pt.PageSize
	offset := logicalAddr % pt.PageSize

	if pageNum < 0 || pageNum >= len(pt.Pages) {
		return -1, false
	}

	page := pt.Pages[pageNum]
	if !page.Valid {
		return -1, false // 页面错误
	}

	page.Referenced = true
	physicalAddr := page.FrameNumber*pt.PageSize + offset
	return physicalAddr, true
}

// Print 打印页表
func (pt *PageTable) Print() {
	fmt.Printf("页表 (页大小: %d):\n", pt.PageSize)
	fmt.Println("页号    帧号    有效    引用    修改")
	fmt.Println("-" + "------" + "------" + "------" + "------" + "------")
	for _, page := range pt.Pages {
		frameStr := "-"
		if page.Valid {
			frameStr = fmt.Sprintf("%d", page.FrameNumber)
		}
		fmt.Printf("%-7d %-7s %-7t %-7t %t\n",
			page.PageNumber, frameStr, page.Valid, page.Referenced, page.Modified)
	}
	fmt.Println()
}

// MemoryExample 内存管理示例
func MemoryExample() {
	fmt.Println("=== 内存管理 (Memory Management) 示例 ===")

	// 首次适应算法
	fmt.Println("\n1. 首次适应算法 (First Fit):")
	mm1 := NewMemoryManager(1000)
	mm1.FirstFit(1, 200)
	mm1.FirstFit(2, 150)
	mm1.FirstFit(3, 300)
	mm1.Print()

	fmt.Println("释放进程2:")
	mm1.Free(2)
	mm1.Print()

	fmt.Println("分配进程4 (100):")
	mm1.FirstFit(4, 100)
	mm1.Print()

	// 最佳适应算法
	fmt.Println("\n2. 最佳适应算法 (Best Fit):")
	mm2 := NewMemoryManager(1000)
	mm2.BestFit(1, 200)
	mm2.BestFit(2, 150)
	mm2.BestFit(3, 300)
	mm2.Free(2)
	mm2.Print()

	fmt.Println("使用最佳适应分配进程4 (100):")
	mm2.BestFit(4, 100)
	mm2.Print()

	// 最差适应算法
	fmt.Println("\n3. 最差适应算法 (Worst Fit):")
	mm3 := NewMemoryManager(1000)
	mm3.WorstFit(1, 200)
	mm3.WorstFit(2, 150)
	mm3.WorstFit(3, 300)
	mm3.Free(2)
	mm3.Print()

	fmt.Println("使用最差适应分配进程4 (100):")
	mm3.WorstFit(4, 100)
	mm3.Print()

	// 分页系统
	fmt.Println("\n4. 分页系统:")
	pageTable := NewPageTable(8, 1024) // 8页，每页1024字节

	// 映射一些页面
	pageTable.MapPage(0, 5)
	pageTable.MapPage(1, 2)
	pageTable.MapPage(3, 7)
	pageTable.MapPage(5, 1)
	pageTable.Print()

	// 地址转换
	fmt.Println("地址转换示例:")
	testAddresses := []int{0, 1024, 2048, 3072, 5120, 7168}
	for _, addr := range testAddresses {
		physAddr, ok := pageTable.TranslateAddress(addr)
		if ok {
			fmt.Printf("  逻辑地址 %d -> 物理地址 %d\n", addr, physAddr)
		} else {
			fmt.Printf("  逻辑地址 %d -> 页面错误!\n", addr)
		}
	}

	fmt.Println("\n5. 内存分配算法比较:")
	fmt.Println("  - 首次适应 (First Fit): 从头开始找第一个合适的块，速度快")
	fmt.Println("  - 最佳适应 (Best Fit): 找最小的合适块，减少大块碎片")
	fmt.Println("  - 最差适应 (Worst Fit): 找最大的块，减少小碎片")
	fmt.Println()
}

// RunAllMemoryExamples 运行所有内存管理相关的示例
func RunAllMemoryExamples() {
	MemoryExample()
}
