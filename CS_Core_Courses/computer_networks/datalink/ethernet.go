package datalink

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// EthernetFrame 以太网帧结构
// 对应 408 考点: 以太网帧格式
type EthernetFrame struct {
	Preamble  string // 前导码 (7 字节) + 帧起始定界符 SFD (1 字节)
	DestMAC   string // 目的 MAC 地址 (6 字节)
	SourceMAC string // 源 MAC 地址 (6 字节)
	Type      uint16 // 类型字段 (2 字节): 0x0800=IPv4, 0x0806=ARP
	Data      []byte // 数据字段 (46-1500 字节)
	FCS       uint32 // 帧校验序列 (4 字节, CRC-32)
}

// NewEthernetFrame 创建以太网帧
func NewEthernetFrame(destMAC, srcMAC string, frameType uint16, data []byte) *EthernetFrame {
	// 如果数据不足 46 字节,填充到 46 字节
	if len(data) < 46 {
		padding := make([]byte, 46-len(data))
		data = append(data, padding...)
	}

	return &EthernetFrame{
		Preamble:  "10101010...", // 前导码模拟
		DestMAC:   destMAC,
		SourceMAC: srcMAC,
		Type:      frameType,
		Data:      data,
		FCS:       calculateFCS(data), // 简化的 FCS 计算
	}
}

// calculateFCS 简化的 FCS 计算 (实际应使用 CRC-32)
func calculateFCS(data []byte) uint32 {
	var sum uint32
	for _, b := range data {
		sum += uint32(b)
	}
	return sum & 0xFFFFFFFF
}

// String 格式化输出以太网帧
func (f *EthernetFrame) String() string {
	var sb strings.Builder
	sb.WriteString("╔═══════════════════ 以太网帧 ═══════════════════╗\n")
	sb.WriteString(fmt.Sprintf("║ 前导码:     %s                     ║\n", f.Preamble))
	sb.WriteString(fmt.Sprintf("║ 目的MAC:    %-32s║\n", f.DestMAC))
	sb.WriteString(fmt.Sprintf("║ 源MAC:      %-32s║\n", f.SourceMAC))
	sb.WriteString(fmt.Sprintf("║ 类型:       0x%04X (%s)%-14s║\n",
		f.Type, getFrameTypeName(f.Type), ""))
	sb.WriteString(fmt.Sprintf("║ 数据长度:   %d 字节%-26s║\n", len(f.Data), ""))
	sb.WriteString(fmt.Sprintf("║ FCS:        0x%08X%-26s║\n", f.FCS, ""))
	sb.WriteString("╚═══════════════════════════════════════════════╝")
	return sb.String()
}

// getFrameTypeName 获取帧类型名称
func getFrameTypeName(frameType uint16) string {
	switch frameType {
	case 0x0800:
		return "IPv4"
	case 0x0806:
		return "ARP"
	case 0x86DD:
		return "IPv6"
	default:
		return "未知"
	}
}

// CSMACD CSMA/CD 协议模拟
// 对应 408 考点: 载波侦听多路访问/冲突检测
type CSMACD struct {
	Channel      bool       // 信道状态: true=忙, false=空闲
	Stations     []*Station // 所有站点
	CollisionLog []string   // 冲突日志
}

// Station 站点
type Station struct {
	Name         string // 站点名称
	MAC          string // MAC 地址
	BackoffCount int    // 退避次数
	MaxRetries   int    // 最大重传次数
}

// NewCSMACD 创建 CSMA/CD 模拟器
func NewCSMACD() *CSMACD {
	return &CSMACD{
		Channel:      false,
		Stations:     make([]*Station, 0),
		CollisionLog: make([]string, 0),
	}
}

// AddStation 添加站点
func (c *CSMACD) AddStation(name, mac string) {
	c.Stations = append(c.Stations, &Station{
		Name:         name,
		MAC:          mac,
		BackoffCount: 0,
		MaxRetries:   16, // 以太网标准: 最多重传 16 次
	})
}

// SendFrame 模拟站点发送帧
func (c *CSMACD) SendFrame(stationName string, destMAC string, data []byte) bool {
	// 查找站点
	var station *Station
	for _, s := range c.Stations {
		if s.Name == stationName {
			station = s
			break
		}
	}

	if station == nil {
		fmt.Printf("错误: 站点 %s 不存在\n", stationName)
		return false
	}

	fmt.Printf("\n[%s] 准备发送数据...\n", station.Name)

	// CSMA/CD 流程
	for attempt := 0; attempt <= station.MaxRetries; attempt++ {
		// 1. 载波侦听 (Carrier Sense)
		fmt.Printf("[%s] 第 %d 次尝试: 侦听信道...", station.Name, attempt+1)
		if c.Channel {
			fmt.Println(" 信道忙,等待...")
			time.Sleep(10 * time.Millisecond)
			continue
		}
		fmt.Println(" 信道空闲")

		// 2. 发送数据
		c.Channel = true
		fmt.Printf("[%s] 开始发送数据...\n", station.Name)

		// 模拟冲突检测 (Collision Detection)
		collision := c.simulateCollision()
		if collision {
			// 检测到冲突
			c.Channel = false
			c.CollisionLog = append(c.CollisionLog,
				fmt.Sprintf("[冲突] %s 在第 %d 次尝试时检测到冲突", station.Name, attempt+1))
			fmt.Printf("[%s] ✗ 检测到冲突!\n", station.Name)

			// 截断二进制指数退避 (Truncated Binary Exponential Backoff)
			k := attempt
			if k > 10 {
				k = 10 // 最多退避 2^10 个时隙
			}
			maxSlots := (1 << k) - 1 // 2^k - 1
			backoffSlots := rand.Intn(maxSlots + 1)
			fmt.Printf("[%s] 执行退避算法: k=%d, 随机退避 %d 个时隙\n",
				station.Name, k, backoffSlots)

			// 模拟退避时间
			time.Sleep(time.Duration(backoffSlots*10) * time.Millisecond)
			continue
		}

		// 3. 发送成功
		c.Channel = false
		frame := NewEthernetFrame(destMAC, station.MAC, 0x0800, data)
		fmt.Printf("[%s] ✓ 发送成功!\n", station.Name)
		fmt.Println(frame)
		return true
	}

	// 超过最大重传次数
	fmt.Printf("[%s] ✗ 超过最大重传次数 (%d),发送失败\n", station.Name, station.MaxRetries)
	return false
}

// simulateCollision 模拟冲突 (30% 概率)
func (c *CSMACD) simulateCollision() bool {
	return rand.Float32() < 0.3
}

// PrintCollisionLog 打印冲突日志
func (c *CSMACD) PrintCollisionLog() {
	fmt.Println("\n【冲突日志】")
	if len(c.CollisionLog) == 0 {
		fmt.Println("  无冲突")
		return
	}
	for i, log := range c.CollisionLog {
		fmt.Printf("  %d. %s\n", i+1, log)
	}
}

// EthernetExample 以太网协议示例
func EthernetExample() {
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("【数据链路层 - 以太网协议示例】")
	fmt.Println(strings.Repeat("─", 50))

	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())

	// 1. 以太网帧结构
	fmt.Println("\n1️⃣  以太网帧结构:")
	data := []byte("Hello, Ethernet!")
	frame := NewEthernetFrame(
		"AA:BB:CC:DD:EE:FF",
		"11:22:33:44:55:66",
		0x0800,
		data,
	)
	fmt.Println(frame)

	// 2. CSMA/CD 协议模拟
	fmt.Println("\n2️⃣  CSMA/CD 协议模拟:")
	fmt.Println("\n场景: 多个站点竞争发送数据")

	csma := NewCSMACD()
	csma.AddStation("站点A", "AA:BB:CC:DD:EE:01")
	csma.AddStation("站点B", "AA:BB:CC:DD:EE:02")

	fmt.Println("\n初始化:")
	for _, s := range csma.Stations {
		fmt.Printf("  • %s (MAC: %s)\n", s.Name, s.MAC)
	}

	// 站点 A 发送数据
	fmt.Println("\n" + strings.Repeat("═", 50))
	testData := []byte("Test data from Station A")
	csma.SendFrame("站点A", "FF:FF:FF:FF:FF:FF", testData)

	// 显示冲突日志
	csma.PrintCollisionLog()

	// 408 考点提示
	fmt.Println("\n📚 408 考点总结:")
	fmt.Println("  ✓ 以太网帧格式: 前导码(8) + 目的MAC(6) + 源MAC(6) + 类型(2) + 数据(46-1500) + FCS(4)")
	fmt.Println("  ✓ 最小帧长: 64 字节 (用于冲突检测)")
	fmt.Println("  ✓ CSMA/CD: 1-坚持 CSMA, 发送时检测冲突")
	fmt.Println("  ✓ 退避算法: 截断二进制指数退避, k=min(重传次数, 10)")
	fmt.Println("  ✓ 最大重传次数: 16 次,超过则放弃")
}
