package models

import (
	"github.com/jmoiron/sqlx"
)

func (m User) GetTableName() string {
	return "user"
}

type User struct {
	Model

	ItemLists []ItemList

	Email string `db:"text"`
}

func (u User) Get() IModel[any]             { return u }
func (u User) SetId(id int) IModel[any]     { u.Id = id; return u }
func (u User) CanUserRead(userId int) bool  { return u.CanUserWrite(userId) }
func (u User) CanUserWrite(userId int) bool { return u.Id == userId }
func (u User) Load(db *sqlx.DB, flags ...bool) (IModel[any], error) {
	if u.Id == 0 { // Load a blank model
		u.ItemLists = []ItemList{{}}
		return u, nil
	}

	return u, nil
}

// GetUserByEmail can retrieve the User based on session ID or email
func GetUserByEmail[M Model](db *sqlx.DB, email string) (*User, error) {
	var user *User
	err := db.Get(user, `
		SELECT id FROM user
		WHERE email = "?"
	`, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
