package gon3

import (
	"fmt"
)

type lexer struct {
	name   string
	input  string
	start  int
	pos    int
	width  int
	tokens chan token
}

func lex(name, input string) (*lexer, chan tokens) {
	l := &lexer{
		name:   name,
		input:  input,
		tokens: make(chan tokens),
	}
	go l.run()
	return l, l.tokens
}

func (l *lexer) run() {
	for state := lexDocument; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *lexer) emit(t tokenType) {
	l.tokens <- token{
		t,
		l.input[l.start:l.pos],
	}
	l.start = l.pos
}

func (l *lexer) emitf(t tokenType, format string, args ...interface{}) {
	l.tokens <- token{
		t,
		fmt.Sprintf(format, args),
	}
	l.start = l.pos
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{
		tokenError,
		fmt.Sprintf(format, args),
	}
	return nil
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() int {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}
