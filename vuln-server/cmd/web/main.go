package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func httpGet(w http.ResponseWriter, url string) {
	// checking for GCP Instance Metadata GET requests
	if strings.Contains(url, "169.254.169.254") ||
		strings.Contains(url, "a9.fe.a9.fe") ||
		strings.Contains(url, "0251.0376.0251.0376") ||
		strings.Contains(url, "10101001.11111110.10101001.11111110") {
		fmt.Fprintf(w, "<h2>GCP Instance Metadata is off limits.. nice try though</h2>\n")
	} else {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(w, "<h3>Unable to make GET request to url: %s</h3>", url)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(w, "<h3>Unable to read HTTP response data</h3>")
		}

		templateData := struct {
			Response string
		}{Response: template.HTMLEscapeString(string(body))}

		t, err := template.ParseFiles("static/curl.html.tmpl")
		if err != nil {
			log.Fatalln("Unable to parse Go template file")
		}
		t.Execute(w, templateData)
	}
}

//func home(w http.ResponseWriter, r *http.Request) {
//switch r.Method {
//case "GET":
//http.ServeFile(w, r, "html/home.html")
//case "POST":
//fmt.Fprintf(w, "POST requests are not supported")
//default:
//fmt.Fprintf(w, "Only GET requests are supported")
//}
//}

func curl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "static/curl.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "<h2>Unable to parse submited data</h2>")
			return
		}

		url := r.FormValue("url")
		if url != "" {
			httpGet(w, url)
		} else {
			http.Redirect(w, r, "/curl", 301)
		}
	default:
		fmt.Fprintf(w, "<h2>Only GET and POST requests are supported</h2>")
	}
}

func main() {
	serverPort := ":3000"

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Unable to get current working directory")
	}

	http.Handle("/", http.FileServer(http.Dir(fmt.Sprintf("%s/static", cwd))))
	http.HandleFunc("/curl", curl)

	log.Printf("Server listening on %s", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}
