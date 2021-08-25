package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
Hello! This is a local server and so we know it is safe to host
confidential stuff here. In /secrets you can find the SSH keys for 
our jump box. Have Fun!
`)
}

func secrets(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("./id_rsa")
	if err != nil {
		fmt.Fprintf(w, "Unable to read id_rsa file...")
	}
	fmt.Fprintf(w, string(data))
}

func main() {
	serverPort := ":8080"
	http.HandleFunc("/", home)
	http.HandleFunc("/secrets", secrets)
	log.Printf("Server listening on %s", serverPort)
	log.Fatalln(http.ListenAndServe(serverPort, nil))
}
