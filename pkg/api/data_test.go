package api

import "testing"

func TestRPCOpts_EndpointURL(t *testing.T) {
	tests := []struct {
		name string
		opts RPCOpts
		want string
	}{
		{
			name: "no params",
			opts: RPCOpts{BaseEndpoint: "/foo/bar"},
			want: "/foo/bar",
		},
		{
			name: "path params only",
			opts: RPCOpts{
				BaseEndpoint:   "/foo/bar",
				PathParameters: []string{"abc"},
			},
			want: "/foo/bar/abc",
		},
		{
			name: "query params only",
			opts: RPCOpts{
				BaseEndpoint:    "/foo/bar",
				QueryParameters: map[string]string{"type": "tls-auth"},
			},
			want: "/foo/bar?type=tls-auth",
		},
		{
			name: "path and query params",
			opts: RPCOpts{
				BaseEndpoint:    "/foo/bar",
				PathParameters:  []string{"abc"},
				QueryParameters: map[string]string{"flag": "1"},
			},
			want: "/foo/bar/abc?flag=1",
		},
		{
			name: "multiple query params are sorted",
			opts: RPCOpts{
				BaseEndpoint: "/foo/bar",
				QueryParameters: map[string]string{
					"z": "last",
					"a": "first",
				},
			},
			want: "/foo/bar?a=first&z=last",
		},
		{
			name: "values are url-escaped",
			opts: RPCOpts{
				BaseEndpoint:    "/foo/bar",
				QueryParameters: map[string]string{"q": "a b&c"},
			},
			want: "/foo/bar?q=a+b%26c",
		},
		{
			name: "existing query string is preserved with &",
			opts: RPCOpts{
				BaseEndpoint:    "/foo/bar?existing=1",
				QueryParameters: map[string]string{"added": "2"},
			},
			want: "/foo/bar?existing=1&added=2",
		},
		{
			name: "empty query map adds no separator",
			opts: RPCOpts{
				BaseEndpoint:    "/foo/bar",
				QueryParameters: map[string]string{},
			},
			want: "/foo/bar",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.opts.EndpointURL(); got != tc.want {
				t.Fatalf("EndpointURL() = %q, want %q", got, tc.want)
			}
		})
	}
}
