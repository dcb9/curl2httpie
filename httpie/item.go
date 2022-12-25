package httpie

import (
	"fmt"

	"github.com/dcb9/curl2httpie/shellwords"
)

const (
	SEP_HEADER    = ":"
	SEP_URL_PARAM = "=="
	SEP_DATA      = "="
	SEP_JSON      = ":="
	SEP_FILE      = "@"
)

type Item struct {
	K string
	V string
	S string
}

func (i *Item) String() string {
	if shellwords.NeedQuote(i.V) {
		return fmt.Sprintf(`'%s%s%s'`, i.K, i.S, i.V)
	}
	return fmt.Sprintf(`%s%s%s`, i.K, i.S, i.V)
}

func NewHeader(key, val string) *Item {
	return &Item{key, val, SEP_HEADER}
}

func NewURLParam(key, val string) *Item {
	return &Item{key, val, SEP_URL_PARAM}
}

func NewDataField(key, val string) *Item {
	return &Item{key, val, SEP_DATA}
}

func NewJSONField(key, val string) *Item {
	return &Item{key, val, SEP_JSON}
}

func NewFileField(key, val string) *Item {
	return &Item{key, val, SEP_FILE}
}
