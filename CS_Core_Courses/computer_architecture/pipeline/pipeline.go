package pipeline

import (
	"fmt"
	"strings"
)

// PipelineStage 流水线阶段
type PipelineStage int

const (
	IF  PipelineStage = iota // 取指令 (Instruction Fetch)
	ID                       // 译码 (Instruction Decode)
	EX                       // 执行 (Execute)
	MEM                      // 访存 (Memory Access)
	WB                       // 写回 (Write Back)
)

func (ps PipelineStage) String() string {
	stages := []string{"IF", "ID", "EX", "MEM", "WB"}
	if ps >= 0 && int(ps) < len(stages) {
		return stages[ps]
	}
	return "??"
}

// InstructionType 指令类型
type InstructionType int

const (
	TypeALU    InstructionType = iota // 算术逻辑运算
	TypeLoad                          // 加载指令
	TypeStore                         // 存储指令
	TypeBranch                        // 分支指令
)

// Instruction 指令定义
// 408 考点：指令在流水线中的需求和相关性
type Instruction struct {
	ID       int             // 指令编号
	Name     string          // 指令名称
	Type     InstructionType // 指令类型
	SrcRegs  []string        // 源寄存器
	DestReg  string          // 目的寄存器
	UseMem   bool            // 是否访存
	IsBranch bool            // 是否是分支指令
}

// PipelineSimulator 流水线模拟器
type PipelineSimulator struct {
	Instructions     []Instruction  // 指令序列
	Timeline         [][]string     // 时间线（每个时钟周期每条指令的状态）
	RegisterState    map[string]int // 寄存器状态（记录哪条指令最后写入）
	Stalls           int            // 暂停周期数
	Cycles           int            // 总周期数
	EnableForwarding bool           // 是否启用转发
}

// NewPipelineSimulator 创建流水线模拟器
func NewPipelineSimulator(instructions []Instruction, enableForwarding bool) *PipelineSimulator {
	return &PipelineSimulator{
		Instructions:     instructions,
		Timeline:         make([][]string, 0),
		RegisterState:    make(map[string]int),
		Stalls:           0,
		Cycles:           0,
		EnableForwarding: enableForwarding,
	}
}

// Run 运行流水线模拟
// 408 考点：模拟流水线执行过程，检测冲突
func (ps *PipelineSimulator) Run() {
	numInstructions := len(ps.Instructions)
	if numInstructions == 0 {
		return
	}

	// 初始化时间线
	// 指令状态：IF, ID, EX, MEM, WB, stall(暂停), done(完成), -(未开始)
	maxCycles := numInstructions*5 + 20 // 预估最大周期数
	ps.Timeline = make([][]string, numInstructions)
	for i := range ps.Timeline {
		ps.Timeline[i] = make([]string, maxCycles)
		for j := range ps.Timeline[i] {
			ps.Timeline[i][j] = "-"
		}
	}

	// 跟踪每条指令当前所在阶段
	instrStage := make([]PipelineStage, numInstructions)
	for i := range instrStage {
		instrStage[i] = -1 // -1 表示未开始
	}

	cycle := 0
	completed := 0

	// 模拟流水线执行
	for completed < numInstructions {
		// 反向遍历（从 WB 到 IF），避免同一周期内的冲突
		for stage := WB; stage >= IF; stage-- {
			for i := 0; i < numInstructions; i++ {
				if instrStage[i] == stage {
					canProgress := true

					// 检查是否可以进入下一阶段
					if stage == ID {
						// 在 ID 阶段检查数据冲突
						canProgress = ps.checkDataHazard(i, cycle, instrStage)
					} else if stage == IF {
						// 检查下一阶段（ID）是否有指令占用
						// 简化处理：如果前一条指令被暂停，当前指令也暂停
						if i > 0 && instrStage[i-1] == IF {
							canProgress = false
						}
					}

					if canProgress {
						// 进入下一阶段
						if stage == WB {
							// 完成执行
							ps.Timeline[i][cycle] = "WB"
							instrStage[i] = -1
							completed++

							// 更新寄存器状态
							if ps.Instructions[i].DestReg != "" {
								ps.RegisterState[ps.Instructions[i].DestReg] = i
							}
						} else {
							ps.Timeline[i][cycle] = stage.String()
							instrStage[i] = stage + 1
						}
					} else {
						// 暂停（插入气泡）
						ps.Timeline[i][cycle] = "stall"
						ps.Stalls++
					}
				}
			}
		}

		// 启动新指令
		for i := 0; i < numInstructions; i++ {
			if instrStage[i] == -1 {
				// 检查是否可以开始（IF 阶段是否空闲）
				canStart := true
				if i > 0 {
					// 确保上一条指令已经离开 IF 阶段
					prevStage := instrStage[i-1]
					if prevStage == 0 || prevStage == -1 {
						canStart = cycle > 0 // 第一个周期除外
					}
				}

				if canStart {
					ps.Timeline[i][cycle] = "IF"
					instrStage[i] = IF + 1
					break // 每个周期只启动一条指令
				}
			}
		}

		cycle++
		if cycle >= maxCycles {
			break
		}
	}

	ps.Cycles = cycle
}

// checkDataHazard 检查数据冲突
// 408 考点：RAW (Read After Write) 数据相关
func (ps *PipelineSimulator) checkDataHazard(instrIndex int, cycle int, instrStage []PipelineStage) bool {
	currentInstr := ps.Instructions[instrIndex]

	// 检查当前指令的源寄存器
	for _, srcReg := range currentInstr.SrcRegs {
		if srcReg == "" {
			continue
		}

		// 查找前面的指令是否会写入该寄存器
		for i := 0; i < instrIndex; i++ {
			prevInstr := ps.Instructions[i]

			// 如果前面的指令写入当前指令需要读取的寄存器
			if prevInstr.DestReg == srcReg {
				prevStage := instrStage[i]

				if ps.EnableForwarding {
					// 启用转发：如果前面指令在 EX 或之后阶段，可以转发
					if prevStage >= 0 && prevStage <= EX {
						// 前面指令还在 IF/ID/EX，需要等待
						return false
					}
				} else {
					// 不启用转发：必须等待前面指令完全完成
					if prevStage >= 0 {
						return false
					}
				}
			}
		}
	}

	return true
}

// PrintTimeline 打印流水线时空图
// 408 考点：流水线时空图是考试常见题型
func (ps *PipelineSimulator) PrintTimeline() {
	fmt.Println("\n流水线时空图：")
	fmt.Println("（stall 表示暂停，- 表示未开始/已完成）")

	// 打印表头
	fmt.Print("指令\\时钟 |")
	for c := 0; c < ps.Cycles; c++ {
		fmt.Printf(" %2d |", c+1)
	}
	fmt.Println()

	// 打印分隔线
	fmt.Print("----------|")
	for c := 0; c < ps.Cycles; c++ {
		fmt.Print("----|")
	}
	fmt.Println()

	// 打印每条指令的时间线
	for i, instr := range ps.Instructions {
		fmt.Printf("%-9s |", instr.Name)
		for c := 0; c < ps.Cycles; c++ {
			stage := ps.Timeline[i][c]
			if stage == "-" {
				fmt.Print("  - |")
			} else if stage == "stall" {
				fmt.Print("  X |") // X 表示暂停
			} else {
				fmt.Printf(" %2s |", stage)
			}
		}
		fmt.Println()
	}
}

// GetStatistics 获取流水线统计信息
// 408 考点：计算加速比、吞吐率、效率
func (ps *PipelineSimulator) GetStatistics() string {
	numInstructions := len(ps.Instructions)
	if numInstructions == 0 || ps.Cycles == 0 {
		return "无统计信息"
	}

	// 非流水线执行时间（假设每条指令 5 个时钟周期）
	nonPipelineTime := numInstructions * 5

	// 加速比
	speedup := float64(nonPipelineTime) / float64(ps.Cycles)

	// 吞吐率（指令数/周期数）
	throughput := float64(numInstructions) / float64(ps.Cycles)

	// 理论最短时间（流水线段数 + 指令数 - 1）
	idealCycles := 5 + numInstructions - 1

	// 效率（实际使用的时空区 / 总时空区）
	// 时空区 = 有效工作的格子数
	usedSlots := 0
	totalSlots := ps.Cycles * 5 // 5 个流水线段
	for i := range ps.Instructions {
		for c := 0; c < ps.Cycles; c++ {
			if ps.Timeline[i][c] != "-" && ps.Timeline[i][c] != "stall" {
				usedSlots++
			}
		}
	}
	efficiency := float64(usedSlots) / float64(totalSlots) * 100

	result := fmt.Sprintf(`
流水线统计信息：
  指令数量:       %d
  执行周期数:     %d
  暂停周期数:     %d
  
性能指标：
  非流水线时间:   %d 周期
  流水线时间:     %d 周期
  理想周期数:     %d 周期
  加速比:         %.2f
  吞吐率:         %.3f 条指令/周期
  效率:           %.2f%%
  
转发机制:         %s
`,
		numInstructions, ps.Cycles, ps.Stalls,
		nonPipelineTime, ps.Cycles, idealCycles,
		speedup, throughput, efficiency,
		getForwardingStatus(ps.EnableForwarding))

	return result
}

func getForwardingStatus(enabled bool) string {
	if enabled {
		return "启用"
	}
	return "禁用"
}

// PipelineExample 流水线示例程序
// 408 考点：演示数据冲突的检测和解决
func PipelineExample() {
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  5 段流水线模拟")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 典型的 408 考试指令序列（存在 RAW 数据相关）
	instructions := []Instruction{
		{ID: 1, Name: "I1", Type: TypeALU, SrcRegs: []string{"R1", "R2"}, DestReg: "R3"},
		{ID: 2, Name: "I2", Type: TypeALU, SrcRegs: []string{"R3", "R4"}, DestReg: "R5"},  // 依赖 I1 的 R3
		{ID: 3, Name: "I3", Type: TypeLoad, SrcRegs: []string{"R5"}, DestReg: "R6"},       // 依赖 I2 的 R5
		{ID: 4, Name: "I4", Type: TypeALU, SrcRegs: []string{"R6", "R7"}, DestReg: "R8"},  // 依赖 I3 的 R6
		{ID: 5, Name: "I5", Type: TypeALU, SrcRegs: []string{"R8", "R9"}, DestReg: "R10"}, // 依赖 I4 的 R8
	}

	// 打印指令序列
	fmt.Println("\n指令序列：")
	for i, instr := range instructions {
		srcRegs := strings.Join(instr.SrcRegs, ", ")
		fmt.Printf("  I%d: %s ← op(%s)\n", i+1, instr.DestReg, srcRegs)
	}

	fmt.Println("\n数据相关分析：")
	fmt.Println("  I2 依赖 I1 (RAW on R3)")
	fmt.Println("  I3 依赖 I2 (RAW on R5)")
	fmt.Println("  I4 依赖 I3 (RAW on R6)")
	fmt.Println("  I5 依赖 I4 (RAW on R8)")

	// 1. 不启用转发机制
	fmt.Println("\n【场景 1：不启用转发机制】")
	fmt.Println("说明：每条指令必须等待前面指令完全完成才能读取寄存器")
	simulator1 := NewPipelineSimulator(instructions, false)
	simulator1.Run()
	simulator1.PrintTimeline()
	fmt.Println(simulator1.GetStatistics())

	// 2. 启用转发机制
	fmt.Println("\n【场景 2：启用转发机制】")
	fmt.Println("说明：EX/MEM/WB 阶段的结果可以直接转发到 ID 阶段")
	simulator2 := NewPipelineSimulator(instructions, true)
	simulator2.Run()
	simulator2.PrintTimeline()
	fmt.Println(simulator2.GetStatistics())

	// 3. 无数据相关的理想情况
	fmt.Println("\n【场景 3：无数据相关的理想流水线】")
	fmt.Println("说明：指令之间没有数据相关，流水线充分利用")
	independentInstructions := []Instruction{
		{ID: 1, Name: "I1", Type: TypeALU, SrcRegs: []string{"R1", "R2"}, DestReg: "R3"},
		{ID: 2, Name: "I2", Type: TypeALU, SrcRegs: []string{"R4", "R5"}, DestReg: "R6"},
		{ID: 3, Name: "I3", Type: TypeALU, SrcRegs: []string{"R7", "R8"}, DestReg: "R9"},
		{ID: 4, Name: "I4", Type: TypeALU, SrcRegs: []string{"R10", "R11"}, DestReg: "R12"},
		{ID: 5, Name: "I5", Type: TypeALU, SrcRegs: []string{"R13", "R14"}, DestReg: "R15"},
	}
	simulator3 := NewPipelineSimulator(independentInstructions, true)
	simulator3.Run()
	simulator3.PrintTimeline()
	fmt.Println(simulator3.GetStatistics())

	fmt.Println("\n" + pipeline408Summary())
}

// pipeline408Summary 408 考试总结
func pipeline408Summary() string {
	return `
╔════════════════════════════════════════════════════════════════╗
║                  408 考试要点总结 - 流水线                     ║
╠════════════════════════════════════════════════════════════════╣
║ 1. 流水线基本原理：                                            ║
║    • 将指令执行分为多个阶段，各阶段并行工作                   ║
║    • 5 段经典流水线：IF → ID → EX → MEM → WB                  ║
║    • 每个时钟周期可以同时处理多条指令的不同阶段               ║
║                                                                ║
║ 2. 流水线冲突类型：                                            ║
║    【结构冲突】硬件资源不足（如只有一个内存端口）             ║
║      解决：增加硬件资源、暂停流水线                           ║
║                                                                ║
║    【数据冲突】指令间存在数据相关                             ║
║      • RAW (Read After Write): 真相关，最常见                 ║
║        例：I1: R3←R1+R2; I2: R5←R3+R4 (I2 读 R3 依赖 I1 写)  ║
║      • WAW (Write After Write): 输出相关，乱序执行时可能出现  ║
║      • WAR (Write After Read): 反相关，乱序执行时可能出现     ║
║      解决：数据转发(forwarding)、暂停(stall)、编译器调度      ║
║                                                                ║
║    【控制冲突】分支指令改变 PC 值                             ║
║      解决：分支预测、延迟分支、提前计算分支目标               ║
║                                                                ║
║ 3. 数据转发机制：                                              ║
║    • EX/MEM → ID: EX 阶段结果直接转发到 ID 阶段              ║
║    • MEM/WB → ID: MEM 阶段结果转发到 ID 阶段                 ║
║    • 减少暂停周期，提高流水线效率                             ║
║    • Load 指令后紧跟使用其结果的指令，仍需插入 1 个气泡       ║
║                                                                ║
║ 4. 性能指标计算：                                              ║
║    【执行时间】                                                ║
║      理想：T = (k + n - 1) × Δt                               ║
║      实际：T = (k + n - 1 + 暂停周期数) × Δt                  ║
║      其中 k=段数, n=指令数, Δt=时钟周期                       ║
║                                                                ║
║    【加速比】                                                  ║
║      S = 不使用流水线时间 / 使用流水线时间                    ║
║      理想加速比：S = k (流水线段数)                           ║
║                                                                ║
║    【吞吐率】                                                  ║
║      TP = n / T (完成的指令数 / 总时间)                       ║
║      最大吞吐率：TPmax = 1 / Δt                               ║
║                                                                ║
║    【效率】                                                    ║
║      E = (有效工作的时空区) / (总时空区) × 100%               ║
║                                                                ║
║ 5. 典型考题：                                                  ║
║    • 画出给定指令序列的流水线时空图                           ║
║    • 识别数据相关，判断是否需要暂停                           ║
║    • 计算执行时间、加速比、吞吐率、效率                       ║
║    • 比较启用/禁用转发机制的性能差异                          ║
║    • 分析分支指令对流水线的影响                               ║
║                                                                ║
║ 6. 提高流水线性能的方法：                                      ║
║    • 数据转发：减少数据冲突导致的暂停                         ║
║    • 指令重排：编译器调整指令顺序，减少相关性                 ║
║    • 分支预测：减少控制冲突的影响                             ║
║    • 增加流水线段数：提高时钟频率（超流水线）                 ║
║    • 多发射：每周期发射多条指令（超标量）                     ║
║                                                                ║
║ 7. 记忆要点：                                                  ║
║    • 流水线不改变单条指令的执行时间，而是提高吞吐率          ║
║    • 气泡（bubble）是指流水线中的空闲槽位                    ║
║    • 暂停（stall）会降低流水线效率                            ║
║    • 理想情况下，n 条指令在 k 段流水线中执行 (k+n-1) 周期    ║
╚════════════════════════════════════════════════════════════════╝
`
}
