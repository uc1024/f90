package jsonpbx

import (
	"encoding/json"
	"io"

	"github.com/uc1024/f90/core/coderx/marshaler"
	"google.golang.org/protobuf/encoding/protojson"
)

var jsonpb = &JSONPb{
	MarshalOptions: protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	},
	UnmarshalOptions: protojson.UnmarshalOptions{
		DiscardUnknown: true,
	},
}

func Marshal(v interface{}) ([]byte, error) {
	return jsonpb.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return jsonpb.Unmarshal(data, v)
}

func NewDecoder(r io.Reader) marshaler.Decoder {
	return json.NewDecoder(r)
}

func NewEncoder(w io.Writer) marshaler.Encoder {
	return jsonpb.NewEncoder(w)
}
