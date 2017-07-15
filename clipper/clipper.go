package main

import "log"
import "time"
import "github.com/atotto/clipboard"
import "github.com/go-redis/redis"

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Could not ping the redis client: %v", err)
	}
	defer client.Close()

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
			err = client.Publish("vnlines", text).Err()
			if err != nil {
				log.Printf("Error publishing to redis")
			}
		}
	}
}
