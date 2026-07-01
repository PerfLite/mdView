<template>
  <div
    v-if="show"
    class="fixed inset-0 z-50 flex items-center justify-center"
    style="background: rgba(0, 0, 0, 0.7); backdrop-filter: blur(4px);"
    @click.self="$emit('close')"
  >
    <div
      class="rounded-xl p-6 w-96 max-w-[90vw] relative"
      :style="{ background: 'var(--bg-editor)', border: '1px solid var(--border)', boxShadow: '0 20px 60px rgba(0,0,0,0.5)' }"
    >
      <button
        class="absolute top-3 right-3 w-3 h-3 rounded-full bg-red-500 hover:bg-red-400 transition-colors cursor-pointer border-none"
        style="box-shadow: 0 0 0 1px rgba(0,0,0,0.15);"
        @click="$emit('close')"
        :title="t.close"
      />
      <h2 class="text-lg font-bold text-center mt-0 mb-5 leading-none" :style="{ color: 'var(--accent)' }">{{ t.settings }}</h2>

      <div class="mb-4">
        <label class="text-xs font-semibold uppercase tracking-wider opacity-60 mb-2 block">{{ t.theme }}</label>
        <div class="grid grid-cols-2 gap-2">
          <button
            v-for="theme in themes"
            :key="theme.name"
            class="toolbar-btn text-left flex items-center gap-2 px-3 py-2"
            :class="{ active: currentTheme === theme.name }"
            @click="$emit('setTheme', theme.name)"
          >
            <span
              class="w-3 h-3 rounded-full inline-block border"
              :style="{ background: theme.accent, borderColor: theme.border }"
            />
            {{ theme.name }}
          </button>
        </div>
      </div>

      <div class="mb-4">
        <label class="text-xs font-semibold uppercase tracking-wider opacity-60 mb-2 block">
          {{ t.zoom }}: {{ zoom }}px
        </label>
        <input
          type="range"
          min="10"
          max="24"
          :value="zoom"
          @input="$emit('update:zoom', +$event.target.value)"
          class="w-full accent-current"
          :style="{ accentColor: 'var(--accent)' }"
        />
        <div class="flex justify-between text-xs opacity-40 mt-1">
          <span>10px</span>
          <span>24px</span>
        </div>
      </div>

      <div>
        <label class="text-xs font-semibold uppercase tracking-wider opacity-60 mb-2 block">{{ t.language }}</label>
        <div class="flex gap-2">
          <button
            class="toolbar-btn text-xs flex-1"
            :class="{ active: locale === 'en' }"
            @click="$emit('setLocale', 'en')"
          >English</button>
          <button
            class="toolbar-btn text-xs flex-1"
            :class="{ active: locale === 'ru' }"
            @click="$emit('setLocale', 'ru')"
          >Русский</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  show: Boolean,
  themes: Array,
  currentTheme: String,
  zoom: Number,
  locale: { type: String, default: 'en' },
  t: { type: Object, default: () => ({}) },
})
defineEmits(['close', 'setTheme', 'update:zoom', 'setLocale'])
</script>
