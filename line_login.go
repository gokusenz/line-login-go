package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	social "github.com/kkdai/line-social-sdk-go"
)

var nounce string
var state string

const lineLoginURL string = "https://access.line.me/oauth2/v2.1/authorize?response_type=code"

type LineToken struct {
	AccessToken  string
	RefreshToken string
}

func browse(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("login.tmpl"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println("Template err:", err)
	}
}

func gotoauthpage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("ParseForm() err: %v\n", err)
		return
	}
	chatbot := r.FormValue("chatbot")

	scope := "profile" //profile | openid | email
	state = GenerateNounce()
	nounce = GenerateNounce()
	redirectURL := fmt.Sprintf("%s/auth", serverURL)
	targetURL := socialClient.GetWebLoinURL(redirectURL, state, scope, social.AuthRequestOptions{Nounce: nounce, BotPrompt: chatbot, Prompt: "consent"})
	http.Redirect(w, r, targetURL, http.StatusSeeOther)
}

func gotoauthOpenIDpage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("ParseForm() err: %v\n", err)
		return
	}
	chatbot := r.FormValue("chatbot")

	scope := "profile openid" //profile | openid | email
	state = GenerateNounce()
	nounce = GenerateNounce()
	redirectURL := fmt.Sprintf("%s/auth", serverURL)
	targetURL := socialClient.GetWebLoinURL(redirectURL, state, scope, social.AuthRequestOptions{Nounce: nounce, BotPrompt: chatbot, Prompt: "consent"})
	http.Redirect(w, r, targetURL, http.StatusSeeOther)
}

func auth(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("ParseForm() err: %v\n", err)
		return
	}
	code := r.FormValue("code")
	inState := r.FormValue("state")
	//Check the state
	if strings.Compare(state, inState) != 0 {
		log.Println("State is not matching.")
		return
	}
	friendshipStatusChanged := r.FormValue("friendship_status_changed")
	log.Println("code:", code, " state:", state, "friend status:", friendshipStatusChanged)

	//Request for access token
	token, err := socialClient.GetAccessToken(fmt.Sprintf("%s/auth", serverURL), code).Do()
	if err != nil {
		log.Println("RequestLoginToken err:", err)
		return
	}

	log.Println("access_token:", token.AccessToken, " refresh_token:", token.RefreshToken)

	//Start to verify token and renew it.
	if result, err := socialClient.TokenVerify(token.AccessToken).Do(); err != nil {
		log.Println("TokenVerify err:", err, result)
		return
	}

	//Start to refresh token and renew it.
	if refresh, err := socialClient.RefreshToken(token.RefreshToken).Do(); err != nil {
		log.Println("RefreshToken err:", err, refresh)
		return
	}

	var payload *social.Payload
	if len(token.IDToken) == 0 {
		// User don't request openID, use access token to get usere profile
		log.Println(" token:", token, " AccessToken:", token.AccessToken)
		res, err := socialClient.GetUserProfile(token.AccessToken).Do()
		if err != nil {
			log.Println("GetUserProfile err:", err)
			return
		}
		payload = &social.Payload{
			Name:    res.DisplayName,
			Picture: res.PictureURL,
		}
	} else {
		//Decode token.IDToken to payload
		payload, err = token.DecodePayload(channelID)
		if err != nil {
			log.Println("DecodeIDToken err:", err)
			return
		}
	}

	//verify access token
	tmpl := template.Must(template.ParseFiles("login_success.tmpl"))
	if err := tmpl.Execute(w, payload); err != nil {
		log.Println("Template err:", err)
	}
}

func AccessToken(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("ParseForm() err: %v\n", err)
		return
	}
	code := r.FormValue("code")
	inState := r.FormValue("state")
	//Check the state
	if strings.Compare(state, inState) != 0 {
		log.Println("State is not matching.")
		return
	}
	//Request for access token
	token, err := socialClient.GetAccessToken(fmt.Sprintf("%s", serverURL), code).Do()
	if err != nil {
		log.Println("RequestLoginToken err:", err)
		return
	}

	log.Println("access_token:", token.AccessToken, " refresh_token:", token.RefreshToken)

	lineToken := LineToken{code, inState}

	js, err := json.Marshal(lineToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
