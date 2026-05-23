package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"time"

	pb "main/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Подключение к серверу
	conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewConfigServiceClient(conn)

	serviceID := "auth-service"

	// 1. Проверка Health
	fmt.Println("\n=== Testing Health Check ===")
	healthResp, err := client.CheckHealth(context.Background(), &pb.HealthRequest{})
	if err != nil {
		log.Printf("Health check error: %v", err)
	} else {
		fmt.Printf("Health status: %s\n", healthResp.Status)
	}

	// 2. Получение первого конфига
	fmt.Println("\n=== Getting initial config ===")
	getResp, err := client.GetConfig(context.Background(), &pb.GetRequest{
		ServiceId: serviceID,
	})
	if err != nil {
		log.Printf("GetConfig error: %v", err)
	} else {
		fmt.Printf("Initial config: version=%d, data=%s\n",
			getResp.Version, string(getResp.Data))
	}

	// 3. Подписка на обновления (Watch) с текущей версией
	fmt.Println("\n=== Starting Watch (waiting for updates) ===")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	currentVersion := getResp.GetVersion()
	fmt.Printf("Subscribing with current_version=%d\n", currentVersion)

	stream, err := client.Watch(ctx, &pb.WatchRequest{
		ServiceId:      serviceID,
		CurrentVersion: currentVersion,
	})
	if err != nil {
		log.Fatalf("Watch error: %v", err)
	}

	// Канал для сигнала остановки
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// Читаем обновления
	go func() {
		for {
			config, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Stream ended")
				return
			}
			if err != nil {
				log.Printf("Stream receive error: %v", err)
				return
			}

			if config.Version == 0 {
				fmt.Printf("[Heartbeat] time=%s\n", time.Now().Format("15:04:05"))
			} else {
				fmt.Printf("\n[CONFIG UPDATE] version=%d, data=%s\n",
					config.Version, string(config.Data))
				fmt.Print("Waiting for next updates...\n")
			}
		}
	}()

	fmt.Println("\n========================================")
	fmt.Println("Client is running. Waiting for config updates...")
	fmt.Println("Now go and change the config (upload new schema/config)")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println("========================================")

	// Ждем сигнал завершения
	<-sigChan
	fmt.Println("\n\nShutting down...")
}
