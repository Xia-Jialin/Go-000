package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("hello~")
	listen, err := net.Listen("tcp", "127.0.0.1:10000")
	if err != nil {
		log.Fatalf("listen error: %v\n", err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("accept error: %v\n", err)
			continue
		}
		// 开始goroutine监听连接
		go handleConn(conn)
	}
}

// 处理一个链接
func handleConn(conn net.Conn) {
	msg := make(chan string, 1)
	defer conn.Close()
	defer close(msg)
	go send(conn, msg)
	go receive(conn, msg)
	select {}
}

func send(conn net.Conn, message <-chan string) {
	for msg := range message {
		fmt.Println(conn, msg)
		_, err := conn.Write([]byte("reviewed" + msg))
		if err != nil {
			fmt.Println("write err %+v", err)
			return
		}
	}
}

func receive(conn net.Conn, message chan string) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Printf("message: %s\n", scanner.Text())
		message <- scanner.Text()
	}
}
