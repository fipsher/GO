package main

import (
	"fmt"
	"math"
)

// LinearSystem - save matrix A and column vector b
type LinearSystem struct {
	Matrix [][]float64
	Vector []float64
}

func chol(A [][]float64, callback ...func(interface{})) [][]float64 {
	n := len(A)
	L := makeL(n)
	total := float64((n + 1) * n * (n - 1) / 6)
	var count float64
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			sum := 0.0
			for k := 0; k < i; k++ {
				sum += L[i][k] * L[j][k]
				count++
				fOpt(callback)(count / total)
			}
			if i == j {
				L[i][i] = math.Sqrt(A[i][j] - sum)
			} else {
				L[j][i] = (A[i][j] - sum) / L[i][i]
			}
		}
	}
	return L
}

// both for backward and forward substitution
func substitution(L MatrixF, b []float64) []float64 {
	n := len(b)
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		y[i] = b[i]
		for j := 0; j < i; j++ {
			y[i] -= L(i, j) * y[j]
		}
		y[i] /= L(i, i)
	}
	return y
}

// Solve given linear system
func (s LinearSystem) doSolve(callback ...func(interface{})) []float64 {
	n := len(s.Matrix)
	fmt.Println("Solving linear system: ", s.Hash())
	printMatrix("A", matrixToF(s.Matrix, FULL), n)
	printVector("b", s.Vector)
	fmt.Println("Decomposing matrix A...")
	L := chol(s.Matrix, fOpt(callback))
	printMatrix("L", matrixToF(L, LOWER), n)
	printMatrix("U", matrixToF(L, UPPER), n)
	fmt.Println("Performing forward substitution...")
	y := substitution(matrixToF(L, LOWER), s.Vector)
	printVector("y", y)
	fmt.Println("Performing backward substitution...")
	x := reverse(substitution(matrixToF(L, BACKWARD), reverse(y)))
	printVector("x", x)
	return x
}
