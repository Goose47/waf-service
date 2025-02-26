package dto

import gen "github.com/Goose47/wafpb/gen/go/waf"

type HTTPParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type HTTPRequest struct {
	ClientIP    string       `json:"client_ip"`
	ClientPort  string       `json:"client_port"`
	ServerIP    string       `json:"server_ip"`
	ServerPort  string       `json:"server_port"`
	URI         string       `json:"uri"`
	Method      string       `json:"method"`
	Proto       string       `json:"proto"`
	Headers     []*HTTPParam `json:"headers"`
	QueryParams []*HTTPParam `json:"query_params"`
	BodyParams  []*HTTPParam `json:"body_params"`
}

// NewHTTPRequest is a factory function for creating DTO from grpc request.
func NewHTTPRequest(req *gen.AnalyzeRequest) *HTTPRequest {
	res := &HTTPRequest{
		ClientIP:   req.ClientIp,
		ClientPort: req.ClientPort,
		ServerIP:   req.ServerIp,
		ServerPort: req.ServerPort,
		URI:        req.Uri,
		Method:     req.Method,
		Proto:      req.Proto,
	}

	res.Headers = make([]*HTTPParam, len(req.Headers))
	fillDTOParams(req.Headers, res.Headers)

	res.QueryParams = make([]*HTTPParam, len(req.QueryParams))
	fillDTOParams(req.QueryParams, res.QueryParams)

	res.BodyParams = make([]*HTTPParam, len(req.BodyParams))
	fillDTOParams(req.BodyParams, res.BodyParams)

	return res
}

func fillDTOParams(params []*gen.AnalyzeRequest_HTTPParam, reqParams []*HTTPParam) {
	for i, param := range params {
		reqParams[i] = &HTTPParam{
			Key:   param.Key,
			Value: param.Value,
		}
	}
}

// ToAnalyzeRequest is a factory method for creating grpc request from DTO.
func (req *HTTPRequest) ToAnalyzeRequest() *gen.AnalyzeRequest {
	res := &gen.AnalyzeRequest{
		ClientIp:   req.ClientIP,
		ClientPort: req.ClientPort,
		ServerIp:   req.ServerIP,
		ServerPort: req.ServerPort,
		Uri:        req.URI,
		Method:     req.Method,
		Proto:      req.Proto,
	}

	res.Headers = make([]*gen.AnalyzeRequest_HTTPParam, len(req.Headers))
	fillParams(res.Headers, req.Headers)

	res.QueryParams = make([]*gen.AnalyzeRequest_HTTPParam, len(req.QueryParams))
	fillParams(res.QueryParams, req.QueryParams)

	res.BodyParams = make([]*gen.AnalyzeRequest_HTTPParam, len(req.BodyParams))
	fillParams(res.BodyParams, req.BodyParams)

	return res
}

func fillParams(params []*gen.AnalyzeRequest_HTTPParam, reqParams []*HTTPParam) {
	for i, param := range reqParams {
		params[i] = &gen.AnalyzeRequest_HTTPParam{
			Key:   param.Key,
			Value: param.Value,
		}
	}
}
