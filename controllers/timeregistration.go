package controllers

import (
	"net/http"

	"gorm.io/gorm"
)

func Timeregistrationclockin(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	}
}

func Timeregistrationclockout(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	}
}

/*func Timeregistrationclockin(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
  	// get user
  	/*var err *error.ErrorResp
  	var user models.User
  	//var users []models.User
  	userId := mux.Vars(r)["id"] // string
  	uintId, err_new := strconv.ParseUint(userId, 10, 32)
  	if err_new != nil {
		fmt.Printf("%T, %v\n", uintId, uintId)
	}
  	user, err = getUser(uint(uintId))
  	if (models.User{} == user) {
    	fmt.Println("No User found, not possible to clockIn!")
    	fmt.Println(err)
  	}
  	// create clockIn record
  	var timeStamp = time.Now()
  	database.Instance.Create(&models.TimeRegister{ClockIn: &timeStamp, UserID: uint(uintId)})

  	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  	w.WriteHeader(http.StatusOK)
  }
}*/