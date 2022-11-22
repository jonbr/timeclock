package models

import (
	"bufio"
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

type measurement struct {
	Width    int `json:"width" gorm:"not null"`
	Height   int
	Thicknes int
	Quantity int
}

type BluePrintDefinition struct {
	Name        string `gorm:"primaryKey;size:191"`
	InventoryID uint
	Inventory   Inventory
}

type BluePrint struct {
	Name                string
	BluePrintDefinition BluePrintDefinition `gorm:"foreignKey:Name"`
	Measurement         measurement         `gorm:"embedded"`
	Quantity            int
	BluePrintIdentifier string
}

type GlassDefinition struct {
	ID          uint `gorm:"primaryKey;autoIncrement:false"`
	InventoryID uint
	Inventory   Inventory
	ColorScheme string `gorm:"not null"`
	LocalName   string `gorm:"not null"`
}

type GlassBox struct {
	BoxID           uint
	GlassDefinition GlassDefinition `gorm:"foreignKey:BoxID"`
	Measurement     measurement     `gorm:"embedded"`
}

func CreateGlassBox(db *gorm.DB, boxID uint, localname string, glassBoxData []byte) ([]GlassBox, *apiError.ErrorResp) {

	idx, err := localnameValidation(localname)
	if err != nil {
		return nil, apiError.New(apiError.WithDetails(err))
	}

	//charsToSkip := []string{"S", "W", "X", "H", "T"}
	// first create Inventory record
	var glassDefinition = GlassDefinition{ID: boxID, InventoryID: 1, ColorScheme: localname[:idx], LocalName: localname}
	inventoryResult := db.Create(&glassDefinition)
	if inventoryResult.Error != nil {
		return nil, apiError.New(apiError.WithDetails(inventoryResult.Error))
	}

	glassBoxes := []GlassBox{}
	var width, height int
	var index = 1
	scanner := bufio.NewScanner(strings.NewReader(string(glassBoxData)))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		if scanner.Text() == "S" || scanner.Text() == "W" || scanner.Text() == "H" || scanner.Text() == "T" || scanner.Text() == "X" {
			continue
		}

		dump.P(scanner.Text())

		switch index {
		case 1:
			width = addMeasurement(scanner.Text())
		case 2:
			height = addMeasurement(scanner.Text())
		case 3:
			glassBoxes = append(glassBoxes, GlassBox{
				BoxID: glassDefinition.ID,
				Measurement: measurement{
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

	// then insert all glasses into GlassBox table
	if err := db.Create(&glassBoxes).Error; err != nil {
		return nil, apiError.New(apiError.WithDetails(err))
	}

	return glassBoxes, nil
}

func CreateBluePrint(db *gorm.DB) {

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
