package lambada

import "net/http"

// toSingleValueHeaders converts the headers to single value headers.
// If any header is multi-valued, only the first value is retained.
func toSingleValueHeaders(h http.Header) map[string]string {
	res := map[string]string{}
	for k, v := range h {
		if len(v) > 0 {
			res[k] = v[0]
		}
	}
	return res
}

// fromSingleValueHeaders returns a http.Header from a map of single-valued headers.
// Header keys are canonicalized (using textproto.CanonicalMIMEHeaderKey) during the copy.
func fromSingleValueHeaders(h map[string]string) http.Header {
	res := make(http.Header)
	for k, v := range h {
		res.Set(k, v)
	}
	return res
}

// canonicalizeHeader returns a copy of h with all the keys canonicalized using http.CanonicalHeaderKey
func canonicalizeHeader(h http.Header) http.Header {
	res := make(http.Header)
	for k, v := range h {
		res[http.CanonicalHeaderKey(k)] = v
	}
	return res
}
