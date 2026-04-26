package main

import (
	// "fmt"
	"net/http"
	// "log"
)

var mem map[string]bool

func initMEM() {
	bufsize := 1024
	mem = make(map[string]bool, bufsize)
}

func getNames(w http.ResponseWriter, r *http.Request) {
	for k := range mem {
		w.Write([]byte(k + "\n"))
	}	
}
func postName(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	mem[name] = true
	w.Write([]byte("ok"))
}

func main() {
	initMEM()
	http.HandleFunc("GET /names", getNames)
	http.HandleFunc("POST /names", postName)
	http.ListenAndServe(":8080", nil)
}
