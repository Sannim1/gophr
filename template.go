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
	templateData["Flash"] = prepareFlash(request)

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

func prepareFlash(request *http.Request) map[string]string {

	flashMessage := request.URL.Query().Get("flash_message")
	if flashMessage == "" {
		return nil
	}

	flash := map[string]string{
		"Message": flashMessage,
	}

	messageType := request.URL.Query().Get("msg_type")
	flash["AlertType"] = getAlertType(messageType)

	return flash
}

func getAlertType(messageType string) (alertType string) {
	validAlerts := map[string]struct{}{
		"success": {},
		"info":    {},
		"danger":  {},
		"warning": {},
	}

	_, isValidAlert := validAlerts[messageType]
	if !isValidAlert {
		return "alert-info"
	}

	return fmt.Sprintf("alert-%s", messageType)
}

var errorTemplate = `
<html>
    <body>
        <h1>Error rendering template: '%s'</h1>
        <p>%s</p>
    </body>
</html>
`
