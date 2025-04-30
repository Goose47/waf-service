package tests

import (
	"github.com/Goose47/wafpb/gen/go/detection"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"tests/internal/suite"
	"time"
)

func TestDetection_CheckIP(t *testing.T) {
	s, ctx := suite.New(t)

	// Generate random ip address.
	ip := gofakeit.IPv4Address()

	checkIP := func() bool {
		req := &detection.CheckIPRequest{
			Ip:        ip,
			Timestamp: timestamppb.New(time.Now()),
		}

		res, err := s.Detection.CheckIP(ctx, req)
		require.NoError(t, err)
		return res.IsSuspicious
	}

	// A random IP address should not be suspicious.
	res := checkIP()
	require.False(t, res)

	// Performing SQLI attack.
	req := newSQLIAnalyzeRequest(t)
	req.ClientIp = ip
	analyzeRes, err := s.Analyzer.Analyze(ctx, req)
	require.NoError(t, err)
	require.True(t, analyzeRes.IsAttack)

	// IP address should be suspicious after performing attack.
	require.Eventually(t, checkIP, time.Second, time.Millisecond)
}
