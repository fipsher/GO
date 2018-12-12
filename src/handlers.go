package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// CalcInfo - status/result of background calculations
type CalcInfo struct {
	status map[int]int
	result map[int][]float64
	cached map[int]string
	mux    sync.Mutex
}

// ReqOutput - calculation info
type ReqOutput struct {
	ID        int
	Status    string
	Progress  int
	Timestamp string
	Result    []float64
}

var info CalcInfo

func initHandlers() {
	rand.Seed(time.Now().UnixNano())
	info = CalcInfo{
		status: make(map[int]int),
		result: make(map[int][]float64),
		cached: make(map[int]string)}
}

func requestReader(r *http.Request) JSONLinSysReader {
	matrix := r.PostFormValue("matrix")
	vector := r.PostFormValue("vector")
	return JSONLinSysReader{matrix, vector}
}

func solveHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/solve" || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	fmt.Println("Processing solve request...")

	reader := requestReader(r)
	x := readAndSolve(reader)

	output, _ := json.Marshal(x)
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func recovery() {
	if r := recover(); r != nil {
		fmt.Println("recovered:", r)
	}
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/start" || r.Method != http.MethodPost {
		http.NotFound(w, r)
	}
	fmt.Println("Processing start request...")

	delay, _ := strconv.Atoi(r.PostFormValue("delay"))
	reader := requestReader(r)
	sys := reader.Read()
	id := rand.Int()
	info.status[id] = 0
	go func(id int) {
		defer recovery()
		step := func(progress interface{}) {
			time.Sleep(time.Duration(delay) * time.Millisecond)
			percent := int(progress.(float64) * 100)
			info.mux.Lock()
			defer info.mux.Unlock()
			if info.status[id] < 0 {
				panic("Calculation was stopped\n")
			}
			info.status[id] = percent
			fmt.Printf("Progress of %d changed to %d%%\n", id, percent)
		}
		cached := func(ts interface{}) {
			tsStr := ts.(string)
			info.mux.Lock()
			defer info.mux.Unlock()
			info.cached[id] = tsStr
			fmt.Printf("Solution of %d retrieved. Cached at %s\n", id, tsStr)
		}
		x := sys.Solve(step, cached)
		info.mux.Lock()
		defer info.mux.Unlock()
		info.status[id] = 100
		info.result[id] = x
	}(id)

	fmt.Fprint(w, id)
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
	}
	fmt.Println("Processing stop request...")

	info.mux.Lock()
	defer info.mux.Unlock()

	idStr := r.URL.Path[len("/stop/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		if _, ok := info.status[id]; ok {
			info.status[id] = -1
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
	}
	// too often to print every time
	//fmt.Println("Processing result request...")

	var res ReqOutput

	info.mux.Lock()
	defer info.mux.Unlock()

	idStr := r.URL.Path[len("/result/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res.Status = "incorrect id"
	} else {
		res.ID = id
		if status, ok := info.status[id]; ok {
			if status < 0 {
				res.Status = "stopped"
			} else {
				res.Status = "calculating"
				res.Progress = status
			}
			if result, ok := info.result[id]; ok {
				res.Status = "done"
				res.Progress = 100
				res.Result = result
			}
			if cached, ok := info.cached[id]; ok {
				res.Status = "cached"
				res.Timestamp = cached
			}
		} else {
			res.Status = "task not found"
		}
	}

	output, _ := json.Marshal(&res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}
