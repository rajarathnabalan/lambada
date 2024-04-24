package lambada

import (
	"net/url"
)

// toURLValues converts a simple map into an url.Values
func toURLValues(v map[string]string) url.Values {
	res := make(url.Values)
	for k, v := range v {
		newV, err := url.QueryUnescape(v)
		if err == nil {
			v = newV
		}
		res.Set(k, v)
	}
	return res
}
