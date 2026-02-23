# 计算机科学基础学习项目 - 新手入门指南

> 本指南帮助你从零开始，系统性地学习和使用这个项目。

## 📌 这个项目是什么？

这是一个**面向考研的计算机科学实践学习项目**，通过代码实现帮助你深入理解：
- **408统考**: 数据结构与算法、操作系统、计算机组成原理、计算机网络
- **数学一**: 高等数学、线性代数、概率论（基础+进阶）
- **英语一**: 词汇学习方法、阅读理解策略

**特点**：不只是看理论，而是通过**可运行的代码**来理解每个概念。

---

## 📂 项目结构一览

```
Structured_CS_Learning/
│
├── CS_Core_Courses/              # 408统考（Go语言实现）
│   ├── data_structures/          # 数据结构
│   │   ├── basic/               # 基础：数组、链表、栈、队列、哈希表
│   │   ├── linear/              # 进阶：串、稀疏矩阵
│   │   ├── tree/                # 树：二叉搜索树、堆
│   │   ├── graph/               # 图：遍历、最短路径、拓扑排序
│   │   └── algorithm/           # 算法：排序、查找、DP、贪心、回溯、KMP
│   ├── operating_system/         # 操作系统
│   │   ├── process/             # 进程管理与调度算法
│   │   ├── memory/              # 内存管理（分页、分段）
│   │   ├── synchronization/     # 进程同步（信号量、互斥锁）
│   │   ├── filesystem/          # 文件系统（inode、目录、分配方式）
│   │   └── scheduling/          # 磁盘调度与死锁处理
│   ├── computer_architecture/    # 计算机组成原理
│   │   ├── cpu/                 # CPU：寄存器、ALU
│   │   ├── memory/              # 存储器：Cache、虚拟内存
│   │   ├── instruction_set/     # 指令系统
│   │   ├── pipeline/            # 流水线
│   │   └── bus/                 # 总线（理论）
│   └── computer_networks/        # 计算机网络
│       ├── application/         # 应用层：HTTP
│       ├── transport/           # 传输层：TCP
│       ├── network/             # 网络层：IP、路由、ARP
│       ├── datalink/            # 数据链路层
│       ├── physical/            # 物理层（理论）
│       └── protocols/           # 协议：DNS
│
├── Math_Foundations/              # 数学一（Python实现）
│   ├── calculus/                 # 高等数学基础
│   ├── linear_algebra/           # 线性代数基础
│   ├── probability/              # 概率论基础
│   ├── calculus_advanced/        # 高数进阶：曲线曲面积分、场论
│   ├── linear_algebra_advanced/  # 线代进阶：二次型、向量空间
│   └── probability_advanced/     # 概率进阶：数理统计、回归
│
└── English_Learning/             # 英语一
    ├── vocabulary/              # 词汇学习方法与工具
    ├── reading/                 # 阅读理解策略
    └── docs/                    # 考试分析
```

---

## 🎯 我该先学什么？（推荐学习路线）

### 路线一：从基础开始（推荐新手）

```
第1周：数据结构基础
    └── CS_Core_Courses/data_structures/basic/
        ├── array.go      → 理解数组的本质
        ├── linkedlist.go → 理解指针和链式存储
        ├── stack.go      → 理解LIFO
        ├── queue.go      → 理解FIFO
        └── hashtable.go  → 理解哈希映射

第2周：树与图
    └── CS_Core_Courses/data_structures/
        ├── tree/         → 二叉树、堆
        └── graph/        → 图的表示与遍历

第3周：操作系统
    └── CS_Core_Courses/operating_system/
        ├── process/      → 进程是什么？如何调度？
        ├── memory/       → 内存如何分配？
        └── synchronization/ → 多进程如何协作？

第4周：计算机组成 + 网络
    └── CS_Core_Courses/
        ├── computer_architecture/cpu/ → CPU如何工作？
        └── computer_networks/         → 网络如何通信？
```

### 路线二：按兴趣模块学习

如果你有特定目标，可以直接跳到对应模块：

| 目标 | 推荐模块 |
|------|----------|
| 准备算法面试 | `data_structures/` 全部 |
| 理解操作系统 | `operating_system/` 全部 |
| 学习网络编程 | `computer_networks/` |
| 机器学习数学基础 | `Math_Foundations/` 全部 |

---

## 🔧 环境准备

### Go语言环境（CS_Core_Courses）

1. **安装Go**：从 https://golang.org/dl/ 下载安装
2. **验证安装**：
   ```bash
   go version
   # 应显示：go version go1.21.x 或更高
   ```

### Python环境（Math_Foundations）

1. **安装Python 3.8+**
2. **安装依赖**：
   ```bash
   cd Math_Foundations
   pip install numpy matplotlib scipy sympy seaborn
   ```

---

## 🚀 如何运行代码

### 运行Go代码

```bash
# 进入项目目录
cd CS_Core_Courses

# 方式1：运行所有示例
go run main.go

# 方式2：单独运行某个模块（需要先了解代码）
# 直接查看某个文件末尾的 XXXExample() 函数
```

### 运行Python代码

```bash
# 进入项目目录
cd Math_Foundations

# 交互式菜单
python main.py

# 或直接运行某个模块
python main.py --calculus        # 高等数学
python main.py --linear-algebra  # 线性代数
python main.py --probability     # 概率论
```

---

## 📖 每个模块如何阅读？

### 阅读代码的通用方法

每个 `.go` 或 `.py` 文件通常包含：

```
1. 数据结构定义（type XXX struct）
   ↓
2. 创建函数（NewXXX）
   ↓
3. 核心方法（增删改查等操作）
   ↓
4. 示例函数（XXXExample）← 从这里开始阅读！
```

**建议**：先看文件末尾的 `XXXExample()` 函数，它演示了该数据结构的完整用法。

### 具体模块阅读指南

#### 📁 data_structures/basic/ —— 基础数据结构

| 文件 | 核心概念 | 阅读重点 |
|------|----------|----------|
| `array.go` | 动态数组 | 扩容机制、时间复杂度 |
| `linkedlist.go` | 链表 | 指针操作、头尾插入删除 |
| `stack.go` | 栈 | LIFO特性、括号匹配应用 |
| `queue.go` | 队列 | FIFO特性、循环队列实现 |
| `hashtable.go` | 哈希表 | 哈希函数、冲突解决 |

#### 📁 data_structures/tree/ —— 树结构

| 文件 | 核心概念 | 阅读重点 |
|------|----------|----------|
| `bst.go` | 二叉搜索树 | 插入/删除/查找、四种遍历 |
| `heap.go` | 堆 | 堆性质、堆排序、优先队列 |

#### 📁 data_structures/graph/ —— 图结构

| 文件 | 核心概念 | 阅读重点 |
|------|----------|----------|
| `graph.go` | 图的算法 | 邻接表、BFS、DFS、Dijkstra |

#### 📁 operating_system/ —— 操作系统

| 目录 | 核心概念 | 阅读重点 |
|------|----------|----------|
| `process/` | 进程管理 | PCB结构、FCFS/SJF/RR调度 |
| `memory/` | 内存管理 | 首次/最佳/最差适应、分页 |
| `synchronization/` | 进程同步 | 信号量、生产者消费者问题 |

#### 📁 computer_architecture/cpu/ —— CPU设计

| 文件 | 核心概念 | 阅读重点 |
|------|----------|----------|
| `register.go` | 寄存器 | 寄存器类型、标志位 |
| `alu.go` | 算术逻辑单元 | 加减乘除、逻辑运算、状态标志 |

#### 📁 computer_networks/ —— 计算机网络

| 目录 | 核心概念 | 阅读重点 |
|------|----------|----------|
| `application/` | HTTP协议 | 请求/响应结构、状态码 |
| `transport/` | TCP协议 | 三次握手、四次挥手、状态机 |

---

## 💡 学习建议

### 1. 先运行，再阅读
```
运行 XXXExample() → 观察输出 → 带着问题读代码
```

### 2. 动手修改
- 修改参数，观察结果变化
- 尝试添加新功能
- 故意制造错误，理解边界情况

### 3. 画图辅助
- 数据结构：画出内存布局
- 算法：手动模拟执行过程
- 操作系统：画状态转换图

### 4. 对照教材
推荐参考书籍：
- 《算法导论》—— 数据结构与算法
- 《现代操作系统》—— 操作系统
- 《计算机组成与设计》—— 计算机组成原理
- 《计算机网络：自顶向下方法》—— 计算机网络

---

## ❓ 常见问题

### Q: 代码运行报错怎么办？
1. 检查Go/Python版本是否符合要求
2. 确保在正确的目录下运行
3. 查看错误信息，通常会指出问题所在

### Q: 看不懂某个函数？
1. 先看函数注释（代码中的中文注释）
2. 运行对应的Example函数看实际效果
3. 在代码中加 `fmt.Println()` 打印中间状态

### Q: 如何检验自己学会了？
- 能否不看代码，自己实现一遍？
- 能否解释清楚时间/空间复杂度？
- 能否用这个数据结构解决一道LeetCode题？

---

## 📞 下一步

1. **选择一个模块开始** —— 建议从 `data_structures/basic/array.go` 开始
2. **运行示例代码** —— `go run main.go` 或查看具体文件
3. **阅读代码和注释** —— 理解实现细节
4. **动手实践** —— 修改代码，解决练习题

---

## 考研专项学习路线

### 408统考路线
```
阶段一: 数据结构基础 → 树与图 → 排序查找算法 → KMP
阶段二: 操作系统(进程→内存→文件→死锁) → 组成原理(CPU→Cache→指令→流水线)
阶段三: 计算机网络(应用→传输→网络→链路) → 真题强化
```

### 数学一路线
```
阶段一: 高等数学基础 → 线性代数基础 → 概率论基础
阶段二: 高数进阶(曲线曲面积分/场论) → 线代进阶(二次型/向量空间) → 概率进阶(数理统计)
阶段三: 真题训练 → 模拟冲刺
```

### 英语一路线
```
持续: 每天背诵考研词汇(参考 English_Learning/vocabulary/)
阶段一: 掌握阅读策略(reading_strategies_guide.md)
阶段二: 练习六大题型(question_types_analysis.md)
阶段三: 真题精做 → 作文模板
```

---

**Happy Learning!**

*如有问题，欢迎在项目中提出Issue。*
