package main

import (
	"bob/core"
	"bob/handler"
	"bob/player"
	"fmt"
	"github.com/alexandrevicenzi/go-sse"
	"net"
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

	//playerRepository := repository.NewPlayerRepository(playerDatabase)

	bobHandler := handler.NewBobHandler(p, s)

	http.Handle("/api/v1/events", s)

	http.HandleFunc("/api/v1/play", bobHandler.HandlePlay)
	http.HandleFunc("/api/v1/pause", bobHandler.HandlePause)
	http.HandleFunc("/api/v1/playback", bobHandler.HandlePlayback)
	http.HandleFunc("/api/v1/playback/seek", bobHandler.HandlePlaybackSeek)
	http.HandleFunc("/api/v1/search", bobHandler.HandleSearch)
	http.HandleFunc("/api/v1/sync", bobHandler.HandleSync)

	l, err := net.Listen("tcp4", "192.168.11.241:5002")
	if err != nil {
		panic(err)
	}
	fmt.Println(http.Serve(l, nil))
}
