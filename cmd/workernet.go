package cmd

import (
	"Bobby/internal/worker"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

func StartWorkerNetwork(workernet worker.WorkerNetwork) {

	workerServer := http.NewServeMux()
	workerServer.HandleFunc(
		"/worker", workernet.HttpHandler,
	)

	err := http.ListenAndServe(fmt.Sprintf(":%d", 4040), workerServer)

	if err != nil {
		log.Error().Err(err).Msg("Unable to start http server.")
	}
}
