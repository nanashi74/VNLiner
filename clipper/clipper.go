package main

import "log"
import "time"
import "net"

import "github.com/atotto/clipboard"

func main() {

	prevText := ""
	for {
		time.Sleep(time.Second / 2)
		text, err := clipboard.ReadAll()
		if err != nil {
			log.Printf("Error reading clipboard: %v", err)
			continue
		}

		if text != prevText {
			prevText = text

			conn, err := net.Dial("tcp", "localhost:8081")
			if err != nil {
				log.Printf("Error connecting: %v", err)
			}

			conn.Write([]byte(text))
			conn.Close()
		}
	}
}
