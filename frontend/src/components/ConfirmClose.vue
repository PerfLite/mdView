<template>
  <div
    v-if="show"
    class="fixed inset-0 z-50 flex items-center justify-center"
    style="background: rgba(0, 0, 0, 0.7); backdrop-filter: blur(4px);"
  >
    <div
      class="rounded-xl p-6 w-80 max-w-[90vw]"
      :style="{ background: 'var(--bg-editor)', border: '1px solid var(--border)', boxShadow: '0 20px 60px rgba(0,0,0,0.5)' }"
    >
      <h3 class="text-base font-semibold mb-4" :style="{ color: 'var(--text)' }">
        {{ t.unsavedTitle }}
      </h3>
      <p class="text-sm mb-5 opacity-70">{{ t.unsavedMessage }}</p>
      <div class="flex gap-2 justify-end">
        <button
          class="toolbar-btn text-xs"
          @click="$emit('cancel')"
        >{{ t.cancel }}</button>
        <button
          class="toolbar-btn text-xs"
          @click="$emit('discard')"
        >{{ t.dontSave }}</button>
        <button
          class="toolbar-btn text-xs active"
          @click="$emit('save')"
        >{{ t.save }}</button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  show: Boolean,
  t: { type: Object, default: () => ({}) },
})
defineEmits(['save', 'discard', 'cancel'])
</script>
