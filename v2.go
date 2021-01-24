package lambada

import (
	"bytes"
	"context"
	"net/http"
)

// makeV2Request converts the API Gateway V2 request stored into req into an http.Request
func makeV2Request(ctx context.Context, req *Request) (*http.Request, error) {
	body, err := bodyToBytes(req.Body, req.IsBase64Encoded)
	if err != nil {
		return nil, err
	}

	// Build the initial request
	httpReq, err := http.NewRequestWithContext(WithRequest(ctx, req), req.RequestContext.HTTP.Method, "", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// Update the request
	httpReq.URL.Path = req.RawPath
	httpReq.URL.RawQuery = toURLValues(req.QueryStringParameters).Encode()
	httpReq.Header = fromSingleValueHeaders(req.Headers)
	httpReq.Host = req.RequestContext.DomainName
	httpReq.RemoteAddr = req.RequestContext.HTTP.SourceIP
	httpReq.Proto = req.RequestContext.HTTP.Protocol
	httpReq.ProtoMajor, httpReq.ProtoMinor, _ = http.ParseHTTPVersion(req.RequestContext.HTTP.Protocol)

	return httpReq, nil
}
