package connector

import "testing"

func TestCurl2Httpie(t *testing.T) {
	cases := []struct {
		in   []string
		want string
	}{
		{
			[]string{"curl", "z.cn"},
			"http GET z.cn",
		}, {
			[]string{"curl", "--url", "local.dev"},
			"http GET local.dev",
		}, {
			[]string{"curl", "--user", "user:password", "local.dev"},
			`http --auth "user:password" GET local.dev`,
		}, {
			[]string{"curl", "local.dev", "--user", "user:password"},
			`http --auth "user:password" GET local.dev`,
		}, {
			[]string{"curl", "-u", "user:password", "local.dev"},
			`http --auth "user:password" GET local.dev`,
		}, {
			[]string{"curl", "-u", "user:password", "-X", "POST", "local.dev"},
			`http --auth "user:password" --form POST local.dev`,
		}, {
			[]string{"curl", "-u", "user:password", "-X", "POST", "-H", "Foo: bar", "--data", "foo=bar", "local.dev"},
			`http --auth "user:password" --form POST local.dev "Foo:bar" "foo=bar"`,
		}, {
			[]string{"curl", "-u", "user:password", "-X", "POST", "-H", "Content-Type: application/json", "--data", "foo=bar", "local.dev"},
			`http --auth "user:password" --json POST local.dev "foo=bar"`,
		}, {
			[]string{"curl", "-u", "user:password", "--request", "GET", "local.dev", "-H", "accept: application/json", "-H", "Authorization: jwtToken"},
			`http --auth "user:password" --json GET local.dev "Authorization:jwtToken"`,
		}, {
			[]string{"curl", "-u", "user:", "local.dev"},
			`http --auth "user:" GET local.dev`,
		}, {
			[]string{"curl", "local.dev", "-u", "user:", "-d", "foo=Bar"},
			`http --auth "user:" --form POST local.dev "foo=Bar"`,
		}, {
			[]string{"curl", "local.dev", "--user-agent", "httpie", "-u", "user:", "-d", "foo=Bar"},
			`http --auth "user:" --form POST local.dev "User-Agent:httpie" "foo=Bar"`,
		}, {
			[]string{"curl", "local.dev", "--referer", "z.cn", "--user-agent", "httpie", "-u", "user:", "-d", "foo=Bar"},
			`http --auth "user:" --form POST local.dev "Referer:z.cn" "User-Agent:httpie" "foo=Bar"`,
		}, {
			[]string{"curl", "local.dev", "-v", "--referer", "z.cn", "--user-agent", "httpie", "-u", "user:", "-d", "foo=Bar"},
			`http --verbose --auth "user:" --form POST local.dev "Referer:z.cn" "User-Agent:httpie" "foo=Bar"`,
		}, {
			[]string{"curl", "local.dev", "--verbose", "--referer", "z.cn", "--user-agent", "httpie", "-u", "user:", "-d", "foo=Bar"},
			`http --verbose --auth "user:" --form POST local.dev "Referer:z.cn" "User-Agent:httpie" "foo=Bar"`,
		}, {
			[]string{"curl", "local.dev", "--verbose", "--cookie", "NAME=VAL", "--referer", "z.cn", "--user-agent", "httpie", "-u", "user:", "-d", "foo=Bar"},
			`http --verbose --auth "user:" --form POST local.dev "Cookie:NAME=VAL" "Referer:z.cn" "User-Agent:httpie" "foo=Bar"`,
		},
	}

	for _, c := range cases {
		got := Curl2Httpie(c.in)
		want := c.want
		if got != want {
			t.Fatalf("Curl2Httpie error got: %s, want: %s, in: %p", got, want, c.in)
		}
	}
}
