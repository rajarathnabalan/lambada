package lambada

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
)

func makeV1Request(ctx context.Context, req *Request) (*http.Request, error) {
	body, err := bodyToBytes(req.Body, req.IsBase64Encoded)
	if err != nil {
		return nil, err
	}

	// Build the initial request
	httpReq, err := http.NewRequestWithContext(withRequest(ctx, req), req.HTTPMethod, "", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// Update the request
	httpReq.URL.Path = req.Path
	httpReq.URL.RawQuery = url.Values(req.MultiValueQueryStringParameters).Encode()
	httpReq.Header = canonicalizeHeader(req.MultiValueHeaders)
	httpReq.Host = req.RequestContext.DomainName
	httpReq.Proto = req.RequestContext.Protocol
	httpReq.ProtoMajor, httpReq.ProtoMinor, _ = http.ParseHTTPVersion(req.RequestContext.Protocol)
	httpReq.RemoteAddr = req.RequestContext.Identity.SourceIP

	return httpReq, nil
}
