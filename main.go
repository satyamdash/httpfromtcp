package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	str_chan := make(chan string)

	go func() {
		var curline string
		buf := make([]byte, 8)
		defer f.Close()
		for {
			n, err := f.Read(buf)
			if n > 0 {
				chunks := string(buf[:n])
				parts := strings.Split(chunks, "\n")
				for i := 0; i < len(parts)-1; i++ {
					str_chan <- curline + parts[i]
					curline = ""
				}

				curline += parts[len(parts)-1]
			}
			if err == io.EOF {
				// file ended; print whatever is left
				if len(curline) > 0 {
					str_chan <- curline
				}
				close(str_chan)
				break
			}
			if err != nil {
				break
			}

		}
	}()

	return str_chan
}

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println(err)
	}

	for line := range getLinesChannel(file) {
		fmt.Println("read:", line)
	}

}
