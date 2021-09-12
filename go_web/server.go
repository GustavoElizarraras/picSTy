package main

import (
	"fmt"
	"io"
	"net/http" // Client type por make requests and receive responses
	"os"
	"time"
)

// Uploading user image
func uploadImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("imgFile")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}

func selectArtwork(w http.ResponseWriter, r *http.Request) {
	artworks := []string{
		"alebrijes", "estanque", "guernica",
		"maya", "mountains", "ninth",
		"starry", "swing", "vetheuil",
	}

	for x, s := range artworks {
		fmt.Println(s)
		fmt.Println(x)
		// if s == r.Form.Get("art") {
		// }
	}

	fmt.Println("no vl")
}

func main() {

	s := http.Server{
		Addr: ":8080",
		// specifying timeouts, set these to handle malicious or broken HTTP clients
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		// Handler:      HelloHandler{},
	}

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/upload", uploadImage)
	http.HandleFunc("/art", selectArtwork)
	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}

}
