package gon3

type lexer interface {
	nextToken() token
}

func newCharMatchLexer(name, input string) *charMatchLexer {
	l := &charMatchLexer{
		name:   name,
		input:  input,
		state:  lexDocument,
		tokens: make(chan token, 2),
	}
	return l
}

func newNaiveRegexLexer(name, input string) *naiveRegexLexer {
	l := &naiveRegexLexer{
		name:  name,
		input: input,
	}
	return l
}
