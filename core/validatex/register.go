package validatex

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func RegisterDefaultValidators(v *validator.Validate, trans ut.Translator) error {
	return Register(defaultFieldValidators, v, trans)
}

func Register(validators []IFieldValidator, v *validator.Validate, trans ut.Translator) (err error) {

	if validators == nil || len(validators) == 0 {
		return
	}

	// * register validation
	for _, t := range validators {
		err = v.RegisterValidation(t.Tag(), t.RegisterFun())
		if err != nil {
			return err
		}
	}

	// * register translation
	if trans != nil {

		for _, t := range validators {

			if t.CustomTransFunc() != nil && t.CustomRegisFunc() != nil {

				err = v.RegisterTranslation(t.Tag(), trans, t.CustomRegisFunc(), t.CustomTransFunc())

			} else if t.CustomTransFunc() != nil && t.CustomRegisFunc() == nil {

				err = v.RegisterTranslation(t.Tag(), trans, registrationFunc(t.Tag(), t.Translation(), t.Override()), t.CustomTransFunc())

			} else if t.CustomTransFunc() == nil && t.CustomRegisFunc() != nil {

				err = v.RegisterTranslation(t.Tag(), trans, t.CustomRegisFunc(), translateFunc)

			} else {
				err = v.RegisterTranslation(t.Tag(), trans, registrationFunc(t.Tag(), t.Translation(), t.Override()), translateFunc)
			}

			if err != nil {
				return
			}
		}

	}

	return
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		if err = ut.Add(tag, translation, override); err != nil {
			return
		}
		return
	}
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		return fe.(error).Error()
	}
	return t
}
