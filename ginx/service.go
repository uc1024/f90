package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/form/v4"
)

type (
	Service struct {
		formDecode *form.Decoder
	}
)

func (service *Service) Form(ctx *gin.Context, req any) {
	// service.formDecode.Decode(req, ctx.Request.Form)
}

// form() 获取特定数据
// validate() 校验数据
// renderError() 错误的返回
// renderSuccess() 正确的返回

// onSuccess
// onError
