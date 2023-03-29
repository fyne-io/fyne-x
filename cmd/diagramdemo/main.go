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

func forceanim() {

	// XXX: very naughty -- accesses shared memory in potentially unsafe
	// ways, this almost certainly has race conditions... don't do this!

	for {
		if forceticks > 0 {
			diagramwidget.Globaldiagram.StepForceLayout(300)
			diagramwidget.Globaldiagram.Refresh()
			forceticks--
			fmt.Printf("forceticks=%d\n", forceticks)
		}

		time.Sleep(time.Millisecond * (1000 / 30))
	}
}

func main() {
	app := app.New()
	w := app.NewWindow("Diagram Demo")

	w.SetMaster()

	diagramWidget := diagramwidget.NewDiagramWidget("Diagram1")
	diagramwidget.Globaldiagram = diagramWidget

	go forceanim()

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
		diagramWidget.StepForceLayout(300)
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
	link0 := diagramwidget.NewDiagramLink(diagramWidget, node0.GetEdgePad(), node1.GetEdgePad(), "Link0")
	link0.AddSourceAnchoredText("sourceRole", "sourceRole")
	link0.AddMidpointAnchoredText("linkName", "Link 0")

	// Link1
	link1 := diagramwidget.NewDiagramLink(diagramWidget, node2.GetEdgePad(), node1.GetEdgePad(), "Link1")
	link1.LinkColor = color.RGBA{255, 64, 64, 255}
	link1.AddTargetDecoration(diagramwidget.NewArrowhead())
	link1.AddTargetDecoration(diagramwidget.NewArrowhead())
	link1.AddMidpointDecoration(diagramwidget.NewArrowhead())
	link1.AddMidpointDecoration(diagramwidget.NewArrowhead())
	link1.AddMidpointAnchoredText("linkName", "Link 1")
	link1.AddSourceDecoration(diagramwidget.NewArrowhead())
	link1.AddSourceDecoration(diagramwidget.NewArrowhead())

	// Link2
	link2 := diagramwidget.NewDiagramLink(diagramWidget, node0.GetEdgePad(), node3.GetEdgePad(), "Link2")
	link2.AddMidpointAnchoredText("linkName", "Link 2")

	// Link3
	link3 := diagramwidget.NewDiagramLink(diagramWidget, node2.GetEdgePad(), node3.GetEdgePad(), "Link3")
	link3.AddSourceAnchoredText("sourceRole", "sourceRole")
	link3.AddMidpointAnchoredText("linkName", "Link 3")
	link3.AddTargetAnchoredText("targetRole", "targetRole")

	// Link4
	link4 := diagramwidget.NewDiagramLink(diagramWidget, node4.GetEdgePad(), node3.GetEdgePad(), "Link4")
	link4.AddMidpointAnchoredText("linkName", "Link 4")

	// Link5
	link5 := diagramwidget.NewDiagramLink(diagramWidget, link4.GetMidPad(), node5.GetEdgePad(), "Link5")
	link5.AddMidpointAnchoredText("linkName", "Link 5")
	link5.AddTargetDecoration(diagramwidget.NewArrowhead())

	w.SetContent(diagramWidget)

	w.ShowAndRun()
}
