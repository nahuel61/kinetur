package main

import (
	"html/template"
	"path/filepath"
	"time"
	_ "time"
	"tp-ISA-go.org/kinetur/pkg/forms"
	"tp-ISA-go.org/kinetur/pkg/models"
)

//templateData struct funciona como estructura para guardar los datos dinamicos que le paso a los templates html
type templateData struct {
	Año               int
	Flash             string
	Form              *forms.Form
	Paciente          *models.Pacientes
	Pacientes         []*models.Pacientes
	AuthenticatedUser *models.Pacientes
	Profesional       *models.Profesionales
	Profesionales     []*models.Profesionales
	CSRFToken         string
}

func humanDate(t time.Time) string {
	// Return the empty string if time has the zero value.
	if t.IsZero() {
		return ""
	}
	// Convert the time to UTC before formatting it.
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Inicializo un nuevo map para que actúe como caché..
	cache := map[string]*template.Template{}
	// Utilice la función filepath.Glob para obtener un slice de todas las rutas de archivo con
	//la extensión '.page.tmpl'. Básicamente, esto nos da un slice de todas las plantillas de 'página' para la aplicación.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}
	// Recorra las páginas una por una.
	for _, page := range pages {
		// Extraiga el nombre del archivo (como 'home.page.tmpl') del archivo completo pat
		//// y asígnelo a la variable de nombre.
		name := filepath.Base(page)
		// Analizar el archivo de template de página en un conjunto de templates.
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Utilice el método ParseGlob para agregar cualquier template de 'diseño' al conjunto de templates
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		// Utilice el método ParseGlob para agregar cualquier template 'parcial' al
		// conjunto de templates
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		// Agrega el conjunto de templates al cahce, usando el nombre de la página
		// (como 'home.page.tmpl') como key.
		cache[name] = ts
	}
	// Devuelve el map.
	return cache, nil
}
