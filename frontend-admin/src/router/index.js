import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录' },
  },
  {
    path: '/',
    component: () => import('@/layout/Layout.vue'),
    redirect: '/dashboard',
    children: [
      { path: 'dashboard', name: 'dashboard', component: () => import('@/views/Dashboard.vue'), meta: { title: '仪表盘' } },
      { path: 'projects', name: 'project-list', component: () => import('@/views/ProjectList.vue'), meta: { title: '项目管理' } },
      { path: 'projects/create', name: 'project-create', component: () => import('@/views/ProjectEdit.vue'), meta: { title: '新建项目' } },
      { path: 'projects/:id/edit', name: 'project-edit', component: () => import('@/views/ProjectEdit.vue'), meta: { title: '编辑项目' } },
      { path: 'categories', name: 'category-manage', component: () => import('@/views/CategoryManage.vue'), meta: { title: '分类管理' } },
      { path: 'tags', name: 'tag-manage', component: () => import('@/views/TagManage.vue'), meta: { title: '标签管理' } },
    ],
  },
  { path: '/:pathMatch(.*)*', redirect: '/dashboard' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.name !== 'login' && !auth.isLoggedIn) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.name === 'login' && auth.isLoggedIn) {
    return { name: 'dashboard' }
  }
  return true
})

router.afterEach((to) => {
  document.title = to.meta?.title ? `${to.meta.title} · 内容管理后台` : '内容管理后台'
})

export default router
