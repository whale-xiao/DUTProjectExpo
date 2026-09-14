<script setup>
// 项目卡片。角标克制：
//   封面右上 = ★推荐（暖橙，全站唯一的强调色）
//   封面左上 = 分类（半透明深色胶囊，克制、不与推荐争抢）
//   底部右侧 = 浏览量（灰字，仅作信息补充）
// 分类名由父组件通过 categoryName 传入（CardDTO 只带 categoryId）。
const props = defineProps({
  card: { type: Object, required: true },
  categoryName: { type: String, default: '' },
})

const tags = () => props.card.tags || []

function fmtViews(n) {
  const v = Number(n) || 0
  if (v >= 10000) return `${(v / 10000).toFixed(1)} 万`
  return v.toLocaleString('zh-CN')
}
</script>

<template>
  <router-link class="card" :to="`/project/${card.id}`">
    <div class="cover">
      <img v-if="card.coverUrl" :src="card.coverUrl" :alt="card.name" loading="lazy" decoding="async" />
      <div v-else class="cover-ph">{{ card.name.slice(0, 1) }}</div>
      <span v-if="categoryName" class="cat">{{ categoryName }}</span>
      <span v-if="card.isRecommended" class="rec">★ 推荐</span>
    </div>
    <div class="body">
      <h3 class="title">{{ card.name }}</h3>
      <p class="summary">{{ card.summary }}</p>
      <div class="foot">
        <div class="tags">
          <i v-for="t in (card.tags || []).slice(0, 3)" :key="t">{{ t }}</i>
          <i v-if="(card.tags || []).length > 3" class="more">+{{ card.tags.length - 3 }}</i>
        </div>
        <span class="views">{{ fmtViews(card.viewCount) }} 浏览</span>
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
  transition: box-shadow var(--dur-base) var(--ease-out);
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
  transition: transform var(--dur-base) var(--ease-out);
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
.cat {
  position: absolute;
  top: 10px;
  left: 10px;
  font-size: 11px;
  line-height: 1;
  padding: 4px 8px;
  border-radius: 999px;
  color: #fff;
  background: rgba(8, 12, 20, 0.55);
  backdrop-filter: blur(4px);
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
.foot {
  display: flex;
  align-items: center;
  gap: 8px;
}
.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  min-width: 0;
  flex: 1;
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
.views {
  flex: none;
  font-size: 11px;
  color: var(--muted);
  white-space: nowrap;
}
</style>
