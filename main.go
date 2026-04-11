// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package main

import (
	"embed"
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

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
	Date        string
	TimeEntry   string
	Version     string
	DesUpdate1  string
	DesUpdate2  string
	Owner       string
	NameRepo    string
	NamePix1    string
	NamePix2    string
	NamePix3    string
	NamePix4    string
	NamePix5    string
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
	name.SetPlaceHolder("App Name - ชื่อโปรแกรม -แอพ จะแสดงในหน้าสโตร์")

	appID := widget.NewEntry()
	appID.SetPlaceHolder("com.example.app - แอพไอดี")

	command := widget.NewEntry()
	command.SetPlaceHolder("binary name - ชื่อโปแกรมตอน Build")

	categories := widget.NewEntry()
	categories.SetPlaceHolder("Utility; - ประเภทโปรแกรม")

	catmenu := widget.NewMultiLineEntry()
	catmenu.SetText(`ประเภทโปรแกรม
	Utility; = ยูทิลิตี้ (ทั่วไป)
	Development; = การพัฒนา
	Game; = เกม
	Graphics; = กราฟิก
	Network; = เครือข่าย
	Office; = สำนักงาน
	Audio; = เสียง
	Video; = วิดีโอ
	System; = ระบบ`)
	catmenu.SetMinRowsVisible(11)

	summary := widget.NewEntry()
	summary.SetPlaceHolder("Short summary - คุณบัติของแอพ")

	description := widget.NewMultiLineEntry()
	description.SetPlaceHolder("Description - รายละเอียดของแอพ")

	developer := widget.NewEntry()
	developer.SetPlaceHolder("Your name - จะแสดงหน้าสโตร์")

	date := widget.NewEntry()
	date.SetPlaceHolder("📅 วันที่ - YYYY-MM-DD")

	timeEntry := widget.NewEntry()
	timeEntry.SetPlaceHolder("⏰ เวลา - HH:MM")

	version := widget.NewEntry()
	version.SetPlaceHolder("ใ่ส่เวอร์ชัน เช่น 1.0.0")

	desUpdate1 := widget.NewEntry()
	desUpdate1.SetPlaceHolder("-สิ่งที่อัพเดท 1")

	desUpdate2 := widget.NewEntry()
	desUpdate2.SetPlaceHolder("-สิ่งที่อัพเดท 2")

	owner := widget.NewEntry()
	owner.SetPlaceHolder("ชื่อเจ้าของ Github [Owner]")

	nameRepo := widget.NewEntry()
	nameRepo.SetPlaceHolder("ชื่อ Repository")

	namePix1 := widget.NewEntry()
	namePix1.SetPlaceHolder("1.ชื่อ รูป ไม่ต้องเติม นามสกุล (เอารูปวางไว้ข้างไฟล์ main โปรเจค)")

	namePix2 := widget.NewEntry()
	namePix2.SetPlaceHolder("2.ชื่อ รูป ไม่ต้องเติม นามสกุล (เอารูปวางไว้ข้างไฟล์ main โปรเจค)")

	namePix3 := widget.NewEntry()
	namePix3.SetPlaceHolder("3.ชื่อ รูป ไม่ต้องเติม นามสกุล (เอารูปวางไว้ข้างไฟล์ main โปรเจค)")

	namePix4 := widget.NewEntry()
	namePix4.SetPlaceHolder("4.ชื่อ รูป ไม่ต้องเติม นามสกุล (เอารูปวางไว้ข้างไฟล์ main โปรเจค)")

	namePix5 := widget.NewEntry()
	namePix5.SetPlaceHolder("5.ชื่อ รูป ไม่ต้องเติม นามสกุล (เอารูปวางไว้ข้างไฟล์ main โปรเจค)")

	// log box
	logBox := widget.NewMultiLineEntry()
	logBox.SetPlaceHolder("Logs will appear here...")
	logBox.Wrapping = fyne.TextWrapWord

	// ============================================================================
	// เลือกแฟ้มเป้าหมาย
	// ============================================================================
	// 🔹 เลือก folder
	projectPath := ""

	selectBtn := widget.NewButton("1 - Select Project Folder", func() {
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
	genscripiconsBtn := widget.NewButton("2 - Generate scrip Icons", func() {

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
	genscripflatpakBtn := widget.NewButton("5 - Generate - Folder and scrip Flatpak - + - File Scrip Build Flatpak", func() {

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
			License:     "GPL-3.0-or-later",
			Developer:   developer.Text,
			Date:        date.Text,
			TimeEntry:   timeEntry.Text,
			Version:     version.Text,
			DesUpdate1:  desUpdate1.Text,
			DesUpdate2:  desUpdate2.Text,
			Owner:       owner.Text,
			NameRepo:    nameRepo.Text,
			NamePix1:    namePix1.Text,
			NamePix2:    namePix2.Text,
			NamePix3:    namePix3.Text,
			NamePix4:    namePix4.Text,
			NamePix5:    namePix5.Text,
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
	buildflatpakBtn := widget.NewButton("7 - Run Build Flatpak", func() {

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
	buildIconsBtn := widget.NewButton("3 - Run Build Icons", func() {

		if projectPath == "" {
			logBox.SetText("❌ select folder first")
			return
		}

		//  run script
		go runScriptbuildIcons(projectPath, logBox)

		logBox.SetText("🚀 Build started in terminal...")
	})

	// ============================================================================
	// ปุ่มเพิ่มวัน เวลา
	// ============================================================================
	nowBtn := widget.NewButton("4 - กด เพื่อ ใส่เวลาปัจจุบัน", func() {
		now := time.Now()

		date.SetText(now.Format("2006-01-02"))
		timeEntry.SetText(now.Format("15:04"))
	})
	// ============================================================================
	// จัดหน้ามัน
	// ============================================================================

	ui := container.NewGridWithColumns(2,

		container.NewVBox(
			selectBtn,
			container.NewGridWithColumns(2, genscripiconsBtn, buildIconsBtn),

			name, appID, command,
			categories, catmenu, developer,
			version,
			container.NewGridWithColumns(3, date, timeEntry, nowBtn),
		),
		container.NewVBox(summary, description,
			container.NewGridWithColumns(2, desUpdate1, desUpdate2),
			container.NewGridWithColumns(2, owner, nameRepo),
			namePix1,
			namePix2,
			namePix3,
			namePix4,
			namePix5,
			genscripflatpakBtn,
			widget.NewLabel("6 - ตรวจเช็คไฟล์ XML ก่อน"),
			buildflatpakBtn,
			//widget.NewLabel("Logs:"),
			logBox,
		),
	)

	w.SetContent(ui)
	w.Resize(fyne.NewSize(1000, 600))
	w.SetFixedSize(true)
	w.ShowAndRun()
}
