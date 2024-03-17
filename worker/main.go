package main

import (
	"bobby-worker/cmd"
	"bobby-worker/internal/app/resources"
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
	fmt.Println(strings.Repeat("=", lineLength))
	padding := math.Floor(float64((lineLength - len(title)) / 2))
	fmt.Printf("%s%s\n", strings.Repeat(" ", int(padding)), title)
	fmt.Printf(
		"%svisit https://bobby.pryter.me/ for more infos.\n", strings.Repeat(" ", int(padding)-1),
	)
	fmt.Println(strings.Repeat("=", lineLength))
	fmt.Print("\n")
}

func main() {
	godotenv.Load()

	log.Logger = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC1123},
	).Level(zerolog.TraceLevel).With().Timestamp().Logger()

	displayAppHeading()

	// init resourceTree
	resources.InitResourceTree(Configs.AppResourcePath)

	//go cmd.StartServingArtifacts(Configs.HTTPServices.Artifacts)

	for {
		restart := cmd.StartWorkerService(Configs.HTTPServices.Worker)
		if !restart {
			break
		}
		time.Sleep(time.Second * 2)
	}
}
