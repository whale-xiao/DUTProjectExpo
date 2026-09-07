<script setup>
// 全部项目页：兼作搜索结果页（?q=）。筛选/排序/搜索都改 URL query，路由变化即刷新。
import { ref, watch, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { listProjects } from '@/api/projects'
import { getCategories } from '@/api/meta'
import SearchBox from '@/components/search/SearchBox.vue'
import ProjectGrid from '@/components/project/ProjectGrid.vue'

const route = useRoute()
const router = useRouter()

const cards = ref([])
const total = ref(0)
const hasMore = ref(false)
const loading = ref(true)
const loadingMore = ref(false)
const error = ref('')
const page = ref(1)
const categories = ref([])

const q = computed(() => route.query.q || '')
const categoryId = computed(() => Number(route.query.category_id) || 0)
const sort = computed(() => route.query.sort || 'recommended')

const SORTS = [
  { value: 'recommended', label: '推荐' },
  { value: 'newest', label: '最新' },
  { value: 'views', label: '浏览最多' },
]

watch(() => route.query, () => fetchList(1), { immediate: true })
onMounted(async () => {
  try {
    categories.value = await getCategories()
  } catch {
    /* 分类加载失败不阻塞列表 */
  }
})

async function fetchList(p) {
  const params = {
    keyword: q.value,
    categoryId: categoryId.value,
    sort: sort.value,
    page: p,
    page_size: 20,
  }
  if (p === 1) loading.value = true
  else loadingMore.value = true
  try {
    const data = await listProjects(params)
    cards.value = p === 1 ? data.list : cards.value.concat(data.list)
    total.value = data.total
    hasMore.value = data.hasMore
    page.value = p
    error.value = ''
  } catch (e) {
    if (p === 1) {
      cards.value = []
      total.value = 0
      error.value = e.message
    }
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

function nav(next) {
  const query = {}
  if (next.q) query.q = next.q
  if (next.category_id) query.category_id = String(next.category_id)
  if (next.sort) query.sort = next.sort
  router.push({ path: '/projects', query })
}

function onSearch(kw) {
  nav({ q: kw, category_id: categoryId.value, sort: sort.value })
}
function pickCategory(id) {
  nav({ q: q.value, category_id: id || 0, sort: sort.value })
}
function changeSort(val) {
  nav({ q: q.value, category_id: categoryId.value, sort: val })
}
function clearAll() {
  router.push({ path: '/projects' })
}
function goDetail(id) {
  router.push(`/project/${id}`)
}

const showFilterNote = computed(() => Boolean(q.value) || categoryId.value > 0)
</script>

<template>
  <div class="container list-page">
    <h1 class="page-title">全部项目</h1>

    <div v-if="q" class="result-note">
      搜索「<strong>{{ q }}</strong>」共 {{ total }} 个结果
      <button class="clear" @click="clearAll">清空筛选 ✕</button>
    </div>

    <div class="toolbar">
      <SearchBox class="tb-search" size="inline" :initial="q" @select="(item) => goDetail(item.id)" @submit="onSearch" />
      <select class="sort" :value="sort" @change="changeSort($event.target.value)">
        <option v-for="s in SORTS" :key="s.value" :value="s.value">{{ s.label }}</option>
      </select>
    </div>

    <div v-if="categories.length" class="chips">
      <button :class="{ on: categoryId === 0 }" @click="pickCategory(0)">全部</button>
      <button
        v-for="c in categories"
        :key="c.id"
        :class="{ on: categoryId === c.id }"
        @click="pickCategory(c.id)"
      >
        {{ c.name }}
      </button>
    </div>

    <div v-if="error && !cards.length" class="error">
      <p>{{ error }}</p>
      <p class="hint">确认后端已启动后刷新。</p>
    </div>

    <ProjectGrid
      v-else
      :cards="cards"
      :loading="loading"
      empty-text="没有匹配的项目，试试清空筛选"
    />

    <div v-if="!loading && hasMore" class="load-more">
      <button :disabled="loadingMore" @click="fetchList(page + 1)">
        {{ loadingMore ? '加载中…' : '加载更多' }}
      </button>
    </div>
    <div v-if="!loading && cards.length && !hasMore" class="the-end">已经到底啦</div>
  </div>
</template>

<style scoped>
.list-page {
  padding-top: 40px;
  padding-bottom: 48px;
}
.page-title {
  margin: 0 0 16px;
  font-size: 26px;
}
.result-note {
  color: var(--sub);
  margin-bottom: 16px;
  font-size: 15px;
}
.result-note strong {
  color: var(--brand-deep);
}
.clear {
  border: none;
  background: var(--brand-soft);
  color: var(--brand);
  border-radius: var(--r-sm);
  padding: 2px 10px;
  margin-left: 10px;
  font-size: 13px;
}
.toolbar {
  display: flex;
  gap: 16px;
  align-items: center;
  margin-bottom: 16px;
}
.tb-search {
  max-width: 380px;
}
.sort {
  border: 1px solid var(--hairline);
  border-radius: var(--r-sm);
  background: var(--surface);
  padding: 6px 10px;
  color: var(--sub);
  font-size: 13px;
}
.chips {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 24px;
}
.chips button {
  border: 1px solid var(--hairline);
  background: var(--surface);
  color: var(--sub);
  border-radius: 999px;
  padding: 4px 14px;
  font-size: 13px;
}
.chips button.on {
  background: var(--brand-soft);
  color: var(--brand);
  border-color: var(--brand);
  font-weight: 500;
}
.error {
  text-align: center;
  color: var(--danger);
  padding: 40px 0;
}
.hint {
  color: var(--muted);
  font-size: 13px;
}
.load-more {
  text-align: center;
  margin-top: 32px;
}
.load-more button {
  border: 1px solid var(--hairline);
  background: var(--surface);
  color: var(--brand);
  border-radius: 999px;
  padding: 8px 32px;
}
.load-more button:disabled {
  color: var(--muted);
  cursor: default;
}
.the-end {
  text-align: center;
  color: var(--muted);
  font-size: 13px;
  margin-top: 28px;
}
</style>
