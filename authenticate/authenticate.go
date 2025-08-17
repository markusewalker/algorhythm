package authenticate

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

var (
	authenticator *spotifyauth.Authenticator
	redirectURI   = "http://127.0.0.1:8080/callback"
)

// SpotifyAuthenticate enables login with your unique client ID and secret, loading env variables from the specified file.
func SpotifyAuthenticate(envFile string) (*spotifyauth.Authenticator, error) {
	err := godotenv.Load(envFile)
	if err != nil {
		return nil, fmt.Errorf("error loading %s file: %v", envFile, err)
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("CLIENT_ID and CLIENT_SECRET must be set")
	}

	authenticator = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserTopRead),
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithClientSecret(clientSecret),
	)

	return authenticator, nil
}
