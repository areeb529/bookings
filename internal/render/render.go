package render

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/areeb529/bookings/internal/config"
	"github.com/areeb529/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var pathToTemplates = "./templates"
var functions = template.FuncMap{}

// NewTemplates sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	//get requested template template from cache
	t, ok := tc[tmpl]
	if !ok {
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	_ = t.Execute(buf, td)

	//render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.hbs from templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.hbs", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		//ts, err := template.New(name).ParseFiles(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.hbs", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.hbs", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}

// var tc = make(map[string]*template.Template)

// func RenderTemplate(w http.ResponseWriter, t string) {
// 	var tmpl *template.Template
// 	var err error

// 	//Check to see if we already have the template in our cache

// 	_, inMap := tc[t]

// 	if !inMap {
// 		// creating template and adding to cache
// 		log.Println("creating template and adding to cache")
// 		err = createTemplateCache(t)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	} else {
// 		log.Println("using cached template")
// 	}

// 	tmpl = tc[t]

// 	err = tmpl.Execute(w, nil)

// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func createTemplateCache(t string) error {
// 	templates := []string{
// 		fmt.Sprintf("./templates/%s", t),
// 		"./templates/base.layout.hbs",
// 	}

// 	// parse the template

// 	tmpl, err := template.ParseFiles(templates...)
// 	if err != nil {
// 		return err
// 	}

// 	// add template to cache (map)
// 	tc[t] = tmpl

// 	return nil
// }
