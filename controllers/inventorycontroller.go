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

		glassBoxResponse, errRespo := models.CreateGlassBox(db, boxID, vars["internalname"], glassBoxRawData)
		if errRespo != nil {
			logger.Log.Error(errRespo)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errRespo)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(glassBoxResponse)
	}
}

func InventoryCreateBluePrint(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var bp models.BluePrint
		if err := json.NewDecoder(r.Body).Decode(&bp); err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(error.New(error.WithDetails(err)))
		}

		if err := bp.CreateBluePrint(db); err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bp)
	}
}

func InventoryUpdateBluePrint(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("enpoint not yet implemented!")
	}
}

func InventoryDeleteBluePrint(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("enpoint not yet implemented!")
	}
}

func InventoryGetBluePrint(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		models.CompareBluePrintWithGlassBox(db)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("enpoint not yet implemented!")
	}
}
