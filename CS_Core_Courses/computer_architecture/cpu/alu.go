package cpu

import (
	"fmt"
	"math/bits"
)

// ALUOperation ALU操作类型
type ALUOperation int

const (
	ALUNop ALUOperation = iota // 无操作
	ALUAdd                     // 加法
	ALUSub                     // 减法
	ALUMul                     // 乘法
	ALUDiv                     // 除法
	ALUAnd                     // 按位与
	ALUOr                      // 按位或
	ALUXor                     // 按位异或
	ALUNot                     // 按位取反
	ALUShl                     // 左移
	ALUShr                     // 右移
	ALUCmp                     // 比较
)

// String 返回操作的字符串表示
func (op ALUOperation) String() string {
	switch op {
	case ALUNop:
		return "NOP"
	case ALUAdd:
		return "ADD"
	case ALUSub:
		return "SUB"
	case ALUMul:
		return "MUL"
	case ALUDiv:
		return "DIV"
	case ALUAnd:
		return "AND"
	case ALUOr:
		return "OR"
	case ALUXor:
		return "XOR"
	case ALUNot:
		return "NOT"
	case ALUShl:
		return "SHL"
	case ALUShr:
		return "SHR"
	case ALUCmp:
		return "CMP"
	default:
		return "UNKNOWN"
	}
}

// ALUResult ALU运算结果
type ALUResult struct {
	Result       int64  // 运算结果
	Zero         bool   // 零标志位
	Carry        bool   // 进位标志位
	Negative     bool   // 负数标志位
	Overflow     bool   // 溢出标志位
	Parity       bool   // 奇偶标志位 (偶数个1为true)
	HasError     bool   // 是否有错误
	ErrorMessage string // 错误信息
}

// ALU 算术逻辑单元
type ALU struct {
	bitWidth int // 位数 (8, 16, 32, 64)
}

// NewALU 创建ALU
func NewALU(bitWidth int) *ALU {
	return &ALU{
		bitWidth: bitWidth,
	}
}

// Execute 执行ALU操作
func (alu *ALU) Execute(op ALUOperation, operand1, operand2 int64) ALUResult {
	result := ALUResult{
		Result:       0,
		Zero:         false,
		Carry:        false,
		Negative:     false,
		Overflow:     false,
		Parity:       false,
		HasError:     false,
		ErrorMessage: "",
	}

	// 掩码用于限制结果位数
	var mask uint64
	switch alu.bitWidth {
	case 8:
		mask = 0xFF
	case 16:
		mask = 0xFFFF
	case 32:
		mask = 0xFFFFFFFF
	case 64:
		mask = ^uint64(0) // 所有位为1
	default:
		result.HasError = true
		result.ErrorMessage = fmt.Sprintf("不支持的位宽: %d", alu.bitWidth)
		return result
	}

	signBit := uint64(1) << 63 // 最高位标志

	switch op {
	case ALUNop:
		result.Result = operand1

	case ALUAdd:
		result.Result = operand1 + operand2
		// 检查进位
		if uint64(operand1)+uint64(operand2) > mask {
			result.Carry = true
		}
		// 检查溢出 (假设两个操作数符号相同，但结果符号不同)
		if uint64(operand1^operand2)&signBit == 0 && uint64(operand1^result.Result)&signBit != 0 {
			result.Overflow = true
		}

	case ALUSub:
		result.Result = operand1 - operand2
		// 检查借位
		if uint64(operand1) < uint64(operand2) {
			result.Carry = true
		}
		// 检查溢出
		if uint64(operand1^operand2)&signBit != 0 && uint64(operand1^result.Result)&signBit != 0 {
			result.Overflow = true
		}

	case ALUMul:
		result.Result = operand1 * operand2
		// 检查溢出 (乘法溢出比较复杂，这里简化处理)
		if operand2 != 0 && (result.Result/operand2) != operand1 {
			result.Overflow = true
		}

	case ALUDiv:
		if operand2 == 0 {
			result.HasError = true
			result.ErrorMessage = "除零错误"
		} else {
			result.Result = operand1 / operand2
		}

	case ALUAnd:
		result.Result = operand1 & operand2

	case ALUOr:
		result.Result = operand1 | operand2

	case ALUXor:
		result.Result = operand1 ^ operand2

	case ALUNot:
		result.Result = ^operand1

	case ALUShl:
		if operand2 < 0 {
			result.HasError = true
			result.ErrorMessage = "左移位数不能为负"
		} else {
			result.Result = operand1 << uint64(operand2)
			// 检查是否有位移出位宽范围
			if operand1>>(uint64(alu.bitWidth)-uint64(operand2)) != 0 {
				result.Carry = true
			}
		}

	case ALUShr:
		if operand2 < 0 {
			result.HasError = true
			result.ErrorMessage = "右移位数不能为负"
		} else {
			result.Result = operand1 >> uint64(operand2)
		}

	case ALUCmp:
		// 比较操作，不返回结果，只设置标志位
		result.Result = 0
		diff := operand1 - operand2
		operand1 = diff // 为了复用下面的标志位设置代码
		operand2 = 0

	default:
		result.HasError = true
		result.ErrorMessage = fmt.Sprintf("不支持的操作: %s", op)
		return result
	}

	// 应用掩码限制结果范围
	result.Result &= int64(mask)

	// 设置标志位
	result.Zero = result.Result == 0
	result.Negative = (result.Result & (1 << (alu.bitWidth - 1))) != 0

	// 计算奇偶标志位 (偶数个1为true)
	ones := bits.OnesCount64(uint64(result.Result))
	result.Parity = (ones % 2) == 0

	return result
}

// PerformAddition 执行加法操作
func (alu *ALU) PerformAddition(a, b int64) ALUResult {
	return alu.Execute(ALUAdd, a, b)
}

// PerformSubtraction 执行减法操作
func (alu *ALU) PerformSubtraction(a, b int64) ALUResult {
	return alu.Execute(ALUSub, a, b)
}

// PerformMultiplication 执行乘法操作
func (alu *ALU) PerformMultiplication(a, b int64) ALUResult {
	return alu.Execute(ALUMul, a, b)
}

// PerformDivision 执行除法操作
func (alu *ALU) PerformDivision(a, b int64) ALUResult {
	return alu.Execute(ALUDiv, a, b)
}

// PerformLogicalOperation 执行逻辑操作
func (alu *ALU) PerformLogicalOperation(op ALUOperation, a, b int64) ALUResult {
	return alu.Execute(op, a, b)
}

// PerformShift 执行移位操作
func (alu *ALU) PerformShift(op ALUOperation, a, shift int64) ALUResult {
	return alu.Execute(op, a, shift)
}

// Compare 比较两个数
func (alu *ALU) Compare(a, b int64) ALUResult {
	return alu.Execute(ALUCmp, a, b)
}

// PrintResult 打印ALU结果
func (result ALUResult) Print(op ALUOperation, operand1, operand2 int64) {
	if result.HasError {
		fmt.Printf("ALU %s: ERROR - %s\n", op, result.ErrorMessage)
		return
	}

	fmt.Printf("ALU %s: %d %s %d = %d (0x%X)\n",
		op, operand1, op, operand2, result.Result, result.Result)

	flags := []string{}
	if result.Zero {
		flags = append(flags, "Z")
	}
	if result.Carry {
		flags = append(flags, "C")
	}
	if result.Negative {
		flags = append(flags, "N")
	}
	if result.Overflow {
		flags = append(flags, "V")
	}
	if result.Parity {
		flags = append(flags, "P")
	}

	if len(flags) > 0 {
		fmt.Printf("  Flags: %v\n", flags)
	}
}

// 示例函数
func ALUExample() {
	fmt.Println("=== 算术逻辑单元 (ALU) 示例 ===")

	// 创建32位ALU
	alu := NewALU(32)

	fmt.Println("1. 算术运算:")
	// 加法
	result := alu.PerformAddition(100, 25)
	result.Print(ALUAdd, 100, 25)

	// 减法
	result = alu.PerformSubtraction(100, 25)
	result.Print(ALUSub, 100, 25)

	// 乘法
	result = alu.PerformMultiplication(100, 25)
	result.Print(ALUMul, 100, 25)

	// 除法
	result = alu.PerformDivision(100, 25)
	result.Print(ALUDiv, 100, 25)

	// 除零错误
	result = alu.PerformDivision(100, 0)
	result.Print(ALUDiv, 100, 0)

	fmt.Println("\n2. 逻辑运算:")
	// 按位与
	result = alu.PerformLogicalOperation(ALUAnd, 0b11001100, 0b10101010)
	result.Print(ALUAnd, 0b11001100, 0b10101010)

	// 按位或
	result = alu.PerformLogicalOperation(ALUOr, 0b11001100, 0b10101010)
	result.Print(ALUOr, 0b11001100, 0b10101010)

	// 按位异或
	result = alu.PerformLogicalOperation(ALUXor, 0b11001100, 0b10101010)
	result.Print(ALUXor, 0b11001100, 0b10101010)

	// 按位取反
	result = alu.PerformLogicalOperation(ALUNot, 0b11001100, 0)
	result.Print(ALUNot, 0b11001100, 0)

	fmt.Println("\n3. 移位运算:")
	// 左移
	result = alu.PerformShift(ALUShl, 0b00001111, 2)
	result.Print(ALUShl, 0b00001111, 2)

	// 右移
	result = alu.PerformShift(ALUShr, 0b11110000, 2)
	result.Print(ALUShr, 0b11110000, 2)

	fmt.Println("\n4. 比较运算:")
	// 大于
	result = alu.Compare(100, 50)
	fmt.Printf("CMP 100, 50: ")
	if result.Zero {
		fmt.Print("等于 ")
	} else if result.Negative {
		fmt.Print("小于 ")
	} else {
		fmt.Print("大于 ")
	}
	fmt.Printf("(Flags: Z=%t, N=%t)\n", result.Zero, result.Negative)

	// 等于
	result = alu.Compare(100, 100)
	fmt.Printf("CMP 100, 100: ")
	if result.Zero {
		fmt.Print("等于 ")
	} else if result.Negative {
		fmt.Print("小于 ")
	} else {
		fmt.Print("大于 ")
	}
	fmt.Printf("(Flags: Z=%t, N=%t)\n", result.Zero, result.Negative)

	// 小于
	result = alu.Compare(50, 100)
	fmt.Printf("CMP 50, 100: ")
	if result.Zero {
		fmt.Print("等于 ")
	} else if result.Negative {
		fmt.Print("小于 ")
	} else {
		fmt.Print("大于 ")
	}
	fmt.Printf("(Flags: Z=%t, N=%t)\n", result.Zero, result.Negative)

	fmt.Println("\n5. 边界测试:")
	// 溢出测试 (假设8位ALU)
	alu8 := NewALU(8)
	result = alu8.PerformAddition(200, 100) // 200 + 100 = 300, 8位会溢出
	fmt.Printf("8位ALU ADD 200, 100: Result=%d, Overflow=%t\n", result.Result, result.Overflow)

	// 进位测试
	result = alu8.PerformAddition(200, 56) // 200 + 56 = 256, 8位有进位
	fmt.Printf("8位ALU ADD 200, 56: Result=%d, Carry=%t\n", result.Result, result.Carry)

	fmt.Println()
}
