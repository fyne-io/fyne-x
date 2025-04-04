package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/x/fyne/widget/diagramwidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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

func main() {
	app := app.New()
	w := app.NewWindow("Diagram Demo")

	w.SetMaster()

	diagramWidget := diagramwidget.NewDiagramWidget("Diagram1")

	scrollContainer := container.NewScroll(diagramWidget)

	go forceanim(diagramWidget)

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

	node5 := diagramwidget.NewDiagramNode(diagramWidget, widget.NewLabel("Node5"), "Node5")
	node5.Move(fyne.NewPos(600, 200))

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

	w.SetContent(scrollContainer)

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
