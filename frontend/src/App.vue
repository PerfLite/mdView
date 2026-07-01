<template>
  <div
    class="flex flex-col h-dvh"
    :style="{ ...themeVars, fontSize: zoom + 'px' }"
    @dragover.prevent
    @drop.prevent="onDropFiles"
  >
    <div style="position: relative; z-index: 10;">
      <Toolbar
        :fileName="fileName"
        :scrollLocked="scrollLocked"
        :viewMode="viewMode"
        :t="t"
        @open="openFile"
        @save="saveFile"
        @toggleLock="scrollLocked = !scrollLocked"
        @toggleSettings="showSettings = !showSettings"
        @toggleAbout="showAbout = !showAbout"
        @exportHTML="exportHTML"
        @exportPDF="exportPDF"
        @setViewMode="viewMode = $event"
      />
    </div>

    <Settings
      :show="showSettings"
      :themes="themes"
      :currentTheme="currentTheme"
      :zoom="zoom"
      :locale="locale"
      :t="t"
      @close="showSettings = false"
      @setTheme="setTheme"
      @update:zoom="v => zoom = v"
      @setLocale="setLocale"
    />

    <About
      :show="showAbout"
      :t="t"
      @close="showAbout = false"
    />

    <div class="flex flex-1 min-h-0">
      <div v-if="viewMode !== 'preview'" :style="viewMode === 'editor' ? {} : { width: editorWidth + 'px' }" class="flex flex-col min-w-0">
        <Editor
          ref="editorRef"
          v-model="markdown"
          :zoom="zoom"
          @scroll="onEditorScroll"
          @save="saveFile"
          @quickSave="quickSave"
        />
      </div>

      <div
        v-if="viewMode === 'split'"
        class="divider"
        @mousedown="startResize"
      />

      <div v-if="viewMode !== 'editor'" class="flex-1 min-w-0 glass">
        <Preview
          ref="previewRef"
          :html="renderedHtml"
          @scroll="onPreviewScroll"
        />
      </div>
    </div>

    <Transition name="toast">
      <div
        v-if="toast"
        class="fixed bottom-6 left-1/2 -translate-x-1/2 px-4 py-2 rounded-lg text-sm font-semibold shadow-lg"
        :style="{ background: 'var(--accent)', color: 'var(--bg)' }"
      >
        {{ toast }}
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useDebounceFn } from './useDebounce'
import { useI18n } from './i18n'
import Editor from './components/Editor.vue'
import Preview from './components/Preview.vue'
import Toolbar from './components/Toolbar.vue'
import Settings from './components/Settings.vue'
import About from './components/About.vue'
import {
  RenderMarkdown,
  OpenFile,
  SaveFile,
  ReadFile,
  GetThemes,
  SetLastOpenPath,
  GetConfig,
  SaveConfig,
  SaveLocale,
  AutoSave,
  LoadAutoSave,
  ClearAutoSave,
  ExportHTML,
  ExportPDF,
  SaveFileAs,
  SetLocale as GoSetLocale,
  GetCLIArgs,
} from '../wailsjs/go/main/App'

const defaultText = `# mdView — Markdown Previewer

> A lightweight, beautiful desktop Markdown editor built with **Wails v2**, **Vue 3**, and **Go**.

## Getting Started

This is a live-preview Markdown editor. Type on the left, see the rendered result on the right. The preview updates in real-time as you type.

## Text Formatting

You can write **bold**, *italic*, ~~strikethrough~~, and \`inline code\`.

### Lists

#### Unordered

- First item
- Second item
- Third item
  - Nested item
  - Another nested item

#### Ordered

1. Step one
2. Step two
3. Step three

## Code Blocks

\`\`\`go
package main

import "fmt"

func main() {
    fmt.Println("Hello from mdView!")
}
\`\`\`

\`\`\`javascript
const greet = (name) => {
    return \`Hello, \${name}!\`;
};

console.log(greet("World"));
\`\`\`

\`\`\`python
def fibonacci(n):
    if n <= 1:
        return n
    return fibonacci(n - 1) + fibonacci(n - 2)

print(fibonacci(10))
\`\`\`

## Tables

| Feature         | Status | Notes                        |
|-----------------|--------|------------------------------|
| Split View      | Done   | Resizable divider            |
| Live Preview    | Done   | 300ms debounce               |
| Themes          | Done   | 4 built-in themes            |
| Zoom            | Done   | Ctrl+Wheel or settings       |
| File Open/Save  | Done   | Native dialogs               |
| Syntax Highlight| Done   | Goldmark + Chroma            |
| Line Numbers    | Done   | Editor gutter                |
| Search/Replace  | Done   | Ctrl+F                       |
| Export HTML/PDF | Done   | Toolbar buttons              |
| Auto-save       | Done   | Periodic temp save           |
| Image Paste     | Done   | Paste or drag images         |
| Drag & Drop     | Done   | Drop .md files               |

## Blockquotes

> "The best way to predict the future is to invent it."
> — Alan Kay

## Links & Images

Check out [Wails](https://wails.io) for building desktop apps with Go and web technologies.

## Horizontal Rule

---

## Task List

- [x] Set up project with Wails v2
- [x] Implement Markdown rendering with goldmark
- [x] Build split-view editor UI
- [x] Add theming system
- [x] Add zoom, scroll lock, line numbers
- [x] Add search/replace
- [x] Add export HTML/PDF
- [x] Add auto-save
- [x] Add image paste/drag

---

*End of test document. Use Ctrl+S to save, Ctrl+O to open, Ctrl+F to search, F11 for fullscreen.*`

const markdown = ref(defaultText)
const renderedHtml = ref('')
const themes = ref([])
const currentTheme = ref('Ocean')
const fileName = ref('')
const filePath = ref('')
const editorWidth = ref(500)
const zoom = ref(14)
const showSettings = ref(false)
const showAbout = ref(false)
const scrollLocked = ref(false)
const viewMode = ref('split')
const isDirty = ref(false)
const toast = ref('')
let toastTimer = null
const locale = ref('en')
const { t } = useI18n(locale)

const editorRef = ref(null)
const previewRef = ref(null)
let syncing = false
let autoSaveTimer = null
let savedContent = defaultText

function onEditorScroll(pct) {
  if (!scrollLocked.value || syncing) return
  syncing = true
  requestAnimationFrame(() => {
    previewRef.value?.scrollTo(pct)
    requestAnimationFrame(() => { syncing = false })
  })
}

function onPreviewScroll(pct) {
  if (!scrollLocked.value || syncing) return
  syncing = true
  requestAnimationFrame(() => {
    editorRef.value?.scrollTo(pct)
    requestAnimationFrame(() => { syncing = false })
  })
}

function showToast(msg) {
  toast.value = msg
  clearTimeout(toastTimer)
  toastTimer = setTimeout(() => toast.value = '', 2000)
}

const themeVars = computed(() => {
  const t = themes.value.find(t => t.name === currentTheme.value)
  if (!t) return {}
  return {
    '--bg': t.background,
    '--bg-editor': t.editor,
    '--bg-preview': t.preview,
    '--text': t.text,
    '--accent': t.accent,
    '--border': t.border,
  }
})

const render = useDebounceFn(async () => {
  renderedHtml.value = await RenderMarkdown(markdown.value)
}, 300)

watch(markdown, (val) => {
  render()
  scheduleAutoSave()
  isDirty.value = val !== savedContent
}, { immediate: true })

async function openFile() {
  try {
    const result = await OpenFile()
    if (!result) return
    const parsed = JSON.parse(result)
    if (parsed.path && parsed.content) {
      markdown.value = parsed.content
      savedContent = parsed.content
      fileName.value = parsed.path.split('/').pop()
      filePath.value = parsed.path
      isDirty.value = false
      showToast('Opened: ' + fileName.value)
    }
  } catch (e) {
    console.error('Open failed:', e)
  }
}

async function saveFile() {
  try {
    if (filePath.value) {
      await SaveFile(filePath.value, markdown.value)
      savedContent = markdown.value
      isDirty.value = false
      showToast('Saved: ' + fileName.value)
      ClearAutoSave()
    } else {
      const result = await SaveFileAs(markdown.value)
      if (result) {
        filePath.value = result
        fileName.value = result.split('/').pop()
        savedContent = markdown.value
        isDirty.value = false
        showToast('Saved: ' + fileName.value)
        ClearAutoSave()
      }
    }
  } catch (e) {
    console.error('Save failed:', e)
  }
}

function setTheme(name) {
  currentTheme.value = name
  SaveConfig(name)
}

function setLocale(lang) {
  locale.value = lang
  SaveLocale(lang)
  GoSetLocale(lang)
}

async function exportHTML() {
  try {
    const html = await RenderMarkdown(markdown.value)
    const name = (fileName.value || 'export').replace(/\.[^.]+$/, '') + '.html'
    await ExportHTML(html, name)
    showToast('Exported: ' + name)
  } catch (e) {
    showToast('Export failed')
  }
}

async function exportPDF() {
  try {
    const html = await RenderMarkdown(markdown.value)
    const name = (fileName.value || 'export').replace(/\.[^.]+$/, '') + '.pdf'
    await ExportPDF(html, name)
    showToast('Exported: ' + name)
  } catch (e) {
    showToast('PDF export failed — no chromium found')
  }
}

function onDropFiles(e) {
  const files = e.dataTransfer?.files
  if (!files?.length) return
  const file = files[0]
  if (file.name.match(/\.(md|markdown|txt)$/i)) {
    const reader = new FileReader()
    reader.onload = async () => {
      markdown.value = reader.result
      fileName.value = file.name
      filePath.value = ''
    }
    reader.readAsText(file)
  }
}

let resizing = false
function startResize(e) {
  resizing = true
  const startX = e.clientX
  const startWidth = editorWidth.value
  function onMove(ev) {
    if (!resizing) return
    editorWidth.value = Math.max(250, startWidth + (ev.clientX - startX))
  }
  function onUp() {
    resizing = false
    document.removeEventListener('mousemove', onMove)
    document.removeEventListener('mouseup', onUp)
  }
  document.addEventListener('mousemove', onMove)
  document.addEventListener('mouseup', onUp)
}

function onWheel(e) {
  if (e.ctrlKey) {
    e.preventDefault()
    zoom.value = Math.min(24, Math.max(10, zoom.value + (e.deltaY > 0 ? -1 : 1)))
  }
}

function toggleFullscreen() {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
  } else {
    document.exitFullscreen()
  }
}

function isKey(e, code, ...keys) {
  return e.code === code || keys.includes(e.key)
}

function onKeyDown(e) {
  if (e.ctrlKey && isKey(e, 'KeyS')) {
    e.preventDefault()
    e.stopPropagation()
    quickSave()
  }
  if (e.ctrlKey && isKey(e, 'KeyO')) {
    e.preventDefault()
    openFile()
  }
  if (e.ctrlKey && isKey(e, 'KeyF')) {
    e.preventDefault()
    editorRef.value?.toggleSearch()
  }
  if (isKey(e, 'F11')) {
    e.preventDefault()
    toggleFullscreen()
  }
}

function scheduleAutoSave() {
  clearTimeout(autoSaveTimer)
  autoSaveTimer = setTimeout(() => {
    AutoSave(markdown.value)
  }, 5000)
}

async function quickSave() {
  await AutoSave(markdown.value)
  savedContent = markdown.value
  isDirty.value = false
  showToast(t.value.saved || 'Saved')
}

onMounted(async () => {
  themes.value = await GetThemes()
  const cfg = await GetConfig()
  if (cfg?.theme) currentTheme.value = cfg.theme
  if (cfg?.locale) locale.value = cfg.locale

  const autoSaved = await LoadAutoSave()

  const args = await GetCLIArgs()
  const fileArg = args.find(a => a.match(/\.(md|markdown|txt)$/i))
  if (fileArg) {
    try {
      const content = await ReadFile(fileArg)
      if (content) {
        markdown.value = content
        savedContent = content
        fileName.value = fileArg.split('/').pop()
        filePath.value = fileArg
      }
    } catch (e) {
      console.error('Failed to open file from args:', e)
    }
  } else if (autoSaved) {
    markdown.value = autoSaved
    savedContent = autoSaved
  }

  render()
  window.addEventListener('wheel', onWheel, { passive: false })
  window.addEventListener('keydown', onKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('wheel', onWheel)
  window.removeEventListener('keydown', onKeyDown)
  clearTimeout(autoSaveTimer)
})
</script>
