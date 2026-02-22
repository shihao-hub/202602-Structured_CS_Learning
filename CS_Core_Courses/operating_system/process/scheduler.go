package process

import (
	"fmt"
	"sort"
)

// Scheduler 调度器接口
type Scheduler interface {
	AddProcess(process *ProcessControlBlock)
	NextProcess() *ProcessControlBlock
	HasNext() bool
	RemoveProcess(pid int) bool
}

// FCFSScheduler 先来先服务调度器
type FCFSScheduler struct {
	queue []*ProcessControlBlock
}

func NewFCFSScheduler() *FCFSScheduler {
	return &FCFSScheduler{
		queue: make([]*ProcessControlBlock, 0),
	}
}

func (s *FCFSScheduler) AddProcess(process *ProcessControlBlock) {
	s.queue = append(s.queue, process)
}

func (s *FCFSScheduler) NextProcess() *ProcessControlBlock {
	if len(s.queue) == 0 {
		return nil
	}

	// 按到达时间排序
	sort.Slice(s.queue, func(i, j int) bool {
		return s.queue[i].ArrivalTime < s.queue[j].ArrivalTime
	})

	process := s.queue[0]
	s.queue = s.queue[1:]
	return process
}

func (s *FCFSScheduler) HasNext() bool {
	return len(s.queue) > 0
}

func (s *FCFSScheduler) RemoveProcess(pid int) bool {
	for i, proc := range s.queue {
		if proc.PID == pid {
			s.queue = append(s.queue[:i], s.queue[i+1:]...)
			return true
		}
	}
	return false
}

// SJFScheduler 最短作业优先调度器
type SJFScheduler struct {
	queue []*ProcessControlBlock
}

func NewSJFScheduler() *SJFScheduler {
	return &SJFScheduler{
		queue: make([]*ProcessControlBlock, 0),
	}
}

func (s *SJFScheduler) AddProcess(process *ProcessControlBlock) {
	s.queue = append(s.queue, process)
}

func (s *SJFScheduler) NextProcess() *ProcessControlBlock {
	if len(s.queue) == 0 {
		return nil
	}

	// 按执行时间排序（最短优先）
	sort.Slice(s.queue, func(i, j int) bool {
		if s.queue[i].RemainingTime == s.queue[j].RemainingTime {
			return s.queue[i].ArrivalTime < s.queue[j].ArrivalTime
		}
		return s.queue[i].RemainingTime < s.queue[j].RemainingTime
	})

	process := s.queue[0]
	s.queue = s.queue[1:]
	return process
}

func (s *SJFScheduler) HasNext() bool {
	return len(s.queue) > 0
}

func (s *SJFScheduler) RemoveProcess(pid int) bool {
	for i, proc := range s.queue {
		if proc.PID == pid {
			s.queue = append(s.queue[:i], s.queue[i+1:]...)
			return true
		}
	}
	return false
}

// PriorityScheduler 优先级调度器
type PriorityScheduler struct {
	queue []*ProcessControlBlock
}

func NewPriorityScheduler() *PriorityScheduler {
	return &PriorityScheduler{
		queue: make([]*ProcessControlBlock, 0),
	}
}

func (s *PriorityScheduler) AddProcess(process *ProcessControlBlock) {
	s.queue = append(s.queue, process)
}

func (s *PriorityScheduler) NextProcess() *ProcessControlBlock {
	if len(s.queue) == 0 {
		return nil
	}

	// 按优先级排序（数字越小优先级越高）
	sort.Slice(s.queue, func(i, j int) bool {
		if s.queue[i].Priority == s.queue[j].Priority {
			return s.queue[i].ArrivalTime < s.queue[j].ArrivalTime
		}
		return s.queue[i].Priority < s.queue[j].Priority
	})

	process := s.queue[0]
	s.queue = s.queue[1:]
	return process
}

func (s *PriorityScheduler) HasNext() bool {
	return len(s.queue) > 0
}

func (s *PriorityScheduler) RemoveProcess(pid int) bool {
	for i, proc := range s.queue {
		if proc.PID == pid {
			s.queue = append(s.queue[:i], s.queue[i+1:]...)
			return true
		}
	}
	return false
}

// RRScheduler 时间片轮转调度器
type RRScheduler struct {
	queue      []*ProcessControlBlock
	timeSlice  int
	currentPos int
}

func NewRRScheduler(timeSlice int) *RRScheduler {
	return &RRScheduler{
		queue:      make([]*ProcessControlBlock, 0),
		timeSlice:  timeSlice,
		currentPos: 0,
	}
}

func (s *RRScheduler) AddProcess(process *ProcessControlBlock) {
	s.queue = append(s.queue, process)
}

func (s *RRScheduler) NextProcess() *ProcessControlBlock {
	if len(s.queue) == 0 {
		return nil
	}

	// 轮转到下一个位置
	if s.currentPos >= len(s.queue) {
		s.currentPos = 0
	}

	process := s.queue[s.currentPos]
	s.currentPos++

	// 如果是最后一个位置，重置
	if s.currentPos >= len(s.queue) {
		s.currentPos = 0
	}

	return process
}

func (s *RRScheduler) HasNext() bool {
	return len(s.queue) > 0
}

func (s *RRScheduler) RemoveProcess(pid int) bool {
	for i, proc := range s.queue {
		if proc.PID == pid {
			s.queue = append(s.queue[:i], s.queue[i+1:]...)
			// 如果删除的当前位置之前的进程，需要调整位置
			if i < s.currentPos {
				s.currentPos--
			}
			if s.currentPos < 0 {
				s.currentPos = 0
			}
			return true
		}
	}
	return false
}

// Simulation 调度模拟器
type Simulation struct {
	processes []*ProcessControlBlock
	scheduler Scheduler
	time      int
}

// NewSimulation 创建调度模拟
func NewSimulation(scheduler Scheduler, processes []*ProcessControlBlock) *Simulation {
	// 将所有进程加入调度器
	for _, proc := range processes {
		scheduler.AddProcess(proc)
	}

	return &Simulation{
		processes: processes,
		scheduler: scheduler,
		time:      0,
	}
}

// Run 运行模拟
func (sim *Simulation) Run() {
	fmt.Printf("=== 开始调度模拟 (时间: %d) ===\n", sim.time)

	for sim.scheduler.HasNext() {
		process := sim.scheduler.NextProcess()

		if process == nil {
			sim.time++
			continue
		}

		fmt.Printf("时间 %d: 调度进程 PID:%d (剩余时间: %d)\n",
			sim.time, process.PID, process.RemainingTime)

		// 设置进程为运行状态
		process.State = StateRunning
		process.WaitTime = sim.time - process.ArrivalTime

		// 根据调度器类型执行
		switch scheduler := sim.scheduler.(type) {
		case *RRScheduler:
			process.Execute(scheduler.timeSlice)
		default:
			// FCFS, SJF, Priority 执行到完成
			process.Execute(process.RemainingTime)
		}

		if process.IsCompleted() {
			fmt.Printf("进程 PID:%d 完成\n", process.PID)
		} else {
			// 对于时间片轮转，如果进程未完成，重新加入队列
			if rrScheduler, ok := sim.scheduler.(*RRScheduler); ok {
				if !process.IsCompleted() {
					process.State = StateReady
					rrScheduler.AddProcess(process)
				}
			}
		}

		sim.time++
	}

	fmt.Printf("=== 调度模拟结束 (总时间: %d) ===\n", sim.time)
}

// 示例函数
func SchedulerExample() {
	fmt.Println("=== 进程调度 (Process Scheduling) 示例 ===")

	// 创建测试进程
	processes := []*ProcessControlBlock{
		NewProcess(1, 0, 2, 10, 0), // PID=1, 优先级=2, 执行时间=10, 到达时间=0
		NewProcess(2, 0, 1, 5, 2),  // PID=2, 优先级=1, 执行时间=5, 到达时间=2
		NewProcess(3, 0, 3, 8, 4),  // PID=3, 优先级=3, 执行时间=8, 到达时间=4
		NewProcess(4, 0, 1, 3, 6),  // PID=4, 优先级=1, 执行时间=3, 到达时间=6
	}

	// FCFS调度
	fmt.Println("\n--- FCFS 调度 ---")
	fcfsScheduler := NewFCFSScheduler()
	fcfsSimulation := NewSimulation(fcfsScheduler, processes)
	fcfsSimulation.Run()

	// 重置进程状态
	for _, proc := range processes {
		proc.RemainingTime = proc.BurstTime
		proc.State = StateReady
		proc.WaitTime = 0
		proc.TurnaroundTime = 0
		proc.ResponseTime = -1
	}

	// SJF调度
	fmt.Println("\n--- SJF 调度 ---")
	sjfScheduler := NewSJFScheduler()
	sjfSimulation := NewSimulation(sjfScheduler, processes)
	sjfSimulation.Run()

	// 重置进程状态
	for _, proc := range processes {
		proc.RemainingTime = proc.BurstTime
		proc.State = StateReady
		proc.WaitTime = 0
		proc.TurnaroundTime = 0
		proc.ResponseTime = -1
	}

	// Priority调度
	fmt.Println("\n--- Priority 调度 ---")
	priorityScheduler := NewPriorityScheduler()
	prioritySimulation := NewSimulation(priorityScheduler, processes)
	prioritySimulation.Run()

	// 重置进程状态
	for _, proc := range processes {
		proc.RemainingTime = proc.BurstTime
		proc.State = StateReady
		proc.WaitTime = 0
		proc.TurnaroundTime = 0
		proc.ResponseTime = -1
	}

	// RR调度
	fmt.Println("\n--- RR 调度 (时间片=3) ---")
	rrScheduler := NewRRScheduler(3)
	rrSimulation := NewSimulation(rrScheduler, processes)
	rrSimulation.Run()

	fmt.Println()
}