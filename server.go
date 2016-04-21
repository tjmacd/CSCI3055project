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
var validPath = regexp.MustCompile("^/(index|upload|comment|view)?$")

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	t, err := ioutil.TempFile("./pictures", "image-")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer t.Close()
	if _, err := io.Copy(t, f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filename := t.Name()
	pics.Pictures = append(pics.Pictures, Picture{Filename: filename, ID: filename[15:]})
	saveXML("picts.xml", pics)
	http.Redirect(w, r, "/index", 302)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	path := "pictures/image-"+r.FormValue("id")
	pic := pics.GetPicture(path)
	if pic.Title == "" {
		pic.Title = pic.Filename[9:]
	}
	bytesStr, _ := encodeImage(pic)
	pic.Image = bytesStr
	err := templates.ExecuteTemplate(w, "view.html", pic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func encodeImage(pic Picture) (string, error) {
	bytes, err := ioutil.ReadFile(pic.Filename)
	return base64.StdEncoding.EncodeToString(bytes), err
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fPics := pics
	for i, pic := range fPics.Pictures {
		bytesStr, _ := encodeImage(pic)
		fPics.Pictures[i].Image = bytesStr
	}
	err := templates.ExecuteTemplate(w, "index.html", fPics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	name := r.FormValue("name")
	message := r.FormValue("message")
	comment := Comment{Name: name, Message: message}
	pics.AddComment(r.FormValue("id"), comment)
	saveXML("picts.xml", pics)
	http.Redirect(w, r, "/view?id="+r.FormValue("id"), http.StatusFound)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index", http.StatusFound)
}

func makeHandler(fn func (http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func main() {
	fmt.Println("Now serving image server...")
	http.HandleFunc("/", makeHandler(rootHandler))
	http.HandleFunc("/index", makeHandler(indexHandler))
	http.HandleFunc("/view", makeHandler(viewHandler))
	http.HandleFunc("/upload", makeHandler(uploadHandler))
	http.HandleFunc("/comment", makeHandler(commentHandler))
	http.ListenAndServe(":8080", nil)
}
