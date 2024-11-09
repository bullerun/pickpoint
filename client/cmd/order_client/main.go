package main

import (
	"OzonHW1/client/internal/app"
	"OzonHW1/client/internal/config"
	"OzonHW1/client/internal/managers"
	order_service "OzonHW1/pkg/order-service/v1"
	"bufio"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const grpcServerHost = "localhost:7001"
const (
	topic             = "pvz.events-log"
	bootstrapServer   = "localhost:9092"
	consumerGroupName = "route256-consumer-group"
)

func main() {
	conn, err := grpc.NewClient(grpcServerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create grpc client: %v", err)
	}
	defer conn.Close()

	conf := config.NewConfig(bootstrapServer)
	fmt.Printf("%+v\n", conf)

	if err != nil {
		log.Fatal(err)
	}

	orderServiceClient := order_service.NewOrderServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)
	quit := make(chan bool, 1)
	cm := managers.NewCommandManager(orderServiceClient, reader, quit)
	app := app.MustLoad(reader, cm)

	go app.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit: // Завершение через команду Exit
		fmt.Println("Received exit command...")
	case sig := <-stop: // Завершение через системный сигнал
		fmt.Println("Received system signal:", sig)
	}
}
