package lambada

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/morelj/lambada/jwtclaims"
)

// Request represents an API Gateway event.
// This struct is both compatible with V1 (Lambda Proxy Integration) and V2 (HTTP API) events and is basically a merge
// of the `APIGatewayProxyRequest` and `APIGatewayV2HTTPRequest` structs defined in the
// `github.com/aws/aws-lambda-go/events` package.
type Request struct {
	// V1 Only
	Resource                        string              `json:"resource"` // The resource path defined in API Gateway
	Path                            string              `json:"path"`     // The url path for the caller
	HTTPMethod                      string              `json:"httpMethod"`
	MultiValueHeaders               map[string][]string `json:"multiValueHeaders"`
	MultiValueQueryStringParameters map[string][]string `json:"multiValueQueryStringParameters"`

	// V2 Only
	Version        string   `json:"version"`
	RouteKey       string   `json:"routeKey"`
	RawPath        string   `json:"rawPath"`
	RawQueryString string   `json:"rawQueryString"`
	Cookies        []string `json:"cookies,omitempty"`

	// V1 + V2
	Headers               map[string]string `json:"headers"`
	QueryStringParameters map[string]string `json:"queryStringParameters,omitempty"`
	PathParameters        map[string]string `json:"pathParameters,omitempty"`
	StageVariables        map[string]string `json:"stageVariables,omitempty"`
	Body                  string            `json:"body,omitempty"`
	IsBase64Encoded       bool              `json:"isBase64Encoded,omitempty"`
	RequestContext        RequestContext    `json:"requestContext"`
}

// RequestContext contains the information to identify the AWS account and resources invoking the Lambda function.
// This struct is both compatible with V1 (Lambda Proxy Integration) and V2 (HTTP API) events and is basically a merge
// of the `APIGatewayProxyRequestContext` and `APIGatewayV2HTTPRequestContext` structs defined in the
// `github.com/aws/aws-lambda-go/events` package.
type RequestContext struct {
	// V1 Only
	ResourceID       string                           `json:"resourceId"`
	OperationName    string                           `json:"operationName,omitempty"`
	Protocol         string                           `json:"protocol"`
	Identity         events.APIGatewayRequestIdentity `json:"identity"`
	ResourcePath     string                           `json:"resourcePath"`
	HTTPMethod       string                           `json:"httpMethod"`
	RequestTime      string                           `json:"requestTime"`
	RequestTimeEpoch int64                            `json:"requestTimeEpoch"`

	// V2 Only
	RouteKey  string                                               `json:"routeKey"`
	Time      string                                               `json:"time"`
	TimeEpoch int64                                                `json:"timeEpoch"`
	HTTP      events.APIGatewayV2HTTPRequestContextHTTPDescription `json:"http"`

	// V1 + V2
	AccountID    string      `json:"accountId"`
	Stage        string      `json:"stage"`
	DomainName   string      `json:"domainName"`
	DomainPrefix string      `json:"domainPrefix"`
	RequestID    string      `json:"requestId"`
	APIID        string      `json:"apiId"` // The API Gateway rest API Id
	Authorizer   *Authorizer `json:"authorizer,omitempty"`
}

// Response contains the response to send back to API Gateway
// This struct is both compatible with V1 (Lambda Proxy Integration) and V2 (HTTP API) events and is basically a merge
// of the `APIGatewayProxyResponse` and `APIGatewayV2HTTPResponse` structs defined in the
// `github.com/aws/aws-lambda-go/events` package.
type Response struct {
	// V2 Only
	Cookies []string `json:"cookies,omitempty"`

	// V1 + V2
	StatusCode        int                 `json:"statusCode"`
	Headers           map[string]string   `json:"headers"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	Body              string              `json:"body"`
	IsBase64Encoded   bool                `json:"isBase64Encoded,omitempty"`
}

// Authorizer contains authorizer details
type Authorizer struct {
	IAM *IAMAuthorizer `json:"iam,omitempty"`
	JWT *JWTAuthorizer `json:"jwt,omitempty"`
}

// IAMAuthorizer contains the details of a request authenticated using the AWS SignV4 authorizer.
type IAMAuthorizer struct {
	AccessKey      string `json:"accessKey,omitempty"`
	AccountID      string `json:"accountId,omitempty"`
	CallerID       string `json:"callerId,omitempty"`
	PrincipalOrgID string `json:"principalOrgId,omitempty"`
	UserARN        string `json:"userArn,omitempty"`
	UserID         string `json:"userId,omitempty"`
}

// JWTAuthorizer contains the details of a request authenticated using the JWT authorizer.
type JWTAuthorizer struct {
	Claims jwtclaims.Claims `json:"claims,omitempty"`
	Scopes interface{}      `json:"scopes,omitempty"`
}
