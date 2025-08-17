package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

func isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("auth")
	return err == nil && cookie.Value == "true"
}

func handleHomePage(w http.ResponseWriter, template *template.Template, userName string) {
	template.Execute(w, map[string]interface{}{"UserName": userName})
}

var (
	cachedArtists []ArtistInfo
	cachedTracks  []TrackInfo
)

type ArtistInfo struct {
	Name        string
	ImageURL    string
	ExternalURL string
}

type TrackInfo struct {
	Name    string
	Artist  string
	TrackID string
}

type AuthURLProvider interface {
	AuthURL(state string) string
}

type SpotifyAuthWrapper struct {
	*spotifyauth.Authenticator
}

func (s *SpotifyAuthWrapper) AuthURL(state string) string {
	return s.Authenticator.AuthURL(state)
}

func SetupRoutes(r *mux.Router, auth AuthURLProvider, envFile string, loginTemplate, homeTemplate, artistsTemplate, tracksTemplate *template.Template) error {
	err := godotenv.Load(envFile)
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	state := os.Getenv("STATE")
	if state == "" {
		return fmt.Errorf("STATE MUST be set")
	}

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !isAuthenticated(r) {
			handleLoginPage(w, loginTemplate)
			return
		}

		userName := ""
		if cookie, err := r.Cookie("username"); err == nil {
			userName = cookie.Value
		}

		handleHomePage(w, homeTemplate, userName)
	}).Name("/")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, auth, state)
	}).Name("/login")

	r.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if realAuth, ok := auth.(*SpotifyAuthWrapper); ok {
			handleCallback(w, r, realAuth.Authenticator, state)
		} else {
			http.Error(w, "Authenticator not supported for callback", http.StatusInternalServerError)
		}
	}).Name("/callback")

	r.HandleFunc("/top-artists", func(w http.ResponseWriter, r *http.Request) {
		handleTopArtistsPage(w, artistsTemplate)
	}).Name("/top-artists")

	r.HandleFunc("/top-tracks", func(w http.ResponseWriter, r *http.Request) {
		handleTopTracksPage(w, tracksTemplate)
	}).Name("/top-tracks")

	return nil
}

// handleLoginPage is a handler for the login page
func handleLoginPage(w http.ResponseWriter, template *template.Template) {
	template.Execute(w, nil)
}

// handleLogin is a handler for the login
func handleLogin(w http.ResponseWriter, r *http.Request, auth AuthURLProvider, state string) {
	url := auth.AuthURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

// handleTopArtistsPage is a handler for the top artists page
func handleTopArtistsPage(w http.ResponseWriter, template *template.Template) {
	template.Execute(w, map[string]any{"TopArtists": cachedArtists})
}

// handleTopTracksPage is a handler for the top tracks page
func handleTopTracksPage(w http.ResponseWriter, template *template.Template) {
	template.Execute(w, map[string]any{"TopTracks": cachedTracks})
}

func handleCallback(w http.ResponseWriter, r *http.Request, auth *spotifyauth.Authenticator, state string) {
	if r.FormValue("state") != state {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	token, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Println("Token error:", err)
		return
	}

	client := spotify.New(auth.Client(r.Context(), token))

	userName := ""
	user, err := client.CurrentUser(r.Context())
	if err == nil {
		userName = user.DisplayName
	}

	topArtists, err := client.CurrentUsersTopArtists(r.Context(), spotify.Limit(5), spotify.Timerange("short_term"))
	if err != nil {
		http.Error(w, "Failed to get top artists", http.StatusInternalServerError)
		log.Println("TopArtists error:", err)
		return
	}

	topTracks, err := client.CurrentUsersTopTracks(r.Context(), spotify.Limit(5), spotify.Timerange("short_term"))
	if err != nil {
		http.Error(w, "Failed to get top tracks", http.StatusInternalServerError)
		log.Println("TopTracks error:", err)
		return
	}

	cachedTracks = nil
	for _, t := range topTracks.Tracks {
		artistName := ""
		if len(t.Artists) > 0 {
			artistName = t.Artists[0].Name
		}

		// Extract track ID from URI (format: spotify:track:TRACKID)
		uriString := string(t.URI)
		trackID := ""

		parts := strings.Split(uriString, ":")
		if len(parts) == 3 {
			trackID = parts[2]
		}

		cachedTracks = append(cachedTracks, TrackInfo{
			Name:    t.Name,
			Artist:  artistName,
			TrackID: trackID,
		})
	}

	cachedArtists = nil
	for _, a := range topArtists.Artists {
		imageURL := ""
		if len(a.Images) > 0 {
			imageURL = a.Images[0].URL
		}

		if len(a.ExternalURLs) == 0 {
			log.Println("No external URL found for artist:", a.Name)
			continue
		}

		cachedArtists = append(cachedArtists, ArtistInfo{
			Name:        a.Name,
			ImageURL:    imageURL,
			ExternalURL: a.ExternalURLs["spotify"],
		})
	}

	// Set authentication and username cookies
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "true",
		Path:   "/",
		MaxAge: 3600,
	})

	http.SetCookie(w, &http.Cookie{
		Name:   "username",
		Value:  userName,
		Path:   "/",
		MaxAge: 3600,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
