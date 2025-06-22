package workflows

import (
	"errors"
	"testing"
	"time"

	"cloudlab/controller/activities"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type AppUpdateWorkflowTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *AppUpdateWorkflowTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
	s.env.SetTestTimeout(30 * time.Second)
}

func (s *AppUpdateWorkflowTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *AppUpdateWorkflowTestSuite) TestAppUpdate_Success() {
	input := AppUpdateInput{
		Url:       "https://github.com/example/cloudlab.git",
		Revision:  "main",
		Namespace: "khuedoan",
		App:       "blog",
		Cluster:   "production",
		NewImages: []activities.Image{
			{Repository: "docker.io/khuedoan/blog", Tag: "abc123def456789"},
		},
	}
	workspace := "/tmp/cloudlab-repos/abc123"

	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return(workspace, nil)
	s.env.OnActivity(activities.UpdateAppVersion, mock.Anything,
		workspace+"/apps", input.Namespace, input.App, input.Cluster, input.NewImages).Return(nil)
	s.env.OnActivity(activities.GitSync, mock.Anything,
		workspace+"/apps/khuedoan/blog/production.yaml").Return(nil)

	s.env.ExecuteWorkflow(AppUpdate, input)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *AppUpdateWorkflowTestSuite) TestAppUpdate_CloneFailure() {
	input := AppUpdateInput{
		Url:       "https://github.com/example/invalid-repo.git",
		Revision:  "main",
		Namespace: "test",
		App:       "app",
		Cluster:   "local",
		NewImages: []activities.Image{
			{Repository: "test/app", Tag: "v1.0.0"},
		},
	}

	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return("", errors.New("repository not found"))

	s.env.ExecuteWorkflow(AppUpdate, input)

	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
	s.Contains(s.env.GetWorkflowError().Error(), "failed to clone repository")
}

func (s *AppUpdateWorkflowTestSuite) TestAppUpdate_UpdateAppVersionFailure() {
	input := AppUpdateInput{
		Url:       "https://github.com/example/cloudlab.git",
		Revision:  "main",
		Namespace: "finance",
		App:       "actualbudget",
		Cluster:   "local",
		NewImages: []activities.Image{
			{Repository: "docker.io/actualbudget/actual-server", Tag: "25.7.0-alpine"},
		},
	}
	workspace := "/tmp/cloudlab-repos/def456"

	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return(workspace, nil)
	s.env.OnActivity(activities.UpdateAppVersion, mock.Anything,
		workspace+"/apps", input.Namespace, input.App, input.Cluster, input.NewImages).Return(
		errors.New("failed to read file: no such file or directory"))

	s.env.ExecuteWorkflow(AppUpdate, input)

	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
	s.Contains(s.env.GetWorkflowError().Error(), "failed to update app version")
}

func (s *AppUpdateWorkflowTestSuite) TestAppUpdate_GitSyncFailure() {
	input := AppUpdateInput{
		Url:       "https://github.com/example/cloudlab.git",
		Revision:  "main",
		Namespace: "khuedoan",
		App:       "notes",
		Cluster:   "production",
		NewImages: []activities.Image{
			{Repository: "ghcr.io/silverbulletmd/silverbullet", Tag: "v3"},
		},
	}
	workspace := "/tmp/cloudlab-repos/ghi789"

	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return(workspace, nil)
	s.env.OnActivity(activities.UpdateAppVersion, mock.Anything,
		workspace+"/apps", input.Namespace, input.App, input.Cluster, input.NewImages).Return(nil)
	s.env.OnActivity(activities.GitSync, mock.Anything,
		workspace+"/apps/khuedoan/notes/production.yaml").Return(
		errors.New("git push failed: authentication required"))

	s.env.ExecuteWorkflow(AppUpdate, input)

	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
	s.Contains(s.env.GetWorkflowError().Error(), "failed to sync changes to git")
}

func (s *AppUpdateWorkflowTestSuite) TestAppUpdate_MultipleImages() {
	input := AppUpdateInput{
		Url:       "https://github.com/example/cloudlab.git",
		Revision:  "develop",
		Namespace: "test",
		App:       "example",
		Cluster:   "local",
		NewImages: []activities.Image{
			{Repository: "zot.zot.svc.cluster.local/example-service", Tag: "newcommithash123"},
			{Repository: "docker.io/redis", Tag: "7.0-alpine"},
		},
	}
	workspace := "/tmp/cloudlab-repos/jkl012"

	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return(workspace, nil)
	s.env.OnActivity(activities.UpdateAppVersion, mock.Anything,
		workspace+"/apps", input.Namespace, input.App, input.Cluster, input.NewImages).Return(nil)
	s.env.OnActivity(activities.GitSync, mock.Anything,
		workspace+"/apps/test/example/local.yaml").Return(nil)

	s.env.ExecuteWorkflow(AppUpdate, input)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *AppUpdateWorkflowTestSuite) TestAppUpdate_RealWorldExample() {
	// Test with realistic data from the actual apps directory
	input := AppUpdateInput{
		Url:       "https://github.com/khuedoan/cloudlab.git",
		Revision:  "main",
		Namespace: "khuedoan",
		App:       "blog",
		Cluster:   "production",
		NewImages: []activities.Image{
			{Repository: "docker.io/khuedoan/blog", Tag: "1234567890abcdef1234567890abcdef12345678"},
		},
	}
	workspace := "/tmp/cloudlab-repos/realworld123"

	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return(workspace, nil)
	s.env.OnActivity(activities.UpdateAppVersion, mock.Anything,
		workspace+"/apps", input.Namespace, input.App, input.Cluster, input.NewImages).Return(nil)
	s.env.OnActivity(activities.GitSync, mock.Anything,
		workspace+"/apps/khuedoan/blog/production.yaml").Return(nil)

	s.env.ExecuteWorkflow(AppUpdate, input)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *AppUpdateWorkflowTestSuite) TestAppUpdate_ActivityTimeout() {
	input := AppUpdateInput{
		Url:       "https://github.com/example/cloudlab.git",
		Revision:  "main",
		Namespace: "test",
		App:       "slow-app",
		Cluster:   "production",
		NewImages: []activities.Image{
			{Repository: "test/slow-app", Tag: "v1.0.0"},
		},
	}

	// Simulate a timeout - we'll just return an error since the test timeout catches this
	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return("", errors.New("timeout"))

	s.env.ExecuteWorkflow(AppUpdate, input)

	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
}

func (s *AppUpdateWorkflowTestSuite) TestAppUpdate_EmptyImages() {
	input := AppUpdateInput{
		Url:       "https://github.com/example/cloudlab.git",
		Revision:  "main",
		Namespace: "test",
		App:       "app",
		Cluster:   "local",
		NewImages: []activities.Image{}, // Empty images array
	}
	workspace := "/tmp/cloudlab-repos/empty123"

	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return(workspace, nil)
	s.env.OnActivity(activities.UpdateAppVersion, mock.Anything,
		workspace+"/apps", input.Namespace, input.App, input.Cluster, input.NewImages).Return(nil)
	s.env.OnActivity(activities.GitSync, mock.Anything,
		workspace+"/apps/test/app/local.yaml").Return(nil)

	s.env.ExecuteWorkflow(AppUpdate, input)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *AppUpdateWorkflowTestSuite) TestAppUpdate_SpecialCharactersInPath() {
	input := AppUpdateInput{
		Url:       "https://github.com/example/cloudlab.git",
		Revision:  "feature/special-branch-name",
		Namespace: "test-namespace",
		App:       "app-with-dashes",
		Cluster:   "staging-env",
		NewImages: []activities.Image{
			{Repository: "registry.example.com/test/app-with-dashes", Tag: "v1.2.3-rc1"},
		},
	}
	workspace := "/tmp/cloudlab-repos/special456"

	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return(workspace, nil)
	s.env.OnActivity(activities.UpdateAppVersion, mock.Anything,
		workspace+"/apps", input.Namespace, input.App, input.Cluster, input.NewImages).Return(nil)
	s.env.OnActivity(activities.GitSync, mock.Anything,
		workspace+"/apps/test-namespace/app-with-dashes/staging-env.yaml").Return(nil)

	s.env.ExecuteWorkflow(AppUpdate, input)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func TestAppUpdateWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(AppUpdateWorkflowTestSuite))
}
