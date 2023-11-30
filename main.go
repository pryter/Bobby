package main

import (
	"Bobby/cmd"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math"
	"os"
	"strings"
	"time"
)

func displayAppHeading() {
	title := fmt.Sprintf("Bobby artifact builder (dist version %s)", Configs.AppVersion)
	lineLength := 56
	log.Info().Msg(strings.Repeat("=", lineLength))
	padding := math.Floor(float64((lineLength - len(title)) / 2))
	log.Info().Msgf("%s%s", strings.Repeat(" ", int(padding)), title)
	log.Info().Msgf("%svisit https://bobby.pryter.me/ for more infos.", strings.Repeat(" ", int(padding)-1))
	log.Info().Msg(strings.Repeat("=", lineLength))
}

func main() {
	godotenv.Load()

	log.Logger = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC1123},
	).Level(zerolog.TraceLevel).With().Timestamp().Logger()

	displayAppHeading()

	log.Info().Msgf("Starting webhook api service on PORT %d", Configs.HTTPServices.Webhook.Port)
	cmd.StartWebhookService(Configs.HTTPServices.Webhook)

}
