package curl

import "fmt"

func NewMethod(method string) *Option {
	return &Option{
		Short: 'X',
		Long: LongName("request"),
		HasArg: true,
		Arg: method,
	}
}

func NewUser(user string) *Option {
	return &Option{
		Short: 'u',
		Long: LongName("user"),
		HasArg: true,
		Arg: user,
	}
}

func NewHeader(key, val string) *Option {
	return &Option{
		Short: 'H',
		Long: LongName("header"),
		HasArg: true,
		Arg: fmt.Sprintf("%s: %v", key, val),
	}
}

func NewJSONHeader() *Option {
	return NewHeader("Content-Type", "application/json")
}

func NewForm(content string) *Option {
	return &Option{
		Short: 'F',
		Long: LongName("form"),
		HasArg: true,
		Arg: content,
	}
}

func NewRequest(request string) *Option {
	return &Option{
		Short: 'X',
		Long: LongName("request"),
		HasArg: true,
		Arg: request,
	}
}

func NewData(data string) *Option {
	return &Option{
		Short: 'd',
		Long: LongName("data"),
		HasArg: true,
		Arg: data,
	}
}
