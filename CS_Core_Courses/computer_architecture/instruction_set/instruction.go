package instruction_set

import (
	"fmt"
)

// AddressingMode 寻址方式
// 408 考点：各种寻址方式的理解和应用
type AddressingMode int

const (
	Immediate        AddressingMode = iota // 立即寻址
	Direct                                 // 直接寻址
	Indirect                               // 间接寻址
	Register                               // 寄存器寻址
	RegisterIndirect                       // 寄存器间接寻址
	Relative                               // 相对寻址
	Base                                   // 基址寻址
	Indexed                                // 变址寻址
)

func (am AddressingMode) String() string {
	modes := []string{
		"立即寻址",
		"直接寻址",
		"间接寻址",
		"寄存器寻址",
		"寄存器间接寻址",
		"相对寻址",
		"基址寻址",
		"变址寻址",
	}
	if am >= 0 && int(am) < len(modes) {
		return modes[am]
	}
	return "未知"
}

// Instruction 指令定义
type Instruction struct {
	Opcode         string         // 操作码
	AddressingMode AddressingMode // 寻址方式
	Operand        int            // 操作数/地址/位移量
	RegisterName   string         // 寄存器名（用于寄存器寻址）
}

// MachineState 机器状态（用于模拟指令执行环境）
type MachineState struct {
	Memory    map[int]int    // 内存：地址 -> 值
	Registers map[string]int // 寄存器：名称 -> 值
	PC        int            // 程序计数器
	BR        int            // 基址寄存器
	IX        int            // 变址寄存器
}

// NewMachineState 创建机器状态
func NewMachineState() *MachineState {
	return &MachineState{
		Memory:    make(map[int]int),
		Registers: make(map[string]int),
		PC:        0x1000, // 假设 PC 初始值为 0x1000
		BR:        0x2000, // 基址寄存器初始值
		IX:        0,      // 变址寄存器初始值
	}
}

// GetEffectiveAddress 计算有效地址
// 408 考点：根据不同寻址方式计算有效地址
func (ms *MachineState) GetEffectiveAddress(instr Instruction) (ea int, needMemoryAccess bool, description string) {
	switch instr.AddressingMode {
	case Immediate:
		// 立即寻址：操作数就是指令中的立即数，无需计算地址
		description = fmt.Sprintf("立即数: %d", instr.Operand)
		return 0, false, description

	case Direct:
		// 直接寻址：EA = A
		ea = instr.Operand
		description = fmt.Sprintf("EA = %d (直接给出)", ea)
		return ea, true, description

	case Indirect:
		// 间接寻址：EA = (A)，需要两次访存
		addr := instr.Operand
		ea = ms.Memory[addr]
		description = fmt.Sprintf("EA = M[%d] = %d (一次间接)", addr, ea)
		return ea, true, description

	case Register:
		// 寄存器寻址：操作数在寄存器中，无需访存
		description = fmt.Sprintf("寄存器 %s 的值", instr.RegisterName)
		return 0, false, description

	case RegisterIndirect:
		// 寄存器间接寻址：EA = (Ri)
		ea = ms.Registers[instr.RegisterName]
		description = fmt.Sprintf("EA = %s = %d (寄存器间接)", instr.RegisterName, ea)
		return ea, true, description

	case Relative:
		// 相对寻址：EA = (PC) + D
		// 注意：这里假设 PC 指向当前指令，有些机器 PC 指向下一条指令
		ea = ms.PC + instr.Operand
		description = fmt.Sprintf("EA = PC + %d = %d + %d = %d",
			instr.Operand, ms.PC, instr.Operand, ea)
		return ea, true, description

	case Base:
		// 基址寻址：EA = (BR) + D
		ea = ms.BR + instr.Operand
		description = fmt.Sprintf("EA = BR + %d = %d + %d = %d",
			instr.Operand, ms.BR, instr.Operand, ea)
		return ea, true, description

	case Indexed:
		// 变址寻址：EA = (IX) + D
		ea = ms.IX + instr.Operand
		description = fmt.Sprintf("EA = IX + %d = %d + %d = %d",
			instr.Operand, ms.IX, instr.Operand, ea)
		return ea, true, description

	default:
		return 0, false, "未知寻址方式"
	}
}

// GetOperand 获取操作数
// 408 考点：完整的操作数获取过程
func (ms *MachineState) GetOperand(instr Instruction) (operand int, memoryAccessCount int, description string) {
	switch instr.AddressingMode {
	case Immediate:
		// 立即寻址：0 次访存
		operand = instr.Operand
		description = fmt.Sprintf("操作数 = %d (立即数，0次访存)", operand)
		return operand, 0, description

	case Direct:
		// 直接寻址：1 次访存
		ea := instr.Operand
		operand = ms.Memory[ea]
		description = fmt.Sprintf("操作数 = M[%d] = %d (1次访存)", ea, operand)
		return operand, 1, description

	case Indirect:
		// 间接寻址：2 次访存
		addr := instr.Operand
		ea := ms.Memory[addr]
		operand = ms.Memory[ea]
		description = fmt.Sprintf("操作数 = M[M[%d]] = M[%d] = %d (2次访存)",
			addr, ea, operand)
		return operand, 2, description

	case Register:
		// 寄存器寻址：0 次访存
		operand = ms.Registers[instr.RegisterName]
		description = fmt.Sprintf("操作数 = %s = %d (0次访存)",
			instr.RegisterName, operand)
		return operand, 0, description

	case RegisterIndirect:
		// 寄存器间接寻址：1 次访存
		ea := ms.Registers[instr.RegisterName]
		operand = ms.Memory[ea]
		description = fmt.Sprintf("操作数 = M[%s] = M[%d] = %d (1次访存)",
			instr.RegisterName, ea, operand)
		return operand, 1, description

	case Relative:
		// 相对寻址：1 次访存
		ea := ms.PC + instr.Operand
		operand = ms.Memory[ea]
		description = fmt.Sprintf("操作数 = M[PC+%d] = M[%d] = %d (1次访存)",
			instr.Operand, ea, operand)
		return operand, 1, description

	case Base:
		// 基址寻址：1 次访存
		ea := ms.BR + instr.Operand
		operand = ms.Memory[ea]
		description = fmt.Sprintf("操作数 = M[BR+%d] = M[%d] = %d (1次访存)",
			instr.Operand, ea, operand)
		return operand, 1, description

	case Indexed:
		// 变址寻址：1 次访存
		ea := ms.IX + instr.Operand
		operand = ms.Memory[ea]
		description = fmt.Sprintf("操作数 = M[IX+%d] = M[%d] = %d (1次访存)",
			instr.Operand, ea, operand)
		return operand, 1, description

	default:
		return 0, 0, "未知寻址方式"
	}
}

// InstructionSetExample 指令集示例程序
// 408 考点：演示各种寻址方式的计算过程
func InstructionSetExample() {
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  指令系统 - 寻址方式演示")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 初始化机器状态
	ms := NewMachineState()

	// 设置内存和寄存器的初始值
	ms.Memory[0x1000] = 100
	ms.Memory[0x1100] = 200
	ms.Memory[0x1200] = 300
	ms.Memory[0x2000] = 0x1100 // 间接寻址指针
	ms.Memory[0x2100] = 400
	ms.Memory[0x2500] = 500

	ms.Registers["R1"] = 50
	ms.Registers["R2"] = 0x1200 // 寄存器间接寻址指针
	ms.Registers["R3"] = 75

	ms.PC = 0x1000 // 程序计数器
	ms.BR = 0x2000 // 基址寄存器
	ms.IX = 0x100  // 变址寄存器

	// 打印初始状态
	fmt.Println("\n【机器初始状态】")
	fmt.Printf("  PC = 0x%04X\n", ms.PC)
	fmt.Printf("  BR = 0x%04X (基址寄存器)\n", ms.BR)
	fmt.Printf("  IX = 0x%04X (变址寄存器)\n", ms.IX)
	fmt.Println("\n  寄存器:")
	fmt.Printf("    R1 = %d\n", ms.Registers["R1"])
	fmt.Printf("    R2 = 0x%04X\n", ms.Registers["R2"])
	fmt.Printf("    R3 = %d\n", ms.Registers["R3"])
	fmt.Println("\n  内存:")
	fmt.Printf("    M[0x1000] = %d\n", ms.Memory[0x1000])
	fmt.Printf("    M[0x1100] = %d\n", ms.Memory[0x1100])
	fmt.Printf("    M[0x1200] = %d\n", ms.Memory[0x1200])
	fmt.Printf("    M[0x2000] = 0x%04X (指针)\n", ms.Memory[0x2000])
	fmt.Printf("    M[0x2100] = %d\n", ms.Memory[0x2100])
	fmt.Printf("    M[0x2500] = %d\n", ms.Memory[0x2500])

	// 定义测试指令
	instructions := []Instruction{
		{Opcode: "MOV", AddressingMode: Immediate, Operand: 99},
		{Opcode: "MOV", AddressingMode: Direct, Operand: 0x1000},
		{Opcode: "MOV", AddressingMode: Indirect, Operand: 0x2000},
		{Opcode: "MOV", AddressingMode: Register, RegisterName: "R1"},
		{Opcode: "MOV", AddressingMode: RegisterIndirect, RegisterName: "R2"},
		{Opcode: "MOV", AddressingMode: Relative, Operand: 0x100},
		{Opcode: "MOV", AddressingMode: Base, Operand: 0x100},
		{Opcode: "MOV", AddressingMode: Indexed, Operand: 0x2400},
	}

	fmt.Println("\n【寻址方式演示】")
	fmt.Println("指令格式: MOV R0, <操作数>")
	fmt.Println("任务: 将操作数加载到 R0 寄存器")

	for i, instr := range instructions {
		fmt.Printf("\n%d. %s 方式\n", i+1, instr.AddressingMode.String())

		// 获取操作数
		operand, accessCount, desc := ms.GetOperand(instr)

		fmt.Printf("   %s\n", desc)
		fmt.Printf("   → R0 ← %d\n", operand)
		fmt.Printf("   访存次数: %d\n", accessCount)
	}

	// 演示一个完整的计算例子
	fmt.Println("\n【完整示例：计算 C = A + B】")
	fmt.Println("\n假设: A 在地址 0x1000, B 在地址 0x1100, C 存放在地址 0x1200")

	fmt.Println("\n不同指令格式的实现:")

	// 三地址指令
	fmt.Println("\n1. 三地址指令 (需要 1 条指令)")
	fmt.Println("   ADD [0x1000], [0x1100], [0x1200]  ; M[0x1200] ← M[0x1000] + M[0x1100]")
	fmt.Println("   访存次数: 3 次 (读A, 读B, 写C)")

	// 二地址指令
	fmt.Println("\n2. 二地址指令 (需要 2 条指令)")
	fmt.Println("   MOV [0x1200], [0x1000]            ; M[0x1200] ← M[0x1000]")
	fmt.Println("   ADD [0x1200], [0x1100]            ; M[0x1200] ← M[0x1200] + M[0x1100]")
	fmt.Println("   访存次数: 5 次 (读A, 写C, 读C, 读B, 写C)")

	// 一地址指令
	fmt.Println("\n3. 一地址指令 (需要 3 条指令，使用累加器 ACC)")
	fmt.Println("   LOAD [0x1000]                     ; ACC ← M[0x1000]")
	fmt.Println("   ADD  [0x1100]                     ; ACC ← ACC + M[0x1100]")
	fmt.Println("   STORE [0x1200]                    ; M[0x1200] ← ACC")
	fmt.Println("   访存次数: 3 次 (读A, 读B, 写C)")

	// 零地址指令
	fmt.Println("\n4. 零地址指令 (需要 4 条指令，使用栈)")
	fmt.Println("   PUSH [0x1000]                     ; 栈顶 ← M[0x1000]")
	fmt.Println("   PUSH [0x1100]                     ; 栈顶 ← M[0x1100]")
	fmt.Println("   ADD                               ; 栈顶 ← 栈顶 + 次栈顶，弹出次栈顶")
	fmt.Println("   POP  [0x1200]                     ; M[0x1200] ← 栈顶，弹出栈顶")
	fmt.Println("   访存次数: 3 次 (读A, 读B, 写C) + 栈操作")

	// 寻址方式对比
	fmt.Println("\n【寻址方式性能对比】")
	fmt.Println("\n速度排序（从快到慢）：")
	fmt.Println("  1. 立即寻址         (0次访存) ⚡")
	fmt.Println("  2. 寄存器寻址       (0次访存) ⚡")
	fmt.Println("  3. 直接寻址         (1次访存)")
	fmt.Println("  4. 寄存器间接寻址   (1次访存)")
	fmt.Println("  5. 相对寻址         (1次访存)")
	fmt.Println("  6. 基址寻址         (1次访存)")
	fmt.Println("  7. 变址寻址         (1次访存)")
	fmt.Println("  8. 间接寻址         (2次访存) 🐌")

	fmt.Println("\n应用场景：")
	fmt.Println("  立即寻址     → 常量（如 MOV R1, #100）")
	fmt.Println("  直接寻址     → 全局变量")
	fmt.Println("  间接寻址     → 指针访问")
	fmt.Println("  寄存器寻址   → 临时变量")
	fmt.Println("  寄存器间接   → 指针、链表遍历")
	fmt.Println("  相对寻址     → 转移指令、位置无关代码")
	fmt.Println("  基址寻址     → 分段管理、程序浮动")
	fmt.Println("  变址寻址     → 数组访问、循环")

	fmt.Println("\n" + instructionSet408Summary())
}

// instructionSet408Summary 408 考试总结
func instructionSet408Summary() string {
	return `
╔════════════════════════════════════════════════════════════════╗
║                 408 考试要点总结 - 指令系统                   ║
╠════════════════════════════════════════════════════════════════╣
║ 1. 指令格式：                                                  ║
║    • 零地址: OP              (栈式机器)                       ║
║    • 一地址: OP | A1         (隐含 ACC)                       ║
║    • 二地址: OP | A1 | A2    (最常用)                         ║
║    • 三地址: OP | A1 | A2 | A3                                ║
║                                                                ║
║    计算 C = A + B 需要的指令数:                                ║
║      三地址: 1 条   二地址: 2 条   一地址: 3 条   零地址: 4 条 ║
║                                                                ║
║ 2. 寻址方式公式（必记）：                                      ║
║    立即寻址:     操作数 = D                                   ║
║    直接寻址:     EA = A,  操作数 = (A)                        ║
║    间接寻址:     EA = (A), 操作数 = ((A))                     ║
║    寄存器寻址:   操作数 = (Ri)                                ║
║    寄存器间接:   EA = (Ri), 操作数 = ((Ri))                   ║
║    相对寻址:     EA = (PC) + D                                ║
║    基址寻址:     EA = (BR) + D                                ║
║    变址寻址:     EA = (IX) + D                                ║
║                                                                ║
║    注: EA = 有效地址, D = 位移量/立即数, (X) = X的内容        ║
║                                                                ║
║ 3. 寻址方式访存次数：                                          ║
║    ┌──────────────┬──────────┬─────────────┐                 ║
║    │ 寻址方式     │ 访存次数 │ 速度        │                 ║
║    ├──────────────┼──────────┼─────────────┤                 ║
║    │ 立即         │    0     │ 最快 ⚡     │                 ║
║    │ 寄存器       │    0     │ 最快 ⚡     │                 ║
║    │ 直接         │    1     │ 快          │                 ║
║    │ 寄存器间接   │    1     │ 快          │                 ║
║    │ 相对/基址/变址│   1     │ 快          │                 ║
║    │ 间接         │    2     │ 慢 🐌       │                 ║
║    └──────────────┴──────────┴─────────────┘                 ║
║                                                                ║
║ 4. 基址寻址 vs 变址寻址：                                      ║
║    ┌─────────┬────────────────┬────────────────┐             ║
║    │ 特性    │ 基址寻址       │ 变址寻址       │             ║
║    ├─────────┼────────────────┼────────────────┤             ║
║    │ 寄存器值│ 系统设置，较大 │ 用户设置，较小 │             ║
║    │ 变化部分│ 指令中的位移量 │ 变址寄存器     │             ║
║    │ 主要用途│ 程序浮动、分段 │ 数组、循环     │             ║
║    └─────────┴────────────────┴────────────────┘             ║
║                                                                ║
║ 5. CISC vs RISC 对比：                                         ║
║    ┌──────────┬──────────────┬──────────────┐                ║
║    │ 特性     │ CISC         │ RISC         │                ║
║    ├──────────┼──────────────┼──────────────┤                ║
║    │ 指令数量 │ 多(200-300)  │ 少(50-100)   │                ║
║    │ 指令格式 │ 变长         │ 定长         │                ║
║    │ 寻址方式 │ 多(10+)      │ 少(3-5)      │                ║
║    │ CPI      │ 2-15         │ ~1           │                ║
║    │ 控制器   │ 微程序       │ 硬布线       │                ║
║    │ 访存指令 │ 任意指令     │ Load/Store   │                ║
║    │ 流水线   │ 困难         │ 容易         │                ║
║    │ 代表     │ x86          │ ARM, RISC-V  │                ║
║    └──────────┴──────────────┴──────────────┘                ║
║                                                                ║
║ 6. 寻址范围计算：                                              ║
║    • 直接寻址: 2^n (n=地址码位数)                             ║
║    • 间接寻址: 2^m (m=存储字长)                               ║
║    • 相对寻址: PC ± 2^(n-1) (n=位移量位数)                    ║
║    • 基址/变址: 取决于寄存器位数 + 位移量位数                 ║
║                                                                ║
║ 7. 典型考题：                                                  ║
║    (1) 计算题: 给定 PC、寄存器、内存状态，计算有效地址和操作数║
║    (2) 对比题: 比较不同寻址方式的速度、灵活性、寻址范围       ║
║    (3) 应用题: 判断某种应用场景适合哪种寻址方式               ║
║    (4) 计算题: 不同地址码指令完成同一任务需要的指令数         ║
║    (5) 分析题: CISC 和 RISC 的优缺点及适用场景                ║
║                                                                ║
║ 8. 解题技巧：                                                  ║
║    • 注意 PC 的值：当前指令地址 or 下一条指令地址？           ║
║    • 间接寻址：记得内容两次                                   ║
║    • 相对寻址：注意符号位，可正可负                           ║
║    • 基址/变址：弄清哪部分由系统管，哪部分由用户控制          ║
╚════════════════════════════════════════════════════════════════╝
`
}
