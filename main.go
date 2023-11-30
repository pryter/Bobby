package main

import (
	"Bobby/cmd"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func main() {
	godotenv.Load()

	log.Logger = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC1123},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()

	log.Info().Msgf("==================================================")
	log.Info().Msgf("   Bobby artifact builder (dist version %s)", Configs.AppVersion)
	log.Info().Msgf("  visit https://bobby.pryter.me/ for more infos.")
	log.Info().Msgf("==================================================")

	log.Info().Msgf("Starting webhook api service on PORT %d", Configs.HTTPServices.Webhook.Port)
	cmd.StartWebhookService(Configs.HTTPServices.Webhook)

}
