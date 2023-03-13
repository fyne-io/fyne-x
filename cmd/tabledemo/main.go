package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"

	"fyne.io/x/fyne/widget/diagramwidget/table"

	"github.com/rocketlaunchr/dataframe-go"
	"github.com/rocketlaunchr/dataframe-go/imports"
)

var mainTable *table.TableWidget

func main() {

	s1 := dataframe.NewSeriesInt64("day", nil, 1, 2, 3, 4, 5, 6, 7, 8)
	s2 := dataframe.NewSeriesFloat64("sales", nil, 50.3, 23.4, 56.2, nil, nil, 84.2, 72, 89)
	s3 := dataframe.NewSeriesString("string!", nil, "foo", "bar", "three", "four", "five", "six", "seven", "eight")
	df := dataframe.NewDataFrame(s1, s2, s3)

	fmt.Print(df.Table())

	app := app.New()
	w := app.NewWindow("Table Demo")

	w.SetMainMenu(
		fyne.NewMainMenu(
			fyne.NewMenu("File",
				fyne.NewMenuItem("Load CSV", func() {
					dialog.ShowFileOpen(func(uri fyne.URIReadCloser, e error) {

						if e != nil {
							dialog.ShowError(e, w)
							return
						}

						content, err := ioutil.ReadAll(uri)
						if err != nil {
							dialog.ShowError(err, w)
							return
						}
						text := string(content)
						reader := strings.NewReader(text)

						ctx := context.Background()
						opts := imports.CSVLoadOptions{
							InferDataTypes: true,
						}
						loaded, err := imports.LoadFromCSV(ctx, reader, opts)

						if err != nil {
							dialog.ShowError(e, w)
							return
						}

						fmt.Printf("loaded new table:\n%s\n", loaded.Table())

						mainTable.ReplaceDataFrame(loaded)

					}, w)
				}),
			),
		),
	)
	w.SetMaster()

	mainTable = table.NewTableWidget(df)

	w.SetContent(mainTable)

	w.ShowAndRun()

}
