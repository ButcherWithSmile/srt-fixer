package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/encoding/charmap"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// This code creates a new Fyne application and window. It also sets the window's icon
	myApp := app.New()
	myWindow := myApp.NewWindow("SRT Fixer")

	r, _ := fyne.LoadResourceFromURLString("https://raw.githubusercontent.com/ButcherWithSmile/srt-fixer/main/icon.png")
	myWindow.SetIcon(r)

	myWindow.Resize(fyne.NewSize(600, 400))

	/* This code creates a file entry widget, a select file button, and a fix SRT button.

	The select file button opens a dialog box that allows the user to select an SRT file.
 	Once the user selects a file, the file path is displayed in the file entry widget.

	The fix SRT button fixes the SRT file by decoding it from Windows-1256 encoding to UTF-8 encoding. 
	The fixed file is saved with the same name as the original file but with the "_fixed" suffix. */
	fileEntry := widget.NewEntry()
	fileEntry.Disable()

	selectFileButton := widget.NewButton("Select SRT File", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				fmt.Println("Error opening file:", err)
				return
			}
			if reader == nil {
				return
			}

			filePath := strings.TrimPrefix(reader.URI().String(), "file://")
			fileEntry.SetText(filePath)

			reader.Close()
		}, myWindow)
	})

	fixButton := widget.NewButton("Fix SRT", func() {
		filePath := fileEntry.Text
		if filePath == "" {
			dialog.ShowError(fmt.Errorf("Please select SRT file"), myWindow)
			return
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		decoder := charmap.Windows1256.NewDecoder()
		decodedData, err := decoder.Bytes(data)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		dir, fileName := filepath.Split(filePath)
		fileName = strings.TrimSuffix(fileName, ".srt")
		newFileName := filepath.Join(dir, fileName+"_fixed.srt")

		err = os.WriteFile(newFileName, decodedData, 0644)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		dialog.ShowInformation("Success", fmt.Sprintf("Fixed file saved as %s", newFileName), myWindow)
	})

	// This code creates an About tab and an App tabs widget
	// The app tabs widget contains two tabs: Main and About
	// The Main tab contains the select file button, file entry widget, and fix SRT button
	aboutTab := container.NewVBox(
		widget.NewLabel("SRT Fixer"),
		widget.NewLabel("Version: 1.0.0"),
		widget.NewLabel("Developer: Hatef PourRajabi"),
		widget.NewLabel("Email: hatef.pr@gmail.com"),
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Main", container.NewVBox(
			selectFileButton,
			fileEntry,
			fixButton,
		)),
		container.NewTabItem("About", aboutTab),
	)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
