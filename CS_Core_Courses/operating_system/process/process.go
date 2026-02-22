package process

import (
	"fmt"
	"time"
)

// ProcessState 进程状态
type ProcessState int

const (
	StateNew ProcessState = iota
	StateReady
	StateRunning
	StateWaiting
	StateTerminated
)

// String 返回状态的字符串表示
func (ps ProcessState) String() string {
	switch ps {
	case StateNew:
		return "New"
	case StateReady:
		return "Ready"
	case StateRunning:
		return "Running"
	case StateWaiting:
		return "Waiting"
	case StateTerminated:
		return "Terminated"
	default:
		return "Unknown"
	}
}

// ProcessControlBlock 进程控制块 (PCB)
type ProcessControlBlock struct {
	PID             int           // 进程ID
	ParentPID       int           // 父进程ID
	State           ProcessState  // 进程状态
	Priority        int           // 优先级
	BurstTime       int           // 需要的CPU时间
	RemainingTime   int           // 剩余时间
	ArrivalTime     int           // 到达时间
	WaitTime        int           // 等待时间
	TurnaroundTime  int           // 周转时间
	ResponseTime    int           // 响应时间
	ContextSwitches int           // 上下文切换次数
	MemoryUsage     int           // 内存使用量
	StartTime       time.Time     // 开始时间
	EndTime         time.Time     // 结束时间
	Children        []int         // 子进程PID列表
}

// NewProcess 创建新进程
func NewProcess(pid, parentPid, priority, burstTime, arrivalTime int) *ProcessControlBlock {
	return &ProcessControlBlock{
		PID:             pid,
		ParentPID:       parentPid,
		State:           StateNew,
		Priority:        priority,
		BurstTime:       burstTime,
		RemainingTime:   burstTime,
		ArrivalTime:     arrivalTime,
		WaitTime:        0,
		TurnaroundTime:  0,
		ResponseTime:    -1, // -1表示还未首次响应
		ContextSwitches: 0,
		MemoryUsage:     0,
		Children:        make([]int, 0),
	}
}

// SetState 设置进程状态
func (pcb *ProcessControlBlock) SetState(state ProcessState) {
	pcb.State = state
}

// Execute 执行进程一个时间片
func (pcb *ProcessControlBlock) Execute(timeSlice int) {
	if pcb.RemainingTime > 0 {
		actualTime := timeSlice
		if pcb.RemainingTime < timeSlice {
			actualTime = pcb.RemainingTime
		}

		pcb.RemainingTime -= actualTime

		// 首次响应
		if pcb.ResponseTime == -1 {
			pcb.ResponseTime = pcb.WaitTime
		}

		if pcb.RemainingTime == 0 {
			pcb.State = StateTerminated
			pcb.EndTime = time.Now()
			pcb.TurnaroundTime = int(pcb.EndTime.Sub(pcb.StartTime).Milliseconds())
		}

		pcb.ContextSwitches++
	}
}

// IsCompleted 检查进程是否完成
func (pcb *ProcessControlBlock) IsCompleted() bool {
	return pcb.RemainingTime == 0
}

// AddChild 添加子进程
func (pcb *ProcessControlBlock) AddChild(childPid int) {
	pcb.Children = append(pcb.Children, childPid)
}

// Print 打印进程信息
func (pcb *ProcessControlBlock) Print() {
	fmt.Printf("PID: %d, ParentPID: %d, State: %s, Priority: %d, BurstTime: %d, RemainingTime: %d, ArrivalTime: %d\n",
		pcb.PID, pcb.ParentPID, pcb.State, pcb.Priority, pcb.BurstTime, pcb.RemainingTime, pcb.ArrivalTime)
	fmt.Printf("  WaitTime: %d, ResponseTime: %d, TurnaroundTime: %d, ContextSwitches: %d\n",
		pcb.WaitTime, pcb.ResponseTime, pcb.TurnaroundTime, pcb.ContextSwitches)
	if len(pcb.Children) > 0 {
		fmt.Printf("  Children: %v\n", pcb.Children)
	}
}

// ProcessManager 进程管理器
type ProcessManager struct {
	processes map[int]*ProcessControlBlock
	nextPID   int
}

// NewProcessManager 创建进程管理器
func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		processes: make(map[int]*ProcessControlBlock),
		nextPID:   1,
	}
}

// CreateProcess 创建进程
func (pm *ProcessManager) CreateProcess(parentPid, priority, burstTime, arrivalTime int) int {
	pid := pm.nextPID
	pm.nextPID++

	process := NewProcess(pid, parentPid, priority, burstTime, arrivalTime)
	process.StartTime = time.Now()
	process.State = StateReady
	pm.processes[pid] = process

	// 如果有父进程，添加到父进程的子进程列表
	if parentPid > 0 {
		if parent, exists := pm.processes[parentPid]; exists {
			parent.AddChild(pid)
		}
	}

	return pid
}

// GetProcess 获取进程
func (pm *ProcessManager) GetProcess(pid int) *ProcessControlBlock {
	return pm.processes[pid]
}

// GetReadyProcesses 获取就绪进程列表
func (pm *ProcessManager) GetReadyProcesses() []*ProcessControlBlock {
	var ready []*ProcessControlBlock
	for _, process := range pm.processes {
		if process.State == StateReady {
			ready = append(ready, process)
		}
	}
	return ready
}

// GetRunningProcess 获取正在运行的进程
func (pm *ProcessManager) GetRunningProcess() *ProcessControlBlock {
	for _, process := range pm.processes {
		if process.State == StateRunning {
			return process
		}
	}
	return nil
}

// GetTerminatedProcesses 获取已终止的进程列表
func (pm *ProcessManager) GetTerminatedProcesses() []*ProcessControlBlock {
	var terminated []*ProcessControlBlock
	for _, process := range pm.processes {
		if process.State == StateTerminated {
			terminated = append(terminated, process)
		}
	}
	return terminated
}

// PrintAllProcesses 打印所有进程信息
func (pm *ProcessManager) PrintAllProcesses() {
	fmt.Println("=== 所有进程信息 ===")
	for _, process := range pm.processes {
		process.Print()
		fmt.Println()
	}
}

// GetAverageTurnaroundTime 获取平均周转时间
func (pm *ProcessManager) GetAverageTurnaroundTime() float64 {
	var totalTurnaround int
	var count int

	for _, process := range pm.processes {
		if process.State == StateTerminated {
			totalTurnaround += process.TurnaroundTime
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return float64(totalTurnaround) / float64(count)
}

// GetAverageWaitTime 获取平均等待时间
func (pm *ProcessManager) GetAverageWaitTime() float64 {
	var totalWait int
	var count int

	for _, process := range pm.processes {
		if process.State == StateTerminated {
			totalWait += process.WaitTime
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return float64(totalWait) / float64(count)
}

// GetCPUUtilization 获取CPU利用率
func (pm *ProcessManager) GetCPUUtilization() float64 {
	var totalBurst int
	var totalActive int

	for _, process := range pm.processes {
		totalBurst += process.BurstTime
		if process.State == StateTerminated {
			totalActive += process.TurnaroundTime
		}
	}

	if totalActive == 0 {
		return 0
	}
	return float64(totalBurst) / float64(totalActive) * 100
}

// 示例函数
func ProcessExample() {
	fmt.Println("=== 进程 (Process) 示例 ===")

	// 创建进程管理器
	pm := NewProcessManager()

	// 创建进程
	fmt.Println("创建进程:")
	pid1 := pm.CreateProcess(0, 2, 5, 0)  // 父进程ID=0(无父进程), 优先级=2, 执行时间=5, 到达时间=0
	pid2 := pm.CreateProcess(0, 1, 3, 1)  // 优先级=1, 执行时间=3, 到达时间=1
	pid3 := pm.CreateProcess(pid1, 3, 8, 2) // pid1的子进程, 优先级=3, 执行时间=8, 到达时间=2
	pid4 := pm.CreateProcess(0, 1, 6, 3)  // 优先级=1, 执行时间=6, 到达时间=3

	fmt.Printf("创建进程 PID: %d (主进程)\n", pid1)
	fmt.Printf("创建进程 PID: %d (主进程)\n", pid2)
	fmt.Printf("创建进程 PID: %d (PID %d的子进程)\n", pid3, pid1)
	fmt.Printf("创建进程 PID: %d (主进程)\n", pid4)

	// 模拟进程执行
	fmt.Println("\n模拟进程执行:")

	// 设置一些进程为运行状态
	process1 := pm.GetProcess(pid1)
	process1.State = StateRunning
	process1.Execute(2) // 执行2个时间单位
	process1.Print()

	// 获取就绪进程
	fmt.Println("\n就绪进程列表:")
	readyProcesses := pm.GetReadyProcesses()
	for _, proc := range readyProcesses {
		fmt.Printf("PID: %d, Priority: %d\n", proc.PID, proc.Priority)
	}

	// 模拟更多执行
	fmt.Println("\n继续模拟执行:")
	process2 := pm.GetProcess(pid2)
	process2.Execute(3) // 执行完成
	process2.Print()

	process1.Execute(3) // 执行完成
	process1.Print()

	// 显示终止的进程
	fmt.Println("\n已终止的进程:")
	terminated := pm.GetTerminatedProcesses()
	for _, proc := range terminated {
		fmt.Printf("PID: %d, TurnaroundTime: %d, WaitTime: %d\n",
			proc.PID, proc.TurnaroundTime, proc.WaitTime)
	}

	// 统计信息
	fmt.Println("\n统计信息:")
	fmt.Printf("平均周转时间: %.2f\n", pm.GetAverageTurnaroundTime())
	fmt.Printf("平均等待时间: %.2f\n", pm.GetAverageWaitTime())
	fmt.Printf("CPU利用率: %.2f%%\n", pm.GetCPUUtilization())

	// 打印所有进程信息
	pm.PrintAllProcesses()
	fmt.Println()
}