package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

const version = "0.0.1"

func init() {
	flag.BoolVar(&args.ShowVersion, "version", false, "print the version and exit")
	flag.BoolVar(&args.Omitempty, "omitempty", true, "omit if google.api is empty")
	flag.BoolVar(&args.AllowDeleteBody, "allow_delete_body", false, "allow delete body")
	flag.BoolVar(&args.AllowEmptyPatchBody, "allow_empty_patch_body", false, "allow empty patch body")
	flag.StringVar(&args.RpcMode, "rpc_mode", "rpcx", "rpc mode, default use rpcx rpc, options: rpcx,official")
	flag.BoolVar(&args.UseEncoding, "use_encoding", false, "use the framework encoding")
	flag.BoolVar(&args.DisableErrorBadRequest, "disable_error_bad_request", false, "disable error bad request")
	flag.BoolVar(&args.DisableClient, "disable_client", true, "disable use client")
}

func main() {
	flag.Parse()
	if args.ShowVersion {
		fmt.Printf("protoc-gen-ginx %v\n", version)
		return
	}

	protogen.Options{ParamFunc: flag.CommandLine.Set}.Run(runProtoGen)
}
