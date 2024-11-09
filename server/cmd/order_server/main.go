package main

import (
	desc "OzonHW1/pkg/order-service/v1"
	"OzonHW1/server/internal/config"
	"OzonHW1/server/internal/infra/kafka/producer"
	"OzonHW1/server/internal/order_service"
	"OzonHW1/server/internal/storage"
	"OzonHW1/server/internal/storage/postgres"
	"OzonHW1/server/internal/tracer"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/loads"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	psqlDSN     = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	grpcHOST    = "localhost:7001"
	httpHOST    = "localhost:7002"
	adminHost   = "localhost:7003"
	indexFile   = "./static/index.html"
	swaggerFile = "./pkg/order-service/v1/order_service.swagger.json"
)
const topic = "pvz.events-log"
const kafkaBroker = "localhost:9092"
const cacheCapacity = 500

func main() {
	ctx := context.Background()
	closer, err := tracer.MustLoad("orderService")
	if err != nil {
		log.Fatal(err)
	}
	defer closer()
	pool, err := pgxpool.New(ctx, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	storageFacade := newStorageFacade(pool)
	prod, err := newSyncProducer()
	if err != nil {
		log.Fatal(err)
	}

	defer prod.Close()

	orderService := order_service.NewImplementation(storageFacade, prod, topic)

	lis, err := net.Listen("tcp", grpcHOST)
	if err != nil {
		log.Fatalf("failed to listen, %v", err)
	}
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor())
	reflection.Register(grpcServer)
	desc.RegisterOrderServiceServer(grpcServer, orderService)
	mux := runtime.NewServeMux()

	err = desc.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, grpcHOST, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})

	if err != nil {
		log.Fatalf("failed to register order service handler: %v", err)
	}

	go func() {
		if err := http.ListenAndServe(httpHOST, mux); err != nil {
			log.Fatalf("failed to listen and serve http: %v", err)
		}
	}()

	go func() {
		specDoc, err := loads.Spec(swaggerFile)
		if err != nil {
			log.Fatalln("Failed to load spec:", err)
		}
		log.Println("Would be serving:", specDoc.Spec().Info.Title)
		adminServer := chi.NewMux()
		adminServer.Handle("/metrics", promhttp.Handler())
		adminServer.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			jsonSpec, err := json.Marshal(specDoc.Spec())
			if err != nil {
				http.Error(w, "Failed to marshal spec", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(jsonSpec)
		})
		adminServer.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if _, err := os.Stat(indexFile); os.IsNotExist(err) {
				http.Error(w, "Index file not found", http.StatusNotFound)
				return
			}

			http.ServeFile(w, r, indexFile)
		})
		if err := http.ListenAndServe(adminHost, adminServer); err != nil {
			log.Fatalf("Failed to listen and serve admin server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve, %v", err)
		}
	}()
	sig := <-stop
	fmt.Println("Received system signal:", sig)
}

func newStorageFacade(pool *pgxpool.Pool) postgres.OrderService {
	txManager := postgres.NewTxManager(pool)
	pgRepository := postgres.NewPgRepository(txManager, cacheCapacity)
	return storage.NewStorageFacade(txManager, pgRepository)
}

func newSyncProducer() (sarama.SyncProducer, error) {
	conf := config.NewConfig(topic, kafkaBroker)
	fmt.Printf("%+v\n", conf)
	return producer.NewSyncProducer(conf.Kafka,
		producer.WithIdempotent(),
		producer.WithRequiredAcks(sarama.WaitForAll),
		producer.WithMaxOpenRequests(1),
		producer.WithMaxRetries(5),
		producer.WithRetryBackoff(10*time.Millisecond),
		producer.WithProducerPartitioner(sarama.NewHashPartitioner),
	)

}
