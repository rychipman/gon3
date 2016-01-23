package gon3

import (
	"easylex"
)

const (
	eof = -1
)

func lexDocument(l *charMatchLexer) stateFn {
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
			l.Emit(tokenTrue)
			return lexDocument
		}
		if matchFalse.MatchOne(l) {
			l.Emit(tokenFalse)
			return lexDocument
		}
		if matchA.MatchOne(l) {
			l.Emit(tokenA)
			return lexDocument
		}
		fallthrough
	default:
		return lexPName
	}
	panic("unreachable")
}

func lexAtStatement(l *charMatchLexer) stateFn {
	easylex.NewMatcher().AcceptRunes("@").MatchOne(l)
	// TODO: assert

	// TODO: implement

	return lexDocument
}

func lexBlankNodeLabel(l *charMatchLexer) stateFn {
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
	l.emit(tokenBlankNodeLabel)
	return lexDocument
}

func lexIRIRef(l *charMatchLexer) stateFn {
	// TODO: implement
	panic("unimplemented")
}

func lexPunctuation(l *charMatchLexer) stateFn {
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

func lexPName(l *charMatchLexer) stateFn {
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
		l.emit(tokenPNameNS)
		return lexDocument
	}
	// accept PN_LOCAL
	if l.Peek() == `\` {
		l.Next()
		matchEscapable.MatchOne(l)
		// TODO: assert
	} else if l.Peek() == "%" {
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
		if l.Peek() == `\` {
			l.Next()
			matchEscapable.MatchOne(l)
			// TODO: assert
			other = true
		} else if l.Peek() == "%" {
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
	l.emit(tokenPNameLN)
	return lexDocument
}
