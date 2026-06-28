package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	bodyData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: filename, Body: bodyData}, nil
}

func veiwHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	p = &Page{Title: title}
	RenderTemplat(w, "view.html", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	page := p.Title
	RenderTemplat(w, "edit.html", page)
}

func RenderTemplat(w http.ResponseWriter, file string, data any) {
	err := tmpl.ExecuteTemplate(w, "templates/"+file, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
func SaveHandle(w http.ResponseWriter, r *http.Request) {
}

var tmpl template.Template

func main() {
	tmpl = *template.Must(template.ParseFiles("edit.html"))
	http.HandleFunc("/save", SaveHandle)
	http.HandleFunc("/edit/", editHandler)
	fmt.Println("starting sever on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
