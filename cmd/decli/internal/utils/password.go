// Package utils provides utility functions for the CLI including secure password handling.
package utils

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

// ReadPassword prompts for a password with the given prompt message and reads it securely.
// The password input is hidden from the terminal (no echo).
func ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)

	// Read password with terminal echo disabled
	passwordBytes, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	// Print newline after password input (since echo is disabled)
	fmt.Println()

	return string(passwordBytes), nil
}

// ReadPasswordWithConfirmation prompts for password twice to confirm.
func ReadPasswordWithConfirmation(prompt, confirmPrompt string) (string, error) {
	password, err := ReadPassword(prompt)
	if err != nil {
		return "", err
	}

	confirm, err := ReadPassword(confirmPrompt)
	if err != nil {
		return "", err
	}

	if password != confirm {
		return "", fmt.Errorf("passwords do not match")
	}

	return password, nil
}
