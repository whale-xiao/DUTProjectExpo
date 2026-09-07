<template>
  <div class="container not-found">
    <p class="code">404</p>
    <p>页面不存在或已被移走</p>
    <router-link to="/" class="home">回到首页</router-link>
  </div>
</template>

<style scoped>
.not-found {
  text-align: center;
  padding: 96px 0;
  color: var(--sub);
}
.code {
  font-size: 64px;
  font-weight: 700;
  color: var(--hairline);
  margin: 0;
}
.home {
  display: inline-block;
  margin-top: 16px;
  color: var(--brand);
  border: 1px solid var(--brand);
  border-radius: 999px;
  padding: 8px 24px;
}
</style>
