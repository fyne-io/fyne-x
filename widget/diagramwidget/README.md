# Fyne Hacks

This repository contains a collection of widgets and programs for the
[Fyne](https://fyne.io/) toolkit. In general, code in here should not be
considered production-ready. Most of it is either things I'm experimenting
with, or prototypes of things I will later upstream. Expect
backwards-incompatible changes without warning.

If you have a question, comment, or patch, you should send it in via [my public
inbox](https://lists.sr.ht/~charles/public-inbox).

## Hacks

### Table Widget

* [code](./table)
* [demo](./cmd/tabledemo)


### Viewport Widget

The viewport widget provides a canvas-like object which can be zoomed and
panned.

* [code](./viewport)
* [demo](./cmd/viewportdemo)

### Graph Widget

The graph widget implements a graph visualization widget which supports
directed and un-directed edges. Widgets can be embedded within nodes, and nodes
can be moved around.

* [code](./graph)
* [demo](./cmd/graphdemo)

