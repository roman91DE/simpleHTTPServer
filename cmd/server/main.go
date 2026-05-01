package main

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	defer func() {
		if err := w.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()
	writer := bufio.NewWriter(w)
	for k := range mem {
		for _, name := range mem[k] {
			if _, err := writer.WriteString(name); err != nil {
				log.Printf("failed to write string: %v", err)
			}
		}
	}
	if err := writer.Flush(); err != nil {
		log.Printf("failed to flush writer: %v", err)
	}
	log.Printf("Memory written to file: %s", file)
}

func getNames(w http.ResponseWriter, r *http.Request) {
	for k := range mem {
		for _, name := range mem[k] {
			if _, err := w.Write([]byte(name)); err != nil {
				log.Printf("failed to write response: %v", err)
			}
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
	log.Printf("Received name: %s from client: %s", name, clientIP)
	mem[clientIP] = []string{name + "\n"}
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func main() {
	initMem()
	http.HandleFunc("GET /names", getNames)
	http.HandleFunc("POST /names", postName)

	srv := &http.Server{Addr: ":8080"}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Starting Server on Port 8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-quit
	log.Println("Shutting down...")
	writeMem()
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}
