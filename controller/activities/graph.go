package activities

import (
	"fmt"
	"strings"

	"github.com/awalterschulze/gographviz"
)

func pruneGraph(dot string, changed []string) (string, error) {
	ast, err := gographviz.ParseString(dot)
	if err != nil {
		return "", fmt.Errorf("Parse DOT failed: %w", err)
	}
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(ast, graph); err != nil {
		return "", fmt.Errorf("Graph analysis failed: %w", err)
	}

	// Build reverse dependency map: target -> dependents
	dependents := map[string][]string{}
	for _, edge := range graph.Edges.Edges {
		dependents[edge.Dst] = append(dependents[edge.Dst], edge.Src)
	}

	// Collect all nodes to keep (changed + all that depend on them)
	keep := map[string]bool{}
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
	for _, node := range changed {
		visit(node)
	}

	// Reconstruct pruned DOT graph
	var b strings.Builder
	b.WriteString("digraph {\n")
	for _, edge := range graph.Edges.Edges {
		if keep[edge.Src] && keep[edge.Dst] {
			b.WriteString(fmt.Sprintf("  %s -> %s;\n", edge.Src, edge.Dst))
		}
	}
	for node := range keep {
		b.WriteString(fmt.Sprintf("  %s ;\n", node))
	}
	b.WriteString("}")
	return b.String(), nil
}
