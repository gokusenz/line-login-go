package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

//LINE Login related configuration
var channelID, channelSecret string

//LINE MessageAPI related configuration
var serverURL string
var botToken, botSecret string

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

	client := &http.Client{}
	ctx := context.Background()
	data := tokenRespone{}
	rc := NewClient(client, code)
	body, err := rc.Do(ctx)
	if err != nil {
		return
	}

	json.Unmarshal([]byte(body), &data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)

}

type tokenRespone struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}
