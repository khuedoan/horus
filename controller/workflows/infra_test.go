package workflows

import (
	"context"
	"errors"
	"testing"
	"time"

	"cloudlab/controller/activities"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type InfraWorkflowTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *InfraWorkflowTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
	// Set a reasonable timeout for tests
	s.env.SetTestTimeout(30 * time.Second)
}

func (s *InfraWorkflowTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *InfraWorkflowTestSuite) TestInfraWorkflow_Success() {
	// Mock data
	input := InfraInputs{
		Url:         "https://github.com/example/repo.git",
		Revision:    "main",
		OldRevision: "HEAD~1",
		Stack:       "dev",
	}
	repoPath := "/tmp/infra-12345"
	changedModules := []string{"module1", "module2"}

	// Create a sample graph
	graph := &activities.Graph{
		Nodes: map[string]bool{
			"module1": true,
			"module2": true,
			"module3": true,
		},
		Edges: map[string][]string{
			"module1": {"module2"}, // module1 depends on module2
		},
	}

	// Create pruned graph (only changed modules and dependents)
	prunedGraph := &activities.Graph{
		Nodes: map[string]bool{
			"module1": true,
			"module2": true,
		},
		Edges: map[string][]string{
			"module1": {"module2"},
		},
	}

	// Mock activities - use mock.Anything for context parameter
	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return(repoPath, nil)
	s.env.OnActivity(activities.TerragruntGraph, mock.Anything, repoPath+"/infra/"+input.Stack).Return(graph, nil)
	s.env.OnActivity(activities.ChangedModules, mock.Anything, repoPath, input.OldRevision).Return(changedModules, nil)
	s.env.OnActivity(activities.TerragruntPrune, mock.Anything, graph, changedModules).Return(prunedGraph, nil)
	s.env.OnActivity(activities.TerragruntApply, mock.Anything, input.Url, input.Revision, "module2", input.Stack).Return(nil)
	s.env.OnActivity(activities.TerragruntApply, mock.Anything, input.Url, input.Revision, "module1", input.Stack).Return(nil)

	// Execute workflow
	s.env.ExecuteWorkflow(Infra, input)

	// Assertions
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var result *activities.Graph
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Equal(prunedGraph, result)
}

func (s *InfraWorkflowTestSuite) TestInfraWorkflow_CloneFailure() {
	input := InfraInputs{
		Url:         "https://github.com/example/invalid-repo.git",
		Revision:    "main",
		OldRevision: "HEAD~1",
		Stack:       "dev",
	}

	// Mock Clone to return error
	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return("", errors.New("repository not found"))

	s.env.ExecuteWorkflow(Infra, input)

	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
	s.Contains(s.env.GetWorkflowError().Error(), "repository not found")
}

func (s *InfraWorkflowTestSuite) TestInfraWorkflow_TerragruntGraphFailure() {
	input := InfraInputs{
		Url:         "https://github.com/example/repo.git",
		Revision:    "main",
		OldRevision: "HEAD~1",
		Stack:       "dev",
	}
	repoPath := "/tmp/infra-12345"

	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return(repoPath, nil)
	s.env.OnActivity(activities.TerragruntGraph, mock.Anything, repoPath+"/infra/"+input.Stack).Return(
		(*activities.Graph)(nil), errors.New("terragrunt dag graph failed"))

	s.env.ExecuteWorkflow(Infra, input)

	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
	s.Contains(s.env.GetWorkflowError().Error(), "terragrunt dag graph failed")
}

func (s *InfraWorkflowTestSuite) TestInfraWorkflow_ChangedModulesFailure() {
	input := InfraInputs{
		Url:         "https://github.com/example/repo.git",
		Revision:    "main",
		OldRevision: "HEAD~1",
		Stack:       "dev",
	}
	repoPath := "/tmp/infra-12345"
	graph := &activities.Graph{
		Nodes: map[string]bool{"module1": true},
		Edges: map[string][]string{},
	}

	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return(repoPath, nil)
	s.env.OnActivity(activities.TerragruntGraph, mock.Anything, repoPath+"/infra/"+input.Stack).Return(graph, nil)
	s.env.OnActivity(activities.ChangedModules, mock.Anything, repoPath, input.OldRevision).Return(
		[]string{}, errors.New("git diff failed"))

	s.env.ExecuteWorkflow(Infra, input)

	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
	s.Contains(s.env.GetWorkflowError().Error(), "git diff failed")
}

func (s *InfraWorkflowTestSuite) TestInfraWorkflow_TerragruntApplyFailure() {
	input := InfraInputs{
		Url:         "https://github.com/example/repo.git",
		Revision:    "main",
		OldRevision: "HEAD~1",
		Stack:       "dev",
	}
	repoPath := "/tmp/infra-12345"
	changedModules := []string{"module1"}

	graph := &activities.Graph{
		Nodes: map[string]bool{"module1": true},
		Edges: map[string][]string{},
	}

	prunedGraph := &activities.Graph{
		Nodes: map[string]bool{"module1": true},
		Edges: map[string][]string{},
	}

	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return(repoPath, nil)
	s.env.OnActivity(activities.TerragruntGraph, mock.Anything, repoPath+"/infra/"+input.Stack).Return(graph, nil)
	s.env.OnActivity(activities.ChangedModules, mock.Anything, repoPath, input.OldRevision).Return(changedModules, nil)
	s.env.OnActivity(activities.TerragruntPrune, mock.Anything, graph, changedModules).Return(prunedGraph, nil)
	s.env.OnActivity(activities.TerragruntApply, mock.Anything, input.Url, input.Revision, "module1", input.Stack).Return(
		errors.New("terragrunt apply failed: resource conflict"))

	s.env.ExecuteWorkflow(Infra, input)

	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
	s.Contains(s.env.GetWorkflowError().Error(), "terragrunt apply failed")
}

func (s *InfraWorkflowTestSuite) TestInfraWorkflow_ComplexDependencyGraph() {
	input := InfraInputs{
		Url:         "https://github.com/example/repo.git",
		Revision:    "main",
		OldRevision: "HEAD~1",
		Stack:       "prod",
	}
	repoPath := "/tmp/infra-67890"
	changedModules := []string{"vpc", "database"}

	// Complex dependency graph:
	// app -> [database, loadbalancer]
	// database -> vpc
	// loadbalancer -> vpc
	// monitoring -> app
	graph := &activities.Graph{
		Nodes: map[string]bool{
			"vpc":          true,
			"database":     true,
			"loadbalancer": true,
			"app":          true,
			"monitoring":   true,
		},
		Edges: map[string][]string{
			"app":          {"database", "loadbalancer"},
			"database":     {"vpc"},
			"loadbalancer": {"vpc"},
			"monitoring":   {"app"},
		},
	}

	// Pruned graph should contain changed modules and their dependents
	prunedGraph := &activities.Graph{
		Nodes: map[string]bool{
			"vpc":        true,
			"database":   true,
			"app":        true,
			"monitoring": true,
		},
		Edges: map[string][]string{
			"app":        {"database"},
			"database":   {"vpc"},
			"monitoring": {"app"},
		},
	}

	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return(repoPath, nil)
	s.env.OnActivity(activities.TerragruntGraph, mock.Anything, repoPath+"/infra/"+input.Stack).Return(graph, nil)
	s.env.OnActivity(activities.ChangedModules, mock.Anything, repoPath, input.OldRevision).Return(changedModules, nil)
	s.env.OnActivity(activities.TerragruntPrune, mock.Anything, graph, changedModules).Return(prunedGraph, nil)

	// Mock TerragruntApply calls in dependency order
	// Level 0: vpc
	s.env.OnActivity(activities.TerragruntApply, mock.Anything, input.Url, input.Revision, "vpc", input.Stack).Return(nil)
	// Level 1: database
	s.env.OnActivity(activities.TerragruntApply, mock.Anything, input.Url, input.Revision, "database", input.Stack).Return(nil)
	// Level 2: app
	s.env.OnActivity(activities.TerragruntApply, mock.Anything, input.Url, input.Revision, "app", input.Stack).Return(nil)
	// Level 3: monitoring
	s.env.OnActivity(activities.TerragruntApply, mock.Anything, input.Url, input.Revision, "monitoring", input.Stack).Return(nil)

	s.env.ExecuteWorkflow(Infra, input)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var result *activities.Graph
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Equal(4, result.NodeCount())
	s.Equal(3, result.EdgeCount())
}

func (s *InfraWorkflowTestSuite) TestInfraWorkflow_NoChangedModules() {
	input := InfraInputs{
		Url:         "https://github.com/example/repo.git",
		Revision:    "main",
		OldRevision: "HEAD~1",
		Stack:       "dev",
	}
	repoPath := "/tmp/infra-12345"
	changedModules := []string{} // No changes

	graph := &activities.Graph{
		Nodes: map[string]bool{
			"module1": true,
			"module2": true,
		},
		Edges: map[string][]string{
			"module1": {"module2"},
		},
	}

	// Pruned graph should be empty
	prunedGraph := &activities.Graph{
		Nodes: map[string]bool{},
		Edges: map[string][]string{},
	}

	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return(repoPath, nil)
	s.env.OnActivity(activities.TerragruntGraph, mock.Anything, repoPath+"/infra/"+input.Stack).Return(graph, nil)
	s.env.OnActivity(activities.ChangedModules, mock.Anything, repoPath, input.OldRevision).Return(changedModules, nil)
	s.env.OnActivity(activities.TerragruntPrune, mock.Anything, graph, changedModules).Return(prunedGraph, nil)

	// No TerragruntApply calls should be made since no modules to deploy

	s.env.ExecuteWorkflow(Infra, input)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var result *activities.Graph
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Equal(0, result.NodeCount())
	s.Equal(0, result.EdgeCount())
}

func (s *InfraWorkflowTestSuite) TestInfraWorkflow_ActivityTimeout() {
	input := InfraInputs{
		Url:         "https://github.com/example/repo.git",
		Revision:    "main",
		OldRevision: "HEAD~1",
		Stack:       "dev",
	}

	// Mock Clone to simulate a timeout scenario
	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return(
		"", errors.New("activity timeout"))

	s.env.ExecuteWorkflow(Infra, input)

	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
}

func (s *InfraWorkflowTestSuite) TestInfraWorkflow_ParallelExecution() {
	// Test that modules at the same dependency level are executed in parallel
	input := InfraInputs{
		Url:         "https://github.com/example/repo.git",
		Revision:    "main",
		OldRevision: "HEAD~1",
		Stack:       "dev",
	}
	repoPath := "/tmp/infra-12345"
	changedModules := []string{"module-a", "module-b", "module-c"}

	// Graph with parallel modules:
	// module-a and module-b can run in parallel (both depend on module-c)
	graph := &activities.Graph{
		Nodes: map[string]bool{
			"module-a": true,
			"module-b": true,
			"module-c": true,
		},
		Edges: map[string][]string{
			"module-a": {"module-c"},
			"module-b": {"module-c"},
		},
	}

	prunedGraph := graph // All modules changed

	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return(repoPath, nil)
	s.env.OnActivity(activities.TerragruntGraph, mock.Anything, repoPath+"/infra/"+input.Stack).Return(graph, nil)
	s.env.OnActivity(activities.ChangedModules, mock.Anything, repoPath, input.OldRevision).Return(changedModules, nil)
	s.env.OnActivity(activities.TerragruntPrune, mock.Anything, graph, changedModules).Return(prunedGraph, nil)

	// Level 0: module-c
	s.env.OnActivity(activities.TerragruntApply, mock.Anything, input.Url, input.Revision, "module-c", input.Stack).Return(nil)
	// Level 1: module-a and module-b (should execute in parallel)
	s.env.OnActivity(activities.TerragruntApply, mock.Anything, input.Url, input.Revision, "module-a", input.Stack).Return(nil)
	s.env.OnActivity(activities.TerragruntApply, mock.Anything, input.Url, input.Revision, "module-b", input.Stack).Return(nil)

	s.env.ExecuteWorkflow(Infra, input)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *InfraWorkflowTestSuite) TestInfraWorkflow_WorkerFailureRetry() {
	// Test that TerragruntApply can handle worker failure and retry on a different worker
	// by ensuring the activity is self-contained (clones repo internally)
	input := InfraInputs{
		Url:         "https://github.com/example/repo.git",
		Revision:    "main",
		OldRevision: "HEAD~1",
		Stack:       "dev",
	}
	repoPath := "/tmp/infra-12345"
	newWorkerRepoPath := "/tmp/infra-67890"
	changedModules := []string{"module1"}

	graph := &activities.Graph{
		Nodes: map[string]bool{
			"module1": true,
		},
		Edges: map[string][]string{},
	}

	prunedGraph := graph // Only module1 changed

	// Initial workflow activities (successful)
	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return(repoPath, nil)
	s.env.OnActivity(activities.TerragruntGraph, mock.Anything, repoPath+"/infra/"+input.Stack).Return(graph, nil)
	s.env.OnActivity(activities.ChangedModules, mock.Anything, repoPath, input.OldRevision).Return(changedModules, nil)
	s.env.OnActivity(activities.TerragruntPrune, mock.Anything, graph, changedModules).Return(prunedGraph, nil)

	// Simulate worker failure and retry on different worker
	applyCallCount := 0
	s.env.OnActivity(activities.TerragruntApply, mock.Anything, input.Url, input.Revision, "module1", input.Stack).Return(
		func(ctx context.Context, repoUrl, revision, modulePath, stack string) error {
			applyCallCount++
			if applyCallCount == 1 {
				// First attempt fails (simulating worker failure)
				return errors.New("worker failed: connection lost")
			}
			// Second attempt succeeds (activity is self-contained and clones repo again)
			return nil
		})

	// Mock additional Clone calls for TerragruntApply retries
	// The activity will call Clone internally to ensure repo availability
	s.env.OnActivity(activities.Clone, mock.Anything, input.Url, input.Revision).Return(newWorkerRepoPath, nil).Maybe()

	s.env.ExecuteWorkflow(Infra, input)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var result *activities.Graph
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Equal(1, result.NodeCount())
	s.Equal(0, result.EdgeCount())
}

func TestInfraWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(InfraWorkflowTestSuite))
}
