package models

import (
	"errors"
	"fmt"
	"strconv"

	apiError "timeclock/error"
	"timeclock/logger"

	"github.com/gookit/goutil/dump"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name          string    `json:"name" gorm:"not null"`
	Username      string    `json:"username" gorm:"unique"`
	Email         string    `json:"email" gorm:"unique"`
	Password      string    `json:"password"`
	Administrator bool      `json:"administrator"`
	Projects      []Project `json:"projects" gorm:"many2many:user_Projects;"`
}

func (user *User) GetUserInternal(db *gorm.DB) *apiError.ErrorResp {
	if result := db.First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return apiError.New(apiError.WithDetails(fmt.Sprintf("User with ID:%s not found", strconv.FormatUint(uint64(user.ID), 10))))
		}
		return apiError.New(apiError.WithDetails(result.Error))
	}

	return nil
}

func (user *User) GetUser(db *gorm.DB) *apiError.ErrorResp {
	dump.P(user.ID)
	if err := db.Model(user).Preload("Projects").Find(&user).Error; err != nil {
		errResponse := apiError.New(apiError.WithDetails(err))
		return errResponse
	}

	return nil
}

func (user *User) GetUsers(db *gorm.DB) ([]User, *apiError.ErrorResp) {
	var users []User
	var errResponse *apiError.ErrorResp

	if err := db.Model(user).Preload("Projects").Find(&users).Error; err != nil {
		errResponse = apiError.New(apiError.WithDetails(err))
	}

	return users, errResponse
}

func (user *User) GetUserByEmail(db *gorm.DB) *apiError.ErrorResp {
	if err := db.Where("email = ?", user.Email).First(&user).Error; err != nil {
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

func (user *User) UpdateUser(db *gorm.DB) *apiError.ErrorResp {
	fmt.Println("---models.UpdateUser---")
	if err := db.Model(user).Updates(user); err != nil {
		logger.Log.Error(err)
		return apiError.New(apiError.WithDetails(err))
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
