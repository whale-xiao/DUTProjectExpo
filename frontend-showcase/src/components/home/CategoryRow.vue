<script setup>
// 首页分类合集行：某一分类下的项目横向滑动（见 docs/12 §6）。
// 数据由 HomeView 一次拉全量后按 categoryId 分组传入，本组件不发请求。
import { ref } from 'vue'

defineProps({
  title: { type: String, required: true },
  items: { type: Array, default: () => [] },
  categoryId: { type: [Number, String], default: '' },
})

const track = ref(null)

// 一次滚动约 80% 可视宽度，且不小于一张卡，避免窄屏滚不动
function scrollBy(dir) {
  const el = track.value
  if (!el) return
  el.scrollBy({ left: dir * Math.max(el.clientWidth * 0.8, 240), behavior: 'smooth' })
}

function fmtViews(n) {
  const v = Number(n) || 0
  return v >= 10000 ? `${(v / 10000).toFixed(1)} 万` : v.toLocaleString('zh-CN')
}
</script>

<template>
  <section class="crow">
    <header class="crow-head">
      <h2>{{ title }}</h2>
      <router-link
        class="more"
        :to="{ path: '/projects', query: categoryId ? { category_id: categoryId } : {} }"
      >
        更多 →
      </router-link>
    </header>

    <div class="crow-wrap">
      <button class="arrow prev" type="button" :aria-label="`向左滚动${title}`" @click="scrollBy(-1)">
        <svg viewBox="0 0 24 24" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M15 5 L8 12 L15 19" />
        </svg>
      </button>

      <ul ref="track" class="crow-track">
        <li v-for="it in items" :key="it.id">
          <router-link class="ccard" :to="`/project/${it.id}`">
            <span class="cc-cover">
              <img v-if="it.coverUrl" :src="it.coverUrl" :alt="it.name" loading="lazy" decoding="async" />
              <span v-else class="ph">{{ it.name.slice(0, 1) }}</span>
              <i v-if="it.isRecommended" class="rec">★</i>
            </span>
            <span class="cc-name">{{ it.name }}</span>
            <span class="cc-meta">{{ fmtViews(it.viewCount) }} 浏览</span>
          </router-link>
        </li>
      </ul>

      <button class="arrow next" type="button" :aria-label="`向右滚动${title}`" @click="scrollBy(1)">
        <svg viewBox="0 0 24 24" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M9 5 L16 12 L9 19" />
        </svg>
      </button>
    </div>
  </section>
</template>

<style scoped>
.crow-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-bottom: 14px;
}
.crow-head h2 {
  margin: 0;
  font-size: 20px;
  color: var(--ink);
}
.more {
  color: var(--brand);
  font-size: 14px;
}

.crow-wrap {
  position: relative;
}

/* 横向轨道：隐藏滚动条，用箭头或触摸滑动 */
.crow-track {
  display: flex;
  gap: 18px;
  list-style: none;
  margin: 0;
  padding: 2px 0 6px;
  overflow-x: auto;
  scroll-snap-type: x proximity;
  scrollbar-width: none;
  -ms-overflow-style: none;
}
.crow-track::-webkit-scrollbar {
  display: none;
}
.crow-track > li {
  flex: 0 0 236px;
  scroll-snap-align: start;
}

.ccard {
  display: block;
}
.cc-cover {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  aspect-ratio: 3 / 2;
  border-radius: var(--r-md);
  overflow: hidden;
  background: var(--bg);
  border: 1px solid var(--hairline);
}
.cc-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform var(--dur-base) var(--ease-out);
}
.ccard:hover .cc-cover img {
  transform: scale(1.04);
}
.ph {
  font-family: var(--serif);
  font-size: 32px;
  font-weight: 700;
  color: var(--muted);
}
.rec {
  position: absolute;
  top: 8px;
  right: 8px;
  font-style: normal;
  font-size: 11px;
  line-height: 1;
  padding: 3px 6px;
  border-radius: 999px;
  background: var(--accent);
  color: #fff;
}
.cc-name {
  display: block;
  margin-top: 9px;
  font-family: var(--serif);
  font-size: 16px;
  font-weight: 700;
  color: var(--ink);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.cc-meta {
  display: block;
  font-size: 12px;
  color: var(--muted);
  margin-top: 2px;
}

/* 左右箭头：悬停时浮现，触屏常显 */
.arrow {
  position: absolute;
  top: 38%;
  transform: translateY(-50%);
  z-index: 2;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--hairline);
  border-radius: 999px;
  background: var(--surface);
  color: var(--sub);
  box-shadow: var(--shadow-hover);
  opacity: 0;
  transition: opacity var(--dur-fast) var(--ease-out);
}
.arrow svg {
  width: 18px;
  height: 18px;
}
.arrow.prev {
  left: -14px;
}
.arrow.next {
  right: -14px;
}
.crow-wrap:hover .arrow,
.arrow:focus-visible {
  opacity: 1;
}

@media (max-width: 639px) {
  .crow-track {
    gap: 12px;
  }
  .crow-track > li {
    flex: 0 0 180px;
  }
  .arrow {
    display: none; /* 触屏直接滑动，不需要箭头 */
  }
}
</style>
