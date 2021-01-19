package lambada

import (
	"net/url"
)

func toURLValues(v map[string]string) url.Values {
	res := make(url.Values)
	for k, v := range v {
		res.Set(k, v)
	}
	return res
}
