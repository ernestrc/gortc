package sip

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// URI represents a RFC-3986 SIP URI: sip:user:password@host:port;uri-parameters?headers
type URI struct {
	Scheme   string
	User     string
	Password string
	Host     string
	Port     uint16
	Params   url.Values
	Headers  url.Values
}

func (u *URI) parsePort(raw string) (err error) {
	var port uint64
	if port, err = strconv.ParseUint(raw, 10, 16); err != nil {
		return fmt.Errorf("invalid URI port: %s", raw)
	}
	u.Port = uint16(port)
	return
}

const (
	scheme = iota
	userpass
	hostport
	params
	headers
)

func (u *URI) parseHost(i int, raw string, state, anchor, anchorAux *int, nextState int) (err error) {
	if *anchorAux == *anchor {
		u.Host = raw[*anchor:i]
	} else {
		u.Host = raw[*anchor:*anchorAux]
		if err = u.parsePort(raw[*anchorAux+1 : i]); err != nil {
			return
		}
	}
	*anchor = i + 1
	*state = nextState

	return
}

func (u *URI) next(i int, r byte, raw string, state, anchor, anchorAux *int) (err error) {
	switch *state {
	case scheme:
		switch r {
		case ':':
			u.Scheme = raw[*anchor:i]
			*anchor = i + 1
			*anchorAux = *anchor
			*state++
		case ' ':
			*anchor++
		}
	case userpass:
		switch r {
		case ':':
			*anchorAux = i
		case '@':
			if *anchorAux == *anchor {
				u.User = raw[*anchor:i]
			} else {
				u.User = raw[*anchor:*anchorAux]
				u.Password = raw[*anchorAux+1 : i]
			}
			*anchor = i + 1
			*anchorAux = *anchor
			*state = hostport
		case '?':
			err = u.parseHost(i, raw, state, anchor, anchorAux, headers)
		case ';':
			err = u.parseHost(i, raw, state, anchor, anchorAux, params)
		}
	case hostport:
		switch r {
		case ':':
			*anchorAux = i
		case '?':
			err = u.parseHost(i, raw, state, anchor, anchorAux, headers)
		case ';':
			err = u.parseHost(i, raw, state, anchor, anchorAux, params)
		}
	case params:
		switch r {
		case '?':
			var params url.Values
			if params, err = url.ParseQuery(strings.Replace(raw[*anchor:i], ";", "&", -1)); err != nil {
				return
			}
			u.Params = params
			*anchor = i + 1
			*state++
		}
	case headers:
		var headers url.Values
		if headers, err = url.ParseQuery(raw[*anchor:i]); err != nil {
			return
		}
		u.Headers = headers
	}

	return
}

func (u *URI) checkValid() error {
	switch u.Scheme {
	case "sips":
	case "sip":
	default:
		return fmt.Errorf("invalid URI scheme: %s", u.Scheme)
	}

	return nil
}

func (u *URI) Parse(raw string) (err error) {
	state := scheme
	anchor := 0
	anchorAux := 0

	for i, l := 0, len(raw); i < l && err == nil; i++ {
		err = u.next(i, raw[i], raw, &state, &anchor, &anchorAux)
	}

	if err != nil {
		return
	}

	// force parse last segment
	switch state {
	case params:
		err = u.next(len(raw), '?', raw, &state, &anchor, &anchorAux)
	default:
		err = u.next(len(raw), ';', raw, &state, &anchor, &anchorAux)
	}
	if err != nil {
		return
	}

	err = u.checkValid()
	return
}

func (u *URI) String() string {
	uri := u.Scheme + ":"

	if u.User != "" {
		uri += u.User
		if u.Password != "" {
			uri += ":"
			uri += u.Password
		}
		uri += "@"
	}

	if u.Port != 0 {
		uri += fmt.Sprintf("%s:%d", u.Host, u.Port)
	} else {
		uri += u.Host
	}

	if u.Params != nil {
		uri += ";"
		uri += strings.Replace(u.Params.Encode(), "&", ";", -1)
	}

	if u.Headers != nil {
		uri += "?" + u.Headers.Encode()
	}

	return uri
}
