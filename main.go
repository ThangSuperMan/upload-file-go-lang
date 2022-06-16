package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("homepage")
	t, err := template.ParseFiles("./views/index.html")

	if err != nil {
		fmt.Println("Error when trying to render the html page")
		fmt.Println(err)
		return
	}

	t.Execute(w, nil)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, fileHeader, err := r.FormFile("image")

	if err != nil {
		fmt.Println("Error Retrieving the File")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Println("Info of file")
	fmt.Println("file name: ", fileHeader.Filename)
	fmt.Println("size file: ", fileHeader.Size)

	tempFile, err := ioutil.TempFile("temp-images", "upload-*.jpeg")
	if err != nil {
		fmt.Println("Error temFile")
		http.Error(w, "Error Internal Server", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	defer tempFile.Close()

	// Read all content of file
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error read content of fille")
		http.Error(w, "Error Internal Server", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// Write bytes of file to the the fils just created
	tempFile.Write(fileBytes)
	http.Redirect(w, r, "http://localhost:3000/?success=true", http.StatusSeeOther)
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		uploadFile(w, r)
	}
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/upload", upload)
	http.ListenAndServe(":3000", nil)
}

func main() {
	setupRoutes()
}
