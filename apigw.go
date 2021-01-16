package lambada

import "github.com/aws/aws-lambda-go/events"

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
	AccountID    string                 `json:"accountId"`
	Stage        string                 `json:"stage"`
	DomainName   string                 `json:"domainName"`
	DomainPrefix string                 `json:"domainPrefix"`
	RequestID    string                 `json:"requestId"`
	APIID        string                 `json:"apiId"` // The API Gateway rest API Id
	Authorizer   map[string]interface{} `json:"authorizer"`
}

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
