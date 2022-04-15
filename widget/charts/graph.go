package charts

import (
	"sync"
)

// GraphRange set the range of the graph.
type GraphRange struct {
	// Set the Y range of the graph.
	YMin, YMax float64
}

// global options for the chart
type sizeOpts struct {
	GraphRange *GraphRange
}

func newSizeOpts() *sizeOpts {
	return &sizeOpts{}
}

// SetGraphRange set the global range of the graph.
func (g *sizeOpts) SetGraphRange(r *GraphRange) {
	g.GraphRange = r
}

// BaseChart struct for any Graph object.
type BaseChart struct {

	// to avoid race condition
	locker *sync.Mutex
}

// Lock is mainly used to prevent race condition when you need to access data.
func (b *BaseChart) Lock() {
	b.getLocker().Lock()
}

// Unlock releases the locker
func (b *BaseChart) Unlock() {
	b.getLocker().Unlock()
}

func (b *BaseChart) getLocker() *sync.Mutex {
	if b.locker == nil {
		b.locker = &sync.Mutex{}
	}
	return b.locker
}
