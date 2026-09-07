<script setup>
// 若依风格后台布局：深色侧边菜单 + 顶栏（面包屑 + 用户下拉）
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import {
  Odometer,
  Files,
  CollectionTag,
  PriceTag,
  ArrowDown,
  SwitchButton,
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const activePath = computed(() => {
  const p = route.path
  if (p.startsWith('/projects')) return '/projects'
  if (p.startsWith('/categories')) return '/categories'
  if (p.startsWith('/tags')) return '/tags'
  return '/dashboard'
})

async function onCommand(cmd) {
  if (cmd === 'logout') {
    await ElMessageBox.confirm('确认退出登录？', '提示', { type: 'warning' })
    auth.logout()
    router.push('/login')
  }
}
</script>

<template>
  <el-container class="layout">
    <el-aside width="210px" class="aside">
      <div class="brand">
        <span class="dot"></span>
        优秀项目展示 · 内容管理
      </div>
      <el-menu
        :default-active="activePath"
        router
        background-color="#001529"
        text-color="#a0a9b4"
        active-text-color="#ffffff"
        class="menu"
      >
        <el-menu-item index="/dashboard">
          <el-icon><Odometer /></el-icon><span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/projects">
          <el-icon><Files /></el-icon><span>项目管理</span>
        </el-menu-item>
        <el-menu-item index="/categories">
          <el-icon><CollectionTag /></el-icon><span>分类管理</span>
        </el-menu-item>
        <el-menu-item index="/tags">
          <el-icon><PriceTag /></el-icon><span>标签管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container class="right">
      <el-header class="header" height="52px">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item>{{ route.meta?.title }}</el-breadcrumb-item>
        </el-breadcrumb>
        <el-dropdown trigger="click" @command="onCommand">
          <span class="user">
            <span class="avatar">{{ (auth.user?.nickname || '管')[0] }}</span>
            {{ auth.user?.nickname || auth.user?.username || '管理员' }}
            <el-icon><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="logout" :icon="SwitchButton">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>

      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.layout {
  height: 100%;
}
.aside {
  background: #001529;
  overflow: hidden;
}
.brand {
  height: 52px;
  display: flex;
  align-items: center;
  gap: 8px;
  color: #fff;
  font-weight: 600;
  font-size: 14px;
  padding: 0 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}
.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #1890ff;
}
.menu {
  border-right: none;
}
.menu :deep(.el-menu-item.is-active) {
  background-color: #1890ff !important;
}
.header {
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}
.user {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  color: #333;
  font-size: 14px;
}
.avatar {
  width: 26px;
  height: 26px;
  border-radius: 50%;
  background: #1890ff;
  color: #fff;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
}
.main {
  background: #f0f2f5;
}
</style>
