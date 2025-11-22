package server

import (
	"context"

	"github.com/awwyushh/chameleon/agentd/internal/composer"
	decisionpkg "github.com/awwyushh/chameleon/agentd/internal/decision"
	honeypot "github.com/awwyushh/chameleon/agentd/internal/honeypot"
	tarpitpkg "github.com/awwyushh/chameleon/agentd/internal/tarpit"
	telemetry "github.com/awwyushh/chameleon/agentd/internal/telemetry"
	templatespkg "github.com/awwyushh/chameleon/agentd/internal/templates"
	pb "github.com/awwyushh/chameleon/agentd/proto"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type AgentServer struct {
	pb.UnimplementedChameleonServer

	Logger        zerolog.Logger
	Engine        *decisionpkg.Engine
	Tarpit        *tarpitpkg.Tarpit
	Templates     *templatespkg.Loader
	Honeypot      *honeypot.Manager
	Telemetry     *telemetry.Forwarder
	AttacksMetric *prometheus.CounterVec

	KillSwitchPtr *bool
}

func (s *AgentServer) Classify(ctx context.Context, req *pb.ClassifyRequest) (*pb.DecisionResponse, error) {
	// 1. Decision engine
	dec, err := s.Engine.Decide(ctx, req.SrcIp, req.Path, req.Method, req.Body)
	if err != nil {
		s.Logger.Error().Err(err).Msg("decision engine error")
		return &pb.DecisionResponse{Action: pb.DecisionResponse_PASS, Message: "internal error"}, nil
	}

	// 2. Tarpit logic
	delay, err := s.Tarpit.ShouldTarpit(ctx, req.SrcIp)
	if err == nil && delay > 0 {
		dec.DelayMs = delay
	}

	// 3. Honeypot spawn
	if dec.Action == decisionpkg.ActSpawn {
		if s.KillSwitchPtr != nil && *s.KillSwitchPtr {
			s.Logger.Warn().Msg("kill-switch active: blocking honeypot spawn")
			dec.Action = decisionpkg.ActDeceive
		} else {
			si, err := s.Honeypot.Spawn(ctx, req.Body, req.SrcIp, dec.Label)
			if err != nil {
				s.Logger.Error().Err(err).Msg("honeypot spawn failed, falling back to deception")
				dec.Action = decisionpkg.ActDeceive
			} else {
				dec.HoneypotHost = "127.0.0.1"
				dec.HoneypotPort = si.HostPort
			}
		}
	}

	// 4. Template-based deception response
	var message string
	if dec.Action == decisionpkg.ActDeceive || dec.Action == decisionpkg.ActSpawn {
		tpl, ok := s.Templates.Get(dec.TemplateID)
		if ok {
			vars := map[string]string{
				"payload_excerpt": truncate(req.Body, 80),
			}
			text, err := composer.ComposeText(tpl.Body, vars)
			if err == nil {
				message = text
			} else {
				message = "Error: Processing request"
			}
		} else {
			// Fallback if template not found
			message = "Error 500: Internal Server Error"
		}
	}

	// 5. Telemetry (non-blocking)
	event := map[string]interface{}{
		"src_ip":     req.SrcIp,
		"path":       req.Path,
		"label":      dec.Label,
		"confidence": dec.Confidence,
		"delay_ms":   dec.DelayMs,
		"template":   dec.TemplateID,
		"action":     dec.Action,
		"hp_port":    dec.HoneypotPort,
	}
	go func() { _ = s.Telemetry.SendEvent(context.Background(), event) }()

	// 6. Prometheus metrics
	if dec.Label != "" {
		s.AttacksMetric.WithLabelValues(dec.Label).Inc()
	}

	// 7. Build gRPC response
	resp := &pb.DecisionResponse{
		Message:      message,
		DelayMs:      int32(dec.DelayMs),
		Label:        dec.Label,
		Confidence:   dec.Confidence,
		HoneypotHost: dec.HoneypotHost,
		HoneypotPort: int32(dec.HoneypotPort),
	}

	switch dec.Action {
	case decisionpkg.ActPass:
		resp.Action = pb.DecisionResponse_PASS
	case decisionpkg.ActDeceive:
		resp.Action = pb.DecisionResponse_DECEIVE
	case decisionpkg.ActSpawn:
		// Use DECEIVE action but with honeypot fields populated
		// The SDK will see honeypot_port > 0 and redirect
		resp.Action = pb.DecisionResponse_DECEIVE
	case decisionpkg.ActTarpit:
		resp.Action = pb.DecisionResponse_TARPIT
	default:
		resp.Action = pb.DecisionResponse_PASS
	}

	if dec.Label != "benign" {
		s.Logger.Info().
			Str("src_ip", req.SrcIp).
			Str("path", req.Path).
			Str("label", dec.Label).
			Float64("conf", dec.Confidence).
			Msg("Attack Detected")
	} else {
		// Optional: Log benign traffic at Debug level so it doesn't flood
		s.Logger.Debug().Str("path", req.Path).Msg("Benign request processed")
	}
	s.Logger.Info().Interface("decision", dec).Msg("Decision after engine")

	return resp, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
