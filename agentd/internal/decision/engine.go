package decision

import (
	"context"
	"time"

	"github.com/awwyushh/chameleon/agentd/internal/ml"
	rulespkg "github.com/awwyushh/chameleon/agentd/internal/rules"
	"github.com/rs/zerolog"
)

type Action int

const (
	ActPass Action = iota
	ActDeceive
	ActTarpit
	ActSpawn
)

type Decision struct {
	Action       Action
	Label        string
	Confidence   float64
	TemplateID   string
	DelayMs      int
	HoneypotHost string
	HoneypotPort int
	Trace        []string
}

type Policy struct {
	MLConfidence      float64
	SpawnOnConfidence float64
}

type Engine struct {
	logger *zerolog.Logger
	rules  *rulespkg.Engine
	ml     *ml.MLClient
	policy *Policy
}

func NewEngine(logger *zerolog.Logger, rules *rulespkg.Engine, mlclient *ml.MLClient, pol Policy) *Engine {
	return &Engine{
		logger: logger,
		rules:  rules,
		ml:     mlclient,
		policy: &pol,
	}
}

func (e *Engine) Decide(ctx context.Context, srcIP, path, method, body string) (*Decision, error) {
	dec := &Decision{Action: ActPass, DelayMs: 0}

	// 1) heuristics (but allow promotion to spawn)
	if rule, ok := e.rules.Match(body); ok {
		dec.Label = rule.Rule.Label
		dec.Confidence = rule.Rule.Confidence
		dec.TemplateID = mapLabelToTemplate(rule.Rule.Label)
		dec.Trace = append(dec.Trace, "heuristic:"+rule.Rule.ID)

		// If heuristic confidence is high enough, spawn; else deceive.
		if rule.Rule.Confidence >= e.policy.SpawnOnConfidence {
			dec.Action = ActSpawn
		} else {
			dec.Action = ActDeceive
		}
		// Return here — we treat heuristics as authoritative (but spawn possible).
		return dec, nil
	}

	// 2) ML fallback
	if e.ml != nil {
		ctx2, cancel := context.WithTimeout(ctx, time.Second*3)
		defer cancel()

		mr, err := e.ml.Predict(ctx2, body)
		if err == nil && mr != nil {
			dec.Label = mr.Label
			dec.Confidence = mr.Confidence
			dec.Trace = append(dec.Trace, "ml:predicted")

			if mr.Label != "benign" && mr.Confidence >= e.policy.MLConfidence {
				// choose to spawn if ML confidence >= spawn threshold
				if mr.Confidence >= e.policy.SpawnOnConfidence {
					dec.Action = ActSpawn
				} else {
					dec.Action = ActDeceive
				}
				dec.TemplateID = mapLabelToTemplate(mr.Label)
				return dec, nil
			}
		} else {
			e.logger.Warn().Err(err).Msg("ml predict error")
		}
	}

	// default: PASS
	dec.Action = ActPass
	dec.Trace = append(dec.Trace, "default:pass")
	
	return dec, nil
}

func mapLabelToTemplate(label string) string {
	switch label {
	case "sqli":
		return "sql_error_1064"
	case "xss":
		return "xss_reflect"
	case "bruteforce":
		return "fake_login_locked"
	default:
		return "default_generic"
	}
}
