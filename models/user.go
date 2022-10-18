package models

import (
  "timeclock/error"

  "gorm.io/gorm"
  //"github.com/gookit/goutil/dump"
)

type User struct {
  gorm.Model
  Name          string
  Administrator bool
}


func (u *User) GetUser(db *gorm.DB) *error.ErrorResp {
  if err := db.First(&u, u.ID).Error; err != nil {
    errResponse := error.New(error.WithDetails(err))
    return errResponse
  }

  return nil
}

/*func DeleteUser(userId string) error {
  if err := db.Delete(&User{}, userId); err != nil {
    return nil
  }

  return nil
}*/
