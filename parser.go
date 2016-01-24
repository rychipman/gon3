package gon3

import (
	"fmt"
	"github.com/rychipman/easylex"
)

type Parser struct {
	// target data structure
	Graph Graph
	// parser state
	lex           lexer
	nextTok       chan easylex.Token
	baseURI       IRI
	namespaces    map[string]IRI //map[prefix]IRI // TODO: create prefix type
	bNodeLabels   map[string]BlankNode
	lastBlankNode BlankNode // TODO: initialize this to -1
	curSubject    RDFTerm   // TODO: create RDFTerm type (or perhaps interface)
	curPredicate  IRI
}

func NewParser(input string) *Parser {
	// initialize parser
	p := &Parser{
		Graph:         []*Triple{},
		lex:           easylex.Lex(input, lexDocument),
		nextTok:       make(chan easylex.Token, 1),
		baseURI:       "", // TODO
		namespaces:    map[string]IRI{},
		bNodeLabels:   map[string]BlankNode{},
		lastBlankNode: BlankNode{-1},
	}
}

func (p *Parser) Parse() (Graph, error) {
	var err error
	for { // while the next token is not an EOF
		err = p.parseStatement()
		if err != nil {
			break
		}
	}
	return p.Graph, err
}

func (p *Parser) peek() easylex.Token {
	select {
	case t := <-p.nextTok:
		p.nextTok <- t
		return t
	default:
		p.nextTok <- p.lex.NextToken()
	}
}

func (p *Parser) next() easylex.Token {
	select {
	case t := <-p.nextTok:
		return t
	default:
		p.nextTok <- p.lex.NextToken()
	}
}

func (p *Parser) emitTriple(subj RDFTerm, pred IRI, obj RDFTerm) { // TODO: work out typing things
	trip := &Triple{
		Subject:   subj,
		Predicate: pred,
		Object:    obj,
	}
	p.Graph = append(p.Graph, trip)
}

func (p *Parser) absIRI(iri string) (IRI, error) {
	// TODO: implement
	// if first char not '<', process as prefixed name
	// else if relative, resolve according to http://www.w3.org/TR/turtle/#sec-iri-references
	// finally, remove unicode escape sequences
	panic("unimplemented")
}

func (p *Parser) blankNode(label string) (BlankNode, error) {
	// TODO: when would we return an error?
	if node, present := p.bNodeLabels[label]; present {
		return node, nil
	}
	newNode := BlankNode{p.lastBlankNode.id + 1, "somelabelname"} // TODO: def not correct
	p.bNodeLabels[label] = newNode
	return newNode, nil
}

func (p *Parser) parseStatement() error {
	tok := p.peek()
	switch tok.Typ {
	case tokenAtPrefix:
		err := p.parsePrefix()
		if err != nil {
			return err
		}
	case tokenAtBase:
		err := p.parseBase()
		if err != nil {
			return err
		}
	case tokenSPARQLBase:
		err := p.ParseSPARQLBase()
		if err != nil {
			return err
		}
	case tokenSPARQLPrefix:
		err := p.ParseSPARQLPrefix()
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
	if tok.Typ != tokenAtPrefix {
		return fmt.Errorf("Expected tokenAtPrefix, got %v", tok)
	}
	// expect PNAME_NS token
	pNameNS := p.next()
	if pNameNS.Typ != tokenPNameNS {
		return fmt.Errorf("Expected tokenPNameNS, got %v", pNameNS)
	}
	// expect IRIREF token
	iriRef := p.next()
	if iriRef.Typ != tokenIRIRef {
		return fmt.Errorf("Expected tokenIRIRef, got %v", iriRef)
	}
	// expect "." token
	period := p.next()
	if period.Typ != tokenEndTriples {
		return fmt.Errorf("Expected tokenEndTriples, got %v", period)
	}
	// map a new namespace in parser state
	key := pNameNS.Val[:len(pNameNS.Val)-1]
	val, err := p.absIRI(iriRef.Val)
	p.namespaces[key] = val
	return err
}

func (p *Parser) parseBase() error {
	// expect '@base' token
	tok := p.next()
	if tok.Typ != tokenAtBase {
		return fmt.Errorf("Expected tokenAtBase, got %v", tok)
	}
	// expect IRIREF token
	iriRef := p.next()
	if iriRef.Typ != tokenIRIRef {
		return fmt.Errorf("Expected tokenIRIRef, got %v", iriRef)
	}
	// expect "." token
	period := p.next()
	if period.Typ != tokenEndTriples {
		return fmt.Errorf("Expected tokenEndTriples, got %v", period)
	}
	// TODO: require iriRef to be an absolute (or maybe prefixed?) iri
	// for now, assume it is an abs iri
	var err error
	p.baseURI, err = p.absIRI(iriRef.Val)
	return err
}

func (p *Parser) parseTriples() error {
	if true { // if "subject predicateobjectlist"
		err := p.parseSubject()
		if err != nil {
			return err
		}
		err = p.parsePredicateObjectList()
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
		if p.peek().Typ != tokenEndTriples {
			err = p.parsePredicateObjectList()
			if err != nil {
				return err
			}
		}
	}
	// expect period token
	tok := p.next()
	if tok.Typ != tokenEndTriples {
		return fmt.Errorf("Expected tokenEndTriples, got %v", tok)
	}
	return nil
}

func (p *Parser) parseSubject() error {
	tok := p.peek()
	// expect a valid subject term, which is one of
	// iri|blanknode|collection
	switch tok.Typ {
	case tokenIRIRef: // TODO: include PrefixedName here
		p.next()
		iri, err := p.absIRI(tok.Val)
		p.curSubject = iri
		return err
	case tokenBlankNodeLabel:
		// TODO: what is the deal with the anon token?
		p.next()
		label := tok.Val // TODO: parse the label out of token value
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
	for p.peek().Typ != tokenEndTriples {
		// expect semicolon token
		tok := p.next()
		if tok.Typ != tokenPredicateListSeparator {
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
	switch tok.Typ {
	case tokenA:
		// TODO: remove magic string
		pred, err := p.absIRI("<http://www.w3.org/1999/02/22-rdf-syntax-ns#type>")
		p.curPredicate = pred
		return err
	case tokenIRIRef: // TODO: include PrefixedName here
		iri, err := p.absIRI(tok.Val)
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
	for p.peek().Typ != tokenEndTriples {
		// expect comma token
		tok := p.next()
		if tok.Typ != tokenObjectListSeparator {
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
	savedSubject := p.curSubject
	savedPredicate := p.curPredicate
	// expect tokenStartCollection
	tok := p.next()
	if tok.Typ != tokenStartCollection {
		return BlankNode{}, fmt.Errorf("Expected tokenStartCollection, got %v", tok)
	}
	// TODO: set curSubject to a new blank node bNode
	// TODO: set curPredicate to rdf:first
	next := p.peek()
	for next.Typ != tokenEndCollection {
		err := p.parseObject()
		if err != nil {
			return BlankNode{}, err
		}
	}

	// TODO: make sure this holds up for empty collections.
	// Also note that empty collections are probably what tokenAnon is.

	// expect tokenEndCollection
	tok = p.next()
	if tok.Typ != tokenEndCollection {
		return BlankNode{}, fmt.Errorf("Expected tokenEndCollection, got %v", tok)
	}
	// TODO: emit triple p.curSubject rdf:rest rdf:nil
	p.curSubject = savedSubject
	p.curPredicate = savedPredicate
	bNode := BlankNode{} // TODO: return bNode created above
	return bNode, nil
}

func (p *Parser) parseObject() error {
	// expect an object
	// where object = iri|blanknode|collection|blanknodepropertylist|literal
	tok := p.peek()
	switch tok.Typ {
	case tokenIRIRef: // TODO: include PrefixedName
		p.next()
		iri, err := p.absIRI(tok.Val)
		p.emitTriple(p.curSubject, p.curPredicate, iri)
		return err
	case tokenBlankNodeLabel:
		// TODO: what is the deal with the anon token?
		p.next()
		label := tok.Val // TODO: parse the label out of token value
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
	if tok.Typ != tokenStartBlankNodePropertyList {
		return BlankNode{}, fmt.Errorf("Expected tokenStartBlankNodePropertyList, got %v", tok)
	}
	// TODO: set curSubject to a new blank node bNode
	err := p.parsePredicateObjectList()
	if err != nil {
		return BlankNode{}, err
	}
	// expect ']' token
	tok = p.next()
	if tok.Typ != tokenEndBlankNodePropertyList {
		return BlankNode{}, fmt.Errorf("Expected tokenEndBlankNodePropertyList, got %v", tok)
	}
	p.curSubject = savedSubject
	p.curPredicate = savedPredicate
	bNode := BlankNode{} // TODO: return the bNode created above
	return bNode, nil
}

func (p *Parser) parseLiteral() (Literal, error) {
	tok := p.peek()
	switch tok.Typ {
	case tokenInteger, tokenDecimal, tokenDouble:
		lit, err := p.parseNumericLiteral()
		return lit, err
	case tokenStringLiteralQuote, tokenStringLiteralSingleQuote, tokenStringLiteralLongQuote, tokenStringLiteralLongSingleQuote:
		lit, err := p.parseRDFLiteral()
		return lit, err
	case tokenTrue, tokenFalse:
		lit, err := p.parseBooleanLiteral()
		return lit, err
	default:
		return nil, fmt.Errorf("Expected a literal token, got %v", tok)
	}
	panic("unreachable")
}

func (p *Parser) parseNumericLiteral() (Literal, error) {
	// TODO: implement
	panic("unimplemented")
}

func (p *Parser) parseRDFLiteral() (Literal, error) {
	// TODO: implement
	panic("unimplemented")
}

func (p *Parser) parseBooleanLiteral() (Literal, error) {
	// TODO: implement
	panic("unimplemented")
}
