package network

import (
	"Bobby/pkg/comm"
	"bobby-worker/internal/app/resources"
	"encoding/json"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"net/url"
	"os"
	"time"
)

func TestConnection(s *spinner.Spinner, hostUrl url.URL) (bool, *websocket.Conn) {
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
		return false, nil
	}
	return true, con
}

func WriteNetworkConfig(s *spinner.Spinner, payload comm.SetupEventWorkerPayload, hostUrl url.URL) {
	s.Suffix = fmt.Sprintf(
		" Registration accpeted with setupID %s", color.GreenString("#%s", payload.SetupId),
	)

	settings, _ := json.Marshal(
		NetworkCredential{
			SetupID:  payload.SetupId,
			HostName: hostUrl.Host,
		},
	)

	_ = os.WriteFile(
		resources.GetResources().Configs.NetworkConfigFile.AbsolutePath, settings, 0777,
	)

	keyFile := resources.GetResources().Configs.MapFile("network-key.pem")
	_ = os.WriteFile(keyFile.AbsolutePath, []byte(payload.Data), 0777)

	s.FinalMSG = "Installation finished. Starting the service...\n"
	s.Stop()
}
