package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"graphics-lab-2/gui"
	"graphics-lab-2/imaging"
	"graphics-lab-2/winapi"
	"log"
	"time"
)

var (
	mw *walk.MainWindow
	filesTV *walk.TableView
	folderLB *walk.Label
)

func getMetadataWithTime(selectedDir string) (metadata []*imaging.ImageMetadata, totalTime time.Duration, err error) {
	var fileNames []string
	startTime := time.Now()
	if fileNames, err = imaging.GetImagesFromDir(selectedDir); err != nil {
		return nil, 0, err
	}
	if metadata, err = imaging.FetchMetadata(fileNames...); err != nil {
		return nil, 0, err
	}
	endTime := time.Now()

	return metadata, endTime.Sub(startTime), nil
}

func browseAndUpdate() {
	var (
		path string
		failMsg string
		metadata []*imaging.ImageMetadata
		totalTime  time.Duration
		err error
	)
	if path, err = winapi.OpenDirectory(); err != nil {
		failMsg = "Can't change the directory: "+err.Error()
		goto fail
	}
	if path == "" {
		return
	}
	if metadata, totalTime, err = getMetadataWithTime(path); err != nil {
		failMsg = "Can't fill the table: "+err.Error()
		goto fail
	}
	filesTV.Model().(*gui.ImageMetadataModel).SetMetadata(metadata)
	_ = folderLB.SetText(path)
	winapi.MessageBox(fmt.Sprintf("Fetched %d files in %v", len(metadata), totalTime),
					  "Folder scan results", win.MB_ICONINFORMATION | win.MB_OK)
	return

fail:
	winapi.MessageBox(failMsg, "Error", win.MB_ICONERROR|win.MB_OK)
}

func main() {
	mw := MainWindow{
		AssignTo: &mw,
		Name:   "Computer Graphics Lab 2",
		Size:   Size{Width: 700, Height: 400},
		Layout: VBox{},
		Children: []Widget{
			VSplitter{
				Children: []Widget {
					GroupBox{
						Title:   "Folder selection",
						Layout: HBox{},
						Children: []Widget{
							Label{
								Text: "Folder selected: ",
								MaxSize: Size{Width: 120, Height: 20},
							},
							Label{
								AssignTo: &folderLB,
								Text: "No folder so far",
							},
							PushButton{
								Alignment: AlignHFarVCenter,
								Text: "Browse",
								MaxSize: Size{Width: 100, Height: 20},
								OnClicked: browseAndUpdate,
							},
						},
					},
					GroupBox{
						Title: "Files in the folder",
						Layout: HBox{},
						Children: []Widget {
							TableView{
								AssignTo:         &filesTV,
								ColumnsOrderable: true,
								Columns: []TableViewColumn{
									{Title: "Name"},
									{Title: "Size"},
									{Title: "X Resolution"},
									{Title: "Y Resolution"},
									{Title: "Color Depth"},
									{Title: "Compression Type"},
								},
								Model: gui.NewImageMetadataModel([]*imaging.ImageMetadata{}),
							},
						},
					},
				},
			},
		},
	}

	if _, err := mw.Run(); err != nil {
		log.Fatal(err)
	}
}
