package models

import "time"

type TimeRegister struct {
  ID        uint
  ClockIn   *time.Time
  ClockOut  *time.Time
  UserID    uint
  User      User
  ProjectID uint
  Project   Project
}