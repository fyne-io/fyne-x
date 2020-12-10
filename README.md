<p align="center">
  <a href="https://pkg.go.dev/fyne.io/fyne-x?tab=doc" title="Go API Reference" rel="nofollow"><img src="https://img.shields.io/badge/go-documentation-blue.svg?style=flat" alt="Go API Reference"></a>
  <a href='http://gophers.slack.com/messages/fyne'><img src='https://img.shields.io/badge/join-us%20on%20slack-gray.svg?longCache=true&logo=slack&colorB=blue' alt='Join us on Slack' /></a>
  <br />
  <a href="https://goreportcard.com/report/fyne.io/x/fyne"><img src="https://goreportcard.com/badge/fyne.io/x/fyne" alt="Code Status" /></a>
  <a href="https://github.com/fyne-io/fyne-x/actions"><img src="https://github.com/fyne-io/fyne-x/workflows/Platform%20Tests/badge.svg" alt="Build Status" /></a>
  <a href='https://coveralls.io/github/fyne-io/fyne-x?branch=develop'><img src='https://coveralls.io/repos/github/fyne-io/fyne-x/badge.svg?branch=develop' alt='Coverage Status' /></a>
</p>

# About

This repository holds community extensions for the [Fyne](https://fyne.io) toolkit.

This is in early development and more information will appear soon.

## Layouts

Community contributed layouts.

`import fyne.io/x/fyne/layout`

## Widgets

Community contributed widgets.

`import fyne.io/x/fyne/widget`

### FileTree

An extension of widget.Tree for displaying a file system hierarchy.

```go
package main

import (
    "fyne.io/fyne"
    "fyne.io/fyne/app"
    "fyne.io/fyne/storage"
    "fyne.io/x/fyne/widget"
    "os"
)

func main() {
    a := app.New()
    w := a.NewWindow("FileTree")
    dir, err := os.Getwd()
    if err != nil {
        fyne.LogError("Could not get current working directory", err)
        return
    }
    w.SetContent(widget.NewFileTree(storage.NewFileURI(dir)))
    w.Resize(fyne.NewSize(400, 300))
    w.ShowAndRun()
}
```

<p align="center" markdown="1" style="max-width: 100%">
  <img src="img/widget-filetree.png" width="1024" height="868" alt="FileTree Widget" style="max-width: 100%" />
</p>
