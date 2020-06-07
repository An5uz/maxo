
package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)



type ItemKind int

const (
	WhiteSpace ItemKind = iota
	LineBreak
	NewLine
	Text
	Digit
	StringQuote
	SingleLineComment
	MultiLineCommentStart
	MultiLineCommentEnd
	Operator
	EOF
	ItemError
)

// item is accumulated while lexing the provided input, and emitted over a
// channel to the parser. Items could also be called tokens as we tokenize the
// input.
type Item struct {
	Position int

	// kind signals how we've classified the data we have accumulated while
	// scanning the string.
	Kind ItemKind

	// value is the segment of data we've accumulated.
	Value string
}

const eof = rune(0)

// stateFn is a function that is specific to a state within the string.
type stateFn func(*Lexer) stateFn

// lex creates a lexer and starts scanning the provided input.
func Lex(input string) *Lexer {
	l := &Lexer{
		input: input,
		state: lexText,
		Items: make(chan Item, 1),
	}

	go l.scan()

	return l
}

// lexer is created to manage an individual scanning/parsing operation.
type Lexer struct {
	input    string    // we'll store the string being parsed
	start    int       // the position we started scanning
	Position int       // the current position of our scan
	width    int       // we'll be using runes which can be double byte
	state    stateFn   // the current state function
	Items    chan Item // the channel we'll use to communicate between the lexer and the parser
}

// emit sends a item over the channel so the parser can collect and manage
// each segment.
func (l *Lexer) emit(k ItemKind) {
	accumulation := l.input[l.start:l.Position]

	i := Item{
		Position: l.start,
		Kind:     k,
		Value:    accumulation,
	}

	l.Items <- i

	l.ignore() // reset our scanner now that we've dispatched a segment
}

// nextItem pulls an item from the lexer's result channel.
func (l *Lexer) nextItem() Item {
	return <-l.Items
}

// ignore resets the start position to the current scan position effectively
// ignoring any input.
func (l *Lexer) ignore() {
	l.start = l.Position
}

// next advances the lexer state to the next rune.
func (l *Lexer) next() (r rune) {
	if l.Position >= len(l.input) {
		l.width = 0
		return eof
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.Position:])
	l.Position += l.width
	return r
}

// backup allows us to step back one run1e which is helpful when you've crossed
// a boundary from one state to another.
func (l *Lexer) backup() {
	l.Position = l.Position - 1
}

// scan will step through the provided text and execute state functions as
// state changes are observed in the provided input.
func (l *Lexer) scan() {
	// When we begin processing, let's assume we're going to process text.
	// One state function will return another until `nil` is returned to signal
	// the end of our process.
	for fn := lexText; fn != nil; {
		fn = fn(l)
	}

	close(l.Items)
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	msg := fmt.Sprintf(format, args...)
	l.Items <- Item{
		Kind:  ItemError,
		Value: msg,
	}

	return nil
}

// lexEOF emits the accumulated data classified by the provided itemKind and
// signals that we've reached the end of our lexing by returning `nil` instead
// of a state function.
func (l *Lexer) lexEOF(k ItemKind) stateFn {

	//	l.backup()
	if l.start > l.Position {
		l.ignore()
	}

	l.emit(k)
	l.emit(EOF)
	return nil
}

// lexText scans what is expected to be text.
func lexText(l *Lexer) stateFn {
	for {
		r := l.next()
		switch {
		case r == eof:
			return l.lexEOF(Text)
		case unicode.IsSpace(r):
			l.backup()

			// emit any text we've accumulated.
			if l.Position > l.start {
				l.emit(Text)
			}
			return lexWhitespace
		}
	}
}

// lexWhitespace scans what is expected to be whitespace.
func lexWhitespace(l *Lexer) stateFn {
	for {
		r := l.next()
		switch {
		case r == eof:
			return l.lexEOF(WhiteSpace)
		case !unicode.IsSpace(r):
			l.backup()
			if l.Position > l.start {
				l.emit(WhiteSpace)
			}
			return lexText
		}
	}
}
// TODO add lexItemKind stuff that is still missing


// ParseSimple is a simple routine to preserve whitespace while reversing the
// characters in words.
func SimpleParser(input string) string {
	var result string
	var word string
	for _, char := range input {
		c := string(char)
		if c == " " {
			// Clean-up the accumulated word
			if len(word) > 0 {
				result += Reverse(word)
			}
			result += " "
			continue
		}
	}

	if len(word) > 0 {
		result += Reverse(word)
	}

	return result
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}