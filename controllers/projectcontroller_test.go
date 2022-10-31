package controllers

import (
   "net/http"
   "net/http/httptest"
   "testing"

   "github.com/gorilla/mux"
   "github.com/steinfletcher/apitest"
   "gorm.io/gorm"
   //jsonpath "github.com/steinfletcher/apitest-jsonpath"
)


func Test_GetProjects(t *testing.T) {
	var db *gorm.DB
   	r := mux.NewRouter()
   	r.Handle("/project/1", GetProjects(db))
   	ts := httptest.NewServer(r)
   	defer ts.Close()
   	apitest.New().
      Handler(r).
      Get("/project/1").
      Expect(t).
      Status(http.StatusOK).
      End()
}