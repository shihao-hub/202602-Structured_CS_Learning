# 线性数据结构进阶 (Linear)

## 408考点对照

| 考点 | 文件 | 权重 |
|------|------|------|
| 串的定义与基本操作 | `string_adt.go` | ★★☆ |
| 稀疏矩阵三元组表示与转置 | `sparse_matrix.go` | ★★☆ |

## 文件说明

### string_adt.go - 串
- 串的顺序存储结构
- 基本操作：求子串、连接、定位、替换、插入、删除、比较
- 注：KMP模式匹配见 `algorithm/string_matching.go`

### sparse_matrix.go - 稀疏矩阵
- 三元组顺序表存储
- 简单转置法 O(Cols*NumTerms)
- 快速转置法 O(Cols+NumTerms)
- 稀疏矩阵加法
