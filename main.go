package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"users-list/controllers"
	"users-list/driver"
	"users-list/models"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

var users []models.User
var db *sql.DB

func init() {
	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db = driver.ConnectDB()
	controller := controllers.Controllers{}
	router := mux.NewRouter()

	router.HandleFunc("/users", controller.GetUsers(db)).Methods("GET")
	router.HandleFunc("/users/{id}", controller.GetUser(db)).Methods("GET")
	router.HandleFunc("/users", controller.AddUser(db)).Methods("POST")
	router.HandleFunc("/users", controller.UpdateUser(db)).Methods("PUT")
	router.HandleFunc("/users/{id}", controller.RemoveUser(db)).Methods("DELETE")

	fmt.Println("Server is running at port 8000")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}))(router)))
}
