package controllers

import (
   "net/http"
   "net/http/httptest"
   "testing"

   "github.com/gorilla/mux"
   "gorm.io/driver/mysql"
   "gorm.io/gorm"
   "github.com/DATA-DOG/go-sqlmock"
   "github.com/steinfletcher/apitest"
   jsonpath "github.com/steinfletcher/apitest-jsonpath"
)


func Test_GetProject(t *testing.T) {
   db, mock, err := sqlmock.New()
   if err != nil {
      t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
   }
   defer db.Close()

   gormDB, err := gorm.Open(mysql.New(mysql.Config{
      Conn: db,
      SkipInitializeWithVersion: true,
   }), &gorm.Config{})
   if err != nil {
      panic(err) // Error here
   }

   var (
       id            = uint(1)
       name          = "Valshlid"
       description   = "Rand byggður kassi við Hlidarenda i RVK"
   )
   mock.ExpectQuery("SELECT(.*)").
      WithArgs(id).
      WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).
      AddRow(uint(id), name, description))

   r := mux.NewRouter()
   r.HandleFunc("/project/{id}", GetProject(gormDB)).Methods("GET")
   ts := httptest.NewServer(r)
   defer ts.Close()
   apitest.New().
      Debug().
      Handler(r).
      Get("/project/1").
      Expect(t).
      Body(`{"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"ID":1,"name":"Valshlid","description":"Rand byggður kassi við Hlidarenda i RVK"}`).
      Status(http.StatusOK).
      //Assert(jsonpath.Equal(`$.message`, "unknown key message")).
      End()
}

func Test_GetProjects(t *testing.T) {
   db, mock, err := sqlmock.New()
   if err != nil {
      t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
   }
   defer db.Close()

   gormDB, err := gorm.Open(mysql.New(mysql.Config{
      Conn: db,
      SkipInitializeWithVersion: true,
   }), &gorm.Config{})
   if err != nil {
      panic(err) // Error here
   }

   rows := sqlmock.NewRows([]string{"id", "name", "description"}).
      AddRow(uint(1), "name1", "description1").
      AddRow(uint(2), "name2", "description2").
      AddRow(uint(3), "name3", "description3")

   mock.ExpectQuery("SELECT").
      WillReturnRows(rows)

   r := mux.NewRouter()
   r.HandleFunc("/project", GetProjects(gormDB)).Methods("GET")
   ts := httptest.NewServer(r)
   defer ts.Close()
   apitest.New().
      Debug().
      Handler(r).
      Get("/project").
      Expect(t).
      Status(http.StatusOK).
      Assert(jsonpath.Len("$", 3)).
      End()
}