package cmd

import (
	"Bobby/internal/events"
	"Bobby/internal/events/pushEvent"
	"errors"
	"fmt"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/rs/zerolog/log"
	"net/http"
)

// HTTPServiceConfig is general config struct for HTTP services
type HTTPServiceConfig struct {
	Port int    `mapstructure:"port"`
	Path string `mapstructre:"path"`
}

// StartWebhookService initiates all webhook services
func StartWebhookService(options HTTPServiceConfig) {

	hook, err := github.New(github.Options.Secret("bah"))

	if err != nil {
		panic(err)
	}

	// create concurrent pool for concurrent build tasks
	pool := events.InitConcurrentPool(events.ConcurrentPoolOptions{MaxConcurrentTasks: 2})

	http.HandleFunc(options.Path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.PushEvent)

		if err != nil {
			// event_not_found errors are negligible
			if errors.Is(err, github.ErrEventNotFound) {
				log.Error().Err(err)
			}

			log.Error().Err(err)
		}

		switch payload.(type) {

		// github's push event case
		case github.PushPayload:
			pushPayload := payload.(github.PushPayload)

			// add event to pool and let ConcurrentPool handle
			pool.Add(func() {
				pushEvent.WebhookPushEvent(pushPayload)
			})
		}
	})

	err = http.ListenAndServe(fmt.Sprintf(":%d", options.Port), nil)

	if err != nil {
		log.Error().Err(err).Msg("Unable to start http server.")
	}
}
