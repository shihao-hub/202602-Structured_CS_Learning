package memory

import "fmt"

func RunAllMemoryExamples() {
	fmt.Println("\n╔══════════════════════════════════════╗")
	fmt.Println("║    计算机组成原理 - 存储器层次模块   ║")
	fmt.Println("╚══════════════════════════════════════╝")
	CacheExample()
	VirtualMemoryExample()
}
