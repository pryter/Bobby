package cmd

import (
	"Bobby/pkg/crypto"
	"bobby-worker/internal/app/network"
	"bobby-worker/internal/app/resources"
	"bobby-worker/internal/app/runtime"
	"bobby-worker/internal/events/pushEvent"
	"bobby-worker/internal/utils"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/url"
	"time"
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
	pubKey, err := resources.GetResources().Configs.MapFile("network-key.pem").Open()

	if err != nil {
		log.Error().Msg("Unable to find network key")
		return false
	}

	macAddr, err := utils.GetMacAddr()
	if err != nil || macAddr[0] == "" {
		log.Error().Msg("Unable to get mac address")
		return false
	}

	log.Debug().Msg("Generating authentication challenge.")
	challengeRawText := macAddr[0] + "|" + networkCreds.SetupID + "|" + time.Now().String()
	challenge, err := crypto.RSAEncrypt(challengeRawText, pubKey)

	if err != nil {
		log.Error().Msg("Unable to generate challenge")
		return false
	}

	log.Debug().Msg("Challenge created. Starting a connection.")
	u := url.URL{
		Scheme:   "ws",
		Host:     networkCreds.HostName,
		Path:     "/worker",
		RawQuery: "sid=" + url.QueryEscape(networkCreds.SetupID) + "&challenge=" + url.QueryEscape(challenge),
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
