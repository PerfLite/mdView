<template>
  <div ref="containerEl" class="cm-container" />
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, shallowRef } from 'vue'
import { EditorView, keymap, lineNumbers, highlightActiveLine, highlightActiveLineGutter, highlightSpecialChars } from '@codemirror/view'
import { EditorState, Compartment } from '@codemirror/state'
import { markdown } from '@codemirror/lang-markdown'
import { search, searchKeymap, highlightSelectionMatches, openSearchPanel, closeSearchPanel } from '@codemirror/search'
import { defaultKeymap, indentWithTab } from '@codemirror/commands'
import { bracketMatching } from '@codemirror/language'
import { oneDark } from '@codemirror/theme-one-dark'

const props = defineProps({
  modelValue: String,
  zoom: { type: Number, default: 14 },
})
const emit = defineEmits(['update:modelValue', 'scroll', 'save', 'quickSave'])

const containerEl = ref(null)
const view = shallowRef(null)
const themeCompartment = new Compartment()

function createEditor(doc) {
  const updateListener = EditorView.updateListener.of((update) => {
    if (update.docChanged) {
      emit('update:modelValue', update.state.doc.toString())
    }
  })

  const scrollListener = EditorView.domEventHandlers({
    scroll: (e) => {
      if (e.target !== view.value?.scrollDOM) return
      const el = view.value.scrollDOM
      const pct = el.scrollTop / (el.scrollHeight - el.clientHeight || 1)
      emit('scroll', pct)
    },
  })

  const tabKeymap = keymap.of([{
    key: 'Tab',
    run: (cmView) => {
      const { state } = cmView
      const { from, to } = state.selection.main
      if (from === to) {
        const line = state.doc.lineAt(from)
        const indent = state.sliceDoc(line.from, from)
        const spaces = indent.match(/^\s*/)[0]
        const newIndent = spaces + '  '
        cmView.dispatch({
          changes: { from, insert: newIndent },
        })
      } else {
        const lines = []
        for (let i = state.doc.lineAt(from).number; i <= state.doc.lineAt(to).number; i++) {
          lines.push(state.doc.line(i))
        }
        const changes = lines.map((l) => ({
          from: l.from,
          insert: '  ',
        }))
        cmView.dispatch({ changes })
      }
      return true
    },
  }])

  const saveKeymap = keymap.of([{
    key: 'Mod-s',
    run: () => {
      emit('quickSave')
      return true
    },
  }])

  const pasteHandler = EditorView.domEventHandlers({
    paste: (e, cmView) => {
      const items = e.clipboardData?.items
      if (!items) return false
      for (const item of items) {
        if (item.type.startsWith('image/')) {
          e.preventDefault()
          const reader = new FileReader()
          reader.onload = () => {
            const md = `![pasted-image](${reader.result})`
            const pos = cmView.state.selection.main.from
            cmView.dispatch({
              changes: { from: pos, insert: md },
              selection: { anchor: pos + md.length },
            })
          }
          reader.readAsDataURL(item.getAsFile())
          return true
        }
      }
      return false
    },
  })

  const dropHandler = EditorView.domEventHandlers({
    drop: (e, cmView) => {
      const files = e.dataTransfer?.files
      if (!files?.length) return false
      const file = files[0]
      if (file.name.match(/\.(md|markdown|txt)$/i)) {
        e.preventDefault()
        const reader = new FileReader()
        reader.onload = () => {
          cmView.dispatch({
            changes: { from: 0, to: cmView.state.doc.length, insert: reader.result },
          })
        }
        reader.readAsText(file)
        return true
      } else if (file.type.startsWith('image/')) {
        e.preventDefault()
        const reader = new FileReader()
        reader.onload = () => {
          const md = `\n![${file.name}](${reader.result})\n`
          const pos = cmView.state.selection.main.from
          cmView.dispatch({
            changes: { from: pos, insert: md },
            selection: { anchor: pos + md.length },
          })
        }
        reader.readAsDataURL(file)
        return true
      }
      return false
    },
  })

  const state = EditorState.create({
    doc: doc || '',
    extensions: [
      EditorView.lineWrapping,
      lineNumbers(),
      highlightActiveLine(),
      highlightActiveLineGutter(),
      highlightSpecialChars(),
      bracketMatching(),
      markdown(),
      oneDark,
      search({ top: true }),
      highlightSelectionMatches(),
      saveKeymap,
      keymap.of([...defaultKeymap, ...searchKeymap, indentWithTab]),
      tabKeymap,
      scrollListener,
      pasteHandler,
      dropHandler,
      updateListener,
      themeCompartment.of(EditorView.theme({
        '&': { height: '100%', fontSize: props.zoom + 'px' },
        '.cm-scroller': { overflow: 'auto', fontFamily: 'var(--font-mono)' },
        '.cm-gutters': { background: 'var(--bg-editor)', borderRight: '1px solid var(--border)' },
        '.cm-activeLineGutter': { background: 'rgba(255,255,255,0.06)', color: 'var(--accent)' },
        '.cm-activeLine': { background: 'rgba(255,255,255,0.04)' },
      })),
    ],
  })

  return new EditorView({
    state,
    parent: containerEl.value,
  })
}

onMounted(() => {
  view.value = createEditor(props.modelValue)
})

onUnmounted(() => {
  view.value?.destroy()
})

watch(() => props.modelValue, (newVal) => {
  const v = view.value
  if (!v) return
  const current = v.state.doc.toString()
  if (newVal !== current) {
    v.dispatch({
      changes: { from: 0, to: v.state.doc.length, insert: newVal },
    })
  }
})

watch(() => props.zoom, (newZoom) => {
  const v = view.value
  if (!v) return
  v.dispatch({
    effects: themeCompartment.reconfigure(EditorView.theme({
      '&': { height: '100%', fontSize: newZoom + 'px' },
      '.cm-scroller': { overflow: 'auto', fontFamily: 'var(--font-mono)' },
      '.cm-gutters': { background: 'var(--bg-editor)', borderRight: '1px solid var(--border)' },
      '.cm-activeLineGutter': { background: 'rgba(255,255,255,0.06)', color: 'var(--accent)' },
      '.cm-activeLine': { background: 'rgba(255,255,255,0.04)' },
    })),
  })
})

function scrollTo(pct) {
  const v = view.value
  if (!v) return
  const el = v.scrollDOM
  const max = el.scrollHeight - el.clientHeight
  el.scrollTop = Math.min(1, Math.max(0, pct)) * max
}

function scrollToLine(lineNum) {
  const v = view.value
  if (!v) return
  const line = v.state.doc.line(Math.min(lineNum, v.state.doc.lines))
  v.dispatch({ selection: { anchor: line.from }, scrollIntoView: true })
}

function toggleSearch() {
  const v = view.value
  if (!v) return
  const panel = v.dom.querySelector('.cm-panels')
  if (panel) {
    closeSearchPanel(v)
  } else {
    openSearchPanel(v)
  }
}

defineExpose({ scrollTo, scrollToLine, toggleSearch })
</script>
