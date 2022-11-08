package models

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	apiError "timeclock/error"
	"timeclock/logger"

	"gorm.io/gorm"
)

type TimeRegister struct {
	ID                    uint
	ClockInTime           *time.Time
	ClockOutTime          *time.Time
	TimeDurationInMinutes float64
	UserID                uint `gorm:"not null"`
	User                  User
	ProjectID             uint `gorm:"not null"`
	Project               Project
}

func (tr *TimeRegister) ClockIn(db *gorm.DB) *apiError.ErrorResp {
	var timeStamp = time.Now()
	tr.ClockInTime = &timeStamp

	lastTimeRegister, err := getLastTimeRegisterRecord(db, tr.UserID)
	if err != nil {
		logger.Log.Error(err)
		return apiError.New(apiError.WithDetails(err))
	}
	// check if user is already clocked-in
	if lastTimeRegister.ClockInTime != nil && lastTimeRegister.ClockOutTime == nil {
		logger.Log.Error(fmt.Sprintf("UserId: %s already clocked-in", strconv.FormatUint(uint64(tr.UserID), 10)))
		return apiError.New(apiError.WithDetails(fmt.Sprintf("UserId: %s already clocked-in", strconv.FormatUint(uint64(tr.UserID), 10))))
	}

	if result := db.Create(&tr); result.Error != nil {
		logger.Log.Error(result.Error)
		return apiError.New(apiError.WithDetails(result.Error))
	}

	return nil
}

func (tr *TimeRegister) ClockOut(db *gorm.DB) *apiError.ErrorResp {
	var timeStamp = time.Now()
	clockOutTime := &timeStamp

	lastTimeRegister, err := getLastTimeRegisterRecord(db, tr.UserID)
	if err != nil {
		logger.Log.Error(err)
		return apiError.New(apiError.WithDetails(err))
	}
	// check if last record clocked-out is set or not
	if lastTimeRegister.ClockOutTime != nil {
		logger.Log.Error(fmt.Sprintf("UserId: %s already clocked-out", strconv.FormatUint(uint64(tr.UserID), 10)))
		return apiError.New(apiError.WithDetails(fmt.Sprintf("UserId: %s already clocked-out", strconv.FormatUint(uint64(tr.UserID), 10))))
	}

	// calculate time difference between clocked-out and clocked-in timestamps in minutes.
	timeDuration := timeDurationMinutes(*lastTimeRegister.ClockInTime, *clockOutTime)
	if result := db.Model(&lastTimeRegister).Updates(TimeRegister{ClockOutTime: clockOutTime, TimeDurationInMinutes: timeDuration}); result.Error != nil {
		logger.Log.Error(result.Error)
		return apiError.New(apiError.WithDetails(result.Error))
	}

	return nil
}

func getLastTimeRegisterRecord(db *gorm.DB, userID uint) (TimeRegister, error) {
	var tr TimeRegister

	result := db.Where("user_id = ?", userID).Last(&tr)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return tr, result.Error
	}

	return tr, nil
}

func timeDurationMinutes(timeOne time.Time, timeTwo time.Time) float64 {
	difference := timeTwo.Sub(timeOne)
	return difference.Minutes()
}
