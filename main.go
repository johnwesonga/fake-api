package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

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

	go func() {
		log.Fatal(http.ListenAndServe(httpListenPort, loggedRouter))
	}()

	log.Println("Starting a fake API")
	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		log.Print("Got SIGINT...")
	case syscall.SIGTERM:
		log.Print("Got SIGTERM...")
	}
	log.Print("The service is shutting down...")
	log.Print("Done")

}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	// fmt.Fprintf(w, "Welcome to the Fake API")
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
