package memory

import (
	"fmt"
)

// PagingSystem 分页系统（扩展版）
// 408考点：分页存储管理的核心概念
type PagingSystem struct {
	PageSize     int        // 页面大小（字节）
	FrameSize    int        // 帧大小（与页面大小相同）
	NumPages     int        // 逻辑地址空间页数
	NumFrames    int        // 物理内存帧数
	PageTable    *PageTable // 页表
	MemoryFrames []bool     // 物理内存帧占用情况（true表示已占用）
	PageFaults   int        // 缺页次数
}

// NewPagingSystem 创建分页系统
func NewPagingSystem(pageSize, numPages, numFrames int) *PagingSystem {
	return &PagingSystem{
		PageSize:     pageSize,
		FrameSize:    pageSize,
		NumPages:     numPages,
		NumFrames:    numFrames,
		PageTable:    NewPageTable(numPages, pageSize),
		MemoryFrames: make([]bool, numFrames),
		PageFaults:   0,
	}
}

// AllocateFrame 分配物理帧
func (ps *PagingSystem) AllocateFrame() (int, bool) {
	for i := 0; i < ps.NumFrames; i++ {
		if !ps.MemoryFrames[i] {
			ps.MemoryFrames[i] = true
			return i, true
		}
	}
	return -1, false // 无空闲帧
}

// LoadPage 加载页面到内存（模拟缺页处理）
// 408考点：缺页中断处理过程
func (ps *PagingSystem) LoadPage(pageNum int) bool {
	if pageNum < 0 || pageNum >= ps.NumPages {
		return false
	}

	// 检查页面是否已在内存中
	if ps.PageTable.Pages[pageNum].Valid {
		return true
	}

	// 缺页，需要分配帧
	frameNum, ok := ps.AllocateFrame()
	if !ok {
		fmt.Printf("  ⚠️  内存已满，需要页面置换\n")
		return false
	}

	// 加载页面到帧
	ps.PageTable.MapPage(pageNum, frameNum)
	ps.PageFaults++
	fmt.Printf("  📄 加载页面 P%d 到帧 F%d（缺页中断 #%d）\n", pageNum, frameNum, ps.PageFaults)
	return true
}

// TranslateAddressDetailed 详细的地址转换过程（用于教学演示）
// 408考点：逻辑地址到物理地址的转换过程
func (ps *PagingSystem) TranslateAddressDetailed(logicalAddr int) {
	fmt.Printf("\n【地址转换过程】\n")
	fmt.Printf("逻辑地址: %d (0x%X)\n", logicalAddr, logicalAddr)

	// 1. 分解逻辑地址
	pageNum := logicalAddr / ps.PageSize
	offset := logicalAddr % ps.PageSize
	fmt.Printf("├─ 页号: %d\n", pageNum)
	fmt.Printf("├─ 页内偏移: %d (0x%X)\n", offset, offset)

	// 2. 检查页号是否合法
	if pageNum < 0 || pageNum >= ps.NumPages {
		fmt.Printf("└─ ❌ 页号越界，访问违例\n")
		return
	}

	// 3. 查页表
	page := ps.PageTable.Pages[pageNum]
	fmt.Printf("├─ 查页表: 页号 %d\n", pageNum)

	// 4. 检查有效位
	if !page.Valid {
		fmt.Printf("├─ ❌ 有效位=0，产生缺页中断\n")
		fmt.Printf("└─ 需要将页面从外存调入内存\n")
		return
	}

	// 5. 获取帧号
	frameNum := page.FrameNumber
	fmt.Printf("├─ 帧号: %d\n", frameNum)

	// 6. 计算物理地址
	physicalAddr := frameNum*ps.FrameSize + offset
	fmt.Printf("├─ 物理地址计算: %d × %d + %d = %d\n",
		frameNum, ps.FrameSize, offset, physicalAddr)
	fmt.Printf("└─ ✓ 物理地址: %d (0x%X)\n", physicalAddr, physicalAddr)
}

// MultiLevelPageTable 多级页表（二级页表示例）
// 408考点：多级页表的组织和地址转换
type MultiLevelPageTable struct {
	PageSize          int         // 页面大小
	Level1Size        int         // 一级页表大小（页目录表项数）
	Level2Size        int         // 二级页表大小（每个二级页表的表项数）
	PageDirectory     []int       // 页目录（存储二级页表的基址）
	SecondLevelTables [][]int     // 二级页表数组
	FrameAllocated    map[int]int // 已分配的帧映射 [页号]->帧号
}

// NewMultiLevelPageTable 创建二级页表
func NewMultiLevelPageTable(pageSize, level1Size, level2Size int) *MultiLevelPageTable {
	return &MultiLevelPageTable{
		PageSize:          pageSize,
		Level1Size:        level1Size,
		Level2Size:        level2Size,
		PageDirectory:     make([]int, level1Size),
		SecondLevelTables: make([][]int, level1Size),
		FrameAllocated:    make(map[int]int),
	}
}

// MapPageMultiLevel 多级页表映射
func (mpt *MultiLevelPageTable) MapPageMultiLevel(pageNum, frameNum int) {
	// 计算一级页表索引和二级页表索引
	level1Index := pageNum / mpt.Level2Size
	level2Index := pageNum % mpt.Level2Size

	// 如果二级页表不存在，创建它
	if mpt.SecondLevelTables[level1Index] == nil {
		mpt.SecondLevelTables[level1Index] = make([]int, mpt.Level2Size)
		for i := range mpt.SecondLevelTables[level1Index] {
			mpt.SecondLevelTables[level1Index][i] = -1 // 初始化为无效
		}
	}

	// 设置映射
	mpt.SecondLevelTables[level1Index][level2Index] = frameNum
	mpt.FrameAllocated[pageNum] = frameNum
}

// TranslateMultiLevel 多级页表地址转换
func (mpt *MultiLevelPageTable) TranslateMultiLevel(logicalAddr int) {
	fmt.Printf("\n【二级页表地址转换】\n")
	fmt.Printf("逻辑地址: %d (0x%X)\n", logicalAddr, logicalAddr)

	// 1. 分解逻辑地址
	pageNum := logicalAddr / mpt.PageSize
	offset := logicalAddr % mpt.PageSize
	fmt.Printf("├─ 页号: %d, 页内偏移: %d\n", pageNum, offset)

	// 2. 分解页号为两级索引
	level1Index := pageNum / mpt.Level2Size
	level2Index := pageNum % mpt.Level2Size
	fmt.Printf("├─ 一级页表索引: %d\n", level1Index)
	fmt.Printf("├─ 二级页表索引: %d\n", level2Index)

	// 3. 查一级页表（页目录）
	if level1Index >= mpt.Level1Size || mpt.SecondLevelTables[level1Index] == nil {
		fmt.Printf("└─ ❌ 一级页表项无效或不存在\n")
		return
	}

	// 4. 查二级页表
	frameNum := mpt.SecondLevelTables[level1Index][level2Index]
	if frameNum == -1 {
		fmt.Printf("└─ ❌ 二级页表项无效，产生缺页中断\n")
		return
	}

	fmt.Printf("├─ 帧号: %d\n", frameNum)

	// 5. 计算物理地址
	physicalAddr := frameNum*mpt.PageSize + offset
	fmt.Printf("└─ ✓ 物理地址: %d (0x%X)\n", physicalAddr, physicalAddr)
}

// PagingExample 分页机制示例（扩展版）
func PagingExample() {
	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║          操作系统 - 分页存储管理 (Paging System)         ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")

	// ============ 示例1: 基本分页系统 ============
	fmt.Println("\n【示例1】基本分页系统与地址转换")
	fmt.Println("─────────────────────────────────────────────────")

	// 创建分页系统：页大小4KB，8个逻辑页，6个物理帧
	pageSize := 4096 // 4KB
	ps := NewPagingSystem(pageSize, 8, 6)

	fmt.Printf("系统配置:\n")
	fmt.Printf("  页面大小: %d 字节 (4KB)\n", pageSize)
	fmt.Printf("  逻辑页数: %d\n", ps.NumPages)
	fmt.Printf("  物理帧数: %d\n", ps.NumFrames)
	fmt.Println()

	// 加载一些页面
	fmt.Println("加载页面到内存:")
	ps.LoadPage(0)
	ps.LoadPage(1)
	ps.LoadPage(3)
	ps.LoadPage(5)

	// 详细演示地址转换过程
	fmt.Println("\n【地址转换示例】")

	// 示例1: 访问逻辑地址 4096（页1的起始地址）
	ps.TranslateAddressDetailed(4096)

	// 示例2: 访问逻辑地址 12300（页3内的某个地址）
	ps.TranslateAddressDetailed(12300)

	// 示例3: 访问未加载的页面
	ps.TranslateAddressDetailed(8192) // 页2，未加载

	// 统计信息
	fmt.Printf("\n系统统计:\n")
	fmt.Printf("  缺页次数: %d\n", ps.PageFaults)
	fmt.Printf("  内存利用率: %.2f%%\n", float64(ps.PageFaults)/float64(ps.NumFrames)*100)

	// ============ 示例2: 多级页表 ============
	fmt.Println("\n\n【示例2】二级页表系统")
	fmt.Println("─────────────────────────────────────────────────")

	// 创建二级页表：页大小4KB，一级页表4项，每个二级页表16项
	// 可表示: 4 × 16 = 64个页面
	mpt := NewMultiLevelPageTable(4096, 4, 16)

	fmt.Printf("二级页表配置:\n")
	fmt.Printf("  页面大小: %d 字节\n", mpt.PageSize)
	fmt.Printf("  一级页表项数（页目录）: %d\n", mpt.Level1Size)
	fmt.Printf("  二级页表项数: %d\n", mpt.Level2Size)
	fmt.Printf("  可表示页数: %d\n", mpt.Level1Size*mpt.Level2Size)
	fmt.Println()

	// 建立映射
	fmt.Println("建立页表映射:")
	mappings := map[int]int{
		0:  10, // 页0 -> 帧10
		1:  5,  // 页1 -> 帧5
		16: 8,  // 页16（第二个二级页表的第一项）-> 帧8
		17: 12, // 页17 -> 帧12
		32: 3,  // 页32（第三个二级页表的第一项）-> 帧3
	}

	for page, frame := range mappings {
		mpt.MapPageMultiLevel(page, frame)
		fmt.Printf("  页 P%-2d -> 帧 F%-2d\n", page, frame)
	}

	// 地址转换演示
	fmt.Println("\n地址转换示例:")

	// 访问页0内的地址
	mpt.TranslateMultiLevel(100) // 页0, 偏移100

	// 访问页16内的地址（跨越到第二个二级页表）
	mpt.TranslateMultiLevel(65536 + 200) // 页16, 偏移200

	// 访问未映射的页
	mpt.TranslateMultiLevel(8192) // 页2，未映射

	// ============ 示例3: 分页机制的优缺点 ============
	fmt.Println("\n\n【示例3】分页机制分析")
	fmt.Println("─────────────────────────────────────────────────")

	fmt.Println("分页存储管理的特点:")
	fmt.Println()
	fmt.Println("优点:")
	fmt.Println("  ✓ 无外部碎片：页面大小固定，物理内存可以充分利用")
	fmt.Println("  ✓ 不需要连续分配：逻辑空间连续，物理空间可以不连续")
	fmt.Println("  ✓ 支持虚拟内存：可以运行大于物理内存的程序")
	fmt.Println("  ✓ 易于共享和保护：以页为单位进行管理")
	fmt.Println()
	fmt.Println("缺点:")
	fmt.Println("  ✗ 存在内部碎片：最后一页可能未被完全利用")
	fmt.Println("  ✗ 页表占用空间：大地址空间需要大页表")
	fmt.Println("  ✗ 地址转换开销：每次访存需要查页表")
	fmt.Println()
	fmt.Println("优化技术:")
	fmt.Println("  • TLB (Translation Lookaside Buffer): 快表，缓存页表项")
	fmt.Println("  • 多级页表: 减少页表占用空间")
	fmt.Println("  • 反置页表: 以物理帧为索引，减少空间开销")

	// ============ 408考点总结 ============
	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    408 考点总结                            ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")

	fmt.Println("\n1. 分页存储管理基本概念:")
	fmt.Println("   • 页 (Page): 逻辑地址空间划分的单位")
	fmt.Println("   • 帧/页框 (Frame): 物理地址空间划分的单位")
	fmt.Println("   • 页表 (Page Table): 页到帧的映射表")
	fmt.Println("   • 页表项 (PTE): 包含帧号、有效位、访问位、修改位等")

	fmt.Println("\n2. 地址转换过程:")
	fmt.Println("   逻辑地址 = 页号 + 页内偏移")
	fmt.Println("   页号 = 逻辑地址 ÷ 页面大小")
	fmt.Println("   页内偏移 = 逻辑地址 % 页面大小")
	fmt.Println("   物理地址 = 帧号 × 页面大小 + 页内偏移")

	fmt.Println("\n3. 页表的组织方式:")
	fmt.Println("   • 单级页表: 简单但页表可能很大")
	fmt.Println("   • 多级页表: 节省空间，但增加访存次数")
	fmt.Println("   • 反置页表: 页表大小与物理内存相关，需要哈希查找")
	fmt.Println("   • 快表 (TLB): 高速缓存，减少页表访问")

	fmt.Println("\n4. 页面大小的影响:")
	fmt.Println("   • 页面太大: 内部碎片增加")
	fmt.Println("   • 页面太小: 页表增大，管理开销增加")
	fmt.Println("   • 常用大小: 4KB, 8KB（需权衡）")

	fmt.Println("\n5. 考试常见题型:")
	fmt.Println("   • 给定逻辑地址，计算物理地址")
	fmt.Println("   • 给定页表和页面大小，进行地址转换")
	fmt.Println("   • 计算页表占用的空间大小")
	fmt.Println("   • 分析内部碎片和外部碎片")
	fmt.Println("   • 比较单级页表和多级页表")
	fmt.Println()

	// ============ 计算示例 ============
	fmt.Println("\n【计算示例】408典型题目")
	fmt.Println("─────────────────────────────────────────────────")
	fmt.Println("题目：某系统采用分页存储管理，逻辑地址32位，页面大小4KB，")
	fmt.Println("      页表项大小4字节。求：")
	fmt.Println("      (1) 逻辑地址空间最多有多少页？")
	fmt.Println("      (2) 页表最多需要多少字节？")
	fmt.Println("      (3) 采用二级页表，页目录和二级页表各有多少项？")
	fmt.Println()

	// 计算
	logicalAddrBits := 32
	pteSize := 4

	totalPages := 1 << (logicalAddrBits - 12) // 2^(32-12) = 2^20
	pageTableSize := totalPages * pteSize     // 页表大小

	fmt.Println("解答:")
	fmt.Printf("(1) 页面大小4KB = 2^12字节，页内偏移需要12位\n")
	fmt.Printf("    页号占 32 - 12 = 20位\n")
	fmt.Printf("    最多页数 = 2^20 = %d 页\n", totalPages)
	fmt.Println()
	fmt.Printf("(2) 页表项数 = %d\n", totalPages)
	fmt.Printf("    页表大小 = %d × %d = %d 字节 = %d MB\n",
		totalPages, pteSize, pageTableSize, pageTableSize/(1024*1024))
	fmt.Println()
	fmt.Printf("(3) 采用二级页表，每个页表页大小也是4KB\n")
	fmt.Printf("    每个页表页可容纳页表项数 = 4096 ÷ 4 = 1024 = 2^10\n")
	fmt.Printf("    二级页表索引需要10位，一级页表索引需要 20 - 10 = 10位\n")
	fmt.Printf("    页目录项数 = 2^10 = 1024\n")
	fmt.Printf("    每个二级页表项数 = 2^10 = 1024\n")
	fmt.Println()
}
