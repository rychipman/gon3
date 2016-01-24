package gon3

type RDFTerm interface {
	// iri, blanknode, literal
}

// This must be a full (i.e. not relative IRI)
type IRI string

func newIRI(iri string) (IRI, error) {
	// TODO: actually implement this
	// if first char not '<', process as prefixed name
	// else if relative, resolve according to http://www.w3.org/TR/turtle/#sec-iri-references
	// finally, remove unicode escape sequences
	return IRI(iri), nil
}

// see http://www.w3.org/TR/rdf11-concepts/#dfn-blank-node
type BlankNode struct {
	id    int
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
