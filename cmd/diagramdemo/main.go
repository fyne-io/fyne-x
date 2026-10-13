package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/x/fyne/widget/diagramwidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var forceticks int = 0

func forceanim(diagramWidget *diagramwidget.DiagramWidget) {

	for {
		if forceticks > 0 {
			fyne.Do(func() {
				diagramwidget.StepForceLayout(diagramWidget, 300)
				diagramWidget.Refresh()
			})
			forceticks--
		}

		time.Sleep(time.Millisecond * (1000 / 30))
	}
}

// ModifiedDragNode illustrates how to extend a node. This particular example
// disables vertical movement and vertical resizing
type ModifiedDragNode struct {
	diagramwidget.BaseDiagramNode
	label *widget.Label
}

// NewModifiedDragNode returns an instance of an ExtendedNode
func NewModifiedDragNode(nodeID string, diagramWidget *diagramwidget.DiagramWidget) diagramwidget.DiagramNode {
	newNode := &ModifiedDragNode{}
	newNode.label = widget.NewLabel("Vertical Changes Not Alloweed")
	diagramwidget.InitializeBaseDiagramNode(newNode, diagramWidget, newNode.label, nodeID)
	newNode.Refresh()
	return newNode
}

// Dragged passes the DragEvent to the diagram for processing after removing any Y value changes
func (en *ModifiedDragNode) Dragged(event *fyne.DragEvent) {
	modifiedDelta := fyne.Delta{
		DX: event.Dragged.DX,
		DY: 0,
	}
	modifiedDragEvent := &fyne.DragEvent{
		PointEvent: event.PointEvent,
		Dragged:    modifiedDelta,
	}
	en.GetDiagram().DiagramNodeDragged(&en.BaseDiagramNode, modifiedDragEvent)
}

// HandleDragged passes the HandleDragged event to the BasseDiagramNode after removing any Y value changes
func (en *ModifiedDragNode) HandleDragged(handle *diagramwidget.Handle, event *fyne.DragEvent) {
	modifiedDelta := fyne.Delta{
		DX: event.Dragged.DX,
		DY: 0,
	}
	modifiedDragEvent := &fyne.DragEvent{
		PointEvent: event.PointEvent,
		Dragged:    modifiedDelta,
	}
	en.BaseDiagramNode.HandleDragged(handle, modifiedDragEvent)
}

// ModifiedPadsNode removes the default rectangle pad on the perimeter and adds two point pads,
// one at each side of the node
type ModifiedPadsNode struct {
	diagramwidget.BaseDiagramNode
	label *widget.Label
}

// NewModifiedPadsNode creates and initializes a ModifiedPadsNode
func NewModifiedPadsNode(nodeID string, diagramWidget *diagramwidget.DiagramWidget) diagramwidget.DiagramNode {
	newNode := &ModifiedPadsNode{}
	newNode.label = widget.NewLabel("Connection pads on left and right")
	diagramwidget.InitializeBaseDiagramNode(newNode, diagramWidget, newNode.label, nodeID)
	// Get rid of the default pad on the perimeter
	newNode.SetConnectionPad(nil, "default")
	newNode.SetConnectionPad(diagramwidget.NewPointPad(newNode.GetProperties().PadStrokeWidth), "left")
	newNode.SetConnectionPad(diagramwidget.NewPointPad(newNode.GetProperties().PadStrokeWidth), "right")
	return newNode
}

func main() {
	app := app.New()
	w := app.NewWindow("Diagram Demo")

	w.SetMaster()

	diagramWidget := diagramwidget.NewDiagramWidget("Diagram1")

	go forceanim(diagramWidget)

	background := canvas.NewCircle(color.RGBA{
		R: 252,
		G: 244,
		B: 3,
		A: 255,
	})

	background.Resize(fyne.NewSize(200, 200))

	diagramWidget.SetBackground(background)

	// Node 0
	node0Label := widget.NewLabel("Node0")
	node0 := diagramwidget.NewDiagramNode(diagramWidget, node0Label, "Node0")
	node0.Move(fyne.NewPos(300, 0))

	// Node 1
	node1Button := widget.NewButton("Node1 Button", func() { fmt.Printf("tapped Node1!\n") })
	node1 := diagramwidget.NewDiagramNode(diagramWidget, node1Button, "Node1")
	node1.Move(fyne.Position{X: 100, Y: 100})

	// Node 2
	node2 := diagramwidget.NewDiagramNode(diagramWidget, nil, "Node2")
	node2Container := container.NewVBox(
		widget.NewLabel("Node2 - with structure"),
		widget.NewButton("Up", func() {
			node2.GetDiagram().DisplaceNode(node2, fyne.Position{X: 0, Y: -10})
			node2.Refresh()
		}),
		widget.NewButton("Down", func() {
			node2.GetDiagram().DisplaceNode(node2, fyne.Position{X: 0, Y: 10})
			node2.Refresh()
		}),
		container.NewHBox(
			widget.NewButton("Left", func() {
				node2.GetDiagram().DisplaceNode(node2, fyne.Position{X: -10, Y: 0})
				node2.Refresh()
			}),
			widget.NewButton("Right", func() {
				node2.GetDiagram().DisplaceNode(node2, fyne.Position{X: 10, Y: 0})
				node2.Refresh()
			}),
		),
	)
	node2.SetInnerObject(node2Container)
	node2.Move(fyne.Position{X: 100, Y: 300})

	// Node 3
	node3 := diagramwidget.NewDiagramNode(diagramWidget, widget.NewButton("Node3: Force layout step", func() {
		diagramwidget.StepForceLayout(diagramWidget, 300)
		diagramWidget.Refresh()
	}), "Node3")
	node3.Move(fyne.Position{X: 400, Y: 100})

	// Node 4
	node4 := diagramwidget.NewDiagramNode(diagramWidget, widget.NewButton("Node4: auto layout", func() {
		forceticks += 100
		diagramWidget.Refresh()
	}), "Node4")
	node4.Move(fyne.Position{X: 400, Y: 400})

	node5 := diagramwidget.NewDiagramNode(diagramWidget, widget.NewButton("Node5: Toggle Background Graphic", func() {
		if diagramWidget.GetBackground() == nil {
			diagramWidget.SetBackground(background)
		} else {
			diagramWidget.SetBackground(nil)
		}
	}), "Node5")
	node5.Move(fyne.NewPos(600, 200))

	node6 := NewModifiedDragNode("Node6", diagramWidget)
	node6.Move(fyne.NewPos(500, 0))

	node7 := NewModifiedPadsNode("Node7", diagramWidget)
	node7Size := node7.Size()
	halfPointPadSize := diagramwidget.PointPadSize / 2
	node7Pads := node7.GetConnectionPads()
	node7DefaultPadPosition := fyne.NewPos(-halfPointPadSize, node7Size.Height/2-halfPointPadSize)
	node7Pads["left"].Move(node7DefaultPadPosition)
	node7RightPadPosition := fyne.NewPos(node7Size.Width-halfPointPadSize, node7Size.Height/2-halfPointPadSize)
	node7Pads["right"].Move(node7RightPadPosition)
	node7.Move(fyne.NewPos(500, 50))
	node7.Refresh()

	// Link0
	link0 := diagramwidget.NewDiagramLink(diagramWidget, "Link0")
	link0.SetSourcePad(node0.GetEdgePad())
	link0.SetTargetPad(node1.GetEdgePad())
	link0.AddSourceAnchoredText("sourceRole", "sourceRole")
	link0.AddMidpointAnchoredText("linkName", "Link 0")
	solidDiamond := createDiamondDecoration()
	solidDiamond.SetSolid(true)
	link0.AddSourceDecoration(solidDiamond)

	// Link1
	link1 := diagramwidget.NewDiagramLink(diagramWidget, "Link1")
	link1.SetSourcePad(node2.GetEdgePad())
	link1.SetTargetPad(node1.GetEdgePad())
	link1.SetForegroundColor(color.RGBA{255, 64, 64, 255})
	link1.AddTargetDecoration(diagramwidget.NewArrowhead())
	link1.AddTargetDecoration(diagramwidget.NewArrowhead())
	link1.AddMidpointDecoration(diagramwidget.NewArrowhead())
	link1.AddMidpointDecoration(diagramwidget.NewArrowhead())
	link1.AddMidpointAnchoredText("linkName", "Link 1")
	link1.AddSourceDecoration(diagramwidget.NewArrowhead())
	link1.AddSourceDecoration(diagramwidget.NewArrowhead())

	// Link2
	link2 := diagramwidget.NewDiagramLink(diagramWidget, "Link2")
	link2.SetSourcePad(node0.GetEdgePad())
	link2.SetTargetPad(node3.GetEdgePad())
	link2.AddMidpointAnchoredText("linkName", "Link 2")
	link2.AddSourceDecoration(createHalfArrowDecoration())

	// Link3
	link3 := diagramwidget.NewDiagramLink(diagramWidget, "Link3")
	link3.SetSourcePad(node2.GetEdgePad())
	link3.SetTargetPad(node3.GetEdgePad())
	link3.AddSourceAnchoredText("sourceRole", "sourceRole")
	link3.AddMidpointAnchoredText("linkName", "Link 3")
	link3.AddTargetAnchoredText("targetRole", "targetRole")
	link3.AddMidpointDecoration(createTriangleDecoration())

	// Link4
	link4 := diagramwidget.NewDiagramLink(diagramWidget, "Link4")
	link4.SetSourcePad(node4.GetEdgePad())
	link4.SetTargetPad(node3.GetEdgePad())
	link4.AddMidpointAnchoredText("linkName", "Link 4")

	// Link5
	link5 := diagramwidget.NewDiagramLink(diagramWidget, "Link5")
	link5.SetSourcePad(link4.GetMidPad())
	link5.SetTargetPad(node5.GetEdgePad())
	link5.AddMidpointAnchoredText("linkName", "Link 5")
	link5.AddTargetDecoration(diagramwidget.NewArrowhead())

	link6 := diagramwidget.NewDiagramLink(diagramWidget, "Link6")
	link6.SetSourcePad(node6.GetEdgePad())
	link6.SetTargetPad(node7.GetConnectionPads()["left"])

	link7 := diagramwidget.NewDiagramLink(diagramWidget, "Link7")
	link7.SetSourcePad(node6.GetEdgePad())
	link7.SetTargetPad(node7.GetConnectionPads()["right"])

	w.SetContent(diagramWidget)

	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}

func createTriangleDecoration() diagramwidget.Decoration {
	points := []fyne.Position{
		{X: 0, Y: 15},
		{X: 15, Y: 0},
		{X: 0, Y: -15},
	}
	polygon := diagramwidget.NewPolygon(points)
	return polygon
}

func createDiamondDecoration() diagramwidget.Decoration {
	points := []fyne.Position{
		{X: 0, Y: 0},
		{X: 8, Y: 4},
		{X: 16, Y: 0},
		{X: 8, Y: -4},
	}
	polygon := diagramwidget.NewPolygon(points)
	return polygon
}

func createHalfArrowDecoration() diagramwidget.Decoration {
	points := []fyne.Position{
		{X: 0, Y: 0},
		{X: 16, Y: 8},
		{X: 16, Y: 0},
	}
	polygon := diagramwidget.NewPolygon(points)
	return polygon
}
