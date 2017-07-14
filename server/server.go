package main

import "fmt"
import "log"
import "net/http"

import "github.com/go-redis/redis"
import "github.com/desertbit/glue"

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Could not ping the redis server.")
	}

	pubsub := client.Subscribe("vnlines")
	defer pubsub.Close()

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("public"))))
	http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("dist"))))

	server := glue.NewServer(glue.Options{
		HTTPListenAddress: ":8080",
	})
	defer server.Release()

	server.OnNewSocket(onNewSocket)

	msgs := pubsub.Channel()
	go handleMessages(msgs, server)

	err = server.Run()
	if err != nil {
		log.Fatalf("Glue Run: %v", err)
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

func handleMessages(msgs <-chan *redis.Message, server *glue.Server) {
	for msg := range msgs {
		sockets := server.Sockets()
		for _, socket := range sockets {
			if socket.IsInitialized() {
				socket.Write(msg.Payload)
			}
		}
	}
}
