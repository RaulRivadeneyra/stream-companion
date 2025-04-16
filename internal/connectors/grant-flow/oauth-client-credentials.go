package grantflow

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/RaulRivadeneyra/stream-companion/internal/connectors/common"
)

const OAUTH2_TOKEN_ENDPOINT = "https://id.twitch.tv/oauth2/token"

type ClientCredentialsRequest struct {
	ClientId     string
	ClientSecret string
	GrantType    common.GrantType
	AuthURL      string
}

type ClientCredentialsReponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   uint32 `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// See Client Credential Grant Flow [docs]
//
// [docs]: https://datatracker.ietf.org/doc/html/rfc6749#section-1.3.4
func GetOAuthClientCredentials(params ClientCredentialsRequest) ClientCredentialsReponse {

	data := url.Values{}
	data.Set("client_id", params.ClientId)
	data.Set("client_secret", params.ClientSecret)
	data.Set("grant_type", string(params.GrantType))

	httpClient := &http.Client{}
	request, err := http.NewRequest(
		http.MethodPost,
		params.AuthURL,
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		log.Fatal("Unable to create new request")
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := httpClient.Do(request)

	if err != nil {
		log.Fatal("No response from request")
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Fatal("Unable to read response body")
	}

	resp := ClientCredentialsReponse{}
	json.Unmarshal(body, &resp)
	return resp
}
