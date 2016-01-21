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
	switch l.peek() {
	case "@":
		return lexAtStatement
	case "[":
		return lexBlankNodePropertyList
	case "(":
		return lexCollection
	case "_":
		return lexBlankNodeLabel
	case "<":
		return lexIRIRef
	case "'":
		// TODO: return proper statefn
	case "\"":
		// TODO: return proper statefn
	default:
		return lexPName
	}
}

func lexAtStatement(l *lexer) stateFn {
	// lex prefix/base directives or langtag
	if !l.accept("@") {
		return l.errorf("lexAtStatement called, but '@' not found")
	}
	if !l.acceptRun(runAlphabet) {
		return l.errorf("Expected Alphabet while lexing '@...', got %s", l.input[l.pos-1:l.pos])
	}
	if l.accept("-") {
		if !l.acceptRun(runAlphanumeric) {
			return l.errorf("Expected Alphanumeric while lexing langtag, got %s", l.input[l.pos-1:l.pos])
		}
		l.emit(tokenLangTag)
		return lexDocument
	}
	if l.atPrefix() {
		l.emit(tokenAtPrefix)
		return lexDocument
	}
	if l.atBase() {
		l.emit(tokenAtBase)
		return lexDocument
	}
	l.emit(tokenLangTag)
	return lexDocument
}

func lexBlankNodePropertyList(l *lexer) stateFn {
	// TODO: implement
}

func lexCollection(l *lexer) stateFn {
	// TODO: implement
}

func lexBlankNodeLabel(l *lexer) stateFn {
	// lex bnode label
	if !l.accept("_") {
		return l.errorf("lexAtStatement called, but '@' not found")
	}
	if !l.accept(":") {
		return l.errorf("Expected ':' while lexing bnode label, got %s", l.input[l.pos-1:l.pos])
	}
	if !l.acceptRun(runPNCharsU + runDigits) {
		return l.errorf("Expected PNCharsU or Digits while lexing bnode label, got %s", l.input[l.pos-1:l.pos])
	}
	if l.acceptRun(runPNChars + ".") {
		for l.acceptRun(runPNChars + ".") {
		}
		l.backup()
		if !l.accept(runPNChars) {
			return l.errorf("Expected PNChars for last char of bnode label, got %s", l.input[l.pos-1:l.pos])
		}
	}
	l.emit(tokenBlankNodeLabel)
	return lexDocument
}

func lexIRIRef(l *lexer) stateFn {
	// TODO: implement
}

func lexPName(l *lexer) stateFn {
	// TODO: implement
}

func (l *lexer) atPrefix() bool {

}

func (l *lexer) atBase() bool {

}
