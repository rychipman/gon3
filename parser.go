package gon3

import (
	"fmt"
	"github.com/rychipman/easylex"
	"net/url"
	"strings"
)

type Parser struct {
	// target data structure
	Graph Graph
	// parser state
	lex           lexer
	nextTok       chan easylex.Token
	baseURI       IRI
	namespaces    map[string]IRI
	bNodeLabels   map[string]BlankNode
	lastBlankNode BlankNode
	curSubject    RDFTerm
	curPredicate  IRI
}

func NewParser(input string) *Parser {
	base, _ := url.Parse("") // TODO: properly initialize baseuri
	// initialize parser
	p := &Parser{
		Graph:         []*Triple{},
		lex:           easylex.Lex(input, lexDocument),
		nextTok:       make(chan easylex.Token, 1),
		baseURI:       IRI{base},
		namespaces:    map[string]IRI{},
		bNodeLabels:   map[string]BlankNode{},
		lastBlankNode: BlankNode{-1, ""},
		curSubject:    IRI{},
		curPredicate:  IRI{},
	}
	return p
}

func (p *Parser) Parse() (Graph, error) {
	var err error
	var done bool
	for { // while the next token is not an EOF
		done, err = p.parseStatement()
		if done || err != nil {
			break
		}
	}
	return p.Graph, err
}

func (p *Parser) peek() easylex.Token {
	for {
		select {
		case t := <-p.nextTok:
			p.nextTok <- t
			return t
		default:
			p.nextTok <- p.lex.NextToken()
		}
	}
}

func (p *Parser) next() easylex.Token {
	for {
		select {
		case t := <-p.nextTok:
			fmt.Printf("TOKEN: %s (type %s)\n", t, t.Typ)
			return t
		default:
			p.nextTok <- p.lex.NextToken()
		}
	}
}

func (p *Parser) expect(typ easylex.TokenType) (easylex.Token, error) {
	tok := p.next()
	if tok.Typ != typ {
		return tok, fmt.Errorf("Expected %s, got %s (type %s)", typ, tok.Val, tok.Typ)
	}
	return tok, nil
}

func (p *Parser) emitTriple(subj RDFTerm, pred IRI, obj RDFTerm) {
	trip := &Triple{
		Subject:   subj,
		Predicate: pred,
		Object:    obj,
	}
	p.Graph = append(p.Graph, trip)
}

func (p *Parser) blankNode(label string) BlankNode {
	if label == "" {
		return p.newBlankNode()
	} else if node, present := p.bNodeLabels[label]; present {
		return node
	}
	newNode := p.newBlankNode()
	p.bNodeLabels[label] = newNode
	return newNode
}

func (p *Parser) newBlankNode() BlankNode {
	id := p.lastBlankNode.Id + 1
	label := fmt.Sprintf("a%d", id)
	b := BlankNode{
		Id:    id,
		Label: label,
	}
	p.lastBlankNode = b
	return b
}

func (p *Parser) resolvePName(pname string) (IRI, error) {
	strs := strings.Split(pname, ":")
	prefix := strs[0]
	name := strs[1]
	if iri, present := p.namespaces[prefix]; present {
		rel, err := iriRefToURL(name)
		if err != nil {
			return IRI{}, err
		}
		resolved := iri.url.ResolveReference(rel) // TODO: make sure this properly resolves weird things
		return IRI{resolved}, nil
	}
	return IRI{}, fmt.Errorf("Prefix %q not found in declared namespaces", prefix)
}

func (p *Parser) resolveIRI(iri string) (IRI, error) {
	rel, err := iriRefToURL(iri)
	if err != nil {
		return IRI{}, err
	}
	url := p.baseURI.url.ResolveReference(rel)
	return IRI{url}, nil
}

func (p *Parser) parseStatement() (bool, error) {
	tok := p.peek()
	switch tok.Typ {
	case easylex.TokenError:
		return false, fmt.Errorf("Received tokenError: %q", tok)
	case easylex.TokenEOF:
		return true, nil
	case tokenAtPrefix:
		err := p.parsePrefix()
		if err != nil {
			return false, err
		}
	case tokenAtBase:
		err := p.parseBase()
		if err != nil {
			return false, err
		}
	case tokenSPARQLBase:
		err := p.parseSPARQLBase()
		if err != nil {
			return false, err
		}
	case tokenSPARQLPrefix:
		err := p.parseSPARQLPrefix()
		if err != nil {
			return false, err
		}
	default:
		err := p.parseTriples()
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

func (p *Parser) parsePrefix() error {
	_, err := p.expect(tokenAtPrefix)
	if err != nil {
		return err
	}
	pNameNS, err := p.expect(tokenPNameNS)
	if err != nil {
		return err
	}
	iriRef, err := p.expect(tokenIRIRef)
	if err != nil {
		return err
	}
	_, err = p.expect(tokenEndTriples)
	if err != nil {
		return err
	}
	// map a new namespace in parser state
	key := pNameNS.Val[:len(pNameNS.Val)-1]
	val, err := newIRIFromString(iriRef.Val)
	p.namespaces[key] = val
	return err
}

func (p *Parser) parseBase() error {
	_, err := p.expect(tokenAtBase)
	if err != nil {
		return err
	}
	iriRef, err := p.expect(tokenIRIRef)
	if err != nil {
		return err
	}
	_, err = p.expect(tokenEndTriples)
	if err != nil {
		return err
	}
	// TODO: validate IRI?
	p.baseURI, err = newIRIFromString(iriRef.Val)
	return err
}

func (p *Parser) parseSPARQLPrefix() error {
	_, err := p.expect(tokenSPARQLPrefix)
	if err != nil {
		return err
	}
	pNameNS, err := p.expect(tokenPNameNS)
	if err != nil {
		return err
	}
	iriRef, err := p.expect(tokenIRIRef)
	if err != nil {
		return err
	}
	// map a new namespace in parser state
	key := pNameNS.Val[:len(pNameNS.Val)-1]
	val, err := newIRIFromString(iriRef.Val)
	p.namespaces[key] = val
	return err
}

func (p *Parser) parseSPARQLBase() error {
	_, err := p.expect(tokenSPARQLBase)
	if err != nil {
		return err
	}
	iriRef, err := p.expect(tokenIRIRef)
	if err != nil {
		return err
	}
	// TODO: validate IRI?
	p.baseURI, err = newIRIFromString(iriRef.Val)
	return err
}

func (p *Parser) parseTriples() error {
	if p.peek().Typ == tokenStartBlankNodePropertyList {
		// if "blanknodepropertylist predicateobjectlist?"
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
	} else {
		// if "subject predicateobjectlist"
		err := p.parseSubject()
		if err != nil {
			return err
		}
		err = p.parsePredicateObjectList()
		if err != nil {
			return err
		}
	}
	_, err := p.expect(tokenEndTriples)
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) parseSubject() error {
	tok := p.peek()
	// expect a valid subject term, which is one of
	// iri|blanknode|collection
	switch tok.Typ {
	case tokenIRIRef:
		p.next()
		iri, err := p.resolveIRI(tok.Val)
		p.curSubject = iri
		return err
	case tokenPNameLN, tokenPNameNS:
		p.next()
		iri, err := p.resolvePName(tok.Val)
		p.curSubject = iri
		return err
	case tokenBlankNodeLabel:
		p.next()
		label := strings.Split(tok.Val, ":")[1]
		bNode := p.blankNode(label)
		p.curSubject = bNode
		return nil
	case tokenAnon:
		p.next()
		p.curSubject = p.blankNode("")
		return nil
	case tokenStartCollection:
		bNode, err := p.parseCollection()
		p.curSubject = bNode
		return err
	default:
		return fmt.Errorf("Expected a subject, got %v (type %s)", tok, tok.Typ)
	}
}

func (p *Parser) parsePredicateObjectList() error {
	err := p.parsePredicate()
	if err != nil {
		return err
	}
	err = p.parseObjectList()
	if err != nil {
		return err
	}
	for p.peek().Typ == tokenPredicateListSeparator {
		_, err = p.expect(tokenPredicateListSeparator)
		if err != nil {
			return err
		}
		tok := p.peek()
		switch tok.Typ {
		case tokenA, tokenIRIRef, tokenPNameLN, tokenPNameNS:
			// if there is a predicate
			err = p.parsePredicate()
			if err != nil {
				return err
			}
			err = p.parseObjectList()
			if err != nil {
				return err
			}
		default:
			// done parsing predicateobjectlist
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
		pred, err := newIRIFromString("<http://www.w3.org/1999/02/22-rdf-syntax-ns#type>")
		p.curPredicate = pred
		return err
	case tokenIRIRef:
		iri, err := p.resolveIRI(tok.Val)
		p.curPredicate = iri
		return err
	case tokenPNameLN, tokenPNameNS:
		iri, err := p.resolvePName(tok.Val)
		p.curPredicate = iri
		return err
	default:
		return fmt.Errorf("Expected predicate, got %v (type %s)", tok, tok.Typ)
	}
}

func (p *Parser) parseObjectList() error {
	err := p.parseObject()
	if err != nil {
		return err
	}
	for p.peek().Typ == tokenObjectListSeparator {
		// expect comma token
		_, err = p.expect(tokenObjectListSeparator)
		if err != nil {
			return err
		}
		err = p.parseObject()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseCollection() (BlankNode, error) {
	savedSubject := p.curSubject
	savedPredicate := p.curPredicate
	_, err := p.expect(tokenStartCollection)
	if err != nil {
		return BlankNode{}, err
	}
	bNode := p.blankNode("")
	p.curSubject = bNode
	p.curPredicate, _ = newIRIFromString("<http://www.w3.org/1999/02/22-rdf-syntax-ns#first>") // TODO: make this a const or something
	next := p.peek()
	for next.Typ != tokenEndCollection {
		err := p.parseObject()
		if err != nil {
			return BlankNode{}, err
		}
		next = p.peek()
	}

	_, err = p.expect(tokenEndCollection)
	if err != nil {
		return BlankNode{}, err
	}
	// TODO: use consts
	rdfRest, _ := newIRIFromString("<http://www.w3.org/1999/02/22-rdf-syntax-ns#rest>")
	rdfNil, _ := newIRIFromString("<http://www.w3.org/1999/02/22-rdf-syntax-ns#nil>")
	p.emitTriple(p.curSubject, rdfRest, rdfNil)
	p.curSubject = savedSubject
	p.curPredicate = savedPredicate
	return bNode, nil
}

func (p *Parser) parseObject() error {
	// expect an object
	// where object = iri|blanknode|collection|blanknodepropertylist|literal
	tok := p.peek()
	switch tok.Typ {
	case tokenIRIRef:
		p.next()
		iri, err := p.resolveIRI(tok.Val)
		p.emitTriple(p.curSubject, p.curPredicate, iri)
		return err
	case tokenPNameLN, tokenPNameNS:
		p.next()
		iri, err := p.resolvePName(tok.Val)
		p.emitTriple(p.curSubject, p.curPredicate, iri)
		return err
	case tokenBlankNodeLabel:
		p.next()
		label := strings.Split(tok.Val, ":")[1]
		bNode := p.blankNode(label)
		p.emitTriple(p.curSubject, p.curPredicate, bNode)
		return nil
	case tokenAnon:
		p.next()
		bNode := p.blankNode("")
		p.emitTriple(p.curSubject, p.curPredicate, bNode)
		return nil
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
		return fmt.Errorf("Expected object, got %v (type %s)", tok, tok.Typ)
	}
}

func (p *Parser) parseBlankNodePropertyList() (BlankNode, error) {
	savedSubject := p.curSubject
	savedPredicate := p.curPredicate
	// expect '[' token
	_, err := p.expect(tokenStartBlankNodePropertyList)
	if err != nil {
		return BlankNode{}, err
	}
	bNode := p.blankNode("bnodeproplist")
	p.curSubject = bNode
	err = p.parsePredicateObjectList()
	if err != nil {
		return BlankNode{}, err
	}
	// expect ']' token
	_, err = p.expect(tokenEndBlankNodePropertyList)
	if err != nil {
		return BlankNode{}, err
	}
	p.curSubject = savedSubject
	p.curPredicate = savedPredicate
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
		return Literal{}, fmt.Errorf("Expected a literal token, got %v", tok)
	}
	panic("unreachable")
}

func (p *Parser) parseNumericLiteral() (Literal, error) {
	// TODO: replace these with consts
	xsdInteger, _ := newIRIFromString("<http://www.w3.org/2001/XMLSchema#integer>")
	xsdDecimal, _ := newIRIFromString("<http://www.w3.org/2001/XMLSchema#decimal>")
	xsdDouble, _ := newIRIFromString("<http://www.w3.org/2001/XMLSchema#double>")
	tok := p.next()
	switch tok.Typ {
	case tokenInteger:
		lit := Literal{
			tok.Val,
			xsdInteger,
			"",
		}
		return lit, nil
	case tokenDecimal:
		lit := Literal{
			tok.Val,
			xsdDecimal,
			"",
		}
		return lit, nil
	case tokenDouble:
		lit := Literal{
			tok.Val,
			xsdDouble,
			"",
		}
		return lit, nil
	default:
		return Literal{}, fmt.Errorf("Expected a numeric literal token, got %s", tok)
	}
}

func (p *Parser) parseRDFLiteral() (Literal, error) {
	tok := p.next()
	stringDT, _ := newIRIFromString("<http://www.w3.org/2001/XMLSchema#string>") // TODO: const
	lit := Literal{
		LexicalForm: "",
		DatatypeIRI: stringDT,
		LanguageTag: "",
	}
	switch tok.Typ {
	case tokenStringLiteralQuote, tokenStringLiteralSingleQuote, tokenStringLiteralLongQuote, tokenStringLiteralLongSingleQuote:
		lit.LexicalForm = lexicalForm(tok.Val)
	default:
		return Literal{}, fmt.Errorf("Expected a string literal token, got %s", tok)
	}
	if p.peek().Typ == tokenLangTag {
		langtag := p.next()
		lit.LanguageTag = langtag.Val[1:]
		// TODO: make this a const
		dIRI, _ := newIRIFromString("<http://www.w3.org/1999/02/22-rdf-syntax-ns#langString>")
		lit.DatatypeIRI = dIRI
	} else if p.peek().Typ == tokenLiteralDatatypeTag {
		p.next()
		tok := p.next()
		switch tok.Typ {
		case tokenIRIRef:
			iri, err := p.resolveIRI(tok.Val)
			if err != nil {
				return Literal{}, err
			}
			lit.DatatypeIRI = iri
		case tokenPNameLN, tokenPNameNS:
			iri, err := p.resolvePName(tok.Val)
			if err != nil {
				return Literal{}, err
			}
			lit.DatatypeIRI = iri
		default:
			return Literal{}, fmt.Errorf("Expected an IRI or PName, got %s (type %s)", tok.Val, tok.Typ)
		}
	}
	return lit, nil
}

func (p *Parser) parseBooleanLiteral() (Literal, error) {
	// TODO: make this a const
	xsdBoolean, _ := newIRIFromString("<http://www.w3.org/2001/XMLSchema#boolean>")
	tok := p.next()
	switch tok.Typ {
	case tokenTrue:
		lit := Literal{
			tok.Val,
			xsdBoolean,
			"",
		}
		return lit, nil
	case tokenFalse:
		lit := Literal{
			tok.Val,
			xsdBoolean,
			"",
		}
		return lit, nil
	default:
		return Literal{}, fmt.Errorf("Expected a boolean literal token, got %s", tok)
	}
}
