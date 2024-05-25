package main

import (
	"time"

	"github.com/DecxBase/gateway/node"
)

func main() {
	loc, _ := time.LoadLocation("Asia/Kathmandu")
	options := &node.ClientOptions{
		Timezone: loc,
	}
	client := node.Client("auth", options).
		SetTitle("Authentication").
		SetDescription("Handles authentication services on api")

	client.Logger.Debug().Msgf("Client: %s", client.Name)
}
