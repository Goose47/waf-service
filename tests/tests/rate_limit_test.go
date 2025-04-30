package tests

import (
	"github.com/Goose47/wafpb/gen/go/detection"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"tests/internal/suite"
	"time"
)

func TestRateLimits_MaxAllowedRequestsPlusOne(t *testing.T) {
	s, ctx := suite.New(t)

	req := newAnalyzeRequest(t)
	req.Method = "GET"
	req.Timestamp = timestamppb.New(time.Now())

	// Requests' timestamps are equal, so these requests will be treated as simultaneous.
	for i := range s.LimiterCfg.MaxRequests + 1 {
		// Call CheckIP to add next fingerprint.
		_, err := s.Detection.CheckIP(ctx, &detection.CheckIPRequest{
			Ip:        req.ClientIp,
			Timestamp: req.Timestamp,
		})
		require.NoError(t, err)

		res, err := s.Analyzer.Analyze(ctx, req)
		require.NoError(t, err)
		require.Equal(t, i == s.LimiterCfg.MaxRequests, res.IsAttack)
	}
}

func TestRateLimits_MaxAllowedRequestsPlusOneAfterDelay(t *testing.T) {
	s, ctx := suite.New(t)

	req := newAnalyzeRequest(t)
	req.Method = "GET"
	req.Timestamp = timestamppb.New(time.Now().Add(-s.LimiterCfg.Per / 2))

	// Requests' timestamps are equal, so these requests will be treated as simultaneous.
	for range s.LimiterCfg.MaxRequests {
		// Call CheckIP to add next fingerprint.
		_, err := s.Detection.CheckIP(ctx, &detection.CheckIPRequest{
			Ip:        req.ClientIp,
			Timestamp: req.Timestamp,
		})
		require.NoError(t, err)

		res, err := s.Analyzer.Analyze(ctx, req)
		require.NoError(t, err)
		require.Equal(t, false, res.IsAttack)
	}

	// Simulate waiting.
	req.Timestamp = timestamppb.New(req.Timestamp.AsTime().Add(s.LimiterCfg.Per))

	// Call CheckIP to add next fingerprint.
	_, err := s.Detection.CheckIP(ctx, &detection.CheckIPRequest{
		Ip:        req.ClientIp,
		Timestamp: req.Timestamp,
	})
	require.NoError(t, err)

	res, err := s.Analyzer.Analyze(ctx, req)
	require.NoError(t, err)
	require.Equal(t, false, res.IsAttack)
}
