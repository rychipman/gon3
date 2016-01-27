package gon3

import (
	"net/url"
)

type RDFTerm interface {
	// iri, blanknode, literal
}

// This must be a full (i.e. not relative IRI)
type IRI struct {
	url *url.URL
}

func newIRIFromString(s string) (IRI, error) {
	url, err := iriRefToURL(s)
	return IRI{url}, err
}

func iriRefToURL(s string) (*url.URL, error) {
	// TODO: implement
	// strip <>, unescape, parse into url
	panic("unimplemented")
}

// see http://www.w3.org/TR/rdf11-concepts/#dfn-blank-node
type BlankNode struct {
	Id    int
	Label string
}

// see http://www.w3.org/TR/rdf11-concepts/#dfn-literal
type Literal struct {
	LexicalForm string
	DatatypeIRI IRI
	LanguageTag string
}

// see http://www.w3.org/TR/rdf11-concepts/#dfn-rdf-triple
type Triple struct {
	Subject   RDFTerm // cannot be a literal
	Predicate IRI
	Object    RDFTerm
}

// An RDF graph is a set of RDF triples
type Graph []*Triple
