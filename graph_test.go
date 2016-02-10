package gon3

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestIsomorphism(t *testing.T) {
	path := "./tests/turtle/"
	testName := "turtle-subm-22"
	tf, err := ioutil.ReadFile(path + testName + ".ttl")
	if err != nil {
		t.Fatalf("Error reading ttl file %s", testName)
	}
	nf, err := ioutil.ReadFile(path + testName + ".nt")
	if err != nil {
		t.Fatalf("Error reading nt file %s", testName)
	}
	ttlGraph, err := NewParser(string(tf)).Parse()
	if err != nil {
		t.Fatalf("Test %s failed: %s", testName, err)
	}
	ntGraph, err := NewParser(string(nf)).Parse()
	if err != nil {
		t.Fatalf("Test %s failed: %s", testName, err)
	}
	if !ntGraph.IsomorphicTo(ttlGraph) {
		fmt.Printf("ttl graph:\n%s\n", ttlGraph)
		fmt.Printf("nt graph:\n%s\n", ntGraph)
		t.Fatalf("Graphs not isomorphic")
	}
}
