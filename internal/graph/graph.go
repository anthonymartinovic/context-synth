package graph

import (
	"fmt"

	"github.com/anthonymartinovic/context-synth/internal/protocol"
)

type CapabilityGraph struct {
	capabilities map[string]protocol.Capability
	adjacency    map[string][]string // depends_on edges: key depends on values
	order        []string            // resolved topological order
}

func Build(capabilities []protocol.Capability) (*CapabilityGraph, error) {
	g := &CapabilityGraph{
		capabilities: make(map[string]protocol.Capability),
		adjacency:    make(map[string][]string),
	}

	for _, cap := range capabilities {
		if _, exists := g.capabilities[cap.ID]; exists {
			return nil, fmt.Errorf("duplicate capability ID %q", cap.ID)
		}
		g.capabilities[cap.ID] = cap
		g.adjacency[cap.ID] = cap.DependsOn
	}

	for id, deps := range g.adjacency {
		for _, dep := range deps {
			if _, exists := g.capabilities[dep]; !exists {
				return nil, fmt.Errorf("capability %q depends on unknown capability %q", id, dep)
			}
		}
	}

	order, err := topologicalSort(g.capabilities, g.adjacency)
	if err != nil {
		return nil, err
	}
	g.order = order

	return g, nil
}

func (g *CapabilityGraph) Order() []string {
	return g.order
}

func (g *CapabilityGraph) Get(id string) (protocol.Capability, bool) {
	cap, ok := g.capabilities[id]
	return cap, ok
}

func topologicalSort(caps map[string]protocol.Capability, adj map[string][]string) ([]string, error) {
	inDegree := make(map[string]int)
	for id := range caps {
		inDegree[id] = 0
	}

	// Build reverse edges: for each dependency, the dependent has an in-degree
	dependents := make(map[string][]string)
	for id, deps := range adj {
		for _, dep := range deps {
			dependents[dep] = append(dependents[dep], id)
			inDegree[id]++
		}
	}

	// Start with nodes that have no dependencies
	var queue []string
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}

	var order []string
	for len(queue) > 0 {
		// Take first element (stable order not guaranteed without sorting, but deterministic enough)
		node := queue[0]
		queue = queue[1:]
		order = append(order, node)

		for _, dependent := range dependents[node] {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}
	}

	if len(order) != len(caps) {
		return nil, fmt.Errorf("capability graph contains a cycle")
	}

	return order, nil
}
