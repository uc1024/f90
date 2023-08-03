package templates

import (
	"bytes"
	_ "embed"
	"html/template"

	"github.com/go-playground/validator/v10"
)

//go:embed code.html
var code_templates_str string
var codeTemplates *template.Template

type CodeTemplatesData struct {
	// Subject string `json:"subject" validate:"required"`
	Code string `json:"code" validate:"required"`
}

func GenerateCodeEmail(message *CodeTemplatesData) (buf []byte, err error) {
	if codeTemplates == nil {
		codeTemplates = template.Must(template.New("code_templates").Parse(code_templates_str))
	}

	err = validator.New().Struct(message)
	if err != nil {
		return nil, err
	}

	tpl := bytes.NewBuffer([]byte{})
	err = codeTemplates.Execute(tpl, message)
	if err != nil {
		return nil, err
	}

	buf = tpl.Bytes()
	return
}
