package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"timeclock/error"
	"timeclock/logger"
	"timeclock/models"
	//"timeclock/utils"

	//"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	//"github.com/gookit/goutil/dump"
)


func GetAll(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		user := &models.User{}
		users, err := user.GetAll(db)
		if err != nil {
			fmt.Println(err)
		} 

		json.NewEncoder(w).Encode(users)
		w.WriteHeader(http.StatusOK)
	}
}

// TODO: merge RegisterUser and CreateUser into one function.
func RegisterUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("---RegisterUser---")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var user models.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(error.New(error.WithDetails(err)))
			return
		}
		if err := user.HashPassword(user.Password); err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(error.New(error.WithDetails(err)))
			return
		}
		if errResp := user.CreateUser(db); errResp != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errResp)
			return
		}
		
		logger.Log.WithFields(logrus.Fields{
			"host":     r.URL.Host,
			"path":     r.URL.Path,
			"header":   r.Header,
			// as you can see, there is a lot the logger can do for us
			// however "body": r.Body will not work, and always log an empty string!
			//"body":     req
			// this is why we'll log our crated struct instead.
		}).Info(user)

		json.NewEncoder(w).Encode(user)
		w.WriteHeader(http.StatusOK)
	}
}

func GetUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		u := &models.User{}
		users, err := u.GetUsers(db)
		if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
			return
		}

		if !u.Administrator {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(error.New(error.WithDetails(fmt.Sprintf("User %s does not have sufficient privledges to create a project!", u.Name))))
			return
		}

		logger.Log.WithFields(logrus.Fields{
			"host":     r.URL.Host,
			"path":     r.URL.Path,
			"header":   r.Header,
			// as you can see, there is a lot the logger can do for us
			// however "body": r.Body will not work, and always log an empty string!
			//"body":     req
			// this is why we'll log our crated struct instead.
		}).Info()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	}
}

func GetUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		/*userId, err := utils.CastStringToUint(mux.Vars(r)["id"])
		if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		fmt.Println("userId: ", userId)*/

		u := &models.User{}
		u.ID = 2 //userId
		if errResp := u.GetUser(db); errResp != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errResp)
			return
		}

		logger.Log.WithFields(logrus.Fields{
			"host":     r.URL.Host,
			"path":     r.URL.Path,
			"header":   r.Header,
			// as you can see, there is a lot the logger can do for us
			// however "body": r.Body will not work, and always log an empty string!
			//"body":     req
			// this is why we'll log our crated struct instead.
		}).Info()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(u)
	}
}

func CreateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var user models.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(error.New(error.WithDetails(err)))
			return
		}
		if errResp := user.CreateUser(db); errResp != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errResp)
			return
		}

		logger.Log.WithFields(logrus.Fields{
			"host":     r.URL.Host,
			"path":     r.URL.Path,
			"header":   r.Header,
			// as you can see, there is a lot the logger can do for us
			// however "body": r.Body will not work, and always log an empty string!
			//"body":     req
			// this is why we'll log our crated struct instead.
		}).Info(user)

		json.NewEncoder(w).Encode(user)
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	}
}