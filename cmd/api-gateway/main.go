package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	// Start all services
	services := []struct {
		name string
		path string
	}{
		{"auth", "services/auth/main.go"},
		{"pds", "services/pds/main.go"},
		{"bgs", "services/bgs/main.go"},
		{"api-gateway", "services/api-gateway/main.go"},
	}

	// Create a channel to receive errors
	errCh := make(chan error, len(services))

	// Start each service
	for _, service := range services {
		go func(name, path string) {
			log.Printf("Starting %s service...", name)
			cmd := exec.Command("go", "run", path)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Printf("Error running %s service: %v", name, err)
				errCh <- err
			}
		}(service.name, service.path)
	}

	// Wait for interrupt signal or error
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigCh:
		log.Println("Received interrupt signal, shutting down...")
	case err := <-errCh:
		log.Printf("Service error: %v, shutting down...", err)
	}
}
