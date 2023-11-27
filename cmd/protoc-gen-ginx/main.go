package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	VERSION             = "0.0.1"
	FILE_SUFFIX         = ".ginx.pb.go"

	DEPRECATION_COMMENT = "// Deprecated: Do not use."
	PLUGINPB_FEATURE    = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
)

var args = NewDefaultArgs()

func init() {
	flag.BoolVar(&args.DisableClient, "disable_client", true, "disable use client")
}

func main() {
	flag.Parse()
	protogen.Options{ParamFunc: flag.CommandLine.Set}.Run(genProto)
}
