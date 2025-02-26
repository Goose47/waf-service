package waf

import (
	"context"
	"fmt"
	"github.com/corazawaf/coraza/v3"
	"log/slog"
	"strconv"
	dtopkg "waf-analyzer/internal/domain/dto"
)

type WAF struct {
	log *slog.Logger
	waf coraza.WAF
}

func MustCreate(log *slog.Logger) *WAF {
	// todo: move paths to config
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
func (waf *WAF) Analyze(ctx context.Context, request *dtopkg.HTTPRequest) (bool, error) {
	const op = "WAF.Analyze"
	log := waf.log.With(slog.String("op", op))

	log.Info("trying to analyze request")

	tx := waf.waf.NewTransaction()
	defer tx.Close()

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
		log.Warn("Attack is found", slog.Int("rule_id", it.RuleID))
		return true, nil
	}

	// todo: maybe accept []byte as body?
	// Content-Type is important to tell coraza which BodyProcessor must be used
	//tx.AddRequestHeader("Content-Type", "application/x-www-form-urlencoded")
	//tx.AddRequestHeader("Content-Type", "application/json")
	//res, _ := io.ReadAll(r.Body)
	//_, _, err = tx.WriteRequestBody(res)

	// Fill POST parameters
	for _, param := range request.QueryParams {
		tx.AddPostRequestArgument(param.Key, param.Value)
	}

	// Process phase 2
	it, err := tx.ProcessRequestBody()
	if err != nil {
		log.Error("process request body", slog.Any("error", err))
		return false, fmt.Errorf("%s: %w", op, err)
	}
	if it != nil {
		log.Warn("Attack is found", slog.Int("rule_id", it.RuleID))
		return true, nil
	}

	return false, nil
}
