package validator

import (
	"log"
	"reflect"
	"strings"
	"ticket/pkg/errors"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	pv "github.com/go-playground/validator/v10"
)

type Validator interface {
	StrcutWithTranslateError(s interface{}) []error
}

type validator struct {
	*pv.Validate
	trans ut.Translator
}

func New(trans ut.Translator) Validator {
	v := pv.New()
	err := RegisterDefaultTranslations(v, trans)
	if err != nil {
		log.Fatal("validator: ", err)
	}
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := ""
		tagJSON := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		tagQuery := strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]
		if tagJSON != "" {
			name = strings.ReplaceAll(tagJSON, "-", "")
		} else if tagQuery != "" {
			name = strings.ReplaceAll(tagQuery, "-", "")
		}
		return name
	})
	return &validator{v, trans}
}

func NewTranslator() ut.Translator {
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	return trans
}

func (v *validator) StrcutWithTranslateError(s interface{}) []error {
	validatorErrs := v.Struct(s).(pv.ValidationErrors)
	if validatorErrs == nil {
		return nil
	}
	out := make([]error, 0)
	for _, e := range validatorErrs {
		out = append(out, errors.ErrValidation.SetError(e))
	}
	return out
}
