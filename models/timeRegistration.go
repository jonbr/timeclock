package models

import "time"

type TimeRegister struct {
	ID        uint
	ClockIn   *time.Time
	ClockOut  *time.Time
	UserID    uint 			`gorm:"not null"`
	User      User
	ProjectID uint 			`gorm:"not null"`
	Project   Project
}