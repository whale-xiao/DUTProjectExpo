<script setup>
// 全局布局：导航 + 路由出口 + 页脚。
//
// 导航有两种形态（见 docs/12 §4.9 / §4.10）：
//   - 首页：头图满屏，导航以 overlay 形式**嵌进图片顶部**（白字、透明底、不占行高），
//           搜索框作为中间一栏并入同一行，不再单独占一条
//   - 其余页：没有头图可嵌，回到常规的浅色行；搜索框由各页自行安排（列表页在页头）
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import SearchBox from '@/components/search/SearchBox.vue'

const route = useRoute()
const router = useRouter()
const overlay = computed(() => route.name === 'home')

function goDetail(id) {
  router.push(`/project/${id}`)
}
function goSearch(kw) {
  router.push({ path: '/projects', query: { q: kw } })
}
</script>

<template>
  <header class="site-nav" :class="{ overlay }">
    <div class="container nav-inner">
      <router-link to="/" class="logo">
        <svg class="logo-mark" viewBox="0 0 24 24" aria-hidden="true">
          <rect x="1" y="1" width="22" height="22" rx="5.5" />
          <path d="M12 5.4 L13.7 10.3 L18.6 12 L13.7 13.7 L12 18.6 L10.3 13.7 L5.4 12 L10.3 10.3 Z" fill="#fff" />
        </svg>
        <span>优秀项目展示系统</span>
      </router-link>

      <div v-if="overlay" class="nav-search">
        <SearchBox
          variant="glass"
          size="inline"
          @select="(item) => goDetail(item.id)"
          @submit="goSearch"
        />
      </div>

      <nav class="nav-links">
        <router-link to="/projects">全部项目</router-link>
      </nav>
    </div>
  </header>

  <main class="site-main">
    <router-view />
  </main>

  <footer class="site-footer">
    <div class="container">© 2026 大连理工大学实训基地 · 优秀项目展示</div>
  </footer>
</template>

<style scoped>
.site-nav {
  background: var(--bg);
  border-bottom: 1px solid var(--hairline);
}
.nav-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  height: var(--nav-h);
}
/* 搜索框占中间一栏。flex:1 吃掉剩余宽度，max-width 防止在大屏上过宽；
   超过 max-width 后剩余空间由 space-between 分摊到两侧，视觉上大致居中。 */
.nav-search {
  flex: 1;
  min-width: 0;
  max-width: 520px;
}
.logo {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 17px;
  font-weight: 700;
  color: var(--brand);
}
.logo-mark {
  width: 20px;
  height: 20px;
}
.logo-mark rect {
  fill: var(--brand);
}
.nav-links a {
  color: var(--sub);
  padding: 6px 10px;
  border-radius: var(--r-sm);
  transition:
    color var(--dur-fast) var(--ease-out),
    background var(--dur-fast) var(--ease-out);
}
.nav-links a:hover {
  color: var(--brand);
  background: var(--brand-soft);
}

/* —— 首页：嵌进头图顶部 ——
   用 absolute 而非 fixed：跟着头图一起滚走，不会在滚动后悬空压住浅色内容。
   文字对比度由 tokens.css 里 --veil-grad 的「顶部那道渐变」保证。 */
.site-nav.overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  z-index: 20; /* 高于头图内部的文字层(3)与箭头(4)，保证可点 */
  background: transparent;
  border-bottom: none;
}
.site-nav.overlay .logo {
  color: #fff;
}
.site-nav.overlay .logo-mark rect {
  fill: rgba(255, 255, 255, 0.22);
}
.site-nav.overlay .nav-links a {
  color: rgba(255, 255, 255, 0.86);
}
.site-nav.overlay .nav-links a:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.16);
}

.site-main {
  min-height: calc(100vh - var(--nav-h) - 56px);
}
.site-footer {
  border-top: 1px solid var(--hairline);
  color: var(--muted);
  font-size: 13px;
  padding: 16px 0;
}
</style>
