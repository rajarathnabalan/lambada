// Package lambada provides a compatibility layer allowing to implement AWS API Gateway V1 and V2 (HTTP APIs) Lambda
// integrations using http.Handler.
// All libraries and frameworks using net/http should work using lambada.
//
// Example:
//
//	package main
//
//	import (
//	    "net/http"
//
//	    "github.com/rajarathnabalan/lambada"
//	)
//
//	func main() {
//	    lambada.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		       w.Write(([]byte)("<html><body><h1>Hello, World!</h1></body></html>"))
//	    }))
//	}
package lambada
