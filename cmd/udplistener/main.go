package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		fmt.Println(err)
	}
	conn, _ := net.DialUDP("udp", nil, addr)

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		if _, err := conn.Write([]byte(input)); err != nil {
			fmt.Println(err)
		}
	}

}
