package web

import (
	"encoding/json"
	"errors"
	"github.com/dimfeld/httptreemux"
	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	entranslation "github.com/go-playground/validator/v10/translations/en"
	"net/http"
	"reflect"
	"strings"
)

var validate *validator.Validate

var translator *ut.UniversalTranslator

func init() {
	validate = validator.New()

	enLocale := en.New()

	translator = ut.New(enLocale, enLocale)

	lang, _ := translator.GetTranslator("en")

	entranslation.RegisterDefaultTranslations(validate, lang)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func Params(r *http.Request) map[string]string {
	return httptreemux.ContextParams(r.Context())
}

func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(val); err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}

	if err := validate.Struct(val); err != nil {
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		lang, _ := translator.GetTranslator("en")

		var fields []FieldError
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Error(),
				Error: verror.Translate(lang),
			}
			fields = append(fields, field)
		}
		return &Error{
			Err:    errors.New("field validation error"),
			Status: http.StatusBadRequest,
			Fields: fields,
		}
	}

	return nil
}
