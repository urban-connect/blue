package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tinygo.org/x/bluetooth"
	"urban-connect.ch/blue/api"
	"urban-connect.ch/blue/detection"
)

func main() {
	quit := make(chan bool)
	detectionChannel := make(chan detection.Device)

	store := detection.NewStore()
	server := api.NewServer(store, detectionChannel)

	if err := bluetooth.DefaultAdapter.Enable(); err != nil {
		fmt.Printf("Failed to enable BLE adapter: %v\n", err)
		return
	}

	go func() {
		fmt.Println("Press Ctrl + C to exit...")

		exit := make(chan os.Signal, 3)
		signal.Notify(exit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-exit

		quit <- true
	}()

	go func() {
		fmt.Println("Starting server on 127.0.0.1:8080...")

		httpServer := &http.Server{
			Addr:              "127.0.0.1:8080",
			Handler:           server,
			ReadHeaderTimeout: 10 * time.Second,
			IdleTimeout:       120 * time.Second,
		}

		if err := httpServer.ListenAndServe(); err != nil {
			fmt.Printf("Failed to start the server: %v\n", err)
		}

		quit <- true
	}()

	if err := store.Watch(detectionChannel, quit); err != nil {
		fmt.Printf("Failed to start the store watcher: %v\n", err)
	}
}
