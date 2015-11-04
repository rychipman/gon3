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
	// TODO: add consts for ^,!,^^,a
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
