package httpie

import "strings"

type Method string

func (m Method) String() string {
	return string(m)
}

func NewMethod(method string) *Method {
	if method == "" {
		method = "GET"
	}

	m := Method(strings.ToUpper(method))
	return &m
}
