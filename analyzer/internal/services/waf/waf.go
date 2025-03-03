// Package waf interacts with coraza API.
package waf

import (
	"bytes"
	"context"
	"fmt"
	"github.com/corazawaf/coraza/v3"
	"log/slog"
	"strconv"
	dtopkg "waf-analyzer/internal/domain/dto"
)

// WAF Loads coraza and provides facade to call coraza API.
type WAF struct {
	log *slog.Logger
	waf coraza.WAF
}

// MustCreate creates WAF instance and panics on error.
func MustCreate(log *slog.Logger) *WAF {
	cfg := coraza.NewWAFConfig().
		WithDirectivesFromFile("/app/secrules/coraza.conf").
		WithDirectivesFromFile("/app/secrules/coreruleset/crs-setup.conf").
		WithDirectivesFromFile("/app/secrules/coreruleset/rules/*.conf").
		WithDirectives(`
			SecRuleEngine On
		`)

	waf, err := coraza.NewWAF(cfg)
	if err != nil {
		panic(err)
	}

	return &WAF{
		log: log,
		waf: waf,
	}
}

// Analyze indicates whether http request contains an attack.
func (waf *WAF) Analyze(_ context.Context, request *dtopkg.HTTPRequest) (bool, error) {
	const op = "WAF.Analyze"
	log := waf.log.With(slog.String("op", op))

	log.Info("trying to analyze request")

	tx := waf.waf.NewTransaction()
	defer func() {
		err := tx.Close()
		if err != nil {
			log.Error(fmt.Sprintf("%s: failed to close transaction", op), slog.Any("error", err))
		}
	}()

	clientPort, err := strconv.Atoi(request.ClientPort)
	if err != nil {
		log.Error("failed to parse client port", slog.Any("error", err))
		return false, fmt.Errorf("%s: %w", op, err)
	}
	serverPort, err := strconv.Atoi(request.ServerPort)
	if err != nil {
		log.Error("failed to parse server port", slog.Any("error", err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	tx.ProcessConnection(request.ClientIP, clientPort, request.ServerIP, serverPort)
	tx.ProcessURI(request.URI, request.Method, request.Proto)
	// Fill Headers
	for _, header := range request.Headers {
		tx.AddRequestHeader(header.Key, header.Value)
	}
	// Fill GET parameters
	for _, param := range request.QueryParams {
		tx.AddGetRequestArgument(param.Key, param.Value)
	}

	// Process phase 1
	if it := tx.ProcessRequestHeaders(); it != nil {
		log.Warn("Attack is found in phase 1", slog.Int("rule_id", it.RuleID))
		return true, nil
	}

	// Fake payload content-type
	tx.AddRequestHeader("Content-Type", "application/x-www-form-urlencoded")

	// Fill POST parameters
	body := make([]byte, 0)
	for _, param := range request.BodyParams {
		body = append(body, param.Key...)
		body = append(body, "="...)
		body = append(body, param.Value...)
		body = append(body, "&"...)
	}
	body = bytes.TrimSuffix(body, []byte("&"))

	it, _, err := tx.WriteRequestBody(body)

	if err != nil {
		log.Error("failed to write request body", slog.Any("error", err))
		return false, fmt.Errorf("%s: %w", op, err)
	}
	if it != nil {
		log.Warn(
			"interrupted while writing request body",
			slog.Int("rule_id", it.RuleID),
			slog.String("data", it.Data),
			slog.String("action", it.Action),
			slog.Int("status", it.Status),
		)
		return true, nil
	}

	// Process phase 2
	it, err = tx.ProcessRequestBody()
	if err != nil {
		log.Error("process request body", slog.Any("error", err))
		return false, fmt.Errorf("%s: %w", op, err)
	}
	if it != nil {
		log.Warn(
			"Attack is found in phase 2",
			slog.Int("rule_id", it.RuleID),
			slog.String("data", it.Data),
			slog.String("action", it.Action),
			slog.Int("status", it.Status),
		)
		return true, nil
	}

	return false, nil
}
