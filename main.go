package main

import (
  "timeclock/controllers"
  "timeclock/database"

  "fmt"
  "log"
  "net/http"

  "github.com/gorilla/mux"
)


func main() {
  // Load Configurations from config.json using Viper
  LoadAppConfig()

  // Initialize Database
  database.Connect(AppConfig.ConnectionString)
  database.Migrate()

  // Initialize the router
  router := mux.NewRouter().StrictSlash(true)
  //router := sw.NewRouter()

  // Register Routes
  RegisterUserRoutes(router)

  // Start the server
  log.Println(fmt.Sprintf("Starting Server on port %s", AppConfig.Port))
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", AppConfig.Port), router))

  // create user
  //user := User{Name: "Helen Bjork", Administrator: false}
  //createUser(db, user)

  // delete user
  //deleteUser(db, 3)

  // create clock_in recore
  //clockIn(db, 1)

  // create clock_out record
  //clockOut(db, 1)

  // createTimeRegistrationRecord
  /*t1 := time.Date(2022, time.Month(10), 8, 7, 0, 0, 0, time.UTC)
  t2 := time.Date(2022, time.Month(10), 8, 15, 10, 30, 0, time.UTC)
  createTimeRegistrationRecord(db, 2, t1, t2)*/
  
  /*date, err := time.Parse("2006-01-02", "2022-10-09")
  if err != nil {
    log.Fatal(err)
  }
  modifyTimeRegistrationRecord(db, 1, date)*/

  // delete record
  //deleteTimeRegistrationRecord()

  //then := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
  //var registers = getTimeRegisterPeriod(db, 1, then, time.Now())
}

func RegisterUserRoutes(r *mux.Router) {
  //r := mux.NewRouter()
  //r.HandleFunc("/", controllers.Home)
  //r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))
  r.HandleFunc("/user", controllers.CreateUserHandler).Methods("POST")
  r.HandleFunc("/user", controllers.GetUserHandler).Methods("GET")
  r.HandleFunc("/user/{id}", controllers.GetUserHandler).Methods("GET")
  //r.HandleFunc("/user{id}", controllers.UpdateUserHandler).Methods("PUT")
  r.HandleFunc("/user{id}", controllers.DeleteUserHandler).Methods("DELETE")

  r.HandleFunc("/timeregistration/clockin/{id}", controllers.Timeregistrationclockin).Methods("POST")
  r.HandleFunc("/timeregistration/clockin/{id}", controllers.Timeregistrationclockout).Methods("POST")

  r.HandleFunc("/project/{id}", controllers.CreateProjectHandler).Methods("POST")
}

// CRUD interface
/*func createUser(db *gorm.DB, user User) {
  if err := db.Create(&user).Error; err != nil { // pass pointer of data to Create
    fmt.Println("Error creating user")
  }
}

func getUser(db *gorm.DB, userID uint) User {
  var user User
  if err := db.First(&user, userID).Error; err != nil {
    fmt.Println(err)
  }

  return user
}

func modifyUser(db *gorm.DB, userID uint) {}

func deleteUser(db *gorm.DB, userID uint) {
  db.Delete(&User{}, userID)
}

func clockIn(db *gorm.DB, userID uint) {
  // get user
  user := getUser(db, userID)
  if (User{} == user) {
    fmt.Println("No User found, not possible to clockIn!")
  }
  // create clockIn record
  var timeStamp = time.Now()
  db.Create(&TimeRegister{ClockIn: &timeStamp, UserID: userID})
}

func clockOut(db *gorm.DB, userID uint) {
  // get user
  user := getUser(db, userID)
  if (User{} == user) {
    fmt.Println("No User found, not possible to clockOut!")
  }
  // retrieve last record where clock_out == NULL
  var lastRecord TimeRegister
  db.Raw("SELECT * FROM time_registers WHERE user_id = ? AND clock_out IS NULL ORDER BY clock_in DESC LIMIT 1", userID).Scan(&lastRecord)
  if (TimeRegister{} == lastRecord) {
    fmt.Println("No record found that matches the query!")
  }

  // update lastRecord with clock_out timestamp
  db.Model(&TimeRegister{}).Where("id = ?", lastRecord.ID).Update("clock_out", time.Now())
}

func createTimeRegistrationRecord(db *gorm.DB, userID uint, clockIn time.Time, clockOut time.Time) {
  user := getUser(db, userID)
  if user.Administrator == false {
    fmt.Println("User not with sufficient privledges to perform this action!")
  }
  // Create record with either clockIn record or both clockIn and ClockOut record
  db.Create(&TimeRegister{ClockIn: &clockIn, ClockOut: &clockOut, UserID: userID})
}

func modifyTimeRegistrationRecord(db *gorm.DB, userID uint, date time.Time) {
  // 1.) first fetch user and check it's access rights
  user := getUser(db, userID)
  if (User{} == user) {
    fmt.Println("No User found, not possible to modifyTimeRegistrationRecord!")
  }
  // update lastRecord with clock_out timestamp
  db.Model(&TimeRegister{}).Where("id = ?", user.ID).Update("clock_out", time.Now())
}

func deleteTimeRegistrationRecord(db *gorm.DB) {
  db.Delete(&User{}, userID)
}

func getTimeRegisterPeriod(db *gorm.DB, userID uint, period_start time.Time, period_end time.Time) ([]TimeRegister) {
  var register[] TimeRegister
  db.Where("user_id = ? AND clock_in > ? AND clock_out < ?", userID, period_start, period_end).Find(&register)

  return register
}*/


