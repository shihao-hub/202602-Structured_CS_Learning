package network

import "fmt"

func RunAllNetworkExamples() {
	fmt.Println("\n╔══════════════════════════════════════╗")
	fmt.Println("║      计算机网络 - 网络层模块         ║")
	fmt.Println("╚══════════════════════════════════════╝")
	IPExample()
	RoutingExample()
	ARPExample()
}
