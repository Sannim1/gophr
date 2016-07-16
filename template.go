package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

var layoutFuncs = template.FuncMap{
	"yield": func() (string, error) {
		return "", fmt.Errorf("yield called inappropriately")
	},
}

var layout = template.Must(
	template.
		New("layout.html").
		Funcs(layoutFuncs).
		ParseFiles("templates/layout.html"),
)

var loadedTemplates = template.Must(
	template.
		New("t").
		ParseGlob("templates/**/*.html"),
)

func RenderTemplate(
	responseWriter http.ResponseWriter,
	request *http.Request,
	templateName string,
	templateData map[string]interface{},
) {
	if templateData == nil {
		templateData = map[string]interface{}{}
	}
	templateData["CurrentUser"] = RequestUser(request)
	templateData["Flash"] = request.URL.Query().Get("flash")

	funcs := template.FuncMap{
		"yield": func() (template.HTML, error) {
			buffer := bytes.NewBuffer(nil)
			err := loadedTemplates.ExecuteTemplate(buffer, templateName, templateData)

			return template.HTML(buffer.String()), err
		},
	}

	layoutClone, _ := layout.Clone()
	layoutClone.Funcs(funcs)

	err := layoutClone.Execute(responseWriter, templateData)

	// err := loadedTemplates.ExecuteTemplate(responseWriter, templateName, templateData)

	if err != nil {
		http.Error(
			responseWriter,
			fmt.Sprintf(errorTemplate, templateName, err),
			http.StatusInternalServerError,
		)
	}
}

var errorTemplate = `
<html>
    <body>
        <h1>Error rendering template: '%s'</h1>
        <p>%s</p>
    </body>
</html>
`
