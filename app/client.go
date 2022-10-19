package app

import (
	"timeclock/controllers"
	"timeclock/models"

	"log"
	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/gorilla/mux"
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
	a.Router.Handle("/project/", controllers.GetProjects(a.DB)).Methods("GET")
	a.Router.Handle("/project/{id}", controllers.GetProject(a.DB)).Methods("GET")
  	a.Router.Handle("/project/{userId}", controllers.CreateProject(a.DB)).Methods("POST")
	//a.Router.Handle("/project/{userId}", controllers.UpdateProject(a.DB)).Methods("PUT")
	//a.Router.Handle("/project/{userId}", controllers.DeleteProject(a.DB)).Methods("DELETE")

	/*a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")*/
}
