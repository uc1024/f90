package formx

import (
	"io"
	"net/url"

	"github.com/go-playground/form/v4"
	"github.com/uc1024/f90/core/coderx/marshaler"
)

var (
	decoder = form.NewDecoder()
	encoder = form.NewEncoder()
)

func init() {
	decoder.SetTagName("json")
	encoder.SetTagName("json")
}

type FormUrl struct{}

// Marshal marshals "v" into byte sequence.
func (*FormUrl) Marshal(v interface{}) ([]byte, error) {
	return Marshal(v)
}

// Unmarshal unmarshals "data" into "v".
// "v" must be a pointer value.
func (*FormUrl) Unmarshal(data []byte, v interface{}) error {
	return Unmarshal(data, v)
}

// NewDecoder returns a Decoder which reads byte sequence from "r".
func (*FormUrl) NewDecoder(r io.Reader) marshaler.Decoder {
	return NewDecoder(r)
}

// NewEncoder returns an Encoder which writes bytes sequence into "w".
// ContentType returns the Content-Type which this marshaler is responsible for.
// The parameter describes the type which is being marshalled, which can sometimes
// affect the content type returned.
// ContentType(v interface{}) string
func (*FormUrl) NewEncoder(w io.Writer) marshaler.Encoder {
	return NewEncoder(w)
}

type DecoderWrapper struct {
	r       io.Reader
	decoder *form.Decoder
}

func (d *DecoderWrapper) Decode(v interface{}) error {
	by, err := io.ReadAll(d.r)
	if err != nil {
		return err
	}
	uv, err := url.ParseQuery(string(by))
	if err != nil {
		return err
	}
	return d.decoder.Decode(v, uv)
}

func Marshal(v interface{}) ([]byte, error) {
	uv, err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	return []byte(uv.Encode()), nil
}

func Unmarshal(data []byte, v interface{}) error {
	uv, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}
	return decoder.Decode(v, uv)
}

func NewDecoder(r io.Reader) marshaler.Decoder {
	return &DecoderWrapper{
		r:       r,
		decoder: decoder,
	}
}

func NewEncoder(w io.Writer) marshaler.Encoder {
	return marshaler.EncoderFunc(func(v interface{}) error {
		uv, err := encoder.Encode(v)
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(uv.Encode()))
		return err
	})
}
