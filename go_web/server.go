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

func userFileForm(r *http.Request) (string, string) {
	// Getting the user file
	r.ParseMultipartForm(32 << 20)
	// Form with name=imgFile
	file, handler, err := r.FormFile("imgFile")
	if err != nil {
		log.Fatal(err)
	}
	// Closing the file
	defer file.Close()
	// Coping the uploaded file into the server
	f, err := os.OpenFile("./client_upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	io.Copy(f, file)
	// path for the uploaded image
	uploadedImg := "/picSTy/go_web/client_upload/" + handler.Filename
	// getting the name of the picture, without the extension, for the struct field Img in main
	userFileName := strings.Split(handler.Filename, ".")[0]
	return uploadedImg, userFileName
}

func selectArtwork(r *http.Request) string {

	var selectedArt string
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
			selectedArt = "/picSTy/go_web/artworks/" + s + ".jpg"
		}
	}
	return selectedArt
}

func sshPythonContainer(pythonArgs string) {
	var (
		commandSSH string
		outSSH     []byte
	)
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
}

func styledImgTemplate(w http.ResponseWriter, userFileName string) {
	// Execute template/styled.html
	type StyledPage struct {
		// Image name
		Img string
	}
	processedImg := StyledPage{Img: userFileName}
	// Parsing the styled user image name
	styledTemplate, err := template.ParseFiles("template/styled.html")
	if err != nil {
		log.Fatal(err)
	}
	styledTemplate.Execute(w, processedImg)
}

// Uploading user image
func formStyleHandler(w http.ResponseWriter, r *http.Request) {
	var (
		selectedArtwork string
		uploadedImg     string
		pythonArgs      string
		userFileName    string
	)
	// Processing the client picture, getting the path and the name with no extension
	uploadedImg, userFileName = userFileForm(r)
	// Getting the selected picture from the HTML form
	selectedArtwork = selectArtwork(r)
	// this variable has the arguments for the python scripts
	pythonArgs = uploadedImg + " " + selectedArtwork
	// stablishing SSH connection to execute the python script for the neural style transfer
	sshPythonContainer(pythonArgs)
	// Executing the template
	styledImgTemplate(w, userFileName)

}

func main() {

	s := http.Server{
		Addr: ":8080",
		// specifying timeouts, set these to handle malicious or broken HTTP clients
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Landing page, serving the index.html in the go_web directory
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)
	// Handling the form of the landing page
	http.HandleFunc("/styled", formStyleHandler)
	// Serving the web server
	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}

}
