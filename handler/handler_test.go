package handler

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const (
	testEnvFile = "test.env"
)

type mockAuthenticator struct{}

func (m *mockAuthenticator) AuthURL(state string) string {
	return "/mock-auth-url?state=" + state
}

func mockTemplate() *template.Template {
	return template.Must(template.New("mock").Parse("{{.}}"))
}

func TestSetupRoutes(t *testing.T) {
	r := mux.NewRouter()
	auth := &spotifyauth.Authenticator{}
	wrapper := &SpotifyAuthWrapper{auth}

	loginTemplate := mockTemplate()
	homeTemplate := mockTemplate()
	artistsTemplate := mockTemplate()
	tracksTemplate := mockTemplate()

	err := SetupRoutes(r, wrapper, testEnvFile, loginTemplate, homeTemplate, artistsTemplate, tracksTemplate)
	if err != nil {
		t.Fatalf("SetupRoutes returned error: %v", err)
	}

	if r.Get("/") == nil {
		t.Error("Expected route '/' to be registered")
	}

	if r.Get("/login") == nil {
		t.Error("Expected route '/login' to be registered")
	}

	if r.Get("/callback") == nil {
		t.Error("Expected route '/callback' to be registered")
	}

	if r.Get("/top-artists") == nil {
		t.Error("Expected route '/top-artists' to be registered")
	}

	if r.Get("/top-tracks") == nil {
		t.Error("Expected route '/top-tracks' to be registered")
	}
}

func TestHandleLoginPage(t *testing.T) {
	w := httptest.NewRecorder()

	mockTemplate := template.Must(template.New("mock").Parse("{{.}}"))
	handleLoginPage(w, mockTemplate)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleLogin(t *testing.T) {
	req := httptest.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	auth := &mockAuthenticator{}
	state := "teststate"

	handleLogin(w, req, auth, state)

	resp := w.Result()
	if resp.StatusCode != http.StatusFound {
		t.Errorf("Expected redirect status, got %d", resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if !strings.Contains(location, state) {
		t.Errorf("Redirect URL does not contain state: %s", location)
	}
}

func TestHandleTopArtistsPage(t *testing.T) {
	w := httptest.NewRecorder()

	cachedArtists = []ArtistInfo{{Name: "Test Artist"}}

	mockTmpl := mockTemplate()
	handleTopArtistsPage(w, mockTmpl)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleTopTracksPage(t *testing.T) {
	w := httptest.NewRecorder()

	cachedTracks = []TrackInfo{{Name: "Test Track", Artist: "Test Artist", TrackID: "123"}}

	mockTmpl := mockTemplate()
	handleTopTracksPage(w, mockTmpl)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleCallbackWrongState(t *testing.T) {
	req := httptest.NewRequest("GET", "/callback?state=wrongstate", nil)
	w := httptest.NewRecorder()

	auth := &spotifyauth.Authenticator{}
	state := "expectedstate"

	handleCallback(w, req, auth, state)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 for state mismatch, got %d", resp.StatusCode)
	}
}
