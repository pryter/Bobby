package main

import (
	"Bobby/internal/events/pushEvent"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

const (
	path = "/webhooks"
)

func main() {
	godotenv.Load()

	hook, err := github.New(github.Options.Secret("bah"))

	if err != nil {
		panic(err)
	}

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.PushEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok events wasn't one of the ones asked to be parsed
			}
		}
		switch payload.(type) {

		case github.PushPayload:
			pushPayload := payload.(github.PushPayload)
			pushEvent.WebhookPushEvent(pushPayload)
		}
	})

	print("Bobby is listening for some fresh pushes.")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
