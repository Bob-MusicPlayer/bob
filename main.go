package main

import (
	"bob/core"
	"bob/player"
	"fmt"
	"github.com/alexandrevicenzi/go-sse"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	environment := &core.Environment{}

	s := sse.NewServer(&sse.Options{
		RetryInterval: 10 * 1000,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, OPTIONS",
			"Access-Control-Allow-Headers": "Keep-Alive,X-Requested-With,Cache-Control,Content-Type,Last-Event-ID",
		},
		ChannelNameFunc: func(request *http.Request) string {
			return request.URL.Path
		},
		Logger: log.New(os.Stdout, "go-sse: ", log.Ldate|log.Ltime|log.Lshortfile),
	})
	defer s.Shutdown()

	configManager := core.NewConfigManager(environment)
	err := configManager.ReadConfig()
	if err != nil {
		panic(err)
	}

	environment.ConfigManager = configManager

	bobForwarder := player.NewBoxForwarder(environment)

	queue := player.NewQueue()

	player.NewPlayer(queue, environment, bobForwarder)

	go func() {
		for {
			time.Sleep(time.Second)
			s.SendMessage("/api/v1/events", sse.SimpleMessage("test"))
		}
	}()

	//bobHandler := handler.NewBobHandler(s)

	http.Handle("/api/v1/events", s)

	fmt.Println(http.ListenAndServe(":5002", nil))
}
