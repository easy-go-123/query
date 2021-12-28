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

// UrlJoin .
// nolint: revive
func UrlJoin(a, b string, noQuery bool) (string, error) {
	a = strings.Trim(a, "\r\n \t")
	b = strings.Trim(b, "\r\n \t")

	if a == "" {
		return "", errFail
	}

	if a[len(a)-1:] != "/" {
		a += "/"
	}

	if noQuery {
		if strings.Contains(a, "&") {
			return "", errFail
		}
	}

	if b == "" {
		return a, nil
	}

	if b[0:1] == "/" {
		b = b[1:]
	}

	if noQuery {
		if strings.Contains(b, "&") {
			return "", errFail
		}
	}

	return a + b, nil
}

// UrlAddQueryString .
// nolint: revive
func UrlAddQueryString(u, k, v string) string {
	vals := url.Values{}
	vals.Set(k, v)

	return UrlAddQuery(u, vals)
}

// UrlAddRawQuery .
// nolint: revive
func UrlAddRawQuery(u, rawQuery string) string {
	rawQuery = strings.TrimPrefix(rawQuery, "?")

	idx := strings.Index(u, "?")
	if idx == -1 {
		return u + "?" + rawQuery
	}

	if idx == len(u)-1 {
		return u + rawQuery
	}

	return u + "&" + rawQuery
}

// UrlAddQuery .
// nolint: revive
func UrlAddQuery(u string, vals url.Values) string {
	return UrlAddRawQuery(u, vals.Encode())
}

// UrlUpdateQueryString 注意：不处理一个key对应多个value的情况.
// nolint: revive
func UrlUpdateQueryString(u, k, v string) (string, error) {
	k = url.QueryEscape(k)
	v = url.QueryEscape(v)

	qp := strings.Index(u, "?")
	if qp == -1 || qp == len(u)-1 {
		return UrlAddQueryString(u, k, v), nil
	}

	idx := strings.Index(u[qp:], k+"=")
	if idx == -1 {
		return UrlAddQueryString(u, k, v), nil
	}

	idx2 := strings.Index(u[qp+idx+1:], k+"=")
	if idx2 != -1 {
		return "", errFail
	}

	idxE := strings.Index(u[qp+idx:], "&")

	vals := url.Values{}
	vals.Set(k, v)

	r := u[:qp+idx] + vals.Encode()

	if idxE != -1 {
		r += u[qp+idxE+1:]
	}

	return r, nil
}
