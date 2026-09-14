import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView,
    meta: { title: '优秀项目展示系统' },
  },
  {
    path: '/projects',
    name: 'projects',
    component: () => import('@/views/ProjectListView.vue'),
    meta: { title: '全部项目' },
  },
  {
    path: '/project/:id',
    name: 'project-detail',
    component: () => import('@/views/ProjectDetailView.vue'),
    meta: { title: '项目详情' },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/views/NotFoundView.vue'),
    meta: { title: '页面不存在' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  },
})

// 标签页标题。首页直接用站名（原来会拼成「优秀项目展示系统 · 优秀项目展示」，重复了一次），
// 其余页用「页面名 · 站名」。顶部导航移除后，站名只在这里出现。
const SITE_NAME = '优秀项目展示系统'

router.afterEach((to) => {
  const base = to.meta?.title
  document.title = !base || to.name === 'home' ? SITE_NAME : `${base} · ${SITE_NAME}`
})

export default router
