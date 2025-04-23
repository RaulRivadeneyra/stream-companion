package main

import (
	"os"
	"path/filepath"

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

	projectPath, _ := os.Getwd()
	content, err := os.ReadFile(filepath.Join(projectPath, "/actions/test-1.lua"))
	if err != nil {
		panic(err)
	}

	node := nodes.NewActionNode()
	node.SetLabel("testNode")
	node.SetValue(string(content))

	globalNode := nodes.NewActionNode()
	globalNode.SetLabel("GLOBAL")

	sv := nodes.NewSharedVariable(globalNode)
	sv.SetName("name")
	sv.SetString("Yujiko")

	svc := nodes.NewSharedVariableCollection()

	svc.AddSharedVariable(sv)

	err = node.Execute(svc)

	if err != nil {
		panic(err)
	}
	content, err = os.ReadFile(filepath.Join(projectPath, "/actions/test-2.lua"))
	if err != nil {
		panic(err)
	}
	node2 := nodes.NewActionNode()
	node2.SetLabel("testNode2")
	node2.SetValue(string(content))

	err = node2.Execute(svc)

	if err != nil {
		panic(err)
	}

}
