package wrap

import (
	"encoding/json"
	"net/url"

	"github.com/go-playground/form/v4"
	"github.com/uc1024/f90/core/coderx/formx"
	"google.golang.org/protobuf/encoding/protojson"
)

type (
	Coder struct {
		formDecode  *form.Decoder
		formEncoder *form.Encoder

		jsonDecode  *json.Decoder
		jsonEncoder *json.Encoder

		protoDecode  *protojson.UnmarshalOptions
		protoEncoder *protojson.MarshalOptions
	}
)

func (coder *Coder) Decoder(req any, values url.Values) error {
	if coder.formDecode == nil {
		return formx.Decoder(req, values)
	}
	return coder.formDecode.Decode(req, values)
}

func (coder *Coder) Encoder(req any) (url.Values, error) {
	if coder.formEncoder == nil {
		return formx.Encoder(req)
	}
	return coder.formEncoder.Encode(req)
}
