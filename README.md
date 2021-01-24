# Lambada - Go net/http compatibility layer for AWS API Gateway Lambda functions

Lambada is a small Go package which provides a layer allowing to use Go's standard library `net/http` package to
handle AWS API Gateway events in AWS Lambda functions.

It basically converts API Gateway events into `http.Request`, calls an `http.Handler` and converts the result written
to the `http.Response` into an API Gateway response.

All libraries using `http.Handler` (e.g. multiplexers) should work using Lambada.

Lambada is compatible with both API Gateway V1 using the Lambada Proxy integration, and API Gateway V2 (HTTP API).

## Installation

Install using `go get`:

```
go get github.com/morelj/lambada
```

## Example

```go
import (
    "net/http"
    "github.com/morelj/lambada"
)

func main() {
    lambada.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	       w.Write(([]byte)("<html><body><h1>Hello, World!</h1></body></html>"))
    }))
}
```
