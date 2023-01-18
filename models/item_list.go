package models

import (
	"github.com/jmoiron/sqlx"
)

func (il ItemList) GetTableName() string { // If I make this return a constant then I have to change the func definition anyway because there were be multiple constants named tableName in the package
	return "item_list"
}
func (il ItemList) GetId() int                   { return il.Id }
func (il ItemList) CanUserRead(userId int) bool  { return il.CanUserWrite(userId) }
func (il ItemList) CanUserWrite(userId int) bool { return il.UserId == userId }

type ItemList struct {
	Model
	UserId int `db:"user_id"`
	User   *User

	Text string `db:"text"`
}

func (il ItemList) Load(db *sqlx.DB) (IModel, error) {
	if u, err := Get(db, User{Model: Model{Id: il.UserId}}); err != nil {
		return &il, err
	} else {
		il.User = u
	}

	return &il, nil
}

func GetItemListIds(db *sqlx.DB, userId int) ([]int, error) {
	itemListIds := []int{}

	err := db.Select(&itemListIds, `
		SELECT id FROM item_list
		WHERE item_list.user_id = ?
	`, userId)
	if err != nil {
		return nil, err
	}

	return itemListIds, nil
}
