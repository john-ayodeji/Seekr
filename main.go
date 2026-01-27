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
	"github.com/john-ayodeji/Seekr/internal/database"
	"github.com/john-ayodeji/Seekr/services/crawler"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	mux := http.NewServeMux()

	//parse env variables
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
	Workers, err := strconv.Atoi(os.Getenv("WORKERS"))
	if err != nil {
		log.Println(err)
	}

	//connect DB and stuffs
	db, err3 := sql.Open("postgres", dbURL)
	if err3 != nil {
		fmt.Println(err3)
		return
	}
	if err := db.Ping(); err != nil {
		fmt.Printf("db ping failed: %v\n", err)
		return
	}
	fmt.Println("db connection successful")

	//parse ApiConfig - initialize global config
	internal.Cfg = &internal.ApiConfig{
		Port:         port,
		RabbitMQ_Url: rabbitMQ_URL,
		Db:           database.New(db),
	}

	//health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("OK")
	})

	//register routes
	RegisterRoutes(mux)

	//Rabbit MQ stuffs
	conn, ch, err := internal.ConnectRabbitMQ(internal.Cfg.RabbitMQ_Url)

	defer conn.Close()
	defer ch.Close()
	fmt.Println("RabbitMQ connection successful")

	// Initialize global Rabbit config
	internal.RabbitCfg = &internal.RabbitConfig{
		URL:        rabbitMQ_URL,
		Connection: conn,
		Channel:    ch,
		Exchange:   "seekr_direct",
	}

	internal.RabbitCfg.CreateExchange()

	_, err1 := internal.RabbitCfg.DeclareAndBindQueue("html_parser.jobs", "page.fetch.success", internal.RabbitCfg.Exchange, true)
	if err1 != nil {
		fmt.Println(err1)
	}
	_, err2 := internal.RabbitCfg.DeclareAndBindQueue("html_fetcher.jobs", "url.fetch.success", internal.RabbitCfg.Exchange, true)
	if err2 != nil {
		fmt.Println(err2)
	}

	_, err4 := internal.RabbitCfg.DeclareAndBindQueue("normalize_text.jobs", "parse.html.success", internal.RabbitCfg.Exchange, true)
	if err4 != nil {
		fmt.Println(err4)
	}

	//RB config parsed to others
	sitemap_handler.RBcfg = *internal.RabbitCfg

	//workers and shii
	for i := 0; i < Workers; i++ {
		go crawler.ProcessHTML(conn)
		go crawler.ProcessParseHTML(conn)
	}

	//Server shii
	addr := fmt.Sprintf("localhost:%d", internal.Cfg.Port)
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
