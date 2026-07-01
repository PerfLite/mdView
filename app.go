package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

type Theme struct {
	Name       string `json:"name"`
	Background string `json:"background"`
	Text       string `json:"text"`
	Accent     string `json:"accent"`
	Editor     string `json:"editor"`
	Preview    string `json:"preview"`
	Border     string `json:"border"`
}

type OpenResult struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

var themes = []Theme{
	{Name: "Ocean", Background: "#0f172a", Text: "#e2e8f0", Accent: "#38bdf8", Editor: "#1e293b", Preview: "#0f172a", Border: "#334155"},
	{Name: "Forest", Background: "#0c1a0c", Text: "#d1fae5", Accent: "#4ade80", Editor: "#14301a", Preview: "#0c1a0c", Border: "#166534"},
	{Name: "Cyberpunk", Background: "#0a0a0f", Text: "#f0f0f0", Accent: "#f0abfc", Editor: "#1a0a2e", Preview: "#0a0a0f", Border: "#7c3aed"},
	{Name: "Midnight", Background: "#18181b", Text: "#fafafa", Accent: "#818cf8", Editor: "#27272a", Preview: "#18181b", Border: "#3f3f46"},
}

type App struct {
	ctx         context.Context
	markdown    goldmark.Markdown
	configDir   string
	isDirty     bool
	locale      string
}

func NewApp() *App {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(highlighting.WithStyle("dracula")),
		),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)
	return &App{markdown: md}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.locale = "en"
	home, err := os.UserHomeDir()
	if err == nil {
		a.configDir = home + "/.config/mdView"
		os.MkdirAll(a.configDir, 0755)
	}
}

func (a *App) GetCLIArgs() []string {
	return os.Args[1:]
}

func (a *App) SetDirty(dirty bool) {
	a.isDirty = dirty
}

func (a *App) SetLocale(locale string) {
	a.locale = locale
}

func (a *App) ConfirmClose() bool {
	if !a.isDirty {
		return true
	}
	title := "Unsaved Changes"
	message := "You have unsaved changes. Do you want to save before closing?"
	cancelBtn := "Cancel"
	if a.locale == "ru" {
		title = "Несохранённые изменения"
		message = "У вас есть несохранённые изменения. Сохранить перед закрытием?"
		cancelBtn = "Отмена"
	}
	dontSave := "Don't Save"
	save := "Save"
	if a.locale == "ru" {
		dontSave = "Не сохранять"
		save = "Сохранить"
	}
	result, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         title,
		Message:       message,
		Buttons:       []string{save, dontSave, cancelBtn},
		DefaultButton: save,
		CancelButton:  cancelBtn,
	})
	if err != nil {
		return true
	}
	switch result {
	case save:
		return true
	case dontSave:
		a.isDirty = false
		return true
	default:
		return false
	}
}

func (a *App) RenderMarkdown(text string) string {
	var buf bytes.Buffer
	if err := a.markdown.Convert([]byte(text), &buf); err != nil {
		return "<p style='color:red'>Error rendering markdown</p>"
	}
	return buf.String()
}

func (a *App) SaveFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func (a *App) SaveFileAs(content string) string {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title: "Save Markdown File",
		Filters: []runtime.FileFilter{
			{DisplayName: "Markdown Files (*.md)", Pattern: "*.md"},
			{DisplayName: "All Files (*.*)", Pattern: "*.*"},
		},
		DefaultFilename: "untitled.md",
	})
	if err != nil || path == "" {
		return ""
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return ""
	}
	return path
}

func (a *App) OpenFile() string {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Open Markdown File",
		Filters: []runtime.FileFilter{
			{DisplayName: "Markdown Files (*.md)", Pattern: "*.md"},
			{DisplayName: "All Files (*.*)", Pattern: "*.*"},
		},
	})
	if err != nil || path == "" {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	a.SetLastOpenPath(path)
	result, _ := json.Marshal(OpenResult{Path: path, Content: string(data)})
	return string(result)
}

func (a *App) ReadFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func (a *App) SetLastOpenPath(path string) {
	os.WriteFile(a.configDir+"/last_open.txt", []byte(path), 0644)
}

func (a *App) GetThemes() []Theme {
	return themes
}

func (a *App) GetConfig() map[string]string {
	path := a.configDir + "/config.json"
	data, err := os.ReadFile(path)
	if err != nil {
		return map[string]string{"theme": "Ocean", "locale": "en"}
	}
	var cfg map[string]string
	json.Unmarshal(data, &cfg)
	if cfg["locale"] == "" {
		cfg["locale"] = "en"
	}
	return cfg
}

func (a *App) SaveConfig(theme string) error {
	existing := a.GetConfig()
	existing["theme"] = theme
	data, _ := json.Marshal(existing)
	return os.WriteFile(a.configDir+"/config.json", data, 0644)
}

func (a *App) SaveLocale(locale string) error {
	existing := a.GetConfig()
	existing["locale"] = locale
	data, _ := json.Marshal(existing)
	return os.WriteFile(a.configDir+"/config.json", data, 0644)
}

func (a *App) AutoSave(content string) error {
	return os.WriteFile(a.configDir+"/autosave.md", []byte(content), 0644)
}

func (a *App) LoadAutoSave() string {
	data, err := os.ReadFile(a.configDir + "/autosave.md")
	if err != nil {
		return ""
	}
	return string(data)
}

func (a *App) ClearAutoSave() {
	os.Remove(a.configDir + "/autosave.md")
}

func (a *App) OpenURL(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}

var htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>mdView Export</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500&family=Nunito:wght@400;600;700&display=swap" rel="stylesheet">
<style>
body{font-family:'Nunito',sans-serif;max-width:800px;margin:2rem auto;padding:0 2rem;color:#e2e8f0;background:#0f172a;line-height:1.8}
h1,h2,h3{color:#38bdf8;margin-top:1.5rem}
h1{border-bottom:1px solid #334155;padding-bottom:.5rem}
code{font-family:'JetBrains Mono',monospace;background:rgba(255,255,255,.06);padding:.15rem .4rem;border-radius:4px;font-size:.9em}
pre{background:#1a1b26;border:1px solid #334155;border-radius:8px;padding:1rem;overflow-x:auto}
pre code{background:none;padding:0;font-size:.85em;line-height:1.6}
blockquote{border-left:3px solid #38bdf8;padding-left:1rem;color:#94a3b8;font-style:italic;margin:1rem 0}
table{border-collapse:collapse;width:100%%;margin:1rem 0}
th,td{border:1px solid #334155;padding:.6rem .8rem;text-align:left}
th{background:rgba(255,255,255,.05);font-weight:600}
a{color:#38bdf8}
img{max-width:100%%;border-radius:8px}
hr{border:none;border-top:1px solid #334155;margin:1.5rem 0}
</style>
</head>
<body>
%s
</body>
</html>`

func (a *App) ExportHTML(html string, path string) error {
	fullHTML := fmt.Sprintf(htmlTemplate, html)
	return os.WriteFile(path, []byte(fullHTML), 0644)
}

func (a *App) ExportPDF(html string, path string) error {
	fullHTML := fmt.Sprintf(htmlTemplate, html)
	tmpHTML := a.configDir + "/export_tmp.html"
	if err := os.WriteFile(tmpHTML, []byte(fullHTML), 0644); err != nil {
		return err
	}
	defer os.Remove(tmpHTML)

	binaries := []string{"chromium", "chromium-browser", "google-chrome", "google-chrome-stable"}
	for _, bin := range binaries {
		cmd := exec.Command(bin, "--headless", "--disable-gpu", "--no-sandbox",
			"--print-to-pdf="+path, tmpHTML)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}
	return fmt.Errorf("no chromium-based browser found for PDF export")
}
