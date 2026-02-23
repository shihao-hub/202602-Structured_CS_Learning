package datalink

import "fmt"

func RunAllDatalinkExamples() {
	fmt.Println("\n╔══════════════════════════════════════╗")
	fmt.Println("║    计算机网络 - 数据链路层模块       ║")
	fmt.Println("╚══════════════════════════════════════╝")
	EthernetExample()
	ErrorDetectionExample()
	SlidingWindowExample()
}
