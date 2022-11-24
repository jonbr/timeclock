package models

import (
	"bufio"
	"bytes"
	"errors"
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
	Thickness     int
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
	count := make(map[Measurement]int)
	var measurements []Measurement
	scanner := bufio.NewScanner(bytes.NewReader(glassBoxData))
	for scanner.Scan() {
		m, err := ParseMeasurement(scanner.Text())
		if err != nil {
			continue
		}
		m.GlassBoxID = boxID

		if n := count[m]; n > 0 {
			count[m] = n + 1
		} else {
			count[m] = 1
			measurements = append(measurements, m)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, apiError.New(apiError.WithDetails(err))
	}

	for i := range measurements {
		(&measurements[i]).Quantity = count[measurements[i]]
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

func ParseMeasurement(line string) (m Measurement, err error) {
	for _, s := range strings.Split(line, " X ") {
		unit := strings.Split(strings.TrimSpace(s), " ")
		if len(unit) != 2 {
			continue
		}

		v, err := strconv.Atoi(unit[0])
		if err != nil {
			return m, err
		}
		switch unit[1] {
		case "W":
			m.Width = v
		case "H":
			m.Height = v
		case "T":
			m.Thickness = v
		}
	}
	if m == (Measurement{}) {
		return m, errors.New("empty")
	}
	return m, nil
}
