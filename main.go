package main

import (
	"fmt"
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

func main() {
	p1 := Page{Title: "TestPage", Body: []byte("this is a simple page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}
