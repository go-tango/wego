package compress

import (
	"bytes"
	"html/template"
)

func parseTmpl(t *template.Template, data map[string]string) (string, error) {
	buf := bytes.NewBufferString("")
	err := t.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type argString []string

func (a argString) Get(i int, args ...string) (r string) {
	if i >= 0 && i < len(a) {
		r = a[i]
	} else if len(args) > 0 {
		r = args[0]
	}
	return
}
