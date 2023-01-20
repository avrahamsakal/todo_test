package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jasonlvhit/gocron"
	"github.com/jmoiron/sqlx"
	"github.com/jordan-borges-lark/todo_test/config"
	"github.com/jordan-borges-lark/todo_test/controllers"
	"github.com/jordan-borges-lark/todo_test/datastores"
	"github.com/jordan-borges-lark/todo_test/helpers"
	"github.com/jordan-borges-lark/todo_test/models"
)

// @TODO: Abstract out these setup sections into separate functions
//
//	in this file called setupConfig(), setupSqlDB()
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
	ballast := make([]byte, 10<<cfg.Ballast_Size)
	go func(b []byte) {
		select {} // Blocks forever to keep GC from collecting ballast
	}(ballast)

	/* === SQL === */

	// Create DB connection
	db, err := datastores.GetSqlDB(
		cfg.Database.Driver_Name,
		cfg.Database.Data_Source_Name,
	)
	if err != nil || db == nil {
		panic(fmt.Sprint(
			"Unable to connect to DB '",
			cfg.Database.Driver_Name,
			"':", err,
		))
	}
	defer db.Close()

	// Keep-alive timer to keep DB connected
	var keepalive func(db *sqlx.DB)
	keepalive = func(db *sqlx.DB) {
		if err := db.Ping(); err != nil {
			panic(err)
		}
		time.Sleep(time.Duration(
			cfg.Database.Keep_Alive_Seconds * int64(time.Second)),
		)
		go keepalive(db) // Without a goroutine here you'll have a stack overflow
	}
	go keepalive(db)

	/* === Load session store === */

	ss := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	/* === Cron jobs === */

	{ // (Limit scope of this block, it's not reused)
		cron := cron{
			SessionExpirationDays: cfg.Session.Expiration_Days,
			Database:              db,
		}
		gocron.Every(1).Day().Do(cron.pruneExpiredSessions)
		cron.startJobs()
	}

	/* === Web server === */

	// Set up muxxer routes

	r := mux.NewRouter()

	// Root/home page default is to create a new list
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If you have no user, a session and session user is
		// created for you by /list/0
		http.Redirect(w, r, "/list/0", http.StatusSeeOther)
		// @TODO: Make home page be the MAX(updated_at) list; if NULL redirs to /list/0
	}).Methods("GET")

	// Set up CRUD routes
	for _, cc := range []controllers.ICrudController[models.IModel[models.Model]]{ // @TODO: With reflection you wouldn't even need to wire this up
		controllers.User[models.User]{},
		controllers.ItemList[models.ItemList]{},
		controllers.ItemListItem[models.ItemListItem]{},
	} {
		//cc.CrudController = controllers.CrudController[models.User]{db, ss, m, cfg}
		crudHandleFunc(r, cc, db, cfg, ss)
	}	

	// @TODO: Add middleware
	//mux.Use(authMiddleware) // Auths user for read/write action on model requested
	//mux.Use(jsonMiddleware)

	fmt.Println("Running on port", os.Getenv("APP_PORT"))
	log.Println("Running on port", os.Getenv("APP_PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("APP_PORT"), r))
}

var crudRouteOverrides = map[string]string{ // @TODO: Pull this from the config?
	models.ItemList{}.GetTableName():     "list",
	models.ItemListItem{}.GetTableName(): "item",
}

func crudHandleFunc[CC controllers.ICrudController[models.IModel[any]]](
	r *mux.Router,
	cc CC,
	db *sqlx.DB,
	cfg *config.Config,
	sessionStore *sessions.CookieStore,
) {
	cc := cc.Get() // todo implement this on both CC and ICC

	// Set up base crud controller
	cc.Database = db
	cc.Config = cfg
	cc.Session = sessionStore
	//cc.SetModel(m) // Should be magically handled with generics

	m := cc.GetModel()
	entityName := m.GetTableName()
	if newName, exists := crudRouteOverrides[entityName]; exists {
		entityName = newName
	}

	// Add routes for each CRUD operation
	for path, handlers := range map[string]map[string]func(w http.ResponseWriter, r *http.Request){
		"/" + helpers.Pluralize(entityName): {
			//"GET": cc.Index, // Not req'd for this project
		}, "/" + entityName: {
			"POST": cc.Create,
		}, "/" + entityName + "/{id}": {
			"GET":    cc.Read,
			"PATCH":  cc.Update,
			"PUT":    cc.Update,
			"DELETE": cc.Delete,
		},
	} {
		for method, handler := range handlers {
			r.HandleFunc(path, handler).Methods(method)
		}
	}
}
