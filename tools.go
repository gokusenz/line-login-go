package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type Payload struct {
	Iss     string `json:"iss"`
	Sub     string `json:"sub"`
	Aud     string `json:"aud"`
	Exp     int    `json:"exp"`
	Iat     int    `json:"iat"`
	Nonce   string `json:"nonce"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateNounce() string {
	return b64.StdEncoding.EncodeToString([]byte(randStringRunes(8)))
}

func DecodeIDToken(idToken string, channelID string) (*Payload, error) {
	splitToken := strings.Split(idToken, ".")
	if len(splitToken) < 3 {
		log.Println("Error: idToken size is wrong, size=", len(splitToken))
		return nil, fmt.Errorf("Error: idToken size is wrong. \n")
	}
	header, payload, signature := splitToken[0], splitToken[1], splitToken[2]
	log.Println("result:", header, payload, signature)

	log.Println("side of payload=", len(payload))
	payload = base64Decode(payload)
	log.Println("side of payload=", len(payload), payload)
	bPayload, err := b64.StdEncoding.DecodeString(payload)
	if err != nil {
		log.Println("base64 decode err:", err)
		return nil, fmt.Errorf("Error: base64 decode. \n")
	}
	log.Println("base64 decode succeess:", string(bPayload))

	retPayload := &Payload{}
	if err := json.Unmarshal(bPayload, retPayload); err != nil {
		return nil, fmt.Errorf("json unmarshal error, %v. \n", err)
	}

	// payload verification
	if strings.Compare(retPayload.Iss, "https://access.line.me") != 0 {
		return nil, fmt.Errorf("Payload verification wrong. Wrong issue organization. \n")
	}
	if strings.Compare(retPayload.Aud, channelID) != 0 {
		return nil, fmt.Errorf("Payload verification wrong. Wrong audience. \n")
	}

	return retPayload, nil
}

func base64Decode(payload string) string {
	rem := len(payload) % 4
	log.Println("rem of payload=", rem)
	if rem > 0 {
		i := 4 - rem
		for ; i > 0; i-- {
			payload = payload + "="
		}
	}
	return payload
}
