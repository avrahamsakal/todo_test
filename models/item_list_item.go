package models

import (
	"github.com/jmoiron/sqlx"
)

func (ili ItemListItem) GetTableName() string {
	return "item_list_item"
}

type ItemListItem struct {
	Model

	ItemListId int `db:"item_list_id"`
	ItemList   ItemList

	Text string `db:"○"` // We never refer to this in the code so it doesn't matter what it's named
}

func (ili ItemListItem) Get() IModel[any]            { return ili }
func (ili ItemListItem) SetId(id int) IModel[any]    { ili.Id = id; return ili }
func (ili ItemListItem) CanUserRead(userId int) bool { return ili.CanUserWrite(userId) }
func (ili ItemListItem) CanUserWrite(userId int) bool {
	return ili.ItemList.UserId == userId
}
func (ili ItemListItem) Load(db *sqlx.DB, flags ...bool) (IModel[any], error) {
	// lazy load by default
	il, err := Read(db, ItemList{Model: Model{Id: ili.ItemListId}})
	if err != nil {
		return ili, err
	}
	ili.ItemList = il

	// Use cascading load instead
	const lazy = /*, flag2, flag3*/ 0 /*, val2, val3*/
	if len(flags) != 0 && !flags[lazy] {
		if m, err := il.Load(db, flags...); err != nil {
			return ili, err
		} else {
			ili.ItemList = m.(ItemList)
		}
	}

	return ili, nil
}

func GetItemListsItems[M Model](db *sqlx.DB, itemListId int) ([]ItemListItem, error) {
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
