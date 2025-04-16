package twitch

import (
	"fmt"
	"os"

	"github.com/RaulRivadeneyra/stream-companion/internal/connectors/common"
	grantflow "github.com/RaulRivadeneyra/stream-companion/internal/connectors/grant-flow"
)

const OAUTH2_TOKEN_ENDPOINT = "https://id.twitch.tv/oauth2/token"

func GetAuthorizationToken(code common.Code, redURL string) grantflow.AuthorizationTokenReponse {
	req := grantflow.AuthorizationTokenRequest{
		ClientId:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		GrantType:    "authorization_code",
		Code:         code,
		AuthURL:      OAUTH2_TOKEN_ENDPOINT,
		RedirectURI:  "http://localhost:3000/oauth/redirect",
	}

	fmt.Println(req)
	return grantflow.GetAuthorizationToken(req)
}

func GetClientCredentials() grantflow.ClientCredentialsReponse {
	req := grantflow.ClientCredentialsRequest{
		ClientId:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		GrantType:    "client_credentials",
		AuthURL:      OAUTH2_TOKEN_ENDPOINT,
	}
	return grantflow.GetOAuthClientCredentials(req)
}
