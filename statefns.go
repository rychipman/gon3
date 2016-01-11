package gon3

import (
	"fmt"
)

type stateFn func(*lexer) stateFn

const (
	runAlphabet     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	runWhitespace   = " \n\t\v\f\r"
	runDigits       = "0123456789"
	runAlphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	// TODO: create valid runs for each token
	runURI     = runAlphanumeric + runWhitespace + "<>/:."
	runQname   = runAlphabet
	runNumeric = runDigits + "-."
)

func lexDocument(l *lexer) stateFn {
	l.acceptRun(runWhitespace)
	switch l.next() {
	// prefix/base directives
	case "@prefix":
		// lex prefix
	case "@base":
		// lex base
	// subject constructions
	case "[":
		// lex blank node prop list
	case "(":
		// lex collection
	case "_":
		// lex bnode label
	case "<":
		// lex iri
	default:
		// lex pname
	}
}

func lexPrefix(l *lexer) stateFn {
	// TODO: implement
}

func lexBase(l *lexer) stateFn {
	// TODO: implement
}

func lexBlankNodePropertyList(l *lexer) stateFn {
	// TODO: implement
}

func lexCollection(l *lexer) stateFn {
	// TODO: implement
}
