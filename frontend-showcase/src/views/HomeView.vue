<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getHome } from '@/api/home'
import SearchBox from '@/components/search/SearchBox.vue'
import ProjectGrid from '@/components/project/ProjectGrid.vue'

const router = useRouter()
const home = ref({ recommended: [], categories: [], latest: [] })
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    home.value = await getHome()
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
})

function goDetail(id) {
  router.push(`/project/${id}`)
}
function goSearch(kw) {
  router.push({ path: '/projects', query: { q: kw } })
}
</script>

<template>
  <!-- 首页即内容：顶部导航下方直接是搜索 + 内容流，不做巨型 hero -->
  <section class="hero">
    <h1 class="slogan">发现实训基地的优秀项目</h1>
    <div class="hero-search">
      <SearchBox size="large" @select="(item) => goDetail(item.id)" @submit="goSearch" />
    </div>
  </section>

  <div v-if="error" class="container error">
    <p>{{ error }}</p>
    <p class="hint">请确认后端已启动（go run . 于 :8080）后刷新。</p>
  </div>

  <section v-else class="container">
    <ProjectGrid v-if="loading" :cards="[]" :loading="true" />

    <template v-else>
      <div v-if="home.recommended.length" class="section">
        <div class="section-head">
          <h2>推荐项目</h2>
          <router-link to="/projects" class="more">查看全部 →</router-link>
        </div>
        <ProjectGrid :cards="home.recommended" />
      </div>

      <div v-if="home.latest.length" class="section">
        <div class="section-head">
          <h2>最新上架</h2>
        </div>
        <ProjectGrid :cards="home.latest" />
      </div>
    </template>
  </section>
</template>

<style scoped>
.hero {
  padding: 72px 24px 56px;
  text-align: center;
}
.slogan {
  margin: 0 0 28px;
  font-family: var(--serif);
  font-size: 34px;
  font-weight: 700;
  color: var(--ink);
}
.hero-search {
  max-width: 560px;
  margin: 0 auto;
}
.section {
  margin-bottom: 48px;
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
</style>
