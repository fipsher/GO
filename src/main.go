package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	initHandlers()
	initDB()
	fs := http.FileServer(http.Dir("client"))


	http.Handle("/", fs)
	http.HandleFunc("/solve", solveHandler)
	http.HandleFunc("/start", startHandler)
	http.HandleFunc("/stop/", stopHandler)
	http.HandleFunc("/result/", resultHandler)

	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
