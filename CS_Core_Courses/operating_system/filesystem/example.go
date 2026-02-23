package filesystem

import "fmt"

// RunAllFilesystemExamples 运行所有文件系统示例
func RunAllFilesystemExamples() {
	fmt.Println("\n╔══════════════════════════════════════╗")
	fmt.Println("║       操作系统 - 文件系统模块         ║")
	fmt.Println("╚══════════════════════════════════════╝")

	InodeExample()
	DirectoryExample()
	FileAllocationExample()
}
