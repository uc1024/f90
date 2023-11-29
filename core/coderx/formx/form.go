package formx

import (
	"net/url"

	"github.com/go-playground/form/v4"
)

var decoder = form.NewDecoder()
var encoder = form.NewEncoder()

func init() {
	decoder.SetTagName("json")
	encoder.SetTagName("json")
}

func Decoder(req any, values url.Values) error {
	return decoder.Decode(req, values)
}

func Encoder(req any) (url.Values, error) {
	return encoder.Encode(req)
}
