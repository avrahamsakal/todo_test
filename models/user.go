package models

import (
	"github.com/jmoiron/sqlx"
)

func (u User) GetTableName() string {
	return "user"
}
type User struct {
	Model
	Email  string `db:"text"`
}
func (u User) CanUserRead(userId int) bool  { return u.CanUserWrite(userId) }
func (u User) CanUserWrite(userId int) bool { return u.Id == userId }


// GetUserByEmail can retrieve the User based on session ID or email
func GetUserByEmail(db *sqlx.DB, email string) (*User, error) {
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
