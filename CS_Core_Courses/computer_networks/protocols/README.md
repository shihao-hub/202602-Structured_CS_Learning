# 网络协议 (Network Protocols)

## 408 考纲映射

本模块涵盖计算机网络 408 考试中常见应用层协议的核心知识点:

## 1. DNS (Domain Name System, 域名系统)

### 功能
- 将域名解析为 IP 地址
- 分布式数据库系统
- 应用层协议,使用 UDP/53 端口

### 记录类型
- **A 记录**: 域名 → IPv4 地址
- **AAAA 记录**: 域名 → IPv6 地址
- **CNAME 记录**: 域名别名 → 规范域名
- **MX 记录**: 邮件服务器地址
- **NS 记录**: 域的权威名称服务器

### 查询方式
- **递归查询**: 客户端 → 本地 DNS 服务器 (服务器负责完整解析)
- **迭代查询**: DNS 服务器之间 (返回下一步查询的服务器地址)

## 2. DHCP (Dynamic Host Configuration Protocol)

### 功能
- 动态分配 IP 地址
- 自动配置网络参数 (子网掩码、网关、DNS)
- 使用 UDP/67 (服务器) 和 UDP/68 (客户端)

### 工作过程 (408 重点)
1. **DHCP Discover**: 客户端广播发现消息
2. **DHCP Offer**: 服务器提供 IP 地址
3. **DHCP Request**: 客户端请求使用该 IP
4. **DHCP ACK**: 服务器确认分配

## 3. FTP (File Transfer Protocol)

### 特点
- 使用 TCP 连接
- **控制连接**: 端口 21 (持久连接)
- **数据连接**: 端口 20 (非持久连接)
- 支持主动模式和被动模式

## 4. SMTP/POP3 (电子邮件协议)

### SMTP (Simple Mail Transfer Protocol)
- 发送邮件协议
- 使用 TCP/25 端口
- 只能传输 7 位 ASCII 码

### POP3 (Post Office Protocol 3)
- 接收邮件协议
- 使用 TCP/110 端口
- 下载并删除模式

### IMAP (Internet Message Access Protocol)
- 更强大的邮件接收协议
- 使用 TCP/143 端口
- 支持在线管理邮件

## 文件说明

- `dns.go`: DNS 域名解析模拟
- `example.go`: 综合示例入口

## 运行示例

```go
package main

import "protocols"

func main() {
    protocols.RunAllProtocolExamples()
}
```

## 408 考点总结

- ✓ DNS 递归查询 vs 迭代查询
- ✓ DHCP 四步交互过程
- ✓ FTP 控制连接与数据连接分离
- ✓ SMTP 使用推 (Push) 模式, POP3 使用拉 (Pull) 模式
- ✓ 常用协议端口号记忆

## 常用端口号速记

| 协议  | 端口 | 传输层协议 |
|-------|------|-----------|
| HTTP  | 80   | TCP       |
| HTTPS | 443  | TCP       |
| FTP   | 21/20| TCP       |
| DNS   | 53   | UDP/TCP   |
| SMTP  | 25   | TCP       |
| POP3  | 110  | TCP       |
| IMAP  | 143  | TCP       |
| DHCP  | 67/68| UDP       |
