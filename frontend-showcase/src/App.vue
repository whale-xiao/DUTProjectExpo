<script setup>
// 顶部导航与页脚为全局布局；路由出口在中间。
</script>

<template>
  <header class="site-nav">
    <div class="container nav-inner">
      <router-link to="/" class="logo">
        <svg class="logo-mark" viewBox="0 0 24 24" aria-hidden="true">
          <rect x="1" y="1" width="22" height="22" rx="5.5" fill="var(--brand)" />
          <path d="M12 5.4 L13.7 10.3 L18.6 12 L13.7 13.7 L12 18.6 L10.3 13.7 L5.4 12 L10.3 10.3 Z" fill="#fff" />
        </svg>
        <span>优秀项目展示系统</span>
      </router-link>
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
  border-bottom: 1px solid var(--hairline);
  background: var(--bg);
}
.nav-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 56px;
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
.nav-links a {
  color: var(--sub);
  padding: 6px 10px;
  border-radius: var(--r-sm);
  transition: color 0.12s var(--ease), background 0.12s var(--ease);
}
.nav-links a:hover {
  color: var(--brand);
  background: var(--brand-soft);
}
.site-main {
  min-height: calc(100vh - 56px - 56px);
}
.site-footer {
  border-top: 1px solid var(--hairline);
  color: var(--muted);
  font-size: 13px;
  padding: 16px 0;
}
</style>
