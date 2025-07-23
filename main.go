package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	_ "github.com/mattn/go-sqlite3"
)

var (
	clientID, clientSecret, refreshToken, resumeID string
	serverPort, serverHost, serverHTTP             string
	templates                                      *template.Template
	client                                         *http.Client
	db                                             *sql.DB
	sessionManager                                 *scs.SessionManager
)

func main() {
	var err error
	db, err = db_init()
	if err != nil {
		log.Fatal(err)
	}

	sessionManager = scs.New()
	sessionManager.Lifetime = 1 * time.Hour
	sessionManager.Cookie.Persist = true
	// sessionManager.Cookie.Secure = false
	// sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Store = memstore.New()

	client = &http.Client{Timeout: 10 * time.Second}
	templates = template.Must(template.New("base").ParseFiles(
		"templates/base.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/index.html",
		"templates/page.html",
		"templates/toggle-switch.html",
	))

	clientID = os.Getenv("HH_CLIENT_ID")
	clientSecret = os.Getenv("HH_CLIENT_SECRET")

	serverHTTP = os.Args[1]
	serverHost = os.Args[2]
	serverPort = os.Args[3]

	if clientID == "" || clientSecret == "" {
		log.Fatal("no credentials provided")
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/page", page)
	http.HandleFunc("/login", login)
	http.HandleFunc("/auth/callback", callback)
	http.HandleFunc("/get-resumes", updateResumesOnDemand)
	http.HandleFunc("POST /toggle-schedule/{id}", toggleResume)

	log.Printf("server starting %s://%s:%s", serverHTTP, serverHost, serverPort)
	err = http.ListenAndServe(":"+serverPort, sessionManager.LoadAndSave(http.DefaultServeMux))
	if err != nil {
		log.Fatal("couldnt stop server", err)
	}

}
