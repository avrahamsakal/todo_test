package models

import (
	"github.com/jmoiron/sqlx"
)

func (il ItemList) GetTableName() string { // If I make this return a constant then I have to change the func definition anyway because there were be multiple constants named tableName in the package
	return "item_list"
}

type ItemList struct {
	Model

	UserId int `db:"user_id"`
	User   User

	ItemListItems []ItemListItem

	Name string `db:"name"`
}

func (il ItemList) Get() IModel[any]             { return il }
func (il ItemList) SetId(id int) IModel[any]     { il.Id = id; return il }
func (il ItemList) CanUserRead(userId int) bool  { return il.CanUserWrite(userId) }
func (il ItemList) CanUserWrite(userId int) bool { return il.UserId == userId }
func (il ItemList) Load(db *sqlx.DB, flags ...bool) (il2 IModel[any], err error) {
	// Add a blank model
	if len(il.ItemListItems) == 0 {
		il.ItemListItems = []ItemListItem{{Text:"Example text"}}
	}
	if il.Id == 0 {
		return il, nil
	}

	il.User, err = Read(db, User{Model: Model{Id: il.UserId}})
	if err != nil {
		return il, err
	}

	il.ItemListItems, err = Where[ItemListItem](db, map[string]interface{}{
		"item_list_item": il.Id,
	}, 50) // @TODO: Need to implement pagination in Where[]() to see the rest
	if err != nil {
		return il, err
	}
	if len(il.ItemListItems) == 0 {
		il.ItemListItems = []ItemListItem{}
	}

	return il, nil
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
