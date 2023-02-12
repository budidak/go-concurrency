package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
)

var pathToTemplates = "./cmd/web/templates"

// TemplateData structure will be used for passing data from our codebase to the templates.
type TemplateData struct {
	StringMap     map[string]string
	IntMap        map[string]int
	FloatMap      map[string]float64
	Data          map[string]any // any = interface{}
	Flash         string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
}

// render executes go templates on the browser
func (app *Config) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) {
	// these are common gotemplates
	partials := []string{
		fmt.Sprintf("%s/base.layout.gotmpl", pathToTemplates),
		fmt.Sprintf("%s/header.partial.gotmpl", pathToTemplates),
		fmt.Sprintf("%s/navbar.partial.gotmpl", pathToTemplates),
		fmt.Sprintf("%s/footer.partial.gotmpl", pathToTemplates),
		fmt.Sprintf("%s/alerts.partial.gotmpl", pathToTemplates),
	}

	// add given gotemplate & common templates to the templateSlice
	var templateSlice []string

	templateSlice = append(templateSlice, fmt.Sprintf("%s/%s", pathToTemplates, t))
	templateSlice = append(templateSlice, partials...)

	if td == nil {
		td = &TemplateData{}
	}

	// parse the templates
	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// execute the parsed file and add data
	if err := tmpl.Execute(w, app.AddDefaultData(td, r)); err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// AddDefaultData adds the session information to the template data
func (app *Config) AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	// we store session data as key:value pairs... and get them from the current request context.
	// "flash", "warning", "error", "userID" are the keys.
	td.Flash = app.Session.PopString(r.Context(), "flash")     // flash message
	td.Warning = app.Session.PopString(r.Context(), "warning") // warning message
	td.Error = app.Session.PopString(r.Context(), "error")     // error message
	td.Now = time.Now()

	if app.IsAuthenticated(r) {
		td.Authenticated = true
		// TODO - get more user information
	}

	return td
}

// IsAuthenticated returns true if userID key exists in session
func (app *Config) IsAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "userID")
}
