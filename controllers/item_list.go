package controllers

import "github.com/jordan-borges-lark/todo_test/models"

type ItemList[M models.IModel[any]] struct {
	CrudController[M]
}
