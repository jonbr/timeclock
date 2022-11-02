package models

import (
	"fmt"

	"timeclock/error"

	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ID 			uint
	Name        string `json:"name" gorm:"unique"`
	Description string `json:"description"`
}

func (p *Project) GetProject(db *gorm.DB) ([]Project, *error.ErrorResp) {
	var projects []Project
	var errResponse *error.ErrorResp

	if p.ID != 0 {
		fmt.Println("p.ID: ", p.ID)
		if err := db.First(&projects, p.ID).Error; err != nil {
			errResponse = error.New(error.WithDetails(err))
		}
	} else {
		if result := db.Find(&projects); result.Error != nil {
			errResponse = error.New(error.WithDetails(result.Error))
		}
	}

	return projects, errResponse
}

func (p *Project) CreateProject(db *gorm.DB) *error.ErrorResp {
	if err := db.Create(&p).Error; err != nil {
		errResponse := error.New(error.WithDetails(err))
		return errResponse
	}

	return nil
}
