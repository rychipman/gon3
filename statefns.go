package gon3

import (
	"fmt"
)

type stateFn func(*lexer) stateFn

const (
	eof = -1
)

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
	case "@":
		if !l.acceptRun(runAlphabet) {
			// TODO: error
		}
		if l.accept("-") {
			if !l.acceptRun(runAlphanumeric) {
				// TODO: error
			}
			l.emit(tokenLangTag)
			return lexDocument
		}
		if l.atPrefix() { // TODO: create this fn
			l.emit(tokenAtPrefix)
			return lexDocument
		}
		if l.atBase() { // TODO: create this fn
			l.emit(tokenAtBase)
			return lexDocument
		}
		l.emit(tokenLangTag)
	case "[":
		// lex blank node prop list
	case "(":
		// lex collection
	case "_":
		// lex bnode label
	case "<":
		// lex iri
	case "'":
	case "\"":
	default:
		// lex pname
	}
}

func (l *lexer) atPrefix() bool {

}

func (l *lexer) atBase() bool {

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
