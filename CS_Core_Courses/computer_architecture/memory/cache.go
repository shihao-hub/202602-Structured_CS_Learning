package memory

import (
	"fmt"
	"math"
	"math/rand"
)

// CacheLine 表示 Cache 中的一行
// 408 考点：Cache 行结构（有效位、标记、数据块）
type CacheLine struct {
	Valid      bool   // 有效位：标识该行是否存有有效数据
	Tag        int    // 标记：用于区分主存中映射到同一 Cache 行的不同块
	Data       []byte // 数据块：实际存储的数据
	AccessTime int    // 访问时间戳：用于 LRU 算法
	LoadTime   int    // 装入时间：用于 FIFO 算法
}

// MappingType Cache 映射方式
type MappingType int

const (
	DirectMapped     MappingType = iota // 直接映射
	FullyAssociative                    // 全相联映射
	SetAssociative                      // 组相联映射
)

// ReplacementPolicy Cache 替换策略
type ReplacementPolicy int

const (
	LRU    ReplacementPolicy = iota // 最近最少使用
	FIFO                            // 先进先出
	Random                          // 随机替换
)

// CacheConfig Cache 配置参数
type CacheConfig struct {
	CacheSize     int               // Cache 总大小（字节）
	BlockSize     int               // 块大小（字节）
	Associativity int               // 相联度（组相联时使用，直接映射为1，全相联为行数）
	MappingType   MappingType       // 映射方式
	Policy        ReplacementPolicy // 替换策略
}

// CacheSimulator Cache 模拟器
type CacheSimulator struct {
	Config      CacheConfig
	Lines       []CacheLine // Cache 行数组
	NumLines    int         // Cache 行数
	NumSets     int         // 组数（直接映射和全相联时为特殊值）
	BlockOffset int         // 块内地址位数
	IndexBits   int         // 索引位数
	TagBits     int         // 标记位数
	Hits        int         // 命中次数
	Misses      int         // 缺失次数
	AccessCount int         // 总访问次数
	CurrentTime int         // 当前时间（用于 LRU 和 FIFO）
}

// NewCacheSimulator 创建 Cache 模拟器
// 408 考点：根据 Cache 配置计算各字段位数
func NewCacheSimulator(config CacheConfig) *CacheSimulator {
	numLines := config.CacheSize / config.BlockSize

	var numSets int
	var indexBits int

	switch config.MappingType {
	case DirectMapped:
		// 直接映射：组数 = 行数，相联度 = 1
		numSets = numLines
		config.Associativity = 1
		indexBits = int(math.Log2(float64(numSets)))
	case FullyAssociative:
		// 全相联：只有 1 组，相联度 = 行数
		numSets = 1
		config.Associativity = numLines
		indexBits = 0
	case SetAssociative:
		// 组相联：组数 = 行数 / 相联度
		numSets = numLines / config.Associativity
		indexBits = int(math.Log2(float64(numSets)))
	}

	blockOffset := int(math.Log2(float64(config.BlockSize)))
	tagBits := 32 - indexBits - blockOffset // 假设 32 位地址

	lines := make([]CacheLine, numLines)
	for i := range lines {
		lines[i] = CacheLine{
			Valid:      false,
			Tag:        -1,
			Data:       make([]byte, config.BlockSize),
			AccessTime: 0,
			LoadTime:   0,
		}
	}

	return &CacheSimulator{
		Config:      config,
		Lines:       lines,
		NumLines:    numLines,
		NumSets:     numSets,
		BlockOffset: blockOffset,
		IndexBits:   indexBits,
		TagBits:     tagBits,
		Hits:        0,
		Misses:      0,
		AccessCount: 0,
		CurrentTime: 0,
	}
}

// ParseAddress 解析内存地址
// 408 考点：地址分解（标记、索引、块内偏移）
func (cs *CacheSimulator) ParseAddress(address int) (tag, index, offset int) {
	offset = address & ((1 << cs.BlockOffset) - 1)
	address >>= cs.BlockOffset

	if cs.IndexBits > 0 {
		index = address & ((1 << cs.IndexBits) - 1)
		address >>= cs.IndexBits
	} else {
		index = 0
	}

	tag = address
	return
}

// Access 访问 Cache
// 408 考点：Cache 访问过程（查找、命中/缺失判断、替换）
func (cs *CacheSimulator) Access(address int) (hit bool, message string) {
	cs.AccessCount++
	cs.CurrentTime++

	tag, index, offset := cs.ParseAddress(address)

	// 确定搜索范围
	startLine := index * cs.Config.Associativity
	endLine := startLine + cs.Config.Associativity

	// 查找 Cache 是否命中
	for i := startLine; i < endLine; i++ {
		if cs.Lines[i].Valid && cs.Lines[i].Tag == tag {
			// Cache 命中
			cs.Hits++
			cs.Lines[i].AccessTime = cs.CurrentTime
			message = fmt.Sprintf("命中 - 地址: 0x%X, Tag: %d, Index: %d, Offset: %d",
				address, tag, index, offset)
			return true, message
		}
	}

	// Cache 缺失，需要替换
	cs.Misses++
	victimLine := cs.selectVictim(startLine, endLine)

	cs.Lines[victimLine].Valid = true
	cs.Lines[victimLine].Tag = tag
	cs.Lines[victimLine].AccessTime = cs.CurrentTime
	cs.Lines[victimLine].LoadTime = cs.CurrentTime

	message = fmt.Sprintf("缺失 - 地址: 0x%X, Tag: %d, Index: %d, Offset: %d, 替换行: %d",
		address, tag, index, offset, victimLine)
	return false, message
}

// selectVictim 选择被替换的 Cache 行
// 408 考点：三种替换算法的实现
func (cs *CacheSimulator) selectVictim(start, end int) int {
	// 首先查找是否有无效行
	for i := start; i < end; i++ {
		if !cs.Lines[i].Valid {
			return i
		}
	}

	// 根据替换策略选择牺牲行
	switch cs.Config.Policy {
	case LRU:
		// LRU：选择最久未使用的行
		minTime := cs.Lines[start].AccessTime
		victim := start
		for i := start + 1; i < end; i++ {
			if cs.Lines[i].AccessTime < minTime {
				minTime = cs.Lines[i].AccessTime
				victim = i
			}
		}
		return victim

	case FIFO:
		// FIFO：选择最早装入的行
		minLoad := cs.Lines[start].LoadTime
		victim := start
		for i := start + 1; i < end; i++ {
			if cs.Lines[i].LoadTime < minLoad {
				minLoad = cs.Lines[i].LoadTime
				victim = i
			}
		}
		return victim

	case Random:
		// 随机替换
		offset := rand.Intn(end - start)
		return start + offset

	default:
		return start
	}
}

// GetStatistics 获取 Cache 统计信息
// 408 考点：命中率计算
func (cs *CacheSimulator) GetStatistics() string {
	hitRate := 0.0
	if cs.AccessCount > 0 {
		hitRate = float64(cs.Hits) / float64(cs.AccessCount) * 100
	}

	return fmt.Sprintf(`
Cache 统计信息：
  总访问次数: %d
  命中次数:   %d
  缺失次数:   %d
  命中率:     %.2f%%
  
Cache 配置：
  映射方式:   %s
  替换策略:   %s
  Cache大小:  %d 字节
  块大小:     %d 字节
  行数:       %d
  组数:       %d
  相联度:     %d
  
地址分解：
  标记位数:   %d 位
  索引位数:   %d 位
  块内偏移:   %d 位
`,
		cs.AccessCount, cs.Hits, cs.Misses, hitRate,
		getMappingTypeName(cs.Config.MappingType),
		getPolicyName(cs.Config.Policy),
		cs.Config.CacheSize, cs.Config.BlockSize,
		cs.NumLines, cs.NumSets, cs.Config.Associativity,
		cs.TagBits, cs.IndexBits, cs.BlockOffset)
}

func getMappingTypeName(mt MappingType) string {
	switch mt {
	case DirectMapped:
		return "直接映射"
	case FullyAssociative:
		return "全相联映射"
	case SetAssociative:
		return "组相联映射"
	default:
		return "未知"
	}
}

func getPolicyName(policy ReplacementPolicy) string {
	switch policy {
	case LRU:
		return "LRU (最近最少使用)"
	case FIFO:
		return "FIFO (先进先出)"
	case Random:
		return "随机替换"
	default:
		return "未知"
	}
}

// CacheExample Cache 示例程序
// 408 考点：用相同的访问序列对比三种映射方式
func CacheExample() {
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  Cache 高速缓存模拟")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 访问序列：典型的 408 考试题型
	// 主存块号: 0, 1, 2, 3, 0, 1, 4, 0
	addresses := []int{
		0x0000, // 块 0
		0x0040, // 块 1
		0x0080, // 块 2
		0x00C0, // 块 3
		0x0000, // 块 0
		0x0040, // 块 1
		0x0100, // 块 4
		0x0000, // 块 0
	}

	// 1. 直接映射
	fmt.Println("\n【1. 直接映射 (Direct Mapped) - LRU】")
	fmt.Println("说明：每个主存块只能映射到唯一的 Cache 行")
	runCacheTest(CacheConfig{
		CacheSize:   256,
		BlockSize:   64,
		MappingType: DirectMapped,
		Policy:      LRU,
	}, addresses)

	// 2. 全相联映射
	fmt.Println("\n【2. 全相联映射 (Fully Associative) - LRU】")
	fmt.Println("说明：主存块可以映射到 Cache 的任意行")
	runCacheTest(CacheConfig{
		CacheSize:   256,
		BlockSize:   64,
		MappingType: FullyAssociative,
		Policy:      LRU,
	}, addresses)

	// 3. 组相联映射 (2 路)
	fmt.Println("\n【3. 二路组相联映射 (2-Way Set Associative) - LRU】")
	fmt.Println("说明：主存块映射到特定组，可放在组内任意一行")
	runCacheTest(CacheConfig{
		CacheSize:     256,
		BlockSize:     64,
		Associativity: 2,
		MappingType:   SetAssociative,
		Policy:        LRU,
	}, addresses)

	// 4. 演示不同替换算法的差异
	fmt.Println("\n【4. 替换算法对比 (直接映射)】")
	policies := []ReplacementPolicy{LRU, FIFO, Random}
	for _, policy := range policies {
		fmt.Printf("\n使用 %s 策略：\n", getPolicyName(policy))
		runCacheTest(CacheConfig{
			CacheSize:   256,
			BlockSize:   64,
			MappingType: DirectMapped,
			Policy:      policy,
		}, addresses)
	}

	fmt.Println("\n" + cache408Summary())
}

func runCacheTest(config CacheConfig, addresses []int) {
	simulator := NewCacheSimulator(config)

	for i, addr := range addresses {
		hit, msg := simulator.Access(addr)
		status := "✗ 缺失"
		if hit {
			status = "✓ 命中"
		}
		fmt.Printf("  访问 %d: %s [%s]\n", i+1, msg, status)
	}

	fmt.Println(simulator.GetStatistics())
}

// cache408Summary 408 考试总结
func cache408Summary() string {
	return `
╔════════════════════════════════════════════════════════════════╗
║                    408 考试要点总结 - Cache                    ║
╠════════════════════════════════════════════════════════════════╣
║ 1. 映射方式特点：                                              ║
║    • 直接映射：硬件简单，冲突率高，命中率较低                 ║
║    • 全相联：命中率最高，硬件复杂，成本高                     ║
║    • 组相联：折中方案，最常用（如 2 路、4 路）                ║
║                                                                ║
║ 2. 地址分解公式：                                              ║
║    地址 = 标记(Tag) | 索引(Index) | 块内偏移(Offset)          ║
║    • 块内偏移位数 = log₂(块大小)                              ║
║    • 索引位数 = log₂(组数)                                    ║
║    • 标记位数 = 地址总位数 - 索引位数 - 块内偏移位数          ║
║                                                                ║
║ 3. Cache 容量计算：                                            ║
║    Cache 总容量 = 行数 × (有效位 + 标记位 + 数据块位)         ║
║                                                                ║
║ 4. 命中率计算：                                                ║
║    命中率 = 命中次数 / 总访问次数 × 100%                      ║
║                                                                ║
║ 5. 平均访问时间：                                              ║
║    Ta = H × Tc + (1-H) × (Tc + Tm)                            ║
║    其中 H=命中率, Tc=Cache访问时间, Tm=主存访问时间           ║
║                                                                ║
║ 6. 替换算法：                                                  ║
║    • LRU：最常考，需要记录访问时间，硬件复杂                  ║
║    • FIFO：记录装入时间，实现简单                             ║
║    • 随机：最简单，性能不可预测                               ║
║                                                                ║
║ 7. 典型考题类型：                                              ║
║    • 给定地址序列，画出 Cache 状态变化，计算命中率            ║
║    • 计算 Cache 各部分位数和总容量                            ║
║    • 分析不同映射方式和替换算法的性能差异                     ║
╚════════════════════════════════════════════════════════════════╝
`
}
