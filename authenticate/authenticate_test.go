package authenticate

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

const (
	testEnvFile = "test.env"
)

func TestSpotifyBadAuthenticate(t *testing.T) {
	err := godotenv.Load(testEnvFile)
	if err != nil {
		t.Fatalf("Error loading test.env file: %v", err)
	}

	os.Setenv("CLIENT_ID", os.Getenv("CLIENT_ID"))
	os.Setenv("CLIENT_SECRET", os.Getenv("CLIENT_SECRET"))

	auth, err := SpotifyAuthenticate(testEnvFile)
	if err != nil {
		t.Fatalf("SpotifyAuthenticate returned error: %v", err)
	}

	if auth == nil {
		t.Error("Expected authenticator, got nil")
	}
}
