package gon3

type stateFn func(*charMatchLexer) stateFn

const (
	eof = -1
)

const (
	runAlphabet     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	runDigits       = "0123456789"
	runAlphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	runNumeric      = runDigits + "-."
	runWhitespace   = " \n\t\v\f\r"    // TODO: check this against TR
	runPNCharsBase  = runAlphabet + "" // TODO: complete with unicode escapes
	runPNCharsU     = runPNCharsBase + "_"
	runPNChars      = runPNCharsU + runDigits + "-" // TODO: complete with unicode escapes
	runEscapable    = "_~.-!$&'()*+,;=/?#@%"
)

func lexDocument(l *charMatchLexer) stateFn {
	l.acceptRun(runWhitespace)
	l.ignore()
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
	// lex bnode label
	if !l.accept("_") {
		return l.errorf("lexAtStatement called, but '@' not found")
	}
	if !l.accept(":") {
		return l.errorf("Expected ':' while lexing bnode label, got %s", l.input[l.pos-1:l.pos])
	}
	if !l.acceptRun(runPNCharsU + runDigits) {
		return l.errorf("Expected PNCharsU or Digits while lexing bnode label, got %s", l.input[l.pos-1:l.pos])
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
}

func lexPunctuation(l *charMatchLexer) stateFn {
	// TODO: implement
}

func lexBooleanLiteral(l *charMatchLexer) stateFn {
	// TODO: implement
}

func lexPName(l *charMatchLexer) stateFn {
	// accept PN_PREFIX
	if acceptRun(runPNCharsBase) {
		if l.acceptRun(runPNChars + ".") {
			for l.acceptRun(runPNChars + ".") {
			}
			l.backup()
			if !l.accept(runPNChars) {
				return l.errorf("Expected PNChars for last char of prefix, got %s", l.input[l.pos-1:l.pos])
			}
		}
	}
	if !accept(":") {
		return l.errorf("Expected ':' in pname, got %s", l.input[l.pos-1:l.pos])
	}
	if l.atWhitespace() {
		l.emit(tokenPNameNS)
		return lexDocument
	}
	// TODO: accept PN_LOCAL
}

func (l *charMatchLexer) atPrefix() bool {
	// TODO: implement
}

func (l *charMatchLexer) atBase() bool {
	// TODO: implement
}

func (l *charMatchLexer) atFalse() bool {
	// TODO: implement
}

func (l *charMatchLexer) atTrue() bool {
	// TODO: implement
}

func (l *charMatchLexer) atA() bool {
	// TODO: implement
}

func (l *charMatchLexer) atWhitespace() bool {
	// TODO: implement
}
