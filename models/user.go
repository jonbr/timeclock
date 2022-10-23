package models

import (
	"fmt"
	"timeclock/error"
	"timeclock/logger"

	"gorm.io/gorm"
	"github.com/gookit/goutil/dump"
)

type User struct {
	gorm.Model
	Name          	string `json:"name" gorm:"not null"`
	Username 		string `json:"username" gorm:"unique"`
	Email 			string `json:"email" gorm:"unique"`
	Password 		string `json:"password"`
	Administrator 	bool   `json:"administrator"`
}


func (u *User) GetUser(db *gorm.DB) *error.ErrorResp {
	fmt.Println("---GetUser---")
	dump.P(u.ID)
	if err := db.First(&u, u.ID).Error; err != nil {
		errResponse := error.New(error.WithDetails(err))
		return errResponse
	}

	return nil
}

func (u *User) GetUsers(db *gorm.DB) ([]User, *error.ErrorResp) {
	var users []User
	var errResponse *error.ErrorResp

	if result := db.Find(&users); result.Error != nil {
		errResponse = error.New(error.WithDetails(result.Error))
	}

	return users, errResponse
}

func (u *User) CreateUser(db *gorm.DB) *error.ErrorResp {
	if result := db.Create(&u); result.Error != nil {
		logger.Log.Error(result.Error)
		return error.New(error.WithDetails(result.Error))
	}


	return nil
}

/*func DeleteUser(userId string) error {
	if err := db.Delete(&User{}, userId); err != nil {
		return nil
	}

	return nil
}*/
