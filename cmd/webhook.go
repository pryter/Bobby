package cmd

import (
	"Bobby/internal/events"
	"Bobby/internal/events/pushEvent"
	"errors"
	"fmt"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"net/http"
)

type HTTPServiceConfig struct {
	Port int    `mapstructure:"port"`
	Path string `mapstructre:"path"`
}

func StartWebhookService(options HTTPServiceConfig) {
	godotenv.Load()

	hook, err := github.New(github.Options.Secret("bah"))

	if err != nil {
		panic(err)
	}

	pool := events.InitConcurrentPool(events.ConcurrentPoolOptions{MaxConcurrentTasks: 2})

	http.HandleFunc(options.Path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.PushEvent)
		if err != nil {
			if errors.Is(err, github.ErrEventNotFound) {
				log.Error().Err(err)
			}
		}

		switch payload.(type) {
		case github.PushPayload:
			pushPayload := payload.(github.PushPayload)
			pool.Add(func() {
				pushEvent.WebhookPushEvent(pushPayload)
			})
		}
	})

	err = http.ListenAndServe(fmt.Sprintf(":%d", options.Port), nil)

	if err != nil {
		log.Fatal().Err(err)
		return
	}
}
