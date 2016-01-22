package gon3

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type charMatchLexer struct {
	name   string
	input  string
	state  stateFn
	start  int
	pos    int
	width  int
	tokens chan token
}

func (l *charMatchLexer) nextToken() token {
	for {
		select {
		case tok := <-l.tokens:
			return tok
		default:
			l.state = l.state(l)
		}
	}
}

func (l *charMatchLexer) emit(t tokenType) {
	l.tokens <- token{
		t,
		l.input[l.start:l.pos],
	}
	l.start = l.pos
}

func (l *charMatchLexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{
		tokenError,
		fmt.Sprintf(format, args),
	}
	return nil
}

func (l *charMatchLexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// ignore resets l.start to the current value of l.pos.
// this ignores all the runes processed since the last
// call to ignore() or emit().
func (l *charMatchLexer) ignore() {
	l.start = l.pos
}

// backup decrements l.pos by the width of the last rune
// processed. backup can only be called once per call to
// next().
func (l *charMatchLexer) backup() {
	l.pos -= l.width
}

// peek returns the value of the rune at l.pos + 1, but
// does not mutate lexer state.
func (l *charMatchLexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// accept will be equivalent to calling next() if the next
// rune is in the provided string. If the next rune is not
// one of those provided, the lexer state is unchanged.
func (l *charMatchLexer) accept(valid string) bool {
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
func (l *charMatchLexer) acceptRun(valid matcher) bool {
	return valid.match(l)
}
