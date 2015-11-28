package gon3

type Parser struct {
	// target data structure
	Graph *Graph
	// parser state
	lex          *lexer // TODO: initialize lexer
	baseURI      IRI
	namespaces   map[prefix]IRI       // TODO: create prefix type
	bNodeLabels  map[string]BlankNode // TODO: update what the bnode type should look like
	curSubject   RDFTerm              // TODO: create RDFTerm type (or perhaps interface)
	curPredicate RDFTerm
}

func (p *Parser) Parse(text string) (*Graph, error) {
	// initialize fields
	p.Graph = Graph{}
	p.namespaces = map[prefix]IRI{}
	bNodeLabels = map[string]BlankNode{}

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
	tok = p.peek()
	switch tok.typ {
	case tokenAtPrefix:
		return p.parsePrefix()
	case tokenAtBase:
		return p.parseBase()
		// TODO: support tokenSPARQLBase, token SPARQLPrefix
	default:
		return p.parseTriples()
	}
}

func (p *Parser) parsePrefix() error {
	// TODO: implement
	// expect '@prefix' token
	// expect PNAME_NS token
	// expect IRIREF token
	// expect '.' token
	// map a new namespace in parser state
}

func (p *Parser) parseBase() error {
	// TODO: implement
	// expect '@base' token
	// expect IRIREF token
	// expect '.' token
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
		err := p.parseBlankNodePropertyList()
		if err != nil {
			return err
		}
		// TODO: parse a predicateobjectlist if we have one
	}
}

func (p *Parser) parseSubject() error {
	// TODO: implement
	// expect a valid subject term, which is one of
	// iri|blanknode|collection
	// parse subject and set p.curSubject
	// return error if exists
}

func (p *Parser) parsePredicateObjectList() error {
	// TODO: figure out details of when semicolon is required
	err := p.parsePredicate()
	if err != nil {
		return err
	}
	err = p.parseObjectList()
	if err != nil {
		return err
	}
	// TODO: continue looping if still have '; predicate objectlist'
}

func (p *Parser) parsePredicate() error {
	// TODO: implement
	// expect token 'a' or an iri
	// if 'a', replace with iri 'http://www.w3.org/1999/02/22-rdf-syntax-ns#type'
	// set curPredicate
	// return error if exists
}

func (p *Parser) parseObjectList() error {
	// TODO: implement
	// expect comma-separated list of objects
	// call parseObject()
	// return an error if we have one
	// pop off predicate at end of objectlist
}

func (p *Parser) parseObject() error {
	// TODO: implement
	// expect an object
	// where object = iri|blanknode|collection|blanknodepropertylist|literal
	// emit a new triple into the graph object
}

func (p *Parser) parseBlankNodePropertyList() error {
	// TODO: implement
	// expect '[' token
	// name a new blank node and push cursubject
	err = p.parsePredicateObjectList()
	if err != nil {
		return err
	}
	// expect ']' token
	// emit a triple
	// return an error if exists
}

func (p *Parser) parseIRI() error {
	// TODO: implement
	// expect IRIREF or prefixedname
}

func (p *Parser) parseBlankNode() error {
	// TODO: implement
	// expech a blank node label or '[  ]'
}

func (p *Parser) parseCollection() error {
	// TODO: implement
	// figure out exactly when triples should be emitted
	// use rdf:first and rdf:rest
	// create blank nodes
}

func (p *Parser) parseLiteral() error {
	// TODO: implement
}
