package shellwords

// Copied from https://github.com/lmika/shellwords/blob/ce258dd729fe9965ab11b0200b9db427ec9fa49d/shellwords.go
/*
Package shellwords splits strings into tokens using similar tokenizing rules as most shells.
This package is based on the Ruby Shellwords package.
A string is made up of tokens separated by whitespace characters (space, tabs and new-lines).
A token can also be made up using string literals delimited by either the single or double
quote characters.
Example:
    package main
    import "fmt"
    import "shellwords"
    func main() {
        tok1 := shellwords.Split("this 'is a' test")    // ["this", "is a", "test"]
    }
*/

import (
	"bytes"
	"strings"
	"text/scanner"
)

// Splitter is used to split an input string into tokens.
type Splitter struct {
	scanner *scanner.Scanner
	buffer  *bytes.Buffer
	hasNext bool
}

// Split is a conveinence method for splitting a string into tokens.  It is equivalent to the code:
//      New(str).Split()
func Split(str string) []string {
	return New(str).Split()
}

// New creates a new splitter for a given string.
func New(str string) *Splitter {
	s := new(scanner.Scanner)
	s.Init(strings.NewReader(str))
	sp := &Splitter{s, new(bytes.Buffer), true}
	sp.scanNext()
	return sp
}

// Next returns the next token from the splitter.  This returns the token and whether or not
// the token was successfully scanned.  If the end of the string is encountered, the second
// return value will be false.
func (s *Splitter) Next() (token string, hasToken bool) {
	token, hasToken = s.buffer.String(), s.hasNext
	s.scanNext()
	return
}

// Scans for the next token
func (s *Splitter) scanNext() {
	s.buffer.Reset()

	if !s.hasNext {
		return
	}

	// Skip the whitespace
	for s.isRuneWhitespace(s.scanner.Peek()) {
		s.scanner.Next()
	}

	if s.scanner.Peek() == scanner.EOF {
		s.hasNext = false
		return
	}

	// Parse the token
	var currStringDelim rune

	c := s.scanner.Next()
	for c != scanner.EOF {
		if s.isRuneWhitespace(c) && (currStringDelim == 0) {
			s.hasNext = true
			return
		} else if s.isRuneStringDelim(c) && (currStringDelim == 0) {
			currStringDelim = c
		} else if s.isRuneStringDelim(c) && (currStringDelim == c) {
			currStringDelim = 0
		} else {
			s.buffer.WriteRune(c)
		}
		c = s.scanner.Next()
	}
}

// Split scans until the end of the string and return the remaining tokens as a slice.  If no
// more tokens are present, this returns an empty slice.
func (s *Splitter) Split() []string {
	tokens := make([]string, 0)

	for token, hasToken := s.Next(); hasToken; token, hasToken = s.Next() {
		tokens = append(tokens, token)
	}

	return tokens
}

// Returns true if the rune is a whitespace
func (s *Splitter) isRuneWhitespace(c rune) bool {
	return (c == ' ') || (c == '\t') || (c == '\r') || (c == '\n')
}

// Returns true if the rune starts a string
func (s *Splitter) isRuneStringDelim(c rune) bool {
	return (c == '\'') || (c == '"')
}

// Returns true if the rune is an escape character
func (s *Splitter) isRuneEscapeChar(c rune) bool {
	return (c == '\\')
}
