package controllers

import (
	"fmt"
	"encoding/json"
	"net/http"
	//"strconv"

	"timeclock/error"
	"timeclock/logger"
	"timeclock/models"
	"timeclock/utils"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"github.com/gookit/goutil/dump"
)

func GetProjects(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &models.Project{}
		projects, err := p.GetProject(db)
		if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
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
		json.NewEncoder(w).Encode(projects)
	}
}

func GetProject(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		dump.P(mux.Vars(r))
		uintParams, err := utils.CastStringToUint(mux.Vars(r))

		if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		p := &models.Project{}
		p.ID = uintParams[1];
		projects, errResp := p.GetProject(db)
		if errResp != nil {
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
			// this is why we'll log our crated struct instead.
		}).Info()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(projects[0])
	}
}

// TODO: use as reference for logging!	
func CreateProject(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var p models.Project
		json.NewDecoder(r.Body).Decode(&p)

		u := &models.User{}
		//u.ID = p.UserID
		/*if errResp := u.GetUser(db); errResp != nil {
			logger.Log.Error(fmt.Sprintf("User with ID: %s not found!", strconv.FormatUint(uint64(p.UserID), 10)))
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(error.New(error.WithDetails(fmt.Sprintf("User with ID: %s not found!", strconv.FormatUint(uint64(p.UserID), 10)))))
			return
		}*/
		if !u.Administrator {
			logger.Log.Error(fmt.Sprintf("User %s does not have sufficient privledges to create a project!", u.Name))
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(error.New(error.WithDetails(fmt.Sprintf("User %s does not have sufficient privledges to create a project!", u.Name))))
			return
		}
			
		//p.UserID = uint(u.ID)
		if errResp := p.CreateProject(db); errResp != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
		}).Info(p)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(p)
	}
}

func UpdateProject(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteProject(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	}
}
