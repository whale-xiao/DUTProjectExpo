<script setup>
// 首页热力榜：按浏览量排序的排行榜（见 docs/12 §6）。
// 与卡片网格形成不同的浏览节奏——左侧榜首放大，右侧 2-N 名紧凑列表。
// 纯展示组件，数据由 HomeView 传入。
import { computed } from 'vue'

const props = defineProps({
  items: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
})

const top = computed(() => props.items.slice(0, 10))
const champ = computed(() => top.value[0] || null)
const rest = computed(() => top.value.slice(1))

function fmtViews(n) {
  const v = Number(n) || 0
  return v >= 10000 ? `${(v / 10000).toFixed(1)} 万` : v.toLocaleString('zh-CN')
}
// 序号统一补零，视觉上对齐
function no(i) {
  return String(i).padStart(2, '0')
}
</script>

<template>
  <section class="rank">
    <header class="rank-head">
      <h2>热力榜</h2>
      <span class="sub">按浏览量排序</span>
      <router-link class="more" :to="{ path: '/projects', query: { sort: 'views' } }">查看全部 →</router-link>
    </header>

    <!-- 骨架 -->
    <div v-if="loading" class="rank-body">
      <div class="sk sk-champ"></div>
      <div class="sk-lines">
        <div v-for="i in 5" :key="i" class="sk sk-row"></div>
      </div>
    </div>

    <div v-else-if="champ" class="rank-body">
      <!-- 榜首：放大展示 -->
      <router-link class="champ" :to="`/project/${champ.id}`">
        <div class="champ-cover">
          <img v-if="champ.coverUrl" :src="champ.coverUrl" :alt="champ.name" loading="lazy" decoding="async" />
          <div v-else class="ph">{{ champ.name.slice(0, 1) }}</div>
          <span class="champ-no">{{ no(1) }}</span>
        </div>
        <div class="champ-info">
          <h3>{{ champ.name }}</h3>
          <p v-if="champ.summary">{{ champ.summary }}</p>
          <span class="views">{{ fmtViews(champ.viewCount) }} 浏览</span>
        </div>
      </router-link>

      <!-- 2-N 名：紧凑列表 -->
      <ol v-if="rest.length" class="rank-list">
        <li v-for="(it, i) in rest" :key="it.id">
          <router-link class="row" :to="`/project/${it.id}`">
            <span class="no">{{ no(i + 2) }}</span>
            <span class="thumb">
              <img v-if="it.coverUrl" :src="it.coverUrl" :alt="it.name" loading="lazy" decoding="async" />
              <span v-else class="ph-sm">{{ it.name.slice(0, 1) }}</span>
            </span>
            <span class="txt">
              <span class="nm">{{ it.name }}</span>
              <span class="sm">{{ it.summary }}</span>
            </span>
            <span class="vw">{{ fmtViews(it.viewCount) }}</span>
          </router-link>
        </li>
      </ol>
    </div>

    <p v-else class="empty">暂无数据</p>
  </section>
</template>

<style scoped>
.rank-head {
  display: flex;
  align-items: baseline;
  gap: 10px;
  margin-bottom: 16px;
}
.rank-head h2 {
  margin: 0;
  font-size: 20px;
  color: var(--ink);
}
.sub {
  font-size: 12px;
  color: var(--muted);
}
.more {
  margin-left: auto;
  color: var(--brand);
  font-size: 14px;
}

.rank-body {
  display: grid;
  grid-template-columns: 340px minmax(0, 1fr);
  gap: 28px;
  align-items: start;
}

/* —— 榜首 —— */
.champ {
  display: block;
  background: var(--surface);
  border: 1px solid var(--hairline);
  border-radius: var(--r-md);
  overflow: hidden;
  box-shadow: var(--shadow-card);
  transition: box-shadow var(--dur-base) var(--ease-out);
}
.champ:hover {
  box-shadow: var(--shadow-hover);
}
.champ-cover {
  position: relative;
  aspect-ratio: 3 / 2;
  background: var(--bg);
  overflow: hidden;
}
.champ-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform var(--dur-base) var(--ease-out);
}
.champ:hover .champ-cover img {
  transform: scale(1.03);
}
.champ-no {
  position: absolute;
  left: 0;
  bottom: 0;
  padding: 2px 12px 4px;
  font-family: var(--serif);
  font-size: 30px;
  font-weight: 700;
  line-height: 1;
  color: #fff;
  background: var(--accent);
  border-top-right-radius: var(--r-md);
}
.champ-info {
  padding: 14px 16px 16px;
}
.champ-info h3 {
  margin: 0 0 6px;
  font-family: var(--serif);
  font-size: 21px;
  font-weight: 700;
  color: var(--ink);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.champ-info p {
  margin: 0 0 10px;
  font-size: 13px;
  color: var(--sub);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.views {
  font-size: 12px;
  color: var(--muted);
}

/* —— 2-N 列表 —— */
.rank-list {
  list-style: none;
  margin: 0;
  padding: 0;
}
.rank-list li + li {
  border-top: 1px solid var(--hairline);
}
.row {
  display: grid;
  grid-template-columns: 34px 76px minmax(0, 1fr) auto;
  align-items: center;
  gap: 12px;
  padding: 9px 8px;
  border-radius: var(--r-sm);
  transition: background var(--dur-fast) var(--ease-out);
}
.row:hover {
  background: var(--brand-soft);
}
.no {
  font-family: var(--serif);
  font-size: 17px;
  font-weight: 700;
  color: var(--muted);
  text-align: center;
}
.row:hover .no {
  color: var(--brand);
}
.thumb {
  width: 76px;
  aspect-ratio: 3 / 2;
  border-radius: var(--r-sm);
  overflow: hidden;
  background: var(--bg);
  display: flex;
  align-items: center;
  justify-content: center;
}
.thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
.ph,
.ph-sm {
  color: var(--muted);
  font-family: var(--serif);
  font-weight: 700;
}
.ph {
  font-size: 34px;
}
.ph-sm {
  font-size: 18px;
}
.txt {
  min-width: 0;
}
.nm {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: var(--ink);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.sm {
  display: block;
  font-size: 12px;
  color: var(--sub);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.vw {
  font-size: 12px;
  color: var(--muted);
  font-variant-numeric: tabular-nums;
}

/* —— 骨架 / 空态 —— */
.rank-body .sk {
  background: var(--bg);
  border-radius: var(--r-md);
}
.sk-champ {
  height: 300px;
}
.sk-lines {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.sk-row {
  height: 56px;
}
.empty {
  color: var(--muted);
  font-size: 14px;
  padding: 32px 0;
}

@media (max-width: 1023px) {
  .rank-body {
    grid-template-columns: 1fr;
    gap: 18px;
  }
  .champ-cover {
    aspect-ratio: 16 / 9;
  }
}

@media (max-width: 639px) {
  .row {
    grid-template-columns: 28px 60px minmax(0, 1fr) auto;
    gap: 8px;
  }
  .thumb {
    width: 60px;
  }
  .sm {
    display: none;
  }
}
</style>
