package filesystem

import "fmt"

// ============================================================
// 索引节点 (inode) 模拟
// 408考点：文件控制块、索引节点结构、直接/间接寻址
// ============================================================

// 常量定义
const (
	BlockSize        = 4096 // 磁盘块大小（字节）
	DirectBlocks     = 12   // 直接块指针数
	IndirectPerBlock = 1024 // 每个间接块包含的指针数（BlockSize/4）
)

// FileType 文件类型
type FileType int

const (
	RegularFile FileType = iota
	Directory
	SymLink
)

func (ft FileType) String() string {
	switch ft {
	case RegularFile:
		return "普通文件"
	case Directory:
		return "目录"
	case SymLink:
		return "符号链接"
	default:
		return "未知"
	}
}

// Permission 文件权限
type Permission struct {
	Read    bool
	Write   bool
	Execute bool
}

func (p Permission) String() string {
	s := ""
	if p.Read {
		s += "r"
	} else {
		s += "-"
	}
	if p.Write {
		s += "w"
	} else {
		s += "-"
	}
	if p.Execute {
		s += "x"
	} else {
		s += "-"
	}
	return s
}

// Inode 索引节点
type Inode struct {
	InodeNumber    int               // inode编号
	Type           FileType          // 文件类型
	Owner          Permission        // 属主权限
	Group          Permission        // 组权限
	Other          Permission        // 其他权限
	LinkCount      int               // 硬链接数
	FileSize       int64             // 文件大小（字节）
	DirectPtr      [DirectBlocks]int // 直接块指针（-1表示未分配）
	SingleIndirect int               // 一次间接块指针
	DoubleIndirect int               // 二次间接块指针
	TripleIndirect int               // 三次间接块指针
}

// NewInode 创建新inode
func NewInode(number int, fileType FileType) *Inode {
	inode := &Inode{
		InodeNumber:    number,
		Type:           fileType,
		Owner:          Permission{true, true, false},
		Group:          Permission{true, false, false},
		Other:          Permission{true, false, false},
		LinkCount:      1,
		SingleIndirect: -1,
		DoubleIndirect: -1,
		TripleIndirect: -1,
	}
	for i := range inode.DirectPtr {
		inode.DirectPtr[i] = -1
	}
	return inode
}

// MaxFileSize 计算inode支持的最大文件大小
// 直接块: 12 * BlockSize
// 一次间接: IndirectPerBlock * BlockSize
// 二次间接: IndirectPerBlock² * BlockSize
// 三次间接: IndirectPerBlock³ * BlockSize
func (inode *Inode) MaxFileSize() int64 {
	direct := int64(DirectBlocks) * int64(BlockSize)
	single := int64(IndirectPerBlock) * int64(BlockSize)
	double := int64(IndirectPerBlock) * int64(IndirectPerBlock) * int64(BlockSize)
	triple := int64(IndirectPerBlock) * int64(IndirectPerBlock) * int64(IndirectPerBlock) * int64(BlockSize)
	return direct + single + double + triple
}

// BlocksNeeded 计算存储指定大小文件所需的数据块数
func BlocksNeeded(fileSize int64) int {
	if fileSize <= 0 {
		return 0
	}
	blocks := int(fileSize / int64(BlockSize))
	if fileSize%int64(BlockSize) != 0 {
		blocks++
	}
	return blocks
}

// AllocateBlocks 模拟为inode分配数据块
func (inode *Inode) AllocateBlocks(fileSize int64, startBlock int) {
	inode.FileSize = fileSize
	needed := BlocksNeeded(fileSize)
	blockNum := startBlock

	// 分配直接块
	for i := 0; i < DirectBlocks && needed > 0; i++ {
		inode.DirectPtr[i] = blockNum
		blockNum++
		needed--
	}

	// 分配一次间接块
	if needed > 0 {
		inode.SingleIndirect = blockNum
		blockNum++ // 间接块本身占一个块
		allocated := needed
		if allocated > IndirectPerBlock {
			allocated = IndirectPerBlock
		}
		needed -= allocated
		blockNum += allocated
	}

	// 分配二次间接块
	if needed > 0 {
		inode.DoubleIndirect = blockNum
		blockNum++
		maxDouble := IndirectPerBlock * IndirectPerBlock
		allocated := needed
		if allocated > maxDouble {
			allocated = maxDouble
		}
		needed -= allocated
	}
}

// Print 打印inode信息
func (inode *Inode) Print() {
	fmt.Printf("inode #%d:\n", inode.InodeNumber)
	fmt.Printf("  类型: %s\n", inode.Type)
	fmt.Printf("  权限: %s%s%s\n", inode.Owner, inode.Group, inode.Other)
	fmt.Printf("  硬链接数: %d\n", inode.LinkCount)
	fmt.Printf("  文件大小: %d 字节\n", inode.FileSize)
	fmt.Printf("  所需数据块: %d\n", BlocksNeeded(inode.FileSize))

	fmt.Printf("  直接块指针: [")
	for i, ptr := range inode.DirectPtr {
		if ptr != -1 {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%d", ptr)
		}
	}
	fmt.Println("]")

	if inode.SingleIndirect != -1 {
		fmt.Printf("  一次间接块: %d\n", inode.SingleIndirect)
	}
	if inode.DoubleIndirect != -1 {
		fmt.Printf("  二次间接块: %d\n", inode.DoubleIndirect)
	}
}

// InodeExample inode示例
func InodeExample() {
	fmt.Println("\n--- 索引节点 (inode) ---")

	inode := NewInode(1, RegularFile)

	fmt.Printf("块大小: %d 字节\n", BlockSize)
	fmt.Printf("直接块数: %d\n", DirectBlocks)
	fmt.Printf("每间接块指针数: %d\n", IndirectPerBlock)
	fmt.Printf("最大文件大小: %.2f GB\n\n",
		float64(inode.MaxFileSize())/(1024*1024*1024))

	// 小文件（只用直接块）
	fmt.Println("【示例1: 小文件 10KB】")
	small := NewInode(1, RegularFile)
	small.AllocateBlocks(10*1024, 100)
	small.Print()

	// 中等文件（需要间接块）
	fmt.Println("\n【示例2: 中等文件 100KB】")
	medium := NewInode(2, RegularFile)
	medium.AllocateBlocks(100*1024, 200)
	medium.Print()

	// 计算不同大小文件的寻址方式
	fmt.Println("\n【文件大小与寻址方式】")
	directMax := int64(DirectBlocks) * int64(BlockSize)
	singleMax := directMax + int64(IndirectPerBlock)*int64(BlockSize)
	doubleMax := singleMax + int64(IndirectPerBlock)*int64(IndirectPerBlock)*int64(BlockSize)

	fmt.Printf("  仅直接寻址:     ≤ %d KB\n", directMax/1024)
	fmt.Printf("  一次间接寻址:   ≤ %.1f MB\n", float64(singleMax)/(1024*1024))
	fmt.Printf("  二次间接寻址:   ≤ %.1f GB\n", float64(doubleMax)/(1024*1024*1024))
}
