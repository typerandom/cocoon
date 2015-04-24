package main

import (
	"errors"
	"reflect"
	"strconv"
	"unicode"
)

type UnsupportedTypeError struct {
	Type    reflect.Type
	Message string
}

func NewUnsupportedTypeError(validatorName string, value interface{}) *UnsupportedTypeError {
	valueType := reflect.TypeOf(value)
	return &UnsupportedTypeError{
		Type:    valueType,
		Message: "Validator with name '" + validatorName + "' on struct '{struct}' and field '{field}' is not supported.",
	}
}

func (this *UnsupportedTypeError) Error() string {
	return this.Message
}

type ValidatorContext struct {
	Value        interface{}
	IsNil        bool
	StopValidate bool
}

func NewValidatorContext(normalizedValue *NormalizedValue) *ValidatorContext {
	return &ValidatorContext{
		Value: normalizedValue.Value,
		IsNil: normalizedValue.IsNil,
	}
}

type ValidatorFilter func(context *ValidatorContext, options []string) error

func IsEmpty(context *ValidatorContext, options []string) error {
	if len(options) > 0 {
		return errors.New("Validator 'empty' does not support any arguments.")
	}

	switch typedValue := context.Value.(type) {
	case string:
		if context.IsNil || len(typedValue) == 0 {
			context.StopValidate = true
		}
		return nil
	case int64:
		if context.IsNil || typedValue == 0 {
			context.StopValidate = true
		}
		return nil
	default:
		if context.IsNil {
			context.StopValidate = true
			return nil
		}
	}

	return NewUnsupportedTypeError("empty", context.Value)
}

func IsNotEmpty(context *ValidatorContext, options []string) error {
	if len(options) > 0 {
		return errors.New("Validator 'not_empty' does not support any arguments.")
	}

	switch typedValue := context.Value.(type) {
	case string:
		if context.IsNil || len(typedValue) == 0 {
			return errors.New("{field} cannot be empty.")
		}
		return nil
	case int64:
		if context.IsNil || typedValue == 0 {
			return errors.New("{field} cannot be empty.")
		}
		return nil
	case float64:
		if context.IsNil || typedValue == 0 {
			return errors.New("{field} cannot be empty.")
		}
		return nil
	default:
		if context.IsNil {
			return errors.New("{field} cannot be empty.")
		}
	}

	return NewUnsupportedTypeError("not_empty", context.Value)
}

func IsMin(context *ValidatorContext, options []string) error {
	if len(options) != 1 {
		return errors.New("Validator 'min' requires a single argument.")
	}

	minValue, err := strconv.Atoi(options[0])

	if err != nil {
		return errors.New("Unable to parse 'min' validator value.")
	}

	switch typedValue := context.Value.(type) {
	case string:
		if context.IsNil || len(typedValue) < minValue {
			return errors.New("{field} cannot be shorter than " + strconv.Itoa(minValue) + " characters.")
		}
		return nil
	case int64:
		if context.IsNil || typedValue < int64(minValue) {
			return errors.New("{field} cannot be less than " + strconv.Itoa(minValue) + ".")
		}
		return nil
	case float64:
		if context.IsNil || typedValue < float64(minValue) {
			return errors.New("{field} cannot be less than " + strconv.Itoa(minValue) + ".")
		}
		return nil
	}

	return NewUnsupportedTypeError("min", context.Value)
}

func IsMax(context *ValidatorContext, options []string) error {
	if len(options) != 1 {
		return errors.New("Validator 'max' requires a single argument.")
	}

	minValue, err := strconv.Atoi(options[0])

	if err != nil {
		return errors.New("Unable to parse 'max' validator value.")
	}

	switch typedValue := context.Value.(type) {
	case string:
		if !context.IsNil && len(typedValue) > minValue {
			return errors.New("{field} is longer than " + strconv.Itoa(minValue) + " characters.")
		}
		return nil
	case int64:
		if !context.IsNil && typedValue > int64(minValue) {
			return errors.New("{field} cannot be greater than " + strconv.Itoa(minValue) + ".")
		}
		return nil
	case float64:
		if !context.IsNil && typedValue > float64(minValue) {
			return errors.New("{field} cannot be greater than " + strconv.Itoa(minValue) + ".")
		}
		return nil
	}

	return NewUnsupportedTypeError("max", context.Value)
}

func IsLowerCase(context *ValidatorContext, options []string) error {
	if len(options) > 0 {
		return errors.New("Validator 'lowercase' does not support any arguments.")
	}

	switch typedValue := context.Value.(type) {
	case string:
		if context.IsNil || len(typedValue) == 0 {
			return nil
		}

		for _, char := range typedValue {
			if unicode.IsLetter(char) && !unicode.IsLower(char) {
				return errors.New("{field} must be in lower case.")
			}
		}

		return nil
	}

	return NewUnsupportedTypeError("lowercase", context.Value)
}

func IsUpperCase(context *ValidatorContext, options []string) error {
	if len(options) > 0 {
		return errors.New("Validator 'uppercase' does not support any arguments.")
	}

	switch typedValue := context.Value.(type) {
	case string:
		if context.IsNil || len(typedValue) == 0 {
			return nil
		}

		for _, char := range typedValue {
			if unicode.IsLetter(char) && !unicode.IsUpper(char) {
				return errors.New("{field} must be in upper case.")
			}
		}

		return nil
	}

	return NewUnsupportedTypeError("uppercase", context.Value)
}

func IsNumeric(context *ValidatorContext, options []string) error {
	if len(options) > 0 {
		return errors.New("Validator 'numeric' does not support any arguments.")
	}

	switch typedValue := context.Value.(type) {
	case string:
		if context.IsNil || len(typedValue) == 0 {
			return errors.New("{field} must be numeric.")
		}

		value, err := strconv.ParseInt(typedValue, 10, 32)

		if err != nil {
			return errors.New("{field} must contain numbers only.")
		}

		context.Value = value

		return nil
	}

	return NewUnsupportedTypeError("numeric", context.Value)
}

/*
IsHex
IsType
IsISO8601
IsUnixTime
IsEmail
IsUrl
IsFilePath
IsType				type(string)
IsInteger
IsDecimal
IsIP
IsRegexMatch
IsUUID
IsNumeric*/

func registerDefaultValidators() {
	registerValidator("empty", IsEmpty)
	registerValidator("not_empty", IsNotEmpty)
	registerValidator("min", IsMin)
	registerValidator("max", IsMax)
	registerValidator("lowercase", IsLowerCase)
	registerValidator("uppercase", IsUpperCase)
	registerValidator("numeric", IsNumeric)
}
