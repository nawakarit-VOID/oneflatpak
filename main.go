package main

import (
	"embed"
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type AppConfig struct {
	Name        string
	AppID       string
	Command     string
	Categories  string
	Summary     string
	Description string
	License     string
	Developer   string
}

// ============================================================================
// ฟังชั้น gen + run template
// ============================================================================
func generateFile(tmplPath, outputPath string, data AppConfig) error {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	//projectPath
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

// ============================================================================
// ฟังชั้น build เป็น flatpak
// ============================================================================
func runScriptbuildflatpak(projectPath string, output *widget.Entry) {

	commands := [][]string{
		{"gnome-terminal", "--", "bash", "-c", "cd '" + projectPath + "' && chmod +x buildflatpak.sh && ./buildflatpak.sh; exec bash"},
		{"x-terminal-emulator", "-e", "bash", "-c", "cd '" + projectPath + "' && chmod +x buildflatpak.sh && ./buildflatpak.sh; exec bash"},
		{"konsole", "-e", "bash", "-c", "cd '" + projectPath + "' && chmod +x buildflatpak.sh && ./buildflatpak.sh; exec bash"},
		{"xfce4-terminal", "-e", "bash", "-c", "cd '" + projectPath + "' && chmod +x buildflatpak.sh && ./buildflatpak.sh; exec bash"},
	}

	for _, c := range commands {
		cmd := exec.Command(c[0], c[1:]...)
		err := cmd.Start()
		if err == nil {
			output.SetText("🚀 opened terminal: " + c[0])
			return
		}
	}

	output.SetText("❌ no terminal found")
}

// ============================================================================
// ฟังชั้น build Icons
// ============================================================================
func runScriptbuildIcons(projectPath string, output *widget.Entry) {

	commands := [][]string{ //ใช้ imagemagick
		{"gnome-terminal", "--", "bash", "-c", "cd '" + projectPath + "' && chmod +x buildicons.sh && ./buildicons.sh; exec bash"},
		{"x-terminal-emulator", "-e", "bash", "-c", "cd '" + projectPath + "' && chmod +x buildicons.sh && ./buildicons.sh; exec bash"},
		{"konsole", "-e", "bash", "-c", "cd '" + projectPath + "' && chmod +x buildicons.sh && ./buildicons.sh; exec bash"},
		{"xfce4-terminal", "-e", "bash", "-c", "cd '" + projectPath + "' && chmod +x buildicons.sh && ./buildicons.sh; exec bash"},
	}

	for _, c := range commands {
		cmd := exec.Command(c[0], c[1:]...)
		err := cmd.Start()
		if err == nil {
			output.SetText("🚀 opened terminal: " + c[0])
			return
		}
	}

	output.SetText("❌ no terminal found")
}

// โหลด icon
func loadIcon(size int) fyne.Resource {
	var file string

	switch {
	case size >= 512:
		file = "icons/icon-512.png" ///ที่อยู่
	case size >= 256:
		file = "icons/icon-256.png"
	case size >= 128:
		file = "icons/icon-128.png"
	default:
		file = "icons/icon-64.png"
	}

	data, _ := iconFS.ReadFile(file)
	return fyne.NewStaticResource(file, data)
}

//go:embed icons/*
var iconFS embed.FS

func main() {

	a := app.NewWithID("com.nawakarit.flatpak")
	icons := loadIcon(64) //เอา data มาใช้
	a.SetIcon(icons)
	w := a.NewWindow("flatpak")
	w.SetIcon(icons)

	// inputs
	name := widget.NewEntry()
	name.SetPlaceHolder("App Name")

	appID := widget.NewEntry()
	appID.SetPlaceHolder("com.example.app")

	command := widget.NewEntry()
	command.SetPlaceHolder("binary name")

	categories := widget.NewEntry()
	categories.SetPlaceHolder("Utility;")

	summary := widget.NewEntry()
	summary.SetPlaceHolder("Short summary")

	description := widget.NewMultiLineEntry()
	description.SetPlaceHolder("Description")

	developer := widget.NewEntry()
	developer.SetPlaceHolder("Your name")

	// 🔥 log box
	logBox := widget.NewMultiLineEntry()
	logBox.SetPlaceHolder("Logs will appear here...")
	logBox.Wrapping = fyne.TextWrapWord

	// ============================================================================
	// เลือกแฟ้มเป้าหมาย
	// ============================================================================
	// 🔹 เลือก folder
	projectPath := ""

	selectBtn := widget.NewButton("Select Project Folder", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri == nil {
				return
			}

			projectPath = uri.Path()
			logBox.SetText("📁 Selected: " + projectPath)
		}, w)
	})
	// ============================================================================
	// Generate scrip Icons Btn
	// ============================================================================
	// 🔧 Generate
	genscripiconsBtn := widget.NewButton("Generate scrip Icons", func() {

		if projectPath == "" {
			logBox.SetText("❌ Please select project folder")
			return
		}
		cfg := AppConfig{}

		generateFile("templates/buildicons.tmpl",
			filepath.Join(projectPath, "buildicons.sh"), cfg) //เอา scrip build ออกมาไว้นอกแฟ้ม flatpak

		logBox.SetText("✅ Generated File - - buildicons - -")
	})
	// ============================================================================
	// Generate scrip flatpak Btn
	// ============================================================================
	// 🔧 Generate
	genscripflatpakBtn := widget.NewButton("Generate - Folder and scrip Flatpak - + - File Scrip Build Flatpak", func() {

		if projectPath == "" {
			logBox.SetText("❌ Please select project folder")
			return
		}

		cfg := AppConfig{
			Name:        name.Text,
			AppID:       appID.Text,
			Command:     command.Text,
			Categories:  categories.Text,
			Summary:     summary.Text,
			Description: description.Text,
			License:     "MIT",
			Developer:   developer.Text,
		}

		flatpakPath := projectPath + "/" + "flatpak"
		os.MkdirAll(flatpakPath, 0755)

		generateFile("templates/desktop.tmpl",
			filepath.Join(flatpakPath, cfg.AppID+".desktop"), cfg)

		generateFile("templates/manifest.tmpl",
			filepath.Join(flatpakPath, cfg.AppID+".json"), cfg)

		generateFile("templates/metainfo.tmpl",
			filepath.Join(flatpakPath, cfg.AppID+".metainfo.xml"), cfg)

		generateFile("templates/buildflatpak.tmpl",
			filepath.Join(projectPath, "buildflatpak.sh"), cfg) //เอา scrip build ออกมาไว้นอกแฟ้ม flatpak

		logBox.SetText("✅ Generated File Flatpak\n---------and---------\nFile Scrip Build Flatpak\n")
	})

	// ============================================================================
	// ปุ่ม Build flatpak
	// ============================================================================
	buildflatpakBtn := widget.NewButton("Run Build", func() {

		if projectPath == "" {
			logBox.SetText("❌ select folder first")
			return
		}

		//  run script
		go runScriptbuildflatpak(projectPath, logBox)

		logBox.SetText("🚀 Build started in terminal...")
	})
	// ============================================================================
	// ปุ่ม Build Icons **ใช้ imagemagick
	// ============================================================================
	buildIconsBtn := widget.NewButton("Run Build", func() {

		if projectPath == "" {
			logBox.SetText("❌ select folder first")
			return
		}

		//  run script
		go runScriptbuildIcons(projectPath, logBox)

		logBox.SetText("🚀 Build started in terminal...")
	})

	/*// 🏗️ Build
	buildBtn := widget.NewButton("Build Flatpak", func() {
		logBox.SetText("🚧 Building...\n")

		cmd := exec.Command(
			"flatpak-builder",
			"--force-clean",
			"build-dir",
			filepath.Join("output", appID.Text+".json"),
		)

		go runCommand(cmd, logBox)
	})*/

	/*	// ▶️ Run
		runBtn := widget.NewButton("Run App", func() {
			logBox.SetText("▶️ Running...\n")

			cmd := exec.Command(
				"flatpak",
				"run",
				appID.Text,
			)

			go runCommand(cmd, logBox)
		})*/

	ui := container.NewVBox(
		selectBtn,
		container.NewHBox(genscripiconsBtn, buildIconsBtn),

		name, appID, command,
		categories, summary, description, developer,
		genscripflatpakBtn, buildflatpakBtn,
		//buildBtn,
		//runBtn,
		widget.NewLabel("Logs:"),
		logBox,
	)

	w.SetContent(ui)
	w.Resize(fyne.NewSize(700, 700))
	w.ShowAndRun()
}
