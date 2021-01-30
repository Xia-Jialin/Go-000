package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	//打开连接:
	conn, err := net.Dial("tcp", "127.0.0.1:10000")
	if err != nil {
		//由于目标计算机积极拒绝而无法创建连接
		fmt.Println("Error dialing", err.Error())
		return // 终止程序
	}
	defer conn.Close()

	inputReader := bufio.NewReader(os.Stdin)
	wr := bufio.NewWriter(conn)

	fmt.Println("First, what is your name?")
	clientName, _ := inputReader.ReadString('\n')
	trimmedClient := strings.Trim(clientName, "\r\n") // Windows 平台下用 "\r\n"，Linux平台下使用 "\n"
	go receive(conn)
	// 给服务器发送信息直到程序退出：
	for {
		fmt.Println("What to send to the server? Type Q to quit.")
		input, _ := inputReader.ReadString('\n')
		trimmedInput := strings.Trim(input, "\r\n")
		if trimmedInput == "Q" {
			return
		}
		wr.Write([]byte(trimmedClient + " says: " + trimmedInput + "\n"))
		wr.Flush()
	}
}

func receive(conn net.Conn) {
	for {
		// 缓冲区
		buf := make([]byte, 1024)
		//接受服务器回发的数据
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("conn.Read err:%v\n", err)
			return
		}

		fmt.Printf("服务器回发的数据为:%v\n", string(buf[:n]))
	}
}
