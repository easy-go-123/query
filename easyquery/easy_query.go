package easyquery

import (
	"errors"
	"net/url"
	"strings"

	"github.com/easy-go-123/query/query"
)

var errFail = errors.New("fail")

func ParseQuery(query string) (url.Values, error) {
	return url.ParseQuery(query)
}

func EncodeValues(values url.Values) string {
	return values.Encode()
}

func Values2Orm(values url.Values, s interface{}) error {
	return query.Unmarshal(values, s)
}

func Orm2Values(s interface{}, values url.Values) error {
	return query.Marshal(s, values)
}

// nolint: revive
func UrlJoin(a, b string) (string, error) {
	a = strings.Trim(a, "\r\n \t")
	b = strings.Trim(b, "\r\n \t")

	if a == "" {
		return "", errFail
	}

	if a[len(a)-1:] != "/" {
		a += "/"
	}

	if b == "" {
		return a, nil
	}

	if b[0:1] == "/" {
		b = b[1:]
	}

	return a + b, nil
}

// nolint: revive
func UrlAddQueryString(u, k, v string) string {
	vals := url.Values{}
	vals.Set(k, v)

	idx := strings.Index(u, "?")
	if idx == -1 {
		return u + "?" + vals.Encode()
	}

	if idx == len(u)-1 {
		return u + vals.Encode()
	}

	return u + "&" + vals.Encode()
}

// nolint: revive
func UrlUpdateQueryString(u, k, v string) string {
	k = url.QueryEscape(k)
	v = url.QueryEscape(v)

	qp := strings.Index(u, "?")
	if qp == -1 || qp == len(u)-1 {
		return UrlAddQueryString(u, k, v)
	}

	idx := strings.Index(u[qp:], k+"=")
	if idx == -1 {
		return UrlAddQueryString(u, k, v)
	}

	idxE := strings.Index(u[qp+idx:], "&")

	vals := url.Values{}
	vals.Set(k, v)

	r := u[:qp+idx] + vals.Encode()

	if idxE != -1 {
		r += u[qp+idxE+1:]
	}

	return r
}
