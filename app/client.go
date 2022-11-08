package app

import (
	"log"

	"timeclock/controllers"
	"timeclock/middlewares"
	"timeclock/models"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Connect(connectionString string) {
	var err error
	a.DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Database...")
}

func (a *App) Migrate() {
	a.DB.AutoMigrate(&models.User{}, &models.TimeRegister{}, &models.Project{})
	log.Println("Database Migration Completed...")
}

func (a *App) InitializeRoutes() {
	a.Router.Handle("/token", controllers.GenerateToken(a.DB)).Methods("POST")

	a.Router.Handle("/user", controllers.GetUsers(a.DB)).Methods("GET") // check token if user has sufficient privledges to execute this operation
	a.Router.Handle("/user/{id}", controllers.GetUser(a.DB)).Methods("GET")
	a.Router.Handle("/user/register", controllers.CreateUser(a.DB)).Methods("POST")
	a.Router.Handle("/user/{id}", controllers.UpdateUser(a.DB)).Methods("PUT")
	a.Router.Handle("/user/{id}", controllers.DeleteUser(a.DB)).Methods("DELETE")

	a.Router.Handle("/project", controllers.GetProjects(a.DB)).Methods("GET")
	a.Router.Handle("/project/{id}", controllers.GetProject(a.DB)).Methods("GET")
	
	a.Router.Handle("/project/", controllers.CreateProject(a.DB)).Methods("POST")
	a.Router.Handle("/project/{id}", controllers.UpdateProject(a.DB)).Methods("PUT")
	a.Router.Handle("/project/{id}", controllers.DeleteProject(a.DB)).Methods("DELETE")

	a.Router.Handle("/timeregistration/clockin", controllers.TimeRegistrationClockIn(a.DB)).Methods("POST")
	a.Router.Handle("/timeregistration/clockout", controllers.TimeRegistrationClockOut(a.DB)).Methods("POST")

	a.Router.Use(middlewares.Auth)
}
