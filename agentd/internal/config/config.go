package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type TarpitCfg struct {
	Threshold     int     `yaml:"threshold"`
	BaseDelayMs   int     `yaml:"base_delay_ms"`
	GrowthFactor  float64 `yaml:"growth_factor"`
	MaxDelayMs    int     `yaml:"max_delay_ms"`
	WindowSeconds int     `yaml:"window_seconds"`
}

type HoneypotCfg struct {
	TTLMinutes        int     `yaml:"ttl_minutes"`
	SpawnOnConfidence float64 `yaml:"spawn_on_confidence"`
	Image             string  `yaml:"image"`
}

type MLConfig struct {
	URL                 string  `yaml:"url"`
	TimeoutSeconds      int     `yaml:"timeout_seconds"`
	ConfidenceThreshold float64 `yaml:"confidence_threshold"`
}

type AggregatorCfg struct {
	URL       string `yaml:"url"`
	JWTSecret string `yaml:"jwt_secret"`
}

type TemplatesCfg struct {
	Path       string `yaml:"path"`
	HMACSecret string `yaml:"hmac_secret"`
}

type AdminCfg struct {
	KillSwitch bool `yaml:"kill_switch"`
}

type MetricsCfg struct {
	Port int `yaml:"port"`
}

type Policy struct {
	Tarpit     TarpitCfg     `yaml:"tarpit"`
	Honeypot   HoneypotCfg   `yaml:"honeypot"`
	ML         MLConfig      `yaml:"ml"`
	Aggregator AggregatorCfg `yaml:"aggregator"`
	Templates  TemplatesCfg  `yaml:"templates"`
	Admin      AdminCfg      `yaml:"admin"`
	Metrics    MetricsCfg    `yaml:"metrics"`
}

func LoadPolicy(path string) (*Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p Policy
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	// default metrics port
	if p.Metrics.Port == 0 {
		p.Metrics.Port = 9090
	}
	// defaults
	if p.Tarpit.WindowSeconds == 0 {
		p.Tarpit.WindowSeconds = 10
	}
	if p.ML.TimeoutSeconds == 0 {
		p.ML.TimeoutSeconds = 3
	}
	return &p, nil
}
