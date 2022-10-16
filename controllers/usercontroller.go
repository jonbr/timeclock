package controllers

import (
	"timeclock/database"
	"timeclock/models"

	"fmt"
	"encoding/json"
	"net/http"
	//"path"

	"github.com/gorilla/mux"
	//"github.com/gookit/goutil/dump"
)

/*func getDataFromRequest(req *http.Request) (string, string) {
   // Extract the user-provided data for the new clichePair
   req.ParseForm()
   form := req.Form
   cliche := form["cliche"][0]    // 1st and only member of a list
   counter := form["counter"][0]  // ditto
   return cliche, counter
}*/

// serves index file
/*func Home(w http.ResponseWriter, r *http.Request) {
		//var newVar := json.NewDecoder(r.Body).Decode("btn")

		dump.P(r.Form())
		//fmt.Println("index: ", index)
    //r.NewRoute().PathPrefix("./web/")
    //r.PathPrefix("./web/")//.Handler(http.FileServer(http.Dir("./web/")))
    p := path.Dir("./web/")
    // set header
    w.Header().Set("Content-type", "text/html")
    //w.Header().Set("Content-Type", "text/javascript")
    http.ServeFile(w, r, p)
}*/

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	database.Instance.Create(&user)
	json.NewEncoder(w).Encode(user)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var users []models.User
	userId := mux.Vars(r)["id"]
	if userId != "" {
		if err := database.Instance.First(&user, userId).Error; err != nil {
    		fmt.Println(err)
    		return
  		}
  		json.NewEncoder(w).Encode(user)
	} else {
		result := database.Instance.Find(&users)
		if result.Error != nil {
			fmt.Println(result.Error)
		}
		json.NewEncoder(w).Encode(users)
	}

  	w.Header().Set("Content-Type", "application/json")
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]

  	if err := models.DeleteUser(userId); err != nil {
  		fmt.Println(err)
  	}	
}

func getUser(userID uint) models.User {
  var user models.User
  if err := database.Instance.First(&user, userID).Error; err != nil {
    fmt.Println(err)
  }

  return user
}