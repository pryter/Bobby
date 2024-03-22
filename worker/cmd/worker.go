package cmd

import (
	"bobby-worker/internal/app/network"
	"bobby-worker/internal/app/resources"
	"bobby-worker/internal/app/runtime"
	"bobby-worker/internal/challenge"
	"bobby-worker/internal/events/pushEvent"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/url"
)

type WorkerServiceOptions struct {
	ServiceBasePath string `mapstructure:"service_base_path"`
}

func StartWorkerService(options WorkerServiceOptions) bool {

	// Check network-credential
	netConfigFile := resources.GetResources().Configs.NetworkConfigFile
	data, err := netConfigFile.Open()

	if err != nil {
		ok := network.NetworkSetup()
		if ok {
			StartWorkerService(options)
			return false
		}
		return ok
	}

	// Parse net config file
	var networkCreds network.NetworkCredential
	err = json.Unmarshal(data, &networkCreds)

	if err != nil {
		log.Error().Msg("Unable to parse config")
	}

	// Create Challenge
	log.Debug().Msg("Preparing resources for creating authentication challenge.")
	challengeBuilder := challenge.Builder{
		PublicKeyPath: resources.GetResources().Configs.MapFile("network-key.pem").AbsolutePath,
		SetupID:       networkCreds.SetupID,
	}

	c, err := challengeBuilder.Generate()

	if err != nil {
		log.Error().Err(err)
		return false
	}

	log.Debug().Msg("Challenge created. Starting a connection.")
	u := url.URL{
		Scheme:   "ws",
		Host:     networkCreds.HostName,
		Path:     "/worker",
		RawQuery: "sid=" + url.QueryEscape(networkCreds.SetupID) + "&challenge=" + url.QueryEscape(c),
	}

	w := runtime.WorkerRuntime{ConnectionUrl: u}

	w.RegisterEvent(
		"push", func(rawPayload json.RawMessage) {
			pushEvent.WebhookPushEvent(
				rawPayload,
				pushEvent.WebhookPushEventOptions{RuntimeBasePath: options.ServiceBasePath},
			)
		},
	)

	ok := w.Run()

	// Restart if not ok
	return !ok
}
