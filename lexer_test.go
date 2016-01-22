package gon3

import (
	"testing"
)

func TestCharMatchLexer(t *testing.T) {
	input := `
@base <www.example.org> .
@prefix : <www.example.org/test/> .
	`
	expected := []tokenType{tokenAtBase, tokenIRIRef, tokenAtPrefix, tokenEndTriples, tokenAtPrefix, tokenPNameNS, tokenIRIRef, tokenEndTriples, tokenEOF}
	c := newCharMatchLexer("testLexer", input)
	toks := []tokenType{}
	var lastTok tokenType
	for lastTok != tokenEOF {
		lastTok = c.nextToken().typ
		toks = append(toks, lastTok)
		if lastTok == tokenError {
			t.Fatalf("Received an error token")
		}
	}
	for i, tokType := range toks {
		if expected[i] != tokType {
			t.Fatalf("Lexed tokens don't match expectations.\nExpected: %v\nActual: %v", expected, toks)
		}
	}
}
