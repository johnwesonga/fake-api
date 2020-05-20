package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/johnwesonga/fake-api/database"
)

const (
	STATIC_DIR = "/static/"
)

func main() {
	httpListenPort := os.Getenv("PORT")
	if httpListenPort == "" {
		httpListenPort = ":8080"
	}

	log.Println("Starting a fake API")

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix(STATIC_DIR).
		Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))

	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/users", getUsersHandler).Methods("GET")
	router.HandleFunc("/user/{name}", getUserHandler).Methods("GET")
	log.Printf("Listening on %s", httpListenPort)

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	log.Fatal(http.ListenAndServe(httpListenPort, loggedRouter))
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, req, "./static/index.html")
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("get all users API request...")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	users := database.GetAllUsers()
	respondWithJSON(w, http.StatusOK, users)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("get user API request...")
	name := mux.Vars(r)["name"]
	user, err := database.GetUser(name)
	if err != nil {
		//log.Println("User not found")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, user)

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
