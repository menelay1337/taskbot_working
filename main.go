package main

import (
	"log"

	tgClient "taskbot1/clients/telegram"
	"taskbot1/config"
	"taskbot1/consumer/event-consumer"
	"taskbot1/events/telegram"
	"taskbot1/storage/mysql"
)

const (
	tgBotHost   = "api.telegram.org"
	batchSize   = 100
)

func main() {
	cfg := config.MustLoad()
	storage, err := mysql.New(cfg.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	storage.Init()

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, cfg.Token),
		storage,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
