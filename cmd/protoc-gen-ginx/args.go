package main

type _args struct {
	DisableClient bool
}


func NewDefaultArgs() *_args {
	return &_args{
		DisableClient: true,
	}
}
