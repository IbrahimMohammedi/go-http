package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		fmt.Println("Error opening tcp:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error listening to conn:", err)
			continue
		}

		fmt.Println("Connection accepted")

		ch := getLinesChannel(conn)

		for line := range ch {
			fmt.Println(line)
		}

		fmt.Println("Connection closed")
		conn.Close()
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	messages := make(chan string)

	go func() {
		defer close(messages)

		data := make([]byte, 8)
		var leftover string

		for {
			n, err := f.Read(data)

			if n > 0 {
				currentChunk := leftover + string(data[:n])
				parts := strings.Split(currentChunk, "\n")

				for i := 0; i < len(parts)-1; i++ {
					messages <- parts[i]
				}

				leftover = parts[len(parts)-1]
			}

			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Error reading from connection:", err)
				break
			}
		}

		if leftover != "" {
			messages <- leftover
		}
	}()

	return messages
}
