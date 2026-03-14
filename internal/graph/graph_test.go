package graph

import (
	"testing"

	"github.com/anthonymartinovic/context-synth/internal/protocol"
)

func TestBuild_LinearDependency(t *testing.T) {
	caps := []protocol.Capability{
		{ID: "a", DependsOn: []string{}},
		{ID: "b", DependsOn: []string{"a"}},
		{ID: "c", DependsOn: []string{"b"}},
	}

	g, err := Build(caps)
	if err != nil {
		t.Fatal(err)
	}

	order := g.Order()
	if len(order) != 3 {
		t.Fatalf("order length = %d, want 3", len(order))
	}

	indexOf := func(id string) int {
		for i, o := range order {
			if o == id {
				return i
			}
		}
		return -1
	}

	if indexOf("a") > indexOf("b") {
		t.Error("a should come before b")
	}
	if indexOf("b") > indexOf("c") {
		t.Error("b should come before c")
	}
}

func TestBuild_NoDependencies(t *testing.T) {
	caps := []protocol.Capability{
		{ID: "x"},
		{ID: "y"},
	}

	g, err := Build(caps)
	if err != nil {
		t.Fatal(err)
	}

	if len(g.Order()) != 2 {
		t.Errorf("order length = %d, want 2", len(g.Order()))
	}
}

func TestBuild_CycleDetection(t *testing.T) {
	caps := []protocol.Capability{
		{ID: "a", DependsOn: []string{"b"}},
		{ID: "b", DependsOn: []string{"a"}},
	}

	_, err := Build(caps)
	if err == nil {
		t.Fatal("expected error for cyclic dependency")
	}
}

func TestBuild_UnknownDependency(t *testing.T) {
	caps := []protocol.Capability{
		{ID: "a", DependsOn: []string{"missing"}},
	}

	_, err := Build(caps)
	if err == nil {
		t.Fatal("expected error for unknown dependency")
	}
}

func TestBuild_DuplicateID(t *testing.T) {
	caps := []protocol.Capability{
		{ID: "a"},
		{ID: "a"},
	}

	_, err := Build(caps)
	if err == nil {
		t.Fatal("expected error for duplicate ID")
	}
}

func TestBuild_Get(t *testing.T) {
	caps := []protocol.Capability{
		{ID: "test", Description: "a test capability"},
	}

	g, err := Build(caps)
	if err != nil {
		t.Fatal(err)
	}

	cap, ok := g.Get("test")
	if !ok {
		t.Fatal("capability not found")
	}
	if cap.Description != "a test capability" {
		t.Errorf("description = %q", cap.Description)
	}

	_, ok = g.Get("nonexistent")
	if ok {
		t.Error("should not find nonexistent capability")
	}
}
