package main

import (
  "timeclock/app"

  "fmt"
  "log"
  "net/http"

  "github.com/gorilla/mux"
  //"github.com/gookit/goutil/dump"
)


func main() {
  a := &app.App{}

  // Load Configurations from config.json using Viper
  LoadAppConfig()

  // Initialize Database
  a.Connect(AppConfig.ConnectionString)
  //a.Migrate()

  // Initialize the router
  a.Router = mux.NewRouter().StrictSlash(true)
  //router := sw.NewRouter()

  // Initialize Routes
  a.InitializeRoutes()

  // Start the server
  log.Println(fmt.Sprintf("Starting Server on port %s", AppConfig.Port))
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", AppConfig.Port), a.Router))
}
