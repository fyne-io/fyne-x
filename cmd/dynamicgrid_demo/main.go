package main

import (
	"errors"
	"fmt"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Dynamic Grid")

	const numTestItems = 128
	names := make([]string, 0, numTestItems)
	for i := 0; i < numTestItems; i++ {
		name := fmt.Sprintf("Button %d", i)
		names = append(names, name)
	}
	w.SetContent(exampleGridler(names))

	w.ShowAndRun()
}

func exampleGridler(names []string) *fyne.Container {
	makeAnItem := func(index int) fyne.CanvasObject {
		// just to spice things up a bit and show more
		// how this works
		if index%10 == 0 {
			return canvas.NewImageFromResource(theme.FyneLogo())
		}
		if index%23 == 0 {
			pb := widget.NewProgressBarInfinite()
			pb.Start()
			return pb
		}
		if index == 42 {
			// fixme this incidentally shows a bug when you slowly change
			// the window size. this widget is not made/meant for different
			// sized widgets, however it would be good to handle them properly
			return container.NewCenter(widget.NewLabel("the meaning of life"))
		}

		name := names[index]
		return widget.NewButton(name, func() { fmt.Println(name) })
	}

	// depending on the widget-heavyness,
	// you likely want a cache
	cache := make(map[int]fyne.CanvasObject)

	dataRequester := func(numItems int) []fyne.CanvasObject {
		objs := make([]fyne.CanvasObject, 0, numItems)

		for _, index := range sampler(len(names), numItems) {
			cached, ok := cache[index]
			if !ok {
				cached = makeAnItem(index)
				cache[index] = cached
			}

			objs = append(objs, cached)
		}

		return objs
	}

	// we always want at least 2 rows and 2 columns
	const imagesPerViewRows = 2
	const imagesPerViewColums = 2
	return NewDynamicGrid(imagesPerViewRows, imagesPerViewColums, dataRequester)
}

// makes a stable same-input same-output sampling
func sampler(maxItems, requiredItems int) []int {
	if requiredItems > maxItems {
		indexes := make([]int, maxItems)
		// in a real scenario we just
		// would not have more to display
		// example file previews
		for index := 0; index < maxItems; index++ {
			indexes[index] = index
		}
		return indexes
	}

	indexes := make([]int, 0, requiredItems)
	step := int(float64(maxItems) / float64(requiredItems))
	found := 0
	for index := 0; found < requiredItems; index += step {
		indexes = append(indexes, index)
		found += 1
	}

	return indexes
}

// ----------------------------------------------------------------
// ----------------------------------------------------------------
// Currently kept here during heavy development, later to be moved
// to the layout folder

// do not directly tinker with Container.Objects or Container.Layout
func NewDynamicGrid(minrows, mincols int, mr MoreRequester) *fyne.Container {
	grid := newDynamicGridLayout(minrows, mincols, mr)
	c := container.New(grid)
	grid.container = c // we need this in the Layout function
	return c
}

// Declare conformity with Layout interface
var _ fyne.Layout = (*dynamicGridLayout)(nil)

type MoreRequester func(numItems int) []fyne.CanvasObject
type dynamicGridLayout struct {
	minrows   int
	mincols   int
	container *fyne.Container
	mr        MoreRequester
}

var ErrorMissingMoreRequester = errors.New("you must provide a MoreRequester function")

// newDynamicGridLayout returns a new dynamic grid layout
func newDynamicGridLayout(minrows, mincols int, mr MoreRequester) *dynamicGridLayout {
	// i dont want to recheck this in hotpath Layout ether
	// fixme go version update replace with max()
	if minrows < 1 {
		minrows = 1
	}
	if mincols < 1 {
		mincols = 1
	}
	if mr == nil {
		// fixme whats the nice way to deal with that?
		// just nop ourselves and return a stackLayout?
		// nil? but then unexplained not working?
		// maybe have NewDynamicGrid just return an error?
		// would mess with cascading
		panic(ErrorMissingMoreRequester)
		return nil
	}

	return &dynamicGridLayout{
		minrows: minrows,
		mincols: mincols,
		mr:      mr,
	}
}

// Layout is called to pack all child objects into a specified size.
// For a DynamicGridLayout this will pack the needed amount of objects
// into a table format with at least the minumum specified columns and rows
// and if less content available to fill as much space as possible.
func (g *dynamicGridLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	padding := theme.Padding()

	childsize := uniformBlocksize(objects).Add(fyne.NewSquareSize(padding / 2))
	numcols := numobjfit(size.Width, childsize.Width)
	numrows := numobjfit(size.Height, childsize.Height)
	maxamount := numrows * numcols

	objects = g.mr(maxamount)
	g.container.Objects = objects //otherwise something is salty and we render nothing
	lenobjects := len(objects)

	if lenobjects < 1 {
		// nothing to do
		return
	}
	if lenobjects < maxamount {
		// stretch over available space but evenly
		numcols, numrows = strechedSizeCalculator(lenobjects, childsize, size)
		maxamount = numcols * numrows
	}

	padWidth := float32(numcols-1) * padding
	padHeight := float32(numrows-1) * padding
	cellWidth := float32(size.Width-padWidth) / float32(numcols)
	cellHeight := float32(size.Height-padHeight) / float32(numrows)

	row, col := 0, 0
	for i, child := range objects {
		// leading edge top left
		x1 := (cellWidth + padding) * float32(col)
		y1 := (cellHeight + padding) * float32(row)
		// trailing edge bottom right
		x2 := (cellWidth+padding)*float32(col+1) - padding
		y2 := (cellHeight+padding)*float32(row+1) - padding

		child.Move(fyne.NewPos(x1, y1))
		child.Resize(fyne.NewSize(x2-x1, y2-y1))

		if (i+1)%numcols == 0 {
			row++
			col = 0
		} else {
			col++
		}
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a DynamicGridLayout this is the size of the largest child object
// multiplied by the minimum number of columns and rows, with
// appropriate padding between children.
func (g *dynamicGridLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	// fixme do we also have to ask for items here?
	// by visual testing the answer is no
	minSize := uniformBlocksize(objects)
	padding := theme.Padding()

	minContentSize := fyne.NewSize(
		minSize.Width*float32(g.mincols),
		minSize.Height*float32(g.minrows),
	)
	return minContentSize.Add(fyne.NewSize(
		padding*fyne.Max(float32(g.mincols-1), 0),
		padding*fyne.Max(float32(g.minrows-1), 0),
	))

}

func numobjfit(fullsize, childsize float32) int {
	// fixme go version update replace with max()
	num := int(math.Floor(float64(fullsize) / float64(childsize)))
	if num < 1 {
		return 1
	}
	return num
}

// fixme find a smart way to do this
// like for real, this cant be it
func strechedSizeCalculator(numobjects int, childsize fyne.Size, size fyne.Size) (int, int) {
	if numobjects < 1 {
		// failsave
		return 1, 1
	}
	var numcols int
	var numrows int
	for {
		childsize = childsize.Add(fyne.NewSquareSize(1))
		numcols = numobjfit(size.Width, childsize.Width)
		numrows = numobjfit(size.Height, childsize.Height)
		if numcols*numrows <= numobjects {
			childsize = childsize.Subtract(fyne.NewSquareSize(1))
			numcols = numobjfit(size.Width, childsize.Width)
			numrows = numobjfit(size.Height, childsize.Height)
			break
		}
	}
	return numcols, numrows
}

func uniformBlocksize(objects []fyne.CanvasObject) fyne.Size {
	childsize := fyne.NewSquareSize(1)
	for _, child := range objects {
		childsize = childsize.Max(child.MinSize())
	}
	return childsize
}
