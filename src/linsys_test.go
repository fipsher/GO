package main

import "testing"

func TestLU(t *testing.T) {
	t.Run("N=1", func(t *testing.T) {
		M := [][]float64{[]float64{4}}
		L := chol(M)
		AssertEqual(len(L), 1, t)
		AssertEqual(len(L[0]), 1, t)
		AssertEqual(L[0][0], 2.0, t)
	})
	t.Run("N=3", func(t *testing.T) {
		M := [][]float64{
			[]float64{1, 2, 3},
			[]float64{2, 5, 7},
			[]float64{3, 7, 11}}
		L := chol(M)
		AssertMatrix(matrixToF(L, LOWER), [][]float64{
			[]float64{1, 0, 0},
			[]float64{2, 1, 0},
			[]float64{3, 1, 1}}, 3, t)
	})
}

func TestSubstitution(t *testing.T) {
	t.Run("N=1", func(t *testing.T) {
		M := [][]float64{[]float64{2}}
		v := []float64{4}
		x := substitution(matrixToF(M, LOWER), v)
		AssertEqual(len(x), 1, t)
		AssertEqual(x[0], 2.0, t)
	})
	t.Run("N=3", func(t *testing.T) {
		M := [][]float64{
			[]float64{1},
			[]float64{2, 1},
			[]float64{3, 1, 1}}
		v := []float64{1, 4, 8}
		x := substitution(matrixToF(M, LOWER), v)
		AssertEqual(len(x), 3, t)
		AssertEqual(x[0], 1.0, t)
		AssertEqual(x[1], 2.0, t)
		AssertEqual(x[2], 3.0, t)
	})
}

func TestSolution(t *testing.T) {
	t.Run("N=1", func(t *testing.T) {
		M := [][]float64{[]float64{2}}
		b := []float64{4}
		s := LinearSystem{M, b}
		x := s.doSolve()
		AssertEqual(len(x), 1, t)
		AssertEqual(x[0], 2.0, t)
	})
	t.Run("N=3", func(t *testing.T) {
		M := [][]float64{
			[]float64{1, 2, 3},
			[]float64{2, 5, 7},
			[]float64{3, 7, 11}}
		b := []float64{34, 80, 121}
		s := LinearSystem{M, b}
		x := s.doSolve()
		AssertEqual(len(x), 3, t)
		AssertEqual(x[0], 3.0, t)
		AssertEqual(x[1], 5.0, t)
		AssertEqual(x[2], 7.0, t)
	})
	t.Run("Callback", func(t *testing.T) {
		called := false
		f := func(interface{}) {
			called = true
		}
		M := [][]float64{
			[]float64{2, 1},
			[]float64{1, 3}}
		b := []float64{3, 4}
		s := LinearSystem{M, b}
		_ = s.doSolve(f)
		AssertEqual(called, true, t)
	})
}
