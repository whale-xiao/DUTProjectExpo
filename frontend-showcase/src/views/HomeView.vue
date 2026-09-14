<script setup>
// 首页：影院头图 + 热力榜 + 推荐 + 分类合集行 + 最新（见 docs/12 §4/§6）。
// 搜索框已并入顶部导航行，见 App.vue。
import { ref, computed, onMounted } from 'vue'
import { getHome } from '@/api/home'
import { listProjects } from '@/api/projects'
import ProjectGrid from '@/components/project/ProjectGrid.vue'
import HeroBanner from '@/components/home/HeroBanner.vue'
import RankBoard from '@/components/home/RankBoard.vue'
import CategoryRow from '@/components/home/CategoryRow.vue'

const home = ref({ recommended: [], categories: [], latest: [] })
const loading = ref(true)
const error = ref('')

// 全量项目（一次拉取，供热力榜与分类合集行复用；按浏览量降序）
const allProjects = ref([])
const allLoading = ref(true)

onMounted(() => {
  loadHome()
  loadAll()
})

async function loadHome() {
  try {
    home.value = await getHome()
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

async function loadAll() {
  try {
    // page_size 上限 100，30 个项目规模一次取全即可，避免每分类各发一次请求
    const page = await listProjects({ sort: 'views', page: 1, page_size: 100 })
    allProjects.value = page.list || []
  } catch {
    allProjects.value = [] // 次要内容，失败不打断首页
  } finally {
    allLoading.value = false
  }
}

// 头图轮播：优先用推荐（最多 5 条，避免一次加载过多大图），无推荐时回落到最新
const featured = computed(() => {
  const rec = home.value.recommended || []
  const src = rec.length ? rec : home.value.latest || []
  return src.slice(0, 5)
})
const featuredEyebrow = computed(() => (home.value.recommended?.length ? '本周主推' : '最新作品'))

// 热力榜取前 10（接口已按 views 降序）
const hotList = computed(() => allProjects.value.slice(0, 10))

// 分类合集行：按分类分组，空分类不渲染
const rows = computed(() => {
  const cats = home.value.categories || []
  return cats
    .map((c) => ({ category: c, items: allProjects.value.filter((p) => p.categoryId === c.id) }))
    .filter((r) => r.items.length > 0)
})
</script>

<template>
  <HeroBanner :items="featured" :loading="loading" :eyebrow="featuredEyebrow" />

  <!-- 搜索框已并入顶部导航行（见 App.vue），此处不再单独占一条 -->

  <div v-if="error" class="container error">
    <p>{{ error }}</p>
    <p class="hint">请确认后端已启动（go run . 于 :8080）后刷新。</p>
  </div>

  <section v-else class="container">
    <ProjectGrid v-if="loading" :cards="[]" :loading="true" />

    <template v-else>
      <div v-if="allLoading || hotList.length" class="section" v-reveal>
        <RankBoard :items="hotList" :loading="allLoading" />
      </div>

      <div v-if="home.recommended.length" class="section">
        <div class="section-head" v-reveal>
          <h2>推荐项目</h2>
          <router-link to="/projects" class="more">查看全部 →</router-link>
        </div>
        <!-- 卡片各自错开浮现，故外层不加 v-reveal，避免父子双重动画 -->
        <ProjectGrid :cards="home.recommended" :categories="home.categories" />
      </div>

      <CategoryRow
        v-for="(r, i) in rows"
        :key="r.category.id"
        v-reveal="i * 60"
        class="section"
        :title="r.category.name"
        :items="r.items"
        :category-id="r.category.id"
      />

      <div v-if="home.latest.length" class="section">
        <div class="section-head" v-reveal>
          <h2>最新上架</h2>
        </div>
        <ProjectGrid :cards="home.latest" :categories="home.categories" />
      </div>
    </template>
  </section>
</template>

<style scoped>
/* 搜索框已移入顶部导航行（见 App.vue），首页不再有独立搜索区 */
.section {
  display: block;
  margin-bottom: 52px;
}
.section-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-bottom: 16px;
}
.section-head h2 {
  margin: 0;
  font-size: 20px;
  color: var(--ink);
}
.more {
  color: var(--brand);
  font-size: 14px;
}
.error {
  padding: 64px 0;
  text-align: center;
  color: var(--danger);
}
.hint {
  color: var(--muted);
  font-size: 13px;
}

@media (max-width: 639px) {
  .section {
    margin-bottom: 36px;
  }
}
</style>
