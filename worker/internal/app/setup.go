package app

import (
	"Bobby/pkg/comm"
	"bobby-worker/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/inancgumus/screen"
	"github.com/manifoldco/promptui"
	"net/url"
	"os"
	"path"
	"time"
)

type NetworkCredential struct {
	SetupID  string `json:"setup_id"`
	HostName string `json:"hostname"`
}

type HTTPServiceConfig struct {
	Port            int    `mapstructure:"port"`
	Path            string `mapstructre:"path"`
	RuntimeBasePath string `mapstructure:"runtime_base_path"`
}

func NetworkSetup(config HTTPServiceConfig) bool {
	screen.Clear()
	prompt := promptui.Select{
		Label: "Select setup mode.",
		Items: []string{"Automatic (recommend)", "Manual"},
	}

	_, result, _ := prompt.Run()

	hostUrl := url.URL{Scheme: "ws"}
	hostUrl.Path = "worker"

	switch result {
	case "Automatic (recommend)":
		hostUrl.Host = "localhost:4040"
		break
	case "Manual":
		manHost := promptui.Prompt{Label: "Main service hostname (hostname:port)"}
		result, _ := manHost.Run()
		hostUrl.Host = result
		break
	default:
		return false
	}

	s := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Connecting to %s", hostUrl.String())
	s.Start()
	time.Sleep(time.Second * 1)
	con, _, err := websocket.DefaultDialer.Dial(hostUrl.String(), nil)

	if err != nil {
		s.FinalMSG = fmt.Sprintf(
			"%s (%s).", color.RedString("✗ Unable to connect to the host"), hostUrl.String(),
		)
		s.Stop()
		color.Red("\nReason: %s", err.Error())
		return false
	}

	defer con.Close()

	s.Suffix = "Registering this comm to main service."

	macAddr, err := utils.GetMacAddr()
	regComm := comm.HostCommand[comm.RegisterCommandPayload]{
		Instruction: "register",
		Payload:     comm.RegisterCommandPayload{MacAddr: macAddr[0]},
	}

	con.WriteMessage(1, regComm.Digest())

	s.Suffix = " Waiting for the host machine to respond (60s)."

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	counter := 60
	cres := make(chan []byte)
	q := make(chan struct{})

	// Wait for response
	go func() {
		_, response, err := con.ReadMessage()
		if err != nil {
			s.FinalMSG = fmt.Sprintf(
				"%s (%s).", color.RedString("✗ Host service responded with error."),
				hostUrl.String(),
			)
			s.Stop()
			q <- struct{}{}
			return
		}

		cres <- response
	}()

	var response []byte

L:
	for {
		select {
		case <-ticker.C:
			if counter <= 0 {
				s.FinalMSG = fmt.Sprintf(
					"%s (%s).", color.RedString("✗ Host service response timeout."),
					hostUrl.String(),
				)
				s.Stop()
				return false
			} else {
				s.Suffix = fmt.Sprintf(
					" Waiting for the host machine to respond (%ds).", counter,
				)
				counter--
			}
			break
		case response = <-cres:
			time.Sleep(time.Second * 1)
			break L
		case <-q:
			return false
		}
	}

	var payload comm.WorkerPayload

	err = json.Unmarshal(response, &payload)

	if payload.SetupId == "" {
		s.FinalMSG = color.RedString("Host response with invalid setupID. Please try again.")
		s.Stop()
		return false
	}

	s.Suffix = fmt.Sprintf(
		" Registration accpeted with setupID %s", color.GreenString("#%s", payload.SetupId),
	)

	settings, _ := json.Marshal(
		NetworkCredential{
			SetupID:  payload.SetupId,
			HostName: hostUrl.Host,
		},
	)

	os.WriteFile(path.Join(config.RuntimeBasePath, "network-credentials.json"), settings, 777)
	s.FinalMSG = "Installation finished. Starting the service...\n"
	s.Stop()
	time.Sleep(time.Second * 2)
	screen.Clear()
	screen.MoveTopLeft()
	return true
}
