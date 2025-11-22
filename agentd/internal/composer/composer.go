package composer

import (
	"bytes"
	"text/template"
	"time"
	"math/rand"
)

func ComposeText(tmpl string, vars map[string]string) (string, error) {
	// basic template fill using text/template
	t, err := template.New("tpl").Funcs(template.FuncMap{
		"now": func() string { return time.Now().Format(time.RFC3339) },
		"randInt": func(max int) int { rand.Seed(time.Now().UnixNano()); return rand.Intn(max) },
	}).Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, vars); err != nil {
		return "", err
	}
	return buf.String(), nil
}
