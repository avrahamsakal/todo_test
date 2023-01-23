package main

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jasonlvhit/gocron"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
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

	// Load .env file into environment
	godotenv.Load()

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

	// Run DB migrations
	datastores.RunSqlMigrations(db)

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
	for _, m := range []models.IModel[any]{ // @TODO: With reflection you wouldn't even need to wire this up
		models.User{},     // i.e. profile page
		models.ItemList{}, // i.e. todo list
		models.ItemListItem{},
		models.Metadata{}, // i.e. admin panel
	} {
		crudHandleFunc(r, controllers.CrudController[models.IModel[any]]{
			controllers.Controller{}, db, ss, m, cfg,
		})
	}

	// Set up static file handler. MUST be the last (i.e. catch-all) route to add!!!
	r.PathPrefix("/").HandlerFunc(catchAllHandler)

	// @TODO: Add middleware
	//mux.Use(authMiddleware) // Auths user for read/write action on model requested

	// Start server and log on failure

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "1818"
	}
	fmt.Println("Running on port", port)
	log.Println("Running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// e.g. /list/12345, /item/12345 /profile, /admin
var modelFriendlyNames = map[string]string{
	models.User{}.GetTableName():         "profile",
	models.ItemList{}.GetTableName():     "list",
	models.ItemListItem{}.GetTableName(): "item",
	models.Metadata{}.GetTableName():     "admin",
}

func crudHandleFunc[CC controllers.ICrudController[models.IModel[any]]](
	r *mux.Router,
	cc CC,
) {
	m := cc.GetModel()

	// Add routes for each CRUD operation
	for _, entityName := range getEntityNames(m) {
		for path, handlers := range map[string]map[string]func(w http.ResponseWriter, r *http.Request){
			"/" + helpers.Pluralize(entityName): {
				"GET": cc.Index, // Not req'd for this project
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
}

// Build list of aliases for model name
func getEntityNames[M models.IModel[any]](m M) []string {
	entityNames := []string{m.GetTableName()}
	if newName, exists := modelFriendlyNames[entityNames[0]]; exists {
		entityNames = append(entityNames, newName)
	}
	typeName := ""
	if T := reflect.TypeOf(m); T.Kind() == reflect.Ptr {
		typeName = "*" + T.Elem().Name()
	} else {
		typeName = T.Name()
	}
	entityNames = append(entityNames, strings.Replace(typeName, "models.", "", 1))
	return entityNames
}

// catchAllHandler returns a static file if file exists; else 404
func catchAllHandler(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join("./", path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", mime.TypeByExtension(path))
	http.ServeFile(w, r, path)
}
