package controllers

import (
	"timeclock/database"
	"timeclock/models"

	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gookit/goutil/dump"
)


func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var project models.Project
	json.NewDecoder(r.Body).Decode(&project)
	dump.P(project)

	userId := mux.Vars(r)["id"]
	fmt.Println("userId: ", userId)
	uintId, err := strconv.ParseUint(userId, 10, 32)
  	if err != nil {
		fmt.Printf("%T, %v\n", uintId, uintId)
	}
  	user := getUser(uint(uintId))
  	if (models.User{} == user) {
    	fmt.Println("No User found, not possible to create a new project!")
  	}
  	if (!user.Administrator) {
  		log.Println(fmt.Sprintf("UserId %s does not have sufficient privledges to create a project!", user.ID))
  	}

  	project.UserID = uint(uintId)

	database.Instance.Create(&project)
	json.NewEncoder(w).Encode(project)
}