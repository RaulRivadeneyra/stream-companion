package grantflow

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/RaulRivadeneyra/stream-companion/internal/connectors/common"
)

type AuthorizationCodeResponse struct {
	Code  common.Code  `json:"code"`
	Scope common.Scope `json:"scope"`
	State common.State `json:"state"`
}

type AuthorizationCodeError struct {
	Error            string       `json:"error"`
	ErrorDescription string       `json:"error_description"`
	State            common.State `json:"state"`
}

type AuthorizationTokenRequest struct {
	ClientId     string
	ClientSecret string
	Code         common.Code
	GrantType    common.GrantType
	AuthURL      string
	RedirectURI  string
}

type AuthorizationTokenReponse struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int64    `json:"expires_in"`
	IdToken      string   `json:"id_token"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

func GetAuthorizationToken(params AuthorizationTokenRequest) AuthorizationTokenReponse {
	data := url.Values{}
	data.Set("client_id", params.ClientId)
	data.Set("client_secret", params.ClientSecret)
	data.Set("code", params.Code)
	data.Set("grant_type", string(params.GrantType))
	data.Set("redirect_uri", params.RedirectURI)

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

	fmt.Println(response)

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Fatal("Unable to read response body")
	}

	resp := AuthorizationTokenReponse{}
	json.Unmarshal(body, &resp)
	return resp
}

// TODO: Implement refresh tokens, e.g. https://dev.twitch.tv/docs/authentication/refresh-tokens/
