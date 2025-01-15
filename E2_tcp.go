package main

import (
	"fmt"
	"net"
	//"bufio"
	"time"
)

var port1 string = ":34933"

func handleClient(conn net.Conn) {
	defer conn.Close()
}

func listen(network string, port string) {
	tcp_addr,_ := net.ResolveTCPAddr(network, port)
	conn, _ := net.ListenTCP(network, tcp_addr)
	
	defer conn.Close()

	buffer := make([]byte, 1024)
	
	n, address, _ := conn.Read(buffer)
	fmt.Println(address.String())
	fmt.Println("the message is: ", string(buffer[:n]))
	//conn2, _ := net.DialTCP(network, nil, tcp_addr)

	
	//conn2, _ := conn.Accept()
	//go handleClient(conn2)
	//defer conn2.Close()
	// buffer := make([]byte, 1024)
	
	// n, address, _ := conn.ReadFrom(buffer)
	// fmt.Println(address.String())
	// fmt.Println("the message is: ", string(buffer[:n]))
	
	//time.Sleep(1*(time.Second))
}

func main() {
	listen("tcp", port1)
}
