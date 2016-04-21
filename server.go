package main

import (
	"html/template"
	"net/http"
	"io/ioutil"
	"io"
	"fmt"
	"encoding/xml"
	//"os"
	"encoding/base64"
	"regexp"
)

type PictureList struct {
	Pictures []Picture `xml:"Picture"`
}

func (l *PictureList) GetPicture(filename string) Picture {
	for _, pic := range l.Pictures {
		if pic.Filename == filename {
			return pic
		}
	}
	return Picture{}
}

func (l *PictureList) AddComment(ID string, comment Comment) {
	for i, pic := range l.Pictures {
		if pic.ID == ID {
			l.Pictures[i].Comments = append(l.Pictures[i].Comments, comment)
			return
		}
	}
}

type Picture struct {
	Title    string
	Filename string
	Image    string
	ID       string
	Comments []Comment `xml:"Comment"`
}

type Comment struct {
	Name    string
	Message string
}

var templates = template.Must(template.ParseFiles("upload.html", "index.html", "view.html"))
var pics = loadXML("picts.xml")
var validPath = regexp.MustCompile("^/(index|upload|comment|view)/([a-zA-Z0-9]+)$")

func loadXML(filename string) PictureList {
	xmlFile, _ := ioutil.ReadFile(filename)
	var pics PictureList
	xml.Unmarshal(xmlFile, &pics)
	return pics
}

func saveXML(filename string, pics PictureList) {
	bytes, _ := xml.Marshal(&pics)
	ioutil.WriteFile(filename, bytes, 0600)
}

func renderTemplate(w http.ResponseWriter, tmpl string, pic *Picture) {
	//fList := template.HTML(data)
	
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		err := templates.ExecuteTemplate(w, "upload.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	f, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer f.Close()
	t, err := ioutil.TempFile("./pictures", "image-")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer t.Close()
	if _, err := io.Copy(t, f); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	filename := t.Name()
	//fmt.Println(filename, filename[15:])
	pics.Pictures = append(pics.Pictures, Picture{Filename: filename, ID: filename[15:]})
	saveXML("picts.xml", pics)
	http.Redirect(w, r, "/", 302)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	path := "pictures/image-"+r.FormValue("id")
	pic := pics.GetPicture(path)
	/*
	if pic == nil {
		http.Error(w, "File not found", http.StatusInternalServerError)
		return
	}*/
	if pic.Title == "" {
		pic.Title = pic.Filename[9:]
	}
	bytes, _ := ioutil.ReadFile(path)
	pic.Image = base64.StdEncoding.EncodeToString(bytes)
	err := templates.ExecuteTemplate(w, "view.html", pic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fPics := pics
	for i, pic := range fPics.Pictures {
		bytes, _ := ioutil.ReadFile(pic.Filename)
		fPics.Pictures[i].Image = base64.StdEncoding.EncodeToString(bytes)
	}
	err := templates.ExecuteTemplate(w, "index.html", fPics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		fmt.Println(r.URL.Path)
		http.NotFound(w, r)
		return
	}
	name := r.FormValue("name")
	message := r.FormValue("message")
	
	comment := Comment{Name: name, Message: message}
	pics.AddComment(m[2], comment)
	saveXML("picts.xml", pics)
	http.Redirect(w, r, "/view?id="+m[2], http.StatusFound)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/comment/", commentHandler)
	http.ListenAndServe(":8080", nil)
}
