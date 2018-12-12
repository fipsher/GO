package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"

	_ "github.com/denisenkom/go-mssqldb"
)

var db *sql.DB

func initCacher() {
	connStr := "server=DESKTOP-1OTLGRK;user id=me;password=aaaAAA123.;database=go"
	var err error
	db, err = sql.Open("sqlserver", connStr)

	if err != nil {
		fmt.Println("  Error open db:", err.Error())
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err)
	}

	if err != nil {
		// not fatal - will not use cache
		log.Println(err)
	}
}

// Hash - calculate hash value of the entire system
func (s LinearSystem) Hash() int32 {
	var hash uint64
	for _, row := range s.Matrix {
		for _, elem := range row {
			hash <<= 1
			hash ^= math.Float64bits(elem)
		}
	}
	for _, elem := range s.Vector {
		hash <<= 1
		hash ^= math.Float64bits(elem)
	}
	return int32(hash ^ (hash >> 32))
}

// Solve with cache lookup
func (s LinearSystem) Solve(callback ...func(interface{})) []float64 {
	hash := s.Hash()
	if db != nil {
		rows := db.QueryRow("SELECT solution, cached FROM solutions WHERE hash = $1", hash)
		var xStr, ts string
		err := rows.Scan(&xStr, &ts)
		if err != sql.ErrNoRows {
			fmt.Println("Solution available in cache!")
			fmt.Println(xStr)
			var x []float64
			json.Unmarshal([]byte(xStr), &x)
			fOpt(callback, 1)(ts)
			fmt.Println(x)
			return x
		}
	}
	x := s.doSolve(fOpt(callback))
	if db != nil {
		fmt.Println("Caching result...")
		xStr, _ := json.Marshal(x)
		db.Exec("INSERT INTO solutions VALUES ($1, $2)",
			hash, string(xStr))
	}
	return x
}

// CREATE TABLE public.solutions
// (
//    hash integer NOT NULL,
//    solution text NOT NULL,
//    cached timestamp DEFAULT CURRENT_TIMESTAMP,
//    CONSTRAINT "solutionsPK" PRIMARY KEY (hash)
// )
