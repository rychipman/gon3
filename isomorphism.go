package gon3

import (
	"fmt"
)

// algorithm: http://www.hpl.hp.com/techreports/2001/HPL-2001-293.pdf
func fastIsomorphic(g1, g2 Graph) bool {
	// label vertices of g1
	// label vertices of g2
	// if sizes not same, return false
	// classify all vertices of both graphs into same class
	for {
		for { // <A>
			// reclassify each vertex by the number of vertices of each class of the current classification it is adjacent to
			// if new classification is same as old, break and go to <B>
			// if any of new classes has different # of members from each graph, fail & backtrack (break and go to <C>)
			// if any new classes has size 2, break and go to <B>
			// set current clasification as new classification and continue
		}
		// <B>
		// if every class has one element from each graph, then return true
		// <C>
		// choose smallest class with more than two vertices.
		// select arbitrary vertex from v1 in this class
		// guess vertex from v2 in this class
		// when we run out of guesses, backtrack
		// <D>
		// generate new classification by putting the pair of vertices

	}
	// if we backtrack through all guesses, graphs are not isomorphic
	return false
}
