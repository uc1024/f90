package threadingx

import "github.com/uc1024/f90/core/rescue"

func GoSafe(fn func()) {
	go RunSafe(fn)
}

func RunSafe(fn func()) {
	defer rescue.Catch()
	fn()
}
