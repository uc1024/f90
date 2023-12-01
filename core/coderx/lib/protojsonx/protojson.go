package protojsonx

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var encoder = protojson.MarshalOptions{
	UseProtoNames:  true, // 使用proto名称
	UseEnumNumbers: true, // 枚举值使用数字
}
var decoder = protojson.UnmarshalOptions{
	DiscardUnknown: true, // 忽略未知字段
}

func Decoder(bytes []byte, value any) error {
	message, ok := value.(proto.Message)
	if !ok {
		return fmt.Errorf("protojsonx: values must be proto.Message")
	}
	return decoder.Unmarshal(bytes, message)
}

func Encoder(value any) ([]byte, error) {
	message, ok := value.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("protojsonx: values must be proto.Message")
	}
	return encoder.Marshal(message)
}
