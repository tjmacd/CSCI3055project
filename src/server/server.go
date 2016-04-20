package main

import (
	"html/template"
	"net/http"
	"io/ioutil"
	"io"
	//"fmt"
	"encoding/xml"
	//"os"
)

type PictureList struct {
	Pictures []Picture `xml:"Picture"`
}

type Picture struct {
	Title    string
	Filename string
	Comments []Comment `xml:"Comment"`
}

type Comment struct {
	Name    string
	Message string
}

var templates = template.Must(template.ParseFiles("upload.html", "index.html"))
var pics = loadXML("picts.xml")

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

func renderTemplate(w http.ResponseWriter, tmpl string, data []byte) {
	fList := template.HTML(data)
	err := templates.ExecuteTemplate(w, tmpl+".html", fList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		renderTemplate(w, "upload", nil)
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
	pics.Pictures = append(pics.Pictures, Picture{Filename: t.Name()})
	saveXML("picts.xml", pics)
	http.Redirect(w, r, "/", 302)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image")
	path := "pictures/image-"+r.FormValue("id")
	http.ServeFile(w, r, path)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	/*
	files, err := ioutil.ReadDir("pictures/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	*/
	var list []byte
	for _, pic := range pics.Pictures {
		filename := pic.Filename[9:]
		list = append(list, "<a href=\"/view?id="+filename[6:]+"\">"+
			"<img src=\"data:image/jpg;base64,"+filename+"\" alt=\""+filename+"\" style=\"width:420px;height:420px;border:0\"></a>\n"...)
		//list = append(list, file.Name()...)
		//list = append(list, "&#13;&#10;"...)
	}
	//fmt.Printf("%s\n", list)
	renderTemplate(w, "index", list)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":8080", nil)
}
