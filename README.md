# 计算机考研系统化学习项目

面向计算机考研（数学一 + 408统考 + 英语一）的代码驱动学习平台。

## 项目定位

通过**可运行的代码**和**系统化文档**帮助考研学生深入理解核心知识点。每个概念不只是看理论，而是通过代码实现来加深理解。

## 覆盖科目

| 科目 | 目录 | 语言 | 内容层次 |
|------|------|------|---------|
| 数学一 | `Math_Foundations/` | Python | 基础 + 进阶 |
| 408统考 | `CS_Core_Courses/` | Go | 基础 + 进阶 |
| 英语一 | `English_Learning/` | Python + Markdown | 词汇 + 阅读 |
| 政治 | `Study_Plans/` | Markdown | 学习计划 |
| **学习计划** | `Study_Plans/` | Markdown | 各科备考规划 |

## 项目结构

```
Structured_CS_Learning/
│
├── CS_Core_Courses/              # 408统考（Go语言）
│   ├── data_structures/          # 数据结构与算法
│   │   ├── basic/               # 基础：数组、链表、栈、队列、哈希表
│   │   ├── linear/              # 进阶：串、稀疏矩阵
│   │   ├── tree/                # 树：BST、堆
│   │   ├── graph/               # 图：遍历、最短路径
│   │   └── algorithm/           # 算法：排序、查找、DP、贪心、回溯、KMP
│   ├── operating_system/         # 操作系统
│   │   ├── process/             # 进程管理与调度
│   │   ├── memory/              # 内存管理（分页、分段）
│   │   ├── synchronization/     # 进程同步
│   │   ├── filesystem/          # 文件系统
│   │   └── scheduling/          # 磁盘调度与死锁
│   ├── computer_architecture/    # 计算机组成原理
│   │   ├── cpu/                 # CPU：寄存器、ALU
│   │   ├── memory/              # 存储：Cache、虚拟内存
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
├── Math_Foundations/              # 数学一（Python）
│   ├── calculus/                 # 高等数学基础
│   ├── linear_algebra/           # 线性代数基础
│   ├── probability/              # 概率论基础
│   ├── calculus_advanced/        # 高数进阶：曲线曲面积分、场论
│   ├── linear_algebra_advanced/  # 线代进阶：二次型、向量空间
│   └── probability_advanced/     # 概率进阶：数理统计、回归
│
├── English_Learning/             # 英语一（Python + Markdown）
│   ├── vocabulary/              # 词汇学习方法与工具
│   ├── reading/                 # 阅读理解策略
│   └── docs/                    # 考试分析
│
├── Study_Plans/                  # 考研学习计划
│   ├── README.md                # 总览与时间规划
│   ├── math-calculus-plan.md    # 高等数学计划
│   ├── math-linear-algebra-plan.md  # 线性代数计划
│   ├── math-probability-plan.md # 概率论计划
│   ├── cs408-data-structure-plan.md      # 数据结构计划
│   ├── cs408-operating-system-plan.md    # 操作系统计划
│   ├── cs408-computer-organization-plan.md  # 组成原理计划
│   ├── cs408-computer-network-plan.md    # 计算机网络计划
│   ├── english-plan.md          # 英语一计划
│   └── politics-plan.md         # 政治计划
│
├── GETTING_STARTED.md            # 新手入门指南
├── CLAUDE.md                     # Claude Code 项目指南
└── README.md                     # 本文件
```

## 快速开始

详见 [GETTING_STARTED.md](GETTING_STARTED.md)

### 408统考 (Go)
```bash
cd CS_Core_Courses
go build ./...     # 编译验证
go run main.go     # 运行所有示例
```

### 数学一 (Python)
```bash
cd Math_Foundations
pip install numpy matplotlib scipy sympy seaborn
python main.py     # 交互式菜单
```

### 英语一
```bash
cd English_Learning
python main.py     # 词汇学习工具
```

### 学习计划
```bash
cd Study_Plans
# 查看各科学习计划
cat README.md      # 总览与时间规划
```

## 项目特色

- **理论+代码双轨**: 每个知识点都有 guide.md (理论) + examples (代码)
- **考点全覆盖**: 每个模块 README 标注对应考研大纲考点
- **中文注释**: 所有代码含详细中文注释
- **可视化**: Python 模块使用 matplotlib 进行数据可视化
- **循序渐进**: 基础 → 进阶分层设计
