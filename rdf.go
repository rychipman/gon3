package gon3

type RDFTerm interface {
	// iri, blanknode, collection
}

// This must be a full (i.e. not relative IRI)
type IRI string

// see http://www.w3.org/TR/rdf11-concepts/#dfn-blank-node
// A blank node can have an internal iri
type BlankNode int

// see http://www.w3.org/TR/rdf11-concepts/#dfn-literal
type Literal struct {
	LexicalForm string
	DatatypeIRI IRI
	LanguageTag
}

// see http://www.w3.org/TR/rdf11-concepts/#dfn-rdf-triple
type Triple struct {
	Subject   IRI // TODO: allow to be a blank node
	Predicate IRI
	Object    IRI // TODO: allow to be a blank node or a literal
}

// An RDF graph is a set of RDF triples
type Graph []*Triple
