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
	envFile = ".env"

	artistsTemplateFile = "views/top_artists.html"
	footerTemplateFile  = "views/shared/footer.html"
	homeTemplateFile    = "views/home.html"
	loginTemplateFile   = "views/login.html"
	tracksTemplateFile  = "views/top_tracks.html"
)

func main() {
	auth, err := authenticate.SpotifyAuthenticate(envFile)
	if err != nil {
		log.Fatalf("Error authenticating: %v", err)
	}

	r := mux.NewRouter()

	loginTemplate := template.Must(template.ParseFiles(loginTemplateFile, footerTemplateFile))
	homeTemplate := template.Must(template.ParseFiles(homeTemplateFile, footerTemplateFile))
	artistsTemplate := template.Must(template.ParseFiles(artistsTemplateFile, footerTemplateFile))
	tracksTemplate := template.Must(template.ParseFiles(tracksTemplateFile, footerTemplateFile))

	handler.SetupRoutes(r, &handler.SpotifyAuthWrapper{Authenticator: auth}, envFile, loginTemplate, homeTemplate, artistsTemplate, tracksTemplate)

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.Handle("/", r)

	fmt.Println("In your browser, please navigate to: http://localhost:8080")
	browser.OpenURL("http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
