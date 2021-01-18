package lambada

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/morelj/lambada/internal/agw"
)

func SetBinary(w http.ResponseWriter) {
	if w, ok := w.(*agw.ResponseWriter); ok {
		w.Binary = true
	}
}

func SetText(w http.ResponseWriter) {
	if w, ok := w.(*agw.ResponseWriter); ok {
		w.Binary = false
	}
}

type contextKeyType struct{}

var contextKey = contextKeyType{}

func GetRequest(ctx context.Context) *Request {
	if res, ok := ctx.Value(contextKey).(*Request); ok {
		return res
	}
	return nil
}

func withRequest(ctx context.Context, req *Request) context.Context {
	return context.WithValue(ctx, contextKey, req)
}

type handler struct {
	httpHandler http.Handler
}

func makeV1Request(ctx context.Context, req *Request) (*http.Request, error) {
	body, err := agw.BodyToBytes(req.Body, req.IsBase64Encoded)
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
	httpReq.Header = agw.CanonicalizeHeader(req.MultiValueHeaders)
	httpReq.Host = req.RequestContext.DomainName
	httpReq.Proto = req.RequestContext.Protocol
	httpReq.ProtoMajor, httpReq.ProtoMinor, _ = http.ParseHTTPVersion(req.RequestContext.Protocol)
	httpReq.RemoteAddr = req.RequestContext.Identity.SourceIP

	return httpReq, nil
}

func makeV2Request(ctx context.Context, req *Request) (*http.Request, error) {
	body, err := agw.BodyToBytes(req.Body, req.IsBase64Encoded)
	if err != nil {
		return nil, err
	}

	// Build the initial request
	httpReq, err := http.NewRequestWithContext(withRequest(ctx, req), req.RequestContext.HTTP.Method, "", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// Update the request
	httpReq.URL.Path = req.RawPath
	httpReq.URL.RawQuery = agw.ToURLValues(req.QueryStringParameters).Encode()
	httpReq.Header = agw.FromSingleValueHeaders(req.Headers)
	httpReq.Host = req.RequestContext.DomainName
	httpReq.RemoteAddr = req.RequestContext.HTTP.SourceIP
	httpReq.Proto = req.RequestContext.HTTP.Protocol
	httpReq.ProtoMajor, httpReq.ProtoMinor, _ = http.ParseHTTPVersion(req.RequestContext.HTTP.Protocol)

	return httpReq, nil
}

func (h handler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	log.Printf("Got payload: %s\n", string(payload))
	// Parse
	var req Request
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}

	w := agw.NewResponseWriter()

	// Find out which version it is
	var httpRequest *http.Request
	var err error
	if req.Version == "2.0" {
		log.Printf("Found V2 request")
		httpRequest, err = makeV2Request(ctx, &req)
	} else {
		log.Printf("Found V1 request")
		httpRequest, err = makeV1Request(ctx, &req)
	}
	if err != nil {
		return nil, err
	}
	h.httpHandler.ServeHTTP(w, httpRequest)

	return json.Marshal(&Response{
		StatusCode:        w.StatusCode(),
		Headers:           agw.ToSingleValueHeaders(w.Header()),
		MultiValueHeaders: w.Header(),
		Body:              agw.BytesToBody(w.Body(), w.Binary),
		IsBase64Encoded:   w.Binary,
	})
}

func Serve(h http.Handler) {
	lambda.StartHandler(handler{httpHandler: h})
}
