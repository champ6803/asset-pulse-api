package usecase

import (
	"fmt"
	"strings"
)

// GetMockLogoURL generates a mock logo URL based on the app name
// This creates a consistent, app-specific logo URL without requiring database storage
func GetMockLogoURL(appName string) *string {
	if appName == "" {
		return nil
	}

	// Normalize app name to lowercase and remove special characters
	normalized := strings.ToLower(appName)
	normalized = strings.ReplaceAll(normalized, " ", "-")
	normalized = strings.ReplaceAll(normalized, ".", "-")
	normalized = strings.ReplaceAll(normalized, "_", "-")

	// Generate logo URL - using logo.clearbit.com as a common service
	// Alternative patterns:
	// - UI Avatars: https://ui-avatars.com/api/?name={appName}&background=random
	// - Clearbit: https://logo.clearbit.com/{domain}
	// - Static CDN: https://cdn.example.com/app-logos/{normalized}.png
	logoURL := fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=random&size=128", appName)
	return &logoURL
}
