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
	exit := make(chan string)
	defer conn.Close()
	defer close(msg)
	defer close(exit)
	go send(conn, msg)
	go receive(conn, msg, exit)
	select {
	case <-exit:
		fmt.Printf("%s退出\n",conn.RemoteAddr().String())
		return
	}
}

func send(conn net.Conn, message chan string) {
	for {
		msg := <-message
		_, err := conn.Write([]byte("reviewed" + msg))
		if err != nil {
			fmt.Printf("write err %+v", err)
			return
		}
	}
}

func receive(conn net.Conn, message chan string, exit chan string) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Printf("message: %s\n", scanner.Text())
		message <- scanner.Text()
	}
	exit <- "q"
}
