package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ValidationErrorToString(err error) string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return err.Error()
	}

	messages := make([]string, 0, len(ve))
	for _, e := range ve {
		messages = append(messages, fmt.Sprintf("%s failed on '%s'.", e.Field(), e.Tag()))
	}

	return strings.Join(messages, "; ")
}
