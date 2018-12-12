package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
