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
	//"github.com/go-sql-driver/mysql"
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

// UserProjects join table for linking projects to users.
type UserProjects struct {
	UserID 			uint
	ProjectID 		uint
}

func (user *User) GetUser(db *gorm.DB) *apiError.ErrorResp {
	if err := db.Model(user).Preload("Projects").First(&user); err.Error != nil {
		if errors.Is(err.Error, gorm.ErrRecordNotFound) {
			return apiError.New(apiError.WithDetails(fmt.Sprintf("User with ID:%s not found", strconv.FormatUint(uint64(user.ID), 10))))
		}
		return apiError.New(apiError.WithDetails(err.Error))
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

func (user *User) CreateUser(db *gorm.DB) *apiError.ErrorResp {
	if result := db.Create(&user); result.Error != nil {
		logger.Log.Error(result.Error)
		return apiError.New(apiError.WithDetails(result.Error))
	}

	return nil
}

// TODO: change all models funcs. to return standard error instead of apiError,
// and have the controller only return an apiError obj.
func (user *User) UpdateUser(db *gorm.DB) error {
	// updating which projects user has relationship with.
	//dump.P(user)
	if len(user.Projects) > 0 {
		fmt.Println("updating which projects user has relationship with")
		if err := updateUserProjects(db, user.ID, user.Projects); err != nil {
			logger.Log.Error(err)
			return err
		}
	}

	result := db.Save(user);
	if result.Error != nil {
		logger.Log.Error(result.Error)
		return result.Error
	}
	if result.RowsAffected < 1 {	
		customError := fmt.Sprintf("Can't update user with id: %s it does not exists!", strconv.FormatUint(uint64(user.ID), 10))
		logger.Log.Error(customError)
		return errors.New(customError)	
	}

	return nil
}

func (user *User) DeleteUser(db *gorm.DB) *apiError.ErrorResp {
	if err := db.Delete(&user).Error; err != nil {
		logger.Log.Error(err.Error)
		return apiError.New(apiError.WithDetails(err.Error))
	}

	return nil
}

func updateUserProjects(db *gorm.DB, userID uint, projects []Project) error {
	userProject := UserProjects{
		UserID: 	userID,
		ProjectID: 	uint(2),
	}
	if result := db.Debug().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&userProject); result.Error != nil {
		dump.P(result.Error)
	}
	fmt.Println("project was not inserted, but instead removed!")
	//continue

	/*for _, project := range projects {
		userProject := UserProjects{
			UserID: 	userID,
			ProjectID: 	project.ID,
		}
		dump.P(userProject)

		if result := db.Debug().Create(userProject); result.Error != nil {
			// record already exists in db, will be removed instead of being inserted.
			if result.Error.(*mysql.MySQLError).Number == 1062 {
				if result := db.Debug().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&userProject); result.Error != nil {
					dump.P(result.Error)
				}
				fmt.Println("project was not inserted, but instead removed!")
				//continue
			}
			dump.P(result.Error)
		}
		fmt.Println("project attached to user")
	}*/

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
