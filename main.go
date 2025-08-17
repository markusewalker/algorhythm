package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"algorhythm/authenticate"
	"algorhythm/handler"

	"github.com/gorilla/mux"
	"github.com/pkg/browser"
)

const (
	envFile             = ".env"
	loginTemplateFile   = "views/login.html"
	artistsTemplateFile = "views/top_artists.html"
	tracksTemplateFile  = "views/top_tracks.html"
)

func main() {
	auth, err := authenticate.SpotifyAuthenticate(envFile)
	if err != nil {
		log.Fatalf("Error authenticating: %v", err)
	}

	r := mux.NewRouter()
	loginTemplate := template.Must(template.ParseFiles(loginTemplateFile))
	homeTemplate := template.Must(template.ParseFiles("views/home.html"))
	artistsTemplate := template.Must(template.ParseFiles(artistsTemplateFile))
	tracksTemplate := template.Must(template.ParseFiles(tracksTemplateFile))

	handler.SetupRoutes(r, &handler.SpotifyAuthWrapper{Authenticator: auth}, envFile, loginTemplate, homeTemplate, artistsTemplate, tracksTemplate)

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("views/assets"))))

	http.Handle("/", r)

	fmt.Println("In your browser, please navigate to: http://localhost:8080")
	browser.OpenURL("http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
