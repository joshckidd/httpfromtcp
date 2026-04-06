package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	f, _ := os.Open("messages.txt")
	var line string
	for {
		b := make([]byte, 8)
		_, err := f.Read(b)

		parts := bytes.Split(b, []byte("\n"))
		line = line + string(parts[0])
		if len(parts) > 1 {
			fmt.Printf("read: %s\n", line)
			line = string(parts[1])
		}
		if err != nil {
			fmt.Printf("read: %s\n", line)
			os.Exit(0)
		}
	}
}
