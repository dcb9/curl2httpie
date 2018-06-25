package connector

import "testing"

func TestCurl2Httpie(t *testing.T) {
	cases := []struct {
		in   []string
		want string
	}{
		{
			[]string{"curl", "-i", "-u", "username", "-d", `'{"scopes":["public_repo"]}'`, `https://api.github.com/authorizations`},
			`echo '{"scopes":["public_repo"]}' | http --auth 'username' --form POST 'https://api.github.com/authorizations'`,
		},
		{
			[]string{"curl", "z.cn"},
			"http GET 'z.cn'",
		}, {
			[]string{"curl", "--url", "local.dev"},
			"http GET 'local.dev'",
		}, {
			[]string{"curl", "--user", "user:password", "local.dev"},
			`http --auth 'user:password' GET 'local.dev'`,
		}, {
			[]string{"curl", "local.dev", "--user", "user:password"},
			`http --auth 'user:password' GET 'local.dev'`,
		}, {
			[]string{"curl", "-u", "user:password", "local.dev"},
			`http --auth 'user:password' GET 'local.dev'`,
		}, {
			[]string{"curl", "-u", "user:password", "-X", "POST", "local.dev"},
			`http --auth 'user:password' --form POST 'local.dev'`,
		}, {
			[]string{"curl", "-u", "user:password", "-X", "POST", "-H", "Foo: bar", "--data", "foo=bar", "local.dev"},
			`http --auth 'user:password' --form POST 'local.dev' 'Foo:bar' 'foo=bar'`,
		}, {
			[]string{"curl", "-u", "user:password", "-X", "POST", "-H", "Content-Type: application/json", "--data", "foo=bar", "local.dev"},
			`http --auth 'user:password' --json POST 'local.dev' 'foo=bar'`,
		}, {
			[]string{"curl", "-u", "user:password", "--request", "GET", "local.dev", "-H", "accept: application/json", "-H", "Authorization: jwtToken"},
			`http --auth 'user:password' --json GET 'local.dev' 'Authorization:jwtToken'`,
		}, {
			[]string{"curl", "-u", "user:", "local.dev"},
			`http --auth 'user:' GET 'local.dev'`,
		}, {
			[]string{"curl", "local.dev", "-u", "user:", "-d", "foo=Bar"},
			`http --auth 'user:' --form POST 'local.dev' 'foo=Bar'`,
		}, {
			[]string{"curl", "local.dev", "--user-agent", "httpie", "-u", "user:", "-d", "foo=Bar"},
			`http --auth 'user:' --form POST 'local.dev' 'User-Agent:httpie' 'foo=Bar'`,
		}, {
			[]string{"curl", "local.dev", "--referer", "z.cn", "--user-agent", "httpie", "-u", "user:", "-d", "foo=Bar"},
			`http --auth 'user:' --form POST 'local.dev' 'Referer:z.cn' 'User-Agent:httpie' 'foo=Bar'`,
		}, {
			[]string{"curl", "local.dev", "-v", "--referer", "z.cn", "--user-agent", "httpie", "-u", "user:", "-d", "foo=Bar"},
			`http --verbose --auth 'user:' --form POST 'local.dev' 'Referer:z.cn' 'User-Agent:httpie' 'foo=Bar'`,
		}, {
			[]string{"curl", "local.dev", "--verbose", "--referer", "z.cn", "--user-agent", "httpie", "-u", "user:", "-d", "foo=Bar"},
			`http --verbose --auth 'user:' --form POST 'local.dev' 'Referer:z.cn' 'User-Agent:httpie' 'foo=Bar'`,
		}, {
			[]string{"curl", "local.dev", "--verbose", "--cookie", "NAME=VAL", "--referer", "z.cn", "--user-agent", "httpie", "-u", "user:", "-d", "foo=Bar"},
			`http --verbose --auth 'user:' --form POST 'local.dev' 'Cookie:NAME=VAL' 'Referer:z.cn' 'User-Agent:httpie' 'foo=Bar'`,
		}, {
			[]string{"curl", "-H", "Host: foo.bar.com", "-H", "Accept: */*", "-H", "User-Agent: debug-MyAppName/ CFNetwork/893.14 Darwin/17.4.0", "-H", "Accept-Language: en-us", "--data", "client_id=foobarfoobarfoobar&client_secret=bazquzbazquz&grant_type=password&password=SomePasswordHere&scope=user&username=first.last%2B1%40domain.com", "--compressed", "https://stage.buildsafely.com/api/oauth/token"},
			`http --form POST 'https://stage.buildsafely.com/api/oauth/token' 'Host:foo.bar.com' 'Accept:*/*' 'User-Agent:debug-MyAppName/ CFNetwork/893.14 Darwin/17.4.0' 'Accept-Language:en-us' 'client_id=foobarfoobarfoobar&client_secret=bazquzbazquz&grant_type=password&password=SomePasswordHere&scope=user&username=first.last%2B1%40domain.com'`,
		},
	}

	for _, c := range cases {
		gotStringer, warningMessages, err := Curl2Httpie(c.in[1:])
		if len(warningMessages) > 0 {
			t.Logf("Curl2Httpie warning messages: %#v in: %#v", warningMessages, c.in)
		}
		if err != nil {
			t.Fatalf("Curl2Httpie error: %s in: %#v", err.Error(), c.in)
			continue
		}

		want := c.want
		if got := gotStringer.String(); got != want {
			t.Errorf("Curl2Httpie error got: %s, want: %s, in: %#v", got, want, c.in)
		}
	}
}
