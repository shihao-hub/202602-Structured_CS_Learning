package filesystem

import "fmt"

// ============================================================
// 文件分配方式
// 408考点：连续分配、链接分配、索引分配的优缺点
// ============================================================

// --- 1. 连续分配 ---
// 文件占据磁盘上连续的块
// 优点：顺序/随机访问快
// 缺点：产生外部碎片、文件扩展困难

// ContiguousBlock 连续分配记录
type ContiguousBlock struct {
	FileName string
	Start    int // 起始块号
	Length   int // 占用块数
}

// ContiguousAllocation 连续分配模拟
type ContiguousAllocation struct {
	TotalBlocks int
	Disk        []string // 每个块的占用者（""表示空闲）
	Files       []ContiguousBlock
}

// NewContiguousAllocation 创建连续分配系统
func NewContiguousAllocation(totalBlocks int) *ContiguousAllocation {
	return &ContiguousAllocation{
		TotalBlocks: totalBlocks,
		Disk:        make([]string, totalBlocks),
	}
}

// Allocate 分配连续空间（首次适应）
func (ca *ContiguousAllocation) Allocate(fileName string, blocks int) bool {
	// 首次适应：找到第一个足够大的连续空闲区
	freeStart := -1
	freeCount := 0

	for i := 0; i < ca.TotalBlocks; i++ {
		if ca.Disk[i] == "" {
			if freeStart == -1 {
				freeStart = i
			}
			freeCount++
			if freeCount >= blocks {
				// 分配
				for j := freeStart; j < freeStart+blocks; j++ {
					ca.Disk[j] = fileName
				}
				ca.Files = append(ca.Files, ContiguousBlock{fileName, freeStart, blocks})
				return true
			}
		} else {
			freeStart = -1
			freeCount = 0
		}
	}
	return false // 没有足够的连续空间
}

// Free 释放文件
func (ca *ContiguousAllocation) Free(fileName string) {
	for i := range ca.Disk {
		if ca.Disk[i] == fileName {
			ca.Disk[i] = ""
		}
	}
	newFiles := make([]ContiguousBlock, 0)
	for _, f := range ca.Files {
		if f.FileName != fileName {
			newFiles = append(newFiles, f)
		}
	}
	ca.Files = newFiles
}

// --- 2. 链接分配 ---
// 每个块包含指向下一块的指针
// 优点：无外部碎片、文件扩展容易
// 缺点：只能顺序访问、指针占空间

// LinkedBlock 链接分配中的块
type LinkedBlock struct {
	BlockNum int
	Next     int // 下一块号，-1表示结尾
}

// LinkedAllocation 链接分配（显式链接/FAT表）
type LinkedAllocation struct {
	TotalBlocks int
	FAT         []int          // FAT表：FAT[i] = i块的下一块号，-1=结尾，-2=空闲
	FileStart   map[string]int // 文件名 → 起始块号
}

// NewLinkedAllocation 创建链接分配系统
func NewLinkedAllocation(totalBlocks int) *LinkedAllocation {
	fat := make([]int, totalBlocks)
	for i := range fat {
		fat[i] = -2 // 全部空闲
	}
	return &LinkedAllocation{
		TotalBlocks: totalBlocks,
		FAT:         fat,
		FileStart:   make(map[string]int),
	}
}

// Allocate 分配（链接方式）
func (la *LinkedAllocation) Allocate(fileName string, blocks int) bool {
	// 找空闲块
	freeBlocks := make([]int, 0)
	for i := 0; i < la.TotalBlocks && len(freeBlocks) < blocks; i++ {
		if la.FAT[i] == -2 {
			freeBlocks = append(freeBlocks, i)
		}
	}
	if len(freeBlocks) < blocks {
		return false
	}

	// 链接空闲块
	la.FileStart[fileName] = freeBlocks[0]
	for i := 0; i < len(freeBlocks)-1; i++ {
		la.FAT[freeBlocks[i]] = freeBlocks[i+1]
	}
	la.FAT[freeBlocks[len(freeBlocks)-1]] = -1 // 结尾标记

	return true
}

// GetFileBlocks 获取文件占用的所有块
func (la *LinkedAllocation) GetFileBlocks(fileName string) []int {
	start, ok := la.FileStart[fileName]
	if !ok {
		return nil
	}

	blocks := make([]int, 0)
	current := start
	for current != -1 {
		blocks = append(blocks, current)
		current = la.FAT[current]
	}
	return blocks
}

// --- 3. 索引分配 ---
// 每个文件有一个索引块，存储所有数据块号
// 优点：支持随机访问、无外部碎片
// 缺点：索引块开销

// IndexAllocation 索引分配
type IndexAllocation struct {
	TotalBlocks int
	Disk        []int          // 块使用状态：0=空闲，1=数据块，2=索引块
	FileIndex   map[string]int // 文件名 → 索引块号
	IndexBlocks map[int][]int  // 索引块号 → 数据块列表
}

// NewIndexAllocation 创建索引分配系统
func NewIndexAllocation(totalBlocks int) *IndexAllocation {
	return &IndexAllocation{
		TotalBlocks: totalBlocks,
		Disk:        make([]int, totalBlocks),
		FileIndex:   make(map[string]int),
		IndexBlocks: make(map[int][]int),
	}
}

// Allocate 分配（索引方式）
func (ia *IndexAllocation) Allocate(fileName string, blocks int) bool {
	// 需要 blocks个数据块 + 1个索引块
	freeBlocks := make([]int, 0)
	for i := 0; i < ia.TotalBlocks && len(freeBlocks) < blocks+1; i++ {
		if ia.Disk[i] == 0 {
			freeBlocks = append(freeBlocks, i)
		}
	}
	if len(freeBlocks) < blocks+1 {
		return false
	}

	// 第一个空闲块作为索引块
	indexBlock := freeBlocks[0]
	dataBlocks := freeBlocks[1 : blocks+1]

	ia.Disk[indexBlock] = 2
	ia.FileIndex[fileName] = indexBlock
	ia.IndexBlocks[indexBlock] = dataBlocks

	for _, b := range dataBlocks {
		ia.Disk[b] = 1
	}
	return true
}

// FileAllocationExample 文件分配方式示例
func FileAllocationExample() {
	fmt.Println("\n--- 文件分配方式 ---")

	// 1. 连续分配
	fmt.Println("\n【连续分配】")
	ca := NewContiguousAllocation(20)
	ca.Allocate("文件A", 3)
	ca.Allocate("文件B", 5)
	ca.Allocate("文件C", 2)
	fmt.Println("磁盘状态:")
	for i, owner := range ca.Disk {
		if owner != "" {
			fmt.Printf("  块%2d: %s\n", i, owner)
		}
	}
	fmt.Println("文件目录:")
	for _, f := range ca.Files {
		fmt.Printf("  %s: 起始块=%d, 长度=%d\n", f.FileName, f.Start, f.Length)
	}

	// 删除文件B，制造碎片
	ca.Free("文件B")
	fmt.Println("\n删除文件B后，尝试分配6块的文件D:")
	ok := ca.Allocate("文件D", 6)
	fmt.Printf("  分配结果: %v (连续空间不足导致失败)\n", ok)

	// 2. 链接分配（FAT）
	fmt.Println("\n【链接分配 (FAT表)】")
	la := NewLinkedAllocation(20)
	la.Allocate("文件A", 3)
	la.Allocate("文件B", 4)

	fmt.Println("FAT表 (非空闲部分):")
	for i, next := range la.FAT {
		if next != -2 {
			nextStr := fmt.Sprintf("%d", next)
			if next == -1 {
				nextStr = "EOF"
			}
			fmt.Printf("  FAT[%d] = %s\n", i, nextStr)
		}
	}
	fmt.Println("文件块链:")
	for name := range la.FileStart {
		blocks := la.GetFileBlocks(name)
		fmt.Printf("  %s: %v\n", name, blocks)
	}

	// 3. 索引分配
	fmt.Println("\n【索引分配】")
	ia := NewIndexAllocation(20)
	ia.Allocate("文件A", 3)
	ia.Allocate("文件B", 4)

	fmt.Println("文件索引信息:")
	for name, idxBlock := range ia.FileIndex {
		dataBlocks := ia.IndexBlocks[idxBlock]
		fmt.Printf("  %s: 索引块=%d, 数据块=%v\n", name, idxBlock, dataBlocks)
	}

	fmt.Println("\n408考点总结:")
	fmt.Println("  ┌────────────┬───────────┬───────────┬───────────┐")
	fmt.Println("  │ 分配方式   │ 顺序访问  │ 随机访问  │ 外部碎片  │")
	fmt.Println("  ├────────────┼───────────┼───────────┼───────────┤")
	fmt.Println("  │ 连续分配   │ 快        │ 快        │ 有        │")
	fmt.Println("  │ 链接分配   │ 较快      │ 慢        │ 无        │")
	fmt.Println("  │ 索引分配   │ 较快      │ 快        │ 无        │")
	fmt.Println("  └────────────┴───────────┴───────────┴───────────┘")
}
