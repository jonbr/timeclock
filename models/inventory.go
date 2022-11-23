package models

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	apiError "timeclock/error"

	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	Type string
}

type Measurement struct {
	Width         int `json:"width" gorm:"not null"`
	Height        int
	Thicknes      int
	Quantity      int
	GlassBoxID    uint    `gorm:"Index"`
	BluePrintName *string `gorm:"Index;size:191"`
}

type BluePrint struct {
	Name        *string     `gorm:"unique"`
	Measurement Measurement `gorm:"foreignKey:Name;references:BluePrintName"`
}

type GlassBox struct {
	BoxID        uint `gorm:"unique"`
	InternalName string
	Measurement  Measurement `gorm:"foreignKey:BoxID;references:GlassBoxID"`
}

type GlassBoxResponse struct {
	BoxID        uint
	InternalName string
	Measurements []Measurement
}

type BluePrintRequest struct {
	Name         string
	Measurements []Measurement
}

func CreateGlassBox(db *gorm.DB, boxID uint, internalName string, glassBoxData []byte) (*GlassBoxResponse, *apiError.ErrorResp) {
	idx, err := localnameValidation(internalName)
	fmt.Println(idx)
	if err != nil {
		return nil, apiError.New(apiError.WithDetails(err))
	}

	measurements := []Measurement{}
	var width, height int
	var index = 1
	scanner := bufio.NewScanner(strings.NewReader(string(glassBoxData)))
	scanner.Split(bufio.ScanWords)
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
			measurements = append(measurements, Measurement{
				Width:      width,
				Height:     height,
				Thicknes:   addMeasurement(scanner.Text()),
				GlassBoxID: boxID,
			})
			index = 0
		}
		index++
	}
	if err := scanner.Err(); err != nil {
		return nil, apiError.New(apiError.WithDetails(err))
	}

	glassBox := GlassBox{BoxID: boxID, InternalName: internalName}
	transactionError := db.Transaction(func(tx *gorm.DB) error {
		// create first all glasses in measurement table
		// followd by inserting into glass_boxes table.
		if err := tx.Create(&measurements).Error; err != nil {
			return err
		}
		if err := tx.Create(&glassBox).Error; err != nil {
			return err
		}

		return nil
	})
	if transactionError != nil {
		return nil, apiError.New(apiError.WithDetails(transactionError))
	}

	return &GlassBoxResponse{
		BoxID:        boxID,
		InternalName: internalName,
		Measurements: measurements,
	}, nil
}

func (bluePrintRequest *BluePrintRequest) CreateBluePrint(db *gorm.DB) *apiError.ErrorResp {
	//bulk create measurements records, fyrst add unique name to the measurement obj.
	for i := range bluePrintRequest.Measurements {
		bluePrintRequest.Measurements[i].BluePrintName = &bluePrintRequest.Name
	}

	transactionError := db.Transaction(func(tx *gorm.DB) error {
		if err := db.Create(&bluePrintRequest.Measurements).Error; err != nil {
			return err
		}

		// bulk create record in BluePrints for each individual measurement record.
		if err := db.Create(&BluePrint{Name: &bluePrintRequest.Name}).Error; err != nil {
			return err
		}

		return nil
	})
	if transactionError != nil {
		return apiError.New(apiError.WithDetails(transactionError))
	}

	return nil
}

func addMeasurement(measurement string) int {
	measurementInt, _ := strconv.Atoi(measurement)
	/*if err != nil {
		//return nil, apiError.New(apiError.WithDetails(err))
	}*/

	return measurementInt
}

// TODO: Add color validation
func localnameValidation(localname string) (int, error) {
	idx := strings.Index(localname, "-")
	customError := fmt.Errorf("Localname incorrect format '%s' must be on format 'color-number'", localname)
	if idx == -1 {
		return -1, customError
	}
	if _, err := strconv.Atoi(localname[idx+1:]); err != nil {
		return -1, customError
	}

	return idx, nil
}
