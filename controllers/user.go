package controllers

import (
	//"encoding/json"
	"net/http"
	//"strconv"

	//"github.com/gorilla/mux"
	//"github.com/jordan-borges-lark/todo_test/models"
	//"github.com/gorilla/mux"
)

type User struct { // implements ICrudController
	CrudController
}

/*func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	// @TODO add to r.Context to tell Update this is a create operation
	u.Update(w, r) // Create must be called without ID, or with nil/zero ID
}*/

func (u *User) Read(w http.ResponseWriter, r *http.Request) {
	
	if r.Header.Get("Content-Type") == "application/json" {
		u.CrudController.Read(w, r)
		return
	}

	user := u.CrudController.ReadBase(w, r)
	OutputJson(w, user)
	// @TODO: Display view using HTML template strings
}

func (u *User) Update(w http.ResponseWriter, r *http.Request) {
	// Auth action
	/*if u.Model.Id != session.getLoggedInUserId() {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}*/
	u.CrudController.Update(w, r)
}

// Different ways to grab user
//vars := mux.Vars(r)
//userId := vars["userId"] // j/k would be insecure
//user, err = models.Get(u.Database, User{Id:userId})
//email := session.getLoggedInUserEmail()
//user, err := models.GetUserByEmail(email)
//user, err := models.Get(u.Database, User{Id:session.getLoggedInUserId()})
//user := models.User{Id: 123, Email: "SESSIONID"}

/*func (u *User) Read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	user, err := models.Get(u.Database, models.User{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // Not really, outputting DB err would be insecure
		return
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(userBytes)
}

func (u *User) Update(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.GetId() == 0 { // Better alternative is Id *int == nil
		u.Create(w, r)
		return
	}
	// Auth action
	/*if user.Id != session.getLoggedInUserId() {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}*//*
	if _, err := models.Update(u.Database, *user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (u *User) Delete(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 Not Found", http.StatusNotFound)
}
*/