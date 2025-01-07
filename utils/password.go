package utils

import (
	"bufio"
	"errors"
	"os"
	"unicode"

	"github.com/manifoldco/promptui"
	"github.com/spf13/pflag"
)

// PasswordPrompt securely prompts user for password input.
// Masks password input and validates password strength.
func PasswordPrompt() string {
	prompt := promptui.Prompt{
		Label:    "Password",
		Validate: validate,
		Mask:     ' ',
	}
	password, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}
	return password
}

// PrivateKeyPrompt securely prompts user for private key input.
// Masks input and performs basic validation on the key format.
func PrivateKeyPrompt() string {
	prompt := promptui.Prompt{
		Label:    "🔑 Private Key",
		Validate: validatePrivateKey,
		Mask:     ' ',
	}
	privateKey, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}
	return privateKey
}

// validate checks if password meets security requirements.
// Ensures password is not empty and meets strength criteria.
func validate(input string) error {
	if input == "" || !strongPassword(input) {
		return errors.New("enter a valid password")
	}
	return nil
}

// validatePrivateKey performs basic validation on private key input.
// Checks for empty input and basic format requirements.
func validatePrivateKey(input string) error {
	if input == "" {
		return errors.New("enter a valid private key")
	}
	return nil
}

// GetPasswordFromFile reads password from specified file path.
// Retrieves password from the first line of the file.
// Fatal error if file cannot be read or is empty.
func GetPasswordFromFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Getting password from the first line of file at described location")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		return scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return ""
}

// AssignPassword determines password source and retrieves password.
// Prioritizes password file if flag is set, otherwise prompts user.
// Warns about security implications of using password file.
func AssignPassword(flagSet *pflag.FlagSet) string {
	if UtilsInterface.IsFlagPassed("password") {
		log.Warn("Password flag is passed")
		log.Warn("This is a unsecure way to use lumino client")
		passwordPath, _ := flagSet.GetString("password")
		return GetPasswordFromFile(passwordPath)
	}
	return PasswordPrompt()
}

// strongPassword validates password strength against security criteria.
// Checks for minimum:
// - 8 characters length
// - One uppercase letter
// - One lowercase letter
// - One number
// - One special character
func strongPassword(input string) bool {
	l, u, n, s := 0, 0, 0, 0
	if len(input) >= 8 {
		for _, char := range input {
			switch {
			case unicode.IsUpper(char):
				u += 1
			case unicode.IsLower(char):
				l += 1
			case unicode.IsNumber(char):
				n += 1
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				s += 1
			}
		}
	}
	return (l >= 1 && u >= 1 && n >= 1 && s >= 1) && (l+u+n+s == len(input))
}
