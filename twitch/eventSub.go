package TwitchAccessTocen

import (
	"context"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
	"log"
)

func AccessToken(ClientID, clientSecret string) string {
	oauth2Config := &clientcredentials.Config{
		ClientID:     ClientID,
		ClientSecret: clientSecret,
		TokenURL:     twitch.Endpoint.TokenURL,
	}

	token, err := oauth2Config.Token(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return token.AccessToken
}
