package main

import (
	"fmt"
	"strings"

	"CS_Core_Courses/computer_architecture/cpu"
	"CS_Core_Courses/computer_architecture/instruction_set"
	archmemory "CS_Core_Courses/computer_architecture/memory"
	"CS_Core_Courses/computer_architecture/pipeline"
	"CS_Core_Courses/computer_networks/application"
	"CS_Core_Courses/computer_networks/datalink"
	"CS_Core_Courses/computer_networks/network"
	"CS_Core_Courses/computer_networks/protocols"
	"CS_Core_Courses/computer_networks/transport"
	"CS_Core_Courses/data_structures/algorithm"
	"CS_Core_Courses/data_structures/basic"
	"CS_Core_Courses/data_structures/linear"
	"CS_Core_Courses/operating_system/filesystem"
	osmemory "CS_Core_Courses/operating_system/memory"
	"CS_Core_Courses/operating_system/process"
	"CS_Core_Courses/operating_system/scheduling"
)

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("欢迎来到计算机科学核心课程学习项目!")
	fmt.Println("考研408统考 · 全模块学习平台")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	fmt.Println("本项目包含以下模块:")
	fmt.Println("1. 数据结构与算法 (Data Structures & Algorithms)")
	fmt.Println("2. 操作系统 (Operating Systems)")
	fmt.Println("3. 计算机组成原理 (Computer Architecture)")
	fmt.Println("4. 计算机网络 (Computer Networks)")
	fmt.Println()

	fmt.Println("以下将运行各模块的示例代码:")
	fmt.Println(strings.Repeat("-", 60))

	// ============================
	// 1. 数据结构与算法
	// ============================
	fmt.Println("\n【模块 1: 数据结构与算法】")
	fmt.Println(strings.Repeat("=", 40))

	// 1.1 基础数据结构
	fmt.Println("\n--- 1.1 基础数据结构 ---")
	basic.RunAllBasicExamples()

	// 1.2 线性结构进阶（串、稀疏矩阵）
	fmt.Println("\n--- 1.2 线性结构进阶 ---")
	linear.RunAllLinearExamples()

	// 1.3 算法（排序、查找、DP、贪心、回溯、KMP）
	fmt.Println("\n--- 1.3 算法 ---")
	algorithm.RunAllAlgorithmExamples()

	// ============================
	// 2. 操作系统
	// ============================
	fmt.Println("\n【模块 2: 操作系统】")
	fmt.Println(strings.Repeat("=", 40))

	// 2.1 进程管理
	fmt.Println("\n--- 2.1 进程管理 ---")
	process.RunAllProcessExamples()

	// 2.2 内存管理（基础 + 分页/分段扩展）
	fmt.Println("\n--- 2.2 内存管理 ---")
	osmemory.RunAllMemoryMgmtExamples()

	// 2.3 文件系统
	fmt.Println("\n--- 2.3 文件系统 ---")
	filesystem.RunAllFilesystemExamples()

	// 2.4 磁盘调度与死锁
	fmt.Println("\n--- 2.4 磁盘调度与死锁 ---")
	scheduling.RunAllSchedulingExamples()

	// ============================
	// 3. 计算机组成原理
	// ============================
	fmt.Println("\n【模块 3: 计算机组成原理】")
	fmt.Println(strings.Repeat("=", 40))

	// 3.1 CPU（寄存器、ALU）
	fmt.Println("\n--- 3.1 CPU ---")
	cpu.RunAllCPUExamples()

	// 3.2 存储器层次（Cache、虚拟内存）
	fmt.Println("\n--- 3.2 存储器层次 ---")
	archmemory.RunAllMemoryExamples()

	// 3.3 指令系统
	fmt.Println("\n--- 3.3 指令系统 ---")
	instruction_set.RunAllInstructionSetExamples()

	// 3.4 流水线
	fmt.Println("\n--- 3.4 流水线 ---")
	pipeline.RunAllPipelineExamples()

	// ============================
	// 4. 计算机网络
	// ============================
	fmt.Println("\n【模块 4: 计算机网络】")
	fmt.Println(strings.Repeat("=", 40))

	// 4.1 应用层（HTTP）
	fmt.Println("\n--- 4.1 应用层 ---")
	application.RunAllApplicationExamples()

	// 4.2 传输层（TCP）
	fmt.Println("\n--- 4.2 传输层 ---")
	transport.RunAllTransportExamples()

	// 4.3 网络层（IP、路由、ARP）
	fmt.Println("\n--- 4.3 网络层 ---")
	network.RunAllNetworkExamples()

	// 4.4 数据链路层（以太网、差错检测、滑动窗口）
	fmt.Println("\n--- 4.4 数据链路层 ---")
	datalink.RunAllDatalinkExamples()

	// 4.5 网络协议（DNS等）
	fmt.Println("\n--- 4.5 网络协议 ---")
	protocols.RunAllProtocolExamples()

	// ============================
	// 结束
	// ============================
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("所有示例代码运行完成!")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	fmt.Println("学习建议:")
	fmt.Println("1. 仔细阅读每个模块的代码实现和中文注释")
	fmt.Println("2. 参考各模块的README文档和408考点对照")
	fmt.Println("3. 尝试修改代码参数并观察结果变化")
	fmt.Println("4. 结合真题练习，巩固知识点")
	fmt.Println()

	fmt.Println("项目结构说明:")
	fmt.Println("- data_structures/  数据结构(基础+线性进阶+算法)")
	fmt.Println("- operating_system/ 操作系统(进程+内存+文件系统+调度)")
	fmt.Println("- computer_architecture/ 计算机组成(CPU+存储+指令+流水线)")
	fmt.Println("- computer_networks/ 计算机网络(应用+传输+网络+链路+协议)")
	fmt.Println()

	fmt.Println("Happy Learning!")
}
