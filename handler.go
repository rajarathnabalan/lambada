package lambada

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
)

type handler struct {
	httpHandler http.Handler
}

func (h handler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {

	log.Printf("Got payload: %s\n", string(payload))
	// Parse
	var req Request
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}

	w := newResponseWriter()

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
		Headers:           toSingleValueHeaders(w.Header()),
		MultiValueHeaders: w.Header(),
		Body:              bytesToBody(w.Body(), w.Binary),
		IsBase64Encoded:   w.Binary,
	})
}

func Serve(h http.Handler) {
	lambda.StartHandler(handler{httpHandler: h})
}
