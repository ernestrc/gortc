package sip

import (
	"net/url"
	"reflect"
	"testing"
)

var cases = []struct {
	raw string
	uri URI
}{
	{
		"sip:ernestrc@unstablebuild.sip.twilio.com;transport=tls", URI{
			Scheme: "sip",
			User:   "ernestrc",
			Host:   "unstablebuild.sip.twilio.com",
			Params: url.Values{"transport": []string{"tls"}},
		},
	},
	{
		"sip:unstablebuild.sip.twilio.com", URI{
			Scheme: "sip",
			Host:   "unstablebuild.sip.twilio.com",
		},
	},
	{
		"sip:unstablebuild.com?ok=sir", URI{
			Scheme: "sip",
			Host:   "unstablebuild.com",
			Headers: url.Values{
				"ok": []string{"sir"},
			},
		},
	},
	{
		"sip:50.23.130.56:5061", URI{
			Scheme: "sip",
			Host:   "50.23.130.56",
			Port:   5061,
		},
	},
	{
		"sips:ernestrc:secret@50.10.23.34:5060;transport=tls", URI{
			Scheme:   "sips",
			User:     "ernestrc",
			Password: "secret",
			Host:     "50.10.23.34",
			Port:     5060,
			Params:   url.Values{"transport": []string{"tls"}},
		},
	},
	{
		"sips:50.10.23.34:5060;transport=tls;transport=udp;user=ernestrc", URI{
			Scheme: "sips",
			Host:   "50.10.23.34",
			Port:   5060,
			Params: url.Values{
				"transport": []string{"tls", "udp"},
				"user":      []string{"ernestrc"},
			},
		},
	},
	{
		"sip:unstable.build:5061?transport=tls&transport=udp&user=ernestrc", URI{
			Scheme: "sip",
			Host:   "unstable.build",
			Port:   5061,
			Headers: url.Values{
				"transport": []string{"tls", "udp"},
				"user":      []string{"ernestrc"},
			},
		},
	},
}

func testOK(t *testing.T, expected, uri URI) {
	if expected.Scheme != uri.Scheme {
		t.Errorf("expected scheme %s found %s", expected.Scheme, uri.Scheme)
	}

	if expected.User != uri.User {
		t.Errorf("expected user %+v found %+v", expected.User, uri.User)
	}

	if expected.Host != uri.Host {
		t.Errorf("expected host %s found %s", expected.Host, uri.Host)
	}

	if expected.Port != uri.Port {
		t.Errorf("expected port %d found %d", expected.Port, uri.Port)
	}

	if !reflect.DeepEqual(expected.Params, uri.Params) {
		t.Errorf("expected Params %v found %v", expected.Params, uri.Params)
	}

	if !reflect.DeepEqual(expected.Headers, uri.Headers) {
		t.Errorf("expected Headers %v found %v", expected.Headers, uri.Headers)
	}
}

func TestParse(t *testing.T) {
	for _, tcase := range cases {
		uri := URI{}
		uri.Parse(tcase.raw)
		testOK(t, tcase.uri, uri)
	}
}

var bad = []struct {
	raw string
	uri URI
}{
	{
		"sip:ernicles:copernicles@1234.com;;?a=b;;", URI{
			Scheme:   "sip",
			Host:     "1234.com",
			User:     "ernicles",
			Password: "copernicles",
			Params:   url.Values{},
			Headers: url.Values{
				"a": []string{"b"},
			},
		},
	},
}

func TestParseTricky(t *testing.T) {
	for _, tcase := range bad {
		uri := URI{}
		uri.Parse(tcase.raw)
		testOK(t, tcase.uri, uri)
	}
}

func TestString(t *testing.T) {
	for _, tcase := range cases {
		if s := tcase.uri.String(); s != tcase.raw {
			t.Errorf("found %s; expected %s", s, tcase.raw)
		}
	}
}

func BenchmarkString(b *testing.B) {
	uri := cases[0]
	for i := 0; i < b.N; i++ {
		uri.uri.String()
	}
}

func BenchmarkParse(b *testing.B) {
	raw := cases[len(cases)-1].raw
	uri := URI{}
	for i := 0; i < b.N; i++ {
		uri.Parse(raw)
	}
}
