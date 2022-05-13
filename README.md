<p align="center">
  <a href="https://pkg.go.dev/fyne.io/x/fyne" title="Go API Reference" rel="nofollow"><img src="https://img.shields.io/badge/go-documentation-blue.svg?style=flat" alt="Go API Reference"></a>
  <a href='http://gophers.slack.com/messages/fyne'><img src='https://img.shields.io/badge/join-us%20on%20slack-gray.svg?longCache=true&logo=slack&colorB=blue' alt='Join us on Slack' /></a>
  <br />
  <a href="https://goreportcard.com/report/fyne.io/x/fyne"><img src="https://goreportcard.com/badge/fyne.io/x/fyne" alt="Code Status" /></a>
  <a href="https://github.com/fyne-io/fyne-x/actions"><img src="https://github.com/fyne-io/fyne-x/workflows/Platform%20Tests/badge.svg" alt="Build Status" /></a>
  <a href='https://coveralls.io/github/fyne-io/fyne-x?branch=master'><img src='https://coveralls.io/repos/github/fyne-io/fyne-x/badge.svg?branch=master' alt='Coverage Status' /></a>
</p>

# About

This repository holds community extensions for the [Fyne](https://fyne.io) toolkit.

This is in early development and more information will appear soon.

## Layouts

Community contributed layouts.

`import "fyne.io/x/fyne/layout"`

## Widgets

Community contributed widgets.

`import "fyne.io/x/fyne/widget"`

### Animated Gif

A widget that will run animated gifs.

```go
gif, err := NewAnimatedGif(storage.NewFileURI("./testdata/gif/earth.gif"))
gif.Start()
```

### FileTree

An extension of widget.Tree for displaying a file system hierarchy.

```go
tree := widget.NewFileTree(storage.NewFileURI("~")) // Start from home directory
tree.Filter = storage.NewExtensionFileFilter([]string{".txt"}) // Filter files
tree.Sorter = func(u1, u2 fyne.URI) bool {
    return u1.String() < u2.String() // Sort alphabetically
}
```

<p align="center" markdown="1" style="max-width: 100%">
  <img src="img/widget-filetree.png" width="1024" height="880" alt="FileTree Widget" style="max-width: 100%" />
</p>

### CompletionEntry

An extension of widget.Entry for displaying a popup menu for completion. The "up" and "down" keys on the keyboard are used to navigate through the menu, the "Enter" key is used to confirm the selection. The options can also be selected with the mouse. The "Escape" key closes the selection list.

```go
entry := widget.NewCompletionEntry([]string{})

// When the use typed text, complete the list.
entry.OnChanged = func(s string) {
    // completion start for text length >= 3
    if len(s) < 3 {
        entry.HideCompletion()
        return
    }

    // Make a search on wikipedia
    resp, err := http.Get(
        "https://en.wikipedia.org/w/api.php?action=opensearch&search=" + entry.Text,
    )
    if err != nil {
        entry.HideCompletion()
        return
    }

    // Get the list of possible completion
    var results [][]string
    json.NewDecoder(resp.Body).Decode(&results)

    // no results
    if len(results) == 0 {
        entry.HideCompletion()
        return
    }

    // then show them
    entry.SetOptions(results[1])
    entry.ShowCompletion()
}
```

<p align="center" markdown="1" style="max-width: 100%">
  <img src="img/widget-completion-entry.png" width="825" height="634" alt="CompletionEntry Widget" style="max-width: 100%" />
</p>

### 7-Segment ("Hex") Display

A skeuomorphic widget simulating a 7-segment "hex" display. Supports setting
digits by value, as well as directly controlling which segments are on or
off.

Check out the [demo](./cmd/hexwidget_demo/main.go) for an example of usage.


![](img/hexwidget_00abcdef.png)

![](img/hexwidget_12345678.png)

```go
h := widget.NewHexWidget()
// show the value 'F' on the display
h.Set(0xf)
```

### Map

An OpenStreetMap widget that can the user can pan and zoom.
To use this in your app and be compliant with their requirements you may need to request
permission to embed in your specific software.

```go
m := NewMap()
```

![](img/map.png)

## Data Binding

Community contributed data sources for binding.

`import fyne.io/x/fyne/data/binding`

### WebString

A `WebSocketString` binding creates a `String` data binding to the specified web socket URL.
Each time a message is read the value will be converted to a `string` and set on the binding.
It is also `Closable` so you should be sure to call `Close()` once you are completed using it.

```go
s, err := binding.NewWebSocketString("wss://demo.piesocket.com/v3/channel_1?api_key=oCdCMcMPQpbvNjUIzqtvF1d2X2okWpDQj4AwARJuAgtjhzKxVEjQU6IdCjwm&notify_self")
l := widget.NewLabelWithData(s)
```

The code above uses a test web sockets server from "PieSocket", you can run the code above
and go to [their test page](https://www.piesocket.com/websocket-tester) to send messages.
The widget will automatically update to the latest data sent through the socket.

### MqttString

A `MqttString` binding creates a `String` data binding to the specified _topic_ associated with
the specified **MQTT** client connection. Each time a message is received the value will be converted
to a `string` and set on the binding. Each time the value is edited, it will be sent back over
**MQTT** on the specified _topic_. It is also a `Closer` so you should be sure to call `Close`
once you are completed using it to disconnect the _topic_ handler from the **MQTT** client connection.

```go
opts := mqtt.NewClientOptions()
opts.AddBroker("tcp://broker.emqx.io:1883")
opts.SetClientID("fyne_demo")
client := mqtt.NewClient(opts)

token := client.Connect()
token.Wait()
if err := token.Error(); err != nil {
    // Handle connection error
}

s, err := binding.NewMqttString(client, "fyne.io/x/string")
```

## Data Validation

Community contributed validators.

`import fyne.io/x/fyne/data/validation`

### Password

A validator for validating passwords. Uses https://github.com/wagslane/go-password-validator
for validation using an entropy system.

```go
pw := validation.NewPassword(70) // Minimum password entropy allowed defined as 70.
```

### Charts

The `widget/charts` package provides some chart widget to draw:

- `BarChar` which is an histrogram view
- `LineChart` wich is simple line chart

You can change colors, override behavior, Y range...

![Chart example](img/chart.png)

See the [`cmd/graph_demo`](cmd/graph_demo) or run from this repository:

```bash
go run cmd/graph_demo/*
```

The demo provides several examples like an overriden view to make the "mouse over" event displaying values and pointer.

![Chart with mouse event](img/chart-mouse.png)

Basic example to create a line chart:

```go
package main

import (
    "fyne.io/fyne/v2"
    "fyne.io/x/fyne/widget/charts"
)

func main() {
    app := app.New()
    w := app.NewWindow("Graphs")

    data := []float64{1.2, 5.3, 2.2, 3.4}
    chart := charts.NewLineChart(nil) // nil = default options
    chart.SetData(data)

    w.SetContent(chart) 
    w.ShowAndRun()
}
```

