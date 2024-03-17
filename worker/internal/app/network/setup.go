package network

import (
	"Bobby/pkg/comm"
	"encoding/json"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/inancgumus/screen"
	"github.com/rs/zerolog/log"
	"net/url"
	"time"
)

type NetworkCredential struct {
	SetupID  string `json:"setup_id"`
	HostName string `json:"hostname"`
}

func NetworkSetup() bool {

	hostUrl := url.URL{Scheme: "ws"}
	hostUrl.Path = "worker"

	ok := SelectInstall(&hostUrl)

	if !ok {
		return false
	}

	s := spinner.New(spinner.CharSets[43], 100*time.Millisecond)

	// Test the connection to provided hostUrl
	ok, con := TestConnection(s, hostUrl)
	if !ok {
		return false
	}
	defer con.Close()

	ok = RegisterWorker(s, con)

	if !ok {
		return false
	}

	// Wait for host to process the request
	ok, hostResponse := WaitForRegisterResponse(s, con, hostUrl)

	if !ok {
		return false
	}

	// Parse host response
	var payload comm.SetupEventWorkerPayload
	err := json.Unmarshal(hostResponse, &payload)

	if err != nil {
		log.Error().Msg("Unable to parse host response.")
		return false
	}

	// If host response with invalid network id.
	if payload.SetupId == "" {
		s.FinalMSG = color.RedString("Host response with invalid setupID. Please try again.")
		s.Stop()
		return false
	}

	// Write received config to disk.
	WriteNetworkConfig(s, payload, hostUrl)

	// Clear screen
	time.Sleep(time.Second * 2)
	screen.Clear()
	screen.MoveTopLeft()
	return true
}
