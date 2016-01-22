package gon3

type naiveRegexLexer struct {
	name  string
	input string
	start int
	pos   int
}

func (l *naiveRegexLexer) nextToken() token {
	l.skipWhitespace()
	switch l.peek() {
	case "@":
		return lexAtStatement
	case "_":
		return lexBlankNodeLabel
	case "<":
		return lexIRIRef
	case "'":
		// TODO: return proper statefn
	case "\"":
		// TODO: return proper statefn
	case "[", "]", "(", ")", ";", ",", ".":
		return lexPunctuation
	case "t", "f", "a":
		if l.atTrue() || l.atFalse() {
			return lexBooleanLiteral
		}
		if l.atA() {
			l.next()
			l.emit(tokenA)
			return lexDocument
		}
		fallthrough
	default:
		return lexPName
	}
}

func (l *naiveRegexLexer) peek() rune {
	if l.pos >= len(l.input) {
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *naiveRegexLexer) skipWhitespace() {
	// TODO: implement
}
