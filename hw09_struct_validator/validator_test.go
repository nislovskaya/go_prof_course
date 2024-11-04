package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	t.Run("Empty struct", func(t *testing.T) {
		input := struct{}{}

		err := Validate(input)

		require.Nil(t, err)
	})

	t.Run("Valid User struct", func(t *testing.T) {
		input := User{
			ID:     "123456789012345678901234567890123456",
			Name:   "John Doe",
			Age:    25,
			Email:  "john.doe@example.com",
			Role:   UserRole("admin"),
			Phones: []string{"12345678901"},
		}

		actualErr := Validate(input)

		require.Nil(t, actualErr)
	})

	t.Run("Invalid User struct - Short ID and Age", func(t *testing.T) {
		input := User{
			ID:     "short_id",
			Name:   "Jane Doe",
			Age:    17,
			Email:  "jane.doe@",
			Role:   UserRole("user"),
			Phones: []string{"123456789"},
		}

		actualErr := Validate(input)

		var validationError ValidationErrors
		if !errors.As(actualErr, &validationError) {
			require.Fail(t, "unexpected return type")
		}

		require.Len(t, actualErr, 5)

		expectedErr := ValidationErrors{
			ValidationError{
				Field: "ID",
				Err:   fmt.Errorf("length must be 36"),
			},
			ValidationError{
				Field: "Age",
				Err:   fmt.Errorf("must be at least 18"),
			},
			ValidationError{
				Field: "Email",
				Err:   fmt.Errorf("does not match pattern ^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"),
			},
			ValidationError{
				Field: "Role",
				Err:   fmt.Errorf("must be one of [admin, stuff]"),
			},
			ValidationError{
				Field: "Phones[0]",
				Err:   fmt.Errorf("length must be 11"),
			},
		}

		require.EqualValues(t, actualErr, expectedErr)
	})

	t.Run("Valid App struct", func(t *testing.T) {
		input := App{
			Version: "1.0.0",
		}

		actualErr := Validate(input)

		require.Nil(t, actualErr)
	})

	t.Run("Invalid App struct - Version Length", func(t *testing.T) {
		input := App{
			Version: "1.0",
		}

		actualErr := Validate(input)

		var validationError ValidationErrors
		if !errors.As(actualErr, &validationError) {
			require.Fail(t, "unexpected return type")
		}

		require.Len(t, actualErr, 1)

		expectedErr := ValidationErrors{
			ValidationError{
				Field: "Version",
				Err:   fmt.Errorf("length must be 5"),
			},
		}

		require.EqualValues(t, actualErr, expectedErr)
	})

	t.Run("Valid Response struct", func(t *testing.T) {
		input := Response{
			Code: 200,
			Body: "Success",
		}

		err := Validate(input)

		require.Nil(t, err)
	})

	t.Run("Invalid Response Code", func(t *testing.T) {
		input := Response{
			Code: 403,
			Body: "Forbidden",
		}

		actualErr := Validate(input)

		var validationError ValidationErrors
		if !errors.As(actualErr, &validationError) {
			require.Fail(t, "unexpected return type")
		}

		require.Len(t, actualErr, 1)

		expectedErr := ValidationErrors{
			ValidationError{
				Field: "Code",
				Err:   fmt.Errorf("must be one of [200, 404, 500]"),
			},
		}

		require.EqualValues(t, actualErr, expectedErr)
	})

	t.Run("Valid Token struct", func(t *testing.T) {
		input := Token{
			Header:    []byte("header"),
			Payload:   []byte("payload"),
			Signature: []byte("signature"),
		}

		actualErr := Validate(input)

		require.Nil(t, actualErr)
	})
}
