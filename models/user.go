package models

import (
	"fmt"
	"timeclock/error"

	"gorm.io/gorm"
	"github.com/gookit/goutil/dump"
)

type User struct {
	gorm.Model
	Name          string
	Administrator bool
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

/*func DeleteUser(userId string) error {
	if err := db.Delete(&User{}, userId); err != nil {
		return nil
	}

	return nil
}*/
