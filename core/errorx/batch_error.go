package errorx

import "bytes"

type errorArray []error

// * 实现error接口
func (ea errorArray) Error() string {
	var buf bytes.Buffer

	// * 错误数组拼接成字符串 \n 分开
	for i := range ea {
		if i > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(ea[i].Error())
	}

	return buf.String()
}

// * batch error 
type BatchError struct {
	errors errorArray
}

// * add errors and nil error will be ignored
func (be *BatchError) Add(errs ...error) {
	for _, v := range errs {
		if v != nil {
			be.errors = append(be.errors, v)
		}
	}
}

// * returns error or one of errors
func (be *BatchError) Err() error {
	switch len(be.errors) {
	case 0:
		return nil
	case 1:
		return be.errors[0]
	default:
		return be.errors
	}
}

// * check errors array is not nil
func (be *BatchError) NotNil() bool {
	return len(be.errors) > 0
}
