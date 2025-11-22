package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	pb "github.com/awwyushh/chameleon/agentd/proto"

	"github.com/awwyushh/chameleon/agentd/internal/config"
	decisionpkg "github.com/awwyushh/chameleon/agentd/internal/decision"
	honeypot "github.com/awwyushh/chameleon/agentd/internal/honeypot"
	mlpkg "github.com/awwyushh/chameleon/agentd/internal/ml"
	rulespkg "github.com/awwyushh/chameleon/agentd/internal/rules"
	server "github.com/awwyushh/chameleon/agentd/internal/server"
	tarpitpkg "github.com/awwyushh/chameleon/agentd/internal/tarpit"
	telemetry "github.com/awwyushh/chameleon/agentd/internal/telemetry"
	templatespkg "github.com/awwyushh/chameleon/agentd/internal/templates"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var killSwitch bool = false

func main() {
	socket := flag.String("socket", "/tmp/chameleon.sock", "unix socket path")
	policyPath := flag.String("policy", "policy.yaml", "policy path")
	flag.Parse()

	// logger setup
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	logger := log.With().Str("component", "agentd").Logger()

	// load policy
	pol, err := config.LoadPolicy(*policyPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load policy")
	}

	// load rules
	rulesEngine, err := rulespkg.LoadRules("./rules")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load rules")
	}

	// ML client
	mlClient := mlpkg.New(pol.ML.URL, pol.ML.TimeoutSeconds)

	// decision engine
	engine := decisionpkg.NewEngine(&logger, rulesEngine, mlClient, decisionpkg.Policy{
		MLConfidence:      pol.ML.ConfidenceThreshold,
		SpawnOnConfidence: pol.Honeypot.SpawnOnConfidence,
	})

	// redis tarpit
	rdb := tarpitpkg.NewRedisClient("localhost:6379")
	tarp := tarpitpkg.New(rdb,
		pol.Tarpit.Threshold,
		pol.Tarpit.BaseDelayMs,
		pol.Tarpit.GrowthFactor,
		pol.Tarpit.MaxDelayMs,
		pol.Tarpit.WindowSeconds,
	)

	// templates
	tloader, err := templatespkg.NewLoader(pol.Templates.Path, pol.Templates.HMACSecret)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load templates")
	}

	// honeypot manager
	hm := honeypot.NewManager(pol.Honeypot.Image, pol.Honeypot.TTLMinutes)

	// telemetry forwarder
	tf := telemetry.NewForwarder(pol.Aggregator.URL, pol.Aggregator.JWTSecret)

	// PROM metrics
	attacks := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "chameleon_attacks_total",
			Help: "Total detected attacks",
		},
		[]string{"label"},
	)
	prometheus.MustRegister(attacks)

	// admin + metrics HTTP server
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/admin/kill-switch", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			killSwitch = true
			w.Write([]byte("kill-switch enabled"))
			return
		}
		if r.Method == "DELETE" {
			killSwitch = false
			w.Write([]byte("kill-switch disabled"))
			return
		}
		w.WriteHeader(405)
	})

	go func() {
		addr := fmt.Sprintf(":%d", pol.Metrics.Port)
		logger.Info().Str("addr", addr).Msg("metrics/admin HTTP started")
		http.ListenAndServe(addr, nil)
	}()

	// gRPC on UDS
	_ = os.Remove(*socket)
	lis, err := net.Listen("unix", *socket)
	if err != nil {
		logger.Fatal().Err(err).Msg("can't listen on UDS")
	}

	grpcServer := grpc.NewServer()

	agentSrv := &server.AgentServer{
		Logger:        logger,
		Engine:        engine,
		Tarpit:        tarp,
		Templates:     tloader,
		Honeypot:      hm,
		Telemetry:     tf,
		AttacksMetric: attacks,
		KillSwitchPtr: &killSwitch,
	}

	pb.RegisterChameleonServer(grpcServer, agentSrv)
	reflection.Register(grpcServer)

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		<-sigc
		grpcServer.GracefulStop()
	}()

	logger.Info().Str("socket", *socket).Msg("gRPC starting")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal().Err(err).Msg("grpc serve failed")
	}
}
