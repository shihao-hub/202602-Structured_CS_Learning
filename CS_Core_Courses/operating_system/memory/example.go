package memory

import "fmt"

// RunAllMemoryMgmtExamples 运行所有内存管理相关的示例
func RunAllMemoryMgmtExamples() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║           操作系统 - 内存管理扩展模块                    ║")
	fmt.Println("║     (Operating System - Memory Management Extended)     ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	// 运行基础内存管理示例（从memory.go）
	MemoryExample()

	// 运行分页机制示例
	PagingExample()

	// 运行分段机制示例
	SegmentationExample()

	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    模块运行完毕                           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝\n")
}
