package gon3

import (
	"strings"
)

type matcher interface {
	match(*charMatchLexer) bool
}

type stringMatcher string

func (s stringMatcher) match(l *charMatchLexer) bool {
	strn := s.(string)
	success := false
	for strings.IndexRune(strn, l.next()) >= 0 {
		if !success {
			success = true
		}
	}
	l.backup()
	return success
}

type unicodeRangeMatcher struct {
	first rune
	last  rune
}

func (u *unicodeRangeMatcher) match(l *charMatchLexer) bool {
	success := false
	for {
		next := l.next()
		if next < u.first || next > u.last {
			break
		}
		if !success {
			success = true
		}
	}
	l.backup()
	return success
}

type unionMatcher struct {
	matchers []matcher
}

func (u *unionMatcher) match(l *charMatchLexer) bool {
	success := false
	for {
		progress := false
		for _, m := range u.matchers {
			progress = m.match(l) || progress
		}
		if !progress {
			break
		} else if !success {
			success = true
		}
	}
	return success
}

type sequentialMatcher struct {
	matchers []matcher
}

func (s *sequentialMatcher) match(l *charMatchLexer) bool {
	success := false
	for _, m := range s.matchers {
		success = m.match(l) || success
	}
	return success
}

var (
	runAlphabet     = stringMatcher("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	runDigits       = stringMatcher("0123456789")
	runAlphanumeric = stringMatcher("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	runWhitespace   = stringMatcher("\u0020\u0009\u000D\u000A") // space, tab, cr, newline

	u00c0_u00d6 = &unicodeRangeMatcher{rune(0x00c0), rune(0x00d6)}
	u00d8_u00f6 = &unicodeRangeMatcher{rune(0x00d8), rune(0x00f6)}
	u00f8_u02ff = &unicodeRangeMatcher{rune(0x00f8), rune(0x02ff)}
	u0370_u037d = &unicodeRangeMatcher{rune(0x0370), rune(0x037d)}
	u037f_u1fff = &unicodeRangeMatcher{rune(0x037f), rune(0x1fff)}
	u200c_u200d = &unicodeRangeMatcher{rune(0x200c), rune(0x200d)}
	u2070_u218f = &unicodeRangeMatcher{rune(0x2070), rune(0x218f)}
	u2c00_u2fef = &unicodeRangeMatcher{rune(0x2c00), rune(0x2fef)}
	u3001_ud7ff = &unicodeRangeMatcher{rune(0x3001), rune(0xd7ff)}
	uf900_ufdcf = &unicodeRangeMatcher{rune(0xf900), rune(0xfdcf)}
	ufdf0_ufffd = &unicodeRangeMatcher{rune(0xfdf0), rune(0xfffd)}

	runPNCharsBase = &unionMatcher{[]matcher{runAlphabet, u00c0_u00d6, u00d8_u00f6, u00f8_u02ff, u0370_u037d, u037f_u1fff, u200c_u200d, u2070_u218f, u2c00_u2fef, u3001_ud7ff, uf900_ufdcf, ufdf0_ufffd}}
	// TODO: 10000 - EFFFF in PN_CHARS_BASE

	runPNCharsU  = &unionMatcher{[]matcher{runPNCharsBase, stringMatcher("_")}}
	runPNChars   = &unionMatcher{[]matcher{runPNCharsU, runDigits, stringMatcher("-\u00B7\u203F\u2040"), &unicodeRangeMatcher{rune(0x0300), rune(0x036f)}}}
	runEscapable = stringMatcher("_~.-!$&'()*+,;=/?#@%")
)
