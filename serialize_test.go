package gon3

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestGoodSerialize(t *testing.T) {

	currentSerializeTests = []string{}
	currentSerializeTests = goodSerializeTests

	verbosity := 1

	for _, testName := range currentSerializeTests {
		testFile := "./tests-out/" + testName + ".ttl"
		outFile := "./tests-out/" + testName + ".out"
		f, err := ioutil.ReadFile(testFile)
		if err != nil {
			t.Fatalf("Error reading test file %s", testFile)
		}
		o, err := ioutil.ReadFile(outFile)
		if verbosity > 100 {
			fmt.Printf("%s", o)
		}
		if err != nil {
			t.Fatalf("Error reading out file %s", outFile)
		}
		if verbosity > 0 {
			fmt.Printf("\nStarting test %s\n", testName)
		}
		p := NewParser(string(f))
		g, err := p.Parse()
		if err != nil {
			t.Fatalf("Test %s failed: %s", testName, err)
		}
		if verbosity > 0 {
			fmt.Printf("Test %s passed.\n", testName)
		}
		if verbosity > 1 {
			fmt.Printf("Graph:\n%s\n", g)
		}
	}
}

// manifest-bad.ttl
// manifest.ttl
// README.txt

var currentSerializeTests []string
var badSerializeTests []string = []string{
	"bad-00",
	"bad-01",
	"bad-02",
	"bad-03",
	"bad-04",
	"bad-05",
	"bad-06",
	"bad-07",
	"bad-08",
	"bad-09",
	"bad-10",
	"bad-11",
	"bad-12",
	"bad-13",
	"bad-14",
}
var goodSerializeTests []string = []string{
	"rdfq-results",
	"rdf-schema",
	"rdfs-namespace",
	"test-00",
	"test-01",
	"test-02",
	"test-03",
	"test-04",
	"test-05",
	"test-06",
	"test-07",
	"test-08",
	"test-09",
	"test-10",
	"test-11",
	"test-12",
	"test-13",
	"test-14",
	"test-15",
	"test-16",
	"test-17",
	"test-18",
	"test-19",
	"test-20",
	"test-21",
	"test-22",
	"test-23",
	"test-24",
	"test-25",
	"test-26",
	"test-27",
	"test-28-out",
	"test-29",
	"test-30",
	"test-30.ttl",
}
