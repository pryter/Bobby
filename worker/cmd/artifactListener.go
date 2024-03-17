package cmd

import (
	"bobby-worker/internal/bucket"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

func pipeJSONStruct(w http.ResponseWriter, d interface{}) {
	b, _ := json.Marshal(d)
	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write(b); err != nil {
		return
	}
}

func pipeHTMLString(w http.ResponseWriter, str string) {
	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write([]byte(str)); err != nil {
		return
	}
}

func pipeFile(w http.ResponseWriter, fileBuffer []byte, filename string) {
	contentDepo := fmt.Sprintf("attachment; filename=%s", filename)
	w.Header().Set("Content-Disposition", contentDepo)
	if _, err := w.Write(fileBuffer); err != nil {
		return
	}
}

type ArtifactServiceOptions struct {
	Port            int    `mapstructure:"port"`
	Path            string `mapstructre:"path"`
	RuntimeBasePath string `mapstructure:"runtime_base_path"`
}

func StartServingArtifacts(options ArtifactServiceOptions) {

	artifactServer := http.NewServeMux()

	fb := bucket.Bucket{RootPath: options.RuntimeBasePath}
	println(fb.RootPath)

	artifactServer.HandleFunc(
		options.Path, func(w http.ResponseWriter, r *http.Request) {

			/*
				====================================
				TODO Implement Authentication System
				====================================
				Contact Authentication Database
				server to verify Bearer token.
				Authentication server might be Google
				Firebase Firestore.
			*/

			query, err := bucket.NewQuery(r.URL.Path)

			if err != nil {
				log.Error().Err(err).Msg("There is an error while parsing query.")
				pipeHTMLString(w, "Error")
				return
			}

			resolved, err := query.Resolve(fb)

			buffer, err := fb.ReadFile(resolved)

			if err != nil {
				log.Error().Err(err).Str(
					"resolved", r.URL.Path,
				).Msg("Unable to access requested file")

				pipeHTMLString(w, "Error")
				return
			}

			pipeFile(w, buffer, resolved.Filename)
		},
	)

	err := http.ListenAndServe(fmt.Sprintf(":%d", options.Port), artifactServer)

	if err != nil {
		log.Error().Err(err).Msg("Unable to start artifact bucket server.")
	}
}
