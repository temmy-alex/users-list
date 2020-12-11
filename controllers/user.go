package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"users-list/models"
	userRepository "users-list/repository/user"
	"users-list/utils"

	"github.com/gorilla/mux"
)

type Controllers struct{}

var users []models.User

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (c Controllers) GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error

		users = []models.User{}
		userRepo := userRepository.UserRepository{}
		users, err := userRepo.GetUsers(db, user, users)

		if err != nil {
			error.Message = "Server Error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, users)
	}
}

func (c Controllers) GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error

		params := mux.Vars(r)

		users = []models.User{}
		userRepo := userRepository.UserRepository{}

		id, _ := strconv.Atoi(params["id"])
		user, err := userRepo.GetUser(db, user, id)

		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "Not Found"
				utils.SendError(w, http.StatusNotFound, error)
				return
			} else {
				error.Message = "Server Error"
				utils.SendError(w, http.StatusInternalServerError, error)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, user)
	}
}

func (c Controllers) AddUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var userID int
		var error models.Error

		json.NewDecoder(r.Body).Decode(&user)

		if user.Name == "" || user.Email == "" {
			error.Message = "Enter missing fields."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		userRepo := userRepository.UserRepository{}
		userID, err := userRepo.AddUser(db, user)

		if err != nil {
			error.Message = "Server Error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		w.Header().Set("Content-Type", "text/plan")
		utils.SendSuccess(w, userID)
	}
}

func (c Controllers) UpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error

		json.NewDecoder(r.Body).Decode(&user)

		if user.ID == 0 || user.Name == "" || user.Email == "" {
			error.Message = "All fields are required"
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		userRepo := userRepository.UserRepository{}
		rowsUpdated, err := userRepo.UpdateUser(db, user)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.SendSuccess(w, rowsUpdated)
	}
}

func (c Controllers) RemoveUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		params := mux.Vars(r)
		userRepo := userRepository.UserRepository{}
		id, _ := strconv.Atoi(params["id"])

		rowsDeleted, err := userRepo.RemoveUser(db, id)

		if err != nil {
			error.Message = "Server error."
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		if rowsDeleted == 0 {
			error.Message = "Not Found"
			utils.SendError(w, http.StatusNotFound, error)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.SendSuccess(w, rowsDeleted)
	}
}
