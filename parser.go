package gon3

type Parser struct {
	// target data structure
	Graph Graph
	// parser state
	baseURI      IRI
	namespaces   map[prefix]IRI       // TODO: create prefix type
	bNodeLabels  map[string]BlankNode // TODO: update what the bnode type should look like
	curSubject   RDFTerm              // TODO: create RDFTerm type (or perhaps interface)
	curPredicate RDFTerm
}
