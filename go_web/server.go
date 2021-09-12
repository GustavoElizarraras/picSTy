package main

import (
	"fmt"
	"io"
	"net/http" // Client type por make requests and receive responses
	"os"
	"os/exec"
	"time"
)

// Uploading user image
func uploadImage(w http.ResponseWriter, r *http.Request) {
	// Getting the user file
	r.ParseMultipartForm(32 << 20)
	// Form with name=imgFile
	file, handler, err := r.FormFile("imgFile")
	if err != nil {
		fmt.Println(err)
	}
	// Closing the file
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	// Coping the uploaded file into the server
	f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	// Processing which artwork will work with the style transfer
	artworks := []string{
		"alebrijes", "estanque", "guernica",
		"maya", "mountains", "ninth",
		"starry", "swing", "vetheuil",
	}

	var selectedImg string

	r.ParseForm()
	for _, s := range artworks {
		if s == r.Form.Get("art") {
			selectedImg = s
		}
	}

	exec.Command("python3", "st_execute.py", selectedImg)

}

func main() {

	s := http.Server{
		Addr: ":8080",
		// specifying timeouts, set these to handle malicious or broken HTTP clients
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/u", uploadImage)
	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}

}
