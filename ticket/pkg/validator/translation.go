package validator

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	pv "github.com/go-playground/validator/v10"
)

// RegisterDefaultTranslations registers a set of default translations
// for all built in tag's in validator; you may add your own as desired.
func RegisterDefaultTranslations(v *pv.Validate, trans ut.Translator) (err error) {
	translations := []struct {
		tag             string
		translation     string
		override        bool
		customRegisFunc pv.RegisterTranslationsFunc
		customTransFunc pv.TranslationFunc
	}{
		{
			tag:         "alphanumfullspecialchars",
			translation: "{0} can only contain alphanumeric characters and special characters $&+,:;=?@#|'<>.^*()%!-",
			override:    true,
		},
		{
			tag:         "alphanumspace",
			translation: "{0} can only contain alphanumeric characters and space",
			override:    true,
		},
		{
			tag:         "alphanumspecialchars",
			translation: "{0} can only contain alphanumeric characters and special characters !@#$&-.,+*_/",
			override:    true,
		},
		{
			tag:         "decimalplace",
			translation: "{0} exceeded the limit of {1} decimals",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				return fmt.Sprintf("%s exceeded the limit of %s decimals", fe.Field(), fe.Param())
			},
		},
		{
			tag:      "minnumeric",
			override: true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				return fmt.Sprintf("%s must be %s or greater", fe.Field(), fe.Param())
			},
		},
		{
			tag:      "maxnumeric",
			override: true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				return fmt.Sprintf("%s must be %s or less", fe.Field(), fe.Param())
			},
		},
		{
			tag:      "gtnumeric",
			override: true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				return fmt.Sprintf("%s must be greater than %s", fe.Field(), fe.Param())
			},
		},
		{
			tag:         "gtnumericfield",
			translation: "{0} must be greater than {1}",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltnumericfield",
			translation: "{0} must be less than {1}",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtdatetime_if",
			translation: "{0} is a required field and must be greater than {1} from now",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				params := strings.Split(fe.Param(), " ")
				if len(params)%3 != 0 {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				t, err := ut.T(fe.Tag(), fe.Field(), params[2])
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "required",
			translation: "{0} is a required field",
			override:    true,
		},
		{
			tag:         "required_if",
			translation: "{0} is a required field",
			override:    true,
		},
		{
			tag: "len",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("len-string", "{0} must be {1} in length", true); err != nil {
					return
				}

				if err = ut.AddCardinal("len-string-character", "{0} character", locales.PluralRuleOne, true); err != nil {
					return
				}

				if err = ut.AddCardinal("len-string-character", "{0} characters", locales.PluralRuleOther, true); err != nil {
					return
				}

				if err = ut.Add("len-number", "{0} must be equal to {1}", true); err != nil {
					return
				}

				if err = ut.Add("len-items", "{0} must contain {1}", true); err != nil {
					return
				}
				if err = ut.AddCardinal("len-items-item", "{0} item", locales.PluralRuleOne, true); err != nil {
					return
				}

				if err = ut.AddCardinal("len-items-item", "{0} items", locales.PluralRuleOther, true); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				var err error
				var t string

				var digits uint64
				var kind reflect.Kind

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err := strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					c, err = ut.C("len-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("len-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					c, err = ut.C("len-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("len-items", fe.Field(), c)

				default:
					t, err = ut.T("len-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "min",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("min-string", "{0} must be at least {1} in length", true); err != nil {
					return
				}

				if err = ut.AddCardinal("min-string-character", "{0} character", locales.PluralRuleOne, true); err != nil {
					return
				}

				if err = ut.AddCardinal("min-string-character", "{0} characters", locales.PluralRuleOther, true); err != nil {
					return
				}

				if err = ut.Add("min-number", "{0} must be {1} or greater", true); err != nil {
					return
				}

				if err = ut.Add("min-items", "{0} must contain at least {1}", true); err != nil {
					return
				}
				if err = ut.AddCardinal("min-items-item", "{0} item", locales.PluralRuleOne, true); err != nil {
					return
				}

				if err = ut.AddCardinal("min-items-item", "{0} items", locales.PluralRuleOther, true); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				var err error
				var t string

				var digits uint64
				var kind reflect.Kind

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err := strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					c, err = ut.C("min-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("min-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					c, err = ut.C("min-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("min-items", fe.Field(), c)

				default:
					t, err = ut.T("min-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "max",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("max-string", "{0} must be a maximum of {1} in length", true); err != nil {
					return
				}

				if err = ut.AddCardinal("max-string-character", "{0} character", locales.PluralRuleOne, true); err != nil {
					return
				}

				if err = ut.AddCardinal("max-string-character", "{0} characters", locales.PluralRuleOther, true); err != nil {
					return
				}

				if err = ut.Add("max-number", "{0} must be {1} or less", true); err != nil {
					return
				}

				if err = ut.Add("max-items", "{0} must contain at maximum {1}", true); err != nil {
					return
				}
				if err = ut.AddCardinal("max-items-item", "{0} item", locales.PluralRuleOne, true); err != nil {
					return
				}

				if err = ut.AddCardinal("max-items-item", "{0} items", locales.PluralRuleOther, true); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				var err error
				var t string

				var digits uint64
				var kind reflect.Kind

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err := strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					c, err = ut.C("max-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("max-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					c, err = ut.C("max-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("max-items", fe.Field(), c)

				default:
					t, err = ut.T("max-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "eq",
			translation: "{0} is not equal to {1}",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ne",
			translation: "{0} should not be equal to {1}",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "lt",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("lt-string", "{0} must be less than {1} in length", true); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-string-character", "{0} character", locales.PluralRuleOne, true); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-string-character", "{0} characters", locales.PluralRuleOther, true); err != nil {
					return
				}

				if err = ut.Add("lt-number", "{0} must be less than {1}", true); err != nil {
					return
				}

				if err = ut.Add("lt-items", "{0} must contain less than {1}", true); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-items-item", "{0} item", locales.PluralRuleOne, true); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-items-item", "{0} items", locales.PluralRuleOther, true); err != nil {
					return
				}

				if err = ut.Add("lt-datetime", "{0} must be less than the current Date & Time", true); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {
					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lt-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lt-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lt-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lt-items", fe.Field(), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
						goto END
					}

					t, err = ut.T("lt-datetime", fe.Field())

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("lt-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "lte",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("lte-string", "{0} must be at maximum {1} in length", true); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-string-character", "{0} character", locales.PluralRuleOne, true); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-string-character", "{0} characters", locales.PluralRuleOther, true); err != nil {
					return
				}

				if err = ut.Add("lte-number", "{0} must be {1} or less", true); err != nil {
					return
				}

				if err = ut.Add("lte-items", "{0} must contain at maximum {1}", true); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-items-item", "{0} item", locales.PluralRuleOne, true); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-items-item", "{0} items", locales.PluralRuleOther, true); err != nil {
					return
				}

				if err = ut.Add("lte-datetime", "{0} must be less than or equal to the current Date & Time", true); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {
					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lte-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lte-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lte-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lte-items", fe.Field(), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
						goto END
					}

					t, err = ut.T("lte-datetime", fe.Field())

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("lte-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "gt",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("gt-string", "{0} must be greater than {1} in length", true); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-string-character", "{0} character", locales.PluralRuleOne, true); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-string-character", "{0} characters", locales.PluralRuleOther, true); err != nil {
					return
				}

				if err = ut.Add("gt-number", "{0} must be greater than {1}", true); err != nil {
					return
				}

				if err = ut.Add("gt-items", "{0} must contain more than {1}", true); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-items-item", "{0} item", locales.PluralRuleOne, true); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-items-item", "{0} items", locales.PluralRuleOther, true); err != nil {
					return
				}

				if err = ut.Add("gt-datetime", "{0} must be greater than the current Date & Time", true); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {
					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gt-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gt-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gt-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gt-items", fe.Field(), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
						goto END
					}

					t, err = ut.T("gt-datetime", fe.Field())

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("gt-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "gte",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("gte-string", "{0} must be at least {1} in length", true); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-string-character", "{0} character", locales.PluralRuleOne, true); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-string-character", "{0} characters", locales.PluralRuleOther, true); err != nil {
					return
				}

				if err = ut.Add("gte-number", "{0} must be {1} or greater", true); err != nil {
					return
				}

				if err = ut.Add("gte-items", "{0} must contain at least {1}", true); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-items-item", "{0} item", locales.PluralRuleOne, true); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-items-item", "{0} items", locales.PluralRuleOther, true); err != nil {
					return
				}

				if err = ut.Add("gte-datetime", "{0} must be greater than or equal to the current Date & Time", true); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {
					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gte-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gte-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gte-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gte-items", fe.Field(), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
						goto END
					}

					t, err = ut.T("gte-datetime", fe.Field())

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("gte-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "eqfield",
			translation: "{0} must be equal to {1}",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "eqcsfield",
			translation: "{0} must be equal to {1}",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "necsfield",
			translation: "{0} cannot be equal to {1}",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtcsfield",
			translation: "{0} must be greater than {1}",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtecsfield",
			translation: "{0} must be greater than or equal to {1}",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltcsfield",
			translation: "{0} must be less than {1}",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltecsfield",
			translation: "{0} must be less than or equal to {1}",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "nefield",
			translation: "{0} cannot be equal to {1}",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtfield",
			translation: "{0} must be greater than {1}",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtefield",
			translation: "{0} must be greater than or equal to {1}",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltfield",
			translation: "{0} must be less than {1}",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltefield",
			translation: "{0} must be less than or equal to {1}",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "alpha",
			translation: "{0} can only contain alphabetic characters",
			override:    true,
		},
		{
			tag:         "alphanum",
			translation: "{0} can only contain alphanumeric characters",
			override:    true,
		},
		{
			tag:         "numeric",
			translation: "{0} must be a valid numeric value",
			override:    true,
		},
		{
			tag:         "number",
			translation: "{0} must be a valid number",
			override:    true,
		},
		{
			tag:         "hexadecimal",
			translation: "{0} must be a valid hexadecimal",
			override:    true,
		},
		{
			tag:         "hexcolor",
			translation: "{0} must be a valid HEX color",
			override:    true,
		},
		{
			tag:         "rgb",
			translation: "{0} must be a valid RGB color",
			override:    true,
		},
		{
			tag:         "rgba",
			translation: "{0} must be a valid RGBA color",
			override:    true,
		},
		{
			tag:         "hsl",
			translation: "{0} must be a valid HSL color",
			override:    true,
		},
		{
			tag:         "hsla",
			translation: "{0} must be a valid HSLA color",
			override:    true,
		},
		{
			tag:         "e164",
			translation: "{0} must be a valid E.164 formatted phone number",
			override:    true,
		},
		{
			tag:         "email",
			translation: "{0} must be a valid email address",
			override:    true,
		},
		{
			tag:         "url",
			translation: "{0} must be a valid URL",
			override:    true,
		},
		{
			tag:         "uri",
			translation: "{0} must be a valid URI",
			override:    true,
		},
		{
			tag:         "base64",
			translation: "{0} must be a valid Base64 string",
			override:    true,
		},
		{
			tag:         "contains",
			translation: "{0} must contain the text '{1}'",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "containsany",
			translation: "{0} must contain at least one of the following characters '{1}'",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludes",
			translation: "{0} cannot contain the text '{1}'",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludesall",
			translation: "{0} cannot contain any of the following characters '{1}'",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludesrune",
			translation: "{0} cannot contain the following '{1}'",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "isbn",
			translation: "{0} must be a valid ISBN number",
			override:    true,
		},
		{
			tag:         "isbn10",
			translation: "{0} must be a valid ISBN-10 number",
			override:    true,
		},
		{
			tag:         "isbn13",
			translation: "{0} must be a valid ISBN-13 number",
			override:    true,
		},
		{
			tag:         "uuid",
			translation: "{0} must be a valid UUID",
			override:    true,
		},
		{
			tag:         "uuid3",
			translation: "{0} must be a valid version 3 UUID",
			override:    true,
		},
		{
			tag:         "uuid4",
			translation: "{0} must be a valid version 4 UUID",
			override:    true,
		},
		{
			tag:         "uuid5",
			translation: "{0} must be a valid version 5 UUID",
			override:    true,
		},
		{
			tag:         "ulid",
			translation: "{0} must be a valid ULID",
			override:    true,
		},
		{
			tag:         "ascii",
			translation: "{0} must contain only ascii characters",
			override:    true,
		},
		{
			tag:         "printascii",
			translation: "{0} must contain only printable ascii characters",
			override:    true,
		},
		{
			tag:         "multibyte",
			translation: "{0} must contain multibyte characters",
			override:    true,
		},
		{
			tag:         "datauri",
			translation: "{0} must contain a valid Data URI",
			override:    true,
		},
		{
			tag:         "latitude",
			translation: "{0} must contain valid latitude coordinates",
			override:    true,
		},
		{
			tag:         "longitude",
			translation: "{0} must contain a valid longitude coordinates",
			override:    true,
		},
		{
			tag:         "ssn",
			translation: "{0} must be a valid SSN number",
			override:    true,
		},
		{
			tag:         "ipv4",
			translation: "{0} must be a valid IPv4 address",
			override:    true,
		},
		{
			tag:         "ipv6",
			translation: "{0} must be a valid IPv6 address",
			override:    true,
		},
		{
			tag:         "ip",
			translation: "{0} must be a valid IP address",
			override:    true,
		},
		{
			tag:         "cidr",
			translation: "{0} must contain a valid CIDR notation",
			override:    true,
		},
		{
			tag:         "cidrv4",
			translation: "{0} must contain a valid CIDR notation for an IPv4 address",
			override:    true,
		},
		{
			tag:         "cidrv6",
			translation: "{0} must contain a valid CIDR notation for an IPv6 address",
			override:    true,
		},
		{
			tag:         "tcp_addr",
			translation: "{0} must be a valid TCP address",
			override:    true,
		},
		{
			tag:         "tcp4_addr",
			translation: "{0} must be a valid IPv4 TCP address",
			override:    true,
		},
		{
			tag:         "tcp6_addr",
			translation: "{0} must be a valid IPv6 TCP address",
			override:    true,
		},
		{
			tag:         "udp_addr",
			translation: "{0} must be a valid UDP address",
			override:    true,
		},
		{
			tag:         "udp4_addr",
			translation: "{0} must be a valid IPv4 UDP address",
			override:    true,
		},
		{
			tag:         "udp6_addr",
			translation: "{0} must be a valid IPv6 UDP address",
			override:    true,
		},
		{
			tag:         "ip_addr",
			translation: "{0} must be a resolvable IP address",
			override:    true,
		},
		{
			tag:         "ip4_addr",
			translation: "{0} must be a resolvable IPv4 address",
			override:    true,
		},
		{
			tag:         "ip6_addr",
			translation: "{0} must be a resolvable IPv6 address",
			override:    true,
		},
		{
			tag:         "unix_addr",
			translation: "{0} must be a resolvable UNIX address",
			override:    true,
		},
		{
			tag:         "mac",
			translation: "{0} must contain a valid MAC address",
			override:    true,
		},
		{
			tag:         "unique",
			translation: "{0} must contain unique values",
			override:    true,
		},
		{
			tag:         "iscolor",
			translation: "{0} must be a valid color",
			override:    true,
		},
		{
			tag:         "oneof",
			translation: "{0} must be one of [{1}]",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				s, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}
				return s
			},
		},
		{
			tag:         "json",
			translation: "{0} must be a valid json string",
			override:    true,
		},
		{
			tag:         "jwt",
			translation: "{0} must be a valid jwt string",
			override:    true,
		},
		{
			tag:         "lowercase",
			translation: "{0} must be a lowercase string",
			override:    true,
		},
		{
			tag:         "uppercase",
			translation: "{0} must be an uppercase string",
			override:    true,
		},
		{
			tag:         "datetime",
			translation: "{0} does not match the {1} format",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "postcode_iso3166_alpha2",
			translation: "{0} does not match postcode format of {1} country",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "postcode_iso3166_alpha2_field",
			translation: "{0} does not match postcode format of country in {1} field",
			override:    true,
			customTransFunc: func(ut ut.Translator, fe pv.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "boolean",
			translation: "{0} must be a valid boolean value",
			override:    true,
		},
	}

	// register function to get tag name from json tags.
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// register translations
	for _, t := range translations {

		if t.customTransFunc != nil && t.customRegisFunc != nil {
			err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, t.customTransFunc)
		} else if t.customTransFunc != nil && t.customRegisFunc == nil {
			err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), t.customTransFunc)
		} else if t.customTransFunc == nil && t.customRegisFunc != nil {
			err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, translateFunc)
		} else {
			err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), translateFunc)
		}

		if err != nil {
			return err
		}
	}

	return err
}

func registrationFunc(tag string, translation string, override bool) pv.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		if err = ut.Add(tag, translation, override); err != nil {
			return err
		}

		return err
	}
}

func translateFunc(ut ut.Translator, fe pv.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		log.Printf("warning: error translating FieldError: %#v", fe)
		return fe.(error).Error()
	}

	return t
}
