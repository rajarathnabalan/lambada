package agw

import (
	"net/url"
)

func ToURLValues(v map[string]string) url.Values {
	res := make(url.Values)
	for k, v := range v {
		res.Set(k, v)
	}
	return res
}
