<script setup>
// 详情页影院头图：作品通栏剧照 + 双渐变遮罩 + 浮层返回 + 叠字标题（见 docs/12 §5）。
import { ref, computed, watch } from 'vue'

const props = defineProps({
  detail: { type: Object, required: true },
})
const emit = defineEmits(['back'])

const imgFailed = ref(false)
watch(
  () => props.detail?.coverUrl,
  () => (imgFailed.value = false),
)
const cover = computed(() => (imgFailed.value ? '' : props.detail?.coverUrl || ''))
const tags = computed(() => props.detail?.tags || [])

function fmtViews(n) {
  return n >= 1000 ? `${(n / 1000).toFixed(1)}k` : String(n)
}
</script>

<template>
  <header class="dhero">
    <div class="dhero-bg" aria-hidden="true">
      <img
        v-if="cover"
        class="dhero-img"
        :src="cover"
        alt=""
        fetchpriority="high"
        decoding="async"
        @error="imgFailed = true"
      />
      <div v-else class="dhero-fallback"></div>
    </div>
    <div class="dhero-veil" aria-hidden="true"></div>

    <button class="dhero-back" type="button" @click="emit('back')">
      <svg viewBox="0 0 24 24" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M15 5 L8 12 L15 19" />
      </svg>
      返回
    </button>

    <div class="container dhero-inner">
      <p class="crumb">
        <span v-if="detail.category">{{ detail.category.name }}</span>
        <span v-if="detail.category && detail.publishedAt" class="dot">·</span>
        <span v-if="detail.publishedAt">{{ detail.publishedAt }} 收录</span>
      </p>

      <h1 class="title">{{ detail.name }}</h1>
      <p v-if="detail.summary" class="summary">{{ detail.summary }}</p>

      <div class="meta">
        <span v-for="t in tags" :key="t" class="chip">{{ t }}</span>
        <span class="views">{{ fmtViews(detail.viewCount) }} 浏览</span>
      </div>
    </div>
  </header>
</template>

<style scoped>
.dhero {
  position: relative;
  height: var(--hero-h);
  overflow: hidden;
  display: flex;
  align-items: flex-end;
}
.dhero-bg,
.dhero-veil {
  position: absolute;
  inset: 0;
}
.dhero-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  animation: dheroIn 0.5s var(--ease-out) both;
}
.dhero-fallback {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, var(--brand-deep) 0%, var(--brand) 100%);
}
.dhero-veil {
  background: var(--veil-grad);
}

/* 返回按钮浮在头图左上角，不占首屏高度 */
.dhero-back {
  position: absolute;
  top: 20px;
  left: 24px;
  z-index: 2;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  height: 34px;
  padding: 0 14px 0 10px;
  border: 1px solid var(--cinema-line);
  border-radius: 999px;
  background: var(--cinema-chip);
  backdrop-filter: blur(6px);
  color: var(--cinema-ink);
  font-size: 13px;
  transition: background var(--dur-fast) var(--ease-out);
}
.dhero-back:hover {
  background: rgba(255, 255, 255, 0.24);
}
.dhero-back svg {
  width: 16px;
  height: 16px;
}

.dhero-inner {
  position: relative;
  z-index: 1;
  width: 100%;
  padding-bottom: 42px;
  color: var(--cinema-ink);
}
.crumb {
  margin: 0 0 10px;
  font-size: 13px;
  letter-spacing: 0.08em;
  color: var(--cinema-sub);
}
.crumb .dot {
  margin: 0 6px;
}
.title {
  margin: 0 0 12px;
  font-family: var(--serif);
  font-size: clamp(26px, 3.6vw, 44px);
  line-height: 1.18;
  font-weight: 700;
  color: var(--cinema-ink);
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
.meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}
.chip {
  font-size: 12px;
  padding: 3px 10px;
  border-radius: 999px;
  color: var(--cinema-ink);
  background: var(--cinema-chip);
  border: 1px solid var(--cinema-line);
  backdrop-filter: blur(4px);
}
.views {
  font-size: 13px;
  color: var(--cinema-sub);
  margin-left: 4px;
}

@keyframes dheroIn {
  from {
    opacity: 0;
    transform: scale(1.04);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@media (max-width: 639px) {
  .dhero {
    height: var(--hero-h-sm);
  }
  .dhero-back {
    top: 12px;
    left: 16px;
    height: 30px;
    padding: 0 12px 0 8px;
    font-size: 12px;
  }
  .dhero-inner {
    padding-bottom: 24px;
  }
  .dhero-inner,
  .dhero-back {
    /* 与全局 container 的 24px 内边距对齐 */
    padding-left: 20px;
    padding-right: 20px;
  }
  .dhero-back {
    padding: 0 12px 0 8px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .dhero-img {
    animation: none;
  }
}
</style>
