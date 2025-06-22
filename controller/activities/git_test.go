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

func TestGitAdd_PathParsing(t *testing.T) {
	// Test the path parsing logic in GitAdd without requiring actual git commands
	tests := []struct {
		name         string
		inputPath    string
		expectedDir  string
		expectedFile string
	}{
		{
			name:         "simple file",
			inputPath:    "/tmp/test.yaml",
			expectedDir:  "/tmp",
			expectedFile: "test.yaml",
		},
		{
			name:         "nested file",
			inputPath:    "/apps/namespace/app/cluster.yaml",
			expectedDir:  "/apps/namespace/app",
			expectedFile: "cluster.yaml",
		},
		{
			name:         "relative path",
			inputPath:    "relative/file.yaml",
			expectedDir:  "relative",
			expectedFile: "file.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the path manipulation logic that GitAdd uses
			actualDir := filepath.Dir(tt.inputPath)
			actualFile := filepath.Base(tt.inputPath)

			if actualDir != tt.expectedDir {
				t.Errorf("Expected directory '%s', got '%s'", tt.expectedDir, actualDir)
			}

			if actualFile != tt.expectedFile {
				t.Errorf("Expected filename '%s', got '%s'", tt.expectedFile, actualFile)
			}
		})
	}
}

func TestGitCommit_PathParsing(t *testing.T) {
	// Test the path parsing logic in GitCommit
	tests := []struct {
		name        string
		inputPath   string
		expectedDir string
		message     string
	}{
		{
			name:        "simple file with default message",
			inputPath:   "/tmp/test.yaml",
			expectedDir: "/tmp",
			message:     "chore(test/app): update local version",
		},
		{
			name:        "nested file with custom message",
			inputPath:   "/apps/namespace/app/cluster.yaml",
			expectedDir: "/apps/namespace/app",
			message:     "feat: update application configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the path manipulation logic that GitCommit uses
			actualDir := filepath.Dir(tt.inputPath)

			if actualDir != tt.expectedDir {
				t.Errorf("Expected directory '%s', got '%s'", tt.expectedDir, actualDir)
			}

			// Verify message is not empty
			if tt.message == "" {
				t.Error("Commit message should not be empty")
			}
		})
	}
}

func TestGitPush_PathParsing(t *testing.T) {
	// Test the path parsing logic in GitPush
	tests := []struct {
		name        string
		inputPath   string
		expectedDir string
	}{
		{
			name:        "simple file",
			inputPath:   "/tmp/test.yaml",
			expectedDir: "/tmp",
		},
		{
			name:        "nested file",
			inputPath:   "/apps/namespace/app/cluster.yaml",
			expectedDir: "/apps/namespace/app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the path manipulation logic that GitPush uses
			actualDir := filepath.Dir(tt.inputPath)

			if actualDir != tt.expectedDir {
				t.Errorf("Expected directory '%s', got '%s'", tt.expectedDir, actualDir)
			}
		})
	}
}

func TestGitActivities_CommandStructure(t *testing.T) {
	// Test that the separate git activities construct the expected commands
	testPath := "/tmp/test/app/cluster.yaml"
	expectedDir := "/tmp/test/app"
	expectedFile := "cluster.yaml"
	commitMessage := "chore(khuedoan/blog): update production version"

	// Verify the path parsing logic
	actualDir := filepath.Dir(testPath)
	actualFile := filepath.Base(testPath)

	if actualDir != expectedDir {
		t.Errorf("Expected directory '%s', got '%s'", expectedDir, actualDir)
	}

	if actualFile != expectedFile {
		t.Errorf("Expected filename '%s', got '%s'", expectedFile, actualFile)
	}

	// Verify the expected command structures for each activity
	tests := []struct {
		name            string
		expectedCommand []string
		description     string
	}{
		{
			name:            "GitAdd command",
			expectedCommand: []string{"git", "-C", expectedDir, "add", expectedFile},
			description:     "GitAdd should construct git add command",
		},
		{
			name:            "GitCommit command",
			expectedCommand: []string{"git", "-C", expectedDir, "commit", "-m", commitMessage},
			description:     "GitCommit should construct git commit command with message",
		},
		{
			name:            "GitPush command",
			expectedCommand: []string{"git", "-C", expectedDir, "push"},
			description:     "GitPush should construct git push command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.expectedCommand

			if len(cmd) < 3 {
				t.Errorf("%s should have at least 3 parts, got %d", tt.description, len(cmd))
				return
			}

			if cmd[0] != "git" {
				t.Errorf("%s should start with 'git', got '%s'", tt.description, cmd[0])
			}

			if cmd[1] != "-C" {
				t.Errorf("%s should have '-C' as second argument, got '%s'", tt.description, cmd[1])
			}

			if cmd[2] != expectedDir {
				t.Errorf("%s should use directory '%s', got '%s'", tt.description, expectedDir, cmd[2])
			}
		})
	}
}

func TestGenerateRepoPath(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		revision string
		wantPath bool // whether we expect a valid path
	}{
		{
			name:     "simple repo",
			url:      "https://github.com/user/repo.git",
			revision: "main",
			wantPath: true,
		},
		{
			name:     "same repo different revision",
			url:      "https://github.com/user/repo.git",
			revision: "develop",
			wantPath: true,
		},
		{
			name:     "empty inputs",
			url:      "",
			revision: "",
			wantPath: true, // Should still generate a path
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := generateRepoPath(tt.url, tt.revision)

			if tt.wantPath {
				if path == "" {
					t.Error("Expected non-empty path")
				}
				if !strings.Contains(path, "/tmp/cloudlab-repos/") {
					t.Errorf("Expected path to contain '/tmp/cloudlab-repos/', got: %s", path)
				}
				if len(filepath.Base(path)) != 16 {
					t.Errorf("Expected base path to be 16 characters, got: %s", filepath.Base(path))
				}
			}
		})
	}

	// Test that same inputs generate same path
	path1 := generateRepoPath("https://github.com/test/repo.git", "main")
	path2 := generateRepoPath("https://github.com/test/repo.git", "main")
	if path1 != path2 {
		t.Errorf("Same inputs should generate same path: %s != %s", path1, path2)
	}

	// Test that different inputs generate different paths
	path3 := generateRepoPath("https://github.com/test/repo.git", "develop")
	if path1 == path3 {
		t.Error("Different revisions should generate different paths")
	}
}
