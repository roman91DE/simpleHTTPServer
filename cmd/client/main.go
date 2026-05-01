package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
)

type clargs struct {
	url   string
	names []string
}

func parseArgs() (error, clargs) {
	if len(os.Args) < 3 {
		return (os.ErrInvalid), clargs{}
	}
	url := os.Args[1]
	names := os.Args[2:]
	return nil, clargs{url: url, names: names}
}

func main() {
	err, clargs := parseArgs()
	if err != nil {
		log.Fatalf("Error parsing arguments: %v", err)
	}
	var wg sync.WaitGroup
	for _, name := range clargs.names {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.PostForm(clargs.url, url.Values{"name": {name}})
			if err != nil {
				log.Fatalf("Error in http.PostForm: %v\n", err)
			}
			log.Printf("Response-Code: %v\n", resp.StatusCode)
		}()
	}
	wg.Wait()
}
