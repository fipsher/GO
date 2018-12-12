package main

import (
	"fmt"
	"log"
	"net/http"
)

func runLocal() {
	reader := FileLinSysReader{"system.txt"}
	//reader := JSONLinSysReader{"[[2,3],[3,5]]", "[8,13]"}
	x := readAndSolve(reader)
	printVector("result", x)
}

func runServer() {
	// initHandlers()
	// initCacher()
	fmt.Println("Starting Cholesky solving server on port 8080...")
	fmt.Println("Visit / for GUI")
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	// http.HandleFunc("/solve", solveHandler)
	// http.HandleFunc("/start", startHandler)
	// http.HandleFunc("/stop/", stopHandler)
	// http.HandleFunc("/result/", resultHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	//runLocal()
	runServer()
}
