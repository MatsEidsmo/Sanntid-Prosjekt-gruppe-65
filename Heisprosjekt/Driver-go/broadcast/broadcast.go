package broadcast


import (
	"fmt"
	"net"
	//"bufio"
	"time"
)

var Port1 string = ":30000"
var Port2 string = ":20020"

func ListenUDP(network string, port string) {
	conn, _ := net.ListenPacket(network, port)

	defer conn.Close()

	buffer := make([]byte, 1024)
	
	n, address, _ := conn.ReadFrom(buffer)
	fmt.Println(address.String())
	fmt.Println("the message is: ", string(buffer[:n]))
	
	time.Sleep(1*(time.Second))
}

func WriteUDP(network string, port string) {
	udp_addr,_ := net.ResolveUDPAddr(network, port)

	conn, _ := net.DialUDP(network, nil, udp_addr)
	defer conn.Close()
	
	_ , err := conn.Write([]byte("Hallo, vi er gruppe 65\n"))

	if err != nil {
		panic(err)
	}
}



