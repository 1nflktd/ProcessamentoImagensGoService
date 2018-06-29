package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func changeBrightness(w http.ResponseWriter, r *http.Request) {
	api := HttpApiNew(w, r)
	if err := api.Init(); err != nil {
		return
	}
	api.Image.changeBrightness(api.Parameters.Intensity)
	if err := api.writeImage(); err != nil {
		return
	}
}

func blurImage(w http.ResponseWriter, r *http.Request) {
	api := HttpApiNew(w, r)
	if err := api.Init(); err != nil {
		return
	}
	api.Image.blur(api.Parameters.Intensity)
	if err := api.writeImage(); err != nil {
		return
	}
}

func sharpenImage(w http.ResponseWriter, r *http.Request) {
	api := HttpApiNew(w, r)
	if err := api.Init(); err != nil {
		return
	}
	api.Image.sharpen(api.Parameters.Intensity)
	if err := api.writeImage(); err != nil {
		return
	}
}

func testConnection(w http.ResponseWriter, r *http.Request) {
	log.Printf("Test solicitation received")

	var retPayload JsonResponse
	retPayload.PayloadBase64 = "Teste";
	if err := json.NewEncoder(w).Encode(retPayload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error %s\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	log.Printf("Test solicitation processed")
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := mux.NewRouter()
	router.HandleFunc("/brightness", changeBrightness).Methods("POST")
	router.HandleFunc("/sharpen", sharpenImage).Methods("POST")
	router.HandleFunc("/blur", blurImage).Methods("POST")
	router.HandleFunc("/test", testConnection).Methods("GET")
	router.HandleFunc("/", mainHandler).Methods("GET")
	log.Printf("Listening on :%s...\n", port)
	log.Fatal(http.ListenAndServe(":" + port, router))
}
