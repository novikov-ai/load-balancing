package main

import (
	"fmt"
	"net/http"
	"sync"
)

type backend struct {
	Path    string
	Port    string
	Handler func(rw http.ResponseWriter, r *http.Request)
}

var (
	backends = map[string]backend{
		"Backend #1": {
			Path: "/1",
			Port: ":1111",
			Handler: func(rw http.ResponseWriter, r *http.Request) {
				fmt.Println("Backend #1 - received request:", r.Method, r.URL.Path)
			},
		},
		"Backend #2": {
			Path: "/2",
			Port: ":2222",
			Handler: func(rw http.ResponseWriter, r *http.Request) {
				fmt.Println("Backend #2 - received request:", r.Method, r.URL.Path)
			},
		},
		"Backend #3": {
			Path: "/3",
			Port: ":3333",
			Handler: func(rw http.ResponseWriter, r *http.Request) {
				fmt.Println("Backend #3 - received request:", r.Method, r.URL.Path)
			},
		},
	}
)

func main() {
	wg := sync.WaitGroup{}

	for _, bb := range backends {
		wg.Add(1)

		http.HandleFunc(bb.Path, bb.Handler)

		addr := fmt.Sprintf("127.0.0.1%s", bb.Port)

		fmt.Println("Backend is listening on:", addr)

		go func() {
			http.ListenAndServe(addr, nil)
			wg.Done()
		}()
	}

	wg.Wait()
}
