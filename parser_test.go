package gon3

import (
	"fmt"
	"github.com/rychipman/easylex"
	"testing"
)

func TestParserControlFlow(t *testing.T) {
	for _, tokSet := range parseFlowTests {
		mock := newTypeMockLexer(tokSet.tokens...)
		p := NewParser("")
		p.lex = mock
		_, err := p.Parse()
		passed := true
		if err != nil {
			passed = false
		}
		if passed != tokSet.shouldPass {
			t.Fatalf("FAIL test %q. Expected pass=%t, got pass=%t", tokSet.name, tokSet.shouldPass, passed)
		}
		fmt.Printf("Passed test %q\n", tokSet.name)
	}
}

var parseFlowTests = []struct {
	name       string
	tokens     []easylex.TokenType
	shouldPass bool
}{
	{
		"Simple @base declaration",
		[]easylex.TokenType{tokenAtBase, tokenIRIRef, tokenEndTriples, easylex.TokenEOF},
		true,
	},
	{
		"Malformed @base declaration",
		[]easylex.TokenType{tokenAtBase, tokenIRIRef, easylex.TokenEOF},
		false,
	},
}
