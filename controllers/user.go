package controllers

import "github.com/jordan-borges-lark/todo_test/models"

type User[M models.IModel[any]] struct { // implements ICrudController
	CrudController[M]
}
