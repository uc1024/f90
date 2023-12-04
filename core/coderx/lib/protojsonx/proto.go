package protojsonx

import (
	"errors"
	"fmt"
	"io"

	"github.com/uc1024/f90/core/coderx/marshaler"
	"google.golang.org/protobuf/proto"
)

var protoMarshaller = &ProtoMarshaller{}

// Marshal marshals "v" into byte sequence.
func Marshal(v interface{}) ([]byte, error) {
	m, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("v must be protojson.Marshaler")
	}
	return protoMarshaller.Marshal(m)
}

// Unmarshal unmarshals "data" into "v".
// "v" must be a pointer value.
func Unmarshal(data []byte, v interface{}) error {
	m, ok := v.(proto.Message)
	if !ok {
		return errors.New("unable to unmarshal non proto field")
	}
	return protoMarshaller.Unmarshal(data, m)
}

// NewDecoder returns a Decoder which reads byte sequence from "r".
func NewDecoder(r io.Reader) marshaler.Decoder {
	return marshaler.DecoderFunc(func(value interface{}) error {
		buffer, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		return protoMarshaller.Unmarshal(buffer, value)
	})
}

// NewEncoder returns an Encoder which writes bytes sequence into "w".
// ContentType returns the Content-Type which this marshaler is responsible for.
// The parameter describes the type which is being marshalled, which can sometimes
// affect the content type returned.
// ContentType(v interface{}) string
func NewEncoder(w io.Writer) marshaler.Encoder {
	return marshaler.EncoderFunc(func(value interface{}) error {
		buffer, err := protoMarshaller.Marshal(value)
		if err != nil {
			return err
		}
		if _, err := w.Write(buffer); err != nil {
			return err
		}
		return nil
	})
}
