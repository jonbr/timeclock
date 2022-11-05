package models

import (
	"fmt"

	"timeclock/error"

	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique"`
	Description string `json:"description"`
}

func (project *Project) GetProject(db *gorm.DB) ([]Project, *error.ErrorResp) {
	var projects []Project
	var errResponse *error.ErrorResp

	if project.ID != 0 {
		fmt.Println("project.ID: ", project.ID)
		if err := db.First(&projects, project.ID).Error; err != nil {
			errResponse = error.New(error.WithDetails(err))
		}
	} else {
		if result := db.Find(&projects); result.Error != nil {
			errResponse = error.New(error.WithDetails(result.Error))
		}
	}

	return projects, errResponse
}

func (project *Project) CreateProject(db *gorm.DB) *error.ErrorResp {
	if err := db.Create(&project).Error; err != nil {
		errResponse := error.New(error.WithDetails(err))
		return errResponse
	}

	return nil
}
