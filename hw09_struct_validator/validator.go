package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var builder strings.Builder
	for _, err := range ve {
		builder.WriteString(fmt.Sprintf("Field '%s': %s\n", err.Field, err.Err))
	}
	return builder.String()
}

type Validator interface {
	Validate(value reflect.Value, fieldName string) []ValidationError
}

type LengthValidator struct {
	expectedLen int
}

func (lv LengthValidator) Validate(value reflect.Value, fieldName string) []ValidationError {
	var validationErrors []ValidationError

	//nolint:exhaustive
	switch value.Kind() {
	case reflect.String, reflect.Int:
		validationErrors = append(validationErrors, lv.validatePlainType(value, fieldName)...)
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			validationErrors = append(validationErrors, lv.Validate(value.Index(i), fmt.Sprintf("%s[%d]", fieldName, i))...)
		}
	default:
		validationErrors = append(validationErrors, ValidationError{
			Field: fieldName,
			Err:   errors.New("unsupported type for validation"),
		})
	}

	return validationErrors
}

func (lv LengthValidator) validatePlainType(value reflect.Value, fieldName string) []ValidationError {
	var validationErrors []ValidationError

	//nolint:exhaustive
	switch value.Kind() {
	case reflect.String, reflect.Int:
		if value.Len() != lv.expectedLen {
			validationErrors = append(validationErrors, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("length must be %d", lv.expectedLen),
			})
		}
	default:
		validationErrors = append(validationErrors, ValidationError{
			Field: fieldName,
			Err:   errors.New("unsupported type for validation"),
		})
	}

	return validationErrors
}

type RegexpValidator struct {
	pattern *regexp.Regexp
}

func (rv RegexpValidator) Validate(value reflect.Value, fieldName string) []ValidationError {
	var validationErrors []ValidationError

	//nolint:exhaustive
	switch value.Kind() {
	case reflect.String:
		if !rv.pattern.MatchString(value.String()) {
			validationErrors = append(validationErrors, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("does not match pattern %s", rv.pattern),
			})
		}
	default:
		validationErrors = append(validationErrors, ValidationError{
			Field: fieldName,
			Err:   errors.New("unsupported type for validation"),
		})
	}

	return validationErrors
}

type Option struct {
	Value string
}

type InValidator struct {
	options []Option
}

func (iv InValidator) Validate(value reflect.Value, fieldName string) []ValidationError {
	var validationErrors []ValidationError

	//nolint:exhaustive
	switch value.Kind() {
	case reflect.String:
		val := value.String()
		if !iv.contains(val) {
			validationErrors = append(validationErrors, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("must be one of [%s]", getOptionValues(iv.options)),
			})
		}
	case reflect.Int:
		numValue := int(value.Int())
		if !iv.containsInt(numValue) {
			validationErrors = append(validationErrors, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("must be one of [%s]", getOptionValues(iv.options)),
			})
		}
	default:
		validationErrors = append(validationErrors, ValidationError{
			Field: fieldName,
			Err:   errors.New("unsupported type for validation"),
		})
	}

	return validationErrors
}

func (iv InValidator) contains(val string) bool {
	for _, option := range iv.options {
		if option.Value == val {
			return true
		}
	}
	return false
}

func (iv InValidator) containsInt(num int) bool {
	for _, option := range iv.options {
		if numOption, err := strconv.Atoi(option.Value); err == nil && numOption == num {
			return true
		}
	}
	return false
}

func getOptionValues(options []Option) string {
	values := make([]string, len(options))
	for i, option := range options {
		values[i] = option.Value
	}
	return strings.Join(values, ", ")
}

type MinValidator struct {
	minValue int
}

func (mv MinValidator) Validate(value reflect.Value, fieldName string) []ValidationError {
	var validationErrors []ValidationError

	//nolint:exhaustive
	switch value.Kind() {
	case reflect.Int:
		if int(value.Int()) < mv.minValue {
			validationErrors = append(validationErrors, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("must be at least %d", mv.minValue),
			})
		}
	default:
		validationErrors = append(validationErrors, ValidationError{
			Field: fieldName,
			Err:   errors.New("unsupported type for validation"),
		})
	}

	return validationErrors
}

type MaxValidator struct {
	maxValue int
}

func (mv MaxValidator) Validate(value reflect.Value, fieldName string) []ValidationError {
	var validationErrors []ValidationError

	//nolint:exhaustive
	switch value.Kind() {
	case reflect.Int:
		if int(value.Int()) > mv.maxValue {
			validationErrors = append(validationErrors, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("must be no more than %d", mv.maxValue),
			})
		}
	default:
		validationErrors = append(validationErrors, ValidationError{
			Field: fieldName,
			Err:   errors.New("unsupported type for validation"),
		})
	}

	return validationErrors
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return errors.New("input is not a struct")
	}

	var validationErrors ValidationErrors

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		validators := createValidators(tag)

		for _, validator := range validators {
			validationErrors = append(validationErrors, validator.Validate(fieldValue, field.Name)...)
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func createValidators(tag string) []Validator {
	var validators []Validator

	validatorsTags := strings.Split(tag, "|")
	for _, validatorTag := range validatorsTags {
		switch {
		case strings.HasPrefix(validatorTag, "len:"):
			expectedLenStr := strings.TrimPrefix(validatorTag, "len:")
			if expectedLen, err := strconv.Atoi(expectedLenStr); err == nil {
				validators = append(validators, LengthValidator{expectedLen: expectedLen})
			}
		case strings.HasPrefix(validatorTag, "regexp:"):
			pattern := strings.TrimPrefix(validatorTag, "regexp:")
			re := regexp.MustCompile(pattern)
			validators = append(validators, RegexpValidator{pattern: re})
		case strings.HasPrefix(validatorTag, "in:"):
			optionsStr := strings.TrimPrefix(validatorTag, "in:")
			options := strings.Split(optionsStr, ",")
			optionList := make([]Option, len(options))
			for i, option := range options {
				optionList[i] = Option{Value: option}
			}
			validators = append(validators, InValidator{options: optionList})
		case strings.HasPrefix(validatorTag, "min:"):
			minValueStr := strings.TrimPrefix(validatorTag, "min:")
			if minValue, err := strconv.Atoi(minValueStr); err == nil {
				validators = append(validators, MinValidator{minValue: minValue})
			}
		case strings.HasPrefix(validatorTag, "max:"):
			maxValueStr := strings.TrimPrefix(validatorTag, "max:")
			if maxValue, err := strconv.Atoi(maxValueStr); err == nil {
				validators = append(validators, MaxValidator{maxValue: maxValue})
			}
		}
	}

	return validators
}
