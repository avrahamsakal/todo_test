package views

import (
	"bytes"
	"html/template"
	"time"

	"github.com/jordan-borges-lark/todo_test/helpers"
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
	elementName := viewName + "/" + viewName
	if layoutName == "" {
		layoutName = "default"
	}

	elem, err := getElement(elementName, "layouts/"+layoutName, values)
	if err != nil {
		return elem, err
	}
	// @TODO: Add a file exists check for the JS and CSS includes here, right now I'm letting it 404
	// @TODO: Move this code to getElement and have it know whether it's a layout?
	/*elem += `
		<script src="/` + elementName + `.js"></script>
		<style src="/` + elementName + `.css"></style>
	`*/

	return elem, nil
}

// getElement injects values into the template matching elementName.htm
func getElement(elementName string, layoutName string, values map[string]interface{}) (string, error) {
	t, err := template.ParseFiles(viewsDirName + "/" + layoutName + ".htm")
	if err != nil {
		return "", err
	}
	/*t, err = t.Clone()
	if err != nil {
		return "", err
	}*/
	t, err = t.ParseFiles(viewsDirName + "/" + elementName + ".htm")
	if err != nil {
		return "", err
	}

	// Add values
	values["LayoutName"] = layoutName
	values["ElementName"] = elementName
	values["NoCache"] = time.Now().String()
	// Should be unreachable
	if values["EntityName"] == nil || values["FriendlyName"] == nil {
		values["EntityName"] = helpers.ToSnakeCase(values["EntityName"].(string))
		values["FriendlyName"] = helpers.ToSnakeCase(values["EntityName"].(string))	
	}
	
	// Apply values to template
	buf := new(bytes.Buffer)
	if err = t.ExecuteTemplate(buf, "default.htm", values); err != nil {
		return "", err
	}

	return buf.String(), nil
}
