package gon3

import (
	"github.com/rychipman/easylex"
	"strings"
)

const (
	eof = -1
)

func lexDocument(l *easylex.Lexer) easylex.StateFn {
	matchWhitespace.MatchRun(l)
	switch l.Peek() {
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
		if matchTrue.MatchOne(l) {
			if isWhitespace(l.Peek()) {
				l.Emit(tokenTrue)
				return lexDocument
			}
		}
		if matchFalse.MatchOne(l) {
			if isWhitespace(l.Peek()) {
				l.Emit(tokenFalse)
				return lexDocument
			}
		}
		if matchA.MatchOne(l) {
			if isWhitespace(l.Peek()) {
				l.Emit(tokenA)
				return lexDocument
			}
		}
		fallthrough
	default:
		return lexPName
	}
	panic("unreachable")
}

func isWhitespace(r rune) bool {
	if strings.IndexRune("\n\r\t\v\f ", r) >= 0 {
		return true
	}
	return false
}

func lexAtStatement(l *easylex.Lexer) easylex.StateFn {
	easylex.NewMatcher().AcceptRunes("@").MatchOne(l)
	// TODO: assert
	if easylex.NewMatcher().AcceptString("prefix").MatchOne(l) {
		if isWhitespace(l.Peek()) {
			l.Emit(tokenAtPrefix)
			return lexDocument
		}
	}
	if easylex.NewMatcher().AcceptString("base").MatchOne(l) {
		if isWhitespace(l.Peek()) {
			l.Emit(tokenAtBase)
			return lexDocument
		}
	}
	matchAlphabet.MatchRun(l)
	// TODO: assert
	for {
		hyphen := easylex.NewMatcher().AcceptRunes("-").MatchOne(l)
		alph := matchAlphaNumeric.MatchRun(l)
		if !hyphen && !alph {
			break
		}
		if hyphen != alph {
			// TODO: error
		}
	}
	l.Emit(tokenLangTag)
	return lexDocument
}

func lexBlankNodeLabel(l *easylex.Lexer) easylex.StateFn {
	easylex.NewMatcher().AcceptRunes("_").MatchOne(l)
	// TODO: assert
	easylex.NewMatcher().AcceptRunes(":").MatchOne(l)
	// TODO: assert
	easylex.NewMatcher().Union(matchPNCharsU).Union(matchDigits).MatchOne(l) // TODO: create these matchers
	// TODO: assert
	for {
		period := matchPeriod.MatchRun(l)
		pnchars := matchPNChars.MatchRun(l)
		if !pnchars {
			if period {
				// TODO: error
			}
			break
		}
	}
	l.Emit(tokenBlankNodeLabel)
	return lexDocument
}

func lexIRIRef(l *easylex.Lexer) easylex.StateFn {
	easylex.NewMatcher().AcceptRunes("<").MatchOne(l)
	// TODO: assert
	iriChars := easylex.NewMatcher().RejectRunes("<>\"{}|^`\\\u0000\u0001\u0002\u0003\u0004\u0005\u0006\u0007\u0008\u0009\u000a\u000b\u000c\u000d\u000e\u000f\u0010\u0011\u0012\u0013\u0014\u0015\u0016\u0017\u0018\u0019\u001a\u001b\u001c\u001d\u001e\u001f\u0020")
	for {
		m1 := iriChars.MatchRun(l)
		if l.Peek() == '\\' {
			l.Next()
			if l.Peek() == 'u' {
				for i := 0; i < 4; i += 1 {
					matchHex.MatchOne(l)
					// TODO: assert
				}
			} else if l.Peek() == 'U' {
				for i := 0; i < 8; i += 1 {
					matchHex.MatchOne(l)
					// TODO: assert
				}
			} else {
				// TODO: error
			}
		} else if !m1 {
			break
		}
	}
	easylex.NewMatcher().AcceptRunes(">").MatchOne(l)
	// TODO: assert
	l.Emit(tokenIRIRef)
	return lexDocument
}

func lexPunctuation(l *easylex.Lexer) easylex.StateFn {
	// [ ] ( ) ; , .
	if matchOpenBracket.MatchOne(l) {
		l.Emit(tokenStartBlankNodePropertyList)
	} else if matchCloseBracket.MatchOne(l) {
		l.Emit(tokenEndBlankNodePropertyList)
	} else if matchOpenParens.MatchOne(l) {
		l.Emit(tokenStartCollection)
	} else if matchCloseParens.MatchOne(l) {
		l.Emit(tokenEndCollection)
	} else if matchSemicolon.MatchOne(l) {
		l.Emit(tokenPredicateListSeparator)
	} else if matchComma.MatchOne(l) {
		l.Emit(tokenObjectListSeparator)
	} else if matchPeriod.MatchOne(l) {
		l.Emit(tokenEndTriples)
	} else {
		// TODO: error
	}
	return lexDocument
}

func lexPName(l *easylex.Lexer) easylex.StateFn {
	// accept PN_PREFIX
	matchPNCharsBase.MatchOne(l)
	for {
		period := matchPeriod.MatchRun(l)
		pnchars := matchPNChars.MatchRun(l)
		if !pnchars {
			if period {
				// TODO: error
			}
			break
		}
	}
	easylex.NewMatcher().AcceptRunes(":").MatchOne(l)
	// TODO: assert
	if matchWhitespace.MatchRun(l) {
		l.Emit(tokenPNameNS)
		return lexDocument
	}
	// accept PN_LOCAL
	if l.Peek() == '\\' {
		l.Next()
		matchEscapable.MatchOne(l)
		// TODO: assert
	} else if l.Peek() == '%' {
		l.Next()
		matchHex.MatchOne(l)
		// TODO: assert
		matchHex.MatchOne(l)
		// TODO: assert
	} else {
		easylex.NewMatcher().AcceptRunes(":").Union(matchPNCharsU).Union(matchDigits).MatchOne(l)
		// TODO: assert
	}
	for {
		period := matchPeriod.MatchRun(l)
		other := false
		if l.Peek() == '\\' {
			l.Next()
			matchEscapable.MatchOne(l)
			// TODO: assert
			other = true
		} else if l.Peek() == '%' {
			l.Next()
			matchHex.MatchOne(l)
			// TODO: assert
			matchHex.MatchOne(l)
			// TODO: assert
			other = true
		} else {
			easylex.NewMatcher().AcceptRunes(":").Union(matchPNCharsU).Union(matchDigits).MatchOne(l)
			// TODO: assert
			other = true
		}
		if !other {
			if period {
				// TODO: error
			}
			break
		}
	}
	l.Emit(tokenPNameLN)
	return lexDocument
}
