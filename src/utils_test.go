package main

import (
	"math"
	"math/rand"
	"testing"
)

func AssertEqual(a interface{}, b interface{}, t *testing.T) {
	switch v := a.(type) {
	case float64:
		eps := 1e-5
		u, ok := b.(float64)
		if ok {
			if math.Abs(u-v) > eps {
				t.Fail()
			}
		} else {
			t.Fail()
		}
	default:
		if a != b {
			t.Fail()
		}
	}
}
func AssertNotEqual(a interface{}, b interface{}, t *testing.T) {
	if a == b {
		t.Fail()
	}
}
func AssertMatrix(A MatrixF, B [][]float64, n int, t *testing.T) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			AssertEqual(A(i, j), B[i][j], t)
		}
	}
}

func TestMatrixF(t *testing.T) {
	M := [][]float64{
		[]float64{1, 2, 3},
		[]float64{2, 5, 7},
		[]float64{3, 7, 11}}
	L := [][]float64{
		[]float64{1},
		[]float64{2, 5},
		[]float64{3, 7, 11}}
	t.Run("FULL", func(t *testing.T) {
		F := matrixToF(M, FULL)
		AssertEqual(F(0, 0), M[0][0], t)
		AssertEqual(F(0, 1), M[0][1], t)
		AssertEqual(F(2, 1), M[2][1], t)
	})
	t.Run("TRANSPOSE", func(t *testing.T) {
		F := matrixToF(M, TRANSPOSE)
		AssertEqual(F(0, 0), M[0][0], t)
		AssertEqual(F(0, 1), M[1][0], t)
		AssertEqual(F(2, 1), M[1][2], t)
	})
	t.Run("LOWER", func(t *testing.T) {
		F := matrixToF(L, LOWER)
		AssertEqual(F(0, 0), L[0][0], t)
		AssertEqual(F(0, 1), 0.0, t)
		AssertEqual(F(2, 1), L[2][1], t)
	})
	t.Run("UPPER", func(t *testing.T) {
		F := matrixToF(L, UPPER)
		AssertEqual(F(0, 0), L[0][0], t)
		AssertEqual(F(0, 1), L[1][0], t)
		AssertEqual(F(2, 1), 0.0, t)
	})
	t.Run("BACKWARD", func(t *testing.T) {
		F := matrixToF(L, BACKWARD)
		AssertEqual(F(0, 0), L[2][2], t)
		AssertEqual(F(0, 1), 0.0, t)
		AssertEqual(F(2, 1), L[1][0], t)
	})
}

func TestReverse(t *testing.T) {
	a := []float64{1, 2, 3}
	b := reverse(a)
	AssertEqual(len(b), 3, t)
	AssertEqual(b[0], a[2], t)
	AssertEqual(b[1], a[1], t)
	AssertEqual(b[2], a[0], t)
}

func TestMakeL(t *testing.T) {
	L := makeL(3)
	AssertEqual(len(L), 3, t)
	AssertEqual(len(L[0]), 1, t)
	AssertEqual(len(L[1]), 2, t)
	AssertEqual(len(L[2]), 3, t)
}


func TestFOpt(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		called := false
		f := func(interface{}) {
			called = true
		}
		fOpt([]func(interface{}){f})(0)
		AssertEqual(called, true, t)
	})
	t.Run("Second", func(t *testing.T) {
		called := false
		notCalled := true
		f := func(interface{}) {
			called = true
		}
		g := func(interface{}) {
			notCalled = false
		}
		fOpt([]func(interface{}){g, f}, 1)(0)
		AssertEqual(called, true, t)
		AssertEqual(notCalled, true, t)
	})
	t.Run("None", func(t *testing.T) {
		fOpt([]func(interface{}){})(0)
	})
}

func TestInitHandlers(t *testing.T) {
	initHandlers()
	AssertNotEqual(rand.Int(), 5577006791947779410, t)
	AssertNotEqual(&info, nil, t)
}
