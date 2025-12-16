package main

import (
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	// 创建一个 UDP 地址
	serverAddr, err := net.ResolveUDPAddr("udp", "43.143.242.172:8888")
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	// 使用 ListenUDP 创建未绑定远程地址的连接
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		log.Fatalf("Failed to listen on UDP: %v", err)
	}
	defer conn.Close()

	// 向中继服务器发送注册消息
	_, err = conn.WriteToUDP([]byte("register"), serverAddr)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	// 接收对方的地址
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	peerAddrStr := string(buf[:n])
	log.Printf("Received peer address: %s", peerAddrStr)

	peerAddr, err := net.ResolveUDPAddr("udp", strings.TrimSpace(peerAddrStr))
	if err != nil {
		log.Fatalf("Failed to resolve peer address: %v", err)
	}
	log.Println(peerAddr)
	// 尝试与对方直接通信
	for range 5 {
		_, err = conn.WriteToUDP([]byte("hello from client"), peerAddr)
		if err != nil {
			log.Printf("Failed to send message to peer: %v", err)
		}

		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Received message from peer: %s", string(buf[:n]))
			break
		}
		log.Printf("Received message from peer: %s", string(buf[:n]))
		time.Sleep(2 * time.Second)
	}
}
