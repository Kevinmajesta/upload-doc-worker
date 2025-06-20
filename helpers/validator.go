package helpers

import (
	"fmt"
	"strings"

	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func FindPgError(err error) *pgconn.PgError {
	for err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return pgErr
		}
		err = errors.Unwrap(err)
	}
	return nil
}

func TranslateErrorMessage(err error) map[string]string {
	errorsMap := make(map[string]string)

	// Validator error
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			switch fieldError.Tag() {
			case "required":
				errorsMap[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errorsMap[field] = "Invalid email format"
			case "unique":
				errorsMap[field] = fmt.Sprintf("%s already exists", field)
			case "min":
				errorsMap[field] = fmt.Sprintf("%s must be at least %s characters", field, fieldError.Param())
			case "max":
				errorsMap[field] = fmt.Sprintf("%s must be at most %s characters", field, fieldError.Param())
			case "numeric":
				errorsMap[field] = fmt.Sprintf("%s must be a number", field)
			default:
				errorsMap[field] = "Invalid value"
			}
		}
	}

	if err != nil {
		pgErr := FindPgError(err)
		if pgErr != nil {
			if pgErr.Code == "23505" {
				if strings.Contains(pgErr.Message, "username") {
					errorsMap["username"] = "Username sudah dipakai"
				} else if strings.Contains(pgErr.Message, "email") {
					errorsMap["email"] = "Email sudah dipakai"
				} else {
					errorsMap["error"] = "Data sudah ada"
				}
			}
		} else if err == gorm.ErrRecordNotFound {
			errorsMap["error"] = "Record not found"
		}
	}

	return errorsMap
}

func IsDuplicateEntryError(err error) bool {
	return FindPgError(err) != nil && FindPgError(err).Code == "23505"
}
