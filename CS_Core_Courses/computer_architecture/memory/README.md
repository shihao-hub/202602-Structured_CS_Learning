# 存储器层次 (Memory Hierarchy)

## 408 考试映射

本模块涵盖计算机组成原理 408 考试中的存储器层次相关内容：

### 1. Cache 高速缓存
- **映射方式**
  - 直接映射 (Direct Mapped)
  - 全相联映射 (Fully Associative)
  - 组相联映射 (Set Associative)
- **替换算法**
  - LRU (Least Recently Used)
  - FIFO (First In First Out)
  - 随机替换 (Random)
- **性能指标**
  - 命中率 (Hit Rate)
  - 缺失率 (Miss Rate)
  - 平均访问时间 (Average Access Time)

### 2. 虚拟存储器
- **地址转换机制**
  - 页表 (Page Table)
  - TLB (Translation Lookaside Buffer)
  - 地址映射过程
- **页面替换算法**
  - FIFO
  - LRU
  - Clock/NRU (Not Recently Used)
  - OPT (Optimal)
- **页面管理**
  - 页表项结构 (有效位、修改位、访问位)
  - 缺页中断处理

## 文件说明

- `cache.go` - Cache 存储器模拟实现
- `virtual_memory.go` - 虚拟存储器机制实现
- `example.go` - 示例程序入口

## 运行示例

```go
package main

import "Structured_CS_Learning/CS_Core_Courses/computer_architecture/memory"

func main() {
    memory.RunAllMemoryExamples()
}
```

## 408 考试重点

1. **Cache 计算题**
   - 给定主存地址，计算 Cache 行号、标记位
   - 计算 Cache 总容量、标记位数
   - 根据访问序列计算命中率

2. **虚拟存储器计算**
   - 虚拟地址到物理地址的转换
   - 页表项结构和页表大小计算
   - TLB 命中对性能的影响

3. **理解概念**
   - 局部性原理 (时间局部性、空间局部性)
   - 存储器层次结构的设计目标
   - Cache 一致性问题
