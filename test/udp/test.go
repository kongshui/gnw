package main

import (
	"log"
	"net"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":8888")
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP: %v", err)
	}
	defer conn.Close()

	log.Println("Relay server started on :8888")

	clients := make([]*net.UDPAddr, 0)

	buf := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		log.Printf("Received '%s' from %s", string(buf[:n]), clientAddr.String())

		// 保存客户端地址
		clients = append(clients, clientAddr)

		// 如果有两个客户端，交换地址
		if len(clients) == 2 {
			log.Printf("Exchanging addresses between %s and %s", clients[0].String(), clients[1].String())

			// 将客户端 1 的地址发送给客户端 2
			_, err = conn.WriteToUDP([]byte(clients[0].String()), clients[1])
			if err != nil {
				log.Printf("Error writing to UDP: %v", err)
			}

			// 将客户端 2 的地址发送给客户端 1
			_, err = conn.WriteToUDP([]byte(clients[1].String()), clients[0])
			if err != nil {
				log.Printf("Error writing to UDP: %v", err)
			}

			// 清空客户端列表
			clients = make([]*net.UDPAddr, 0)
		}
	}
}
