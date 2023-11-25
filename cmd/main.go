package main

import (
	"Bobby/cmd/service"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	service.StartWebhookService()
}
