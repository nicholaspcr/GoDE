package validation

import (
	"strings"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid username", "john_doe", false},
		{"valid with numbers", "user123", false},
		{"valid with hyphen", "john-doe", false},
		{"valid minimum length", "abc", false},
		{"valid maximum length", strings.Repeat("a", 64), false},
		{"too short", "ab", true},
		{"too long", strings.Repeat("a", 65), true},
		{"empty", "", true},
		{"with spaces", "john doe", true},
		{"with special chars", "john@doe", true},
		{"with dots", "john.doe", true},
		{"only whitespace", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid password", "ValidPass123!", false},
		{"minimum length", "12345678", false},
		{"maximum length", strings.Repeat("a", 72), false},
		{"too short", "short", true},
		{"empty", "", true},
		{"too long", strings.Repeat("a", 73), true},
		{"7 characters", "1234567", true},
		{"8 characters", "12345678", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid email", "user@example.com", false},
		{"valid with subdomain", "user@mail.example.com", false},
		{"valid with numbers", "user123@example.com", false},
		{"valid with dots", "first.last@example.com", false},
		{"valid with plus", "user+tag@example.com", false},
		{"invalid no @", "userexample.com", true},
		{"invalid no domain", "user@", true},
		{"invalid no TLD", "user@example", true},
		{"invalid no local part", "@example.com", true},
		{"empty", "", true},
		{"just @", "@", true},
		{"multiple @", "user@@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    *api.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &api.User{
				Ids:      &api.UserIDs{Username: "john_doe"},
				Email:    "john@example.com",
				Password: "validpass123",
			},
			wantErr: false,
		},
		{
			name:    "nil user",
			user:    nil,
			wantErr: true,
		},
		{
			name: "nil user IDs",
			user: &api.User{
				Ids:      nil,
				Email:    "john@example.com",
				Password: "validpass123",
			},
			wantErr: true,
		},
		{
			name: "invalid username",
			user: &api.User{
				Ids:      &api.UserIDs{Username: "ab"},
				Email:    "john@example.com",
				Password: "validpass123",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			user: &api.User{
				Ids:      &api.UserIDs{Username: "john_doe"},
				Email:    "invalid-email",
				Password: "validpass123",
			},
			wantErr: true,
		},
		{
			name: "empty email",
			user: &api.User{
				Ids:      &api.UserIDs{Username: "john_doe"},
				Email:    "",
				Password: "validpass123",
			},
			wantErr: true,
		},
		{
			name: "invalid password (too short)",
			user: &api.User{
				Ids:      &api.UserIDs{Username: "john_doe"},
				Email:    "john@example.com",
				Password: "short",
			},
			wantErr: true,
		},
		{
			name: "whitespace trimmed email",
			user: &api.User{
				Ids:      &api.UserIDs{Username: "john_doe"},
				Email:    "  john@example.com  ",
				Password: "validpass123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUser(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateStringLength(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		min     int
		max     int
		field   string
		wantErr bool
	}{
		{"within range", "hello", 1, 10, "test", false},
		{"exact min", "a", 1, 10, "test", false},
		{"exact max", "1234567890", 1, 10, "test", false},
		{"too short", "", 1, 10, "test", true},
		{"too long", "12345678901", 1, 10, "test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStringLength(tt.value, tt.min, tt.max, tt.field)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateNonEmpty(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		field   string
		wantErr bool
	}{
		{"non-empty", "hello", "test", false},
		{"empty", "", "test", true},
		{"whitespace only", "   ", "test", true},
		{"tab and space", "\t  \n", "test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNonEmpty(tt.value, tt.field)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
