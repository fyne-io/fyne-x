package charts

import "testing"

func TestGraph_SetRange(t *testing.T) {
	g := &graph{}

	if err := g.SetGraphRange(&GraphRange{0, 10}); err != nil {
		t.Errorf("SetRange failed: %s", err)
	}

	if g.graphRange.YMin != 0 {
		t.Errorf("SetRange failed: Min != 0")
	}

	if g.graphRange.YMax != 10 {
		t.Errorf("SetRange failed: Max != 10")
	}
}

func TestGraph_nilRange(t *testing.T) {
	g := &graph{}

	if err := g.SetGraphRange(nil); err != nil {
		t.Errorf("SetRange failed: %s", err)
	}
}

func TestGraph_badRage(t *testing.T) {
	g := &graph{}

	if err := g.SetGraphRange(&GraphRange{10, 0}); err == nil {
		t.Errorf("SetRange failed: %s", err)
	}
}
