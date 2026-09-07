<script setup>
defineProps({
  card: { type: Object, required: true },
})
</script>

<template>
  <router-link class="card" :to="`/project/${card.id}`">
    <div class="cover">
      <img v-if="card.coverUrl" :src="card.coverUrl" :alt="card.name" loading="lazy" decoding="async" />
      <div v-else class="cover-ph">{{ card.name.slice(0, 1) }}</div>
      <span v-if="card.isRecommended" class="rec">★ 推荐</span>
    </div>
    <div class="body">
      <h3 class="title">{{ card.name }}</h3>
      <p class="summary">{{ card.summary }}</p>
      <div class="tags">
        <i v-for="t in card.tags.slice(0, 3)" :key="t">{{ t }}</i>
        <i v-if="card.tags.length > 3" class="more">+{{ card.tags.length - 3 }}</i>
      </div>
    </div>
  </router-link>
</template>

<style scoped>
.card {
  display: block;
  background: var(--surface);
  border: 1px solid var(--hairline);
  border-radius: var(--r-md);
  overflow: hidden;
  box-shadow: var(--shadow-card);
  transition: box-shadow 0.12s var(--ease);
}
.card:hover {
  box-shadow: var(--shadow-hover);
}
.cover {
  position: relative;
  aspect-ratio: 3 / 2;
  background: var(--bg);
  overflow: hidden;
}
.cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.12s var(--ease);
}
.card:hover .cover img {
  transform: scale(1.02);
}
.cover-ph {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 44px;
  color: var(--muted);
  font-family: var(--serif);
  font-weight: 700;
}
.rec {
  position: absolute;
  top: 10px;
  right: 10px;
  background: var(--accent);
  color: #fff;
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 999px;
}
.body {
  padding: 14px 16px 16px;
}
.title {
  margin: 0 0 6px;
  font-family: var(--serif);
  font-size: 20px;
  font-weight: 700;
  color: var(--ink);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.summary {
  margin: 0 0 10px;
  color: var(--sub);
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}
.tags i {
  font-style: normal;
  font-size: 11px;
  color: var(--sub);
  background: var(--bg);
  border: 1px solid var(--hairline);
  padding: 1px 6px;
  border-radius: 6px;
}
</style>
