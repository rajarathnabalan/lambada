package lambada

import (
	"encoding/base64"
	"strings"
)

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

// isBinary attempts to detect whether content represented by contentType and contentEncoding is binary or text.
func isBinary(contentType, contentEncoding string) binaryStatus {
	// A non-empty content encoding means binary (i.e. deflate, gzip, br)
	if contentEncoding != "" {
		return bsBinary
	}

	// Discard anything before ;
	base, _, _ := strings.Cut(contentType, ";")
	mediaType := strings.TrimSpace(strings.ToLower(base))
	// Split between type and subtype (e.g text/plain)
	typ, subTyp, _ := strings.Cut(mediaType, "/")

	for _, detector := range binDetectors {
		status := detector.isBinary(mediaType, typ, subTyp)
		switch status {
		case bsBinary:
			return bsBinary
		case bsText:
			return bsText
		}
	}

	// We were not able to guess if this is binary or not, return the default
	return bsUnknown
}

type binaryStatus int8

const (
	bsUnknown binaryStatus = iota
	bsText
	bsBinary
)

type binDetector interface {
	isBinary(mediaType, typ, subTyp string) binaryStatus
}

type typeBinDetector struct {
	typ    string
	status binaryStatus
}

func (m *typeBinDetector) isBinary(mediaType, typ, subTyp string) binaryStatus {
	if typ == m.typ {
		return m.status
	}
	return bsUnknown
}

type mediaTypeBinDetector struct {
	mediaType string
	status    binaryStatus
}

func (m *mediaTypeBinDetector) isBinary(mediaType, typ, subTyp string) binaryStatus {
	if mediaType == m.mediaType {
		return m.status
	}
	return bsUnknown
}

type mediaTypeSuffixBinDetector struct {
	suffix string
	status binaryStatus
}

func (m *mediaTypeSuffixBinDetector) isBinary(mediaType, typ, subTyp string) binaryStatus {
	if strings.HasSuffix(mediaType, m.suffix) {
		return m.status
	}
	return bsUnknown
}

var binDetectors = []binDetector{
	// text/ is known to be text
	&typeBinDetector{
		typ:    "text",
		status: bsText,
	},

	// Media types ending with +json or +xml are known to be text
	&mediaTypeSuffixBinDetector{
		suffix: "+json",
		status: bsText,
	},
	&mediaTypeSuffixBinDetector{
		suffix: "+xml",
		status: bsText,
	},

	// Some application/ types are text
	&mediaTypeBinDetector{
		mediaType: "application/json",
		status:    bsText,
	},
	&mediaTypeBinDetector{
		mediaType: "application/x-javascript",
		status:    bsText,
	},
	&mediaTypeBinDetector{
		mediaType: "application/javascript",
		status:    bsText,
	},
	// application/ is considered otherwise to be binary
	&typeBinDetector{
		typ:    "application",
		status: bsBinary,
	},

	// The following is usually binary (except image/svg+xml, handled with the +xml rule)
	&typeBinDetector{
		typ:    "multipart",
		status: bsBinary,
	},
	&typeBinDetector{
		typ:    "image",
		status: bsBinary,
	},
	&typeBinDetector{
		typ:    "audio",
		status: bsBinary,
	},
	&typeBinDetector{
		typ:    "video",
		status: bsBinary,
	},
	&typeBinDetector{
		typ:    "font",
		status: bsBinary,
	},
}
