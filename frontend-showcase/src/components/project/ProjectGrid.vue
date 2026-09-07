<script setup>
import ProjectCard from './ProjectCard.vue'

defineProps({
  cards: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  emptyText: { type: String, default: '暂无项目' },
  skeletonCount: { type: Number, default: 6 },
})
</script>

<template>
  <div v-if="loading" class="grid">
    <div v-for="i in skeletonCount" :key="i" class="skeleton">
      <div class="sk-cover"></div>
      <div class="sk-line w60"></div>
      <div class="sk-line w90"></div>
    </div>
  </div>

  <div v-else-if="cards.length" class="grid">
    <ProjectCard v-for="card in cards" :key="card.id" :card="card" />
  </div>

  <div v-else class="empty">
    <p>{{ emptyText }}</p>
    <slot name="empty-action"></slot>
  </div>
</template>

<style scoped>
.grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px 24px;
}
@media (max-width: 1023px) {
  .grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
@media (max-width: 639px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
.skeleton {
  border-radius: var(--r-md);
  border: 1px solid var(--hairline);
  overflow: hidden;
  background: var(--surface);
}
.sk-cover {
  aspect-ratio: 3 / 2;
  background: var(--bg);
}
.sk-line {
  height: 14px;
  margin: 12px 16px;
  background: var(--bg);
  border-radius: 6px;
}
.w60 {
  width: 60%;
}
.w90 {
  width: 90%;
}
.empty {
  text-align: center;
  color: var(--muted);
  padding: 64px 0;
}
.empty p {
  font-size: 15px;
}
</style>
