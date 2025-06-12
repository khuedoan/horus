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
	s.True(prunedGraph.Nodes["database"]) // changed module should be included
	s.True(prunedGraph.Nodes["app"])      // dependent should be included
	s.False(prunedGraph.Nodes["vpc"])     // non-dependent should be pruned
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
	s.Empty(prunedGraph.Nodes) // no changes means empty graph
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
	s.True(prunedGraph.Nodes["database"])   // changed module
	s.True(prunedGraph.Nodes["app"])        // direct dependent
	s.True(prunedGraph.Nodes["monitoring"]) // transitive dependent
	s.False(prunedGraph.Nodes["vpc"])       // not a dependent
	s.False(prunedGraph.Nodes["cache"])     // not a dependent
}

// Test using TestActivityEnvironment for activities that need proper context
func (s *ActivityTestSuite) TestTerragruntPrune_WithActivityEnvironment() {
	// Test data
	graph := &Graph{
		Nodes: map[string]bool{
			"vpc":        true,
			"database":   true,
			"app":        true,
			"monitoring": true,
		},
		Edges: map[string][]string{
			"app":        {"database", "vpc"},
			"database":   {"vpc"},
			"monitoring": {"app"},
		},
	}

	changedModules := []string{"database"}

	s.env.RegisterActivity(PruneGraph)

	val, err := s.env.ExecuteActivity(PruneGraph, graph, changedModules)
	s.NoError(err)

	var result *Graph
	s.NoError(val.Get(&result))

	// Only database (changed) and its dependents (app, monitoring) should be included
	// vpc is not included because nothing depends on it
	expectedNodes := []string{"database", "app", "monitoring"}
	actualNodes := result.GetNodes()
	s.ElementsMatch(expectedNodes, actualNodes)

	s.Contains(result.Nodes, "database")
	s.Contains(result.Nodes, "app")
	s.Contains(result.Nodes, "monitoring")
	s.NotContains(result.Nodes, "vpc")
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
	assert.Empty(t, graph.Nodes)
}

func TestNewGraphFromDot_InvalidFormat(t *testing.T) {
	dotString := `not a valid dot format`

	graph, err := NewGraphFromDot(dotString)

	assert.NoError(t, err) // Should not error, just ignore invalid lines
	assert.Empty(t, graph.Nodes)
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

	// All nodes should be present somewhere in the levels
	allNodes := make(map[string]bool)
	for _, level := range levels {
		for _, node := range level {
			allNodes[node] = true
		}
	}
	assert.True(t, allNodes["a"])
	assert.True(t, allNodes["b"])
	assert.True(t, allNodes["c"])
}

func TestExtractQuoted(t *testing.T) {
	// Test the extractQuoted function indirectly through NewGraphFromDot
	dotString := `digraph {
		"hello" -> "world";
		"test";
	}`

	graph, err := NewGraphFromDot(dotString)

	assert.NoError(t, err)
	assert.True(t, graph.Nodes["hello"])
	assert.True(t, graph.Nodes["world"])
	assert.True(t, graph.Nodes["test"])
}

func TestGraph_AddEdge_CreatesNodes(t *testing.T) {
	graph := NewGraph()

	graph.AddEdge("a", "b")

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
	// Test hasCorrectRevision with non-existent directory
	nonExistentPath := "/tmp/non-existent-repo-12345"
	hasCorrect := hasCorrectRevision(context.Background(), nonExistentPath, "main")
	assert.False(t, hasCorrect)
}
