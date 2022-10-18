package models


import (
	"timeclock/error"

	//"github.com/gookit/goutil/dump"
	"gorm.io/gorm"
)


type Project struct {
	gorm.Model
	Name 		string
	Description string 
	UserID    	uint
  	User      	User
}

func (p *Project) GetProject(db *gorm.DB) ([]Project, *error.ErrorResp) {
	var projects []Project
	var errResponse *error.ErrorResp

	//dump.P(p)

	if p.ID != 0 {
		if err := db.First(&projects, p.ID).Error; err != nil {
	    errResponse = error.New(error.WithDetails(err))
	    //return errResponse
	  }
	} else {
		if result := db.Find(&projects); result.Error != nil {
			errResponse = error.New(error.WithDetails(result.Error))
			//return nil, errResponse
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
