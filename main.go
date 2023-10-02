package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Declare global variables
var tpl *template.Template
var word string
var str string

// Define a struct for JSON data
type Colors struct {
	Colors []string `json:"colors"`
}

// Function for generating ASCII art
func ascii_art(argument string, fonts string) (string, int) {
	// Read the ASCII art font file
	banner, err := ioutil.ReadFile("fonts/" + fonts + ".txt")
	if err != nil {
		return "Error 500\nInternal Error", 500
	}
	// Split the ASCII art into lines
	split := strings.Split(string(banner), "\n")
	if fonts == "thinkertoy" {
		split = strings.Split(string(banner), "\r\n")
	}

	// Split the input text into lines

	myting := strings.Split(strings.ReplaceAll(argument, "\r", ""), "\\n")

	// Loop through each line and generate ASCII art

	for word := 0; word < len(myting); word++ {
		if word == 0 && len(myting) >= 3 {
			if len(myting[0]) == 0 && len(myting[1]) == 0 && len(myting[2]) == 0 {
				word += 1
			}
		}
		for k := 0; k < 8; k++ {
			if len(myting[word]) == 0 && len(myting) >= 2 {
				k = 7
			}
			for i := 0; i < len(myting[word]); i++ {

				str += split[(int(myting[word][i])-32)*9+1+k]

			}

			if len(myting[word]) != 0 {
				str += "\n"
			}

			if len(myting[word]) == 0 && len(myting) >= 2 {
				str += "\n"
				// This would check for a new line which in this case is a backslash n (" \n")
				if len(myting) == 2 && word != len(myting)-1 {
					if len(myting[word+1]) == 0 {
						word++
					}
				}

			}
		}
	}

	// Write the generated ASCII art to files
	
	err = os.WriteFile("download.doc", []byte(str), 0666)
	if err != nil {
		panic(err)
	}
	err1 := os.WriteFile("download.txt", []byte(str), 0666)
	if err1 != nil {
		panic(err1)
	}
	return str, 200
}

// Function for handling file download

func download(w http.ResponseWriter, r *http.Request) {

	formatType := r.FormValue("fileformat")

	f, _ := os.Open("download." + formatType)
	defer f.Close()

	file, _ := f.Stat()
	fsize := file.Size()

	sfSize := strconv.Itoa(int(fsize))
	
	// Set HTTP headers for file download

	w.Header().Set("Content-Disposition", "attachment; filename=asciiresults."+formatType)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", sfSize)

	io.Copy(w, f)
}

// Initialize the HTML template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

// Function for rendering input text

func render(s string) (string, int) {
	a := []rune(s)
	noerr, _ := errorcheck(s)
	if len(s) >= 128 {
		return "Too long", 400
	}
	if s == "" {
		return "Enter a text!", 200
	}
	if noerr {
		for i, _ := range s {
			if a[i] == 13 && a[i+1] == 10 {
				a[i] = 92
				a[i+1] = 110
			}
		}
		return string(a), 200
	} else {
		return "Bad request", 400
	}
}

// Function for checking errors in input text

func errorcheck(s string) (bool, int) {
	a := []rune(s)
	for i, _ := range s {
		if a[i] <= 127 {
			continue
		} else {
			return false, 400
		}
	}
	return true, 200
}

// Function for processing input data

func processor(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		text := r.FormValue("ascii-data")
		fonts := r.FormValue("fonts")
		final, renderstatus := render(text)
		str = ""
		ColorsJason := "./color.json"
		TextColor := r.FormValue("color") // Get the selected color from the form
		bgcolor := "#a78295"

		file, err := os.Open(ColorsJason)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer file.Close()
		ColorsArry, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		colors := Colors{}

		err = json.Unmarshal(ColorsArry, &colors)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		found := false
		for _, color := range colors.Colors {
			if color == TextColor {
				found = true
				break
			}
		}

		if found == true {
			bgcolor = "#ffffff"
		}

		if fonts == "standard" || fonts == "shadow" || fonts == "thinkertoy" {
			if renderstatus == 400 {
				fmt.Printf("%s did a bad request (400)\nWith the text: %s\n", r.RemoteAddr, text)
				data, _ := ascii_art(final, fonts)

				d := struct {
					First           string
					SelectedColor   string
					Backgroundcolor string
				}{
					First:           data,
					SelectedColor:   TextColor,
					Backgroundcolor: bgcolor,
				}
				w.WriteHeader(http.StatusBadRequest)
				tpl.ExecuteTemplate(w, "index.html", d)
				return
			}

			data, statuscode := ascii_art(final, fonts)
			if statuscode == 500 {
				fmt.Printf("%s got internal error (500) from %s\n", r.RemoteAddr, text)
				data, _ := ascii_art(final, fonts)
				d := struct {
					First           string
					SelectedColor   string
					Backgroundcolor string
				}{
					First:           data,
					SelectedColor:   TextColor, // Pass the selected color to the template
					Backgroundcolor: bgcolor,
				}
				w.WriteHeader(http.StatusInternalServerError)
				tpl.ExecuteTemplate(w, "index.html", d)
				return
			}
			if statuscode == 200 {
				fmt.Printf("%s/n sent the text: %s\nWith the font: %s\n %s\n", r.RemoteAddr, text, fonts, TextColor)
				d := struct {
					First           string
					SelectedColor   string
					Backgroundcolor string
				}{
					First:           data,
					SelectedColor:   TextColor, // Pass the selected color to the template
					Backgroundcolor: bgcolor,
				}
				tpl.ExecuteTemplate(w, "index.html", d)
				return
			}

		} else {
			fmt.Printf("New connection from %s\n", r.RemoteAddr)
			d := struct {
				SelectedColor string
			}{
				SelectedColor: TextColor, // Pass the selected color to the template
			}
			tpl.ExecuteTemplate(w, "index.html", d)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	tpl.ExecuteTemplate(w, "500.html", nil)
	return
}

// Main function

func main() {
	// Define HTTP routes and handlers
	http.HandleFunc("/", index)
	http.HandleFunc("/ascii-art", processor)
	http.HandleFunc("/right", download)
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))
	
	// Start the HTTP server

	fmt.Println("HTTP SERVER RUNNING AT: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Function for handling the root path

func index(w http.ResponseWriter, r *http.Request) {
	userAgent := r.Header.Get("User-Agent")
	if strings.Contains(userAgent, "curl") {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	if r.URL.Path != "/ascii-art" && r.URL.Path != "/" {
		fmt.Printf("%s got a 404 with the path: %s\n", r.RemoteAddr, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		tpl.ExecuteTemplate(w, "404.html", nil)
		return
	}

	if r.Method == "POST" {
		// Handle POST request to the root path ("/")
		fmt.Printf("%s made a POST request to the root path\n", r.RemoteAddr)
		w.WriteHeader(http.StatusBadRequest)
		tpl.ExecuteTemplate(w, "400.html", nil)
		return
	}

	// Handle GET request to the root path ("/")
	fmt.Printf("New connection from %s\n", r.RemoteAddr)
	tpl.ExecuteTemplate(w, "index.html", nil)
}
