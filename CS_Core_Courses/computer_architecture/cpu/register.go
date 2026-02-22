package cpu

import (
	"fmt"
)

// RegisterType 寄存器类型
type RegisterType int

const (
	RegisterGeneral             RegisterType = iota // 通用寄存器
	RegisterProgramCounter                          // 程序计数器
	RegisterInstructionRegister                     // 指令寄存器
	RegisterAccumulator                             // 累加器
	RegisterStackPointer                            // 栈指针
	RegisterBasePointer                             // 基址指针
	RegisterFlag                                    // 标志寄存器
)

// String 返回寄存器类型的字符串表示
func (rt RegisterType) String() string {
	switch rt {
	case RegisterGeneral:
		return "General"
	case RegisterProgramCounter:
		return "PC"
	case RegisterInstructionRegister:
		return "IR"
	case RegisterAccumulator:
		return "ACC"
	case RegisterStackPointer:
		return "SP"
	case RegisterBasePointer:
		return "BP"
	case RegisterFlag:
		return "FLAG"
	default:
		return "Unknown"
	}
}

// Register 寄存器结构
type Register struct {
	Name  string
	Type  RegisterType
	Value int64
	Size  int // 位数
}

// NewRegister 创建新寄存器
func NewRegister(name string, regType RegisterType, size int) *Register {
	return &Register{
		Name:  name,
		Type:  regType,
		Value: 0,
		Size:  size,
	}
}

// SetValue 设置寄存器值
func (r *Register) SetValue(value int64) {
	// 根据寄存器大小进行掩码操作
	switch r.Size {
	case 8:
		r.Value = value & 0xFF
	case 16:
		r.Value = value & 0xFFFF
	case 32:
		r.Value = value & 0xFFFFFFFF
	case 64:
		r.Value = value // 64位不需要掩码
	default:
		r.Value = value
	}
}

// GetValue 获取寄存器值
func (r *Register) GetValue() int64 {
	return r.Value
}

// GetSignedValue 获取有符号值
func (r *Register) GetSignedValue() int64 {
	// 根据寄存器大小进行符号扩展
	switch r.Size {
	case 8:
		if r.Value&0x80 != 0 {
			return r.Value | ^int64(0xFF)
		}
	case 16:
		if r.Value&0x8000 != 0 {
			return r.Value | ^int64(0xFFFF)
		}
	case 32:
		if r.Value&0x80000000 != 0 {
			return r.Value | ^int64(0xFFFFFFFF)
		}
	}
	return r.Value
}

// Increment 递增
func (r *Register) Increment() {
	r.SetValue(r.Value + 1)
}

// Decrement 递减
func (r *Register) Decrement() {
	r.SetValue(r.Value - 1)
}

// Add 加法
func (r *Register) Add(value int64) {
	r.SetValue(r.Value + value)
}

// IsZero 检查是否为零
func (r *Register) IsZero() bool {
	return r.Value == 0
}

// Print 打印寄存器信息
func (r *Register) Print() {
	fmt.Printf("Register: %-4s Type: %-8s Size: %-2d Value: 0x%X (%d)\n",
		r.Name, r.Type, r.Size, r.Value, r.Value)
}

// RegisterFile 寄存器文件
type RegisterFile struct {
	registers map[string]*Register
}

// NewRegisterFile 创建寄存器文件
func NewRegisterFile() *RegisterFile {
	rf := &RegisterFile{
		registers: make(map[string]*Register),
	}

	// 创建常用寄存器
	// 通用寄存器 (32位)
	for i := 0; i < 8; i++ {
		name := fmt.Sprintf("R%d", i)
		rf.registers[name] = NewRegister(name, RegisterGeneral, 32)
	}

	// 特殊寄存器
	rf.registers["PC"] = NewRegister("PC", RegisterProgramCounter, 32)
	rf.registers["IR"] = NewRegister("IR", RegisterInstructionRegister, 32)
	rf.registers["ACC"] = NewRegister("ACC", RegisterAccumulator, 32)
	rf.registers["SP"] = NewRegister("SP", RegisterStackPointer, 32)
	rf.registers["BP"] = NewRegister("BP", RegisterBasePointer, 32)
	rf.registers["FLAG"] = NewRegister("FLAG", RegisterFlag, 16)

	return rf
}

// GetRegister 获取寄存器
func (rf *RegisterFile) GetRegister(name string) *Register {
	return rf.registers[name]
}

// SetRegister 设置寄存器值
func (rf *RegisterFile) SetRegister(name string, value int64) bool {
	if reg, exists := rf.registers[name]; exists {
		reg.SetValue(value)
		return true
	}
	return false
}

// GetRegisterValue 获取寄存器值
func (rf *RegisterFile) GetRegisterValue(name string) int64 {
	if reg, exists := rf.registers[name]; exists {
		return reg.GetValue()
	}
	return 0
}

// CopyRegister 复制寄存器
func (rf *RegisterFile) CopyRegister(dest, src string) bool {
	if destReg, destExists := rf.registers[dest]; destExists {
		if srcReg, srcExists := rf.registers[src]; srcExists {
			destReg.SetValue(srcReg.GetValue())
			return true
		}
	}
	return false
}

// PrintAll 打印所有寄存器
func (rf *RegisterFile) PrintAll() {
	fmt.Println("=== Register File ===")
	for _, reg := range rf.registers {
		reg.Print()
	}
	fmt.Println()
}

// PrintGeneralRegisters 打印通用寄存器
func (rf *RegisterFile) PrintGeneralRegisters() {
	fmt.Println("=== General Registers ===")
	for i := 0; i < 8; i++ {
		name := fmt.Sprintf("R%d", i)
		reg := rf.registers[name]
		fmt.Printf("R%d: 0x%-8X ", i, reg.GetValue())
		if (i+1)%4 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

// PrintSpecialRegisters 打印特殊寄存器
func (rf *RegisterFile) PrintSpecialRegisters() {
	fmt.Println("=== Special Registers ===")
	specialRegs := []string{"PC", "IR", "ACC", "SP", "BP", "FLAG"}
	for _, name := range specialRegs {
		reg := rf.registers[name]
		fmt.Printf("%-4s: 0x%-8X ", name, reg.GetValue())
	}
	fmt.Println()
}

// Flag 标志位定义
const (
	FlagZero      = 1 << 0 // 零标志位
	FlagCarry     = 1 << 1 // 进位标志位
	FlagNegative  = 1 << 2 // 负数标志位
	FlagOverflow  = 1 << 3 // 溢出标志位
	FlagParity    = 1 << 4 // 奇偶标志位
	FlagInterrupt = 1 << 5 // 中断标志位
)

// SetFlag 设置标志位
func (rf *RegisterFile) SetFlag(flag int, set bool) {
	flagReg := rf.registers["FLAG"]
	if set {
		flagReg.SetValue(flagReg.GetValue() | int64(flag))
	} else {
		flagReg.SetValue(flagReg.GetValue() &^ int64(flag))
	}
}

// GetFlag 获取标志位
func (rf *RegisterFile) GetFlag(flag int) bool {
	flagReg := rf.registers["FLAG"]
	return (flagReg.GetValue() & int64(flag)) != 0
}

// UpdateFlags 根据计算结果更新标志位
func (rf *RegisterFile) UpdateFlags(result int64) {
	// 零标志位
	rf.SetFlag(FlagZero, result == 0)

	// 负数标志位 (假设32位)
	rf.SetFlag(FlagNegative, (result&0x80000000) != 0)

	// 奇偶标志位 (假设8位结果)
	parity := true
	for i := 0; i < 8; i++ {
		if (result>>i)&1 == 1 {
			parity = !parity
		}
	}
	rf.SetFlag(FlagParity, !parity)
}

// 示例函数
func RegisterExample() {
	fmt.Println("=== 寄存器 (Register) 示例 ===")

	// 创建寄存器文件
	rf := NewRegisterFile()

	fmt.Println("1. 寄存器基本操作:")
	// 设置寄存器值
	rf.SetRegister("R0", 0x12345678)
	rf.SetRegister("R1", -100)
	rf.SetRegister("PC", 0x1000)

	// 打印寄存器
	fmt.Println("R0:", rf.GetRegister("R0"))
	fmt.Println("R1:", rf.GetRegister("R1"))
	fmt.Println("R1 signed:", rf.GetRegister("R1").GetSignedValue())

	fmt.Println("\n2. 寄存器操作:")
	// 寄存器操作
	r0 := rf.GetRegister("R0")
	fmt.Println("R0 original:", r0.GetValue())
	r0.Increment()
	fmt.Println("R0 after increment:", r0.GetValue())
	r0.Add(0x100)
	fmt.Println("R0 after adding 0x100:", r0.GetValue())

	fmt.Println("\n3. 寄存器复制:")
	// 复制寄存器
	rf.CopyRegister("R2", "R0")
	fmt.Println("R2 after copy from R0:", rf.GetRegister("R2").GetValue())

	fmt.Println("\n4. 标志位操作:")
	// 标志位操作
	rf.UpdateFlags(0)
	fmt.Println("Zero flag:", rf.GetFlag(FlagZero))
	fmt.Println("Negative flag:", rf.GetFlag(FlagNegative))

	rf.UpdateFlags(-1)
	fmt.Println("After setting result = -1:")
	fmt.Println("Zero flag:", rf.GetFlag(FlagZero))
	fmt.Println("Negative flag:", rf.GetFlag(FlagNegative))
	fmt.Println("Parity flag:", rf.GetFlag(FlagParity))

	fmt.Println("\n5. 寄存器文件状态:")
	// 打印寄存器状态
	rf.PrintGeneralRegisters()
	rf.PrintSpecialRegisters()

	// 模拟程序执行过程
	fmt.Println("\n6. 模拟程序执行:")
	fmt.Println("初始状态:")
	rf.PrintAll()

	// 模拟执行几条指令
	fmt.Println("执行: MOV R0, #42")
	rf.SetRegister("R0", 42)
	rf.UpdateFlags(42)

	fmt.Println("执行: ADD R0, #10")
	r0 = rf.GetRegister("R0")
	r0.Add(10)
	rf.UpdateFlags(r0.GetValue())

	fmt.Println("执行: MOV R1, R0")
	rf.CopyRegister("R1", "R0")

	fmt.Println("执行: INC PC")
	rf.GetRegister("PC").Increment()

	fmt.Println("执行后状态:")
	rf.PrintAll()
	fmt.Println()
}
