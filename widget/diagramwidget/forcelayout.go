package diagramwidget

import (
	"math"

	"fyne.io/fyne/v2"

	"fyne.io/x/fyne/widget/diagramwidget/geometry/r2"
)

// adjacent returns true if there is at least one edge between n1 and n2
func (g *DiagramWidget) adjacent(n1, n2 *DiagramNode) bool {
	// TODO: expensive, may be worth caching?
	for _, e := range g.Links {
		if ((e.Origin == n1) && (e.Target == n2)) || ((e.Origin == n2) && (e.Target == n1)) {
			return true
		}
	}

	return false
}

func calculateDistance(n1, n2 *DiagramNode) float64 {
	return r2.MakeLineFromEndpoints(n1.R2Center(), n2.R2Center()).Length()
}

// calculateForce calculates the force between the given pair of nodes.
//
// The force is calculated at n1.
func (g *DiagramWidget) calculateForce(n1, n2 *DiagramNode, targetLength float64) r2.Vec2 {
	// spring constant for linear spring
	k := float64(0.01)
	d := calculateDistance(n1, n2)

	v := n2.R2Center().Add(n1.R2Center().Scale(-1)).Unit().Scale(-1)

	if g.adjacent(n1, n2) {
		// adjacent nodes act like springs, and want to be close to the given
		// length.

		// avoid bouncing
		delta := math.Abs(d - targetLength)
		if delta < 0.05*targetLength {
			return r2.V2(0, 0)
		}

		if d < targetLength {
			return v.Scale(1*d*k + k*math.Pow(d, 1/(d+1)))
		} else {
			return v.Scale(-1*d*k - 0.01*k*math.Pow(d, 2))
		}
	} else {
		if d > 1.2*targetLength {
			return r2.V2(0, 0)
		}
		// non-adjacent nodes repel, at a rate falling of with distance.
		return v.Scale(50 * math.Sqrt(1/(d+0.1)))
		// return r2.V2(0, 0*math.Sqrt(1))
	}
}

// StepForceLayout calculates one step of force directed graph layout, with
// the target distance between adjacent nodes being targetLength.
func (g *DiagramWidget) StepForceLayout(targetLength float64) {
	deltas := make(map[string]r2.Vec2)

	// calculate all the deltas from the current state
	for k, nk := range g.Nodes {
		deltas[k] = r2.V2(0, 0)

		for j, nj := range g.Nodes {
			if j == k {
				continue
			}
			deltas[k] = deltas[k].Add(g.calculateForce(nk, nj, targetLength))
		}
	}

	// flip into current state
	for k, nk := range g.Nodes {
		nk.Displace(fyne.Position{X: float32(deltas[k].X), Y: float32(deltas[k].Y)})
	}

}
