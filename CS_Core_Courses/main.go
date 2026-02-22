package main

import (
	"fmt"
	"strings"

	"CS_Core_Courses/computer_architecture/cpu"
	"CS_Core_Courses/computer_networks/application"
	"CS_Core_Courses/computer_networks/transport"
	"CS_Core_Courses/data_structures/basic"
	"CS_Core_Courses/operating_system/process"
)

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("欢迎来到计算机科学核心课程学习项目!")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	fmt.Println("本项目包含以下模块:")
	fmt.Println("1. 数据结构与算法 (Data Structures & Algorithms)")
	fmt.Println("2. 操作系统 (Operating Systems)")
	fmt.Println("3. 计算机组成原理 (Computer Architecture)")
	fmt.Println("4. 计算机网络 (Computer Networks)")
	fmt.Println()

	// 运行所有示例
	fmt.Println("以下将运行各模块的示例代码:")
	fmt.Println(strings.Repeat("-", 60))

	// 1. 数据结构示例
	fmt.Println("\n【模块 1: 数据结构与算法】")
	fmt.Println(strings.Repeat("=", 40))
	basic.RunAllBasicExamples()

	// 2. 操作系统示例
	fmt.Println("\n【模块 2: 操作系统】")
	fmt.Println(strings.Repeat("=", 40))
	process.RunAllProcessExamples()

	// 3. 计算机组成原理示例
	fmt.Println("\n【模块 3: 计算机组成原理】")
	fmt.Println(strings.Repeat("=", 40))
	cpu.RunAllCPUExamples()

	// 4. 计算机网络示例
	fmt.Println("\n【模块 4: 计算机网络】")
	fmt.Println(strings.Repeat("=", 40))
	application.RunAllApplicationExamples()
	transport.RunAllTransportExamples()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("所有示例代码运行完成!")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	fmt.Println("学习建议:")
	fmt.Println("1. 仔细阅读每个模块的代码实现")
	fmt.Println("2. 参考各模块的README文档")
	fmt.Println("3. 尝试修改代码并观察结果变化")
	fmt.Println("4. 完成课后练习题")
	fmt.Println("5. 将知识应用到实际项目中")
	fmt.Println()

	fmt.Println("项目结构说明:")
	fmt.Println("- README.md: 项目总体说明")
	fmt.Println("- data_structures/: 数据结构模块")
	fmt.Println("- operating_system/: 操作系统模块")
	fmt.Println("- computer_architecture/: 计算机组成原理模块")
	fmt.Println("- computer_networks/: 计算机网络模块")
	fmt.Println()

	fmt.Println("学习提示:")
	fmt.Println("- 每个模块都可以独立运行学习")
	fmt.Println("- 代码中包含详细的中文注释")
	fmt.Println("- 如遇问题，欢迎提出issue讨论")
	fmt.Println()

	fmt.Println("Happy Learning! 🚀")
}
