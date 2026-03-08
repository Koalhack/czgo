// Package template/message provides method for commit msg format/render.
package template

import (
	"bytes"
	"text/template"
)

func RenderCommitMsg(tmplStr string, data any) (string, error) {
	tmpl, err := template.New("message").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
