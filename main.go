package main

import (
	"Test_task_Halolab/cmd/api"
	"context"
	"fmt"
	"os/signal"
	"syscall"
)

// @title           Underwater sensor data aggregator service
// @version         1.0
// @description     A sensor data management service API in Go.

// @host      localhost:8080
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := api.NewServer()

	err := server.StartServer(ctx)
	if err != nil {
		fmt.Println("failed to start app:", err)
	}
}
