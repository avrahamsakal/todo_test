// @TODO: Use this package to implement
// view models and have strongly-typed
// parameters passed to the template engine
package views

import (
	"bytes"
	"html/template"
)

const viewsDirName = "./views"

// GetView gets a string version of rendered view,
//
//	then passes it as a parameter to layout,
//	and returns the string version of rendered layout
func GetView(
	viewName string,
	layoutName string,
	values map[string]interface{},
) (string, error) {
	values["entityName"] = viewName
	elem, err := getElement(viewName, values)
	if err != nil {
		return elem, err
	}

	if layoutName == "" {
		layoutName = "layout"
	}
	layout, err := getElement(layoutName, map[string]interface{}{
		"body":       elem,
		"entityName": viewName,
	})
	if err != nil {
		return layout, err
	}

	return layout, nil
}

// getElement injects values into the template matching elementName.htm
func getElement(elementName string, values map[string]interface{}) (string, error) {
	//htm, err := getHtmlFile(elementName)

	t, err := template.ParseFiles(
		viewsDirName + "/" + elementName + ".htm",
	)
	if err != nil {
		return "", err
	}

	// Apply values to template
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, values); err != nil {
		return "", err
	}
	return buf.String(), nil

	//@TODO: It would be super trivial to test whether string "elem"
	// has any parameters that weren't found in what we passed for
	// values argument, then see if it matches another view name,
	// so we can recursively call GetView and create a whole
	// element tree, e.g.:
	//		<div>{{element_name("Text", 0.994, {"k":"v"})}}</div>
	// If we detect this, we search if element_name.htm
	// file exists, add these values to the values collection,
	// and run GetView("element_name"), which could itself have elements.

	// @TODO: Also, we need to error out if the view template was
	// expecting parameters and they weren't supplied and
	// it's also not the name of a view
}

/*func getHtmlFile(elementName string) (string, error) {
	if byt, err := os.ReadFile(
		viewsDirName + "/" + elementName + ".htm",
	); err != nil {
		return "", nil
	} else {
		return string(byt), nil
	}
}*/
