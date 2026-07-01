# mdView

[![License: GPL-3.0](https://img.shields.io/badge/License-GPL--3.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg?logo=go)](https://go.dev)
[![Vue](https://img.shields.io/badge/Vue-3-4FC08D.svg?logo=vue.js)](https://vuejs.org)
[![Wails](https://img.shields.io/badge/Wails-v2-DD5B47.svg)](https://wails.io)
[![Platform](https://img.shields.io/badge/Platform-Linux-lightgrey.svg)](https://github.com/PerfLite/mdView)
[![Version](https://img.shields.io/badge/Version-1.0.0-green.svg)](https://github.com/PerfLite/mdView/releases)

[![RU](https://img.shields.io/badge/📖-README%20на%20русском-blue)](https://github.com/PerfLite/mdView/blob/main/README_ru.md)

A lightweight Markdown editor with live preview, built for developers.

## Features

- Split-view editor with resizable panels
- Live Markdown preview (300ms debounce)
- Syntax highlighting (Goldmark + Chroma)
- 4 built-in themes: Ocean, Forest, Cyberpunk, Midnight
- Zoom control (Ctrl+Wheel or Settings)
- Search/Replace (Ctrl+F)
- Line numbers and active line highlighting
- Word wrap
- File Open/Save with native dialogs
- Quick Save (Ctrl+S)
- Export HTML and PDF
- Auto-save (every 5 seconds)
- Image paste and drag-and-drop
- Markdown file drag-and-drop
- View modes: Editor only, Split, Preview only
- English and Russian language support
- Configurable via Settings panel

## Tech Stack

| Component | Technology |
|-----------|------------|
| Framework | Wails v2 + Go |
| UI | Vue 3 + Tailwind CSS |
| Editor | CodeMirror 6 |
| Markdown | Goldmark (Go) |
| Syntax Highlight | Chroma |

## Installation

### deb package (Debian/Ubuntu)

Install dependencies first:

```bash
sudo apt install libgtk-3-0 libwebkit2gtk-4.0-37
```

Then install the package:

```bash
sudo dpkg -i mdview_1.0.0_amd64.deb
```

### pacman (Arch Linux)

Install dependencies:

```bash
sudo pacman -S gtk3 webkit2gtk-4.1
```

Build from source:

```bash
git clone https://github.com/PerfLite/mdView.git
cd mdView
wails build
sudo cp build/bin/mdView /usr/bin/mdview
sudo cp mdview.desktop /usr/share/applications/
sudo cp icon_256.png /usr/share/icons/hicolor/256x256/apps/mdview.png
```

### AppImage (any Linux)

No installation required. Self-contained.

```bash
chmod +x mdView-x86_64.AppImage
./mdView-x86_64.AppImage
```

### Build from source

Requirements: Go 1.21+, Node.js 18+, Wails v2

```bash
git clone https://github.com/PerfLite/mdView.git
cd mdView
wails build
```

Binary will be at `build/bin/mdView`.

## Set as Default Editor

After installing:

```bash
xdg-mime default mdview.desktop text/markdown
xdg-mime default mdview.desktop text/x-markdown
```

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| Ctrl+S | Quick save |
| Ctrl+O | Open file |
| Ctrl+F | Toggle search |
| F11 | Toggle fullscreen |
| Ctrl+Wheel | Zoom in/out |

## Configuration

Config files are stored at `~/.config/mdView/`:

- `config.json` — theme and language settings
- `autosave.md` — auto-saved content
- `last_open.txt` — last opened file path

## License

[GPL-3.0](LICENSE)
