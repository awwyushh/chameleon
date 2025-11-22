package rules

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

type CompiledRule struct {
	Rule    Rule
	Regex   *regexp.Regexp
}

type Engine struct {
	SQLi       []CompiledRule
	XSS        []CompiledRule
	Bruteforce []CompiledRule
}

// LoadRules loads all yaml files in ../rules directory
func LoadRules(path string) (*Engine, error) {
	engine := &Engine{}

	err := filepath.Walk(path, func(p string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(p) != ".yaml" {
			return nil
		}

		content, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}

		var rules []Rule
		if err := yaml.Unmarshal(content, &rules); err != nil {
			return fmt.Errorf("YAML unmarshal fail: %w", err)
		}

		for _, r := range rules {
			reg, err := regexp.Compile(r.Pattern)
			if err != nil {
				return fmt.Errorf("invalid regex in rule %s: %v", r.ID, err)
			}

			cr := CompiledRule{Rule: r, Regex: reg}

			switch r.Label {
			case "sqli":
				engine.SQLi = append(engine.SQLi, cr)
			case "xss":
				engine.XSS = append(engine.XSS, cr)
			case "bruteforce":
				engine.Bruteforce = append(engine.Bruteforce, cr)
			default:
				return fmt.Errorf("unknown rule label: %s", r.Label)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return engine, nil
}

// Match returns the first matching rule
func (e *Engine) Match(body string) (*CompiledRule, bool) {
	for _, r := range e.SQLi {
		if r.Regex.MatchString(body) {
			return &r, true
		}
	}
	for _, r := range e.XSS {
		if r.Regex.MatchString(body) {
			return &r, true
		}
	}
	for _, r := range e.Bruteforce {
		if r.Regex.MatchString(body) {
			return &r, true
		}
	}
	return nil, false
}
