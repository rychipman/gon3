package gon3

import (
	"fmt"
)

type stateFn func(*lexer) stateFn

const (
	atPrefix       = "@prefix"
	atBase         = "@base"
	atKeywords     = "@keywords"
	atA            = "@a"
	atHas          = "@has"
	atIs           = "@is"
	atOf           = "@of"
	atForAll       = "@forAll"
	atForSome      = "@forSome"
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
	// TODO
}

func lexBarename(l *lexer) stateFn {
	// TODO
}

func lexExplicitURI(l *lexer) stateFn {
	// TODO
}
