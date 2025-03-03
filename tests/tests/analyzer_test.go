package tests

import (
	"fmt"
	"github.com/Goose47/wafpb/gen/go/analyzer"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"net/url"
	"strconv"
	"testing"
	"tests/internal/suite"
)

func newAnalyzeRequest(t *testing.T) *analyzer.AnalyzeRequest {
	randomURL, err := url.Parse(gofakeit.URL())
	require.NoError(t, err)

	return &analyzer.AnalyzeRequest{
		ClientIp:   gofakeit.IPv4Address(),
		ClientPort: strconv.Itoa(gofakeit.IntRange(1024, 65535)),
		ServerIp:   gofakeit.IPv4Address(),
		ServerPort: strconv.Itoa(80),
		Uri:        fmt.Sprintf("/%s", randomURL.Path),
		Method:     "",
		Proto:      "HTTP/1.1",
		Headers: []*analyzer.AnalyzeRequest_HTTPParam{
			{
				Key:   "User-Agent",
				Value: gofakeit.UserAgent(),
			},
		},
		QueryParams: make([]*analyzer.AnalyzeRequest_HTTPParam, 0),
		BodyParams:  make([]*analyzer.AnalyzeRequest_HTTPParam, 0),
	}
}

func newXSSAnalyzeRequest(t *testing.T) *analyzer.AnalyzeRequest {
	req := newAnalyzeRequest(t)
	req.Uri = "/protected?username=alert(%27XSS!%27)"
	req.Method = "GET"
	return req
}

func newSQLIAnalyzeRequest(t *testing.T) *analyzer.AnalyzeRequest {
	req := newAnalyzeRequest(t)
	req.Method = "POST"
	req.BodyParams = append(
		req.BodyParams,
		&analyzer.AnalyzeRequest_HTTPParam{
			Key:   "username",
			Value: "' OR 1=1; DROP TABLE users; --",
		},
	)
	req.Headers = append(
		req.Headers,
		&analyzer.AnalyzeRequest_HTTPParam{
			Key:   "Content-Length",
			Value: "39",
		},
		&analyzer.AnalyzeRequest_HTTPParam{
			Key:   "Content-Type",
			Value: "application/x-www-form-urlencoded",
		},
	)

	return req
}

func TestAnalyzer_GetRequest(t *testing.T) {
	s, ctx := suite.New(t)
	req := newAnalyzeRequest(t)
	req.Method = "GET"

	res, err := s.Analyzer.Analyze(ctx, req)
	require.NoError(t, err)
	require.False(t, res.IsAttack)
}

func TestAnalyzer_PostRequest(t *testing.T) {
	s, ctx := suite.New(t)

	req := newAnalyzeRequest(t)
	req.Method = "POST"

	username := gofakeit.Username()
	req.BodyParams = append(
		req.BodyParams,
		&analyzer.AnalyzeRequest_HTTPParam{
			Key:   "username",
			Value: username,
		},
	)
	req.Headers = append(
		req.Headers,
		&analyzer.AnalyzeRequest_HTTPParam{
			Key:   "Content-Length",
			Value: fmt.Sprintf("%d", len(username)),
		},
		&analyzer.AnalyzeRequest_HTTPParam{
			Key:   "Content-Type",
			Value: "application/x-www-form-urlencoded",
		},
	)

	res, err := s.Analyzer.Analyze(ctx, req)
	require.NoError(t, err)
	require.False(t, res.IsAttack)
}

func TestAnalyzer_GetXSSAttack(t *testing.T) {
	s, ctx := suite.New(t)

	req := newXSSAnalyzeRequest(t)

	res, err := s.Analyzer.Analyze(ctx, req)
	require.NoError(t, err)
	require.True(t, res.IsAttack)
}

func TestAnalyzer_PostSQLIAttack(t *testing.T) {
	s, ctx := suite.New(t)

	req := newSQLIAnalyzeRequest(t)

	res, err := s.Analyzer.Analyze(ctx, req)
	require.NoError(t, err)
	require.True(t, res.IsAttack)
}
