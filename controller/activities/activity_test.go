package activities

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type ActivityTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestActivityEnvironment
}

func (s *ActivityTestSuite) SetupTest() {
	s.env = s.NewTestActivityEnvironment()
	s.env.SetTestTimeout(30 * time.Second)
}

// Test graph pruning logic with direct calls
func (s *ActivityTestSuite) TestPruneGraph_Success() {
	ctx := context.Background()
	originalGraph := &Graph{
		Nodes: map[string]bool{
			"vpc":      true,
			"database": true,
			"app":      true,
		},
		Edges: map[string][]string{
			"database": {"vpc"},
			"app":      {"database"},
		},
	}
	changedFiles := []string{"database"}

	prunedGraph, err := PruneGraph(ctx, originalGraph, changedFiles)

	s.NoError(err)
	s.Equal(2, prunedGraph.NodeCount()) // database and app (which depends on database)
	s.Equal(1, prunedGraph.EdgeCount()) // app -> database
	s.True(prunedGraph.Nodes["database"])
	s.True(prunedGraph.Nodes["app"])
	s.False(prunedGraph.Nodes["vpc"]) // vpc should be pruned as it's not changed and no dependents
}

func (s *ActivityTestSuite) TestPruneGraph_EmptyChanges() {
	ctx := context.Background()
	originalGraph := &Graph{
		Nodes: map[string]bool{
			"vpc":      true,
			"database": true,
		},
		Edges: map[string][]string{
			"database": {"vpc"},
		},
	}
	changedFiles := []string{} // No changes

	prunedGraph, err := PruneGraph(ctx, originalGraph, changedFiles)

	s.NoError(err)
	s.Equal(0, prunedGraph.NodeCount())
	s.Equal(0, prunedGraph.EdgeCount())
}

func (s *ActivityTestSuite) TestPruneGraph_ComplexDependencies() {
	ctx := context.Background()
	// Complex graph: monitoring -> app -> [database, cache] -> vpc
	originalGraph := &Graph{
		Nodes: map[string]bool{
			"vpc":        true,
			"database":   true,
			"cache":      true,
			"app":        true,
			"monitoring": true,
		},
		Edges: map[string][]string{
			"database":   {"vpc"},
			"cache":      {"vpc"},
			"app":        {"database", "cache"},
			"monitoring": {"app"},
		},
	}
	changedFiles := []string{"database"} // Only database changed

	prunedGraph, err := PruneGraph(ctx, originalGraph, changedFiles)

	s.NoError(err)
	// Should include: database (changed), app (depends on database), monitoring (depends on app)
	s.Equal(3, prunedGraph.NodeCount())
	s.True(prunedGraph.Nodes["database"])
	s.True(prunedGraph.Nodes["app"])
	s.True(prunedGraph.Nodes["monitoring"])
	s.False(prunedGraph.Nodes["vpc"])   // not a dependent
	s.False(prunedGraph.Nodes["cache"]) // not a dependent
}

// Test using TestActivityEnvironment for activities that need proper context
func (s *ActivityTestSuite) TestTerragruntPrune_WithActivityEnvironment() {
	originalGraph := &Graph{
		Nodes: map[string]bool{
			"vpc":      true,
			"database": true,
			"app":      true,
		},
		Edges: map[string][]string{
			"database": {"vpc"},
			"app":      {"database"},
		},
	}
	changedFiles := []string{"database"}

	// Register the activity first
	s.env.RegisterActivity(TerragruntPrune)

	// Use the activity environment to execute the activity
	val, err := s.env.ExecuteActivity(TerragruntPrune, originalGraph, changedFiles)
	s.NoError(err)

	var result *Graph
	err = val.Get(&result)
	s.NoError(err)
	s.Equal(2, result.NodeCount())
	s.Equal(1, result.EdgeCount())
}

func TestActivityTestSuite(t *testing.T) {
	suite.Run(t, new(ActivityTestSuite))
}

// Additional comprehensive unit tests for graph functions
func TestNewGraphFromDot_EmptyGraph(t *testing.T) {
	dotString := `digraph {
	}`

	graph, err := NewGraphFromDot(dotString)

	assert.NoError(t, err)
	assert.Equal(t, 0, graph.NodeCount())
	assert.Equal(t, 0, graph.EdgeCount())
}

func TestNewGraphFromDot_InvalidFormat(t *testing.T) {
	dotString := `not a valid dot format`

	graph, err := NewGraphFromDot(dotString)

	assert.NoError(t, err) // Should not error, just ignore invalid lines
	assert.Equal(t, 0, graph.NodeCount())
}

func TestGraph_TopologicalSort_CyclicGraph(t *testing.T) {
	// Create a graph with a cycle: A -> B -> C -> A
	graph := &Graph{
		Nodes: map[string]bool{
			"a": true,
			"b": true,
			"c": true,
		},
		Edges: map[string][]string{
			"a": {"b"},
			"b": {"c"},
			"c": {"a"}, // Creates cycle
		},
	}

	levels := graph.TopologicalSort()

	// Should handle cycles gracefully by putting remaining nodes in final level
	assert.Greater(t, len(levels), 0)

	// All nodes should be present somewhere
	allNodes := make(map[string]bool)
	for _, level := range levels {
		for _, node := range level {
			allNodes[node] = true
		}
	}
	assert.Len(t, allNodes, 3)
	assert.True(t, allNodes["a"])
	assert.True(t, allNodes["b"])
	assert.True(t, allNodes["c"])
}

func TestExtractQuotedString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`"hello world"`, "hello world"},
		{`""`, ""},
		{`hello`, ""},            // No quotes
		{`"hello`, ""},           // Missing closing quote
		{`hello"`, ""},           // Missing opening quote
		{`  "hello"  `, "hello"}, // With whitespace
	}

	for _, test := range tests {
		result := extractQuotedString(test.input)
		assert.Equal(t, test.expected, result, "Input: %s", test.input)
	}
}

func TestGraph_AddEdge_CreatesNodes(t *testing.T) {
	graph := NewGraph()

	graph.AddEdge("a", "b")

	assert.Equal(t, 2, graph.NodeCount())
	assert.Equal(t, 1, graph.EdgeCount())
	assert.True(t, graph.Nodes["a"])
	assert.True(t, graph.Nodes["b"])
	assert.Contains(t, graph.Edges["a"], "b")
}

func TestGraph_GetNodes(t *testing.T) {
	graph := &Graph{
		Nodes: map[string]bool{
			"vpc":      true,
			"database": true,
			"app":      true,
		},
		Edges: map[string][]string{},
	}

	nodes := graph.GetNodes()

	assert.Len(t, nodes, 3)
	assert.Contains(t, nodes, "vpc")
	assert.Contains(t, nodes, "database")
	assert.Contains(t, nodes, "app")
}

func TestClone_PathGeneration(t *testing.T) {
	// Test that generateRepoPath creates deterministic paths
	url1 := "https://github.com/example/repo.git"
	revision1 := "main"

	path1 := generateRepoPath(url1, revision1)
	path2 := generateRepoPath(url1, revision1)

	// Same inputs should generate same path
	assert.Equal(t, path1, path2)

	// Different inputs should generate different paths
	path3 := generateRepoPath(url1, "develop")
	assert.NotEqual(t, path1, path3)

	path4 := generateRepoPath("https://github.com/other/repo.git", revision1)
	assert.NotEqual(t, path1, path4)

	// Paths should be under /tmp/cloudlab-repos/
	assert.True(t, strings.HasPrefix(path1, "/tmp/cloudlab-repos/"))
}

func TestClone_CheckRepoStatus(t *testing.T) {
	// Test checkRepoStatus with non-existent directory
	nonExistentPath := "/tmp/non-existent-repo-12345"
	exists, hash := checkRepoStatus(context.Background(), nonExistentPath, "main")
	assert.False(t, exists)
	assert.Empty(t, hash)
}
