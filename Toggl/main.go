package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"toggl/app"
	"toggl/app/config"
)

func main() {
	// Load configuration
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	// Create a new instance of the app
	app, err := app.NewApp(config)
	if err != nil {
		log.Fatalf("failed to create app: %s", err)
	}

	// Start the app
	go func() {
		if err := app.Start(); err != nil {
			log.Fatalf("failed to start app: %s", err)
		}
	}()

	// Wait for SIGINT or SIGTERM signal to gracefully stop the app
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Received interrupt signal. Stopping app...")

	// Stop the app
	if err := app.Stop(); err != nil {
		log.Fatalf("failed to stop app: %s", err)
	}

	log.Println("App stopped gracefully.")
}
