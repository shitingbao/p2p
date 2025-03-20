package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

// 客户端信息
type ClientInfo struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

var clients = make(map[string]ClientInfo)

func handleConnection(conn net.Conn) {
	log.Println("connect ", conn.RemoteAddr().String())
	ctx := context.Background()
	defer conn.Close()
	defer ctx.Done()

	go func() {
		// 获取客户端 IP 和端口
		addr := conn.RemoteAddr().(*net.TCPAddr)
		client := ClientInfo{IP: addr.IP.String(), Port: fmt.Sprintf("%d", addr.Port)}

		// 存储客户端信息
		clients[addr.String()] = client
		fmt.Printf("Client connected: %s\n", addr.String())

		// 等待两个客户端连接
		if len(clients) == 2 {
			clientList := []ClientInfo{}
			for _, c := range clients {
				clientList = append(clientList, c)
			}

			// 发送对方的 IP/端口信息
			for _, c := range clients {
				conn, err := net.Dial("tcp", c.IP+":"+c.Port)
				if err != nil {
					fmt.Println("Failed to send client info:", err)
					continue
				}
				data, _ := json.Marshal(clientList)
				conn.Write(data)
			}
		}
	}()

	r := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			log.Println("read done ")
			return
		default:
			res, err := r.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					log.Println("read eof ")
					return
				}
				log.Println("read err ", err)
				return
			}

			log.Println("read msg ", res)
		}
	}

}

func main() {
	listenPort := ":8080"
	listener, err := net.Listen("tcp", listenPort)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	fmt.Println("Server started on port ", listenPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn)
	}
}
