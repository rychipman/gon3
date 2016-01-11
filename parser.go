package gon3

import (
	"fmt"
)

type Parser struct {
	// target data structure
	Graph *Graph
	// parser state
	lex           *lexer // TODO: initialize lexer
	baseURI       IRI
	namespaces    map[prefix]IRI // TODO: create prefix type
	bNodeLabels   map[string]BlankNode
	lastBlankNode BlankNode // TODO: initialize this to -1
	curSubject    RDFTerm   // TODO: create RDFTerm type (or perhaps interface)
	curPredicate  RDFTerm
}

func (p *Parser) emitTriple(subj, pred, obj RDFTerm) { // TODO: work out typing things
	trip := Triple{
		Subject:   subj,
		Predicate: pred,
		Object:    obj,
	}
	append(p.Graph, trip)
}

func (p *Parser) absIRI(iri string) (IRI, error) {
	// if first char not '<', process as prefixed name
	// else if relative, resolve according to http://www.w3.org/TR/turtle/#sec-iri-references
	// finally, remove unicode escape sequences
}

func (p *Parser) blankNode(label string) (BlankNode, error) {
	// TODO: when would we return an error?
	if node, present := p.bNodeLabels[label]; present {
		return node, nil
	}
	newNode := p.lastBlankNode + 1 // TODO: def not correct
	p.bNodeLabels[label] = newNode
	return newNode, nil
}

func (p *Parser) Parse(text string) (*Graph, error) {
	// initialize fields
	p.Graph = Graph{}
	p.namespaces = map[prefix]IRI{}
	p.bNodeLabels = map[string]BlankNode{}

	err := nil
	for { // while the next token is not an EOF
		err = p.parseStatement()
		if err != nil {
			break
		}
	}
	return p.Graph, err
}

func (p *Parser) parseStatement() error {
	// TODO: check if triples or directive
	tok := p.peek()
	switch tok.typ {
	case tokenAtPrefix:
		err := p.parsePrefix()
		if err != nil {
			return err
		}
	case tokenAtBase:
		err := p.parseBase()
		// TODO: support tokenSPARQLBase, token SPARQLPrefix
		if err != nil {
			return err
		}
	default:
		err := p.parseTriples()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parsePrefix() error {
	tok := p.next()
	if tok.typ != tokenAtPrefix {
		return fmt.Errorf("Expected tokenAtPrefix, got %v", tok)
	}
	// expect PNAME_NS token
	pNameNS := p.next()
	if pNameNS.typ != tokenPNameNS {
		return fmt.Errorf("Expected tokenPNameNS, got %v", pNameNS)
	}
	// expect IRIREF token
	iriRef := p.next()
	if iriRef.typ != tokenIRIRef {
		return fmt.Errorf("Expected tokenIRIRef, got %v", iriRef)
	}
	// map a new namespace in parser state
	key := pNameNS.val[:len(pNameNS.val)-1]
	val, err := p.absIRI(iriRef.val)
	p.namespaces[key] = val
	return err
}

func (p *Parser) parseBase() error {
	// expect '@base' token
	tok := p.next()
	if tok.typ != tokenAtBase {
		return fmt.Errorf("Expected tokenAtBase, got %v", tok)
	}
	// expect IRIREF token
	iriRef := p.next()
	if iriRef.typ != tokenIRIRef {
		return fmt.Errorf("Expected tokenIRIRef, got %v", iriRef)
	}
	// TODO: require iriRef to be an absolute (or maybe prefixed?) iri
	// for now, assume it is an abs iri
	p.baseURI, err = p.absIRI(iriRef.val)
	return err
}

func (p *Parser) parseTriples() error {
	if true { // if "subject predicateobjectlist"
		err := p.parseSubject()
		if err != nil {
			return err
		}
		err = parsePredicateObjectList()
		if err != nil {
			return err
		}
	} else { // if "blanknodepropertylist predicateobjectlist?"
		bNode, err := p.parseBlankNodePropertyList()
		if err != nil {
			return err
		}
		p.curSubject = bNode
		// parse a predicateobjectlist if we have one
		if p.peek().typ != tokenEndTriples {
			err = p.parsePredicateObjectList()
			if err != nil {
				return err
			}
		}
	}
	// expect period token
	tok := p.next()
	if tok.typ != tokenEndTriples {
		return fmt.Errorf("Expected tokenEndTriples, got %v", tok)
	}
	return nil
}

func (p *Parser) parseSubject() error {
	// TODO: implement
	tok := p.peek()
	// expect a valid subject term, which is one of
	// iri|blanknode|collection
	switch tok.typ {
	case tokenIRIRef: // TODO: include PrefixedName here
		p.next()
		iri, err := p.absIRI(tok.val)
		p.curSubject = iri
		return err
	case tokenBlankNodeLabel:
		// TODO: what is the deal with the anon token?
		p.next()
		label := tok.val // TODO: parse the label out of token value
		bNode, err := p.blankNode(label)
		p.curSubject = bNode
		return err
	case tokenStartCollection:
		bNode, err := p.parseCollection()
		p.curSubject = bNode
		return err
	default:
		return fmt.Errorf("Expected a subject, got %v", tok)
	}
}

func (p *Parser) parsePredicateObjectList() error {
	// http://www.w3.org/TR/turtle/#predicate-lists
	// TODO: figure out details of when semicolon is required
	// currently, this assumes that there will not be a semicolon
	// unless the list is continued
	err := p.parsePredicate()
	if err != nil {
		return err
	}
	err = p.parseObjectList()
	if err != nil {
		return err
	}
	for p.peek().typ != tokenEndTriples {
		// expect semicolon token
		tok := p.next()
		if tok.typ != tokenPredicateListSeparator {
			return fmt.Errorf("Expected tokenPredicateListSeparator, got %v", tok)
		}
		err := p.parsePredicate()
		if err != nil {
			return err
		}
		err = p.parseObjectList()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parsePredicate() error {
	// expect token 'a' or an iri
	tok := p.next()
	switch tok.typ {
	case tokenA:
		// TODO: remove magic string
		pred, err := p.absIRI("<http://www.w3.org/1999/02/22-rdf-syntax-ns#type>")
		p.curPredicate = pred
		return err
	case tokenIRIRef: // TODO: include PrefixedName here
		iri, err := p.absIRI(tok.val)
		p.curPredicate = iri
		return err
	default:
		return fmt.Errorf("Expected predicate, got %v", tok)
	}
}

func (p *Parser) parseObjectList() error {
	err := p.parseObject()
	if err != nil {
		return err
	}
	for p.peek().typ != tokenEndTriples {
		// expect comma token
		tok := p.next()
		if tok.typ != tokenObjectListSeparator {
			return fmt.Errorf("Expected tokenObjectListSeparator, got %v", tok)
		}
		err := p.parseObject()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseCollection() (BlankNode, error) {
	// TODO: implement
	savedSubject := p.curSubject
	savedPredicate := p.curPredicate
	// expect tokenStartCollection
	tok := p.next()
	if tok.typ != tokenStartCollection {
		return fmt.Errorf("Expected tokenStartCollection, got %v", tok)
	}
	// set curSubject to a new blank node bNode
	// set curPredicate to rdf:first
	next = p.peek()
	for next.typ != tokenEndCollection {
		err := p.parseObject()
		if err != nil {
			return err
		}
	}
	// TODO make sure this holds up for empty collections
	// expect tokenEndCollection
	tok := p.next()
	if tok.typ != tokenEndCollection {
		return fmt.Errorf("Expected tokenEndCollection, got %v", tok)
	}
	// emit triple p.curSubject rdf:rest rdf:nil
	p.curSubject = savedSubject
	p.curPredicate = savedPredicate
	return bNode, nil
}

func (p *Parser) parseObject() error {
	// TODO: implement
	// expect an object
	// where object = iri|blanknode|collection|blanknodepropertylist|literal
	tok := p.peek()
	switch tok.typ {
	case tokenIRIRef: // TODO: include PrefixedName
		p.next()
		iri, err := p.absIRI(tok.val)
		p.emitTriple(p.curSubject, p.curPredicate, iri)
		return err
	case tokenBlankNodeLabel:
		// TODO: what is the deal with the anon token?
		p.next()
		label := tok.val // TODO: parse the label out of token value
		bNode, err := p.blankNode(label)
		p.emitTriple(p.curSubject, p.curPredicate, bNode)
		return err
	case tokenStartCollection:
		bNode, err := p.parseCollection()
		p.emitTriple(p.curSubject, p.curPredicate, bNode)
		return err
	case tokenStartBlankNodePropertyList:
		bNode, err := p.parseBlankNodePropertyList()
		p.emitTriple(p.curSubject, p.curPredicate, bNode)
		return err
	case tokenInteger, tokenDecimal, tokenDouble, tokenTrue, tokenFalse, tokenStringLiteralQuote, tokenStringLiteralSingleQuote, tokenStringLiteralLongQuote, tokenStringLiteralLongSingleQuote:
		lit, err := p.parseLiteral()
		p.emitTriple(p.curSubject, p.curPredicate, lit)
		return err
	default:
		return fmt.Errorf("Expected object, got %v", tok)
	}
}

func (p *Parser) parseBlankNodePropertyList() (BlankNode, error) {
	savedSubject := p.curSubject
	savedPredicate := p.curPredicate
	// expect '[' token
	tok := p.next()
	if tok.typ != tokenStartBlankNodePropertyList {
		return fmt.Errorf("Expected tokenStartBlankNodePropertyList, got %v", tok)
	}
	// set curSubject to a new blank node bNode
	err := p.parsePredicateObjectList()
	if err != nil {
		return err
	}
	// expect ']' token
	tok = p.next()
	if tok.typ != tokenEndBlankNodePropertyList {
		return fmt.Errorf("Expected tokenEndBlankNodePropertyList, got %v", tok)
	}
	p.curSubject = savedSubject
	p.curPredicate = savedPredicate
	return bNode, nil
}

func (p *Parser) parseLiteral() (Literal, error) {
	tok := p.peek()
	switch tok.typ {
	case tokenInteger, tokenDecimal, tokenDouble:
		lit, err := p.parseNumericLiteral()
		return lit, err
	case tokenStringLiteralQuote, tokenStringLiteralSingleQuote, tokenStringLiteralLongQuote, tokenStringLiteralLongSingleQuote:
		lit, err := p.parseRDFLiteral()
		return lit, err
	case tokenTrue, tokenFalse:
		lit, err := p.parseBooleanLiteral()
		return lit, err
	}
}

func (p *Parser) parseNumericLiteral() (Literal, error) {
	// TODO: implement
}

func (p *Parser) parseRDFLiteral() (Literal, error) {
	// TODO: implement
}

func (p *Parser) parseBooleanLiteral() (Literal, error) {
	// TODO: implement
}
