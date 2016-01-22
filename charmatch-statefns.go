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
	runWhitespace   = "\u0020\u0009\u000D\u000A" // space, tab, cr, newline

	runPNCharsBase = runAlphabet + u00c0_u00d6 + u00d8_u00f6 + u00f8_u02ff + u0370_u037d + u037f_u1fff + u200c_u200d + u2070_u218f + u2c00_u2fef + u3001_ud7ff + uf900_ufdcf + ufdf0_ufffd // TODO
	// TODO: 10000 - EFFFF in PN_CHARS_BASE
	runPNCharsU  = runPNCharsBase + "_"
	runPNChars   = runPNCharsU + runDigits + "-" + "\u00B7" + "\u0300\u0301\u0302\u0303\u0304\u0305\u0306\u0307\u0308\u0309\u030A\u030B\u030C\u030D\u030E\u030F\u0310\u0311\u0312\u0313\u0314\u0315\u0316\u0317\u0318\u0319\u031A\u031B\u031C\u031D\u031E\u031F\u0320\u0321\u0322\u0323\u0324\u0325\u0326\u0327\u0328\u0329\u032A\u032B\u032C\u032D\u032E\u032F\u0330\u0331\u0332\u0333\u0334\u0335\u0336\u0337\u0338\u0339\u033A\u033B\u033C\u033D\u033E\u033F\u0340\u0341\u0342\u0343\u0344\u0345\u0346\u0347\u0348\u0349\u034A\u034B\u034C\u034D\u034E\u034F\u0350\u0351\u0352\u0353\u0354\u0355\u0356\u0357\u0358\u0359\u035A\u035B\u035C\u035D\u035E\u035F\u0360\u0361\u0362\u0363\u0364\u0365\u0366\u0367\u0368\u0369\u036A\u036B\u036C\u036D\u036E\u036F" + "\u2034\u2040"
	runEscapable = "_~.-!$&'()*+,;=/?#@%"
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
