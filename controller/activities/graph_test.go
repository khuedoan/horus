package activities

import (
	"context"
	"reflect"
	"sort"
	"testing"
)

func TestPruneGraphSimple(t *testing.T) {
	dot := `
digraph {
	"A" -> "B";
	"B" -> "C";
	"D" -> "B";
	"E" -> "F";
	"C";
	"F";
}
`
	graph, err := NewGraphFromDot(dot)
	if err != nil {
		t.Fatalf("Failed to create graph from dot: %v", err)
	}

	testCases := []struct {
		name          string
		changed       []string
		expectedNodes []string
		expectedEdges []Edge
	}{
		{
			name:          "Prune to C and its dependencies",
			changed:       []string{"C"},
			expectedNodes: []string{"A", "B", "C", "D"},
			expectedEdges: []Edge{
				{Src: "A", Dest: "B"},
				{Src: "B", Dest: "C"},
				{Src: "D", Dest: "B"},
			},
		},
		{
			name:          "Prune to F and its dependencies",
			changed:       []string{"F"},
			expectedNodes: []string{"E", "F"},
			expectedEdges: []Edge{
				{Src: "E", Dest: "F"},
			},
		},
		{
			name:          "Prune to B and its dependencies",
			changed:       []string{"B"},
			expectedNodes: []string{"A", "B", "D"},
			expectedEdges: []Edge{
				{Src: "A", Dest: "B"},
				{Src: "D", Dest: "B"},
			},
		},
		{
			name:          "No nodes changed",
			changed:       []string{},
			expectedNodes: []string{},
			expectedEdges: []Edge{},
		},
		{
			name:          "Changed node not in graph",
			changed:       []string{"Z"},
			expectedNodes: []string{},
			expectedEdges: []Edge{},
		},
		{
			name:          "Multiple changed nodes",
			changed:       []string{"C", "F"},
			expectedNodes: []string{"A", "B", "C", "D", "E", "F"},
			expectedEdges: []Edge{
				{Src: "A", Dest: "B"},
				{Src: "B", Dest: "C"},
				{Src: "D", Dest: "B"},
				{Src: "E", Dest: "F"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prunedGraph, err := PruneGraph(context.Background(), graph, tc.changed)
			if err != nil {
				t.Fatalf("PruneGraph failed: %v", err)
			}

			prunedNodes := make([]string, 0, len(prunedGraph.Nodes))
			for _, n := range prunedGraph.Nodes {
				prunedNodes = append(prunedNodes, n.Name)
			}
			sort.Strings(prunedNodes)
			sort.Strings(tc.expectedNodes)

			if !reflect.DeepEqual(prunedNodes, tc.expectedNodes) {
				t.Errorf("Expected nodes %v, but got %v", tc.expectedNodes, prunedNodes)
			}

			prunedEdges := make([]Edge, len(prunedGraph.Edges))
			for i, e := range prunedGraph.Edges {
				prunedEdges[i] = *e
			}

			// Sort edges for consistent comparison
			sort.Slice(prunedEdges, func(i, j int) bool {
				if prunedEdges[i].Src != prunedEdges[j].Src {
					return prunedEdges[i].Src < prunedEdges[j].Src
				}
				return prunedEdges[i].Dest < prunedEdges[j].Dest
			})
			sort.Slice(tc.expectedEdges, func(i, j int) bool {
				if tc.expectedEdges[i].Src != tc.expectedEdges[j].Src {
					return tc.expectedEdges[i].Src < tc.expectedEdges[j].Src
				}
				return tc.expectedEdges[i].Dest < tc.expectedEdges[j].Dest
			})

			if !reflect.DeepEqual(prunedEdges, tc.expectedEdges) {
				t.Errorf("Expected edges %v, but got %v", tc.expectedEdges, prunedEdges)
			}
		})
	}
}

func TestPruneGraphRealWorld(t *testing.T) {
	realWorldDot := `digraph {
	"aks-windows-node-exporter" ;
	"azuresql" ;
	"azuresql" -> "core";
	"azuresqlusers" ;
	"azuresqlusers" -> "azuresql";
	"bootstrap-va" ;
	"bootstrap-va" -> "cluster-va";
	"bootstrap2-va" ;
	"bootstrap2-va" -> "cluster2-va";
	"cluster-va" ;
	"cluster-va" -> "core";
	"cluster2-va" ;
	"core" ;
	"db/auror-integration" ;
	"db/auror-integration" -> "core";
	"db/auror-integration" -> "azuresql";
	"db/doc-chat" ;
	"db/doc-chat" -> "core";
	"db/doc-chat" -> "azuresql";
	"dems-cluster-identity" ;
	"dems-cluster-identity" -> "cluster-va";
	"dems-search-grpc/cosmosdb-cassandra" ;
	"dems-search-grpc/cosmosdb-cassandra" -> "core";
	"doc-chat/openai" ;
	"doc-chat/openai" -> "core";
	"doc-chat/openai-fallback" ;
	"doc-chat/openai-fallback" -> "core";
	"doc-chat/search-service-va" ;
	"doc-chat/search-service-va" -> "core";
	"ecom/arkham-hsm-als-endpoint" ;
	"ecom/arkham-hsm-legacy-endpoint" ;
	"ecom/redis" ;
	"ecom/redis" -> "core";
	"ecom/redis-case" ;
	"ecom/redis-case" -> "core";
	"ecom/redis-legacy-endpoint" ;
	"ecom/redis-legacy-endpoint" -> "core";
	"ecom/redis-legacy-endpoint" -> "ecom/redis";
	"ecom/redis-sharon" ;
	"ecom/redis-sharon" -> "core";
	"ecom/redis-webhooks-premium" ;
	"ecom/redis-webhooks-premium" -> "core";
	"endpoints/azuresql-legacy-endpoint-tx" ;
	"endpoints/azuresql-legacy-endpoint-tx" -> "core";
	"endpoints/azuresql-legacy-endpoint-tx" -> "azuresql";
	"endpoints/azuresql-legacy-endpoint-va" ;
	"endpoints/azuresql-legacy-endpoint-va" -> "core";
	"endpoints/azuresql-legacy-endpoint-va" -> "azuresql";
	"endpoints/storage-accounts" ;
	"endpoints/storage-accounts" -> "storage-accounts/ingestion";
	"enterprise/app-identity/auror" ;
	"enterprise/app-identity/auror" -> "core";
	"enterprise/keyvault/auror" ;
	"enterprise/keyvault/auror" -> "core";
	"enterprise/keyvault/auror" -> "enterprise/app-identity/auror";
	"enterprise/redis/auror" ;
	"enterprise/redis/auror" -> "core";
	"espio/az-openai" ;
	"espio/az-openai" -> "core";
	"espio/espio-redis" ;
	"espio/espio-redis" -> "core";
	"espio/openai" ;
	"espio/openai-b" ;
	"espio/openai-b" -> "core";
	"eventgrid-subscription" ;
	"eventgrid-subscription" -> "core";
	"evp/hyperscale" ;
	"evp/hyperscale" -> "core";
	"evp/hyperscaleusers" ;
	"evp/hyperscaleusers" -> "evp/hyperscale";
	"performance/lakehouse" ;
	"performance/lakehouse" -> "core";
	"performance/redis-jarvis" ;
	"performance/redis-jarvis" -> "core";
	"performance/redis-pipeline" ;
	"performance/redis-pipeline" -> "core";
	"performance/redis-starhopper" ;
	"performance/redis-starhopper" -> "core";
	"pes/keyvault" ;
	"pes/keyvault" -> "core";
	"pes/keyvault" -> "dems-cluster-identity";
	"ratelimit/redis" ;
	"ratelimit/redis" -> "core";
	"sage/datafactory" ;
	"sage/datafactory" -> "core";
	"sage/datafactory/alerts" ;
	"sage/datafactory/alerts" -> "sage/datafactory";
	"sage/datafactory/alerts" -> "core";
	"sage/datafactory/evidence-domain-migration-internal-pipeline" ;
	"sage/datafactory/evidence-domain-migration-internal-pipeline" -> "sage/datafactory";
	"sage/datafactory/evidence-domain-migration-main-pipeline" ;
	"sage/datafactory/evidence-domain-migration-main-pipeline" -> "sage/datafactory";
	"sage/endpoints/hyperscale-legacy-endpoint-tx" ;
	"sage/endpoints/hyperscale-legacy-endpoint-tx" -> "core";
	"sage/endpoints/hyperscale-legacy-endpoint-tx" -> "sage/hyperscale";
	"sage/endpoints/hyperscale-legacy-endpoint-va" ;
	"sage/endpoints/hyperscale-legacy-endpoint-va" -> "core";
	"sage/endpoints/hyperscale-legacy-endpoint-va" -> "sage/hyperscale";
	"sage/hyperscale" ;
	"sage/hyperscale" -> "core";
	"sage/hyperscale/named-replica" ;
	"sage/hyperscale/named-replica" -> "sage/hyperscale";
	"sage/hyperscale/named-replica" -> "core";
	"sage/hyperscaleusers" ;
	"sage/hyperscaleusers" -> "sage/hyperscale";
	"sage/redis" ;
	"sage/redis" -> "core";
	"servicebus-premium" ;
	"servicebus-premium" -> "core";
	"sonic/rev-storage" ;
	"sonic/rev-storage" -> "core";
	"sonic/sonic" ;
	"sonic/sonic" -> "core";
	"sonic/sonic-redis" ;
	"sonic/sonic-redis" -> "core";
	"sonic/translation" ;
	"sonic/translation" -> "core";
	"storage-accounts/case-share" ;
	"storage-accounts/ingestion" ;
	"storage-accounts/rtiworker" ;
	"storage-accounts/sage" ;
	"system-status/cosmosdb-cassandra" ;
	"system-status/cosmosdb-cassandra" -> "core";
	"user-settings/cosmosdb-cassandra" ;
	"user-settings/cosmosdb-cassandra" -> "core";
	"visionsearchpoc/vision" ;
	"visionsearchpoc/vision" -> "core";
	"visualization/cosmosdb-cassandra" ;
	"visualization/cosmosdb-cassandra" -> "core";
	"visualization/redis-cluster" ;
	"visualization/redis-cluster" -> "core";
	"visualization/redis-cluster-rtm" ;
	"visualization/redis-cluster-rtm" -> "core";
	"vm-apps/lsln-500" ;
	"vm-apps/lsln-500" -> "core";
	"vm-apps/lsln-500" -> "vm-apps/solr8-j11-lb";
	"vm-apps/solr8-j11-lb" ;
	"vm-apps/solr8-j11-lb" -> "core";
	"webhooks/cosmosdb-cassandra-dispatch" ;
	"webhooks/cosmosdb-cassandra-dispatch" -> "core";
	"xshare/azuresql" ;
	"xshare/azuresql" -> "core";
	"xshare/azuresqlusers" ;
	"xshare/azuresqlusers" -> "xshare/azuresql";
}`

	graph, err := NewGraphFromDot(realWorldDot)
	if err != nil {
		t.Fatalf("Failed to create graph from real-world DOT: %v", err)
	}

	// Test case: cluster-va changed
	prunedGraph, err := PruneGraph(context.Background(), graph, []string{"cluster-va"})
	if err != nil {
		t.Fatalf("PruneGraph failed: %v", err)
	}

	// digraph {
	//   "bootstrap-va" -> "cluster-va";
	//   "dems-cluster-identity" -> "cluster-va";
	//   "pes/keyvault" -> "dems-cluster-identity";
	// }
	expectedNodes := []string{"bootstrap-va", "cluster-va", "dems-cluster-identity", "pes/keyvault"}
	expectedEdges := []Edge{
		{Src: "bootstrap-va", Dest: "cluster-va"},
		{Src: "dems-cluster-identity", Dest: "cluster-va"},
		{Src: "pes/keyvault", Dest: "dems-cluster-identity"},
	}

	prunedNodes := make([]string, 0, len(prunedGraph.Nodes))
	for _, n := range prunedGraph.Nodes {
		prunedNodes = append(prunedNodes, n.Name)
	}
	sort.Strings(prunedNodes)
	sort.Strings(expectedNodes)

	if !reflect.DeepEqual(prunedNodes, expectedNodes) {
		t.Errorf("Expected nodes %v, but got %v", expectedNodes, prunedNodes)
	}

	prunedEdges := make([]Edge, len(prunedGraph.Edges))
	for i, e := range prunedGraph.Edges {
		prunedEdges[i] = *e
	}

	// Sort edges for consistent comparison
	sort.Slice(prunedEdges, func(i, j int) bool {
		if prunedEdges[i].Src != prunedEdges[j].Src {
			return prunedEdges[i].Src < prunedEdges[j].Src
		}
		return prunedEdges[i].Dest < prunedEdges[j].Dest
	})
	sort.Slice(expectedEdges, func(i, j int) bool {
		if expectedEdges[i].Src != expectedEdges[j].Src {
			return expectedEdges[i].Src < expectedEdges[j].Src
		}
		return expectedEdges[i].Dest < expectedEdges[j].Dest
	})

	if !reflect.DeepEqual(prunedEdges, expectedEdges) {
		t.Errorf("Expected edges %v, but got %v", expectedEdges, prunedEdges)
	}
}

func TestTopologicalSort(t *testing.T) {
	testCases := []struct {
		name           string
		nodes          []string
		edges          []Edge
		expectedLevels [][]string
	}{
		{
			name:  "Simple linear dependency",
			nodes: []string{"A", "B", "C"},
			edges: []Edge{
				{Src: "B", Dest: "A"},
				{Src: "C", Dest: "B"},
			},
			expectedLevels: [][]string{
				{"A"},
				{"B"},
				{"C"},
			},
		},
		{
			name:  "Parallel dependencies",
			nodes: []string{"A", "B", "C", "D"},
			edges: []Edge{
				{Src: "C", Dest: "A"},
				{Src: "C", Dest: "B"},
				{Src: "D", Dest: "C"},
			},
			expectedLevels: [][]string{
				{"A", "B"},
				{"C"},
				{"D"},
			},
		},
		{
			name:  "Complex dependency graph",
			nodes: []string{"A", "B", "C", "D", "E", "F"},
			edges: []Edge{
				{Src: "C", Dest: "A"},
				{Src: "C", Dest: "B"},
				{Src: "D", Dest: "C"},
				{Src: "E", Dest: "C"},
				{Src: "F", Dest: "D"},
				{Src: "F", Dest: "E"},
			},
			expectedLevels: [][]string{
				{"A", "B"},
				{"C"},
				{"D", "E"},
				{"F"},
			},
		},
		{
			name:           "No dependencies",
			nodes:          []string{"A", "B", "C"},
			edges:          []Edge{},
			expectedLevels: [][]string{{"A", "B", "C"}},
		},
		{
			name:           "Single node",
			nodes:          []string{"A"},
			edges:          []Edge{},
			expectedLevels: [][]string{{"A"}},
		},
		{
			name:  "Real world example: bootstrap depends on cluster",
			nodes: []string{"bootstrap", "cluster"},
			edges: []Edge{
				{Src: "bootstrap", Dest: "cluster"},
			},
			expectedLevels: [][]string{
				{"cluster"},
				{"bootstrap"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create graph
			graph := &Graph{
				Nodes: make([]*Node, len(tc.nodes)),
				Edges: make([]*Edge, len(tc.edges)),
			}

			for i, nodeName := range tc.nodes {
				graph.Nodes[i] = &Node{Name: nodeName}
			}

			for i, edge := range tc.edges {
				graph.Edges[i] = &Edge{Src: edge.Src, Dest: edge.Dest}
			}

			// Get topological sort
			levels := graph.TopologicalSort()

			// Verify number of levels
			if len(levels) != len(tc.expectedLevels) {
				t.Errorf("Expected %d levels, but got %d", len(tc.expectedLevels), len(levels))
				return
			}

			// Verify each level
			for levelIndex, expectedLevel := range tc.expectedLevels {
				actualLevel := levels[levelIndex]

				// Sort both slices for comparison
				sort.Strings(expectedLevel)
				sort.Strings(actualLevel)

				if !reflect.DeepEqual(actualLevel, expectedLevel) {
					t.Errorf("Level %d: expected %v, but got %v", levelIndex, expectedLevel, actualLevel)
				}
			}
		})
	}
}
