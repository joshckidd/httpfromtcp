package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

func main() {
	f, _ := net.Listen("tcp", "127.0.0.1:42069")
	defer f.Close()

	for {
		c, _ := f.Accept()
		fmt.Println("Connection accepted.")
		ch := getLinesChannel(c)

		for i := range ch {
			fmt.Println(i)
		}
		fmt.Println("Connection closed.")
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		var line string
		for {
			b := make([]byte, 8)
			_, err := f.Read(b)

			parts := bytes.Split(b, []byte("\n"))
			line = line + string(parts[0])
			if len(parts) > 1 {
				ch <- line
				line = string(parts[1])
			}
			if err != nil {
				ch <- line
				close(ch)
				break
			}
		}
	}()

	return ch
}
