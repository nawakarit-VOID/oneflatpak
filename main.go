package main

import (
	"os"
	"path/filepath"
	"text/template"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type AppConfig struct {
	Name        string
	AppID       string
	Command     string
	Icon        string
	Categories  string
	Summary     string
	Description string
	License     string
	Developer   string
}

func generateFile(tmplPath, outputPath string, data AppConfig) error {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

func main() {
	a := app.New()
	w := a.NewWindow("Flatpak GUI Builder")

	// form inputs
	name := widget.NewEntry()
	name.SetPlaceHolder("App Name")

	appID := widget.NewEntry()
	appID.SetPlaceHolder("com.example.app")

	command := widget.NewEntry()
	command.SetPlaceHolder("binary name")

	icon := widget.NewEntry()
	icon.SetPlaceHolder("icon name ไม่ต้องใส่ .png")

	categories := widget.NewEntry()
	categories.SetPlaceHolder("Utility;")

	summary := widget.NewEntry()
	summary.SetPlaceHolder("Short summary")

	description := widget.NewMultiLineEntry()
	description.SetPlaceHolder("Description")

	developer := widget.NewEntry()
	developer.SetPlaceHolder("Your name")

	output := widget.NewLabel("")

	generateBtn := widget.NewButton("Generate Files", func() {
		cfg := AppConfig{
			Name:        name.Text,
			AppID:       appID.Text,
			Command:     command.Text,
			Icon:        icon.Text,
			Categories:  categories.Text,
			Summary:     summary.Text,
			Description: description.Text,
			License:     "MIT",
			Developer:   developer.Text,
		}

		os.MkdirAll("output", 0755)

		err1 := generateFile("templates/desktop.tmpl",
			filepath.Join("output", cfg.AppID+".desktop"), cfg)

		err2 := generateFile("templates/manifest.tmpl",
			filepath.Join("output", cfg.AppID+".json"), cfg)

		err3 := generateFile("templates/metainfo.tmpl",
			filepath.Join("output", cfg.AppID+".metainfo.xml"), cfg)

		if err1 != nil || err2 != nil || err3 != nil {
			output.SetText("❌ Error generating files")
		} else {
			output.SetText("✅ Files generated in /output")
		}
	})

	form := container.NewVBox(
		name, appID, command, icon,
		categories, summary, description, developer,
		generateBtn,
		output,
	)

	w.SetContent(form)
	w.Resize(fyne.NewSize(400, 600))
	w.ShowAndRun()
}
