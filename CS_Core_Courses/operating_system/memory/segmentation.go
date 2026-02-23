package memory

import (
	"fmt"
)

// Segment 段结构
// 408考点：段式存储管理的基本单位
type Segment struct {
	SegmentNumber int    // 段号
	Base          int    // 基址（段在物理内存中的起始地址）
	Limit         int    // 段长（段的大小）
	Name          string // 段名（如代码段、数据段、栈段）
}

// SegmentTable 段表
// 408考点：段号到物理地址的映射表
type SegmentTable struct {
	Segments []*Segment // 段表项
}

// NewSegmentTable 创建段表
func NewSegmentTable() *SegmentTable {
	return &SegmentTable{
		Segments: make([]*Segment, 0),
	}
}

// AddSegment 添加段
func (st *SegmentTable) AddSegment(segNum, base, limit int, name string) {
	segment := &Segment{
		SegmentNumber: segNum,
		Base:          base,
		Limit:         limit,
		Name:          name,
	}
	st.Segments = append(st.Segments, segment)
}

// GetSegment 获取段信息
func (st *SegmentTable) GetSegment(segNum int) (*Segment, bool) {
	for _, seg := range st.Segments {
		if seg.SegmentNumber == segNum {
			return seg, true
		}
	}
	return nil, false
}

// TranslateAddress 段式地址转换
// 408考点：逻辑地址（段号，段内偏移）-> 物理地址
// 物理地址 = 基址 + 段内偏移（需检查段内偏移 < 段长）
func (st *SegmentTable) TranslateAddress(segNum, offset int) (int, error) {
	// 查段表
	segment, ok := st.GetSegment(segNum)
	if !ok {
		return -1, fmt.Errorf("段号 %d 不存在", segNum)
	}

	// 检查段内偏移是否越界
	if offset < 0 || offset >= segment.Limit {
		return -1, fmt.Errorf("段内偏移 %d 越界（段长为 %d）", offset, segment.Limit)
	}

	// 计算物理地址
	physicalAddr := segment.Base + offset
	return physicalAddr, nil
}

// TranslateAddressDetailed 详细的段式地址转换（用于教学演示）
func (st *SegmentTable) TranslateAddressDetailed(segNum, offset int) {
	fmt.Printf("\n【段式地址转换】\n")
	fmt.Printf("逻辑地址: (段号=%d, 段内偏移=%d)\n", segNum, offset)

	// 查段表
	segment, ok := st.GetSegment(segNum)
	if !ok {
		fmt.Printf("└─ ❌ 段号 %d 不存在，访问违例\n", segNum)
		return
	}

	fmt.Printf("├─ 段信息: %s\n", segment.Name)
	fmt.Printf("├─ 基址: %d (0x%X)\n", segment.Base, segment.Base)
	fmt.Printf("├─ 段长: %d 字节\n", segment.Limit)

	// 检查越界
	if offset < 0 || offset >= segment.Limit {
		fmt.Printf("└─ ❌ 段内偏移越界！(偏移=%d, 段长=%d)\n", offset, segment.Limit)
		return
	}

	// 计算物理地址
	physicalAddr := segment.Base + offset
	fmt.Printf("├─ 物理地址计算: %d + %d = %d\n", segment.Base, offset, physicalAddr)
	fmt.Printf("└─ ✓ 物理地址: %d (0x%X)\n", physicalAddr, physicalAddr)
}

// Print 打印段表
func (st *SegmentTable) Print() {
	fmt.Println("\n段表内容:")
	fmt.Printf("%-6s %-12s %-12s %-12s\n", "段号", "段名", "基址", "段长")
	fmt.Println("─────────────────────────────────────────────────")
	for _, seg := range st.Segments {
		fmt.Printf("%-6d %-12s %-12d %-12d\n",
			seg.SegmentNumber, seg.Name, seg.Base, seg.Limit)
	}
	fmt.Println()
}

// SegmentPageSystem 段页式存储管理系统
// 408考点：段页式结合，先分段再分页
type SegmentPageSystem struct {
	SegmentTable      *SegmentTable           // 段表
	PageTablesPerSeg  map[int]*PageTable      // 每个段的页表 [段号]->页表
	PageSize          int                     // 页面大小
}

// NewSegmentPageSystem 创建段页式系统
func NewSegmentPageSystem(pageSize int) *SegmentPageSystem {
	return &SegmentPageSystem{
		SegmentTable:     NewSegmentTable(),
		PageTablesPerSeg: make(map[int]*PageTable),
		PageSize:         pageSize,
	}
}

// AddSegmentWithPages 添加段（带分页）
func (sps *SegmentPageSystem) AddSegmentWithPages(segNum int, segName string, segSize int) {
	// 计算需要多少页
	numPages := (segSize + sps.PageSize - 1) / sps.PageSize

	// 添加段（基址在段页式中可以是段页表的地址，这里简化为0）
	sps.SegmentTable.AddSegment(segNum, 0, segSize, segName)

	// 为该段创建页表
	sps.PageTablesPerSeg[segNum] = NewPageTable(numPages, sps.PageSize)
}

// TranslateSegmentPageAddress 段页式地址转换
// 408考点：两级转换
// 1. 段号 -> 查段表 -> 得到该段的页表基址
// 2. 段内偏移 -> 分解为页号和页内偏移
// 3. 页号 -> 查页表 -> 得到帧号
// 4. 物理地址 = 帧号 × 页大小 + 页内偏移
func (sps *SegmentPageSystem) TranslateSegmentPageAddress(segNum, offset int) {
	fmt.Printf("\n【段页式地址转换】\n")
	fmt.Printf("逻辑地址: (段号=%d, 段内偏移=%d)\n", segNum, offset)

	// 1. 查段表
	segment, ok := sps.SegmentTable.GetSegment(segNum)
	if !ok {
		fmt.Printf("└─ ❌ 段号 %d 不存在\n", segNum)
		return
	}
	fmt.Printf("├─ 段信息: %s (段长=%d)\n", segment.Name, segment.Limit)

	// 2. 检查段内偏移是否越界
	if offset < 0 || offset >= segment.Limit {
		fmt.Printf("└─ ❌ 段内偏移越界\n")
		return
	}

	// 3. 获取该段的页表
	pageTable, ok := sps.PageTablesPerSeg[segNum]
	if !ok {
		fmt.Printf("└─ ❌ 段 %d 的页表不存在\n", segNum)
		return
	}

	// 4. 段内偏移分解为页号和页内偏移
	pageNum := offset / sps.PageSize
	pageOffset := offset % sps.PageSize
	fmt.Printf("├─ 段内分页: 页号=%d, 页内偏移=%d\n", pageNum, pageOffset)

	// 5. 查页表
	if pageNum >= len(pageTable.Pages) {
		fmt.Printf("└─ ❌ 页号越界\n")
		return
	}

	page := pageTable.Pages[pageNum]
	if !page.Valid {
		fmt.Printf("└─ ❌ 页表项无效，产生缺页中断\n")
		return
	}

	frameNum := page.FrameNumber
	fmt.Printf("├─ 帧号: %d\n", frameNum)

	// 6. 计算物理地址
	physicalAddr := frameNum*sps.PageSize + pageOffset
	fmt.Printf("└─ ✓ 物理地址: %d (0x%X)\n", physicalAddr, physicalAddr)
}

// SegmentationExample 分段存储管理示例
func SegmentationExample() {
	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║        操作系统 - 分段存储管理 (Segmentation System)     ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")

	// ============ 示例1: 纯段式存储管理 ============
	fmt.Println("\n【示例1】纯段式存储管理")
	fmt.Println("─────────────────────────────────────────────────")

	// 创建段表
	st := NewSegmentTable()

	// 添加段（模拟一个进程的内存布局）
	st.AddSegment(0, 1000, 5000, "代码段 (Code)")
	st.AddSegment(1, 6000, 3000, "数据段 (Data)")
	st.AddSegment(2, 9000, 2000, "栈段 (Stack)")
	st.AddSegment(3, 11000, 1000, "堆段 (Heap)")

	st.Print()

	// 地址转换示例
	fmt.Println("地址转换示例:")

	// 访问代码段中的指令
	st.TranslateAddressDetailed(0, 100) // 段0，偏移100

	// 访问数据段中的数据
	st.TranslateAddressDetailed(1, 500) // 段1，偏移500

	// 访问栈段
	st.TranslateAddressDetailed(2, 100) // 段2，偏移100

	// 越界访问
	st.TranslateAddressDetailed(1, 5000) // 段1的段长只有3000，越界

	// 不存在的段
	st.TranslateAddressDetailed(5, 100) // 段5不存在

	// ============ 示例2: 段页式存储管理 ============
	fmt.Println("\n\n【示例2】段页式存储管理")
	fmt.Println("─────────────────────────────────────────────────")

	pageSize := 1024 // 页大小1KB
	sps := NewSegmentPageSystem(pageSize)

	fmt.Printf("系统配置: 页面大小 = %d 字节\n\n", pageSize)

	// 添加段（自动分页）
	sps.AddSegmentWithPages(0, "代码段", 5000)  // 需要5页
	sps.AddSegmentWithPages(1, "数据段", 3000)  // 需要3页
	sps.AddSegmentWithPages(2, "栈段", 2000)   // 需要2页

	// 为各段的页建立映射（模拟加载到物理内存）
	// 代码段：页0->帧10, 页1->帧11, 页2->帧12, 页3->帧13, 页4->帧14
	for i := 0; i < 5; i++ {
		sps.PageTablesPerSeg[0].MapPage(i, 10+i)
	}

	// 数据段：页0->帧20, 页1->帧21, 页2->帧22
	for i := 0; i < 3; i++ {
		sps.PageTablesPerSeg[1].MapPage(i, 20+i)
	}

	// 栈段：页0->帧30, 页1->帧31
	for i := 0; i < 2; i++ {
		sps.PageTablesPerSeg[2].MapPage(i, 30+i)
	}

	// 打印段表
	sps.SegmentTable.Print()

	// 地址转换示例
	fmt.Println("段页式地址转换示例:")

	// 访问代码段第2页内的某个位置
	sps.TranslateSegmentPageAddress(0, 2100) // 段0, 偏移2100 = 页2, 页内偏移52

	// 访问数据段第1页
	sps.TranslateSegmentPageAddress(1, 1500) // 段1, 偏移1500 = 页1, 页内偏移476

	// 访问栈段
	sps.TranslateSegmentPageAddress(2, 500) // 段2, 偏移500 = 页0, 页内偏移500

	// ============ 示例3: 分段与分页的对比 ============
	fmt.Println("\n\n【示例3】分段 vs 分页 vs 段页式")
	fmt.Println("─────────────────────────────────────────────────")

	fmt.Println("\n分段存储管理 (Segmentation):")
	fmt.Println("  特点:")
	fmt.Println("    • 按程序逻辑结构划分（代码、数据、栈等）")
	fmt.Println("    • 段的大小不固定，由逻辑需求决定")
	fmt.Println("    • 逻辑地址 = (段号, 段内偏移)")
	fmt.Println("  优点:")
	fmt.Println("    ✓ 符合程序逻辑结构，便于编程和共享")
	fmt.Println("    ✓ 段内地址连续，访问局部性好")
	fmt.Println("    ✓ 易于实现保护和共享（按段设置权限）")
	fmt.Println("  缺点:")
	fmt.Println("    ✗ 产生外部碎片（段大小不一）")
	fmt.Println("    ✗ 段太大时，分配困难")

	fmt.Println("\n分页存储管理 (Paging):")
	fmt.Println("  特点:")
	fmt.Println("    • 按物理大小均匀划分")
	fmt.Println("    • 页的大小固定（如4KB）")
	fmt.Println("    • 逻辑地址 = (页号, 页内偏移)")
	fmt.Println("  优点:")
	fmt.Println("    ✓ 无外部碎片")
	fmt.Println("    ✓ 易于管理和分配")
	fmt.Println("    ✓ 支持虚拟内存")
	fmt.Println("  缺点:")
	fmt.Println("    ✗ 有内部碎片（最后一页未用满）")
	fmt.Println("    ✗ 不符合程序逻辑结构")

	fmt.Println("\n段页式存储管理 (Segmentation + Paging):")
	fmt.Println("  特点:")
	fmt.Println("    • 结合两者优点：先分段，段内再分页")
	fmt.Println("    • 逻辑地址 = (段号, 段内偏移)")
	fmt.Println("    • 段内偏移 = (页号, 页内偏移)")
	fmt.Println("  优点:")
	fmt.Println("    ✓ 既符合逻辑结构，又无外部碎片")
	fmt.Println("    ✓ 便于共享和保护")
	fmt.Println("  缺点:")
	fmt.Println("    ✗ 地址转换复杂（需查两次表）")
	fmt.Println("    ✗ 管理开销大")

	// ============ 408考点总结 ============
	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    408 考点总结                            ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")

	fmt.Println("\n1. 分段存储管理核心概念:")
	fmt.Println("   • 段 (Segment): 程序按逻辑功能划分的单位")
	fmt.Println("   • 段表 (Segment Table): 段号到物理地址的映射")
	fmt.Println("   • 段表项内容: 段号、基址、段长、权限等")
	fmt.Println("   • 地址转换: 物理地址 = 基址 + 段内偏移")
	fmt.Println("   • 越界检查: 段内偏移必须 < 段长")

	fmt.Println("\n2. 段式地址转换过程:")
	fmt.Println("   ① 从逻辑地址获取段号和段内偏移")
	fmt.Println("   ② 用段号查段表，得到段的基址和段长")
	fmt.Println("   ③ 检查段内偏移是否越界（偏移 < 段长）")
	fmt.Println("   ④ 计算物理地址 = 基址 + 段内偏移")

	fmt.Println("\n3. 段页式地址转换过程:")
	fmt.Println("   ① 用段号查段表，得到该段的页表基址")
	fmt.Println("   ② 将段内偏移分解为页号和页内偏移")
	fmt.Println("   ③ 用页号查页表，得到帧号")
	fmt.Println("   ④ 计算物理地址 = 帧号 × 页大小 + 页内偏移")
	fmt.Println("   访存次数：至少3次（段表、页表、实际数据）")

	fmt.Println("\n4. 三种方式的比较:")
	fmt.Println("┌─────────┬──────────┬──────────┬──────────┐")
	fmt.Println("│ 比较项  │   分页   │   分段   │  段页式  │")
	fmt.Println("├─────────┼──────────┼──────────┼──────────┤")
	fmt.Println("│ 逻辑性  │ 不符合   │ 符合     │ 符合     │")
	fmt.Println("│ 外部碎片│ 无       │ 有       │ 无       │")
	fmt.Println("│ 内部碎片│ 有       │ 无       │ 有(小)   │")
	fmt.Println("│ 大小    │ 固定     │ 不固定   │ 页固定   │")
	fmt.Println("│ 共享保护│ 较难     │ 容易     │ 容易     │")
	fmt.Println("│ 地址转换│ 简单     │ 简单     │ 复杂     │")
	fmt.Println("└─────────┴──────────┴──────────┴──────────┘")

	fmt.Println("\n5. 考试常见题型:")
	fmt.Println("   • 给定段表，进行段式地址转换")
	fmt.Println("   • 判断段内偏移是否越界")
	fmt.Println("   • 比较分页、分段、段页式的优缺点")
	fmt.Println("   • 段页式系统的地址转换计算")
	fmt.Println("   • 分析外部碎片和内部碎片")
	fmt.Println()

	// ============ 计算示例 ============
	fmt.Println("\n【计算示例】408典型题目")
	fmt.Println("─────────────────────────────────────────────────")
	fmt.Println("题目：某段页式系统，段表和页表如下：")
	fmt.Println()
	fmt.Println("段表：")
	fmt.Println("  段号  页表起址  段长(页)")
	fmt.Println("  0     100      5")
	fmt.Println("  1     200      3")
	fmt.Println()
	fmt.Println("页表（段0）：")
	fmt.Println("  页号  帧号")
	fmt.Println("  0     10")
	fmt.Println("  1     15")
	fmt.Println("  2     20")
	fmt.Println()
	fmt.Println("若页大小为1KB，逻辑地址为 (1, 1500)，求物理地址。")
	fmt.Println("（1为段号，1500为段内偏移）")
	fmt.Println()

	fmt.Println("解答:")
	fmt.Println("① 用段号1查段表 -> 该段页表起址=200，段长=3页")
	fmt.Println("② 段内偏移1500 ÷ 1024 = 页号1，余数476（页内偏移）")
	fmt.Println("③ 在段1的页表中，页号1对应的帧号需要查表（此处未给出）")
	fmt.Println("④ 假设段1页号1的帧号为25，则物理地址 = 25 × 1024 + 476 = 26076")
	fmt.Println()
	fmt.Println("注意：实际考试会给出完整的页表信息。")
	fmt.Println()
}
