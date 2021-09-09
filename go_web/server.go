package main

import (
	"fmt"
	"io"
	"net/http" // Client type por make requests and receive responses
	"os"
	"time"
)

// type HelloHandler struct{}

// func (hh HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Hello!\n"))
// }

// type Page struct{
// 	Title string
// 	Body []byte
// }

// func loadPage(title string) {

// }

// func picstyHandler(w http.ResponseWriter, r *http.Request) {
// 	// picsty, _:= ioutil.ReadFile("index.html")
// 	t, _ := template.ParseFiles("index.html")
// }

// Uploading image
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
	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}

}
