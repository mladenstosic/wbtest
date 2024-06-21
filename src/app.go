package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Router *mux.Router
}

type User struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	DateOfBirth string `json:"date_of_birth"`
}

type output struct {
	Response_code int    `json:"response_code"`
	Message       string `json:"message"`
	Response      *User  `json:"response,omitempty"`
}

// Init app
func (app *App) Init(db string) error {

	var err error

	// Open db
	app.DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to %s database (err:%s)", db, err.Error())
	}

	// Migrate the schema
	app.DB.AutoMigrate(&User{})

	// Define router
	app.Router = mux.NewRouter()

	// Register handler to get the user data
	app.Router.Handle("/{id}", app).Methods("GET")

	// Register handler to save user data
	app.Router.Handle("/save", app).Methods("POST")

	return nil
}

// Run app
func (app *App) Run(port string) error {
	return http.ListenAndServe(port, app.Router)
}

// HTTP handler
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set the header to json
	w.Header().Set("Content-Type", "application/json")

	// Define user and output
	var (
		user User
		out  output
	)

	// Defer the output processing
	defer func() {
		// Marshal response
		output, errOut := json.Marshal(out)
		if errOut != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("cannot marshal output\n"))
			log.Println("cannot marshal output", out)
		}

		// Log the response
		log.Println(out.Message)
		w.Write([]byte(output))
	}()

	// Diferentiate between endpoints
	switch r.Method {
	// Get the user
	case "GET":

		// Get the id from url
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			out.Response_code = 400
			out.Message = "cannot find user id"
			return
		}

		// Search for user from id
		app.DB.Model(&User{}).Select("*").Where("Id = ?", id).Scan(&user)
		if user.Id == "" {
			w.WriteHeader(http.StatusNotFound)
			out.Response_code = 404
			out.Message = "user id:" + id + " not found"
			return
		}

		// Return user data
		out.Response_code = 200
		out.Message = "user id:" + id + " found"
		out.Response = &user
		return

	// Save the user
	case "POST":

		//TODO: test, email, date of birth

		// Get the data and unmarshal it
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			out.Response_code = 400
			out.Message = "bad data format"
			return
		}

		// Check for user details
		// ID field is required, here we could check if all required user data exists or if its formated properly
		if user.Id == "" {
			w.WriteHeader(http.StatusBadRequest)
			out.Response_code = 400
			out.Message = "bad data format, please add id(required),name,email,date_of_birth"
			return
		}

		// Save the user
		err = app.DB.Create(user).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			out.Response_code = 500
			out.Message = "cannot save user with id:" + user.Id + " error:" + err.Error()
			return
		}

		// User saved
		out.Response_code = 200
		out.Message = "user with id:" + user.Id + " saved"
		out.Response = &user
	}
}
