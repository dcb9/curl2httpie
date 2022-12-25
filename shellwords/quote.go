package shellwords

import (
	"fmt"
	"strings"
)

func NeedQuote(s string) bool {
	return -1 != strings.IndexFunc(s, func(r rune) bool {
		switch r {
		case '&', '@', '#', '[', ']', '{', '}', ' ', '(', ')', '*':
			return true
		}
		return false
	})
}

func AddQuoteIfNeeded(s string) string {
	if NeedQuote(s) {
		return fmt.Sprintf("'%s'", s)
	}
	return fmt.Sprintf("%s", s)
}
