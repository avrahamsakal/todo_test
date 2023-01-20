package models

import (
	"github.com/jmoiron/sqlx"
)

func (m User) GetTableName() string {
	return "user"
}

type User struct {
	Model
	Email string `db:"text"`
}

// func (m User) GetId() int                { return m.Id }
func (m User) Get() IModel[any]             { return m }
func (m User) SetId(id int) IModel[any]     { m.Id = id; return m }
func (u User) CanUserRead(userId int) bool  { return u.CanUserWrite(userId) }
func (u User) CanUserWrite(userId int) bool { return u.Id == userId }
func (m User) Load(db *sqlx.DB, flags ...bool) (IModel[any], error) {
	return m, nil
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
