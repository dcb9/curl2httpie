package connector

import "testing"

func TestHttpie2Curl(t *testing.T) {
	cases := []struct {
		in   []string
		want string
	}{
		{
			[]string{"http", ":/foo"},
			`curl localhost/foo`,
		},
		{
			[]string{"http", ":3000/bar"},
			`curl localhost:3000/bar`,
		},
		{
			[]string{"http", ":"},
			`curl localhost/`,
		},
		{
			[]string{"http", "example.org", "id==1"},
			`curl example.org?id=1`,
		},
		{
			[]string{"http", "--auth", "username", "example.org", "id==1"},
			`curl --user username example.org?id=1`,
		},
		{
			[]string{"http", "--auth", "username", "example.org", "id==1", "foo:bar"},
			`curl --user username --header 'foo: bar' example.org?id=1`,
		},
		{
			[]string{"http", "--form", "--auth", "username", "example.org", "id==1", "foo:bar", "foo=bar"},
			`curl --user username --header 'foo: bar' --data foo=bar example.org?id=1`,
		},
		{
			[]string{"http", "--auth", "username", "example.org", "id==1", "foo:bar", "foo=bar"},
			`curl --user username --header 'foo: bar' --header 'Content-Type: application/json' --data '{"foo":"bar"}' example.org?id=1`,
		},
		{
			[]string{"http", "-f", "--auth", "username", "example.org", "id==1", "foo:bar", "foo=bar", "file@test_obj.json"},
			`curl --user username --header 'foo: bar' --form 'file=@"test_obj.json"' --data foo=bar example.org?id=1`,
		},
		{
			[]string{"http", "--auth", "username", "example.org", "id==1", "foo:bar", "foo=bar", `a:={"foo": "bar"}`},
			`curl --user username --header 'foo: bar' --header 'Content-Type: application/json' --data '{"a":{"foo":"bar"},"foo":"bar"}' example.org?id=1`,
		},
		{
			[]string{"http", "--auth", "username", "POST", "example.org", "id==1", "foo:bar", "foo=bar", `a:={"foo": "bar"}`},
			`curl --request POST --user username --header 'foo: bar' --header 'Content-Type: application/json' --data '{"a":{"foo":"bar"},"foo":"bar"}' example.org?id=1`,
		},
		{
			[]string{"http", "PUT", "z.cn"},
			`curl --request PUT z.cn`,
		},
		{
			[]string{"http", "z.cn"},
			"curl z.cn",
		},
		{
			[]string{"http", "--auth", "username", "--auth-type", "basic", "example.org", "id==1"},
			`curl --user username --basic example.org?id=1`,
		},
		{
			[]string{"http", "--auth", "username", "--auth-type", "digest", "example.org", "id==1"},
			`curl --user username --digest example.org?id=1`,
		},
		{
			[]string{"http", "--auth", "username", "--auth-type", "digest", "--proxy", "http:http://foo.bar:3128", "example.org", "id==1"},
			`curl --user username --digest --proxy http:http://foo.bar:3128 example.org?id=1`,
		},
		{
			[]string{"http", "--auth", "username", "--auth-type", "digest", "--proxy", "http:http://foo.bar:3128", "example.org", "id==1"},
			`curl --user username --digest --proxy http:http://foo.bar:3128 example.org?id=1`,
		},
		{
			[]string{"http", "--auth", "username", "--auth-type", "digest", "--proxy", "http:http://foo.bar:3128", "--follow", "example.org", "id==1"},
			`curl --user username --digest --proxy http:http://foo.bar:3128 --location example.org?id=1`,
		},
		{
			[]string{"http", "--auth", "username", "--auth-type", "digest", "--proxy", "http:http://foo.bar:3128", "--follow", "--max-redirects", "10", "example.org", "id==1"},
			`curl --user username --digest --proxy http:http://foo.bar:3128 --location --max-redirs 10 example.org?id=1`,
		},
		{
			[]string{"http", "--auth", "username", "--auth-type", "digest", "--proxy", "http:http://foo.bar:3128", "--follow", "--max-redirects", "10", "--timeout", "30", "example.org", "id==1"},
			`curl --user username --digest --proxy http:http://foo.bar:3128 --location --max-redirs 10 --max-time 30 example.org?id=1`,
		},
		{
			[]string{"https", "--auth", "username", "--auth-type", "digest", "--proxy", "http:http://foo.bar:3128", "--follow", "--max-redirects", "10", "--timeout", "30", "example.org", "id==1"},
			`curl --user username --digest --proxy http:http://foo.bar:3128 --location --max-redirs 10 --max-time 30 https://example.org?id=1`,
		},
		{
			[]string{"https", "pie.dev"},
			`curl https://pie.dev`,
		},
		{
			[]string{"http", "pie.dev"},
			`curl pie.dev`,
		},
		{
			[]string{"https", "pie.dev", "key==mykey", "secret==mysecret"},
			`curl 'https://pie.dev?key=mykey&secret=mysecret'`,
		},
		{
			[]string{"http", "-a", "username:password", "pie.dev"},
			`curl --user username:password pie.dev`,
		},
		{
			[]string{"http", "pie.dev", "-a", "username:password"},
			`curl --user username:password pie.dev`,
		},
	}

	// cases = []struct {
	// 	in   []string
	// 	want string
	// }{
	// 	{
	// 		[]string{"http", "pie.dev", "-a", "username:password"},
	// 		`curl --user username:password pie.dev`,
	// 	},
	// }

	for _, c := range cases {
		gotStringer, warningMessages, err := Httpie2Curl(c.in)
		if len(warningMessages) > 0 {
			t.Logf("Httpie2Curl warning messages: %#v in: %#v", warningMessages, c.in)
		}
		if err != nil {
			t.Fatalf("Httpie2Curl error: %s in: %#v", err.Error(), c.in)
			continue
		}

		want := c.want
		if got := gotStringer.String(); got != want {
			t.Errorf("Httpie2Curl error\ngot:\n%s\nwant:\n%s\nin:\n%#v", got, want, c.in)
		}
	}
}
