package lambada

import "encoding/base64"

// bodyToBytes converts the API Gateway request body into a byte slice.
// If isBase64 is true, the body is first decoded from Base64.
func bodyToBytes(body string, isBase64 bool) ([]byte, error) {
	if isBase64 {
		return base64.StdEncoding.DecodeString(body)
	}
	return ([]byte)(body), nil
}

// bytesToBody converts a byte slice to a API Gateway response body.
// If isBase64 is true, the bytes are encoded to Base64.
func bytesToBody(bytes []byte, isBase64 bool) string {
	if isBase64 {
		return base64.StdEncoding.EncodeToString(bytes)
	}
	return string(bytes)
}
