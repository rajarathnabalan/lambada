package lambada

import "net/http"

func toSingleValueHeaders(h http.Header) map[string]string {
	res := map[string]string{}
	for k, v := range h {
		if len(v) > 0 {
			res[k] = v[0]
		}
	}
	return res
}

func fromSingleValueHeaders(h map[string]string) http.Header {
	res := make(http.Header)
	for k, v := range h {
		res.Set(k, v)
	}
	return res
}

func canonicalizeHeader(h http.Header) http.Header {
	res := make(http.Header)
	for k, v := range h {
		res[http.CanonicalHeaderKey(k)] = v
	}
	return res
}
