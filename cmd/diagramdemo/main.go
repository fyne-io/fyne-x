package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/x/fyne/widget/diagramwidget"
	"fyne.io/x/fyne/widget/diagramwidget/arrowhead"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var forceticks int = 0

var globaldiagram *diagramwidget.DiagramWidget

func forceanim() {

	// XXX: very naughty -- accesses shared memory in potentially unsafe
	// ways, this almost certainly has race conditions... don't do this!

	for {
		if forceticks > 0 {
			globaldiagram.StepForceLayout(300)
			globaldiagram.Refresh()
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

	g := diagramwidget.NewDiagram()

	go forceanim()

	l := widget.NewLabel("teeexxttt")
	n := diagramwidget.NewDiagramNode(g, l)
	g.Nodes["node0"] = n
	n1 := n

	b := widget.NewButton("button", func() { fmt.Printf("tapped!\n") })
	n = diagramwidget.NewDiagramNode(g, b)
	n.Move(fyne.Position{X: 200, Y: 200})
	g.Nodes["node1"] = n
	n2 := n

	n = diagramwidget.NewDiagramNode(g, nil)
	c := container.NewVBox(
		widget.NewLabel("Fancy node!"),
		widget.NewButton("Up", func() {
			n.Displace(fyne.Position{X: 0, Y: -10})
			n.Refresh()
		}),
		widget.NewButton("Down", func() {
			n.Displace(fyne.Position{X: 0, Y: 10})
			n.Refresh()
		}),
		container.NewHBox(
			widget.NewButton("Left", func() {
				n.Displace(fyne.Position{X: -10, Y: 0})
				n.Refresh()
			}),
			widget.NewButton("Right", func() {
				n.Displace(fyne.Position{X: 10, Y: 0})
				n.Refresh()
			}),
		),
	)
	n.InnerObject = c
	n.Move(fyne.Position{X: 300, Y: 300})
	g.Nodes["node2"] = n
	n3 := n

	n = diagramwidget.NewDiagramNode(g, widget.NewButton("force layout step", func() {
		g.StepForceLayout(300)
		g.Refresh()
	}))
	n.Move(fyne.Position{X: 400, Y: 200})
	g.Nodes["node4"] = n
	n4 := n

	n = diagramwidget.NewDiagramNode(g, widget.NewButton("auto layout", func() {
		forceticks += 100
		g.Refresh()
	}))
	n.Move(fyne.Position{X: 400, Y: 500})
	g.Nodes["node5"] = n
	n5 := n

	globaldiagram = g

	g.Links["edge0"] = diagramwidget.NewDiagramLink(g, n1, n2)
	edge1 := diagramwidget.NewDiagramLink(g, n3, n2)
	g.Links["edge1"] = edge1
	edge1.LinkColor = color.RGBA{255, 64, 64, 255}
	edge1.TargetDecorations = append(g.Links["edge1"].TargetDecorations, arrowhead.NewArrowhead())
	edge1.MidpointDecorations = append(g.Links["edge1"].TargetDecorations, arrowhead.NewArrowhead())
	edge1.SourceDecorations = append(g.Links["edge1"].SourceDecorations, arrowhead.NewArrowhead())
	g.Links["edge2"] = diagramwidget.NewDiagramLink(g, n1, n4)
	g.Links["edge3"] = diagramwidget.NewDiagramLink(g, n3, n4)
	g.Links["edge4"] = diagramwidget.NewDiagramLink(g, n5, n4)

	w.SetContent(g)

	w.ShowAndRun()
}
