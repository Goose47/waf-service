package dto

import gen "github.com/Goose47/wafpb/gen/go/analyzer"

type HTTPParam struct {
	Key   string
	Value string
}

type HTTPRequest struct {
	ClientIP    string      `json:"client_ip"`
	ClientPort  string      `json:"client_port"`
	ServerIP    string      `json:"server_ip"`
	ServerPort  string      `json:"server_port"`
	URI         string      `json:"uri"`
	Method      string      `json:"method"`
	Proto       string      `json:"proto"`
	Headers     []HTTPParam `json:"headers"`
	QueryParams []HTTPParam `json:"query_params"`
	BodyParams  []HTTPParam `json:"body_params"`
}

func NewHTTPRequest(
	clientIP string,
	clientPort string,
	serverIP string,
	serverPort string,
	URI string,
	method string,
	proto string,
	headers []*gen.AnalyzeRequest_HTTPParam,
	queryParams []*gen.AnalyzeRequest_HTTPParam,
	bodyParams []*gen.AnalyzeRequest_HTTPParam,
) *HTTPRequest {
	req := &HTTPRequest{
		ClientIP:   clientIP,
		ClientPort: clientPort,
		ServerIP:   serverIP,
		ServerPort: serverPort,
		URI:        URI,
		Method:     method,
		Proto:      proto,
	}

	req.Headers = make([]HTTPParam, len(headers))
	fillParams(headers, req.Headers)

	req.QueryParams = make([]HTTPParam, len(queryParams))
	fillParams(queryParams, req.QueryParams)

	req.BodyParams = make([]HTTPParam, len(bodyParams))
	fillParams(bodyParams, req.BodyParams)

	return req
}

func fillParams(params []*gen.AnalyzeRequest_HTTPParam, reqParams []HTTPParam) {
	reqParams = make([]HTTPParam, len(params))
	for i, param := range params {
		reqParams[i] = HTTPParam{
			Key:   param.Key,
			Value: param.Value,
		}
	}
}
