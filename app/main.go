package main

import (
	tgClient "app/clients/telegram"
	"app/config"
	"app/consumer/event-consumer"
	"app/events/telegram"
	"app/storage/files"

	"flag"
	"log"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}
	cfg := config.Get()

	eventsProcessor := telegram.New(
		tgClient.New(cfg.TgBotHost, cfg.Token),
		files.New(cfg.StoragePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, cfg.BatchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	var token *string = flag.String("tg-bot-token", "", "Token for access to telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
