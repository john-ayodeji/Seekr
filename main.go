package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/john-ayodeji/Seekr/handlers/sitemap_handler"
	"github.com/john-ayodeji/Seekr/internal"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

type ApiConfig struct {
	Port         int
	RabbitMQ_Url string
}

func main() {
	mux := http.NewServeMux()
	godotenv.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ExitSignal := make(chan os.Signal)

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Println(err)
	}

	rabbitMQ_URL := os.Getenv("RABBITMQ_URL")
	dbURL := os.Getenv("DB_URL")

	cfg := ApiConfig{
		Port:         port,
		RabbitMQ_Url: rabbitMQ_URL,
	}

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("OK")
	})
	RegisterRoutes(mux)

	//Rabbit MQ stuffs
	conn, ch, err := internal.ConnectRabbitMQ(cfg.RabbitMQ_Url)

	defer conn.Close()
	defer ch.Close()
	fmt.Println("RabbitMQ connection successful")

	rConfig := internal.RabbitCfg
	rConfig = &internal.RabbitConfig{
		URL:        rabbitMQ_URL,
		Connection: conn,
		Channel:    ch,
		Exchange:   "seekr_direct",
	}

	rConfig.CreateExchange()

	_, err1 := rConfig.DeclareAndBindQueue("html_parser.jobs", "page.fetch.success", rConfig.Exchange, true)
	if err1 != nil {
		fmt.Println(err1)
	}
	_, err2 := rConfig.DeclareAndBindQueue("html_fetcher.jobs", "url.fetch.success", rConfig.Exchange, true)
	if err2 != nil {
		fmt.Println(err1)
	}

	//RB config parsed to others
	sitemap_handler.RBcfg = *rConfig

	//Server shii
	addr := fmt.Sprintf("localhost:%d", cfg.Port)
	server := http.Server{Addr: addr, Handler: mux}
	fmt.Printf("Server started on port %v\nURL: https://%v\n", port, addr)
	if err := server.ListenAndServe(); err != nil {
		log.Println("Server failed to start")
		os.Exit(1)
	}

	//Shutdown signal and shii
	signal.Notify(ExitSignal, os.Interrupt)
	log.Println("Server shutting down")
	server.Shutdown(ctx)
	log.Println("Server has been shut down")
}
