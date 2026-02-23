# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

这是一个面向计算机考研（数学一 + 408统考 + 英语一）的**代码驱动学习平台**。通过可运行的代码和系统化文档帮助考研学生深入理解核心知识点。

**核心理念**：不只是看理论，而是通过代码实现来加深理解。

## 项目结构

```
Structured_CS_Learning/
├── CS_Core_Courses/          # 408统考（Go语言）
│   ├── data_structures/      # 数据结构：basic → linear → tree → graph → algorithm
│   ├── operating_system/     # 操作系统：process → memory → synchronization → filesystem → scheduling
│   ├── computer_architecture/ # 计算机组成：cpu → memory → instruction_set → pipeline → bus
│   └── computer_networks/    # 计算机网络：application → transport → network → datalink → physical → protocols
├── Math_Foundations/         # 数学一（Python）
│   ├── calculus/             # 高等数学基础
│   ├── linear_algebra/       # 线性代数基础
│   ├── probability/          # 概率论基础
│   ├── calculus_advanced/    # 高数进阶
│   ├── linear_algebra_advanced/ # 线代进阶
│   └── probability_advanced/ # 概率进阶
├── English_Learning/         # 英语一（Python + Markdown）
│   ├── vocabulary/           # 词汇学习
│   ├── reading/              # 阅读策略
│   └── docs/                 # 考试分析
└── Study_Plans/              # 学习计划（Markdown）
```

## 常用命令

### Go 代码（CS_Core_Courses）
```bash
cd CS_Core_Courses
go build ./...           # 编译验证
go run main.go           # 运行所有示例
```

### Python 代码（Math_Foundations）
```bash
cd Math_Foundations
pip install numpy matplotlib scipy sympy seaborn
python main.py           # 交互式菜单
python main.py --calculus        # 特定模块
```

### Python 代码（English_Learning）
```bash
cd English_Learning
python main.py           # 词汇学习工具
```

## 代码架构与组织模式

### 统一的模块结构
每个知识点模块都遵循 **理论 + 代码双轨** 的架构：

1. **guide.md**（理论文档）- 解释概念、原理、考点
2. **代码实现文件** - 可运行的代码示例

### Go 代码组织规范
每个 `.go` 文件遵循统一结构：
```
1. 数据结构定义（type XXX struct）
2. 构造函数（NewXXX）
3. 核心方法（增删改查等操作）
4. 示例函数（XXXExample）← 文件末尾，从这开始阅读！
```

### Python 代码组织规范
- 使用 `main.py` 作为入口点，提供交互式菜单
- 支持命令行参数直接运行特定模块
- 使用 matplotlib/seaborn 进行数据可视化

### 层次组织原则
- **基础 → 进阶**：每个科目都按难度分层
- **主题分组**：相关概念放在同一目录下
- **独立模块**：每个数据结构/概念独立成文件

## 文件命名规范

| 语言 | 规范 | 示例 |
|------|------|------|
| Go 模块 | kebab-case | `array.go`, `linkedlist.go`, `binarysearchtree.go` |
| Python 模块 | snake_case | `heap_sort.py`, `binary_search_tree.py` |
| Markdown 文档 | kebab-case | `reading_strategies_guide.md` |

## 编码约定

### 中文注释
- **所有代码必须包含详细的中文注释**
- 函数注释说明功能、参数、返回值
- 复杂逻辑添加行内注释

### 示例函数
- 每个模块文件末尾必须有 `XXXExample()` 函数
- 示例函数应展示该数据结构/算法的完整用法
- 这是用户阅读代码的入口点

### 考研导向
- 所有内容对应考研大纲考点
- 模块 README 标注对应考点
- 注重基础概念的理解而非工程优化

## Git 提交规范

遵循 Conventional Commits 格式（中文）：

```
<type>(<scope>): <subject>

[可选正文]

[可选页脚]
```

**Type**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`, `ci`, `revert`

**Scope 示例**: `math`, `data-structure`, `os`, `network`, `english`, `study-plan`

**Subject**: ≤50字，中文，动词开头，无句号
