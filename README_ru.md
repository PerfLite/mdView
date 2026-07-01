# mdView

[![License: GPL-3.0](https://img.shields.io/badge/License-GPL--3.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg?logo=go)](https://go.dev)
[![Vue](https://img.shields.io/badge/Vue-3-4FC08D.svg?logo=vue.js)](https://vuejs.org)
[![Wails](https://img.shields.io/badge/Wails-v2-DD5B47.svg)](https://wails.io)
[![Platform](https://img.shields.io/badge/Platform-Linux-lightgrey.svg)](https://github.com/PerfLite/mdView)
[![Version](https://img.shields.io/badge/Version-1.0.0-green.svg)](https://github.com/PerfLite/mdView/releases)

Лёгкий редактор Markdown с предпросмотром, создан для разработчиков.

## Возможности

- Разделённый вид редактора с регулируемыми панелями
- Предпросмотр Markdown в реальном времени (задержка 300мс)
- Подсветка синтаксиса (Goldmark + Chroma)
- 4 встроенные темы: Ocean, Forest, Cyberpunk, Midnight
- Управление масштабом (Ctrl+Wheel или Настройки)
- Поиск/Замена (Ctrl+F)
- Номера строк и подсветка активной строки
- Перенос длинных строк
- Открытие/Сохранение файлов через нативные диалоги
- Быстрое сохранение (Ctrl+S)
- Экспорт в HTML и PDF
- Автосохранение (каждые 5 секунд)
- Вставка изображений и drag-and-drop
- Drag-and-drop файлов .md
- Режимы отображения: Только редактор, Разделённый, Только превью
- Поддержка английского и русского языков
- Настройка через панель «Настройки»

## Технологии

| Компонент | Технология |
|-----------|------------|
| Фреймворк | Wails v2 + Go |
| Интерфейс | Vue 3 + Tailwind CSS |
| Редактор | CodeMirror 6 |
| Markdown | Goldmark (Go) |
| Подсветка | Chroma |

## Установка

### deb-пакет (Debian/Ubuntu)

Сначала установите зависимости:

```bash
sudo apt install libgtk-3-0 libwebkit2gtk-4.0-37
```

Затем установите пакет:

```bash
sudo dpkg -i mdview_1.0.0_amd64.deb
```

### pacman (Arch Linux)

Установите зависимости:

```bash
sudo pacman -S gtk3 webkit2gtk-4.1
```

Соберите из исходников:

```bash
git clone https://github.com/PerfLite/mdView.git
cd mdView
wails build
sudo cp build/bin/mdView /usr/bin/mdview
sudo cp mdview.desktop /usr/share/applications/
sudo cp icon_256.png /usr/share/icons/hicolor/256x256/apps/mdview.png
```

### AppImage (любой Linux)

Установка не требуется. Самодостаточный.

```bash
chmod +x mdView-x86_64.AppImage
./mdView-x86_64.AppImage
```

### Сборка из исходников

Требования: Go 1.21+, Node.js 18+, Wails v2

```bash
git clone https://github.com/PerfLite/mdView.git
cd mdView
wails build
```

Бинарник будет в `build/bin/mdView`.

## Назначение редактором по умолчанию

После установки:

```bash
xdg-mime default mdview.desktop text/markdown
xdg-mime default mdview.desktop text/x-markdown
```

## Горячие клавиши

| Комбинация | Действие |
|------------|----------|
| Ctrl+S | Быстрое сохранение |
| Ctrl+O | Открыть файл |
| Ctrl+F | Открыть/закрыть поиск |
| F11 | Полноэкранный режим |
| Ctrl+Wheel | Масштаб |

## Конфигурация

Файлы конфигурации хранятся в `~/.config/mdView/`:

- `config.json` — настройки темы и языка
- `autosave.md` — автосохранённое содержимое
- `last_open.txt` — путь последнего открытого файла

## Лицензия

[GPL-3.0](LICENSE)
