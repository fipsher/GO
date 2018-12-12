package main

import "fmt"

// MatrixF - Matrix accessor function
type MatrixF func(int, int) float64

// FuncType - Enum of available function types
type FuncType int

// I just want to create this "enum", and not write comments to every line!!!
const (
	FULL      FuncType = iota
	TRANSPOSE          //full base
	LOWER              //lower base
	UPPER              //lower base
	BACKWARD           //lower base
)

func matrixToF(M [][]float64, t FuncType) MatrixF {
	switch t {
	case FULL:
		return func(i, j int) float64 {
			return M[i][j]
		}
	case TRANSPOSE:
		return func(i, j int) float64 {
			return M[j][i]
		}
	case LOWER:
		return func(i, j int) float64 {
			if j < len(M[i]) {
				return M[i][j]
			}
			return 0.0
		}
	case UPPER:
		return func(i, j int) float64 {
			return matrixToF(M, LOWER)(j, i)
		}
	case BACKWARD:
		return func(i, j int) float64 {
			n := len(M)
			return matrixToF(M, LOWER)(n-j-1, n-i-1)
		}
	default:
		panic("Wrong function type") // throw
	}
}

// I don't care about the performance any more...
func reverse(a []float64) []float64 {
	n := len(a)
	b := make([]float64, n)
	for i := 0; i < n; i++ {
		b[i] = a[n-i-1]
	}
	return b
}

func makeL(n int) [][]float64 {
	A := make([][]float64, n)
	for i := range A {
		A[i] = make([]float64, i+1)
	}
	return A
}

func printMatrix(name string, A MatrixF, n int) {
	fmt.Printf("%s = \n", name)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			fmt.Printf("%6.3g", A(i, j))
		}
		fmt.Println()
	}
}
func printVector(name string, b []float64) {
	fmt.Printf("%s = \n", name)
	n := len(b)
	for i := 0; i < n; i++ {
		fmt.Printf("%6.3g", b[i])
	}
	fmt.Println()
}

func readAndSolve(reader LinearSystemReader) []float64 {
	sys := reader.Read()
	x := sys.Solve()
	return x
}

func fOpt(callback []func(interface{}), index ...int) func(interface{}) {
	i := 0
	if len(index) > 0 {
		i = index[0]
	}
	if len(callback) > i {
		return callback[i]
	}
	return func(interface{}) {}
}
