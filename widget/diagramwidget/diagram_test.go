package diagramwidget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestDependencies(t *testing.T) {
	app := test.NewApp()
	assert.NotNil(t, app)
	diagram := NewDiagramWidget("Diagram1")
	node1ID := "Node1"
	node1 := NewDiagramNode(diagram, nil, node1ID)
	node1.Move(fyne.NewPos(100, 100))
	node2ID := "Node2"
	node2 := NewDiagramNode(diagram, nil, node2ID)
	node2.Move(fyne.NewPos(200, 100))
	assert.Equal(t, 0, len(diagram.diagramElementLinkDependencies))
	linkID := "Link1"
	link := NewDiagramLink(diagram, linkID)
	link.SetSourcePad(node1.GetDefaultConnectionPad())
	link.SetTargetPad(node2.GetDefaultConnectionPad())
	assert.NotNil(t, link)
	assert.Equal(t, 2, len(diagram.diagramElementLinkDependencies))

	node1Dependencies := diagram.diagramElementLinkDependencies[node1ID]
	assert.Equal(t, 1, len(node1Dependencies))
	assert.Equal(t, link, node1Dependencies[0].link)
	assert.Equal(t, node1.GetDefaultConnectionPad(), node1Dependencies[0].pad)

	node2Dependencies := diagram.diagramElementLinkDependencies[node2ID]
	assert.Equal(t, 1, len(node2Dependencies))
	assert.Equal(t, link, node2Dependencies[0].link)
	assert.Equal(t, node2.GetDefaultConnectionPad(), node2Dependencies[0].pad)

	// Now test the dependency management when a node is deleted
	diagram.RemoveElement(node2ID)
	assert.Nil(t, diagram.GetDiagramElement(node2ID))
	assert.Nil(t, diagram.GetDiagramElement(linkID))
	assert.Equal(t, 0, len(diagram.diagramElementLinkDependencies))

}
