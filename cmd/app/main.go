package main

import (
	"fmt"
	"os"

	"github.com/RaulRivadeneyra/stream-companion/internal/nodes"
	envmanager "github.com/RaulRivadeneyra/stream-companion/internal/pkg/env-manager"
)

const OAUTH2_TOKEN_ENDPOINT = "https://id.twitch.tv/oauth2/token"

func main() {
	envmanager.LoadEnvs()
	// mux := routes.NewRouter()
	//
	// log.Println("Listening...")
	// http.ListenAndServe(":3000", mux)
	node := nodes.NewActionNode()

	content, err := os.ReadFile("ABS PATH TO ../../actions/test-1.lua")
	if err != nil {
		panic(err)
	}
	node.SetCode(string(content))

	inputs := map[string]any{
		"name": "Yujiko",
	}

	result, _ := node.Execute(inputs)

	fmt.Println(result)
}
