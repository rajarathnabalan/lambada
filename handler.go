package lambada

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
)

// LambdaHandler is a Lambada lambda handler function which can be used with lambda.Start.
type LambadaHandler func(ctx context.Context, req Request) (Response, error)

// NewHandler returns a Lambda function handler which can be used with lambda.Start.
// The returned lambda handler wraps incoming requets into http.Request, calls the provided http.Handler and converts
// the response into an API Gateway response.
func NewHandler(h http.Handler, options ...Option) LambadaHandler {
	opts := newOptions(options...)

	return func(ctx context.Context, req Request) (Response, error) {
		opts.requestLogger.Printf("Got request: %s\n", marshalJSON(&req))

		w := newResponseWriter(opts.outputMode, opts.defaultBinary)

		// Find out which version it is
		var httpRequest *http.Request
		var err error
		if req.Version == "2.0" {
			httpRequest, err = makeV2Request(ctx, &req)
		} else {
			httpRequest, err = makeV1Request(ctx, &req)
		}
		if err != nil {
			return Response{}, err
		}

		// Let the handler process the request
		h.ServeHTTP(w, httpRequest)
		w.finalize()

		res := Response{
			StatusCode:        w.statusCode,
			Headers:           toSingleValueHeaders(w.lockedHeader),
			MultiValueHeaders: w.lockedHeader,
			Body:              bytesToBody(w.body.Bytes(), w.binary),
			IsBase64Encoded:   w.binary,
		}
		opts.responseLogger.Printf("Response: %s\n", marshalJSON(&res))
		return res, nil
	}
}

// Serve starts the Lambda handler using the http.Handler to serve incoming requests.
// Serve calls lambda.Start(NewHandler(h)) under the hood.
func Serve(h http.Handler) {
	ServeWithOptions(h)
}

// ServeWithOptions starts the lambda handler using the http.Handler and options to serve incoming requests.
// ServeWithOptions calls lambda.Start(NewHandler(h, options...)) under the hood.
func ServeWithOptions(h http.Handler, options ...Option) {
	lambda.Start(NewHandler(h, options...))
}

func marshalJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("<failed to marshal json: %v>", err)
	}
	return string(data)
}
