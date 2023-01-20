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
	views "github.com/jordan-borges-lark/todo_test/views/crud"
)

// Not an entity
type ICrudController[M models.IModel[any]] interface { // implements IController
	Controller
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)

	GetDatabase() *sqlx.DB
	GetModel() M
	SetModel(M)
}

type CrudController[M models.IModel[any]] struct {
	Database *sqlx.DB
	Session  *sessions.CookieStore // Temporary hack to get it to compile, this needs to be a generic session store
	Model    M
	Config   *config.Config
}

func (cc *CrudController[M]) GetDatabase() *sqlx.DB {
	return cc.Database
}
func (cc *CrudController[M]) GetModel() M  { return cc.Model }
func (cc *CrudController[M]) SetModel(m M) { cc.Model = m }

func (cc *CrudController[M]) Create(w http.ResponseWriter, r *http.Request) {
	// @TODO add to r.Context to tell Update this is a create operation
	// @TODO: atm Create must be called without ID, or with nil/zero ID, otherwise it becomes like a regular Update attempt
	cc.Update(w, r)
}

func (cc *CrudController[M]) Read(w http.ResponseWriter, r *http.Request) {
	m := cc.ReadBase(w, r)
	// JSON
	if r.Header.Get("Content-Type") == "application/json" {
		OutputJson(w, m)
		return
	}
	// @TODO: XML can by done just by replacing json.Marshal with xml.Marshal

	// HTML
	OutputHtml(w, m)
}
func (cc *CrudController[M]) ReadBase(w http.ResponseWriter, r *http.Request) M {
	m, err := cc.GetModelByIdFromRequest(r)
	if err != nil { // Auth
		http.Error(w, err.Message, err.StatusCode)
	}

	return m
}

func (cc *CrudController[M]) Update(w http.ResponseWriter, r *http.Request) {
	OutputJson(w, cc.UpdateBase(w, r))
}

/*func toModel[M models.Model](m M) *models.Model {
	m2 := models.Model(m)
	return &m2
}*/

func (cc *CrudController[M]) UpdateBase(w http.ResponseWriter, r *http.Request) *sql.Result {
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	m := cc.GetModel()
	
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	// Load collections (needed for auth)
	m2, err := m.Load(cc.Database, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	} else {
		m = m2.(M)
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

func (cc *CrudController[M]) Delete(w http.ResponseWriter, r *http.Request) {
	OutputJson(w, cc.DeleteBase(w, r))
}
func (cc *CrudController[M]) DeleteBase(w http.ResponseWriter, r *http.Request) *sql.Result {
	// Auth
	m, err := cc.GetModelByIdFromRequest(r)
	if err != nil {
		http.Error(w, err.Message, err.StatusCode)
		return nil
	}
	//m := m2.Get()

	// Enterprise data retention requirements?
	// Need to soft-delete by setting deleted_at = NOW()
	//if config.Database.ShouldSoftDelete {
	now := time.Now()
	m.SetDeletedAt(&now)

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

func OutputHtml[M models.IModel[any]](w http.ResponseWriter, m M) {
	body, err := views.GetCrudView(m)
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

func (cc *CrudController[M]) GetModelByIdFromRequest(r *http.Request) (M, *HttpError) {
	// Get lazy-loaded model
	vars := mux.Vars(r)
	m := cc.GetModel()
	if id, _ := strconv.Atoi(vars["id"]); id == 0 {
		return m, &HttpError{StatusCode: http.StatusBadRequest, Message: "400 Bad Request\nInvalid value for argument 'id'"}
	} else {
		m.SetId(id)
	}

	if m2, err := models.Read(cc.GetDatabase(), m); err != nil { // @TODO: Only echo DB errors when config.Environment == "dev"
		return m2, &HttpError{StatusCode: http.StatusInternalServerError, Message: err.Error()}
	} else {
		m = m2
	}

	if m.GetId() == 0 {
		return m, &HttpError{StatusCode: http.StatusNotFound, Message: "404 Not Found"}
	}

	// Load model (required for auth to work)
	if m2, err := m.Load(cc.Database, false); err != nil { // @TODO: Only echo DB errors when config.Environment == "dev"
		return m2.(M), &HttpError{StatusCode: http.StatusInternalServerError, Message: err.Error()}
	} else {
		m = m2.(M)
	}

	// Authorize user for reading this model
	//if !m.CanUserRead(cc.Session.GetLoggedInUserId()) {
	if !m.CanUserRead(123) {
		return m, &HttpError{StatusCode: http.StatusForbidden, Message: "403 Forbidden"}
	}

	return m, nil
}
