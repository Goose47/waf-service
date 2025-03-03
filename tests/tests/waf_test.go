package tests

import (
	"github.com/Goose47/wafpb/gen/go/analyzer"
	"github.com/Goose47/wafpb/gen/go/waf"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"testing"
	"tests/internal/suite"
	"time"
)

func newWAFRequest(t *testing.T) *waf.AnalyzeRequest {
	req := newAnalyzeRequest(t)

	return &waf.AnalyzeRequest{
		ClientIp:    req.ClientIp,
		ClientPort:  req.ClientPort,
		ServerIp:    req.ServerIp,
		ServerPort:  req.ServerPort,
		Uri:         req.Uri,
		Method:      req.Method,
		Proto:       req.Proto,
		Headers:     mapParams(req.Headers),
		QueryParams: mapParams(req.QueryParams),
		BodyParams:  mapParams(req.BodyParams),
	}
}

func newXSSWAFRequest(t *testing.T) *waf.AnalyzeRequest {
	req := newAnalyzeRequest(t)
	req.Uri = "/protected?username=alert(%27XSS!%27)"
	req.Method = "GET"

	return &waf.AnalyzeRequest{
		ClientIp:    req.ClientIp,
		ClientPort:  req.ClientPort,
		ServerIp:    req.ServerIp,
		ServerPort:  req.ServerPort,
		Uri:         req.Uri,
		Method:      req.Method,
		Proto:       req.Proto,
		Headers:     mapParams(req.Headers),
		QueryParams: mapParams(req.QueryParams),
		BodyParams:  mapParams(req.BodyParams),
	}
}

func mapParams(params []*analyzer.AnalyzeRequest_HTTPParam) []*waf.AnalyzeRequest_HTTPParam {
	mappedParams := make([]*waf.AnalyzeRequest_HTTPParam, len(params))
	for i, param := range params {
		mappedParams[i] = &waf.AnalyzeRequest_HTTPParam{
			Key:   param.Key,
			Value: param.Value,
		}
	}
	return mappedParams
}

func TestWAF_MultipleRequests(t *testing.T) {
	s, ctx := suite.New(t)

	// Generate random ip address.
	ip := gofakeit.IPv4Address()

	goodReq := newWAFRequest(t)
	goodReq.ClientIp = ip
	goodReq.Method = "GET"

	// Good requests should not be attacks.
	res, err := s.Waf.Analyze(ctx, goodReq)
	require.NoError(t, err)
	require.Equal(t, float32(0), res.AttackProbability)

	// First attack will be analyzed in background.
	xssReq := newXSSWAFRequest(t)
	xssReq.ClientIp = ip

	res, err = s.Waf.Analyze(ctx, xssReq)
	require.NoError(t, err)
	require.Equal(t, float32(0), res.AttackProbability)

	// Next attack will be blocked.
	require.Eventually(t, func() bool {
		res, err = s.Waf.Analyze(ctx, xssReq)
		require.NoError(t, err)
		return res.AttackProbability > 0
	}, time.Second, time.Millisecond)

	// Requests that are not attacks will pass.
	res, err = s.Waf.Analyze(ctx, goodReq)
	require.NoError(t, err)
	require.Equal(t, float32(0), res.AttackProbability)
}
