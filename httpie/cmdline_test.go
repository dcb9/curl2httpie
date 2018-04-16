package httpie

import (
	"testing"
)

func TestNewCmdLineByArgs(t *testing.T) {
	// todo test all flags
	cases := []struct {
		in   []string
		want *CmdLine
	}{
		{
			[]string{"http", "z.cn"},
			&CmdLine{
				URL: "z.cn",
			},
		},
		{
			[]string{"http", "-f", "POST", "z.cn"},
			&CmdLine{
				URL: "z.cn",
				Flags: []*Flag{
					FormFlag,
				},
				Method: NewMethod("POST"),
			},
		},
		{
			[]string{"http", "-f", "--auth", "user:pass", "POST", "z.cn"},
			&CmdLine{
				URL: "z.cn",
				Flags: []*Flag{
					FormFlag,
					AuthFlagWithArg("user:pass"),
				},
				Method: NewMethod("POST"),
			},
		},
		{

			[]string{"http", "z.cn", "X-API-Token:123"},
			&CmdLine{
				URL: "z.cn",
				Items: []*Item{
					NewHeader("X-API-Token", "123"),
				},
			},
		},
		{

			[]string{"http", "z.cn", `X-API-Token\::123`},
			&CmdLine{
				URL: "z.cn",
				Items: []*Item{
					NewHeader(`X-API-Token\:`, "123"),
				},
			},
		},
		{

			[]string{"http", "POST", "example.org", "foo\\==bar"},
			&CmdLine{
				Method: NewMethod("POST"),
				URL:    "example.org",
				Items: []*Item{
					NewDataField("foo\\=", "bar"),
				},
			},
		},
		{

			[]string{"http", "POST", "example.org", "foo=bar"},
			&CmdLine{
				Method: NewMethod("POST"),
				URL:    "example.org",
				Items: []*Item{
					NewDataField("foo", "bar"),
				},
			},
		},
		{

			[]string{"http", "--form", "PUT", "example.org", "X-API-Token:123", "name=John"},
			&CmdLine{
				Method: NewMethod("PUT"),
				URL:    "example.org",
				Items: []*Item{
					NewHeader("X-API-Token", "123"),
					NewDataField("name", "John"),
				},
				Flags: []*Flag{
					FormFlag,
				},
			},
		},
		{

			[]string{"http", "--form", "PUT", "example.org", "X-API-Token:123", "name=John", "id==1"},
			&CmdLine{
				Method: NewMethod("PUT"),
				URL:    "example.org",
				Items: []*Item{
					NewHeader("X-API-Token", "123"),
					NewDataField("name", "John"),
					NewURLParam("id", "1"),
				},
				Flags: []*Flag{
					FormFlag,
				},
			},
		},
		{

			[]string{"http", "--form", "PUT", "example.org", "X-API-Token:123", `foo\==bar`, "id==1"},
			&CmdLine{
				Method: NewMethod("PUT"),
				URL:    "example.org",
				Items: []*Item{
					NewHeader("X-API-Token", "123"),
					NewDataField("foo\\=", "bar"),
					NewURLParam("id", "1"),
				},
				Flags: []*Flag{
					FormFlag,
				},
			},
		},
	}

	for _, c := range cases {
		got, err := NewCmdLineByArgs(c.in[1:])
		if err != nil {
			t.Fatalf("NewCmdLineByArgs error: %s in: %#v", err.Error(), c.in)
		}

		want := c.want
		if got.String() != want.String() {
			t.Errorf("NewCmdLineByArgs error got: %s, want: %s, in: %#v", got.String(), want.String(), c.in)
		}
	}
}
