package main

import ( 
	"html/template"
	"io/ioutil"
	"net/http"
	"log"
)

type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename,p.Body,0600)
}

func loadPage(title string) (*Page,error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title,Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w,"view",p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p=&Page{Title: title}
	}
	renderTemplate(w,"edit",p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+ title, http.StatusFound) 
}

var templates = template.Must(template.ParseFiles("edit.html","view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html",p)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/view/",viewHandler)
	http.HandleFunc("/edit/",editHandler)
	http.HandleFunc("/save/",saveHandler)
        log.Fatal(http.ListenAndServe(":8080",nil))
}

