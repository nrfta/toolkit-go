package shared

import (
	"fmt"

	"github.com/google/uuid"
)

// GenerateID generates a unique ID for testing
func GenerateID() string {
	return uuid.New().String()
}

// Ptr returns a pointer to the given value
func Ptr[T any](v T) *T {
	return &v
}

// Contains checks if a slice contains an element
func Contains[T comparable](slice []T, elem T) bool {
	for _, item := range slice {
		if item == elem {
			return true
		}
	}
	return false
}

// ExtractField extracts a field from a slice of structs
func ExtractIDs[T any](items []T, idFunc func(T) string) []string {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = idFunc(item)
	}
	return ids
}

// BuildTestName creates a test name with index
func BuildTestName(prefix string, index int) string {
	return fmt.Sprintf("%s-%d", prefix, index)
}
