package gon3

import (
	"crypto/sha1"
	"fmt"
	"sort"
)

// algorithm: http://www.hpl.hp.com/techreports/2001/HPL-2001-293.pdf
// also: http://blog.datagraph.org/2010/03/rdf-isomorphism
func Isomorphic(g1, g2 Graph) bool {
	grounded1, ungrounded1 := hashNodes(g1, bNodesIn(g1), map[Term]string{})
	grounded2, ungrounded2 := hashNodes(g2, bNodesIn(g2), map[Term]string{})

	// check that grounded hashes are the same up to this point
	for term, _ := range grounded1 {
		_, present := grounded2[term]
		if !present {
			return false
		}
	}
	for term, _ := range grounded2 {
		_, present := grounded1[term]
		if !present {
			return false
		}
	}

	// map nodes in g1 to nodes in g2
	bijection := map[Term]Term{}
	for node := range g1.IterNodes() {
		for other, hash := range ungrounded2 {
			nodeHash := ungrounded1[node.(BlankNode)]
			if nodeHash == hash {
				bijection[node] = other
				delete(ungrounded2, other)
				break
			}
		}
	}

	// if all nodes accounted for on both sides, we have a bijection

	panic("unimplemented")
}

func hashNodes(g Graph, nodes []BlankNode, hashes map[Term]string) (grounded map[Term]string, ungrounded map[BlankNode]string) {
	grounded = map[Term]string{}
	for k, v := range hashes {
		grounded[k] = v
	}
	ungrounded = map[BlankNode]string{}
	for {
		isUniqueHash := map[string]bool{}
		numGrounded := len(grounded)
		// mark as many nodes as possible as grounded
		for _, node := range nodes {
			if _, present := grounded[node]; !present {
				isGrounded, hash := nodeHash(node, g, grounded)
				if isGrounded {
					grounded[node] = hash
				}
				ungrounded[node] = hash
				if _, present := isUniqueHash[hash]; !present {
					isUniqueHash[hash] = true
				} else {
					isUniqueHash[hash] = false
				}
			}
		}
		// mark any node with a unique hash as grounded
		for node, hash := range ungrounded {
			if isUniqueHash[hash] {
				grounded[node] = hash
			}
		}
		// break if we haven't added any grounded nodes
		if len(grounded) == numGrounded {
			break
		}
	}
	return
}

func nodeHash(node BlankNode, g Graph, hashes map[Term]string) (bool, string) {
	tripleSignatures := []string{}
	grounded := true
	for trip := range g.IterTriples() {
		if trip.includes(node) {
			tripleSignatures = append(tripleSignatures, hashString(trip, hashes, node))
			for _, term := range []Term{trip.Subject, trip.Predicate, trip.Object} {
				_, present := hashes[term]
				if !term.Equals(node) && !present {
					grounded = false
				}
			}
		}
	}
	trips := sort.StringSlice(tripleSignatures)
	trips.Sort()
	hash := sha1.Sum([]byte(fmt.Sprintf("%v", trips)))
	return grounded, string(hash[:])
}

func hashString(trip *Triple, hashes map[Term]string, node BlankNode) string {
	str := ""
	for _, term := range []Term{trip.Subject, trip.Predicate, trip.Object} {
		hash, grounded := hashes[term]
		switch {
		case node.Equals(term):
			str += "itself"
		case grounded:
			str += hash
		case isBlankNode(node):
			str += "a blank node"
		default:
			str += node.String()
		}
	}
	return str
}

func bNodesIn(g Graph) []BlankNode {
	bNodes := []BlankNode{}
	for trip := range g.IterTriples() {
		for _, term := range []Term{trip.Subject, trip.Predicate, trip.Object} {
			if isBlankNode(term) {
				bNodes = append(bNodes, term.(BlankNode))
			}
		}
	}
	return bNodes
}
