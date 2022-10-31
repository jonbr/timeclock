package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"timeclock/error"
	"timeclock/logger"
	"timeclock/models"
	"timeclock/utils"

	"github.com/gookit/goutil/dump"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		user := &models.User{}
		users, err := user.GetUsers(db)
		if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
			return
		}

		logger.Log.WithFields(logrus.Fields{
			"host":   r.URL.Host,
			"path":   r.URL.Path,
			"header": r.Header,
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
		userId, err := utils.CastStringToUint(mux.Vars(r))
		if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		u := &models.User{}
		u.ID = userId[0]
		if errResp := u.GetUser(db); errResp != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errResp)
			return
		}

		logger.Log.WithFields(logrus.Fields{
			"host":   r.URL.Host,
			"path":   r.URL.Path,
			"header": r.Header,
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
			"host":   r.URL.Host,
			"path":   r.URL.Path,
			"header": r.Header,
			// as you can see, there is a lot the logger can do for us
			// however "body": r.Body will not work, and always log an empty string!
			//"body":     req
			// this is why we'll log our crated struct instead.
		}).Info(user)

		json.NewEncoder(w).Encode(user)
		w.WriteHeader(http.StatusOK)
	}
}

// TODO: need to refactor and finish func.
func UpdateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Println("---UpdateUser---")

		userId, err := utils.CastStringToUint(mux.Vars(r))
		if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		user := &models.User{
			//ID: 2,
			Name: "jinzhu",
			/*Projects: []models.Project{
				ID: 2,
			},*/
		}
		/*user.ID = userId[0]
		fmt.Println("userId: ", userId[0])
		if errResp := user.GetUser(db); errResp != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errResp)
			return
		}*/
		fmt.Println("userId: ", userId)
		dump.P(user)
		/*if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(error.New(error.WithDetails(err)))
			return
		}*/

		if errResp := user.UpdateUser(db); errResp != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errResp)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

func DeleteUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Endpoint: DeleteUser not yet been implemented!")
	}
}
