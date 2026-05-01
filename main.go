package main

import (
	// "fmt"
	"bufio"
	"log"
	"net/http"
	"os"
)

var mem map[string][]string

func initMem() {
	bufsize := 1024
	mem = make(map[string][]string, bufsize)
}

func writeMem() {
	file := "mem.txt"
	w, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()
	writer := bufio.NewWriter(w)
	for k := range mem {
		for _, name := range mem[k] {
			writer.WriteString(name)
		}
	}
	writer.Flush()
}

func getNames(w http.ResponseWriter, r *http.Request) {
	for k := range mem {
		for _, name := range mem[k] {
			w.Write([]byte(name))
		}
	}
}
func postName(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	clientIP := r.RemoteAddr
	mem[clientIP] = []string{name + "\n"}
	w.Write([]byte("ok"))
}

func main() {
	initMem()
	defer writeMem()
	http.HandleFunc("GET /names", getNames)
	http.HandleFunc("POST /names", postName)
	http.ListenAndServe(":8080", nil)
}
