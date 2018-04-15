package connector

import "testing"

func TestHttpie2Curl(t *testing.T) {
	cases := []struct {
		in   []string
		want string
	}{
		{
			[]string{"http", "example.org", "id==1"},
			`curl example.org?id=1`,
		},
		{
			[]string{"http", "--auth", "username", "example.org", "id==1"},
			`curl --user "username" example.org?id=1`,
		},
		{
			[]string{"http", "--auth", "username", "example.org", "id==1", "foo:bar"},
			`curl --user "username" --header "foo: bar" example.org?id=1`,
		},
		{
			[]string{"http", "--form", "--auth", "username", "example.org", "id==1", "foo:bar", "foo=bar"},
			`curl --user "username" --header "foo: bar" --data "foo=bar" example.org?id=1`,
		},
		{
			[]string{"http", "--auth", "username", "example.org", "id==1", "foo:bar", "foo=bar"},
			`curl --user "username" --header "foo: bar" --header "Content-Type: application/json" --data "{\"foo\":\"bar\"}" example.org?id=1`,
		},
		{
			[]string{"http", "-f", "--auth", "username", "example.org", "id==1", "foo:bar", "foo=bar", "file@test_obj.json"},
			`curl --user "username" --header "foo: bar" --form "file=@\"test_obj.json\"" --data "foo=bar" example.org?id=1`,
		},
		{
			[]string{"http", "--auth", "username", "example.org", "id==1", "foo:bar", "foo=bar", `a:={"foo": "bar"}`},
			`curl --user "username" --header "foo: bar" --header "Content-Type: application/json" --data "{\"a\":{\"foo\":\"bar\"},\"foo\":\"bar\"}" example.org?id=1`,
		},
		{
			[]string{"http", "--auth", "username", "POST", "example.org", "id==1", "foo:bar", "foo=bar", `a:={"foo": "bar"}`},
			`curl --request "POST" --user "username" --header "foo: bar" --header "Content-Type: application/json" --data "{\"a\":{\"foo\":\"bar\"},\"foo\":\"bar\"}" example.org?id=1`,
		},
		{
			[]string{"http","PUT", "z.cn"},
			`curl --request "PUT" z.cn`,
		},
		{
			[]string{"http", "z.cn"},
			"curl z.cn",
		},
	}

	for _, c := range cases {
		gotStringer, warningMessages, err := Httpie2Curl(c.in[1:])
		if len(warningMessages) > 0 {
			t.Logf("Httpie2Curl warning messages: %#v in: %#v", warningMessages, c.in)
		}
		if err != nil {
			t.Fatalf("Httpie2Curl error: %s in: %#v", err.Error(), c.in)
			continue
		}

		want := c.want
		if got := gotStringer.String(); got != want {
			t.Errorf("Httpie2Curl error got: %s, want: %s, in: %#v", got, want, c.in)
		}
	}
}

