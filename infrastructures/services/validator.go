package services

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en2 "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
	"reflect"
	"strings"
)

type Validation struct {
	validator  *validator.Validate
	translator ut.Translator
}

func NewValidation() *Validation {
	validate := validator.New()

	// Using translator to translate error message
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = en2.RegisterDefaultTranslations(validate, trans)

	// Overriding "eq" translation
	err := validate.RegisterTranslation("eq", trans, func(ut ut.Translator) error {
		return ut.Add("eq", "{0} doesn't match!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("eq", fe.Field())

		return t
	})

	if err != nil {
		panic(fmt.Errorf("new_validation: overriding eq trans: %v", err))
	}

	return &Validation{
		validator:  validate,
		translator: trans,
	}
}

func FieldValue(payload interface{}, field string) interface{} {
	r := reflect.ValueOf(payload)
	f := reflect.Indirect(r).FieldByName(field)

	return f.Interface()
}

func translateError(field string, err error, trans ut.Translator) string {
	if err == nil {
		return ""
	}

	var validatorErrors validator.ValidationErrors
	errors.As(err, &validatorErrors)

	for _, e := range validatorErrors {
		translated := e.Translate(trans)

		return fmt.Sprintf("%s:%s%s", field, field, translated)
	}

	return ""
}

func Validate(s interface{}, schema map[string]string, v *Validation) {
	var messages []string
	isError := false

	for field, rule := range schema {
		value := FieldValue(s, field)

		if err := v.validator.Var(value, rule); err != nil {
			translatedErr := translateError(field, err, v.translator)
			isError = true

			messages = append(messages, translatedErr)
		}
	}

	if isError {
		panic(fiber.NewError(
			fiber.StatusBadRequest,
			"Invalid request payload!\n"+
				strings.Join(messages, ";"),
		))
	}
}
