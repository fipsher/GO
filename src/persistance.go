package main

import (
	// "database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"math"

	redis "github.com/go-redis/redis"
)

// var db *sql.DB

var db * redis.Client

func initDB() {
	db = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := db.Ping().Result()
	if err != nil {
		db = nil
	}

	fmt.Println(pong, err)
   }

// Hash - calculate hash value of the entire system
func (s LinearSystem) Hash() string {
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
	inthash := int(hash ^ (hash >> 32))

	return strconv.Itoa(inthash)
}

// Solve with cache lookup
func (s LinearSystem) Solve(callback ...func(interface{})) []float64 {

	hash := s.Hash()
	if db != nil {
		val, err := db.Get(hash).Result()
		
		// rows.Scan(&xStr, &ts)
		// err := rows.Scan(&xStr, &ts)
		if err != redis.Nil {
			fmt.Println("Solution available in cache!")
			var x []float64
			json.Unmarshal([]byte(val), &x)
			fOpt(callback, 1)
			return x
		}
	}
	x := s.doSolve(fOpt(callback))
	if db != nil {
		fmt.Println("Caching result...")
		xStr, _ := json.Marshal(x)
		db.Set(hash, string(xStr), 0)
	}
	return x
}