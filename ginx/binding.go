package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	JSON          = binding.JSON
	XML           = binding.XML
	Form          = binding.Form
	Query         = queryBinding{}
	FormPost      = binding.FormPost
	FormMultipart = binding.FormMultipart
	ProtoBuf      = binding.ProtoBuf
	MsgPack       = binding.MsgPack
	YAML          = binding.YAML
	// Uri           = binding.Uri
	Uri    = uriBinding{}
	Header = binding.Header
	TOML   = binding.TOML
)

var _validator = validator.New()

// Default returns the appropriate Binding instance based on the HTTP method
// and the content type.
func BindingDefault(method, contentType string) binding.Binding {
	if method == http.MethodGet {
		return Form
	}

	switch contentType {
	case binding.MIMEJSON:
		return JSON
	case binding.MIMEXML, binding.MIMEXML2:
		return XML
	case binding.MIMEPROTOBUF:
		return ProtoBuf
	case binding.MIMEMSGPACK, binding.MIMEMSGPACK2:
		return MsgPack
	case binding.MIMEYAML:
		return YAML
	case binding.MIMETOML:
		return TOML
	case binding.MIMEMultipartPOSTForm:
		return FormMultipart
	default: // case MIMEPOSTForm:
		return Form
	}
}

func Bind(c *gin.Context, obj any) error {
	return c.ShouldBindWith(obj, BindingDefault(c.Request.Method, c.ContentType()))
}

func BindUri(c *gin.Context, obj any) error {
	m := make(map[string][]string)
	for _, v := range c.Params {
		m[v.Key] = []string{v.Value}
	}
	return Uri.BindUri(m, obj)
}

func BindQuery(c *gin.Context, obj any) error {
	return c.ShouldBindWith(obj, Query)
}

func Validate() *validator.Validate {
	return _validator
}

type uriBinding struct{}

func (uriBinding) Name() string {
	return "uri"
}

func (uriBinding) BindUri(m map[string][]string, obj any) error {
	if err := binding.MapFormWithTag(obj, m, "json"); err != nil {
		return err
	}
	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}

type queryBinding struct{}

func (queryBinding) Name() string {
	return "query"
}

func (queryBinding) Bind(req *http.Request, obj any) error {
	values := req.URL.Query()
	if err := binding.MapFormWithTag(obj, values, "json"); err != nil {
		return err
	}
	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}
