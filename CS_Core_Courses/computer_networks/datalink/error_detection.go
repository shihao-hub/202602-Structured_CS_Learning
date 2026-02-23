package datalink

import (
	"fmt"
	"strings"
)

// CRC 循环冗余校验
// 对应 408 考点: CRC 校验码计算

// CRCCalculator CRC 计算器
type CRCCalculator struct {
	Polynomial string // 生成多项式 (二进制字符串)
	PolyBits   int    // 生成多项式位数
}

// NewCRCCalculator 创建 CRC 计算器
// polynomial: 生成多项式,如 "11001" 表示 x^4 + x^3 + 1
func NewCRCCalculator(polynomial string) *CRCCalculator {
	return &CRCCalculator{
		Polynomial: polynomial,
		PolyBits:   len(polynomial),
	}
}

// Calculate 计算 CRC 校验码
// data: 原始数据 (二进制字符串)
// 返回: CRC 校验码 (二进制字符串)
func (c *CRCCalculator) Calculate(data string) string {
	// 1. 在数据后面添加 r 个 0 (r = 生成多项式位数 - 1)
	r := c.PolyBits - 1
	paddedData := data + strings.Repeat("0", r)

	fmt.Printf("\n【CRC 计算过程】\n")
	fmt.Printf("原始数据:     %s\n", data)
	fmt.Printf("生成多项式:   %s\n", c.Polynomial)
	fmt.Printf("补 %d 个 0:    %s\n", r, paddedData)

	// 2. 模 2 除法
	remainder := c.modulo2Division(paddedData, c.Polynomial)

	fmt.Printf("余数 (CRC):   %s\n", remainder)

	return remainder
}

// modulo2Division 模 2 除法 (多项式除法)
func (c *CRCCalculator) modulo2Division(dividend, divisor string) string {
	// 转换为字节切片便于操作
	result := []byte(dividend)
	divisorLen := len(divisor)

	// 逐位异或
	for i := 0; i <= len(result)-divisorLen; i++ {
		// 如果当前位为 1,执行异或
		if result[i] == '1' {
			for j := 0; j < divisorLen; j++ {
				if result[i+j] == divisor[j] {
					result[i+j] = '0'
				} else {
					result[i+j] = '1'
				}
			}
		}
	}

	// 返回最后 r 位 (余数)
	r := divisorLen - 1
	return string(result[len(result)-r:])
}

// Encode 编码: 原始数据 + CRC 校验码
func (c *CRCCalculator) Encode(data string) string {
	crc := c.Calculate(data)
	return data + crc
}

// Verify 验证接收到的数据是否正确
func (c *CRCCalculator) Verify(receivedData string) bool {
	// 对接收到的数据进行模 2 除法,余数为 0 则无错
	remainder := c.modulo2Division(receivedData, c.Polynomial)
	// 检查余数是否全为 0
	for _, bit := range remainder {
		if bit == '1' {
			return false
		}
	}
	return true
}

// HammingCode 海明码
// 对应 408 考点: 海明码编码、检错、纠错

// HammingCodec 海明码编解码器
type HammingCodec struct {
	DataBits   int // 数据位数
	ParityBits int // 校验位数 (满足 2^r >= m + r + 1)
	TotalBits  int // 总位数
}

// NewHammingCodec 创建海明码编解码器
// dataBits: 数据位数
func NewHammingCodec(dataBits int) *HammingCodec {
	// 计算所需的校验位数 r: 2^r >= m + r + 1
	parityBits := 0
	for (1 << parityBits) < (dataBits + parityBits + 1) {
		parityBits++
	}

	return &HammingCodec{
		DataBits:   dataBits,
		ParityBits: parityBits,
		TotalBits:  dataBits + parityBits,
	}
}

// Encode 海明码编码
// data: 数据位 (二进制字符串)
// 返回: 海明码 (二进制字符串)
func (h *HammingCodec) Encode(data string) string {
	if len(data) != h.DataBits {
		fmt.Printf("错误: 数据位数应为 %d,实际为 %d\n", h.DataBits, len(data))
		return ""
	}

	fmt.Printf("\n【海明码编码过程】\n")
	fmt.Printf("数据位数: %d, 校验位数: %d, 总位数: %d\n", h.DataBits, h.ParityBits, h.TotalBits)

	// 初始化海明码 (1-based 索引,位置 0 不使用)
	hamming := make([]byte, h.TotalBits+1)
	dataIndex := 0

	// 1. 填充数据位 (跳过 2^i 位置)
	for i := 1; i <= h.TotalBits; i++ {
		if isPowerOfTwo(i) {
			hamming[i] = '0' // 校验位先设为 0
		} else {
			hamming[i] = data[dataIndex]
			dataIndex++
		}
	}

	fmt.Printf("初始布局: ")
	h.printHammingBits(hamming)

	// 2. 计算校验位
	for p := 0; p < h.ParityBits; p++ {
		pos := 1 << p // 2^p
		parity := 0

		// 统计该校验位负责的位
		for i := 1; i <= h.TotalBits; i++ {
			if (i & pos) != 0 { // i 的二进制表示中第 p 位为 1
				if hamming[i] == '1' {
					parity ^= 1
				}
			}
		}

		hamming[pos] = byte('0' + parity)
	}

	fmt.Printf("填充校验位: ")
	h.printHammingBits(hamming)

	// 返回海明码 (去掉位置 0)
	return string(hamming[1:])
}

// Decode 海明码解码与纠错
// received: 接收到的海明码
// 返回: (原始数据, 错误位置)
func (h *HammingCodec) Decode(received string) (string, int) {
	if len(received) != h.TotalBits {
		fmt.Printf("错误: 海明码位数应为 %d,实际为 %d\n", h.TotalBits, len(received))
		return "", -1
	}

	fmt.Printf("\n【海明码解码与纠错】\n")
	fmt.Printf("接收码字: %s\n", received)

	// 转换为 1-based 数组
	hamming := make([]byte, h.TotalBits+1)
	hamming[0] = '0'
	copy(hamming[1:], []byte(received))

	// 1. 计算校验子 (Syndrome)
	syndrome := 0
	for p := 0; p < h.ParityBits; p++ {
		pos := 1 << p
		parity := 0

		for i := 1; i <= h.TotalBits; i++ {
			if (i & pos) != 0 {
				if hamming[i] == '1' {
					parity ^= 1
				}
			}
		}

		if parity != 0 {
			syndrome |= pos
		}
	}

	fmt.Printf("校验子 (Syndrome): %d (二进制: %b)\n", syndrome, syndrome)

	// 2. 纠错
	errorPos := syndrome
	if errorPos > 0 {
		fmt.Printf("检测到错误位置: 第 %d 位\n", errorPos)
		// 翻转错误位
		if hamming[errorPos] == '0' {
			hamming[errorPos] = '1'
		} else {
			hamming[errorPos] = '0'
		}
		fmt.Printf("纠正后码字: ")
		h.printHammingBits(hamming)
	} else {
		fmt.Println("✓ 未检测到错误")
	}

	// 3. 提取数据位
	var data strings.Builder
	for i := 1; i <= h.TotalBits; i++ {
		if !isPowerOfTwo(i) {
			data.WriteByte(hamming[i])
		}
	}

	return data.String(), errorPos
}

// isPowerOfTwo 判断是否为 2 的幂
func isPowerOfTwo(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}

// printHammingBits 打印海明码位
func (h *HammingCodec) printHammingBits(bits []byte) {
	for i := 1; i <= h.TotalBits; i++ {
		if isPowerOfTwo(i) {
			fmt.Printf("[%c]", bits[i]) // 校验位用方括号
		} else {
			fmt.Printf(" %c ", bits[i]) // 数据位
		}
	}
	fmt.Println()
}

// ErrorDetectionExample 差错检测示例
func ErrorDetectionExample() {
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("【数据链路层 - 差错检测示例】")
	fmt.Println(strings.Repeat("─", 50))

	// 1. CRC 校验示例
	fmt.Println("\n1️⃣  CRC (循环冗余校验):")
	fmt.Println("\n例 1: CRC-4 (生成多项式: x^4 + x^3 + 1 = 11001)")
	crc := NewCRCCalculator("11001")
	data1 := "101101"
	encoded1 := crc.Encode(data1)
	fmt.Printf("\n编码结果: %s\n", encoded1)

	// 验证正确数据
	fmt.Println("\n验证正确接收:")
	isValid := crc.Verify(encoded1)
	fmt.Printf("数据: %s, 校验结果: %v\n", encoded1, isValid)

	// 验证错误数据
	fmt.Println("\n验证错误接收 (第 5 位出错):")
	errorData := "1011010001" // 人为制造错误
	isValid = crc.Verify(errorData)
	fmt.Printf("数据: %s, 校验结果: %v\n", errorData, isValid)

	// 2. 海明码示例
	fmt.Println("\n\n2️⃣  海明码 (Hamming Code):")
	fmt.Println("\n例 1: 4 位数据的海明码编码")
	hamming := NewHammingCodec(4)
	data2 := "1011"
	encoded2 := hamming.Encode(data2)
	fmt.Printf("\n编码结果: %s\n", encoded2)

	// 无错误解码
	fmt.Println("\n场景 1: 无错误接收")
	decoded, errorPos := hamming.Decode(encoded2)
	fmt.Printf("解码数据: %s, 错误位置: %d\n", decoded, errorPos)

	// 单比特错误解码与纠正
	fmt.Println("\n场景 2: 第 3 位出错")
	// 翻转第 3 位
	errorBits := []byte(encoded2)
	if errorBits[2] == '0' {
		errorBits[2] = '1'
	} else {
		errorBits[2] = '0'
	}
	errorEncoded := string(errorBits)
	decoded, errorPos = hamming.Decode(errorEncoded)
	fmt.Printf("解码数据: %s, 错误位置: %d\n", decoded, errorPos)

	// 408 考点提示
	fmt.Println("\n📚 408 考点总结:")
	fmt.Println("  ✓ CRC: 生成多项式除法,余数作为校验码")
	fmt.Println("  ✓ CRC 可检测所有奇数位错,所有双比特错,所有小于 r 位的突发错")
	fmt.Println("  ✓ 海明码: 2^r >= m + r + 1 (m 为数据位, r 为校验位)")
	fmt.Println("  ✓ 海明码校验位位置: 2^0, 2^1, 2^2, ... (1, 2, 4, 8, ...)")
	fmt.Println("  ✓ 海明码可检测 2 位错,纠正 1 位错")
	fmt.Println("  ✓ 海明距离 d=3 时,检错能力 d-1=2,纠错能力 (d-1)/2=1")
}
