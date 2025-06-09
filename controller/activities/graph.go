package activities

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/awalterschulze/gographviz"
)

// Node represents a node in the graph.
type Node struct {
	Name string
	// Attributes can be added here if needed
}

// Edge represents a directed edge in the graph.
type Edge struct {
	Src  string
	Dest string
	// Attributes can be added here if needed
}

// Graph represents a serializable graph.
type Graph struct {
	Nodes []*Node
	Edges []*Edge
}

// PruneGraph takes a graph and a list of changed nodes, and returns a new graph
// containing only the changed nodes and their dependents.
func PruneGraph(ctx context.Context, graph *Graph, changed []string) (*Graph, error) {
	// Build reverse dependency map: target -> dependents
	dependents := make(map[string][]string)
	for _, edge := range graph.Edges {
		dependents[edge.Dest] = append(dependents[edge.Dest], edge.Src)
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
	for _, nodeName := range changed {
		// Make sure we have the node in the graph before visiting
		var found bool
		for _, n := range graph.Nodes {
			if n.Name == nodeName {
				found = true
				break
			}
		}
		if found {
			visit(nodeName)
		}
	}

	// Reconstruct pruned graph
	prunedGraph := &Graph{Nodes: []*Node{}, Edges: []*Edge{}}
	for _, node := range graph.Nodes {
		if keep[node.Name] {
			prunedGraph.Nodes = append(prunedGraph.Nodes, node)
		}
	}
	for _, edge := range graph.Edges {
		if keep[edge.Src] && keep[edge.Dest] {
			prunedGraph.Edges = append(prunedGraph.Edges, edge)
		}
	}

	return prunedGraph, nil
}

// NewGraphFromDot creates a Graph from a DOT string.
func NewGraphFromDot(dot string) (*Graph, error) {
	ast, err := gographviz.ParseString(dot)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DOT string: %w", err)
	}

	g := gographviz.NewGraph()
	if err := gographviz.Analyse(ast, g); err != nil {
		return nil, fmt.Errorf("failed to analyse graph: %w", err)
	}

	graph := &Graph{Nodes: []*Node{}, Edges: []*Edge{}}
	nodeSet := make(map[string]bool)

	addNode := func(name string) {
		if !nodeSet[name] {
			graph.Nodes = append(graph.Nodes, &Node{Name: name})
			nodeSet[name] = true
		}
	}

	unquote := func(s string) string {
		res, err := strconv.Unquote(s)
		if err != nil {
			return s
		}
		return res
	}

	for _, edge := range g.Edges.Edges {
		src := unquote(edge.Src)
		dest := unquote(edge.Dst)
		addNode(src)
		addNode(dest)
		graph.Edges = append(graph.Edges, &Edge{Src: src, Dest: dest})
	}

	for _, node := range g.Nodes.Nodes {
		name := unquote(node.Name)
		addNode(name)
	}
	return graph, nil
}

// ToDot converts a Graph to a DOT string.
func (g *Graph) ToDot() string {
	var b strings.Builder
	b.WriteString("digraph {\n")
	for _, edge := range g.Edges {
		b.WriteString(fmt.Sprintf("  %q -> %q;\n", edge.Src, edge.Dest))
	}
	for _, node := range g.Nodes {
		b.WriteString(fmt.Sprintf("  %q;\n", node.Name))
	}
	b.WriteString("}")
	return b.String()
}
