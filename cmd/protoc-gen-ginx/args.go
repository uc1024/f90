package main

type _args struct {
	ShowVersion            bool
	Omitempty              bool
	AllowDeleteBody        bool
	AllowEmptyPatchBody    bool
	RpcMode                string
	UseEncoding            bool
	DisableErrorBadRequest bool
	DisableClient          bool
}

var args = NewDefaultArgs()

func NewDefaultArgs() *_args {
	return &_args{
		ShowVersion:            false,
		Omitempty:              true,
		AllowDeleteBody:        false,
		AllowEmptyPatchBody:    false,
		RpcMode:                "rpcx",
		UseEncoding:            false,
		DisableErrorBadRequest: false,
		DisableClient:          true,
	}
}
