# Fyne DiagramWidget

This package contains a collection of widgets for the [Fyne](https://fyne.io/) 
toolkit. The code here is intended to be production ready, but may be lacking
some desirable functional features. If you have suggestions for changes to 
existing functionality or addition of new functionality, please look at the existing
issues in the repository to see if your idea is already on the table. If it is not,
feel free to open an issue. 

This collection should be considered a work in progress. When changes are made,
serious consideration will be given to backward compatibility, but compatibility
is not guaranteed. 

## Diagram Widget

The DiagramWidget is intended to be incorporated into a Fyne application. It provides a
drawing area within which a diagram can be created. The diagram itself is a collection of 
DiagramElement widgets (an interface). There are two types of DiagramElements: DiagramNode widgets and DiagramLink widgets. DiagramNode widgets are thin wrappers around a user-supplied CanvasObject.
Any valid CanvasObject can be used. DiagramLinks are line-based connections between DiagramElements.
Note that links can connect to other links as well as nodes.

While some provisions have been made for automatic layout, layouts are for the convenience
of the author and are on-demand only. The design intent is that users will place the diagram elements for human readability. 

DiagramElements are essentially self-managed from a layout perspective. DiagramNodes have no size
constraints imposed by the DiagramWidget and can be placed anywhere. DiagramLinks connect 
DiagramElements. The DiagramWidget keeps track of the DiagramElements to which each DiagramLink 
is connected and calls the Refresh() method on the link when the connected diagram element is moved 
or resized. 

* [demo](../../cmd/diagramdemo/main.go)

### DiagramElement Interface

A DiagramElement is the base interface for any element of the diagram being managed by the 
DiagramWidget. It provides a common interface for DiagramNode and DiagramLink widgets. The DiagramElement
interface provides operations for retrieving the DiagramWidget, the ID of the DiagramElement, and
for showing and hiding the handles that are used for graphically manipulating the diagram element.
The specifics of what handles do are different for nodes and links - these are described below in the
sections for their respective widgets.

### DiagramNode Widget

The DiagramNode widget is a wrapper around a user-supplied CanvasObject. In addition to the user-supplied
CanvasObject, the node displays a border and, when selected, handles at the corners and edge mid-points that can be used to manipulate the size of the node. The node can be selected and dragged to a new position with a mouse by clicking in the border area around the canvas object. 

### DiagramLink Widget

The DiagramLink widget provides a directed line-based connection between two DiagramElements. 
The link is defined in terms of LinkPoints that are connected by LinkSegments (both of which
are widgets in their own right). The link maintains an array of points, with the point at index
[0] being the point at which the link connects to the source DiagramElement and the point at the 
last index being the point at which the link connects to the target DiagramElement. The link also
maintains an array of line segments, with the segment at index [0] connecting points [0] and [1], 
the segment at index [1] connecting the points [1] and [2], etc. The current implementation only
has a single segment, but interfaces will be added shortly to enable the addition and removal of
points and segments.

Many visual languages (formalized diagrams) utilize graphical decorations on lines. The link
provides the ability to add an arbitrary number of graphic decorations at three points along 
the link: the source end, the target end, and the midpoint. Decorations are stacked in the order
they are added at the indicated point. The location of the source and target points is obvious,
but the midpoint bears some discussion. If there is only one line segment, the midpoint is the
midpoint of this segment. If there is more than one line segment, the "midpoint" is defined to
be the next to last point in the array of points. For a two-segment link, this will be the point
at which the two segments join. For a multi-segment link, this will be the point at which the 
next-to-last and last segments join.

Also common in visual languages are textual annotations associated with either the link as a whole 
or to the ends of the link. For this purpose, the link allows the association of one or more 
AnchoredText widgets with each of the reference points on the link: source, target, and midpoint.
These widgets keep track of their position relative to the link's reference points. They can 
be moved interactively with the mouse to a new position. When the reference point on the link
moves, the anchored text will also move, maintaining its relative position. 

Users do not create AnchoredText widgets directly: the link itself creates and manages them. 
the user calls Add<position>AnchoredText(key, text) to add an anchored text. The key is expected
to be unique at the position and can be used to update the text later. The AnchoredText can also
be directly edited in the diagram.  

When a link connects to another link, it connects at the midpoint of the source or target link.