package event

import (
	"Bobby/internal/token"
	"github.com/go-playground/webhooks/v6/github"
)

func WebhookPushEvent(payload github.PushPayload) {
	install := payload.Installation

	/*
		TODO
		[X] generate app access token
		[ ] clone or pull repository
		[ ] build project
		[ ] create commit check run
		[ ] provide artifacts url
	*/

	token.GenerateAccessToken(install.ID, payload.Repository.ID)

}
