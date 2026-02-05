# Fyne MultiTab

A powerful and flexible multi-tab widget for Fyne applications, providing enhanced tab management with icons, tooltips, close buttons, keyboard navigation, and event callbacks.

## Features

- **Tab Management**: Add, remove, insert, and move tabs dynamically
- **Visual Enhancements**: Support for tab icons and tooltips
- **User Interaction**: Close buttons and keyboard navigation
- **Event Handling**: Callbacks for tab changes and removals
- **Customization**: Configurable options for closing, reordering, and navigation

## Installation

To use this package in your Fyne project:

```bash
go get github.com/Aswanidev-vs/fyne-multitab
```

## Basic Usage

Here's a simple example of how to create and use the MultiTab widget:

```go
package main

import (
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/widget"
    "github.com/Aswanidev-vs/fyne-multitab/multitab"
)

func main() {
    a := app.New()
    w := a.NewWindow("MultiTab Example")

    tabs := multitab.New()

    // Add some tabs
    tabs.AddTab("Tab 1", widget.NewLabel("Content of Tab 1"))
    tabs.AddTab("Tab 2", widget.NewLabel("Content of Tab 2"))

    w.SetContent(tabs)
    w.ShowAndRun()
}
```

This creates a window with two tabs, each containing a simple label.

## API Reference

### Creating a New Tabs Widget

#### `New() *Tabs`

Creates a new instance of the Tabs widget with default settings.

**Returns:** A pointer to the newly created Tabs widget.

**Example:**
```go
tabs := multitab.New()
```

### Adding Tabs

#### `AddTab(title string, content fyne.CanvasObject)`

Adds a new tab with the specified title and content.

**Parameters:**
- `title`: The title to display on the tab
- `content`: The Fyne canvas object to display as the tab's content

**Example:**
```go
tabs.AddTab("Home", widget.NewLabel("Welcome to the application"))
```

#### `AddTabWithIcon(title string, content fyne.CanvasObject, icon fyne.Resource)`

Adds a new tab with title, content, and an icon.

**Parameters:**
- `title`: The title to display on the tab
- `content`: The Fyne canvas object to display as the tab's content
- `icon`: The icon resource to display on the tab

**Example:**
```go
import "fyne.io/fyne/v2/theme"

tabs.AddTabWithIcon("Settings", settingsWidget, theme.SettingsIcon())
```

#### `AddTabWithIconAndTooltip(title string, content fyne.CanvasObject, icon fyne.Resource, tooltip string)`

Adds a new tab with title, content, icon, and tooltip.

**Parameters:**
- `title`: The title to display on the tab
- `content`: The Fyne canvas object to display as the tab's content
- `icon`: The icon resource to display on the tab
- `tooltip`: The tooltip text to show on hover

**Example:**
```go
tabs.AddTabWithIconAndTooltip("Logs", logsWidget, theme.InfoIcon(), "View application logs")
```

### Removing Tabs

#### `RemoveTab(index int)`

Removes the tab at the specified index.

**Parameters:**
- `index`: The index of the tab to remove (0-based)

**Example:**
```go
// Remove the first tab
tabs.RemoveTab(0)
```

### Tab Navigation and Selection

#### `SetActive(index int)`

Sets the active (selected) tab to the one at the specified index.

**Parameters:**
- `index`: The index of the tab to make active (0-based)

**Example:**
```go
// Switch to the second tab
tabs.SetActive(1)
```

#### `ActiveIndex() int`

Returns the index of the currently active tab.

**Returns:** The index of the active tab, or -1 if no tabs exist.

**Example:**
```go
currentTab := tabs.ActiveIndex()
fmt.Printf("Current active tab: %d\n", currentTab)
```

#### `TabCount() int`

Returns the total number of tabs.

**Returns:** The number of tabs in the widget.

**Example:**
```go
totalTabs := tabs.TabCount()
fmt.Printf("Total tabs: %d\n", totalTabs)
```

### Configuration Options

#### `SetAllowClose(allow bool)`

Enables or disables the close buttons on tabs.

**Parameters:**
- `allow`: `true` to show close buttons, `false` to hide them

**Example:**
```go
// Enable close buttons
tabs.SetAllowClose(true)
```

#### `SetAllowReorder(allow bool)`

Enables or disables tab reordering (currently not implemented in the renderer).

**Parameters:**
- `allow`: `true` to allow reordering, `false` to disable

**Example:**
```go
tabs.SetAllowReorder(true)
```

#### `SetKeyboardNavigation(enable bool)`

Enables or disables keyboard navigation for tabs.

**Parameters:**
- `enable`: `true` to enable keyboard shortcuts, `false` to disable

**Example:**
```go
tabs.SetKeyboardNavigation(true)
```

### Event Callbacks

#### `OnTabChange(callback TabChangeCallback)`

Sets a callback function to be called when the active tab changes.

**Parameters:**
- `callback`: A function that takes an `int` parameter (the new active tab index)

**Callback Type:**
```go
type TabChangeCallback func(index int)
```

**Example:**
```go
tabs.OnTabChange(func(index int) {
    fmt.Printf("Switched to tab %d\n", index)
})
```

#### `OnTabRemoved(callback TabRemovedCallback)`

Sets a callback function to be called when a tab is removed.

**Parameters:**
- `callback`: A function with no parameters

**Callback Type:**
```go
type TabRemovedCallback func()
```

**Example:**
```go
tabs.OnTabRemoved(func() {
    fmt.Println("A tab was removed")
    if tabs.TabCount() == 0 {
        // Close the window if no tabs left
        // window.Close()
    }
})
```

### Retrieving Tab Information

#### `GetTabTitle(index int) string`

Returns the title of the tab at the specified index.

**Parameters:**
- `index`: The index of the tab (0-based)

**Returns:** The title of the tab, or an empty string if the index is invalid.

**Example:**
```go
title := tabs.GetTabTitle(0)
fmt.Printf("First tab title: %s\n", title)
```

#### `GetTabIcon(index int) fyne.Resource`

Returns the icon of the tab at the specified index.

**Parameters:**
- `index`: The index of the tab (0-based)

**Returns:** The icon resource, or `nil` if no icon is set or index is invalid.

**Example:**
```go
icon := tabs.GetTabIcon(0)
if icon != nil {
    // Use the icon
}
```

#### `GetTabTooltip(index int) string`

Returns the tooltip of the tab at the specified index.

**Parameters:**
- `index`: The index of the tab (0-based)

**Returns:** The tooltip text, or an empty string if no tooltip is set or index is invalid.

**Example:**
```go
tooltip := tabs.GetTabTooltip(0)
fmt.Printf("Tooltip: %s\n", tooltip)
```

### Advanced Tab Manipulation

#### `InsertTab(index int, title string, content fyne.CanvasObject)`

Inserts a new tab at the specified position.

**Parameters:**
- `index`: The position to insert the tab (0-based)
- `title`: The title of the new tab
- `content`: The content of the new tab

**Example:**
```go
// Insert a tab at position 1
tabs.InsertTab(1, "New Tab", widget.NewLabel("Inserted content"))
```

#### `MoveTab(from, to int)`

Moves a tab from one position to another.

**Parameters:**
- `from`: The current index of the tab to move
- `to`: The new index for the tab

**Example:**
```go
// Move the first tab to the third position
tabs.MoveTab(0, 2)
```

## Keyboard Navigation

When keyboard navigation is enabled (`SetKeyboardNavigation(true)`), the following shortcuts are available:

- **Tab**: Switch to the next tab
- **T**: Create a new tab
- **W**: Close the current tab

**Example:**
```go
tabs.SetKeyboardNavigation(true)
// Now users can use Tab, T, and W keys for navigation
```

## Advanced Example

Here's a more comprehensive example showcasing various features:

```go
package main

import (
    "fmt"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"
    "github.com/Aswanidev-vs/fyne-multitab/multitab"
)

func main() {
    a := app.New()
    w := a.NewWindow("Enhanced MultiTab Example")
    w.Resize(fyne.NewSize(900, 600))

    tabs := multitab.New()

    // Configure tabs
    tabs.SetAllowClose(true)
    tabs.SetKeyboardNavigation(true)

    // Tab change callback
    tabs.OnTabChange(func(index int) {
        fmt.Printf("Switched to tab %d: %s\n", index, tabs.GetTabTitle(index))
    })

    // Tab removed callback
    tabs.OnTabRemoved(func() {
        if tabs.TabCount() == 0 {
            w.Close()
        }
    })

    // Home tab with icon and tooltip
    homeContent := container.NewVBox(
        widget.NewLabel("Welcome to Enhanced MultiTab"),
        widget.NewButton("Add New Tab", func() {
            count := tabs.TabCount() + 1
            tabs.AddTab(fmt.Sprintf("Tab %d", count), widget.NewLabel(fmt.Sprintf("Content for tab %d", count)))
        }),
    )

    tabs.AddTabWithIconAndTooltip("Home", homeContent, theme.HomeIcon(), "Welcome page")

    // Settings tab
    settingsContent := container.NewVBox(
        widget.NewLabel("Tab Configuration:"),
        widget.NewCheck("Allow closing tabs", func(b bool) {
            tabs.SetAllowClose(b)
        }),
        widget.NewCheck("Keyboard navigation", func(b bool) {
            tabs.SetKeyboardNavigation(b)
        }),
    )

    tabs.AddTabWithIconAndTooltip("Settings", settingsContent, theme.SettingsIcon(), "Configure application")

    w.SetContent(tabs)
    w.ShowAndRun()
}
```

This example demonstrates:
- Creating tabs with icons and tooltips
- Configuring widget options
- Using callbacks for tab events
- Dynamic tab creation
- Interactive settings

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This project is licensed under the MIT License.
