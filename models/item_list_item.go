package models

import (
	"github.com/jmoiron/sqlx"
)

func (ili ItemListItem) GetTableName() string {
	return "item_list_item"
}
func (ili ItemListItem) GetId() int                  { return ili.Id }
func (ili ItemListItem) CanUserRead(userId int) bool { return ili.CanUserWrite(userId) }
func (ili ItemListItem) CanUserWrite(userId int) bool {
	return ili.ItemList != nil && ili.ItemList.UserId == userId
}

type ItemListItem struct {
	Model
	ItemListId int `db:"item_list_id"`
	ItemList   *ItemList

	Text string `db:"text"`
}

func (ili ItemListItem) Load(db *sqlx.DB) (*ItemListItem, error) {
	// lazy load by default
	if il, err := Get(db, ItemList{Model: Model{Id: ili.ItemListId}}); err != nil {
		return &ili, err
	} else {
		ili.ItemList = il
	}
	/* // Use cascading load instead
	} else if m, err := il.Load(db); err != nil { 
		return &ili, err
	} else {
		ili.ItemList = m.(*ItemList)
	}*/

	return &ili, nil
}

func GetItemListsItems(db *sqlx.DB, itemListId int) ([]ItemListItem, error) {
	items := []ItemListItem{}

	err := db.Select(&items, `
		SELECT * FROM `+ItemListItem{}.GetTableName()+`
		WHERE item_list_id = ?
	`, itemListId)
	if err != nil {
		return nil, err
	}

	return items, nil
}
