package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
	template.Execute(w, map[string]interface{}{
		"UserName":          userName,
		"TopGenre":          cachedHomeStats.TopGenre,
		"ListeningSessions": cachedHomeStats.ListeningSessions,
		"UniqueArtists":     cachedHomeStats.UniqueArtists,
		"FreshDiscoveries":  cachedHomeStats.FreshDiscoveries,
	})
}

var (
	cachedArtists   []ArtistInfo
	cachedTracks    []TrackInfo
	cachedHomeStats = HomeStats{
		TopGenre:          "Not enough data yet",
		ListeningSessions: 0,
		UniqueArtists:     0,
		FreshDiscoveries:  0,
	}
)

type HomeStats struct {
	TopGenre          string
	ListeningSessions int
	UniqueArtists     int
	FreshDiscoveries  int
}

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

// handleCallback is a handler for the Spotify authentication callback.
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

	topArtists, err := client.CurrentUsersTopArtists(r.Context(), spotify.Limit(3), spotify.Timerange("short_term"))
	if err != nil {
		http.Error(w, "Failed to get top artists", http.StatusInternalServerError)
		log.Println("TopArtists error:", err)
		return
	}

	topTracks, err := client.CurrentUsersTopTracks(r.Context(), spotify.Limit(3), spotify.Timerange("short_term"))
	if err != nil {
		http.Error(w, "Failed to get top tracks", http.StatusInternalServerError)
		log.Println("TopTracks error:", err)
		return
	}

	recentlyPlayed, recentErr := client.PlayerRecentlyPlayed(r.Context())
	if recentErr != nil {
		log.Println("RecentlyPlayed error:", recentErr)
	}

	cachedTracks = nil
	for _, t := range topTracks.Tracks {
		artistName := ""
		if len(t.Artists) > 0 {
			artistName = t.Artists[0].Name
		}

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

		if len(cachedTracks) == 3 {
			break
		}
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

	cachedHomeStats = deriveHomeStats(topArtists.Artists, topTracks.Tracks, recentlyPlayed)

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

// deriveHomeStats is a helper function to calculate the home page stats based on the user's top artists,
// tracks, and recently played items.
func deriveHomeStats(artists []spotify.FullArtist, tracks []spotify.FullTrack, recent []spotify.RecentlyPlayedItem) HomeStats {
	stats := HomeStats{
		TopGenre:          "Not enough data yet",
		ListeningSessions: len(tracks),
		UniqueArtists:     0,
		FreshDiscoveries:  0,
	}

	// Looking to find the most common genre amongst the top artists.
	genreCounts := map[string]int{}
	for _, artist := range artists {
		for _, genre := range artist.Genres {
			normalized := strings.TrimSpace(genre)
			if normalized == "" {
				continue
			}

			genreCounts[normalized]++
		}
	}

	maxGenreCount := 0
	for genre, count := range genreCounts {
		if count > maxGenreCount {
			maxGenreCount = count
			stats.TopGenre = genre
		}
	}

	// Looking to see if there are any fresh discoveries in the top tracks.
	now := time.Now()
	for _, track := range tracks {
		releaseDate := strings.TrimSpace(track.Album.ReleaseDate)
		if releaseDate != "" && isSameMonthRelease(releaseDate, now) {
			stats.FreshDiscoveries++
		}
	}

	if len(recent) > 0 {
		monthlyArtists := map[string]struct{}{}
		stats.ListeningSessions = 0

		for _, item := range recent {
			if item.PlayedAt.Year() != now.Year() || item.PlayedAt.Month() != now.Month() {
				continue
			}

			stats.ListeningSessions++

			for _, artist := range item.Track.Artists {
				name := strings.TrimSpace(artist.Name)
				if name != "" {
					monthlyArtists[name] = struct{}{}
				}
			}
		}

		stats.UniqueArtists = len(monthlyArtists)
		return stats
	}

	uniqueArtists := map[string]struct{}{}
	for _, track := range tracks {
		for _, artist := range track.Artists {
			name := strings.TrimSpace(artist.Name)
			if name != "" {
				uniqueArtists[name] = struct{}{}
			}
		}
	}

	stats.UniqueArtists = len(uniqueArtists)

	return stats
}

// isSameMonthRelease is a helper function to determine if a track's release date is in the same
// month and year as the current date. It handles both "YYYY-MM-DD" and "YYYY-MM" formats.
func isSameMonthRelease(releaseDate string, now time.Time) bool {
	if len(releaseDate) >= 7 {
		year, yearErr := strconv.Atoi(releaseDate[0:4])
		month, monthErr := strconv.Atoi(releaseDate[5:7])
		if yearErr == nil && monthErr == nil {
			return year == now.Year() && month == int(now.Month())
		}
	}

	return false
}
