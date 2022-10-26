package models

import (
	apiError "timeclock/error"
	"timeclock/logger"

	"golang.org/x/crypto/bcrypt"
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


func (u *User) GetUser(db *gorm.DB) *apiError.ErrorResp {
	dump.P(u.ID)
	if err := db.First(&u, u.ID).Error; err != nil {
		errResponse := apiError.New(apiError.WithDetails(err))
		return errResponse
	}

	return nil
}

func (u *User) GetUsers(db *gorm.DB) ([]User, *apiError.ErrorResp) {
	var users []User
	var errResponse *apiError.ErrorResp

	if result := db.Find(&users); result.Error != nil {
		errResponse = apiError.New(apiError.WithDetails(result.Error))
	}

	return users, errResponse
}

func (user *User) GetUserByEmail(db *gorm.DB, email string) *apiError.ErrorResp {
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		errResponse := apiError.New(apiError.WithDetails(err))
		return errResponse
	}

	return nil
}

func (u *User) CreateUser(db *gorm.DB) *apiError.ErrorResp {
	if result := db.Create(&u); result.Error != nil {
		logger.Log.Error(result.Error)
		return apiError.New(apiError.WithDetails(result.Error))
	}


	return nil
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) *apiError.ErrorResp {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return apiError.New(apiError.WithDetails(err))
	}
	return nil
}

/*func DeleteUser(userId string) error {
	if err := db.Delete(&User{}, userId); err != nil {
		return nil
	}

	return nil
}*/
