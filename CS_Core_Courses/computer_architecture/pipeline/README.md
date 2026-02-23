# 流水线技术 (Pipeline)

## 408 考试映射

本模块涵盖计算机组成原理 408 考试中的指令流水线相关内容。

### 1. 流水线基本概念
- **5 段流水线结构**
  - IF (Instruction Fetch): 取指令
  - ID (Instruction Decode): 译码
  - EX (Execute): 执行
  - MEM (Memory Access): 访存
  - WB (Write Back): 写回
- **流水线性能指标**
  - 吞吐率 (Throughput)
  - 加速比 (Speedup)
  - 效率 (Efficiency)

### 2. 流水线冲突
- **结构冲突 (Structural Hazard)**
  - 硬件资源不足导致的冲突
  - 解决方法：增加硬件资源、暂停流水线
- **数据冲突 (Data Hazard)**
  - RAW (Read After Write): 写后读，真相关
  - WAW (Write After Write): 写后写，输出相关
  - WAR (Write After Read): 读后写，反相关
  - 解决方法：转发（forwarding）、暂停（stall）
- **控制冲突 (Control Hazard)**
  - 分支指令导致的冲突
  - 解决方法：分支预测、延迟分支

### 3. 流水线性能计算
- **理论加速比**: Sn = n（n 为流水线段数）
- **实际加速比**: S = (不使用流水线时间) / (使用流水线时间)
- **吞吐率**: TP = 指令数 / 总时间
- **效率**: E = (实际使用的时空区) / (总的时空区) × 100%

## 文件说明

- `pipeline.go` - 5 段流水线模拟实现
- `example.go` - 示例程序入口

## 运行示例

```go
package main

import "Structured_CS_Learning/CS_Core_Courses/computer_architecture/pipeline"

func main() {
    pipeline.RunAllPipelineExamples()
}
```

## 408 考试重点

1. **流水线时空图**
   - 画出指令在流水线各阶段的执行情况
   - 标识冲突和气泡（bubble）

2. **性能计算题**
   - 计算流水线执行时间
   - 计算加速比和吞吐率
   - 分析冲突对性能的影响

3. **冲突识别与解决**
   - 识别数据相关类型（RAW/WAW/WAR）
   - 判断是否可以通过转发解决
   - 计算需要插入的气泡数

4. **超标量与超流水线**
   - 超标量：多发射，增加并行度
   - 超流水线：增加流水线段数，提高时钟频率
