package gon3

import (
	"fmt"
	"net/url"
)

type RDFTerm interface {
	// iri, blanknode, literal
}

// This must be a full (i.e. not relative IRI)
type IRI struct {
	url *url.URL
}

func (i IRI) String() string {
	return fmt.Sprintf("<%s>", i.url)
}

func newIRIFromString(s string) (IRI, error) {
	url, err := iriRefToURL(s)
	return IRI{url}, err
}

func iriRefToURL(s string) (*url.URL, error) {
	// TODO: implement
	// strip <>, unescape, parse into url
	return url.Parse("")
}

// see http://www.w3.org/TR/rdf11-concepts/#dfn-blank-node
type BlankNode struct {
	Id    int
	Label string
}

func (b BlankNode) String() string {
	label := b.Label
	if label == "" {
		label = "anonbnode"
	}
	return fmt.Sprintf("_:%s", label)
}

// see http://www.w3.org/TR/rdf11-concepts/#dfn-literal
type Literal struct {
	LexicalForm string
	DatatypeIRI IRI
	LanguageTag string
}

func (l Literal) String() string {
	if l.LanguageTag != "" {
		return fmt.Sprintf("%q@%s", l.LexicalForm, l.LanguageTag)
	}
	return fmt.Sprintf("%q^^%s", l.LexicalForm, l.DatatypeIRI)
}

// see http://www.w3.org/TR/rdf11-concepts/#dfn-rdf-triple
type Triple struct {
	Subject   RDFTerm // cannot be a literal
	Predicate IRI
	Object    RDFTerm
}

func (t *Triple) String() string {
	return fmt.Sprintf("%s %s %s .", t.Subject, t.Predicate, t.Object)
}

// An RDF graph is a set of RDF triples
type Graph []*Triple

func (g Graph) String() string {
	str := ""
	for i, t := range g {
		if i > 0 {
			str += "\n"
		}
		str = fmt.Sprintf("%s%s", str, t)
	}
	return str
}
