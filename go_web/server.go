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
		pythonArgs  string
		commandSSH  string
		outSSH      []byte
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
	// Coping the uploaded file into the server
	f, err := os.OpenFile("./client_upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	// string of the path for the uploaded image
	uploadedImg = "/picSTy/go_web/client_upload/" + handler.Filename
	// Processing which artwork will work with the style transfer
	artworks := []string{
		"alebrijes", "estanque", "guernica",
		"waterfall", "mountains", "ninth",
		"starry", "swing", "vetheuil",
	}

	// processing the form values
	r.ParseForm()
	for _, s := range artworks {
		// validating the string of the selected image
		if s == r.Form.Get("art") {
			// path to the selected artwork
			selectedImg = "/picSTy/go_web/artworks/" + s + ".jpg"
		}
	}
	// this variable has the arguments for the python scripts
	pythonArgs = uploadedImg + " " + selectedImg

	// SSH connection to the Python container
	client, err := goph.New("root", "172.19.0.3", goph.Password("root"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	commandSSH = format("python3 /picSTy/stylet_py/styling.py {{.}}", pythonArgs)
	// Execute your command.

	outSSH, err = client.Run(commandSSH)
	// In case of an error within the Python container, this allows to print it in the terminal
	fmt.Println(string(outSSH))
	if err != nil {
		log.Fatal(err)
	}

	// Execute a template for the HandleFunc
	type StyledPage struct {
		Img string
	}
	// getting the name of the picture, without the extension, for the struct field Img
	processedImg := StyledPage{Img: strings.Split(handler.Filename, ".")[0]}

	styledTemplate, err := template.ParseFiles("template/styled.html")
	if err != nil {
		log.Fatal(err)
	}
	styledTemplate.Execute(w, processedImg)
}

func main() {

	s := http.Server{
		Addr: ":8080",
		// specifying timeouts, set these to handle malicious or broken HTTP clients
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Landing page
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)
	// Handling the form of the landing page
	http.HandleFunc("/styled", uploadImage)
	// Serving the web server
	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}

}
