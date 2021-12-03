package easyquery

import (
	"net/url"

	"github.com/easy-go-123/query/query"
)

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
