package gon3

import (
	"fmt"
	"github.com/rychipman/easylex"
	"io/ioutil"
	"testing"
)

func TestLexer(t *testing.T) {
	for _, testName := range positiveTests {
		testFile := "./tests/turtle/lex/" + testName
		b, err := ioutil.ReadFile(testFile)
		if err != nil {
			t.Fatalf("Error reading test file %s", testFile)
		}
		fmt.Printf("\nStarting test %s\n", testName)
		l := easylex.Lex(string(b), lexDocument)
		for {
			token := l.NextToken()
			tokenType := token.Typ
			if tokenType == easylex.TokenEOF {
				fmt.Printf("Test %s passed.\n", testName)
				break
			}
			if tokenType == easylex.TokenError {
				t.Fatalf("Test %s failed.", testFile)
			}
			//fmt.Printf("%s -- %s\n", tokenType, token)
		}
	}
}

var positiveTests []string = []string{
	"turtle-syntax-base-01.ttl",
	"turtle-syntax-base-02.ttl",
	"turtle-syntax-base-03.ttl",
	"turtle-syntax-base-04.ttl",
	"turtle-syntax-blank-label.ttl",
	"turtle-syntax-bnode-01.ttl",
	"turtle-syntax-bnode-02.ttl",
	"turtle-syntax-bnode-03.ttl",
	"turtle-syntax-bnode-04.ttl",
	"turtle-syntax-bnode-05.ttl",
	"turtle-syntax-bnode-06.ttl",
	"turtle-syntax-bnode-07.ttl",
	"turtle-syntax-bnode-08.ttl",
	"turtle-syntax-bnode-09.ttl",
	"turtle-syntax-bnode-10.ttl",
	"turtle-syntax-datatypes-01.ttl",
	"turtle-syntax-datatypes-02.ttl",
	"turtle-syntax-file-01.ttl",
	"turtle-syntax-file-02.ttl",
	"turtle-syntax-file-03.ttl",
	"turtle-syntax-kw-01.ttl",
	"turtle-syntax-kw-02.ttl",
	"turtle-syntax-kw-03.ttl",
	"turtle-syntax-lists-01.ttl",
	"turtle-syntax-lists-02.ttl",
	"turtle-syntax-lists-03.ttl",
	"turtle-syntax-lists-04.ttl",
	"turtle-syntax-lists-05.ttl",
	"turtle-syntax-ln-colons.ttl",
	"turtle-syntax-ln-dots.ttl",
	"turtle-syntax-ns-dots.ttl",
	"turtle-syntax-number-01.ttl",
	"turtle-syntax-number-02.ttl",
	"turtle-syntax-number-03.ttl",
	"turtle-syntax-number-04.ttl",
	"turtle-syntax-number-05.ttl",
	"turtle-syntax-number-06.ttl",
	"turtle-syntax-number-07.ttl",
	"turtle-syntax-number-08.ttl",
	"turtle-syntax-number-09.ttl",
	"turtle-syntax-number-10.ttl",
	"turtle-syntax-number-11.ttl",
	"turtle-syntax-pname-esc-01.ttl",
	"turtle-syntax-pname-esc-02.ttl",
	"turtle-syntax-pname-esc-03.ttl",
	"turtle-syntax-prefix-01.ttl",
	"turtle-syntax-prefix-02.ttl",
	"turtle-syntax-prefix-03.ttl",
	"turtle-syntax-prefix-04.ttl",
	"turtle-syntax-prefix-05.ttl",
	"turtle-syntax-prefix-06.ttl",
	"turtle-syntax-prefix-07.ttl",
	"turtle-syntax-prefix-08.ttl",
	"turtle-syntax-prefix-09.ttl",
	"turtle-syntax-str-esc-01.ttl",
	"turtle-syntax-str-esc-02.ttl",
	"turtle-syntax-str-esc-03.ttl",
	"turtle-syntax-string-01.ttl",
	"turtle-syntax-string-02.ttl",
	"turtle-syntax-string-03.ttl",
	"turtle-syntax-string-04.ttl",
	"turtle-syntax-string-05.ttl",
	"turtle-syntax-string-06.ttl",
	"turtle-syntax-string-07.ttl",
	"turtle-syntax-string-08.ttl",
	"turtle-syntax-string-09.ttl",
	"turtle-syntax-string-10.ttl",
	"turtle-syntax-string-11.ttl",
	"turtle-syntax-struct-01.ttl",
	"turtle-syntax-struct-02.ttl",
	"turtle-syntax-struct-03.ttl",
	"turtle-syntax-struct-04.ttl",
	"turtle-syntax-struct-05.ttl",
	"turtle-syntax-uri-01.ttl",
	"turtle-syntax-uri-02.ttl",
	"turtle-syntax-uri-03.ttl",
	"turtle-syntax-uri-04.ttl",
}
