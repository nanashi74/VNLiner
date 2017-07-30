package main

import "fmt"
import "log"
import "net"
import "net/http"
import "bufio"
import "io/ioutil"

import "github.com/desertbit/glue"

func main() {

	inputSocket, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	defer inputSocket.Close()

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("public"))))
	http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("dist"))))

	server := glue.NewServer(glue.Options{
		HTTPListenAddress: ":8080",
	})
	defer server.Release()

	server.OnNewSocket(onNewSocket)

	go server.Run()

	for {
		conn, err := inputSocket.Accept()
		if err != nil {
			log.Printf("Connection Error: %v", err)
		}

		go handleMessages(&conn, server)
	}

}

func onNewSocket(s *glue.Socket) {
	s.OnClose(func() {
		log.Printf("socket closed with remote address: %s", s.RemoteAddr())
	})

	s.OnRead(func(data string) {
		fmt.Printf("Received message")
	})

	log.Printf("Socket connected: %s", s.RemoteAddr())
}

func handleMessages(conn *net.Conn, server *glue.Server) {

	conts, err := ioutil.ReadAll(bufio.NewReader(*conn))
	if err != nil {
		log.Printf("Read Error: %v", err)
	}

	sockets := server.Sockets()
	for _, socket := range sockets {
		if socket.IsInitialized() {
			socket.Write(string(conts))
		}
	}
}
