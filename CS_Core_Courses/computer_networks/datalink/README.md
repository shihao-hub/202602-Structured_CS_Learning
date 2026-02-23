# 数据链路层 (Data Link Layer)

## 408 考纲映射

本模块涵盖计算机网络 408 考试中数据链路层的核心知识点:

### 1. 以太网协议
- 以太网帧格式 (前导码、目的/源 MAC 地址、类型、数据、FCS)
- CSMA/CD 协议 (载波侦听多路访问/冲突检测)
- 冲突检测与退避算法

### 2. 差错检测
- 循环冗余校验 (CRC)
- 海明码 (Hamming Code)
- 单比特差错检测与纠正

### 3. 可靠传输协议
- 停止-等待协议 (Stop-and-Wait)
- 回退 N 帧 (Go-Back-N, GBN)
- 选择重传 (Selective Repeat, SR)
- 滑动窗口机制

## 文件说明

- `ethernet.go`: 以太网协议与 CSMA/CD 实现
- `error_detection.go`: CRC 和海明码差错检测
- `sliding_window.go`: 滑动窗口协议实现
- `example.go`: 综合示例入口

## 运行示例

```go
package main

import "datalink"

func main() {
    datalink.RunAllDatalinkExamples()
}
```

## 学习建议

1. 理解以太网帧结构和 CSMA/CD 工作原理
2. 掌握 CRC 和海明码的计算方法
3. 区分三种滑动窗口协议的特点
4. 重点关注 GBN 和 SR 的窗口大小与序号范围关系
5. 结合 408 真题练习帧格式和协议计算题
