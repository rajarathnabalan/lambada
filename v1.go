package lambada

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
)

// makeV1Request converts the API Gateway V1 request stored into req into an http.Request
func makeV1Request(ctx context.Context, req *Request) (*http.Request, error) {
	body, err := bodyToBytes(req.Body, req.IsBase64Encoded)
	if err != nil {
		return nil, err
	}

	// Build the initial request
	httpReq, err := http.NewRequestWithContext(WithRequest(ctx, req), req.HTTPMethod, "", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// Update the request
	httpReq.URL.Path = req.Path

	if len(req.MultiValueQueryStringParameters) == 0 && len(req.QueryStringParameters) > 0 {
		httpReq.URL.RawQuery = toURLValues(req.QueryStringParameters).Encode()
	} else {
		httpReq.URL.RawQuery = url.Values(req.MultiValueQueryStringParameters).Encode()
	}

	if len(req.MultiValueHeaders) == 0 && len(req.Headers) > 0 {
		httpReq.Header = fromSingleValueHeaders(req.Headers)
	} else {
		httpReq.Header = canonicalizeHeader(req.MultiValueHeaders)
	}

	if req.RequestContext.DomainName == "" {
		req.RequestContext.DomainName = httpReq.Header.Get("host")
		port := httpReq.Header.Get("x-forward-port")
		if req.RequestContext.DomainName != "" && port != "" {
			req.RequestContext.DomainName += ":" + port
		}
	}
	httpReq.Host = req.RequestContext.DomainName

	if req.RequestContext.Protocol == "" {
		req.RequestContext.Protocol = httpReq.Header.Get("x-forward-proto")
	}
	httpReq.Proto = req.RequestContext.Protocol
	httpReq.ProtoMajor, httpReq.ProtoMinor, _ = http.ParseHTTPVersion(req.RequestContext.Protocol)

	if req.RequestContext.Identity.SourceIP == "" {
		req.RequestContext.Identity.SourceIP = httpReq.Header.Get("x-forward-for")
	}
	httpReq.RemoteAddr = req.RequestContext.Identity.SourceIP

	return httpReq, nil
}
