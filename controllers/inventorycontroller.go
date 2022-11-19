package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"timeclock/error"
	"timeclock/logger"
	"timeclock/models"
	"timeclock/utils"

	"github.com/gookit/goutil/dump"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func InventoryGlass(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/(text|json); charset=UTF-8")

		boxID, err := utils.CastStringToUint(mux.Vars(r))
		if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		glassBoxRawData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(error.New(error.WithDetails(err)))
		}
		dump.P(string(glassBoxRawData))

		var glassBoxes []models.GlassBox
		var errRespo *error.ErrorResp
		if glassBoxes, errRespo = models.InventoryGlassCreate(db, boxID[0], glassBoxRawData); errRespo != nil {
			logger.Log.Error(errRespo)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errRespo)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(glassBoxes)
	}
}

func InventoryBluePrint(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		//w.Header().Set("Access-Control-Allow-Origin", "*")

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Endpoint not yet implemented!")
	}
}
