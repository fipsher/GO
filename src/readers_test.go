package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func ReaderTester(r LinearSystemReader, t *testing.T) {
	s := r.Read()
	AssertMatrix(matrixToF(s.Matrix, FULL), [][]float64{
		[]float64{1, 2},
		[]float64{3, 4}}, 2, t)
	AssertEqual(len(s.Vector), 2, t)
	AssertEqual(s.Vector[0], 5.0, t)
	AssertEqual(s.Vector[1], 6.0, t)
}

func TestJSON(t *testing.T) {
	r := JSONLinSysReader{"[[1,2],[3,4]]", "[5,6]"}
	ReaderTester(r, t)
}

func TestFile(t *testing.T) {
	content := []byte("2\n1 2\n3 4\n5\n6")
	tmpfile, _ := ioutil.TempFile("", "system")
	defer os.Remove(tmpfile.Name())
	tmpfile.Write(content)
	tmpfile.Close()
	r := FileLinSysReader{tmpfile.Name()}
	ReaderTester(r, t)
}

type TestReader LinearSystem

func (r TestReader) Read() LinearSystem {
	return LinearSystem(r)
}
func TestReadAndSolve(t *testing.T) {
	r := TestReader{[][]float64{
		[]float64{2, 1},
		[]float64{1, 3}},
		[]float64{3, 4}}
	x := readAndSolve(r)
	AssertEqual(len(x), 2, t)
	AssertEqual(x[0], 1.0, t)
	AssertEqual(x[1], 1.0, t)
}
