package network

import (
	"fmt"
	"math"
	"strings"
)

// IPv4Header IPv4 数据报头部结构
// 对应 408 考点: IP 数据报格式
type IPv4Header struct {
	Version        uint8  // 版本号 (4 位)
	IHL            uint8  // 首部长度 (4 位, 单位: 4字节)
	TypeOfService  uint8  // 服务类型 (8 位)
	TotalLength    uint16 // 总长度 (16 位, 单位: 字节)
	Identification uint16 // 标识 (16 位)
	Flags          uint8  // 标志 (3 位: 保留位, DF, MF)
	FragmentOffset uint16 // 片偏移 (13 位, 单位: 8字节)
	TTL            uint8  // 生存时间 (8 位)
	Protocol       uint8  // 协议 (8 位: 6=TCP, 17=UDP, 1=ICMP)
	HeaderChecksum uint16 // 首部检验和 (16 位)
	SourceIP       string // 源 IP 地址 (32 位)
	DestIP         string // 目的 IP 地址 (32 位)
}

// NewIPv4Header 创建 IPv4 头部
func NewIPv4Header(srcIP, dstIP string, totalLen uint16, protocol uint8) *IPv4Header {
	return &IPv4Header{
		Version:        4,
		IHL:            5, // 20 字节 = 5 * 4
		TypeOfService:  0,
		TotalLength:    totalLen,
		Identification: 12345,
		Flags:          0,
		FragmentOffset: 0,
		TTL:            64,
		Protocol:       protocol,
		HeaderChecksum: 0,
		SourceIP:       srcIP,
		DestIP:         dstIP,
	}
}

// String 格式化输出 IP 头部信息
func (h *IPv4Header) String() string {
	var sb strings.Builder
	sb.WriteString("═══════════════ IPv4 数据报头部 ═══════════════\n")
	sb.WriteString(fmt.Sprintf("版本: %d | 首部长度: %d (字节: %d)\n", h.Version, h.IHL, h.IHL*4))
	sb.WriteString(fmt.Sprintf("总长度: %d 字节 | 标识: %d\n", h.TotalLength, h.Identification))
	sb.WriteString(fmt.Sprintf("标志: DF=%d MF=%d | 片偏移: %d\n", (h.Flags>>1)&1, h.Flags&1, h.FragmentOffset))
	sb.WriteString(fmt.Sprintf("TTL: %d | 协议: %d\n", h.TTL, h.Protocol))
	sb.WriteString(fmt.Sprintf("源IP: %s | 目的IP: %s\n", h.SourceIP, h.DestIP))
	sb.WriteString("═══════════════════════════════════════════════")
	return sb.String()
}

// SubnetCalculator 子网计算器
// 对应 408 考点: 子网划分与子网掩码
type SubnetCalculator struct {
	IPAddress      string // IP 地址
	SubnetMask     string // 子网掩码
	NetworkAddress string // 网络地址
	BroadcastAddr  string // 广播地址
	FirstHost      string // 第一个可用主机地址
	LastHost       string // 最后一个可用主机地址
	TotalHosts     int    // 总主机数
	UsableHosts    int    // 可用主机数
}

// ipToUint32 将 IP 字符串转换为 uint32
func ipToUint32(ip string) uint32 {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return 0
	}
	var result uint32
	for _, part := range parts {
		var octet uint32
		fmt.Sscanf(part, "%d", &octet)
		result = (result << 8) | octet
	}
	return result
}

// uint32ToIP 将 uint32 转换为 IP 字符串
func uint32ToIP(val uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		(val>>24)&0xFF,
		(val>>16)&0xFF,
		(val>>8)&0xFF,
		val&0xFF)
}

// CalculateSubnet 计算子网信息
func CalculateSubnet(ipAddr, mask string) *SubnetCalculator {
	calc := &SubnetCalculator{
		IPAddress:  ipAddr,
		SubnetMask: mask,
	}

	ipUint := ipToUint32(ipAddr)
	maskUint := ipToUint32(mask)

	// 计算网络地址: IP AND 子网掩码
	networkUint := ipUint & maskUint
	calc.NetworkAddress = uint32ToIP(networkUint)

	// 计算广播地址: 网络地址 OR (NOT 子网掩码)
	broadcastUint := networkUint | ^maskUint
	calc.BroadcastAddr = uint32ToIP(broadcastUint)

	// 第一个可用主机地址: 网络地址 + 1
	calc.FirstHost = uint32ToIP(networkUint + 1)

	// 最后一个可用主机地址: 广播地址 - 1
	calc.LastHost = uint32ToIP(broadcastUint - 1)

	// 计算主机位数
	hostBits := 0
	temp := ^maskUint
	for temp > 0 {
		hostBits++
		temp >>= 1
	}

	// 总主机数 = 2^主机位数
	calc.TotalHosts = int(math.Pow(2, float64(hostBits)))
	// 可用主机数 = 总主机数 - 2 (减去网络地址和广播地址)
	calc.UsableHosts = calc.TotalHosts - 2

	return calc
}

// String 格式化输出子网信息
func (c *SubnetCalculator) String() string {
	var sb strings.Builder
	sb.WriteString("╔═══════════════ 子网计算结果 ═══════════════╗\n")
	sb.WriteString(fmt.Sprintf("║ IP 地址:         %s\n", c.IPAddress))
	sb.WriteString(fmt.Sprintf("║ 子网掩码:        %s\n", c.SubnetMask))
	sb.WriteString("╠═════════════════════════════════════════════╣\n")
	sb.WriteString(fmt.Sprintf("║ 网络地址:        %s\n", c.NetworkAddress))
	sb.WriteString(fmt.Sprintf("║ 广播地址:        %s\n", c.BroadcastAddr))
	sb.WriteString(fmt.Sprintf("║ 第一个主机地址:  %s\n", c.FirstHost))
	sb.WriteString(fmt.Sprintf("║ 最后一个主机地址:%s\n", c.LastHost))
	sb.WriteString("╠═════════════════════════════════════════════╣\n")
	sb.WriteString(fmt.Sprintf("║ 总主机数:        %d\n", c.TotalHosts))
	sb.WriteString(fmt.Sprintf("║ 可用主机数:      %d\n", c.UsableHosts))
	sb.WriteString("╚═════════════════════════════════════════════╝")
	return sb.String()
}

// IPFragment IP 分片结构
// 对应 408 考点: IP 分片与重组
type IPFragment struct {
	FragmentID     int    // 分片编号
	Offset         uint16 // 片偏移 (单位: 8字节)
	Length         uint16 // 数据长度
	MoreFragments  bool   // MF 标志 (More Fragments)
	Data           string // 数据内容描述
	Identification uint16 // 标识 (用于重组)
}

// String 格式化输出分片信息
func (f *IPFragment) String() string {
	mf := "0"
	if f.MoreFragments {
		mf = "1"
	}
	return fmt.Sprintf("分片#%d | 标识=%d | 偏移=%d(x8) | 长度=%d字节 | MF=%s | 数据: %s",
		f.FragmentID, f.Identification, f.Offset, f.Length, mf, f.Data)
}

// FragmentPacket 模拟 IP 分片过程
// dataSize: 数据部分大小 (字节)
// mtu: 最大传输单元 (字节)
// 返回: 分片列表
func FragmentPacket(dataSize, mtu int, identification uint16) []*IPFragment {
	const ipHeaderSize = 20 // IP 头部固定 20 字节

	// 每个分片的最大数据长度 = MTU - IP头部长度
	maxDataPerFragment := mtu - ipHeaderSize
	// 确保数据长度是 8 字节的倍数 (片偏移单位是 8 字节)
	maxDataPerFragment = (maxDataPerFragment / 8) * 8

	var fragments []*IPFragment
	remainingData := dataSize
	offset := 0
	fragmentID := 1

	for remainingData > 0 {
		// 当前分片的数据长度
		currentDataLen := maxDataPerFragment
		if remainingData < maxDataPerFragment {
			currentDataLen = remainingData
		}

		// 判断是否还有更多分片
		moreFragments := remainingData > maxDataPerFragment

		fragment := &IPFragment{
			FragmentID:     fragmentID,
			Offset:         uint16(offset / 8), // 片偏移以 8 字节为单位
			Length:         uint16(currentDataLen),
			MoreFragments:  moreFragments,
			Data:           fmt.Sprintf("数据[%d-%d]", offset, offset+currentDataLen-1),
			Identification: identification,
		}

		fragments = append(fragments, fragment)

		remainingData -= currentDataLen
		offset += currentDataLen
		fragmentID++
	}

	return fragments
}

// IPExample IP 模块示例
func IPExample() {
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("【网络层 - IP 协议示例】")
	fmt.Println(strings.Repeat("─", 50))

	// 1. IPv4 头部示例
	fmt.Println("\n1️⃣  IPv4 数据报头部结构:")
	header := NewIPv4Header("192.168.1.100", "10.0.0.5", 1500, 6)
	fmt.Println(header)

	// 2. 子网划分示例
	fmt.Println("\n2️⃣  子网划分与计算:")
	fmt.Println("\n示例 1: C 类网络")
	subnet1 := CalculateSubnet("192.168.1.100", "255.255.255.0")
	fmt.Println(subnet1)

	fmt.Println("\n示例 2: 子网划分 (/26)")
	subnet2 := CalculateSubnet("172.16.10.50", "255.255.255.192")
	fmt.Println(subnet2)

	// 3. IP 分片示例
	fmt.Println("\n3️⃣  IP 分片模拟:")
	fmt.Println("\n场景: 3800 字节数据报,经过 MTU=1500 字节的链路")
	fragments := FragmentPacket(3800, 1500, 54321)
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	for _, frag := range fragments {
		fmt.Printf("║ %s\n", frag)
	}
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	// 408 考点提示
	fmt.Println("\n📚 408 考点总结:")
	fmt.Println("  ✓ IP 数据报格式 (20 字节固定头部)")
	fmt.Println("  ✓ 子网划分: 网络地址 = IP AND 掩码")
	fmt.Println("  ✓ IP 分片: 片偏移以 8 字节为单位, MF 标志")
	fmt.Println("  ✓ 分片数据长度必须是 8 字节的倍数 (最后一片除外)")
}
