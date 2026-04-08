package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	u, _ := net.ResolveUDPAddr("udp", "localhost:42069")

	c, _ := net.DialUDP("udp", nil, u)
	defer c.Close()

	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		l, _ := r.ReadString('\n')
		_, err := c.Write([]byte(l))
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
