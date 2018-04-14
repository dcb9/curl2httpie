package curl

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
