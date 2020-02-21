package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	social "github.com/kkdai/line-social-sdk-go"
)

//LINE Login related configuration
var channelID, channelSecret string

//LINE MessageAPI related configuration
var serverURL string
var botToken, botSecret string
var socialClient *social.Client

var state string

type LineToken struct {
	AccessToken  string
	RefreshToken string
}

func main() {
	var err error
	serverURL = "https://line-login-soical.herokuapp.com/"
	channelID = "1653859637"
	channelSecret = "350c50d8c1e4435726d64450f45142c3"

	if socialClient, err = social.New(channelID, channelSecret); err != nil {
		log.Println("Social SDK:", socialClient, " err:", err)
		return
	}

	// fs := http.FileServer(http.Dir("static"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))

	//For LINE login
	http.HandleFunc("/", accessToken)
	// http.HandleFunc("/", accessToken)
	// http.HandleFunc("/gotoauthOpenIDpage", gotoauthOpenIDpage)
	// http.HandleFunc("/gotoauthpage", gotoauthpage)
	// http.HandleFunc("/auth", auth)

	//provide by Heroku
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func accessToken(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	//Request for access token
	token, err := socialClient.GetAccessToken(fmt.Sprintf("%s", serverURL), code).Do()
	if err != nil {
		log.Println("RequestLoginToken err:", err)
		// token = fmt.Sprintf("%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// log.Println("access_token:", token.AccessToken, " refresh_token:", token.RefreshToken)

	lineToken := LineToken{token.AccessToken, token.RefreshToken}

	js, err := json.Marshal(lineToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
