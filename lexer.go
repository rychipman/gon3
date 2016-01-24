package gon3

type lexer interface {
	NextToken easylex.Token
}

type mockLexer struct {
	tokens []easylex.Token
	pos int
}

func newMockLexer(args ...easylex.Token) *mockLexer {
	return &mockLexer{
		tokens: args
		pos:0,
	}
}

func (m *mockLexer) NextToken() easylex.Token {
	ret := m.tokens[m.pos]
	m.pos += 1
	return ret
}
