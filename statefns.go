package gon3

import (
	"fmt"
	"github.com/rychipman/easylex"
	"strings"
)

const (
	eof = -1
)

func lexDocument(l *easylex.Lexer) easylex.StateFn {
	matchWhitespace.MatchRun(l)
	l.Ignore()
	switch l.Peek() {
	case easylex.EOF:
		l.Emit(easylex.TokenEOF)
		return nil
	case '#':
		return lexComment
	case '@':
		return lexAtStatement
	case '_':
		return lexBlankNodeLabel
	case '<':
		return lexIRIRef
	case '"', '\'':
		return lexRDFLiteral
	case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// TODO: handle decimal vs period
		return lexNumericLiteral
	case '^', '[', ']', '(', ')', ';', ',', '.':
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

func lexComment(l *easylex.Lexer) easylex.StateFn {
	easylex.NewMatcher().AcceptRunes("#").MatchOne(l)
	// TODO: assert
	easylex.NewMatcher().RejectRunes("\n").MatchRun(l)
	l.Ignore()
	return lexDocument
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

func lexRDFLiteral(l *easylex.Lexer) easylex.StateFn {
	if matchLongQuote.MatchOne(l) {
		fmt.Printf("hi")
		for {
			if matchLongQuote.MatchOne(l) {
				break
			}
			q := matchQuote.MatchOne(l)
			q = matchQuote.MatchOne(l) || q
			ch := true
			if easylex.NewMatcher().RejectRunes(`"\`).MatchOne(l) {
				// do nothing
			} else if l.Peek() == '\\' {
				l.Next()
				switch l.Peek() {
				case 'u':
					l.Next()
					for i := 0; i < 4; i += 1 {
						matchHex.MatchOne(l)
						// TODO: assert
					}
				case 'U':
					l.Next()
					for i := 0; i < 8; i += 1 {
						matchHex.MatchOne(l)
						// TODO: assert
					}
				case 't', 'b', 'n', 'r', 'f', '"', '\'', '\\':
					l.Next()
				default:
					ch = false
				}
			}
			if q && !ch {
				// TODO: error
			}
		}
		l.Emit(tokenStringLiteralLongQuote)
		return lexDocument
	}
	if matchLongSingleQuote.MatchOne(l) {
		for {
			if matchLongSingleQuote.MatchOne(l) {
				break
			}
			q := matchSingleQuote.MatchOne(l)
			q = matchSingleQuote.MatchOne(l) || q
			ch := true
			if easylex.NewMatcher().RejectRunes(`'\`).MatchOne(l) {
				// do nothing
			} else if l.Peek() == '\\' {
				l.Next()
				switch l.Peek() {
				case 'u':
					l.Next()
					for i := 0; i < 4; i += 1 {
						matchHex.MatchOne(l)
						// TODO: assert
					}
				case 'U':
					l.Next()
					for i := 0; i < 8; i += 1 {
						matchHex.MatchOne(l)
						// TODO: assert
					}
				case 't', 'b', 'n', 'r', 'f', '"', '\'', '\\':
					l.Next()
				default:
					ch = false
				}
			}
			if !q && !ch {
				break
			}
		}
		l.Emit(tokenStringLiteralLongSingleQuote)
		return lexDocument
	}
	if matchQuote.MatchOne(l) {
		for {
			if easylex.NewMatcher().RejectRunes("\u0022\u005c\u000a\u000d").MatchOne(l) {
				// do nothing
			} else if l.Peek() == '\\' {
				l.Next()
				switch l.Peek() {
				case 'u':
					l.Next()
					for i := 0; i < 4; i += 1 {
						matchHex.MatchOne(l)
						// TODO: assert
					}
				case 'U':
					l.Next()
					for i := 0; i < 8; i += 1 {
						matchHex.MatchOne(l)
						// TODO: assert
					}
				case 't', 'b', 'n', 'r', 'f', '"', '\'', '\\':
					l.Next()
				default:
					break
				}
			} else {
				break
			}
		}
		matchQuote.MatchOne(l)
		// TODO: assert
		l.Emit(tokenStringLiteralQuote)
		return lexDocument
	}
	if matchSingleQuote.MatchOne(l) {
		for {
			if easylex.NewMatcher().RejectRunes("\u0027\u005c\u000a\u000d").MatchOne(l) {
				// do nothing
			} else if l.Peek() == '\\' {
				l.Next()
				switch l.Peek() {
				case 'u':
					l.Next()
					for i := 0; i < 4; i += 1 {
						matchHex.MatchOne(l)
						// TODO: assert
					}
				case 'U':
					l.Next()
					for i := 0; i < 8; i += 1 {
						matchHex.MatchOne(l)
						// TODO: assert
					}
				case 't', 'b', 'n', 'r', 'f', '"', '\'', '\\':
					l.Next()
				default:
					break
				}
			} else {
				break
			}
		}
		matchSingleQuote.MatchOne(l)
		// TODO: assert
		l.Emit(tokenStringLiteralSingleQuote)
		return lexDocument
	}
	panic("unreachable")
}

func lexNumericLiteral(l *easylex.Lexer) easylex.StateFn {
	easylex.NewMatcher().AcceptRunes("+-").MatchOne(l)
	if matchDigits.MatchRun(l) {
		if easylex.NewMatcher().AcceptRunes("eE").MatchOne(l) {
			easylex.NewMatcher().AcceptRunes("+-").MatchOne(l)
			matchDigits.MatchRun(l)
			// TODO: assert
			l.Emit(tokenDouble)
			return lexDocument
		} else if matchPeriod.MatchOne(l) {
			if matchDigits.MatchRun(l) {
				if isWhitespace(l.Peek()) {
					l.Emit(tokenDecimal)
					return lexDocument
				}
			}
			easylex.NewMatcher().AcceptRunes("eE").MatchOne(l)
			// TODO: assert
			easylex.NewMatcher().AcceptRunes("+-").MatchOne(l)
			matchDigits.MatchRun(l)
			// TODO: assert
			l.Emit(tokenDouble)
			return lexDocument
		} else {
			l.Emit(tokenInteger)
			return lexDocument
		}
	} else {
		matchPeriod.MatchOne(l)
		// TODO: assert
		matchDigits.MatchRun(l)
		// TODO: assert
		if easylex.NewMatcher().AcceptRunes("eE").MatchOne(l) {
			easylex.NewMatcher().AcceptRunes("+-").MatchOne(l)
			matchDigits.MatchRun(l)
			// TODO: assert
			l.Emit(tokenDouble)
			return lexDocument
		}
		l.Emit(tokenDecimal)
		return lexDocument
	}
}

func lexPunctuation(l *easylex.Lexer) easylex.StateFn {
	// ^ [ ] ( ) ; , .
	if matchDoubleCaret.MatchOne(l) {
		l.Emit(tokenLiteralDatatypeTag)
	} else if matchOpenBracket.MatchOne(l) {
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
	if l.Peek() == '<' || isWhitespace(l.Peek()) {
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
		} else if l.Peek() == ':' {
			l.Next()
			other = true
		} else {
			other = matchPNChars.MatchRun(l)
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
