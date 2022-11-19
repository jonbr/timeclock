package models

import (
	"bufio"
	"strconv"
	"strings"
	"time"
	apiError "timeclock/error"

	"gorm.io/gorm"
)

type Inventory struct {
	ID        uint `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Type      string
}

type measurements struct {
	Width    int `json:"width" gorm:"not null"`
	Height   int
	Thicknes int
}

type BluePrint struct {
	Inventory           Inventory    `gorm:"embedded"`
	Measurements        measurements `gorm:"embedded"`
	Quantity            int
	BluePrintIdentifier string
}

type GlassBox struct {
	BoxID        int
	Inventory    Inventory    `gorm:"foreignKey:BoxID"`
	Measurements measurements `gorm:"embedded"`
}

func InventoryGlassCreate(db *gorm.DB, glassBoxID uint, glassBoxData []byte) ([]GlassBox, *apiError.ErrorResp) {
	scanner := bufio.NewScanner(strings.NewReader(string(glassBoxData)))
	scanner.Split(bufio.ScanWords)

	//charsToSkip := []string{"S", "W", "X", "H", "T"}

	glassBoxes := []GlassBox{}
	var width, height int

	var index = 1
	for scanner.Scan() {
		if scanner.Text() == "S" || scanner.Text() == "W" || scanner.Text() == "H" || scanner.Text() == "T" || scanner.Text() == "X" {
			continue
		}

		switch index {
		case 1:
			width = addMeasurement(scanner.Text())
		case 2:
			height = addMeasurement(scanner.Text())
		case 3:
			glassBoxes = append(glassBoxes, GlassBox{
				BoxID: int(glassBoxID),
				Measurements: measurements{
					Width:    width,
					Height:   height,
					Thicknes: addMeasurement(scanner.Text()),
				},
			})
			index = 0
		}
		index++
	}
	if err := scanner.Err(); err != nil {
		return nil, apiError.New(apiError.WithDetails(err))
	}

	// first create Inventory record
	if err := db.Create(&Inventory{ID: glassBoxID, Type: "glass"}).Error; err != nil {
		return nil, apiError.New(apiError.WithDetails(err))
	}

	// then insert all glasses into GlassBox table
	if err := db.Create(&glassBoxes).Error; err != nil {
		return nil, apiError.New(apiError.WithDetails(err))
	}

	return glassBoxes, nil
}

func addMeasurement(measurement string) int {
	measurementInt, _ := strconv.Atoi(measurement)
	/*if err != nil {
		//return nil, apiError.New(apiError.WithDetails(err))
	}*/

	return measurementInt
}
