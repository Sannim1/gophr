package main

import(
    "fmt"
    "bytes"
    "net/http"
    "html/template"
)

var layoutFuncs = template.FuncMap {
    "yield": func () (string, error)  {
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

func RenderTemplate (
    responseWriter http.ResponseWriter,
    request *http.Request,
    templateName string,
    templateData interface{},
) {
    funcs := template.FuncMap {
        "yield": func () (template.HTML, error) {
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
