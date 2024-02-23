package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, maxLength int, minLength int) error {
	n := len(value)
	if n > maxLength || n < minLength {
		return fmt.Errorf("%s length must be between %d and %d", value, minLength, maxLength)
	}
	return nil
}

func ValidateUsername(username string) error {
	if err := ValidateString(username, 100, 3); err != nil {
		return err
	}

	if !isValidUsername(username) {
		return fmt.Errorf("%s is not a valid username. A valid username must contain only letters, numbers, and underscores", username)
	}
	return nil
}

func ValidatePassword(password string) error {
	return ValidateString(password, 100, 3)
}

func ValidateEmail(email string) error {
	if err := ValidateString(email, 200, 3); err != nil {
		return err
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("%s is not a valid email address", email)
	}
	return nil
}

func ValidateFullName(fullName string) error {
	if err := ValidateString(fullName, 200, 3); err != nil {
		return err
	}

	if !isValidFullName(fullName) {
		return fmt.Errorf("%s is not a valid full name. A valid full name must contain only letters and spaces", fullName)
	}

	return nil
}
