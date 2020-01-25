package main

import (
	"net/http"
	"fmt"
	"os"
	"io"
 	"path/filepath"
 	"html/template"
)



func HandlerIndex(w http.ResponseWriter, r *http.Request) {
	var tmp = template.Must(template.ParseFiles("views/Index.html"))
    var error = tmp.ExecuteTemplate(w,"Index",nil)
    if error != nil {
        http.Error(w, error.Error(), http.StatusInternalServerError)
    }
}


func HandlerUpload(w http.ResponseWriter, r *http.Request) {
	var tmp = template.Must(template.ParseFiles("views/Upload.html"))
	var error = tmp.ExecuteTemplate(w,"Upload",nil)
    if error != nil {
        http.Error(w, error.Error(), http.StatusInternalServerError)
    }
}


func HandlerProsesUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseMultipartForm(1024); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		alias := r.FormValue("alias")
	
		uploadedFile, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer uploadedFile.Close()
	
		dir, err := os.Getwd()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		filename := handler.Filename
		if alias != "" {
			filename = fmt.Sprintf("%s%s", alias, filepath.Ext(handler.Filename))
		}
		fileLocation := filepath.Join(dir, "files", filename)
		targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer targetFile.Close()
	
		if _, err := io.Copy(targetFile, uploadedFile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", 301)
	}

}




func main() {
	http.HandleFunc("/", HandlerIndex)
	http.HandleFunc("/upload", HandlerUpload)
	http.HandleFunc("/upload/proses", HandlerProsesUpload)

	fmt.Println("server started at localhost:8000")
	http.ListenAndServe(":8000", nil)
}