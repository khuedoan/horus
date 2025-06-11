package activities

import (
	"context"
	"fmt"
	"strings"
)

// Graph represents a simple directed graph with efficient operations.
type Graph struct {
	Nodes map[string]bool     `json:"nodes"`
	Edges map[string][]string `json:"edges"` // source -> []destinations
}

// NewGraph creates a new empty graph.
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]bool),
		Edges: make(map[string][]string),
	}
}

// AddNode adds a node to the graph.
func (g *Graph) AddNode(name string) {
	g.Nodes[name] = true
}

// AddEdge adds a directed edge from src to dest.
func (g *Graph) AddEdge(src, dest string) {
	g.AddNode(src)
	g.AddNode(dest)
	g.Edges[src] = append(g.Edges[src], dest)
}

// GetNodes returns all node names.
func (g *Graph) GetNodes() []string {
	nodes := make([]string, 0, len(g.Nodes))
	for name := range g.Nodes {
		nodes = append(nodes, name)
	}
	return nodes
}

// NodeCount returns the number of nodes in the graph.
func (g *Graph) NodeCount() int {
	return len(g.Nodes)
}

// EdgeCount returns the number of edges in the graph.
func (g *Graph) EdgeCount() int {
	count := 0
	for _, dests := range g.Edges {
		count += len(dests)
	}
	return count
}

// PruneGraph takes a graph and a list of changed nodes, and returns a new graph
// containing only the changed nodes and their dependents.
func PruneGraph(ctx context.Context, graph *Graph, changed []string) (*Graph, error) {
	// Build reverse dependency map: target -> dependents
	dependents := make(map[string][]string)
	for src, dests := range graph.Edges {
		for _, dest := range dests {
			dependents[dest] = append(dependents[dest], src)
		}
	}

	// Collect all nodes to keep (changed + all that depend on them)
	keep := make(map[string]bool)
	var visit func(string)
	visit = func(node string) {
		if keep[node] {
			return
		}
		keep[node] = true
		for _, dep := range dependents[node] {
			visit(dep)
		}
	}

	// Only visit nodes that actually exist in the graph
	for _, nodeName := range changed {
		if graph.Nodes[nodeName] {
			visit(nodeName)
		}
	}

	// Create pruned graph
	prunedGraph := NewGraph()
	for node := range keep {
		prunedGraph.AddNode(node)
	}
	for src, dests := range graph.Edges {
		if keep[src] {
			for _, dest := range dests {
				if keep[dest] {
					prunedGraph.AddEdge(src, dest)
				}
			}
		}
	}

	return prunedGraph, nil
}

// NewGraphFromDot creates a Graph from a DOT string using a simple parser.
func NewGraphFromDot(dot string) (*Graph, error) {
	graph := NewGraph()

	lines := strings.Split(dot, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") || line == "digraph {" || line == "}" {
			continue
		}

		// Remove trailing semicolon if present
		line = strings.TrimSuffix(line, ";")
		line = strings.TrimSpace(line)

		// Parse edges: "A" -> "B"
		if strings.Contains(line, "->") {
			parts := strings.Split(line, "->")
			if len(parts) == 2 {
				src := extractQuotedString(strings.TrimSpace(parts[0]))
				dest := extractQuotedString(strings.TrimSpace(parts[1]))
				if src != "" && dest != "" {
					graph.AddEdge(src, dest)
				}
			}
		} else {
			// Parse standalone nodes: "C"
			nodeName := extractQuotedString(line)
			if nodeName != "" {
				graph.AddNode(nodeName)
			}
		}
	}

	return graph, nil
}

// extractQuotedString extracts the content between quotes from a string like "hello"
func extractQuotedString(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return ""
}

// ToDot converts a Graph to a DOT string.
func (g *Graph) ToDot() string {
	var b strings.Builder
	b.WriteString("digraph {\n")

	for src, dests := range g.Edges {
		for _, dest := range dests {
			b.WriteString(fmt.Sprintf("  %q -> %q;\n", src, dest))
		}
	}

	// Write standalone nodes (those not in any edge)
	edgeNodes := make(map[string]bool)
	for src, dests := range g.Edges {
		edgeNodes[src] = true
		for _, dest := range dests {
			edgeNodes[dest] = true
		}
	}

	for node := range g.Nodes {
		if !edgeNodes[node] {
			b.WriteString(fmt.Sprintf("  %q;\n", node))
		}
	}

	b.WriteString("}")
	return b.String()
}

// TopologicalSort returns modules grouped by dependency levels for parallel execution.
// Edge from A to B means A depends on B, so B must run before A.
func (g *Graph) TopologicalSort() [][]string {
	// Build adjacency list and in-degree count
	adjList := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialize all nodes with in-degree 0
	for node := range g.Nodes {
		inDegree[node] = 0
		adjList[node] = []string{}
	}

	// Build the graph and calculate in-degrees
	// Edge from Src to Dest means Src depends on Dest
	// So Dest should run before Src
	for src, dests := range g.Edges {
		for _, dest := range dests {
			adjList[dest] = append(adjList[dest], src)
			inDegree[src]++
		}
	}

	var levels [][]string
	remaining := make(map[string]bool)
	for node := range g.Nodes {
		remaining[node] = true
	}

	// Process nodes level by level
	for len(remaining) > 0 {
		var currentLevel []string

		// Find all nodes with in-degree 0 (no dependencies)
		for nodeName := range remaining {
			if inDegree[nodeName] == 0 {
				currentLevel = append(currentLevel, nodeName)
			}
		}

		// If no nodes found with in-degree 0, there's a cycle
		if len(currentLevel) == 0 {
			// Return remaining nodes as the final level to handle cycles gracefully
			var cycleNodes []string
			for nodeName := range remaining {
				cycleNodes = append(cycleNodes, nodeName)
			}
			if len(cycleNodes) > 0 {
				levels = append(levels, cycleNodes)
			}
			break
		}

		// Add current level
		levels = append(levels, currentLevel)

		// Remove processed nodes and update in-degrees
		for _, nodeName := range currentLevel {
			delete(remaining, nodeName)
			for _, dependent := range adjList[nodeName] {
				if remaining[dependent] {
					inDegree[dependent]--
				}
			}
		}
	}

	return levels
}
