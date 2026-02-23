package datalink

import (
	"fmt"
	"strings"
)

// Frame 数据帧
type Frame struct {
	SeqNum int    // 序号
	Data   string // 数据
	ACK    int    // 确认号 (用于捎带确认)
}

// String 格式化输出帧
func (f *Frame) String() string {
	return fmt.Sprintf("帧[Seq=%d, Data=%s]", f.SeqNum, f.Data)
}

// StopAndWait 停止-等待协议
// 对应 408 考点: 最简单的滑动窗口协议 (窗口大小 = 1)
type StopAndWait struct {
	SenderSeq   int  // 发送方序号
	ReceiverSeq int  // 接收方期望序号
	Timeout     bool // 超时标志
}

// NewStopAndWait 创建停止-等待协议
func NewStopAndWait() *StopAndWait {
	return &StopAndWait{
		SenderSeq:   0,
		ReceiverSeq: 0,
		Timeout:     false,
	}
}

// Send 发送帧
func (sw *StopAndWait) Send(data string) {
	frame := &Frame{SeqNum: sw.SenderSeq, Data: data}
	fmt.Printf("发送方: 发送 %s\n", frame.String())
}

// Receive 接收帧
func (sw *StopAndWait) Receive(frame *Frame) bool {
	if frame.SeqNum == sw.ReceiverSeq {
		fmt.Printf("接收方: 正确接收 %s, 发送 ACK %d\n", frame, sw.ReceiverSeq)
		sw.ReceiverSeq = 1 - sw.ReceiverSeq // 0/1 交替
		return true
	}
	fmt.Printf("接收方: 收到重复帧 %s, 丢弃, 重发 ACK %d\n", frame, 1-sw.ReceiverSeq)
	return false
}

// ACK 确认
func (sw *StopAndWait) ACK() {
	fmt.Printf("发送方: 收到 ACK %d\n", sw.SenderSeq)
	sw.SenderSeq = 1 - sw.SenderSeq // 0/1 交替
}

// GoBackN 回退 N 帧协议
// 对应 408 考点: GBN 协议,发送窗口 > 1,接收窗口 = 1
type GoBackN struct {
	WindowSize  int      // 窗口大小
	SeqNumBits  int      // 序号位数
	MaxSeqNum   int      // 最大序号 (2^n - 1)
	SendBase    int      // 发送窗口基序号
	NextSeqNum  int      // 下一个待发送序号
	ExpectedSeq int      // 接收方期望序号
	SentFrames  []*Frame // 已发送但未确认的帧
}

// NewGoBackN 创建 GBN 协议
func NewGoBackN(windowSize, seqNumBits int) *GoBackN {
	maxSeqNum := (1 << seqNumBits) - 1 // 2^n - 1
	// 408 考点: GBN 窗口大小 <= 2^n - 1
	if windowSize > maxSeqNum {
		fmt.Printf("警告: 窗口大小 %d 超过最大值 %d, 自动调整\n", windowSize, maxSeqNum)
		windowSize = maxSeqNum
	}

	return &GoBackN{
		WindowSize:  windowSize,
		SeqNumBits:  seqNumBits,
		MaxSeqNum:   maxSeqNum,
		SendBase:    0,
		NextSeqNum:  0,
		ExpectedSeq: 0,
		SentFrames:  make([]*Frame, 0),
	}
}

// CanSend 是否可以发送
func (gbn *GoBackN) CanSend() bool {
	return (gbn.NextSeqNum - gbn.SendBase) < gbn.WindowSize
}

// Send 发送帧
func (gbn *GoBackN) Send(data string) bool {
	if !gbn.CanSend() {
		fmt.Printf("发送方: 窗口已满 [%d, %d), 无法发送\n",
			gbn.SendBase, gbn.SendBase+gbn.WindowSize)
		return false
	}

	frame := &Frame{SeqNum: gbn.NextSeqNum % (gbn.MaxSeqNum + 1), Data: data}
	gbn.SentFrames = append(gbn.SentFrames, frame)
	fmt.Printf("发送方: 发送 %s, 窗口 [%d, %d)\n",
		frame, gbn.SendBase, gbn.SendBase+gbn.WindowSize)
	gbn.NextSeqNum++
	return true
}

// Receive 接收帧
func (gbn *GoBackN) Receive(frame *Frame) int {
	expectedSeq := gbn.ExpectedSeq % (gbn.MaxSeqNum + 1)
	if frame.SeqNum == expectedSeq {
		fmt.Printf("接收方: 正确接收 %s, 发送 ACK %d\n", frame, expectedSeq)
		gbn.ExpectedSeq++
		return expectedSeq
	}
	// 收到失序帧,丢弃
	lastACK := (gbn.ExpectedSeq - 1) % (gbn.MaxSeqNum + 1)
	if lastACK < 0 {
		lastACK = gbn.MaxSeqNum
	}
	fmt.Printf("接收方: 收到失序帧 %s (期望 %d), 丢弃, 重发 ACK %d\n",
		frame, expectedSeq, lastACK)
	return lastACK
}

// ACK 确认
func (gbn *GoBackN) ACK(ackNum int) {
	// 累积确认: ACK n 表示 n 及之前的所有帧都正确接收
	fmt.Printf("发送方: 收到 ACK %d (累积确认)\n", ackNum)
	// 更新窗口基序号
	ackedCount := 0
	for i := gbn.SendBase; i <= ackNum; i++ {
		ackedCount++
	}
	gbn.SendBase = ackNum + 1
	// 移除已确认的帧
	if ackedCount > 0 && ackedCount <= len(gbn.SentFrames) {
		gbn.SentFrames = gbn.SentFrames[ackedCount:]
	}
	fmt.Printf("发送方: 窗口前移到 [%d, %d)\n", gbn.SendBase, gbn.SendBase+gbn.WindowSize)
}

// Timeout 超时,重传所有已发送但未确认的帧
func (gbn *GoBackN) Timeout() {
	fmt.Printf("发送方: 超时! 回退重传从 %d 开始的所有帧\n", gbn.SendBase)
	for _, frame := range gbn.SentFrames {
		fmt.Printf("发送方: 重传 %s\n", frame)
	}
	gbn.NextSeqNum = gbn.SendBase + len(gbn.SentFrames)
}

// SelectiveRepeat 选择重传协议
// 对应 408 考点: SR 协议,发送窗口 = 接收窗口,窗口大小 <= 2^(n-1)
type SelectiveRepeat struct {
	WindowSize   int            // 窗口大小
	SeqNumBits   int            // 序号位数
	MaxSeqNum    int            // 最大序号
	SendBase     int            // 发送窗口基序号
	NextSeqNum   int            // 下一个待发送序号
	RecvBase     int            // 接收窗口基序号
	RecvBuffer   map[int]*Frame // 接收缓存 (序号 -> 帧)
	ACKedFrames  map[int]bool   // 已确认的帧
}

// NewSelectiveRepeat 创建 SR 协议
func NewSelectiveRepeat(windowSize, seqNumBits int) *SelectiveRepeat {
	maxSeqNum := (1 << seqNumBits) - 1
	maxWindowSize := 1 << (seqNumBits - 1) // 2^(n-1)
	// 408 考点: SR 窗口大小 <= 2^(n-1)
	if windowSize > maxWindowSize {
		fmt.Printf("警告: SR 窗口大小 %d 超过最大值 %d, 自动调整\n", windowSize, maxWindowSize)
		windowSize = maxWindowSize
	}

	return &SelectiveRepeat{
		WindowSize:  windowSize,
		SeqNumBits:  seqNumBits,
		MaxSeqNum:   maxSeqNum,
		SendBase:    0,
		NextSeqNum:  0,
		RecvBase:    0,
		RecvBuffer:  make(map[int]*Frame),
		ACKedFrames: make(map[int]bool),
	}
}

// CanSend 是否可以发送
func (sr *SelectiveRepeat) CanSend() bool {
	return (sr.NextSeqNum - sr.SendBase) < sr.WindowSize
}

// Send 发送帧
func (sr *SelectiveRepeat) Send(data string) bool {
	if !sr.CanSend() {
		fmt.Printf("发送方: 窗口已满 [%d, %d), 无法发送\n",
			sr.SendBase, sr.SendBase+sr.WindowSize)
		return false
	}

	seqNum := sr.NextSeqNum % (sr.MaxSeqNum + 1)
	frame := &Frame{SeqNum: seqNum, Data: data}
	fmt.Printf("发送方: 发送 %s, 窗口 [%d, %d)\n",
		frame, sr.SendBase, sr.SendBase+sr.WindowSize)
	sr.NextSeqNum++
	return true
}

// Receive 接收帧
func (sr *SelectiveRepeat) Receive(frame *Frame) {
	// 检查是否在接收窗口内
	inWindow := false
	for i := 0; i < sr.WindowSize; i++ {
		expectedSeq := (sr.RecvBase + i) % (sr.MaxSeqNum + 1)
		if frame.SeqNum == expectedSeq {
			inWindow = true
			break
		}
	}

	if !inWindow {
		fmt.Printf("接收方: %s 不在窗口 [%d, %d) 内, 丢弃\n",
			frame, sr.RecvBase, sr.RecvBase+sr.WindowSize)
		return
	}

	// 缓存帧并发送 ACK
	sr.RecvBuffer[frame.SeqNum] = frame
	fmt.Printf("接收方: 接收 %s, 缓存, 发送 ACK %d\n", frame, frame.SeqNum)

	// 如果是窗口基序号,交付并滑动窗口
	if frame.SeqNum == sr.RecvBase%(sr.MaxSeqNum+1) {
		fmt.Printf("接收方: 交付数据并滑动窗口\n")
		for {
			currentSeq := sr.RecvBase % (sr.MaxSeqNum + 1)
			if _, exists := sr.RecvBuffer[currentSeq]; !exists {
				break
			}
			fmt.Printf("接收方: 交付 帧[Seq=%d]\n", currentSeq)
			delete(sr.RecvBuffer, currentSeq)
			sr.RecvBase++
		}
		fmt.Printf("接收方: 窗口前移到 [%d, %d)\n", sr.RecvBase, sr.RecvBase+sr.WindowSize)
	}
}

// ACK 确认 (单独确认,非累积)
func (sr *SelectiveRepeat) ACK(ackNum int) {
	fmt.Printf("发送方: 收到 ACK %d (单独确认)\n", ackNum)
	sr.ACKedFrames[ackNum] = true

	// 如果是窗口基序号,滑动窗口
	if ackNum == sr.SendBase%(sr.MaxSeqNum+1) {
		for {
			currentSeq := sr.SendBase % (sr.MaxSeqNum + 1)
			if !sr.ACKedFrames[currentSeq] {
				break
			}
			delete(sr.ACKedFrames, currentSeq)
			sr.SendBase++
		}
		fmt.Printf("发送方: 窗口前移到 [%d, %d)\n", sr.SendBase, sr.SendBase+sr.WindowSize)
	}
}

// SlidingWindowExample 滑动窗口协议示例
func SlidingWindowExample() {
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("【数据链路层 - 滑动窗口协议示例】")
	fmt.Println(strings.Repeat("─", 50))

	// 1. 停止-等待协议
	fmt.Println("\n1️⃣  停止-等待协议 (Stop-and-Wait):")
	fmt.Println("特点: 发送窗口 = 1, 接收窗口 = 1, 序号 0/1 交替")
	sw := NewStopAndWait()
	fmt.Println("\n场景: 正常发送")
	sw.Send("数据A")
	sw.Receive(&Frame{SeqNum: 0, Data: "数据A"})
	sw.ACK()
	sw.Send("数据B")
	sw.Receive(&Frame{SeqNum: 1, Data: "数据B"})
	sw.ACK()

	// 2. 回退 N 帧协议
	fmt.Println("\n\n2️⃣  回退 N 帧协议 (Go-Back-N):")
	fmt.Println("特点: 发送窗口 > 1, 接收窗口 = 1, 累积确认")
	gbn := NewGoBackN(4, 3) // 窗口大小 4, 序号 0-7 (3 位)
	fmt.Println("\n场景 1: 连续发送")
	for i := 0; i < 5; i++ {
		gbn.Send(fmt.Sprintf("数据%d", i))
	}

	fmt.Println("\n场景 2: 正常接收")
	gbn.Receive(&Frame{SeqNum: 0, Data: "数据0"})
	ack0 := gbn.Receive(&Frame{SeqNum: 1, Data: "数据1"})
	gbn.ACK(ack0)

	fmt.Println("\n场景 3: 丢失帧 2, 回退重传")
	gbn.Receive(&Frame{SeqNum: 3, Data: "数据3"}) // 失序
	gbn.Timeout()

	// 3. 选择重传协议
	fmt.Println("\n\n3️⃣  选择重传协议 (Selective Repeat):")
	fmt.Println("特点: 发送窗口 = 接收窗口, 单独确认, 选择重传")
	sr := NewSelectiveRepeat(4, 3) // 窗口大小 4, 序号 0-7
	fmt.Println("\n场景 1: 连续发送")
	for i := 0; i < 5; i++ {
		sr.Send(fmt.Sprintf("数据%d", i))
	}

	fmt.Println("\n场景 2: 失序接收 (帧 1 丢失)")
	sr.Receive(&Frame{SeqNum: 0, Data: "数据0"})
	sr.ACK(0)
	sr.Receive(&Frame{SeqNum: 2, Data: "数据2"}) // 失序,缓存
	sr.ACK(2)
	sr.Receive(&Frame{SeqNum: 3, Data: "数据3"}) // 失序,缓存
	sr.ACK(3)

	fmt.Println("\n场景 3: 接收丢失的帧 1")
	sr.Receive(&Frame{SeqNum: 1, Data: "数据1"}) // 填补空缺,交付
	sr.ACK(1)

	// 408 考点提示
	fmt.Println("\n📚 408 考点总结:")
	fmt.Println("  ✓ 停止-等待: 效率低,序号 0/1 交替")
	fmt.Println("  ✓ GBN: 发送窗口 Ws ∈ [1, 2^n-1], 接收窗口 Wr = 1")
	fmt.Println("  ✓ GBN: 累积确认, ACK n 表示 n 及之前都收到")
	fmt.Println("  ✓ SR: 发送窗口 Ws = 接收窗口 Wr, Ws ≤ 2^(n-1)")
	fmt.Println("  ✓ SR: 单独确认,选择重传出错帧")
	fmt.Println("  ✓ 窗口大小限制: 避免新旧帧混淆")
}
