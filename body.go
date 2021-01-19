package lambada

import "encoding/base64"

func bodyToBytes(body string, isBase64 bool) ([]byte, error) {
	if isBase64 {
		return base64.StdEncoding.DecodeString(body)
	}
	return ([]byte)(body), nil
}

func bytesToBody(bytes []byte, isBase64 bool) string {
	if isBase64 {
		return base64.StdEncoding.EncodeToString(bytes)
	}
	return string(bytes)
}
