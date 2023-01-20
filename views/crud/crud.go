package crud

import (
	"github.com/jordan-borges-lark/todo_test/helpers"
	"github.com/jordan-borges-lark/todo_test/models"
	"github.com/jordan-borges-lark/todo_test/views"
)

// CRUD view model
type CrudViewModel struct {
	EntityName  string
	Fields      []Field
	Collections []Collection
}
type Field struct {
	Key   string
	Value interface{}
}
type Collection []CrudViewModel

// Convert the Model map to a CrudViewModel and pass it to GetView/getElement
func GetCrudView[M models.IModel[any]](m M) (string, error) {
	mMap := helpers.ToMap(m)
	vm := getViewModelForModel(m.GetTableName(), mMap)
	vmMap := helpers.ToMap(vm)

	return views.GetView("crud", "", vmMap)
}

// getViewModelForModel recursively populates ViewModel based on Model relationships
func getViewModelForModel(entityName string, m map[string]interface{}) *CrudViewModel {
	vm := CrudViewModel{
		EntityName:  entityName,
		Fields:      []Field{},
		Collections: []Collection{},
	}

	for k, v := range m {
		// Objects treated as array of size 1
		if m, isMap := v.(map[string]interface{}); isMap {
			v = []map[string]interface{}{m}
		}
		array, isArray := v.([]map[string]interface{})

		// Field

		if !isArray {
			vm.Fields = append(vm.Fields, Field{k, v})
			continue
		}

		// Collection

		// Recursively populates until it finds no
		// more collections, which it won't unless
		// you somehow use .Load() ad inf.
		var collection = make([]CrudViewModel, len(array))
		for _, collectionsMap := range array {
			collection = append(
				collection,
				*getViewModelForModel(k, collectionsMap),
			)
		}

		// Add blank collection element to create a blank row
		blankViewModel := &CrudViewModel{}
		if len(collection) > 0 {
			blankViewModel = &collection[0]
		}
		idField := &Field{}
		for _, f := range blankViewModel.Fields {
			if f.Key == "id" {
				idField = &f
				break
			}
		}
		idField.Key = "id"
		idField.Value = 0
		collection = append(
			collection,
			*blankViewModel,
		)

		vm.Collections = append(vm.Collections, collection)
	}

	return &vm
}
