package rules

type Rule struct {
	ID         string  `yaml:"id"`
	Label      string  `yaml:"label"`
	Confidence float64 `yaml:"confidence"`
	Pattern    string  `yaml:"pattern"`
}

type RuleSet struct {
	SQLi       []Rule
	XSS        []Rule
	Bruteforce []Rule
}
