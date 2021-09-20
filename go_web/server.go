package main

import (
	"fmt"
	"io"
	"log"
	"net/http" // Client type por make requests and receive responses
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/melbahja/goph"
)

func format(s string, v interface{}) string {
	t, b := new(template.Template), new(strings.Builder)
	template.Must(t.Parse(s)).Execute(b, v)
	return b.String()
}

// Uploading user image
func uploadImage(w http.ResponseWriter, r *http.Request) {
	var (
		selectedImg string
		uploadedImg string
		files       string
		command     string
		// ssh         string
	)

	// Getting the user file
	r.ParseMultipartForm(32 << 20)
	// Form with name=imgFile
	file, handler, err := r.FormFile("imgFile")
	if err != nil {
		fmt.Println(err)
	}
	// Closing the file
	defer file.Close()
	// fmt.Fprintf(w, "%v", handler.Header)
	// Coping the uploaded file into the server
	f, err := os.OpenFile("./client_upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	uploadedImg = "/picSTy/go_web/client_upload/" + handler.Filename
	// Processing which artwork will work with the style transfer
	artworks := []string{
		"alebrijes", "estanque", "guernica",
		"waterfall", "mountains", "ninth",
		"starry", "swing", "vetheuil",
	}

	r.ParseForm()
	for _, s := range artworks {
		if s == r.Form.Get("art") {
			selectedImg = "/picSTy/go_web/artworks/" + s + ".jpg"
		}
	}

	files = uploadedImg + " " + selectedImg
	fmt.Println(files)
	var o []byte
	client, err := goph.New("root", "172.19.0.3", goph.Password("root"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(client)
	defer client.Close()
	command = format("python3 /picSTy/stylet_py/styling.py {{.}}", files)
	// Execute your command.
	fmt.Println(command)
	o, err = client.Run(command)
	fmt.Println(string(o))
	if err != nil {
		log.Fatal(err)
	}

	type Page struct {
		Img string
	}
	x := Page{Img: strings.Split(handler.Filename, ".")[0]}
	fmt.Println(x)
	t, err := template.ParseFiles("template/styled.html")
	if err != nil {
		log.Fatal(err)
	}

	t.Execute(w, x)
}

func main() {

	s := http.Server{
		Addr: ":8080",
		// specifying timeouts, set these to handle malicious or broken HTTP clients
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fs := http.FileServer(http.Dir("."))

	http.Handle("/", fs)

	http.HandleFunc("/styled", uploadImage)
	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}

}

// func landingHandler(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "./landing_page/index.html")
// }
