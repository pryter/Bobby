package worker

import (
	"Bobby/internal/app"
	"Bobby/pkg/comm"
	"Bobby/pkg/crypto"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"os"
	"path"
	"regexp"
)

func isValidMACAddress(mac string) bool {
	// Define the regex pattern
	var macRegex = regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
	return macRegex.MatchString(mac)
}

// registerAction occurs when the command register received from worker
func registerAction(conn *websocket.Conn, command comm.RawHostCommand) {
	// Generate uuid
	id := uuid.New()
	workerP := comm.SetupEventWorkerPayload{
		SetupId: id.String(),
	}

	// Resolve register payload
	var payload comm.RegisterCommandPayload
	err := command.ResolvePayload(&payload)
	if err != nil {
		return
	}

	// Check mac addr validity
	if !isValidMACAddress(payload.MacAddr) {
		return
	}

	// Generate RSA key pair
	private, public, err := crypto.CreateRSAKeyPair()
	if err != nil {
		return
	}

	// Add public key to worker payload
	workerP.Data = string(public)
	// Save private key to file
	resource := app.GetResources()
	workerFolder := path.Join(resource.WorkerData.GetAbsolutePath(), id.String())
	_ = os.Mkdir(workerFolder, 0777)

	// Write key
	err = os.WriteFile(path.Join(workerFolder, "key"), private, 0777)
	if err != nil {
		return
	}

	// Write Secret
	err = os.WriteFile(path.Join(workerFolder, "secret"), []byte(payload.MacAddr), 0777)
	if err != nil {
		return
	}

	// Parse entire worker payload
	d, err := workerP.Digest()
	if err != nil {
		return
	}

	// Send worker payload
	err = conn.WriteMessage(websocket.TextMessage, d)
	if err != nil {
		return
	}
}
