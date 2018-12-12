package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// LinearSystemReader - read linear system data from different sources
type LinearSystemReader interface {
	Read() LinearSystem
}

// FileLinSysReader - reader linear system from file
type FileLinSysReader struct {
	FileName string
}

// JSONLinSysReader - read linear system from JSON strings
type JSONLinSysReader struct {
	Matrix string
	Vector string
}

func (r FileLinSysReader) Read() LinearSystem {
	f, _ := os.Open(r.FileName)
	defer f.Close()

	var n int
	fmt.Fscan(f, &n)

	A := make([][]float64, n)
	for i := range A {
		A[i] = make([]float64, n)
		for j := range A[i] {
			fmt.Fscan(f, &A[i][j])
		}
	}

	b := make([]float64, n)
	for i := range b {
		fmt.Fscan(f, &b[i])
	}

	return LinearSystem{A, b}
}

func (r JSONLinSysReader) Read() LinearSystem {
	var A [][]float64
	json.Unmarshal([]byte(r.Matrix), &A)
	var b []float64
	json.Unmarshal([]byte(r.Vector), &b)
	return LinearSystem{A, b}
}
