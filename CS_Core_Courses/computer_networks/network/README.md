# 网络层 (Network Layer)

## 408 考纲映射

本模块涵盖计算机网络 408 考试中网络层的核心知识点:

### 1. IP 协议
- IPv4 地址与子网划分
- IP 数据报格式
- IP 分片与重组
- 子网掩码与 CIDR

### 2. 路由算法
- 距离向量路由算法 (RIP - Bellman-Ford)
- 链路状态路由算法 (OSPF - Dijkstra)
- 路由表的建立与更新

### 3. ARP 协议
- 地址解析协议 (ARP)
- ARP 缓存机制
- ARP 请求与应答过程

## 文件说明

- `ip.go`: IP 协议实现,包括子网计算和分片模拟
- `routing.go`: 路由算法实现 (距离向量和链路状态)
- `arp.go`: ARP 协议模拟
- `example.go`: 综合示例入口

## 运行示例

```go
package main

import "network"

func main() {
    network.RunAllNetworkExamples()
}
```

## 学习建议

1. 先理解 IP 地址结构和子网划分
2. 掌握两种主要路由算法的区别
3. 理解 ARP 在局域网中的作用
4. 结合 408 真题进行练习
