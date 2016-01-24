package gon3

import (
	"fmt"
	"github.com/rychipman/easylex"
	"testing"
)

func TestParser(t *testing.T) {
	for _, tokSet := range parseTests {
		mock := newMockLexer(tokSet.tokens...)
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

var parseTests = []struct {
	name       string
	tokens     []easylex.Token
	shouldPass bool
}{}
