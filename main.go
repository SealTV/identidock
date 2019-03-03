package main

import (
	"identidock/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis"
)

func main() {
	opt := redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	}

	r := redis.NewClient(&opt)
	pong, err := r.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(pong)

	logger := log.New(os.Stdout, "", 0)
	h := server.NewServer(r)
	go func() {
		http.ListenAndServe(":5000", h)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGKILL)

	<-stop
	logger.Println("Server gracefully stopped")
	os.Exit(0)
}
