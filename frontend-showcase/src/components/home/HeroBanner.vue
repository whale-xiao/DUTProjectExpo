<script setup>
// 首页影院头图：主推作品轮播（自动播放 + 左右箭头 + 圆点导航）。见 docs/12 §4。
// 纯展示组件——数据由 HomeView 传入，自身不发请求（对齐 09 §6.2 分层约定）。
//
// 无障碍与性能：
//  - 悬停 / 聚焦 / 标签页隐藏时暂停自动播放（WCAG 2.2.2 的"可暂停"机制）
//  - prefers-reduced-motion 下完全关闭自动播放与滑动过渡
//  - 只有首屏第一张图带 fetchpriority=high，其余 lazy；非活动帧对读屏隐藏
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
  items: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  eyebrow: { type: String, default: '本周主推' },
})

const INTERVAL = 6000

const index = ref(0)
const paused = ref(false)
const reduced = ref(false)
const failed = ref([]) // 按索引记录加载失败的图，回落到该帧的品牌渐变

const count = computed(() => props.items?.length || 0)
const active = computed(() => props.items?.[index.value] || null)
const multi = computed(() => count.value > 1)
const coverOf = (it, i) => (failed.value[i] ? '' : it?.coverUrl || '')

// 某一帧的图挂了：只标记该帧回落到品牌渐变，不影响其余帧
function markFailed(i) {
  if (!failed.value.includes(i)) failed.value = [...failed.value, i]
}

// 数据刷新后把索引夹回有效范围
watch(count, (n) => {
  if (!n) index.value = 0
  else if (index.value > n - 1) index.value = n - 1
  failed.value = []
})

// —— 轮播控制 ——
let timer = null

function go(i) {
  if (!count.value) return
  index.value = (i + count.value) % count.value
  schedule()
}
function next() {
  go(index.value + 1)
}
function prev() {
  go(index.value - 1)
}

const autoplay = computed(() => multi.value && !paused.value && !reduced.value)

function schedule() {
  clearInterval(timer)
  if (autoplay.value) timer = setInterval(next, INTERVAL)
}

watch(autoplay, schedule)

// 切换标签页时暂停，避免后台空转与回来时"跳帧"
function onVisibility() {
  paused.value = document.hidden
}

let mq = null
function onMqChange(e) {
  reduced.value = e.matches
}

onMounted(() => {
  mq = window.matchMedia('(prefers-reduced-motion: reduce)')
  reduced.value = mq.matches
  mq.addEventListener('change', onMqChange)
  document.addEventListener('visibilitychange', onVisibility)
  schedule()
})

onBeforeUnmount(() => {
  clearInterval(timer)
  mq?.removeEventListener('change', onMqChange)
  document.removeEventListener('visibilitychange', onVisibility)
})
</script>

<template>
  <section
    class="hero"
    @mouseenter="paused = true"
    @mouseleave="paused = false"
    @focusin="paused = true"
    @focusout="paused = false"
    @keydown.left.prevent="prev"
    @keydown.right.prevent="next"
  >
    <!-- 图片轨：只放图，文字在下方静态层，保证滑动时文字不跟着抖 -->
    <div class="hero-viewport">
      <div class="hero-track" :style="{ transform: `translateX(-${index * 100}%)` }">
        <div
          v-for="(it, i) in items"
          :key="it.id"
          class="hero-slide"
          :aria-hidden="i !== index ? 'true' : null"
        >
          <img
            v-if="coverOf(it, i)"
            class="hero-img"
            :src="coverOf(it, i)"
            alt=""
            :fetchpriority="i === 0 ? 'high' : 'auto'"
            :loading="i === 0 ? 'eager' : 'lazy'"
            decoding="async"
            @error="markFailed(i)"
          />
          <div v-else class="hero-fallback"></div>
        </div>
      </div>
    </div>

    <!-- 当前帧的整块点击热区（在遮罩之下、文字之上） -->
    <router-link
      v-if="active && !loading"
      class="hero-hit"
      :to="`/project/${active.id}`"
      :aria-label="`进入作品 ${active.name}`"
    ></router-link>

    <div class="hero-veil" aria-hidden="true"></div>

    <div class="container hero-bottom">
      <!-- 骨架：占位高度与成功态一致，避免 CLS -->
      <div v-if="loading" class="hero-sk">
        <div class="sk sk-eyebrow"></div>
        <div class="sk sk-title"></div>
        <div class="sk sk-line"></div>
        <div class="sk sk-line short"></div>
        <div class="sk sk-btn"></div>
      </div>

      <template v-else-if="active">
        <!-- key 绑定 id：切帧时重放淡入动画 -->
        <div :key="active.id" class="hero-text">
          <p class="eyebrow">{{ eyebrow }}</p>
          <h1 class="title">
            <router-link :to="`/project/${active.id}`">{{ active.name }}</router-link>
          </h1>
          <p v-if="active.summary" class="summary">{{ active.summary }}</p>
          <div v-if="active.tags?.length" class="tags">
            <i v-for="t in active.tags.slice(0, 3)" :key="t">{{ t }}</i>
          </div>
          <div class="actions">
            <router-link class="btn btn-solid" :to="`/project/${active.id}`">进入作品 →</router-link>
            <router-link class="btn btn-ghost" to="/projects">查看全部项目</router-link>
          </div>
        </div>

        <div v-if="multi" class="hero-dots" role="tablist" aria-label="主推作品切换">
          <button
            v-for="(it, i) in items"
            :key="it.id"
            type="button"
            class="dot"
            :class="{ on: i === index }"
            role="tab"
            :aria-selected="i === index"
            :aria-label="`第 ${i + 1} 个主推作品：${it.name}`"
            @click="go(i)"
          ></button>
        </div>
      </template>

      <!-- 无推荐也无最新（含接口失败）：保留品牌标语，头图下方搜索框照常可用 -->
      <template v-else>
        <div class="hero-text">
          <p class="eyebrow">大连理工大学实训基地</p>
          <h1 class="title">发现实训基地的优秀项目</h1>
        </div>
      </template>
    </div>

    <!-- 左右箭头：仅多帧时出现 -->
    <template v-if="multi && !loading">
      <button class="hero-arrow prev" type="button" aria-label="上一个主推作品" @click="prev">
        <svg viewBox="0 0 24 24" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M15 5 L8 12 L15 19" />
        </svg>
      </button>
      <button class="hero-arrow next" type="button" aria-label="下一个主推作品" @click="next">
        <svg viewBox="0 0 24 24" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M9 5 L16 12 L9 19" />
        </svg>
      </button>
    </template>
  </section>
</template>

<style scoped>
/* 满屏：图片正好占满「导航栏以下」的第一屏；下方内容靠滚动才出现 */
.hero {
  position: relative;
  height: var(--hero-h-screen);
  overflow: hidden;
  display: flex;
  align-items: flex-end;
}

/* —— 图片轨 —— */
.hero-viewport {
  position: absolute;
  inset: 0;
  overflow: hidden;
}
.hero-track {
  display: flex;
  height: 100%;
  transition: transform var(--dur-slide) var(--ease-out);
  will-change: transform;
}
.hero-slide {
  flex: 0 0 100%;
  height: 100%;
  position: relative;
}
.hero-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
.hero-fallback {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, var(--brand-deep) 0%, var(--brand) 100%);
}

/* —— 点击热区 / 遮罩 / 内容层 —— */
.hero-hit {
  position: absolute;
  inset: 0;
  z-index: 1;
}
.hero-veil {
  position: absolute;
  inset: 0;
  z-index: 2;
  background: var(--veil-grad);
  pointer-events: none;
}
.hero-bottom {
  position: relative;
  z-index: 3;
  width: 100%;
  /* 留出空间给下方上浮的搜索框，避免首屏底部过分拥挤 */
  padding-bottom: 64px;
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 28px;
  /* 容器本身不拦截点击，让空白处的点击落到底层热区 */
  pointer-events: none;
}
.hero-text {
  flex: 1;
  min-width: 0;
  color: var(--cinema-ink);
  animation: textIn 0.5s var(--ease-out) both;
}
.hero-bottom a,
.hero-bottom button {
  pointer-events: auto;
}

/* —— 文案 —— */
.eyebrow {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin: 0 0 12px;
  font-size: 13px;
  letter-spacing: 0.16em;
  color: var(--cinema-sub);
}
.eyebrow::before {
  content: '';
  width: 22px;
  height: 2px;
  border-radius: 2px;
  background: var(--accent);
}
.title {
  margin: 0 0 12px;
  font-family: var(--serif);
  font-size: clamp(28px, 4vw, 48px);
  line-height: 1.15;
  font-weight: 700;
  color: var(--cinema-ink);
}
.title a {
  color: inherit;
}
.title a:hover {
  text-decoration: underline;
  text-underline-offset: 6px;
  text-decoration-thickness: 2px;
}
.summary {
  margin: 0 0 14px;
  max-width: 620px;
  font-size: 15px;
  line-height: 1.7;
  color: var(--cinema-sub);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 20px;
}
.tags i {
  font-style: normal;
  font-size: 12px;
  padding: 3px 10px;
  border-radius: 999px;
  color: var(--cinema-ink);
  background: var(--cinema-chip);
  border: 1px solid var(--cinema-line);
  backdrop-filter: blur(4px);
}
.actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}
.btn {
  display: inline-flex;
  align-items: center;
  height: 42px;
  padding: 0 22px;
  border-radius: 999px;
  font-size: 14px;
  font-weight: 600;
  transition:
    background var(--dur-base) var(--ease-out),
    transform var(--dur-base) var(--ease-out);
}
.btn-solid {
  background: var(--surface);
  color: var(--brand-deep);
}
.btn-solid:hover {
  transform: translateY(-1px);
}
.btn-ghost {
  color: var(--cinema-ink);
  border: 1px solid var(--cinema-line);
}
.btn-ghost:hover {
  background: var(--cinema-chip);
}

/* —— 圆点 —— */
.hero-dots {
  flex: none;
  display: flex;
  align-items: center;
  gap: 8px;
  padding-bottom: 14px;
  pointer-events: auto;
}
.dot {
  width: 8px;
  height: 8px;
  padding: 0;
  border: none;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.36);
  transition:
    width var(--dur-base) var(--ease-out),
    background var(--dur-base) var(--ease-out);
}
.dot:hover {
  background: rgba(255, 255, 255, 0.62);
}
.dot.on {
  width: 22px;
  background: var(--cinema-ink);
}

/* —— 左右箭头 —— */
.hero-arrow {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  z-index: 4;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--cinema-line);
  border-radius: 999px;
  background: var(--cinema-chip);
  backdrop-filter: blur(6px);
  color: var(--cinema-ink);
  opacity: 0;
  transition:
    opacity var(--dur-fast) var(--ease-out),
    background var(--dur-base) var(--ease-out);
}
.hero-arrow svg {
  width: 20px;
  height: 20px;
}
.hero-arrow.prev {
  left: 24px;
}
.hero-arrow.next {
  right: 24px;
}
.hero-arrow:hover {
  background: rgba(255, 255, 255, 0.24);
}
/* 悬停 / 键盘聚焦时才现身，避免静止时干扰画面 */
.hero:hover .hero-arrow,
.hero:focus-within .hero-arrow {
  opacity: 1;
}

/* —— 骨架 —— */
.hero-sk .sk {
  border-radius: var(--r-sm);
  background: rgba(255, 255, 255, 0.16);
}
.sk-eyebrow {
  width: 96px;
  height: 14px;
  margin-bottom: 14px;
}
.sk-title {
  width: min(420px, 60%);
  height: 46px;
  margin-bottom: 14px;
}
.sk-line {
  width: min(560px, 80%);
  height: 15px;
  margin-bottom: 10px;
}
.sk-line.short {
  width: min(360px, 52%);
  margin-bottom: 22px;
}
.sk-btn {
  width: 132px;
  height: 42px;
  border-radius: 999px;
}

@keyframes textIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 639px) {
  /* 移动端同样满屏（--hero-h-screen 已随 dvh 自适应地址栏） */
  .hero-bottom {
    padding-bottom: 44px;
    gap: 16px;
  }
  .summary {
    font-size: 14px;
  }
  .actions {
    gap: 10px;
  }
  .btn {
    height: 38px;
    padding: 0 18px;
    font-size: 13px;
  }
  .hero-arrow {
    width: 36px;
    height: 36px;
    opacity: 1; /* 触屏没有 hover，直接常显 */
  }
  .hero-arrow.prev {
    left: 12px;
  }
  .hero-arrow.next {
    right: 12px;
  }
  .hero-dots {
    padding-bottom: 10px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .hero-track {
    transition: none;
  }
  .hero-text {
    animation: none;
  }
  .dot,
  .btn,
  .hero-arrow {
    transition: none;
  }
}
</style>
