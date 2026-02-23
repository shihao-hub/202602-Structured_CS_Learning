# 算法模块 (Algorithm)

## 408考点对照

| 考点 | 文件 | 权重 |
|------|------|------|
| 排序算法（8种） | `sorting.go` | ★★★ |
| 查找算法（顺序/二分/分块） | `searching.go` | ★★★ |
| KMP字符串匹配 | `string_matching.go` | ★★★ |
| 动态规划 | `dynamic_programming.go` | ★★☆ |
| 贪心算法（Prim/Kruskal/哈夫曼） | `greedy.go` | ★★★ |
| 回溯法 | `backtracking.go` | ★★☆ |

## 文件说明

### sorting.go - 八大排序算法
- 冒泡排序、选择排序、插入排序、希尔排序
- 归并排序、快速排序、堆排序、基数排序
- 每种算法附带稳定性分析和复杂度统计

### searching.go - 查找算法
- 顺序查找（含哨兵优化）
- 二分查找（迭代版/递归版）
- 插值查找、分块查找

### dynamic_programming.go - 动态规划
- 0/1背包、完全背包
- LCS（最长公共子序列）、LIS（最长递增子序列）
- Floyd全源最短路径、矩阵链乘法

### greedy.go - 贪心算法
- 活动选择问题
- 哈夫曼编码（含编码树构建）
- Prim/Kruskal最小生成树（含并查集）

### backtracking.go - 回溯法
- N皇后问题（含棋盘可视化）
- 图着色问题
- 子集和问题、全排列

### string_matching.go - 字符串匹配
- 朴素模式匹配
- KMP算法（含next/nextval数组详解）
- 多位置匹配

## 运行方式

所有示例通过 `RunAllAlgorithmExamples()` 统一调用。
