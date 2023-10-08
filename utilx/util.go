package utilx

import "github.com/uc1024/f90/core/stringx"

func NoError(err error) {
	if err != nil {
		panic(err)
	}
}

func NoFalse(b bool) {
	if !b {
		panic("false")
	}
}

func MustSucc[T any](s T, err error) T {
	if err != nil {
		NoError(err)
	}
	return s
}

func RandomNum(size int) string {
	return stringx.Randn(size, "0123456789")
}
