package gon3

import (
	"fmt"
)

type lexer struct {
	name   string
	input  string
	state  stateFn
	start  int
	pos    int
	width  int
	tokens chan token
}

func lex(name, input string) *lexer {
	l := &lexer{
		name:   name,
		input:  input,
		state:  lexDocument,
		tokens: make(chan tokens, 2),
	}
	return l
}

func (l *lexer) nextToken() token {
	for {
		select {
		case tok := <-l.tokens:
			return tok
		default:
			l.state = l.state(l)
		}
	}
}

func (l *lexer) emit(t tokenType) {
	l.tokens <- token{
		t,
		l.input[l.start:l.pos],
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

// ignore resets l.start to the current value of l.pos.
// this ignores all the runes processed since the last
// call to ignore() or emit().
func (l *lexer) ignore() {
	l.start = l.pos
}

// backup decrements l.pos by the width of the last rune
// processed. backup can only be called once per call to
// next().
func (l *lexer) backup() {
	l.pos -= l.width
}

// peek returns the value of the rune at l.pos + 1, but
// does not mutate lexer state.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// accept will be equivalent to calling next() if the next
// rune is in the provided string. If the next rune is not
// one of those provided, the lexer state is unchanged.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun will continue accepting tokens until it encounters
// one not included in the string of valid tokens.
// acceptRun is preferable to blind calls to next() because it
// forces explicit declaration of what characters are allowed.
func (l *lexer) acceptRun(valid string) bool {
	success := false
	for strings.IndexRune(valid, l.next()) >= 0 {
		if !success {
			success = true
		}
	}
	l.backup()
	return success
}
