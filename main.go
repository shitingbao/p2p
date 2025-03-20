package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// 服务器地址
const serverAddr = "47.99.104.79:8080"

// 客户端信息结构
type ClientInfo struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

func main() {
	// 连接到服务器
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		return
	}

	// 读取服务器发送的对方客户端信息
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)
	var peers []ClientInfo
	json.Unmarshal(buffer[:n], &peers)
	// 拿到信息就可以断开了
	conn.Close()

	// 选择另一个客户端进行连接
	var remote ClientInfo
	for _, p := range peers {
		if p.IP != conn.LocalAddr().(*net.TCPAddr).IP.String() {
			remote = p
			break
		}
	}

	fmt.Printf("Trying to connect to %s:%s\n", remote.IP, remote.Port)

	// 通过 UDP 进行 NAT 穿透
	udpAddr, _ := net.ResolveUDPAddr("udp", remote.IP+":"+remote.Port)
	udpConn, _ := net.DialUDP("udp", nil, udpAddr)
	defer udpConn.Close()

	// 发送测试包，尝试穿透 NAT
	for i := 0; i < 5; i++ {
		udpConn.Write([]byte("test"))
		fmt.Println("Sent test packet")
		time.Sleep(2 * time.Second)
	}

	// 监听 UDP 响应
	buf := make([]byte, 1024)
	n, _, err = udpConn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("No response, NAT blocking the connection")
	} else {
		fmt.Println("Received response from peer:", string(buf[:n]))
	}
}
