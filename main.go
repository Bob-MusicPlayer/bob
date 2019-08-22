package main

import (
	"bob/core"
	"bob/handler"
	"bob/player"
	"bob/util"
	"fmt"
	"net/http"
)

func main() {
	environment := &core.Environment{}

	hub := util.NewHub()
	go hub.Run()

	configManager := core.NewConfigManager(environment)
	err := configManager.ReadConfig()
	if err != nil {
		panic(err)
	}

	environment.ConfigManager = configManager

	bobForwarder := player.NewBoxForwarder(environment)

	queue := player.NewQueue()

	p := player.NewPlayer(queue, environment, bobForwarder)

	worker := util.NewWorker(hub, p, environment)
	bobHandler := handler.NewBobHandler(hub, worker)

	http.HandleFunc("/api/v1/connect", bobHandler.HandleConnect)

	fmt.Println(http.ListenAndServe(":5002", nil))
}
