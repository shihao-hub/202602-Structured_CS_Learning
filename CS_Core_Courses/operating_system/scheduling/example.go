package scheduling

import "fmt"

// RunAllSchedulingExamples 运行所有调度相关的示例
func RunAllSchedulingExamples() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║       操作系统 - 磁盘调度与死锁处理模块                  ║")
	fmt.Println("║    (Operating System - Scheduling & Deadlock Module)    ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	// 运行磁盘调度算法示例
	DiskSchedulerExample()

	// 运行死锁处理示例
	DeadlockExample()

	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    模块运行完毕                           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝\n")
}
