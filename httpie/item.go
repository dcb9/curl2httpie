package httpie

import (
	"fmt"
)

type Itemer interface {
	Item() string
}

type QueryItemer interface {
	ToQuery() string
}

type Item struct {
	K string
	V string
	S string
}

type QueryItem Item


func (i *Item) Item() string {
	return fmt.Sprintf(`"%s%s%s"`, i.K, i.S, i.V)
}

func (i *QueryItem) Item() string {
	return fmt.Sprintf(`"%s%s%s"`, i.K, i.S, i.V)
}

func (i *QueryItem) ToQuery() string {
	return fmt.Sprintf(`%s=%s`, i.K, i.V)
}

func NewHeader(key, val string) *Item {
	return &Item{key, val, ":"}
}

func NewURLParam(key, val string) *QueryItem {
	return &QueryItem{key, val, "=="}
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
