package gon3

import (
	"fmt"
)

type stateFn func(*lexer) stateFn

const (
	atPrefix       = "@prefix"
	atBase         = "@base"
	endStatement   = "."
	endProperty    = ";"
	listSeparator  = ","
	beginPropList  = "["
	endPropList    = "]"
	beginFormula   = "{"
	endFormula     = "}"
	beginPathlist  = "("
	endPathlist    = ")"
	equals         = "="
	implies        = "=>"
	reverseImplies = "<="
	eof            = -1
	// TODO: add consts for ^,!,^^,a,<,>
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
	for {
		// I am a state function. I return another state function
		if l.next() == eof {
			break
		}
	}
	// reached EOF
	// TODO: emit error token if unparsed text in doc
}

func lexWhitespace(l *lexer) stateFn {

}

func lexNumericLiteral(l *lexer) stateFn {
	// TODO
}

func lexString(l *lexer) stateFn {
	// TODO
}

func lexVariable(l *lexer) stateFn {
	// TODO
}

func lexLangcode(l *lexer) stateFn {
	// TODO
}

func lexQname(l *lexer) stateFn {
	if l.input[l.pos] != l.qnamePrefix {
		l.errorf("Expected qname prefix %s, got %s", l.qnamePrefix, l.input[l.pos])
		return nil
	}
	l.pos += len(l.qnamePrefix)
	l.acceptRun(runQname)
	l.emit(tokenQname)
	return lexWhitespace
}

func lexBarename(l *lexer) stateFn {
	// TODO
}

func lexExplicitURI(l *lexer) stateFn {
	l.acceptRun(runURI)
	l.emit(tokenExplicitURI)
	return lexWhitespace
}
