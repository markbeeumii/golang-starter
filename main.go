package main

import (
	"log"
	"myproject/api"
	"myproject/config"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config.InitDB()
	router := mux.NewRouter()
	api.SetRoutes(router)
	log.Println("Starting server...8000")
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		log.Fatal(err)
	}
}
