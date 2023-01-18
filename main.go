package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jasonlvhit/gocron"
	"github.com/jmoiron/sqlx"
	"github.com/jordan-borges-lark/todo_test/config"
	"github.com/jordan-borges-lark/todo_test/controllers"
	"github.com/jordan-borges-lark/todo_test/helpers"
	"github.com/jordan-borges-lark/todo_test/models"
)

func main() {
	/* === Config === */

	// Load config

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" { // There's no null coalescing operator ??/?:/|| or ternaries in golang, so this pattern is idiomatic for go
		appEnv = "dev"
	}
	cfg := &config.Config{}
	if err := cfg.Load(appEnv); err != nil {
		panic(err)
	}

	// Environment-specific startup logic

	switch appEnv {
	case "dev":
		fmt.Println("Debug statements")
	case "prod":
	}

	/* === Golang-specific memory ballast === */

	// Create a large heap allocation; 30 means 10 GiB
	ballast := make([]byte, 10<<cfg.BallastSize)
	go func(b []byte) {
		select {}
	}(ballast)

	/* === SQL === */

	// Create DB connection
	db := GetDB(
		cfg.Database.DriverName,
		cfg.Database.DataSourceName,
	)
	defer db.Close()

	// Keep-alive timer to keep DB connected
	var keepalive func(db *sqlx.DB)
	keepalive = func(db *sqlx.DB) {
		if err := db.Ping(); err != nil {
			panic(err)
		}
		time.Sleep(time.Duration(
			cfg.Database.KeepAliveSeconds * int64(time.Second)),
		)
		go keepalive(db) // Without a goroutine here you'll have a stack overflow
	}
	go keepalive(db)

	/* === Cron jobs === */
	{
		cron := cron{}
		gocron.Every(1).Day().Do(cron.sampleCronJob)
		cron.startJobs()
	}

	/* === Web server === */

	// Set up muxxer routes

	r := mux.NewRouter()

	// Home page default is to create a new list
	// @TODO: Make home page be the MAX(updated_at) list?
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/list/0", http.StatusSeeOther)
	}).Methods("GET")

	// Set up CRUD routes

	var overrides = map[string]string{
		models.ItemList{}.GetTableName():     "list",
		models.ItemListItem{}.GetTableName(): "item",
	}

	for m, cc := range map[interface{}]interface{}{
		models.User{}:         controllers.User{}, // With reflection you wouldn't even need to wire this up
		models.ItemList{}:     controllers.ItemList{},
		models.ItemListItem{}: controllers.ItemListItem{},
	} {
		// Set up base crud controller
		m := m.(models.Model)
		cc := cc.(controllers.CrudController)
		cc.Database = db
		cc.Config = cfg
		cc.Session = &sessions.Session{}
		cc.SetModel(&m)

		entityName := m.GetTableName()
		if newName, exists := overrides[entityName]; exists {
			entityName = newName
		}

		// Add routes for each CRUD operation
		for path, handlers := range map[string]map[string]func(w http.ResponseWriter, r *http.Request) {
			"/"+helpers.Pluralize(entityName): {
				//"GET": cc.Index, // Not req'd for this project
			}, "/"+entityName: {
				"POST": cc.Create,
			}, "/"+entityName+"/{id}": {
				"GET": cc.Read,
				"PATCH": cc.Update,
				"PUT": cc.Update,
				"DELETE": cc.Delete,
			},
		}{
			for method, handler := range handlers {
				r.HandleFunc(path, handler).Methods(method)
			}
		}
	}

	// @TODO: Add middleware
	//mux.Use(authMiddleware)
	//mux.Use(jsonMiddleware)

	fmt.Println("Running on port", os.Getenv("APP_PORT"))
	log.Println("Running on port", os.Getenv("APP_PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("APP_PORT"), r))
}
