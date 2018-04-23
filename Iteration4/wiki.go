package main

import (
	"./dao"
	"./models"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("tmpl/view.html", "tmpl/edit.html", "tmpl/list.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var validListPath = regexp.MustCompile("^/(list)(/)?$")
var wikiDao *dao.WikiDAO

func main() {
	wikiDao = &dao.WikiDAO{Server: "127.0.0.1", Database: "gowiki"}
	wikiDao.Connect()
	http.HandleFunc("/", makeHandler(listHandler))
	http.HandleFunc("/list/", makeHandler(listHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}

func makeHandler(fn func(w http.ResponseWriter, r *http.Request, title string)) func(w http.ResponseWriter, p *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			r.URL.Path = "/list"
		}
		if validListPath.FindString(r.URL.Path) == r.URL.Path {
			fn(w, r, "")
		} else {
			m := validPath.FindStringSubmatch(r.URL.Path)
			if m == nil {
				http.NotFound(w, r)
				return
			}
			fn(w, r, m[2])
		}
	}
}

func listHandler(w http.ResponseWriter, r *http.Request, title string) {
	if r.Method == "POST" {
		searchTerm := r.FormValue("searchterm")
		println(searchTerm) //ToDo: For debugging pusposes, will be deletede
		pages, err := wikiDao.FindEntries(searchTerm)
		if err != nil {
			http.NotFound(w, r)
		}
		renderListTemplate(w, "list", pages)
	} else {
		pages, err := wikiDao.ListAllEntries()
		if err != nil {
			http.NotFound(w, r)
		}
		renderListTemplate(w, "list", pages)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := wikiDao.LoadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := wikiDao.LoadPage(title)
	if err != nil {
		p = &models.Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	pageID := r.FormValue("id")
	fmt.Printf("Page Id:%v\n", pageID)
	var p *models.Page
	var err error
	p = &models.Page{Title: title, Body: []byte(body)}
	if pageID == "" {
		err = wikiDao.SavePage(p, true)
	} else {
		p.ID = bson.ObjectIdHex(pageID)
		fmt.Printf("Page Id:%v\n", p.ID)
		err = wikiDao.SavePage(p, false)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *models.Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderListTemplate(w http.ResponseWriter, tmpl string, p *[]models.Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
