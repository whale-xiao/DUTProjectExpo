<script setup>
// 项目详情：作品文章（blocks 驱动渲染），封面大图 + 窄列正文。
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getProject } from '@/api/projects'
import BlockRenderer from '@/components/detail/BlockRenderer.vue'

const route = useRoute()
const router = useRouter()

const detail = ref(null)
const loading = ref(true)
const error = ref('')

watch(() => route.params.id, load, { immediate: true })

async function load() {
  loading.value = true
  error.value = ''
  detail.value = null
  try {
    detail.value = await getProject(Number(route.params.id))
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

function goBack() {
  if (window.history.length > 1) router.back()
  else router.push('/projects')
}

function fmtViews(n) {
  return n >= 1000 ? `${(n / 1000).toFixed(1)}k` : String(n)
}
</script>

<template>
  <div class="container detail">
    <div v-if="loading" class="loading">
      <div class="sk title"></div>
      <div class="sk line w50"></div>
      <div class="sk cover"></div>
      <div class="sk line"></div>
      <div class="sk line"></div>
      <div class="sk line w80"></div>
    </div>

    <div v-else-if="error && !detail" class="not-found">
      <p>这篇作品不存在或已下架</p>
      <button @click="router.push('/projects')">返回全部项目</button>
    </div>

    <template v-else-if="detail">
      <button class="back" @click="goBack">← 返回</button>

      <h1 class="name">{{ detail.name }}</h1>
      <p class="summary">{{ detail.summary }}</p>

      <div class="meta">
        <span v-for="t in detail.tags" :key="t" class="tag">{{ t }}</span>
        <span v-if="detail.category" class="meta-item">{{ detail.category.name }}</span>
        <span class="meta-item">{{ fmtViews(detail.viewCount) }} 浏览</span>
        <span class="meta-item">{{ detail.publishedAt }} 收录</span>
      </div>

      <div v-if="detail.coverUrl" class="cover">
        <img :src="detail.coverUrl" :alt="detail.name" />
      </div>

      <div v-if="detail.members.length" class="members">
        <strong>项目成员</strong>
        <span v-for="m in detail.members" :key="m.name" class="member">
          {{ m.name }}<i v-if="m.role"> · {{ m.role }}</i>
        </span>
      </div>

      <article class="article">
        <BlockRenderer v-for="(b, i) in detail.blocks" :key="i" :block="b" />
      </article>

      <div class="detail-foot">
        <button class="back" @click="goBack">← 返回全部项目</button>
      </div>
    </template>
  </div>
</template>

<style scoped>
.detail {
  max-width: 1180px;
  padding-top: 32px;
  padding-bottom: 56px;
}
.back {
  border: none;
  background: transparent;
  color: var(--sub);
  font-size: 14px;
  padding: 4px 0;
}
.back:hover {
  color: var(--brand);
}
.name {
  margin: 16px 0 8px;
  font-family: var(--serif);
  font-size: 38px;
  line-height: 1.25;
  font-weight: 700;
}
.summary {
  margin: 0 0 16px;
  max-width: 640px;
  color: var(--sub);
  font-size: 15px;
}
.meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 24px;
}
.tag {
  font-size: 12px;
  color: var(--sub);
  background: var(--bg);
  border: 1px solid var(--hairline);
  border-radius: 6px;
  padding: 2px 8px;
}
.meta-item {
  color: var(--muted);
  font-size: 13px;
}
.cover {
  margin: 0 auto 24px;
  max-width: 100%;
  border-radius: var(--r-md);
  overflow: hidden;
  background: var(--bg);
}
.cover img {
  width: 100%;
  display: block;
}
.members {
  color: var(--sub);
  font-size: 14px;
  margin-bottom: 8px;
}
.members strong {
  color: var(--ink);
  margin-right: 10px;
}
.member {
  margin-right: 14px;
}
.member i {
  font-style: normal;
  color: var(--muted);
}
.article {
  max-width: 720px;
  margin: 8px auto 0;
}
.detail-foot {
  max-width: 720px;
  margin: 40px auto 0;
}
.loading .sk {
  background: var(--bg);
  border-radius: var(--r-md);
}
.loading .title {
  width: 320px;
  height: 40px;
  margin: 16px 0;
}
.loading .cover {
  height: 260px;
  aspect-ratio: auto;
}
.loading .line {
  height: 16px;
  margin: 12px 0;
}
.w50 {
  width: 50%;
}
.w80 {
  width: 80%;
}
.not-found {
  text-align: center;
  padding: 96px 0;
  color: var(--sub);
}
.not-found button {
  margin-top: 16px;
  border: 1px solid var(--brand);
  color: var(--brand);
  background: transparent;
  padding: 8px 20px;
  border-radius: 999px;
}
</style>
