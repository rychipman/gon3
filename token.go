package gon3

import (
	"fmt"
)

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF
	tokenNumericLiteral
	tokenString
	tokenVariable
	tokenLangcode
	tokenQname
	tokenBarename
	tokenExplicituri
	// TODO: allow document children in any order
	// TODO: add quickvariable support (is this underscore namespace?)
	// TODO: add boolean support
	// TODO: add more numeric token types
)

type token struct {
	typ tokenType
	val string
}

func (t token) String() string {
	switch t.typ {
	case tokenError:
		return t.val
	case tokenEOF:
		return "EOF"
	}
	if len(t.val) > 10 {
		return fmt.Sprintf("%.10q...", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}
