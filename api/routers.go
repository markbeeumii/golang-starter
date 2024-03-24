package api

import (
	"github.com/gorilla/mux"
)

func SetRoutes(router *mux.Router) {
	router.HandleFunc("/users", CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", GetOneUser).Methods("GET")
	router.HandleFunc("/users", FindAll).Methods("GET")
}

// func StartServer() {
// 	log.Println("Starting server...")
// 	log.Fatal(http.ListenAndServe(":8000", nil))
// 	http.ListenAndServe(":8000", nil)
// }
