package gon3

import (
	"unicode/utf8"
)

type naiveRegexLexer struct {
	name  string
	input string
	start int
	pos   int
}

func (l *naiveRegexLexer) nextToken() token {
	l.skipWhitespace()
	switch l.peek() {
	}
	panic("unimplemented")
}

func (l *naiveRegexLexer) peek() rune {
	if l.pos >= len(l.input) {
		return eof
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return r
}

func (l *naiveRegexLexer) skipWhitespace() {
	// TODO: implement
	panic("unimplemented")
}
