package main

import (
	"log"
	"net/http"

	envmanager "github.com/RaulRivadeneyra/stream-companion/internal/pkg/env-manager"
	"github.com/RaulRivadeneyra/stream-companion/internal/routes"
)

const OAUTH2_TOKEN_ENDPOINT = "https://id.twitch.tv/oauth2/token"

func main() {
	envmanager.LoadEnvs()

	mux := routes.NewRouter()

	log.Println("Listening...")
	http.ListenAndServe(":3000", mux)
}
