package main

import (
	"html/template"
	"net/http"
	"io/ioutil"
	"io"
	"fmt"
)

var templates = template.Must(template.ParseFiles("upload.html", "index.html"))

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
	t, err := ioutil.TempFile(".", "image-")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer t.Close()
	if _, err := io.Copy(t, f); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/view?id="+t.Name()[6:], 302)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, "pictures/image-"+r.FormValue("id"))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("pictures/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var list []byte
	for _, file := range files {
		filename := file.Name()
		list = append(list, "<a href=\"/view?id="+filename[6:]+"\">"+
			"<img src=\"data:image/jpg;base64,"+filename+"\" alt=\""+filename+"\" style=\"width:420px;height:420px;border:0\"></a>\n"...)
		//list = append(list, file.Name()...)
		//list = append(list, "&#13;&#10;"...)
	}
	fmt.Printf("%s\n", list)
	renderTemplate(w, "index", list)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/view", viewHandler)
	http.ListenAndServe(":8080", nil)
}
