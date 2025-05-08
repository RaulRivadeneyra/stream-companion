package routes

import (
	"fmt"
	"net/http"

	"github.com/RaulRivadeneyra/stream-companion/commerce/connectors/twitch"
)

const IndexURI = "/"
const OAuthRedirectURI = "/oauth/redirect"
const TokenRedirectURI = "/oauth/token/redirect"

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(IndexURI, indexHandler)
	mux.HandleFunc(OAuthRedirectURI, oauthRedirectHandler)
	mux.HandleFunc(TokenRedirectURI, tokenRedirectHandler)
	return mux
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authorized?!")
	// fmt.Println(r)
}

func oauthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	if queryParams.Has("error") {
		fmt.Println(
			queryParams.Get("error"),
			queryParams.Get("error_description"),
			queryParams.Get("state"),
		)
		return
	}

	code := queryParams.Get("code")
	scope := queryParams.Get("scope")
	state := queryParams.Get("state")

	fmt.Println(code, scope, state)

	resp := twitch.GetAuthorizationToken(code, "http://localhost:3000"+TokenRedirectURI)
	fmt.Println(resp)

	fmt.Fprintln(w, "Authorization successful, you can close this window!")

}

func tokenRedirectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Token?!", r)
}
