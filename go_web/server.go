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
	fmt.Fprintf(w, "%v", handler.Header)
	// Coping the uploaded file into the server
	f, err := os.OpenFile("./client_upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	uploadedImg = "../client_upload/" + handler.Filename
	// Processing which artwork will work with the style transfer
	artworks := []string{
		"alebrijes", "estanque", "guernica",
		"waterfall", "mountains", "ninth",
		"starry", "swing", "vetheuil",
	}

	r.ParseForm()
	for _, s := range artworks {
		if s == r.Form.Get("art") {
			selectedImg = "../artwork/" + s + ".jpg"
		}
	}

	client, err := goph.New("root", "172.19.0.3", goph.Password(""))
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
	files = " " + uploadedImg + " " + selectedImg
	command = format("python3 /picSTy/stylet_py/styling.py {{.}}", files)
	// Execute your command.
	_, err = client.Run(command)

	// ssh = "sshpass -p 'root' ssh root@172.19.0.3"
	// fmt.Println(uploadedImg, selectedImg, ssh)
	// cmd := exec.Command("go", "version")
	// cmd := exec.Command(ssh, "python3", "/picSTy/stylet_py/styling.py", uploadedImg, selectedImg)

	// err = cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

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
