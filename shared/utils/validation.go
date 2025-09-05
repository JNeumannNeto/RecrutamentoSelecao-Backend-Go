package utils

import (
	"regexp"
	"strings"
)

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidPassword(password string) bool {
	return len(password) >= 6
}

func IsValidUUID(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return uuidRegex.MatchString(uuid)
}

func IsValidRole(role string) bool {
	validRoles := []string{"admin", "candidate"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

func IsValidJobStatus(status string) bool {
	validStatuses := []string{"open", "closed"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func IsValidApplicationStatus(status string) bool {
	validStatuses := []string{"applied", "reviewing", "interview", "rejected", "accepted"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func IsValidProficiencyLevel(level string) bool {
	validLevels := []string{"beginner", "intermediate", "advanced", "expert"}
	for _, validLevel := range validLevels {
		if level == validLevel {
			return true
		}
	}
	return false
}

func SanitizeString(input string) string {
	return strings.TrimSpace(input)
}

func IsEmptyOrWhitespace(input string) bool {
	return strings.TrimSpace(input) == ""
}

func ValidateRequiredFields(fields map[string]string) []string {
	var errors []string
	for fieldName, fieldValue := range fields {
		if IsEmptyOrWhitespace(fieldValue) {
			errors = append(errors, fieldName+" is required")
		}
	}
	return errors
}
