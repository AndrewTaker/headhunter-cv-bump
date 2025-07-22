package main

import (
	"log"
	"net/http"
	"os"
)

var (
	clientID, clientSecret, refreshToken, resumeID string
	serverPort, serverHost, serverHTTP             string
)

const (
	tokenURL = "https://hh.ru/oauth/token"
	bumpURL  = "https://api.hh.ru/resumes/%s/publish"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func main() {
	clientID = os.Getenv("HH_CLIENT_ID")
	clientSecret = os.Getenv("HH_CLIENT_SECRET")

	serverHTTP = os.Args[1]
	serverHost = os.Args[2]
	serverPort = os.Args[3]

	if clientID == "" || clientSecret == "" {
		log.Fatal("no credentials provided")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("QWEQWEQ"))
	})

	log.Printf("server starting %s://%s:%s", serverHTTP, serverHost, serverPort)
	err := http.ListenAndServe(":"+serverPort, nil)
	if err != nil {
		log.Fatal("couldnt stop server", err)
	}

}
