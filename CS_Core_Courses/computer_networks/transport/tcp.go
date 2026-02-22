package transport

import (
	"fmt"
	"math/rand"
	"time"
)

// TCPState TCP状态
type TCPState int

const (
	TCPClosed TCPState = iota
	TCPListen
	TCPSynSent
	TCPSynReceived
	TCPEstablished
	TCPFinWait1
	TCPFinWait2
	TCPClosing
	TCPTimeWait
	TCPCloseWait
	TCPLastAck
)

// String 返回状态的字符串表示
func (state TCPState) String() string {
	switch state {
	case TCPListen:
		return "LISTEN"
	case TCPSynSent:
		return "SYN_SENT"
	case TCPSynReceived:
		return "SYN_RECEIVED"
	case TCPEstablished:
		return "ESTABLISHED"
	case TCPFinWait1:
		return "FIN_WAIT_1"
	case TCPFinWait2:
		return "FIN_WAIT_2"
	case TCPClosing:
		return "CLOSING"
	case TCPTimeWait:
		return "TIME_WAIT"
	case TCPCloseWait:
		return "CLOSE_WAIT"
	case TCPLastAck:
		return "LAST_ACK"
	case TCPClosed:
		return "CLOSED"
	default:
		return "UNKNOWN"
	}
}

// TCPFlag TCP标志位
type TCPFlag uint8

const (
	FlagFIN TCPFlag = 1 << 0
	FlagSYN TCPFlag = 1 << 1
	FlagRST TCPFlag = 1 << 2
	FlagPSH TCPFlag = 1 << 3
	FlagACK TCPFlag = 1 << 4
	FlagURG TCPFlag = 1 << 5
)

// TCPPacket TCP数据包
type TCPPacket struct {
	SourcePort     uint16
	DestPort       uint16
	SequenceNumber uint32
	AckNumber      uint32
	DataOffset     uint8 // 4位
	Reserved       uint8 // 3位
	Flags          TCPFlag
	WindowSize     uint16
	Checksum       uint16
	UrgentPointer  uint16
	Options        []byte
	Data           []byte
}

// NewTCPPacket 创建TCP包
func NewTCPPacket(srcPort, dstPort uint16, flags TCPFlag, data []byte) *TCPPacket {
	// 生成随机序列号
	rand.Seed(time.Now().UnixNano())
	seqNum := rand.Uint32()

	return &TCPPacket{
		SourcePort:     srcPort,
		DestPort:       dstPort,
		SequenceNumber: seqNum,
		AckNumber:      0,
		DataOffset:     20, // 基本头部20字节，无选项
		Reserved:       0,
		Flags:          flags,
		WindowSize:     65535, // 最大窗口大小
		Checksum:       0,     // 实际应用中需要计算
		UrgentPointer:  0,
		Options:        nil,
		Data:           data,
	}
}

// HasFlag 检查是否有指定标志位
func (p *TCPPacket) HasFlag(flag TCPFlag) bool {
	return p.Flags&flag != 0
}

// String 返回TCP包的字符串表示
func (p *TCPPacket) String() string {
	flags := ""
	if p.HasFlag(FlagFIN) {
		flags += "FIN "
	}
	if p.HasFlag(FlagSYN) {
		flags += "SYN "
	}
	if p.HasFlag(FlagRST) {
		flags += "RST "
	}
	if p.HasFlag(FlagPSH) {
		flags += "PSH "
	}
	if p.HasFlag(FlagACK) {
		flags += "ACK "
	}
	if p.HasFlag(FlagURG) {
		flags += "URG "
	}

	return fmt.Sprintf("TCP %d->%d Seq=%d Ack=%d Flags=[%s] Len=%d",
		p.SourcePort, p.DestPort, p.SequenceNumber, p.AckNumber, flags, len(p.Data))
}

// TCPConnection TCP连接
type TCPConnection struct {
	LocalPort       uint16
	RemotePort      uint16
	State           TCPState
	LocalSeqNum     uint32
	RemoteSeqNum    uint32
	WindowSize      uint16
	ReceiveBuffer   []byte
	SendBuffer      []byte
	UnackedBytes    int
	LastAckSent     uint32
	LastAckReceived uint32
}

// NewTCPConnection 创建TCP连接
func NewTCPConnection(localPort, remotePort uint16) *TCPConnection {
	rand.Seed(time.Now().UnixNano())
	return &TCPConnection{
		LocalPort:       localPort,
		RemotePort:      remotePort,
		State:           TCPClosed,
		LocalSeqNum:     rand.Uint32(),
		RemoteSeqNum:    0,
		WindowSize:      65535,
		ReceiveBuffer:   make([]byte, 0, 65535),
		SendBuffer:      make([]byte, 0, 65535),
		UnackedBytes:    0,
		LastAckSent:     0,
		LastAckReceived: 0,
	}
}

// SetState 设置连接状态
func (conn *TCPConnection) SetState(state TCPState) {
	conn.State = state
}

// HandlePacket 处理接收到的TCP包
func (conn *TCPConnection) HandlePacket(packet *TCPPacket) []TCPPacket {
	var responses []TCPPacket

	switch conn.State {
	case TCPClosed:
		if packet.HasFlag(FlagSYN) {
			// 收到SYN，发送SYN+ACK，进入SYN_RECEIVED状态
			conn.RemoteSeqNum = packet.SequenceNumber + 1
			synAck := NewTCPPacket(conn.LocalPort, conn.RemotePort, FlagSYN|FlagACK, nil)
			synAck.SequenceNumber = conn.LocalSeqNum
			synAck.AckNumber = conn.RemoteSeqNum
			responses = append(responses, *synAck)
			conn.SetState(TCPSynReceived)
		}

	case TCPSynSent:
		if packet.HasFlag(FlagSYN) && packet.HasFlag(FlagACK) {
			// 收到SYN+ACK，发送ACK，进入ESTABLISHED状态
			if packet.AckNumber == conn.LocalSeqNum+1 {
				conn.RemoteSeqNum = packet.SequenceNumber + 1
				conn.LocalSeqNum++
				ack := NewTCPPacket(conn.LocalPort, conn.RemotePort, FlagACK, nil)
				ack.SequenceNumber = conn.LocalSeqNum
				ack.AckNumber = conn.RemoteSeqNum
				responses = append(responses, *ack)
				conn.SetState(TCPEstablished)
			}
		}

	case TCPSynReceived:
		if packet.HasFlag(FlagACK) {
			// 收到ACK，进入ESTABLISHED状态
			if packet.AckNumber == conn.LocalSeqNum+1 {
				conn.LocalSeqNum++
				conn.SetState(TCPEstablished)
			}
		}

	case TCPEstablished:
		if packet.HasFlag(FlagFIN) {
			// 收到FIN，发送ACK，进入CLOSE_WAIT状态
			conn.RemoteSeqNum++
			ack := NewTCPPacket(conn.LocalPort, conn.RemotePort, FlagACK, nil)
			ack.SequenceNumber = conn.LocalSeqNum
			ack.AckNumber = conn.RemoteSeqNum
			responses = append(responses, *ack)
			conn.SetState(TCPCloseWait)
		} else if packet.HasFlag(FlagACK) {
			// 处理ACK
			if packet.AckNumber > conn.LastAckReceived {
				conn.LastAckReceived = packet.AckNumber
				conn.UnackedBytes = int(conn.LocalSeqNum - packet.AckNumber)
			}
		} else if len(packet.Data) > 0 {
			// 处理数据
			conn.ReceiveBuffer = append(conn.ReceiveBuffer, packet.Data...)
			conn.RemoteSeqNum += uint32(len(packet.Data))

			// 发送ACK
			ack := NewTCPPacket(conn.LocalPort, conn.RemotePort, FlagACK, nil)
			ack.SequenceNumber = conn.LocalSeqNum
			ack.AckNumber = conn.RemoteSeqNum
			responses = append(responses, *ack)
		}

	case TCPFinWait1:
		if packet.HasFlag(FlagACK) {
			if packet.AckNumber == conn.LocalSeqNum+1 {
				conn.LocalSeqNum++
				conn.SetState(TCPFinWait2)
			}
		}
		if packet.HasFlag(FlagFIN) {
			conn.RemoteSeqNum++
			ack := NewTCPPacket(conn.LocalPort, conn.RemotePort, FlagACK, nil)
			ack.SequenceNumber = conn.LocalSeqNum
			ack.AckNumber = conn.RemoteSeqNum
			responses = append(responses, *ack)
			conn.SetState(TCPClosing)
		}

	case TCPFinWait2:
		if packet.HasFlag(FlagFIN) {
			conn.RemoteSeqNum++
			ack := NewTCPPacket(conn.LocalPort, conn.RemotePort, FlagACK, nil)
			ack.SequenceNumber = conn.LocalSeqNum
			ack.AckNumber = conn.RemoteSeqNum
			responses = append(responses, *ack)
			conn.SetState(TCPTimeWait)
		}

	case TCPCloseWait:
		// 应用层调用Close()后发送FIN
		// 这里由外部触发

	case TCPLastAck:
		if packet.HasFlag(FlagACK) {
			if packet.AckNumber == conn.LocalSeqNum+1 {
				conn.SetState(TCPClosed)
			}
		}

	case TCPClosing:
		if packet.HasFlag(FlagACK) {
			if packet.AckNumber == conn.LocalSeqNum+1 {
				conn.SetState(TCPTimeWait)
			}
		}

	case TCPTimeWait:
		// 等待2MSL后关闭连接
		// 这里简化处理
		conn.SetState(TCPClosed)
	}

	return responses
}

// Connect 主动连接
func (conn *TCPConnection) Connect() TCPPacket {
	// 发送SYN
	syn := NewTCPPacket(conn.LocalPort, conn.RemotePort, FlagSYN, nil)
	syn.SequenceNumber = conn.LocalSeqNum
	conn.SetState(TCPSynSent)
	return *syn
}

// Close 关闭连接
func (conn *TCPConnection) Close() []TCPPacket {
	var responses []TCPPacket

	switch conn.State {
	case TCPEstablished:
		// 发送FIN
		fin := NewTCPPacket(conn.LocalPort, conn.RemotePort, FlagFIN, nil)
		fin.SequenceNumber = conn.LocalSeqNum
		conn.LocalSeqNum++
		responses = append(responses, *fin)
		conn.SetState(TCPFinWait1)

	case TCPCloseWait:
		// 发送FIN
		fin := NewTCPPacket(conn.LocalPort, conn.RemotePort, FlagFIN, nil)
		fin.SequenceNumber = conn.LocalSeqNum
		conn.LocalSeqNum++
		responses = append(responses, *fin)
		conn.SetState(TCPLastAck)
	}

	return responses
}

// SendData 发送数据
func (conn *TCPConnection) SendData(data []byte) TCPPacket {
	if conn.State != TCPEstablished {
		return TCPPacket{}
	}

	packet := NewTCPPacket(conn.LocalPort, conn.RemotePort, FlagPSH|FlagACK, data)
	packet.SequenceNumber = conn.LocalSeqNum
	packet.AckNumber = conn.RemoteSeqNum
	conn.LocalSeqNum += uint32(len(data))
	conn.UnackedBytes += len(data)

	return *packet
}

// String 返回连接状态的字符串表示
func (conn *TCPConnection) String() string {
	return fmt.Sprintf("TCP Connection %d<->%d State: %s LocalSeq: %d RemoteSeq: %d Unacked: %d",
		conn.LocalPort, conn.RemotePort, conn.State, conn.LocalSeqNum, conn.RemoteSeqNum, conn.UnackedBytes)
}

// TCPServer 简单TCP服务器
type TCPServer struct {
	ListenPort  uint16
	Connections map[string]*TCPConnection
}

// NewTCPServer 创建TCP服务器
func NewTCPServer(port uint16) *TCPServer {
	return &TCPServer{
		ListenPort:  port,
		Connections: make(map[string]*TCPConnection),
	}
}

// Listen 监听连接
func (server *TCPServer) Listen() {
	fmt.Printf("TCP Server listening on port %d\n", server.ListenPort)
	// 在实际实现中，这里会开始监听网络接口
}

// HandleConnection 处理连接
func (server *TCPServer) HandleConnection(packet *TCPPacket) []TCPPacket {
	key := fmt.Sprintf("%d:%d", packet.SourcePort, packet.DestPort)
	conn, exists := server.Connections[key]

	if !exists {
		if packet.HasFlag(FlagSYN) {
			// 新连接
			conn = NewTCPConnection(server.ListenPort, packet.SourcePort)
			server.Connections[key] = conn
		} else {
			return []TCPPacket{}
		}
	}

	return conn.HandlePacket(packet)
}

// 示例函数
func TCPExample() {
	fmt.Println("=== TCP 协议示例 ===")

	fmt.Println("1. TCP 连接建立 (三次握手):")

	// 客户端
	client := NewTCPConnection(12345, 80)
	fmt.Printf("客户端初始状态: %s\n", client.State)

	// 服务器
	server := NewTCPServer(80)
	server.Listen()

	// 1. 客户端发送SYN
	synPacket := client.Connect()
	fmt.Printf("1. 客户端发送: %s\n", synPacket.String())
	fmt.Printf("   客户端状态: %s\n", client.State)

	// 2. 服务器处理SYN，发送SYN+ACK
	serverResponses := server.HandleConnection(&synPacket)
	for _, resp := range serverResponses {
		fmt.Printf("2. 服务器发送: %s\n", resp.String())
	}

	// 3. 客户端处理SYN+ACK，发送ACK
	if len(serverResponses) > 0 {
		clientResponses := client.HandlePacket(&serverResponses[0])
		for _, resp := range clientResponses {
			fmt.Printf("3. 客户端发送: %s\n", resp.String())
		}
		fmt.Printf("   客户端状态: %s\n", client.State)
	}

	fmt.Println("\n2. TCP 数据传输:")
	// 客户端发送数据
	if client.State == TCPEstablished {
		data := []byte("Hello, TCP Server!")
		dataPacket := client.SendData(data)
		fmt.Printf("客户端发送数据: %s\n", dataPacket.String())

		// 服务器处理数据包
		serverDataResponses := server.HandleConnection(&dataPacket)
		for _, resp := range serverDataResponses {
			fmt.Printf("服务器响应: %s\n", resp.String())
		}
	}

	fmt.Println("\n3. TCP 连接断开 (四次挥手):")
	// 客户端主动关闭
	closePackets := client.Close()
	for _, packet := range closePackets {
		fmt.Printf("客户端发送: %s\n", packet.String())
		fmt.Printf("   客户端状态: %s\n", client.State)
	}

	// 服务器处理FIN
	if len(closePackets) > 0 {
		serverCloseResponses := server.HandleConnection(&closePackets[0])
		for _, resp := range serverCloseResponses {
			fmt.Printf("服务器发送: %s\n", resp.String())
		}

		// 客户端处理服务器的ACK
		if len(serverCloseResponses) > 0 {
			client.HandlePacket(&serverCloseResponses[0])
			fmt.Printf("客户端状态: %s\n", client.State)

			// 服务器发送FIN
			key := fmt.Sprintf("%d:%d", client.LocalPort, client.RemotePort)
			if conn, exists := server.Connections[key]; exists {
				serverFinPackets := conn.Close()
				for _, finPacket := range serverFinPackets {
					fmt.Printf("服务器发送: %s\n", finPacket.String())
				}

				// 客户端处理FIN
				if len(serverFinPackets) > 0 {
					clientFinResponses := client.HandlePacket(&serverFinPackets[0])
					for _, resp := range clientFinResponses {
						fmt.Printf("客户端发送: %s\n", resp.String())
						fmt.Printf("   客户端状态: %s\n", client.State)
					}
				}
			}
		}
	}

	fmt.Println("\n4. TCP 标志位说明:")
	flags := map[TCPFlag]string{
		FlagFIN: "FIN - 连接终止",
		FlagSYN: "SYN - 同步序列号",
		FlagRST: "RST - 重置连接",
		FlagPSH: "PSH - 推送功能",
		FlagACK: "ACK - 确认应答",
		FlagURG: "URG - 紧急指针",
	}

	for flag, desc := range flags {
		fmt.Printf("  0x%02X: %s\n", flag, desc)
	}

	fmt.Println("\n5. TCP 状态转换:")
	states := []TCPState{
		TCPClosed, TCPSynSent, TCPSynReceived,
		TCPEstablished, TCPFinWait1, TCPFinWait2,
		TCPClosing, TCPTimeWait, TCPCloseWait,
		TCPLastAck,
	}

	for _, state := range states {
		fmt.Printf("  %d: %s\n", state, state.String())
	}

	fmt.Println()
}
