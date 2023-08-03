package validatex

import (
	"fmt"
	"strings"

	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"
)

var defualtZhTrans = func() ut.Translator {
	z := zhongwen.New()
	uni := ut.New(z, z)
	trans, b := uni.GetTranslator("zh")
	if !b {
		panic("not zh")
	}
	return trans
}()

var DefualtZhTrans = defualtZhTrans

func RegisterDefaultTranslations(valid *validator.Validate) error {
	err := zh.RegisterDefaultTranslations(valid, defualtZhTrans)
	return err
}

type TranslateError struct{}

func (t TranslateError) Translate(err error) error {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return errs
	}

	detail := []string{}
	for k, _ := range errs {
		item := errs[k]
		detail = append(detail, item.Translate(defualtZhTrans))
	}

	return fmt.Errorf("%s", strings.Join(detail, "\n"))
}
