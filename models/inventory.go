package models

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	apiError "timeclock/error"

	"github.com/gookit/goutil/dump"
	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	Type string
}

type BluePrint struct {
	Name         *string       `gorm:"primaryKey;size:191"`
	Measurements []Measurement `gorm:"foreignKey:BluePrintName"`
}

type GlassBox struct {
	BoxID        uint `gorm:"primaryKey;autoIncrement:false"`
	InternalName string
	Measurements []Measurement `gorm:"foreignKey:GlassBoxID"`
}

type Measurement struct {
	Width         int `json:"width" gorm:"not null"`
	Height        int
	Thickness     int
	Quantity      int
	GlassBoxID    *uint
	BluePrintName *string `gorm:"size:191"`
}

type GlassBoxResponse struct {
	BoxID        *uint
	InternalName string
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
		m.GlassBoxID = &boxID

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
		if err := tx.Create(&glassBox).Error; err != nil {
			return err
		}
		if err := tx.Create(&measurements).Error; err != nil {
			return err
		}

		return nil
	})
	if transactionError != nil {
		return nil, apiError.New(apiError.WithDetails(transactionError))
	}

	return &GlassBoxResponse{
		BoxID:        &boxID,
		InternalName: internalName,
		Measurements: measurements,
	}, nil
}

func (bluePrint *BluePrint) CreateBluePrint(db *gorm.DB) *apiError.ErrorResp {
	// Add unique name to the measurement obj.
	for i := range bluePrint.Measurements {
		bluePrint.Measurements[i].BluePrintName = bluePrint.Name
	}

	transactionError := db.Transaction(func(tx *gorm.DB) error {
		if err := db.Omit("Measurements").Create(&bluePrint).Error; err != nil {
			return err
		}

		// bulk create record in BluePrints for each individual measurement record.
		if err := db.Create(&bluePrint.Measurements).Error; err != nil {
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

func CompareBluePrintWithGlassBox(db *gorm.DB) {
	fmt.Println("---CompareBluePrintWithGlassBox---")

	//var bluePrint BluePrint
	var bluePrints []BluePrint

	db.Where("name = ?", "VAL-1B").Preload("Measurements").Find(&bluePrints)

	dump.P(bluePrints)
}
