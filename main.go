package main

import (
	"bob/core"
	"bob/handler"
	"bob/player"
	"fmt"
	"github.com/alexandrevicenzi/go-sse"
	"net/http"
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
		Logger: nil,
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

	p := player.NewPlayer(queue, environment, bobForwarder)

	bobHandler := handler.NewBobHandler(p, s)

	http.Handle("/api/v1/events", s)

	http.HandleFunc("/api/v1/play", bobHandler.HandlePlay)
	http.HandleFunc("/api/v1/pause", bobHandler.HandlePause)
	http.HandleFunc("/api/v1/playback", bobHandler.HandlePlayback)
	http.HandleFunc("/api/v1/search", bobHandler.HandleSearch)

	fmt.Println(http.ListenAndServe(":5002", nil))
}
