package gon3

type stateFn func(*charMatchLexer) stateFn

const (
	eof = -1
)

func lexDocument(l *charMatchLexer) stateFn {
	l.acceptRun(runWhitespace)
	l.ignore()
	switch l.peek() {
	case '@':
		return lexAtStatement
	case '_':
		return lexBlankNodeLabel
	case '<':
		return lexIRIRef
	case '\'':
		// TODO: return proper statefn
	case '"':
		// TODO: return proper statefn
	case '[', ']', '(', ')', ';', ',', '.':
		return lexPunctuation
	case 't', 'f', 'a':
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
	panic("unreachable")
}

func lexAtStatement(l *charMatchLexer) stateFn {
	// lex prefix/base directives or langtag
	if !l.accept("@") {
		return l.errorf("lexAtStatement called, but '@' not found")
	}
	if !l.acceptRun(runAlphabet) {
		return l.errorf("Expected Alphabet while lexing '@...', got %s", l.input[l.pos-1:l.pos])
	}
	if l.accept("-") {
		if !l.acceptRun(runAlphanumeric) {
			return l.errorf("Expected Alphanumeric while lexing langtag, got %s", l.input[l.pos-1:l.pos])
		}
		l.emit(tokenLangTag)
		return lexDocument
	}
	if l.atPrefix() {
		l.emit(tokenAtPrefix)
		return lexDocument
	}
	if l.atBase() {
		l.emit(tokenAtBase)
		return lexDocument
	}
	l.emit(tokenLangTag)
	return lexDocument
}

func lexBlankNodeLabel(l *charMatchLexer) stateFn {

	newMatcher()
		.acceptRunes("_")
		.matchOne()

	newMatcher()
		.acceptRunes(":")
		.matchOne()

	newMatcher()
		.union(mPNCharsU)
		.acceptRunes("0123456789")
		.matchOne()

	bNodeLabelMatcher := &sequentialMatcher{
		[]matcher{
			&stringMatcher{false, "_"},
			&stringMatcher{false, ":"},
			&unionMatcher{
				[]matcher{
					//PN_CHARS_U
					&stringMatcher{false, runDigits},
				},
			},
		},
	}

	if l.acceptRun(runPNChars + ".") {
		for l.acceptRun(runPNChars + ".") {
		}
		l.backup()
		if !l.accept(runPNChars) {
			return l.errorf("Expected PNChars for last char of bnode label, got %s", l.input[l.pos-1:l.pos])
		}
	}
	l.emit(tokenBlankNodeLabel)
	return lexDocument
}

func lexIRIRef(l *charMatchLexer) stateFn {
	// TODO: implement
	panic("unimplemented")
}

func lexPunctuation(l *charMatchLexer) stateFn {
	// TODO: implement
	panic("unimplemented")
}

func lexBooleanLiteral(l *charMatchLexer) stateFn {
	// TODO: implement
	panic("unimplemented")
}

func lexPName(l *charMatchLexer) stateFn {
	// accept PN_PREFIX
	if l.acceptRun(runPNCharsBase) {
		if l.acceptRun(runPNChars + ".") {
			for l.acceptRun(runPNChars + ".") {
			}
			l.backup()
			if !l.accept(runPNChars) {
				return l.errorf("Expected PNChars for last char of prefix, got %s", l.input[l.pos-1:l.pos])
			}
		}
	}
	if !l.accept(":") {
		return l.errorf("Expected ':' in pname, got %s", l.input[l.pos-1:l.pos])
	}
	if l.atWhitespace() {
		l.emit(tokenPNameNS)
		return lexDocument
	}
	// TODO: accept PN_LOCAL
	panic("unfinished")
}

func (l *charMatchLexer) atPrefix() bool {
	// TODO: implement
	panic("unimplemented")
}

func (l *charMatchLexer) atBase() bool {
	// TODO: implement
	panic("unimplemented")
}

func (l *charMatchLexer) atFalse() bool {
	// TODO: implement
	panic("unimplemented")
}

func (l *charMatchLexer) atTrue() bool {
	// TODO: implement
	panic("unimplemented")
}

func (l *charMatchLexer) atA() bool {
	// TODO: implement
	panic("unimplemented")
}

func (l *charMatchLexer) atWhitespace() bool {
	// TODO: implement
	panic("unimplemented")
}
