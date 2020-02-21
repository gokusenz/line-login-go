package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	social "github.com/kkdai/line-social-sdk-go"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

//LINE Login related configuration
var channelID, channelSecret string

//LINE MessageAPI related configuration
var serverURL string
var botToken, botSecret string
var socialClient *social.Client

func main() {
	var err error
	serverURL = "https://line-login-soical.herokuapp.com/"
	channelID = "1653859637"
	channelSecret = "350c50d8c1e4435726d64450f45142c3"

	if socialClient, err = social.New(channelID, channelSecret); err != nil {
		log.Println("Social SDK:", socialClient, " err:", err)
		return
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

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
