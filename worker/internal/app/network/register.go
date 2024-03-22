package network

import (
	"Bobby/pkg/comm"
	"Bobby/pkg/utils"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/url"
	"time"
)

func RegisterWorker(s *spinner.Spinner, con *websocket.Conn) bool {
	s.Suffix = "Registering this worker to main service."

	macAddr, err := utils.GetMacAddr()

	if err != nil {
		log.Error().Msg("Unable to get machine address")
		return false
	}

	regComm := comm.HostCommand[comm.RegisterCommandPayload]{
		Instruction: "register",
		Payload:     comm.RegisterCommandPayload{MacAddr: macAddr[0]},
	}

	con.WriteMessage(1, regComm.Digest())

	s.Suffix = " Waiting for the host machine to respond (60s)."
	return true
}

func WaitForRegisterResponse(
	s *spinner.Spinner,
	con *websocket.Conn,
	hostUrl url.URL,
) (bool, []byte) {
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
				return false, nil
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
			return false, nil
		}
	}

	return true, response
}
