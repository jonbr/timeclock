package models

import (
	"fmt"
	"strconv"

	"timeclock/error"
	"timeclock/logger"

	"gorm.io/gorm"
	"github.com/gookit/goutil/dump"
)

type Project struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique"`
	Description string `json:"description"`
}

// GetProject returns idividual record if project.id is attached to the project object
// otherwise returns all project records.
func (project *Project) GetProject(db *gorm.DB) ([]Project, *error.ErrorResp) {
	var projects []Project

	if project.ID != 0 {
		if err := db.First(&projects, project.ID).Error; err != nil {
			return nil, error.New(error.WithDetails(err))
		}
	} else {
		if result := db.Find(&projects); result.Error != nil {
			return nil, error.New(error.WithDetails(result.Error))
		}
	}

	return projects, nil
}

func (project *Project) CreateProject(db *gorm.DB) *error.ErrorResp {
	if err := db.Create(&project).Error; err != nil {
		return error.New(error.WithDetails(err))
	}

	return nil
}

// TODO: look into some validation before and after db action.
func (project *Project) UpdateProject(db *gorm.DB) *error.ErrorResp {
	result := db.Model(&project).Updates(Project{Name: project.Name, Description: project.Description})
	dump.P(result)
	if result.Error != nil {
		return error.New(error.WithDetails(result.Error))
	}
	if result.RowsAffected < 1 {
		customError := fmt.Sprintf("Can't update project with id: %s it does not exists!", strconv.FormatUint(uint64(project.ID), 10))
		logger.Log.Error(customError)
		return error.New(error.WithDetails(customError))
	}

	return nil
}

// TODO: implement logging and returing to user error handling as is done here.
func (project *Project) DeleteProject(db *gorm.DB) *error.ErrorResp {
	result := db.Delete(&project);
	if result.Error != nil {
		return error.New(error.WithDetails(result.Error))
	}
	if result.RowsAffected < 1 {
		customError := fmt.Sprintf("Can't delete project with id: %s it does not exists!", strconv.FormatUint(uint64(project.ID), 10))
		logger.Log.Error(customError)
		return error.New(error.WithDetails(customError))
    }

	return nil
}
