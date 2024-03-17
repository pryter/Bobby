package network

import (
	"github.com/manifoldco/promptui"
	"net/url"
)

func SelectInstall(hostUrl *url.URL) bool {
	prompt := promptui.Select{
		Label: "Select network mode.",
		Items: []string{"Automatic (recommend)", "Manual"},
	}

	_, result, _ := prompt.Run()

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

	return true
}
