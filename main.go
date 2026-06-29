package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

// page layout for a wiki page
type Page struct {
	Title string
	Body  []byte
}

// method to save on the page
func (p *Page) save() error {
	// use title as file name to write
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

// read file to load the page
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	bodyData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: bodyData}, nil
}

// vieww handler for viewing a page
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	RenderTemplat(w, "view.html", p)
}

// handler to edit a page
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	RenderTemplat(w, "edit.html", p)
}

// helper function to render response template
func RenderTemplat(w http.ResponseWriter, file string, data *Page) {
	err := tmpl.ExecuteTemplate(w, file, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("%e", err), http.StatusInternalServerError)
		return
	}
}

// handle saving of pages
func SaveHandle(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	page := &Page{Title: title, Body: []byte(body)}
	err := page.save()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
// 	sm := validPath.FindStringSubmatch(r.URL.Path)
// 	if sm == nil {
// 		http.NotFound(w, r)
// 		return "", errors.New("Invalide title page")
// 	}
// 	return sm[2], nil
// }

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

var tmpl *template.Template

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func main() {
	tmpl = template.Must(template.ParseFiles("templates/view.html", "templates/edit.html"))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/save/", makeHandler(SaveHandle))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	fmt.Println("starting sever on http://localhost:8080/view/golang")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
