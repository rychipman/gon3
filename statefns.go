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
	runDigits       = "0123456789"
	runAlphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	runNumeric      = runDigits + "-."
	runWhitespace   = " \n\t\v\f\r"    // TODO: check this against TR
	runPNCharsBase  = runAlphabet + "" // TODO: complete with unicode escapes
	runPNCharsU     = runPNCharsBase + "_"
	runPNChars      = runPNCharsU + runDigits + "-" // TODO: complete with unicode escapes
	runEscapable    = "_~.-!$&'()*+,;=/?#@%"
)

func lexDocument(l *lexer) stateFn {
	l.acceptRun(runWhitespace)
	l.ignore()
	switch l.next() {
	case "@":
		// lex prefix/base directives or langtag
		if !l.acceptRun(runAlphabet) {
			l.errorf("Expected Alphabet while lexing '@...', got %s", l.input[l.pos-1:l.pos])
		}
		if l.accept("-") {
			if !l.acceptRun(runAlphanumeric) {
				l.errorf("Expected Alphanumeric while lexing langtag, got %s", l.input[l.pos-1:l.pos])
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
		return lexDocument
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
