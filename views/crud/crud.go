package crud

import (
	"fmt"
	"strings"

	"github.com/jordan-borges-lark/todo_test/helpers"
	"github.com/jordan-borges-lark/todo_test/models"
	"github.com/jordan-borges-lark/todo_test/views"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// CRUD view model
type CrudViewModel struct {
	Id           string // int treats 0 as blank??
	EntityName   string
	FriendlyName string
	Fields       []Field
	Collections  []Collection
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

func GetCrudViewAll[M models.IModel[any]](m []M) (string, error) {
	// Index should pass to the view a blank viewmodel with m[0] in its .Collections
	return GetCrudView(m[0])
}

// getViewModelForModel recursively populates ViewModel based on Model relationships
func getViewModelForModel(entityName string, m map[string]interface{}) *CrudViewModel {
	// Convert model name to entityName (snake case) format
	if strings.ToUpper(entityName[0:1]) == entityName[0:1] {
		entityName = helpers.ToSnakeCase(strings.TrimSuffix(entityName, "s"))
	}
	friendlyName := cases.Title(language.English).String(
		strings.ReplaceAll(entityName, "_", " "),
	)

	vm := CrudViewModel{
		Id:           fmt.Sprint(int(m["Id"].(float64))),
		EntityName:   entityName,
		FriendlyName: friendlyName,
		Fields:       []Field{},
		Collections:  []Collection{},
	}

	/*
		   Hit list:
		x  Change ItemListItem Text property to a Unicode bullet point
		x  Add routes for "ItemListItem" (i.e. reflection)
		x  Add value friendlyName Split(entityName, "_")[-1]
		o  Add value friendlyNameTitle Title(friendly)
		x  Change EntryName in TitleCase to snake case
		x  Change EntryName to TitleCase and call it EntryNameTitle
		x  Add 16x16 minicons named after friendlyName
		x  Change h2 label to friendlyNameTitle
		x  Add link from 16x16 icons to entityName/id route
		x  Objects (but NOT collections of one) don't need Load to add a blank
		x  if v.(string) needs to say if m["Id"] == 0 then v = "Example " + k
		   Trace out delete.onclick text.onchange event to ensure create/delete
		   Once Save functionality is working, can mysqldump the DB;
		   		otherwise just put the create statements in text files
	*/

	for k, v := range m {
		if k == "Id" && entityName != "item_list" {
			continue
		}

		// Objects treated as array of size 1
		m, isMap := v.(map[string]interface{})
		if isMap {
			v = []map[string]interface{}{m}
		}
		// Unbox maps
		if ifaceArray, isArray := v.([]interface{}); isArray {
			vMap := make([]map[string]interface{}, len(ifaceArray))
			for i, iface := range ifaceArray {
				vMap[i] = iface.(map[string]interface{})
			}
			v = vMap
		}

		array, isArray := v.([]map[string]interface{})

		// Field

		if !isArray {
			// For now we only handle string fields, but we need
			// radios for booleans, calendars for dates, etc.
			if _, isStr := v.(string); isStr {
				if m["Id"] == "0" && v == "" { // Fill in blank model with example values
					v = "Example " + k
				}
				vm.Fields = append(vm.Fields, Field{k, v})
			}
			continue
		}

		// Collection

		// Recursively populates until it finds no
		// more collections, which it won't unless
		// you somehow use .Load() ad inf.
		collection := []CrudViewModel{}
		for _, collectionsMap := range array {
			collection = append(
				collection,
				*getViewModelForModel(k, collectionsMap),
			)
		}

		// If this is a real collection and not a virtual collection
		// created for an object, then add a blank collection element
		// to create a blank row
		if !isMap {
			blankViewModel := &CrudViewModel{}
			if len(collection) > 0 {
				blankViewModel = &collection[0]
			}
			var idField *Field
			found := false
			for _, f := range blankViewModel.Fields {
				if f.Key == "Id" {
					idField, found = &f, true
					break
				}
			}
			if !found {
				blankViewModel.Fields = append(
					blankViewModel.Fields, Field{
						Key: "Id", Value: 0,
					})
			} else {
				idField.Value = 0
			}

			collection = append(
				collection,
				*blankViewModel,
			)
		}

		vm.Collections = append(vm.Collections, collection)
	}

	return &vm
}
