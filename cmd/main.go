package main

import (
	"log"
	"net/http"

	"github.com/RaulRivadeneyra/stream-companion/api/routes"
	envmanager "github.com/RaulRivadeneyra/stream-companion/pkg/env-manager"
)

const OAUTH2_TOKEN_ENDPOINT = "https://id.twitch.tv/oauth2/token"

func main() {
	envmanager.LoadEnvs()

	mux := routes.NewRouter()

	log.Println("Listening...")
	http.ListenAndServe(":3000", mux)
}
