package cmd

import (
	"Bobby/internal/worker"
	"errors"
	"fmt"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/rs/zerolog/log"
	"net/http"
)

// HTTPServiceConfig is general config struct for HTTP services
type HTTPServiceConfig struct {
	Port            int    `mapstructure:"port"`
	Path            string `mapstructre:"path"`
	RuntimeBasePath string `mapstructure:"runtime_base_path"`
}

// StartWebhookService initiates all webhook services
func StartWebhookService(tunnel worker.PayloadTunnel, options HTTPServiceConfig) {

	webhookServer := http.NewServeMux()
	hook, err := github.New(github.Options.Secret("bah"))

	if err != nil {
		panic(err)
	}

	webhookServer.HandleFunc(
		options.Path, func(w http.ResponseWriter, r *http.Request) {
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

				workerPayload := worker.WorkerPayload{
					// Only one setup available at the moment.
					SetupId: "ab2kd",
					Data:    pushPayload,
				}

				tunnel.Tunnel <- workerPayload
			}
		},
	)

	err = http.ListenAndServe(fmt.Sprintf(":%d", options.Port), webhookServer)

	if err != nil {
		log.Error().Err(err).Msg("Unable to start http server.")
	}
}
