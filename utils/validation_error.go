package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationErrorToString(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var messages []string
		for _, e := range ve {
			fieldName := e.Field()
			messages = append(messages, fmt.Sprintf("%s failed on '%s'.", fieldName, e.Tag()))
		}
		return strings.Join(messages, "; ")
	}

	return err.Error()
}
