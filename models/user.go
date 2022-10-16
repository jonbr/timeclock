package models

import "gorm.io/gorm"

var db *gorm.DB

type User struct {
  gorm.Model
  Name  string
  Administrator bool
}

func DeleteUser(userId string) error {
  if err := db.Delete(&User{}, userId); err != nil {
    return nil
  }

  return nil
}
