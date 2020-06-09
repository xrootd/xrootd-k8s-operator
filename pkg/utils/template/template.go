package template

import (
	"bytes"
	"path/filepath"
	"text/template"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("template")

func ApplyTemplate(tmplPath string, tmplData interface{}) (string, error) {
	tmpl, err := template.New(filepath.Base(tmplPath)).Funcs(TemplateFunctions).ParseFiles(tmplPath)
	if err != nil {
		log.Error(err, "Cannot open template file", "path", tmplPath)
		return "", err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, tmplData)
	if err != nil {
		log.Error(err, "Cannot apply template", "path", tmplPath)
	}
	return buf.String(), err
}
