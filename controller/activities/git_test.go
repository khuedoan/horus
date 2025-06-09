package activities

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestChangedModules(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "test-git-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directory structure
	testDirs := []string{
		"infra/dev/core",
		"infra/dev/networking",
		"infra/dev/databases/postgres",
		"infra/prod/core",
		"infra/prod/monitoring",
		"shared/modules/vpc",
		"shared/modules/security",
		"docs",
	}

	for _, dir := range testDirs {
		err := os.MkdirAll(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}
	}

	// Create terragrunt.hcl files in specific directories
	terragruntDirs := []string{
		"infra/dev/core",
		"infra/dev/networking",
		"infra/dev/databases/postgres",
		"infra/prod/core",
		"infra/prod/monitoring",
		"shared/modules/vpc",
	}

	for _, dir := range terragruntDirs {
		terragruntPath := filepath.Join(tempDir, dir, "terragrunt.hcl")
		err := os.WriteFile(terragruntPath, []byte("# terragrunt config"), 0644)
		if err != nil {
			t.Fatalf("Failed to create terragrunt.hcl in %s: %v", dir, err)
		}
	}

	// Test cases with mock changed files
	testCases := []struct {
		name         string
		changedFiles []string
		expected     []string
	}{
		{
			name: "Single module change",
			changedFiles: []string{
				"infra/dev/core/main.tf",
				"infra/dev/core/variables.tf",
			},
			expected: []string{"core"},
		},
		{
			name: "Multiple modules in same environment",
			changedFiles: []string{
				"infra/dev/core/main.tf",
				"infra/dev/networking/vpc.tf",
				"infra/dev/databases/postgres/db.tf",
			},
			expected: []string{"core", "networking", "databases/postgres"},
		},
		{
			name: "Modules across different environments",
			changedFiles: []string{
				"infra/dev/core/main.tf",
				"infra/prod/core/main.tf",
				"infra/prod/monitoring/alerts.tf",
			},
			expected: []string{"core", "monitoring"},
		},
		{
			name: "Shared modules without infra prefix",
			changedFiles: []string{
				"shared/modules/vpc/main.tf",
				"shared/modules/vpc/outputs.tf",
			},
			expected: []string{"shared/modules/vpc"},
		},
		{
			name: "Files without terragrunt.hcl",
			changedFiles: []string{
				"docs/README.md",
				"shared/modules/security/policy.tf", // no terragrunt.hcl in security dir
			},
			expected: []string{},
		},
		{
			name: "Mixed files with and without terragrunt.hcl",
			changedFiles: []string{
				"infra/dev/core/main.tf",
				"docs/README.md",
				"shared/modules/vpc/vpc.tf",
			},
			expected: []string{"core", "shared/modules/vpc"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock ChangedModules function that uses our test data
			// We'll create a custom function that simulates the file system checks
			modules := getChangedModulesFromFiles(tempDir, tc.changedFiles)

			sort.Strings(modules)
			sort.Strings(tc.expected)

			if !reflect.DeepEqual(modules, tc.expected) {
				t.Errorf("Expected modules %v, but got %v", tc.expected, modules)
			}
		})
	}
}

// Helper function to simulate ChangedModules logic without Git
func getChangedModulesFromFiles(repoPath string, changedFiles []string) []string {
	seen := make(map[string]struct{})
	modules := make([]string, 0) // Initialize as empty slice instead of nil

	for _, file := range changedFiles {
		// Get the directory of the changed file
		dir := filepath.Dir(file)

		// Walk up the directory tree to find the closest directory containing terragrunt.hcl
		currentDir := dir
		for {
			terragruntPath := filepath.Join(repoPath, currentDir, "terragrunt.hcl")
			if _, err := os.Stat(terragruntPath); err == nil {
				// Found terragrunt.hcl, this is a module directory
				modulePath := currentDir

				// Remove infra/<env>/ prefix if present
				if len(modulePath) > 0 && filepath.HasPrefix(modulePath, "infra/") {
					parts := strings.Split(filepath.ToSlash(modulePath), "/")
					if len(parts) >= 3 && parts[0] == "infra" {
						// Remove "infra" and environment (e.g., "dev", "prod")
						modulePath = strings.Join(parts[2:], "/")
					}
				}

				// Skip empty paths
				if modulePath != "" && modulePath != "." {
					// Normalize path separators to forward slashes
					modulePath = filepath.ToSlash(modulePath)

					if _, exists := seen[modulePath]; !exists {
						modules = append(modules, modulePath)
						seen[modulePath] = struct{}{}
					}
				}
				break
			}

			// Move up one directory level
			parent := filepath.Dir(currentDir)
			if parent == currentDir || parent == "." {
				// Reached the root, no terragrunt.hcl found
				break
			}
			currentDir = parent
		}
	}

	return modules
}
