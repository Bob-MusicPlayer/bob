package main

import (
	"bob/handler"
	"bob/util"
	"fmt"
	"net/http"
)

func main() {
	hub := util.NewHub()
	go hub.Run()

	worker := util.NewWorker(hub)

	go func() {
		for {
			select {
			case <-hub.OnRegister:
				fmt.Println("New Client connected")
				fmt.Println(fmt.Sprintf("%d Clients connected", len(hub.Clients)))
			case <-hub.OnUnregister:
				fmt.Println("Client leaved")
				fmt.Println(fmt.Sprintf("%d Clients connected", len(hub.Clients)))
			}

		}
	}()

	bobHandler := handler.NewBobHandler(hub, worker)

	http.HandleFunc("/api/v1/connect", bobHandler.HandleConnect)

	fmt.Println(http.ListenAndServe(":5002", nil))
}
