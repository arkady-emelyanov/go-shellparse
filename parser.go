package shellparse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type parser struct {
	input     string
	vars      map[string]string
	pos       int
	start     int
	width     int
	lastError error
	collected []string
}

const (
	runeEOF           = -1
	runeVarStart      = '$'
	runeVarLeftDelim  = '{'
	runeVarRightDelim = '}'
	runeEscape        = '\\'
	runeComment       = '#'
)

// step back in input steam
// can be used only once per next call
func (p *parser) backup() {
	p.pos -= p.width
}

// discard currently processing string
func (p *parser) ignore() {
	p.start = p.pos
}

// lookup next rune in input without adjusting
// current position.
func (p *parser) peek() rune {
	if p.pos >= len(p.input) {
		return runeEOF
	}
	r, _ := utf8.DecodeRuneInString(p.input[p.pos:])
	return r
}

// fetch next rune from input and update parser
// internal position.
func (p *parser) next() rune {
	if p.pos >= len(p.input) {
		p.width = 0
		return runeEOF
	}

	r, w := utf8.DecodeRuneInString(p.input[p.pos:])
	p.width = w
	p.pos += p.width
	return r
}

// collect currently processing string
func (p *parser) collect() {
	p.collected = append(p.collected, p.input[p.start:p.pos])
}

// instead of collect, extract is used to exclude
// currently processing string. It's like ignore
// but returns string back.
func (p *parser) extract() string {
	s := p.input[p.start:p.pos]
	p.ignore()
	return s
}

// once extracted and processed, string can
// be injected as parsed.
func (p *parser) insert(s string) {
	p.collected = append(p.collected, s)
}

// calling this method remembers error and break FSM processing
func (p *parser) errorf(format string, args ...interface{}) stateFn {
	p.lastError = fmt.Errorf(format, args...)
	return nil
}

type stateFn func(p *parser) stateFn

// FSM for splitting string into words, considering quotes,
// variables, escape chars and friends.
//
// Primary usage is to parse and decode string into []string
// suitable for passing to exec.Cmd.
//
func splitWordsFsm(input string) ([]string, error) {
	p := &parser{input: input}
	for state := splitWordsInitial; state != nil; {
		state = state(p)
	}

	err := p.lastError
	return p.collected, err
}

// FSM for replacing ${VAR} with actual content provided as map
func replaceVarsFsm(input string, vars map[string]string) (string, error) {
	p := &parser{input: input, vars: vars}
	for state := replaceVarsInitial; state != nil; {
		state = state(p)
	}

	col := strings.Join(p.collected, "")
	err := p.lastError
	return col, err
}

// FSM for de-escaping words, remove unnecessary escape chars and
// unnecessary quotes from parsed string.
//
// some examples:
// * `\\'it was escaped\\'` -> `it was escaped`
// * `it\'s ok` -> `it's ok`
//
func unescapeWordsFsm(input string) (string, error) {
	p := &parser{input: input}
	for state := unescapeVarsInitial; state != nil; {
		state = state(p)
	}

	col := strings.Join(p.collected, "")
	err := p.lastError
	return col, err
}

// initial state for unescape fsm
func unescapeVarsInitial(p *parser) stateFn {
	var startQuote rune

	first := true
Loop:
	for {
		r := p.next()

		switch {
		case r == runeEOF:
			p.collect() // collect everything left in buffer
			break Loop

		case r == runeEscape:
			next := p.peek()
			// all escape sequences should be unescaped,
			// except some special chars..
			if next == runeEscape || isQuote(next) || isStringPart(next) {
				p.backup()  // shift escape char back
				p.collect() // collect everything before escape char
				p.next()    // step forward to escape char
				p.ignore()  // ignore escape char
			}

		case isQuote(r) && p.pos >= 1 && first == true:
			startQuote = r
			first = false
			p.ignore() // skip first quote

		case isStringPart(r):
			first = false

		case r == startQuote:
			if p.peek() == runeEOF {
				p.backup()  // shift quote back
				p.collect() // collect everything before quote
				break Loop  // stop processing
			}
		}
	}

	return nil
}

// initial state for var replacement fsm
func replaceVarsInitial(p *parser) stateFn {
	var prev rune

	for {
		r := p.next()
		if r == runeEOF {
			p.collect() // collect everything left in buffer
			break
		}

		// `\${` escaped
		if r == runeVarStart && prev != runeEscape {
			if p.peek() == runeVarLeftDelim {
				p.backup()  // shift $ char back
				p.collect() // collect everything before $ char
				p.next()    // step forward to $ char
				p.next()    // step forward to { char
				p.ignore()  // ignore both
				return replaceVar
			}
		}

		prev = r
	}

	return nil
}

// replace ${VAR} with actual content
func replaceVar(p *parser) stateFn {
	for {
		r := p.next()
		if r == runeEOF {
			return p.errorf("unexpected EOF while looking for: %#U", runeVarRightDelim)
		}

		if r == runeVarRightDelim {
			p.backup()
			break
		}
	}

	k := p.extract()
	v, ok := p.vars[k]
	if !ok {
		return p.errorf("unknown variable: %#v", k)
	}

	p.insert(v) // insert directly without altering pos/start
	p.next()    // step forward for end right delimiter
	p.ignore()  // skip it

	return replaceVarsInitial
}

// initial state for word splitting fsm
func splitWordsInitial(p *parser) stateFn {
	switch r := p.next(); {
	case r == runeEOF:
		break

	case r == runeEscape:
		next := p.peek()
		if next == runeEscape || isQuote(next) { // possible situations: \', \\, \"
			p.next() // since next rune is escaped, treating it as part of string
		}
		return splitWordsInitial

	case r == runeComment:
		return ignoreLine

	case isSpace(r) || isEndOfLine(r):
		p.ignore()
		return splitWordsInitial

	case isQuote(r):
		p.backup() // quote should be part of string
		return collectQuotedString

	case isStringPart(r):
		return collectString

	default:
		return p.errorf("unexpected char: %#U (%d)", r, p.pos)
	}

	return nil
}

func ignoreLine(p *parser) stateFn {
	for {
		r := p.next()
		if isEndOfLine(r) || r == runeEOF {
			p.ignore()
			break
		}
	}
	return splitWordsInitial
}

// parse string, string is everything between non string chars
func collectString(p *parser) stateFn {
	var prev rune

	for {
		r := p.next()
		if !isStringPart(r) && prev == runeEscape {
			// everything after escape rune is skipped
			continue
		}

		if !isStringPart(r) || isSpace(r) || isEndOfLine(r) || r == runeEOF {
			p.backup()
			p.collect()
			break
		}

		prev = r
	}
	return splitWordsInitial
}

// parse string enclosed with '"' or "'"
// for situations when delimiter nested into another quote: 'echo "D'oh"',
// nested quote should be quoted: 'echo "D\'oh"'
func collectQuotedString(p *parser) stateFn {
	var previousRune rune

	quote := p.next() // get quote
	for {
		r := p.next()

		// expression: `\\` should be ignored
		if r == runeEscape && p.peek() == runeEscape {
			p.next()
			continue
		}

		// expression: `"it\"s"`, emit as `it"s`
		if r == quote && previousRune == runeEscape {
			previousRune = r
			continue
		}

		// expression: `"it's"`, emit as `it's`
		if isQuote(r) && r != quote {
			previousRune = r
			continue
		}

		// final quote found
		if r == quote {
			p.collect() // collect string
			break
		}

		if r == runeEOF {
			return p.errorf("unexpected EOF while looking for enclosing: %#U (%d)", quote, p.start)
		}

		previousRune = r
	}

	return splitWordsInitial
}

// chars treated as quote
func isQuote(r rune) bool {
	return r == '"' || r == '\''
}

// chars treated as space
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// chars treated as EOL marker
func isEndOfLine(r rune) bool {
	return r == '\n'
}

// chars allowed as string contents
func isStringPart(r rune) bool {
	return !isEndOfLine(r) && !isSpace(r) && !isQuote(r) && unicode.IsPrint(r)
}
