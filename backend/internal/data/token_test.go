package data_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/validator"
)

func TestValidateToken(t *testing.T) {
	tokens := []string{
		"",
		"shortToken",
		strings.Repeat("a", 27),
	}

	for _, token := range tokens {
		t.Run(fmt.Sprintf("Incorrect token %s", token), func(t *testing.T) {
			v := validator.New()
			data.ValidateTokenPlaintext(v, token)
			if len(v.Errors) == 0 {
				t.Fatalf("Expected to have validation errors")
			}
		})
	}

	t.Run("validates correct token", func(t *testing.T) {
		v := validator.New()
		data.ValidateTokenPlaintext(v, strings.Repeat("a", 26))
		if len(v.Errors) != 0 {
			t.Fatalf("Expected to have no validation errors, but got %v", v.Errors)
		}
	})
}
