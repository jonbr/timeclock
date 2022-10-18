package controllers

import (
	"timeclock/error"
	"timeclock/models"
	"timeclock/utils"

	"fmt"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	//"github.com/gookit/goutil/dump"
)

func GetProjects(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &models.Project{}
	  	projects, err := p.GetProject(db)
	  	if err != nil {
	  		w.WriteHeader(http.StatusNotFound)
	  		json.NewEncoder(w).Encode(err)
	  		return
	  	}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(projects)
	}
}


func GetProject(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		tmpVar := mux.Vars(r)["id"]
		fmt.Println("tmpVar: ", tmpVar)

		projectId, err := utils.CastStringToUint(mux.Vars(r)["id"])
		fmt.Println("projectId: ", projectId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
	  		json.NewEncoder(w).Encode(err)
	  		return
		}

	  	p := &models.Project{}
	  	p.ID = projectId;
	  	projects, errResp := p.GetProject(db)
	  	if errResp != nil {
	  		w.WriteHeader(http.StatusNotFound)
	  		json.NewEncoder(w).Encode(errResp)
	  		return
	  	}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(projects)
	}
}

func CreateProject(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("---CreateProject---")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		
		//var errResp *error.ErrorResp
		var p models.Project
		json.NewDecoder(r.Body).Decode(&p)

		
		//u.ID = utils.CastStringToUint(mux.Vars(r)["userId"])

		userId, err := utils.CastStringToUint(mux.Vars(r)["userId"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
	  		json.NewEncoder(w).Encode(err)
	  		return
		}

		u := &models.User{}
		u.ID = userId
	  	if errResp := u.GetUser(db); errResp != nil {
  			w.WriteHeader(http.StatusNotFound)
  			json.NewEncoder(w).Encode(errResp)
  			return
	  	} 
	  	if !u.Administrator {
  			w.WriteHeader(http.StatusUnauthorized)
  			json.NewEncoder(w).Encode(error.New(error.WithDetails(fmt.Sprintf("User %s does not have sufficient privledges to create a project!", u.Name))))
  			return
  		}
	  		
  		p.UserID = uint(u.ID)
  		if errResp := p.CreateProject(db); errResp != nil {
  			w.WriteHeader(http.StatusInternalServerError)
  			json.NewEncoder(w).Encode(errResp)
  			return
  		}

  		w.WriteHeader(http.StatusOK)
  		json.NewEncoder(w).Encode(p)
	}
}
