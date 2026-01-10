package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	m, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	defer m.Close()
	ch := getLinesChannel(m)
	for c := range ch {
		fmt.Printf("read: %s\n", c)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	messages := make(chan string)

	data := make([]byte, 8)
	var leftover string
	go func() {
		for {
			n, err := f.Read(data)

			if n > 0 {
				current_chunk := leftover + string(data[:n])

				parts := strings.Split(current_chunk, "\n")

				for i := 0; i < len(parts)-1; i++ {
					messages <- parts[i]
				}

				leftover = parts[len(parts)-1]
			}

			if err == io.EOF {
				break
			}

			if err != nil {
				fmt.Println("Error while reading the txt file: ", err)
				break
			}
		}
		if leftover != "" {
			messages <- leftover
		}
		close(messages)
	}()

	return messages
}
