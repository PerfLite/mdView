<template>
  <div
    ref="el"
    class="preview-content h-full overflow-y-auto p-6"
    v-html="html"
    @scroll="onScroll"
  />
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({ html: String })
const emit = defineEmits(['scroll'])

const el = ref(null)

function onScroll() {
  const t = el.value
  const pct = t.scrollTop / (t.scrollHeight - t.clientHeight || 1)
  emit('scroll', pct)
}

function scrollTo(pct) {
  if (!el.value) return
  const max = el.value.scrollHeight - el.value.clientHeight
  el.value.scrollTop = Math.min(1, Math.max(0, pct)) * max
}

defineExpose({ scrollTo })
</script>
