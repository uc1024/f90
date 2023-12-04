package protojsonx

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uc1024/f90/core/coderx/lib/examplepb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var message = &examplepb.ABitOfEverything{
	SingleNested:        &examplepb.ABitOfEverything_Nested{},
	RepeatedStringValue: nil,
	MappedStringValue:   nil,
	MappedNestedValue:   nil,
	RepeatedEnumValue:   nil,
	TimestampValue:      &timestamppb.Timestamp{},
	Uuid:                "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",
	Nested: []*examplepb.ABitOfEverything_Nested{
		{
			Name:   "foo",
			Amount: 12345,
		},
	},
	Uint64Value: 0xFFFFFFFFFFFFFFFF,
	EnumValue:   examplepb.NumericEnum_ONE,
	OneofValue: &examplepb.ABitOfEverything_OneofString{
		OneofString: "bar",
	},
	MapValue: map[string]examplepb.NumericEnum{
		"a": examplepb.NumericEnum_ONE,
		"b": examplepb.NumericEnum_ZERO,
	},
}

func TestMarshal(t *testing.T) {
	// Marshal
	buffer, err := Marshal(message)
	if err != nil {
		t.Fatalf("Marshalling returned error: %s", err.Error())
	}

	// Unmarshal
	unmarshalled := &examplepb.ABitOfEverything{}
	err = Unmarshal(buffer, unmarshalled)
	if err != nil {
		t.Fatalf("Unmarshalling returned error: %s", err.Error())
	}

	if !proto.Equal(unmarshalled, message) {
		t.Errorf(
			"Unmarshalled didn't match original message: (original = %v) != (unmarshalled = %v)",
			unmarshalled,
			message,
		)
	}
}

func TestMarshal2(t *testing.T) {
	buf, err := Marshal(message)
	assert.NoError(t, err)
	reads := bytes.NewReader(buf)
	decoder := NewDecoder(reads)
	unmarshalled := &examplepb.ABitOfEverything{}
	err = decoder.Decode(unmarshalled)
	assert.NoError(t, err)
	if !proto.Equal(unmarshalled, message) {
		t.Errorf(
			"Unmarshalled didn't match original message: (original = %v) != (unmarshalled = %v)",
			unmarshalled,
			message,
		)
	}
}
