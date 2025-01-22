package main

import (
	"fmt"
	"net"
	"time"
	//"bufio"
)

var port1TCP string = "10.100.23.204:33546"

func handleClient(conn net.Conn) {
	defer conn.Close()
}



func listen2(network string, port string) {
	tcp_addr, err := net.ResolveTCPAddr(network, port)
	if err != nil {
		fmt.Println("error")
	}
	conn, err := net.Dial(network, tcp_addr.String())
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("error")
	}
	fmt.Println("the message is: ", string(buffer[:n]))
	
	time.Sleep(1*(time.Second))
}

func writeTCP(network string, adr_port string) {
	tcp_addr, err := net.ResolveTCPAddr(network, adr_port)
	if err != nil {
		fmt.Println("error")
	}
	conn, err := net.Dial(network, tcp_addr.String())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	_ , err1 := conn.Write([]byte("Kommer meld"))

	if err1 != nil {
		panic(err1)
	}

}

func main() {
	
	go writeTCP("tcp", port1TCP)
	listen2("tcp", port1TCP)
	
}
