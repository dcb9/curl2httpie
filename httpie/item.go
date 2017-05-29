package httpie

import "fmt"

type Item struct {
	K string
	V string
	S string
}

func (i *Item) String() string {
	return fmt.Sprintf(`"%s%s%s"`, i.K, i.S, i.V)
}

func NewHeader(key, val string) *Item {
	return &Item{key, val, ":"}
}

func NewURLParam(key, val string) *Item {
	return &Item{key, val, "=="}
}

func NewDataField(key, val string) *Item {
	return &Item{key, val, "="}
}

func NewJSONField(key, val string) *Item {
	return &Item{key, val, ":="}
}

func NewFileField(key, val string) *Item {
	return &Item{key, val, "@"}
}
