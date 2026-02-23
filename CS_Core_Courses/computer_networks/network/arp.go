package network

import (
	"fmt"
	"strings"
	"time"
)

// ARPEntry ARP 表项
// 对应 408 考点: ARP 缓存
type ARPEntry struct {
	IPAddress  string    // IP 地址
	MACAddress string    // MAC 地址
	Timestamp  time.Time // 时间戳
	TTL        int       // 生存时间 (秒)
}

// IsExpired 检查表项是否过期
func (e *ARPEntry) IsExpired() bool {
	return time.Since(e.Timestamp).Seconds() > float64(e.TTL)
}

// String 格式化输出 ARP 表项
func (e *ARPEntry) String() string {
	age := int(time.Since(e.Timestamp).Seconds())
	status := "有效"
	if e.IsExpired() {
		status = "过期"
	}
	return fmt.Sprintf("IP: %-15s | MAC: %-17s | 年龄: %3ds/%3ds | 状态: %s",
		e.IPAddress, e.MACAddress, age, e.TTL, status)
}

// ARPTable ARP 表 (缓存)
// 对应 408 考点: ARP 地址解析协议
type ARPTable struct {
	Entries map[string]*ARPEntry // IP -> ARP 表项映射
}

// NewARPTable 创建 ARP 表
func NewARPTable() *ARPTable {
	return &ARPTable{
		Entries: make(map[string]*ARPEntry),
	}
}

// Add 添加 ARP 表项
func (t *ARPTable) Add(ip, mac string, ttl int) {
	t.Entries[ip] = &ARPEntry{
		IPAddress:  ip,
		MACAddress: mac,
		Timestamp:  time.Now(),
		TTL:        ttl,
	}
}

// Lookup 查找 MAC 地址
func (t *ARPTable) Lookup(ip string) (string, bool) {
	entry, exists := t.Entries[ip]
	if !exists {
		return "", false
	}

	// 检查是否过期
	if entry.IsExpired() {
		delete(t.Entries, ip) // 删除过期表项
		return "", false
	}

	return entry.MACAddress, true
}

// CleanExpired 清除过期表项
func (t *ARPTable) CleanExpired() int {
	count := 0
	for ip, entry := range t.Entries {
		if entry.IsExpired() {
			delete(t.Entries, ip)
			count++
		}
	}
	return count
}

// PrintTable 打印 ARP 表
func (t *ARPTable) PrintTable() {
	fmt.Println("\n【ARP 缓存表】")
	if len(t.Entries) == 0 {
		fmt.Println("  (空)")
		return
	}

	fmt.Println("┌─────────────────┬───────────────────┬──────────────┬────────┐")
	fmt.Println("│  IP 地址        │  MAC 地址         │   年龄/TTL   │  状态  │")
	fmt.Println("├─────────────────┼───────────────────┼──────────────┼────────┤")
	for _, entry := range t.Entries {
		age := int(time.Since(entry.Timestamp).Seconds())
		status := "✓"
		if entry.IsExpired() {
			status = "✗ 过期"
		}
		fmt.Printf("│ %-15s │ %-17s │ %4ds / %4ds │ %-6s │\n",
			entry.IPAddress, entry.MACAddress, age, entry.TTL, status)
	}
	fmt.Println("└─────────────────┴───────────────────┴──────────────┴────────┘")
}

// ARPMessage ARP 报文
// 对应 408 考点: ARP 请求与应答
type ARPMessage struct {
	OpCode         string // 操作码: REQUEST 或 REPLY
	SenderIP       string // 发送方 IP
	SenderMAC      string // 发送方 MAC
	TargetIP       string // 目标 IP
	TargetMAC      string // 目标 MAC (请求时为全 0)
}

// NewARPRequest 创建 ARP 请求
func NewARPRequest(senderIP, senderMAC, targetIP string) *ARPMessage {
	return &ARPMessage{
		OpCode:    "REQUEST",
		SenderIP:  senderIP,
		SenderMAC: senderMAC,
		TargetIP:  targetIP,
		TargetMAC: "00:00:00:00:00:00", // 未知,填充 0
	}
}

// NewARPReply 创建 ARP 应答
func NewARPReply(senderIP, senderMAC, targetIP, targetMAC string) *ARPMessage {
	return &ARPMessage{
		OpCode:    "REPLY",
		SenderIP:  senderIP,
		SenderMAC: senderMAC,
		TargetIP:  targetIP,
		TargetMAC: targetMAC,
	}
}

// String 格式化输出 ARP 报文
func (m *ARPMessage) String() string {
	var sb strings.Builder
	sb.WriteString("╔══════════════════ ARP 报文 ══════════════════╗\n")
	sb.WriteString(fmt.Sprintf("║ 操作码:     %-28s║\n", m.OpCode))
	sb.WriteString("╟──────────────────────────────────────────────╢\n")
	sb.WriteString(fmt.Sprintf("║ 发送方 IP:  %-28s║\n", m.SenderIP))
	sb.WriteString(fmt.Sprintf("║ 发送方 MAC: %-28s║\n", m.SenderMAC))
	sb.WriteString(fmt.Sprintf("║ 目标 IP:    %-28s║\n", m.TargetIP))
	sb.WriteString(fmt.Sprintf("║ 目标 MAC:   %-28s║\n", m.TargetMAC))
	sb.WriteString("╚══════════════════════════════════════════════╝")
	return sb.String()
}

// Host 主机 (用于 ARP 模拟)
type Host struct {
	Name       string     // 主机名
	IPAddress  string     // IP 地址
	MACAddress string     // MAC 地址
	ARPTable   *ARPTable  // ARP 表
}

// NewHost 创建主机
func NewHost(name, ip, mac string) *Host {
	return &Host{
		Name:       name,
		IPAddress:  ip,
		MACAddress: mac,
		ARPTable:   NewARPTable(),
	}
}

// SendARPRequest 发送 ARP 请求 (广播)
func (h *Host) SendARPRequest(targetIP string) *ARPMessage {
	fmt.Printf("\n[%s] 发送 ARP 请求 (广播): 谁是 %s? 请告诉 %s\n",
		h.Name, targetIP, h.IPAddress)
	return NewARPRequest(h.IPAddress, h.MACAddress, targetIP)
}

// ReceiveARPRequest 接收 ARP 请求
func (h *Host) ReceiveARPRequest(msg *ARPMessage) *ARPMessage {
	// 检查目标 IP 是否是自己
	if msg.TargetIP != h.IPAddress {
		return nil // 不是发给自己的,忽略
	}

	fmt.Printf("\n[%s] 收到 ARP 请求: %s (%s) 询问 %s 的 MAC 地址\n",
		h.Name, msg.SenderIP, msg.SenderMAC, msg.TargetIP)

	// 更新 ARP 表 (学习发送方的 IP-MAC 映射)
	h.ARPTable.Add(msg.SenderIP, msg.SenderMAC, 120)
	fmt.Printf("[%s] 更新 ARP 表: %s -> %s\n", h.Name, msg.SenderIP, msg.SenderMAC)

	// 发送 ARP 应答 (单播)
	fmt.Printf("[%s] 发送 ARP 应答 (单播): %s 的 MAC 地址是 %s\n",
		h.Name, h.IPAddress, h.MACAddress)
	return NewARPReply(h.IPAddress, h.MACAddress, msg.SenderIP, msg.SenderMAC)
}

// ReceiveARPReply 接收 ARP 应答
func (h *Host) ReceiveARPReply(msg *ARPMessage) {
	fmt.Printf("\n[%s] 收到 ARP 应答: %s 的 MAC 地址是 %s\n",
		h.Name, msg.SenderIP, msg.SenderMAC)

	// 更新 ARP 表
	h.ARPTable.Add(msg.SenderIP, msg.SenderMAC, 120)
	fmt.Printf("[%s] 更新 ARP 表: %s -> %s\n", h.Name, msg.SenderIP, msg.SenderMAC)
}

// ARPExample ARP 协议示例
func ARPExample() {
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("【网络层 - ARP 协议示例】")
	fmt.Println(strings.Repeat("─", 50))

	// 创建两个主机
	hostA := NewHost("主机A", "192.168.1.10", "AA:BB:CC:DD:EE:01")
	hostB := NewHost("主机B", "192.168.1.20", "AA:BB:CC:DD:EE:02")

	fmt.Println("\n初始状态:")
	fmt.Printf("主机A: IP=%s, MAC=%s\n", hostA.IPAddress, hostA.MACAddress)
	fmt.Printf("主机B: IP=%s, MAC=%s\n", hostB.IPAddress, hostB.MACAddress)

	// 场景: 主机 A 想要发送数据给主机 B,但不知道 B 的 MAC 地址
	fmt.Println("\n" + strings.Repeat("═", 50))
	fmt.Println("场景: 主机A 要发送数据给 192.168.1.20,但不知道其 MAC 地址")
	fmt.Println(strings.Repeat("═", 50))

	// 1. 主机 A 查找 ARP 表
	fmt.Println("\n步骤 1: 主机A 查找 ARP 表")
	targetIP := "192.168.1.20"
	mac, found := hostA.ARPTable.Lookup(targetIP)
	if found {
		fmt.Printf("✓ 找到缓存: %s -> %s\n", targetIP, mac)
	} else {
		fmt.Printf("✗ 未找到缓存,需要发送 ARP 请求\n")
	}

	// 2. 主机 A 发送 ARP 请求 (广播)
	fmt.Println("\n步骤 2: 主机A 发送 ARP 请求 (广播)")
	request := hostA.SendARPRequest(targetIP)
	fmt.Println(request)

	// 3. 主机 B 接收 ARP 请求并回复
	fmt.Println("\n步骤 3: 主机B 接收 ARP 请求")
	reply := hostB.ReceiveARPRequest(request)
	if reply != nil {
		fmt.Println(reply)
	}

	// 4. 主机 A 接收 ARP 应答
	fmt.Println("\n步骤 4: 主机A 接收 ARP 应答")
	hostA.ReceiveARPReply(reply)

	// 5. 显示最终的 ARP 表
	fmt.Println("\n步骤 5: 查看最终的 ARP 表")
	fmt.Println("\n主机A 的 ARP 表:")
	hostA.ARPTable.PrintTable()

	fmt.Println("\n主机B 的 ARP 表:")
	hostB.ARPTable.PrintTable()

	// 6. 再次通信,命中缓存
	fmt.Println("\n" + strings.Repeat("═", 50))
	fmt.Println("场景 2: 主机A 再次发送数据给主机B (命中缓存)")
	fmt.Println(strings.Repeat("═", 50))
	mac, found = hostA.ARPTable.Lookup(targetIP)
	if found {
		fmt.Printf("✓ 命中 ARP 缓存: %s -> %s\n", targetIP, mac)
		fmt.Println("✓ 直接封装以太网帧发送,无需 ARP 请求")
	}

	// 408 考点提示
	fmt.Println("\n📚 408 考点总结:")
	fmt.Println("  ✓ ARP 功能: IP 地址 → MAC 地址 (网络层 → 数据链路层)")
	fmt.Println("  ✓ ARP 请求: 广播方式,目标 MAC 为全 F (FF:FF:FF:FF:FF:FF)")
	fmt.Println("  ✓ ARP 应答: 单播方式,直接回复给请求方")
	fmt.Println("  ✓ ARP 缓存: 减少网络流量,有超时机制")
	fmt.Println("  ✓ 工作层次: 网络层协议,使用数据链路层服务")
}
