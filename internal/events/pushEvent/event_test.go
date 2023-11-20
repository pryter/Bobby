package pushEvent

import (
	"Bobby/internal/utils"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/joho/godotenv"
	"path/filepath"
	"testing"
)

func TestWebhookPushEvent(t *testing.T) {

	godotenv.Load(filepath.Join(utils.GetProjectRoot(), ".env"))

	repository := github.PushPayload{}.Repository
	repository.ID = 571145096
	repository.CloneURL = "https://github.com/pryter/anri.git"

	installation := github.PushPayload{}.Installation
	installation.ID = 44151598

	_ = github.PushPayload{
		Ref:          "testRef",
		Repository:   repository,
		Installation: installation,
	}

	t.Run("environment_setup", TestEnvironmentSetup)
	t.Run("cli_tests", TestCliFactory)

}
