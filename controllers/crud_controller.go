package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/jordan-borges-lark/todo_test/config"
	"github.com/jordan-borges-lark/todo_test/models"
	"github.com/jordan-borges-lark/todo_test/views"
)

// Not an entity
type ICrudController interface { // implements IController
	Controller
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)

	GetDatabase() *sqlx.DB
	GetModel() models.IModel
	SetModel(*models.Model)
}
type CrudController struct {
	Database *sqlx.DB
	Session  *sessions.CookieStore // Temporary hack to get it to compile, this needs to be a generic session store
	Model    *models.Model
	Config   *config.Config
}

func (cc *CrudController) GetDatabase() *sqlx.DB {
	return cc.Database
}
func (cc *CrudController) GetModel() models.Model   { return *cc.Model }
func (cc *CrudController) SetModel(m *models.Model) { cc.Model = m }

func (cc *CrudController) Create(w http.ResponseWriter, r *http.Request) {
	// @TODO add to r.Context to tell Update this is a create operation
	// @TODO: atm Create must be called without ID, or with nil/zero ID, otherwise it becomes like a regular Update attempt
	cc.Update(w, r)
}

func (cc *CrudController) Read(w http.ResponseWriter, r *http.Request) {
	m := cc.ReadBase(w, r)
	// JSON
	if r.Header.Get("Content-Type") == "application/json" {
		OutputJson(w, m)
		return
	}
	// @TODO: XML can by done just by replacing json.Marshal with xml.Marshal

	// HTML
	OutputHtml(w, m, "crud")
}
func (cc *CrudController) ReadBase(w http.ResponseWriter, r *http.Request) *models.Model {
	if m, err := cc.GetModelByIdFromRequest(r); err != nil { // Auth
		http.Error(w, err.Message, err.StatusCode)
		return nil
	} else {
		return m
	}
}

func (cc *CrudController) Update(w http.ResponseWriter, r *http.Request) {
	OutputJson(w, cc.UpdateBase(w, r))
}
func (cc *CrudController) UpdateBase(w http.ResponseWriter, r *http.Request) *sql.Result {
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	m := cc.GetModel()
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	// Load collections (needed for auth)
	if m2, err := m.Load(cc.Database); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	} else {
		m = m2.(models.Model)
	}

	// Auth
	if !m.CanUserWrite(123) {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return nil
	}

	// Update
	if result, err := models.Update(cc.Database, m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	} else {
		return &result
	}
}

func (cc *CrudController) Delete(w http.ResponseWriter, r *http.Request) {
	OutputJson(w, cc.DeleteBase(w, r))
}
func (cc *CrudController) DeleteBase(w http.ResponseWriter, r *http.Request) *sql.Result {
	// Auth
	m, err := cc.GetModelByIdFromRequest(r)
	if err != nil {
		http.Error(w, err.Message, err.StatusCode)
		return nil
	}

	// Enterprise data retention requirements?
	// Need to soft-delete by setting deleted_at = NOW()
	//if config.Database.ShouldSoftDelete {
		now := time.Now()
		m.DeletedAt = &now
		
		// Update
		if result, err := models.Update(cc.Database, m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil
		} else {
			return &result
		}
	//} // else @TODO
}

func OutputJson(w http.ResponseWriter, value interface{}) {
	if byt, err := json.Marshal(value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.Write(byt)
	}
}

func OutputHtml(w http.ResponseWriter, m *models.Model, viewName string) {
	//fields := models.GetDBFields(m)

	// This is called a marshal-unmarshal cycle; needed to convert struct => map[string]iface{}
	var fields map[string]interface{}
	byt, _ := json.Marshal(m)
	json.Unmarshal(byt, &fields)

	// Migrate collection fields to a separate view value "collections"
	collections := map[string][]map[string]interface{}{}
	for k, v := range fields {
		if array, isArray := v.([]map[string]interface{}); isArray {
			// Add a model with ID "0" to each collection
			// to create a blank row at the end (in the view)
			collections[k] = append(array, map[string]interface{}{"id":0})
			// Migrate value from "fields" to "collections"
			delete(fields, k)
		}
	}

	values := map[string]interface{}{
		"name": m.GetTableName(),
		"fields": fields,
		"collections": collections,
	}
	body, err := views.GetView(viewName, "", values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(body))
}

type HttpError struct {
	StatusCode int
	Message    string
}

func (he HttpError) Error() string { return he.Message }

func (cc *CrudController) GetModelByIdFromRequest(r *http.Request) (*models.Model, *HttpError) {
	// Get lazy-loaded model
	vars := mux.Vars(r)
	m := cc.Model
	m.Id, _ = strconv.Atoi(vars["id"])
	if m.Id == 0 {
		return m, &HttpError{StatusCode: http.StatusBadRequest, Message: "400 Bad Request\nInvalid value for argument 'id'"}
	}

	m, err := models.Get(cc.GetDatabase(), *m)
	if err != nil { // @TODO: Only echo DB errors when config.Environment == "dev"
		return m, &HttpError{StatusCode: http.StatusInternalServerError, Message: err.Error()}
	} else if m.Id == 0 {
		return m, &HttpError{StatusCode: http.StatusNotFound, Message: "404 Not Found"}
	}

	// Load model (required for auth to work)
	m2, err := m.Load(cc.Database)
	if err != nil { // @TODO: Only echo DB errors when config.Environment == "dev"
		return m, &HttpError{StatusCode: http.StatusInternalServerError, Message: err.Error()}
	}
	m = m2.(*models.Model)

	// Authorize user for reading this model
	//if !m.CanUserRead(cc.Session.GetLoggedInUserId()) {
	if !m.CanUserRead(123) {
		return nil, &HttpError{StatusCode: http.StatusForbidden, Message: "403 Forbidden"}
	}

	return m, nil
}
