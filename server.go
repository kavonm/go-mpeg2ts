package main

import (
	"log"
	"net/http"
)

func server() {
	http.HandleFunc("/test.ts", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./test.ts")
	})
	http.HandleFunc("/test2.ts", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./test2.ts")
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
