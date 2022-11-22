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

		vars := mux.Vars(r)
		dump.P(vars)

		boxID, err := utils.CastParamToUint(vars["boxid"])
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
		//dump.P(string(glassBoxRawData))

		var glassBoxes []models.GlassBox
		var errRespo *error.ErrorResp
		if glassBoxes, errRespo = models.CreateGlassBox(db, boxID, vars["localname"], glassBoxRawData); errRespo != nil {
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

		var bp models.BluePrint
		if err := json.NewDecoder(r.Body).Decode(&bp); err != nil {
			logger.Log.Error(err)
		}

		models.CreateBluePrint(db)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bp)
	}
}
