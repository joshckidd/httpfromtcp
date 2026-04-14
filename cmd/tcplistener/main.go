package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"net"
)

func main() {
	f, _ := net.Listen("tcp", "127.0.0.1:42069")
	defer f.Close()

	for {
		c, _ := f.Accept()
		fmt.Println("Connection accepted.")

		req, err := request.RequestFromReader(c)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Request line:")
			fmt.Printf("- Method: %s\n", req.RequestLine.Method)
			fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
			fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		}

		fmt.Println("Connection closed.")
	}

}
