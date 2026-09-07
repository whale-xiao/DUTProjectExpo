<script setup>
// Reddit 式搜索框：输入即联想、点击直达、回车进结果页；支持键盘 ↑↓ Enter Esc。
import { ref, watch, computed, onBeforeUnmount } from 'vue'
import { suggestProjects } from '@/api/projects'

const props = defineProps({
  placeholder: { type: String, default: '搜索项目名、技术栈、关键词…' },
  size: { type: String, default: 'large' }, // large | inline
  initial: { type: String, default: '' },
})
const emit = defineEmits(['select', 'submit'])

const keyword = ref(props.initial)
watch(
  () => props.initial,
  (v) => (keyword.value = v),
)

const suggestions = ref([])
const loading = ref(false)
const open = ref(false)
const activeIndex = ref(-1)
let timer = null
let ctrl = null

watch(keyword, (kw) => {
  clearTimeout(timer)
  kw = (kw || '').trim()
  if (!kw) {
    suggestions.value = []
    loading.value = false
    open.value = false
    activeIndex.value = -1
    return
  }
  activeIndex.value = -1
  timer = setTimeout(() => fetchSuggest(kw), 200)
})

async function fetchSuggest(kw) {
  ctrl?.abort?.()
  ctrl = new AbortController()
  loading.value = true
  open.value = true
  try {
    suggestions.value = await suggestProjects(kw, 6, ctrl.signal)
  } catch {
    suggestions.value = [] // 网络失败静默收起，不打断首页内容
  } finally {
    loading.value = false
  }
}

function choose(item) {
  open.value = false
  emit('select', item)
}

function submitKeyword() {
  const kw = keyword.value.trim()
  if (!kw) return
  open.value = false
  emit('submit', kw)
}

function onKeydown(e) {
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    open.value = true
    if (suggestions.value.length) activeIndex.value = (activeIndex.value + 1) % suggestions.value.length
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    if (activeIndex.value <= 0) activeIndex.value = suggestions.value.length - 1
    else activeIndex.value -= 1
  } else if (e.key === 'Enter') {
    e.preventDefault()
    if (activeIndex.value >= 0 && suggestions.value[activeIndex.value]) {
      choose(suggestions.value[activeIndex.value])
    } else {
      submitKeyword()
    }
  } else if (e.key === 'Escape') {
    open.value = false
  }
}

function onBlur() {
  // 延迟关闭，保证联想项 mousedown 先触发选择
  setTimeout(() => (open.value = false), 150)
}

function highlight(text, kw) {
  const t = text || ''
  if (!kw) return [{ text: t, hit: false }]
  const lt = t.toLowerCase()
  const lk = kw.toLowerCase()
  const out = []
  let i = 0
  for (;;) {
    const idx = lt.indexOf(lk, i)
    if (idx < 0) {
      out.push({ text: t.slice(i), hit: false })
      break
    }
    if (idx > i) out.push({ text: t.slice(i, idx), hit: false })
    out.push({ text: t.slice(idx, idx + kw.length), hit: true })
    i = idx + kw.length
  }
  return out
}

const showEmptyHint = computed(() => !loading.value && suggestions.value.length === 0)

onBeforeUnmount(() => {
  clearTimeout(timer)
  ctrl?.abort?.()
})
</script>

<template>
  <div class="searchbox" :class="size">
    <div class="searchbox-bar">
      <svg class="ico" aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
        <circle cx="11" cy="11" r="7"></circle>
        <path d="M20.5 20.5 L16.5 16.5"></path>
      </svg>
      <input
        v-model="keyword"
        class="input"
        :class="{ 'has-value': keyword }"
        type="text"
        :placeholder="placeholder"
        @keydown="onKeydown"
        @blur="onBlur"
        @focus="keyword && (open = true)"
      />
      <button v-if="keyword" class="clear" aria-label="清空" @mousedown.prevent="keyword = ''">×</button>
    </div>

    <div v-if="open && keyword" class="searchbox-drop">
      <template v-if="loading">
        <div class="drop-loading">正在搜索…</div>
      </template>

      <template v-else-if="suggestions.length">
        <button
          v-for="(item, i) in suggestions"
          :key="item.id"
          type="button"
          class="drop-item"
          :class="{ active: i === activeIndex }"
          @mousedown.prevent="choose(item)"
        >
          <span class="drop-thumb">
            <img v-if="item.coverUrl" :src="item.coverUrl" alt="" />
            <span v-else class="drop-thumb-ph">{{ item.name.slice(0, 1) }}</span>
          </span>
          <span class="drop-body">
            <span class="drop-title">
              <template v-for="(p, idx) in highlight(item.name, keyword)" :key="idx">
                <mark v-if="p.hit">{{ p.text }}</mark>
                <template v-else>{{ p.text }}</template>
              </template>
              <span v-if="item.isRecommended" class="rec">★</span>
            </span>
            <span class="drop-summary">{{ item.summary }}</span>
          </span>
          <span class="drop-tags">
            <i v-for="t in item.tags.slice(0, 3)" :key="t">{{ t }}</i>
          </span>
        </button>
        <div class="drop-foot" @mousedown.prevent="submitKeyword">回车查看「{{ keyword }}」的全部结果 →</div>
      </template>

      <template v-else-if="showEmptyHint">
        <div class="drop-empty">没有匹配「{{ keyword }}」的项目，回车查看全部结果</div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.searchbox {
  position: relative;
  width: 100%;
}
.searchbox.large .searchbox-bar {
  border-radius: var(--r-lg);
}
.searchbox.inline .searchbox-bar {
  border-radius: var(--r-md);
}
.searchbox-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 16px;
  height: 48px;
  background: var(--surface);
  border: 1px solid var(--hairline);
  box-shadow: var(--shadow-card);
  transition: border-color 0.12s var(--ease), box-shadow 0.12s var(--ease);
}
.searchbox-bar:focus-within {
  border-color: var(--brand);
  box-shadow: 0 0 0 3px var(--brand-soft);
}
.ico {
  width: 18px;
  height: 18px;
  color: var(--muted);
  flex: none;
}
.input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  font-size: 15px;
  color: var(--ink);
}
.input::placeholder {
  color: var(--muted);
}
.clear {
  border: none;
  background: var(--hairline);
  color: var(--sub);
  width: 22px;
  height: 22px;
  border-radius: 50%;
  line-height: 1;
  font-size: 14px;
}
.searchbox-drop {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  right: 0;
  background: var(--surface);
  border: 1px solid var(--hairline);
  border-radius: var(--r-md);
  box-shadow: 0 8px 24px rgba(30, 36, 51, 0.14);
  overflow: hidden;
  z-index: 30;
}
.drop-loading,
.drop-empty {
  padding: 14px 16px;
  color: var(--muted);
  font-size: 13px;
}
.drop-item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 10px 14px;
  border: none;
  background: transparent;
  text-align: left;
}
.drop-item.active,
.drop-item:hover {
  background: var(--brand-soft);
}
.drop-thumb {
  width: 46px;
  height: 34px;
  border-radius: var(--r-sm);
  overflow: hidden;
  background: var(--bg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--muted);
  font-size: 16px;
  flex: none;
}
.drop-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.drop-body {
  flex: 1;
  min-width: 0;
}
.drop-title {
  display: block;
  font-weight: 600;
  font-size: 14px;
  color: var(--ink);
}
.drop-title mark {
  background: var(--brand-soft);
  color: var(--brand-deep);
  border-radius: 2px;
}
.rec {
  color: var(--accent);
  margin-left: 4px;
  font-size: 12px;
}
.drop-summary {
  display: block;
  color: var(--sub);
  font-size: 12px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.drop-tags {
  display: flex;
  gap: 4px;
  flex: none;
}
.drop-tags i {
  font-style: normal;
  font-size: 11px;
  color: var(--sub);
  background: var(--bg);
  border: 1px solid var(--hairline);
  padding: 1px 6px;
  border-radius: 6px;
}
.drop-foot {
  padding: 9px 14px;
  font-size: 13px;
  color: var(--brand);
  border-top: 1px solid var(--hairline);
  cursor: pointer;
}
</style>
