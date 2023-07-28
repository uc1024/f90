package utilx

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
