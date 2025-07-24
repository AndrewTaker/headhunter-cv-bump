package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	_ "github.com/mattn/go-sqlite3"
)

var (
	clientID, clientSecret, resumeID, redirectURL string
	isProd                                        bool
	serverPort, serverHost, serverHTTP            string
	templates                                     *template.Template
	client                                        *http.Client
	db                                            *sql.DB
	sessionManager                                *scs.SessionManager
)

func main() {
	log.Println("change")
	clientID = os.Getenv("HH_CLIENT_ID")
	clientSecret = os.Getenv("HH_CLIENT_SECRET")
	redirectURL = os.Getenv("HH_REDIRECT_URL")

	serverHTTP = os.Args[1]
	serverHost = os.Args[2]
	serverPort = os.Args[3]
	isProd, _ = strconv.ParseBool(os.Args[4])

	if clientID == "" || clientSecret == "" || redirectURL == "" {
		log.Fatal("main: no credentials provided")
	}

	var err error
	db, err = db_init()
	if err != nil {
		log.Fatal("main: ", err)
	}

	sessionManager = scs.New()
	sessionManager.Lifetime = 1 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = isProd
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Store = memstore.New()

	client = &http.Client{Timeout: 10 * time.Second}
	templates = template.Must(
		template.New("base").
			Funcs(template.FuncMap{
				"formatTime": func(t HHTime) string {
					return t.Format(timeLayout)
				},
			}).
			ParseFiles(
				"templates/base.html",
				"templates/header.html",
				"templates/info.html",
				"templates/modal.html",
				"templates/toggle-switch.html",
			),
	)

	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/invalidate", invalidateUserData)
	http.HandleFunc("/auth/callback", callback)
	http.HandleFunc("/get-resumes", updateResumesOnDemand)
	http.HandleFunc("/open-modal", openModal)
	http.HandleFunc("/close-modal", closeModal)
	http.HandleFunc("POST /toggle-schedule/{id}", toggleResume)

	log.Printf("server starting %s://%s:%s", serverHTTP, serverHost, serverPort)
	err = http.ListenAndServe(":"+serverPort, sessionManager.LoadAndSave(http.DefaultServeMux))
	if err != nil {
		log.Fatal("main: couldnt start server", err)
	}

}
