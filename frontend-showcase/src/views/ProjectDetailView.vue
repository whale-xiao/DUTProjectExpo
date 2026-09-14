<script setup>
// 项目详情：作品资料页——通栏剧照头图 + 「左资料卡 / 右正文」两栏（见 docs/12 §5）。
// 正文仍完全由后端 blocks[] 驱动，前端不感知正文结构。
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getProject } from '@/api/projects'
import BlockRenderer from '@/components/detail/BlockRenderer.vue'
import DetailHero from '@/components/detail/DetailHero.vue'
import MetaCard from '@/components/detail/MetaCard.vue'

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
</script>

<template>
  <!-- 加载中：头图占位高度与成功态一致，避免 CLS -->
  <div v-if="loading" class="dhero-sk">
    <div class="container">
      <div class="sk sk-crumb"></div>
      <div class="sk sk-title"></div>
      <div class="sk sk-line"></div>
    </div>
  </div>

  <div v-else-if="error && !detail" class="container not-found">
    <p>这篇作品不存在或已下架</p>
    <button @click="router.push('/projects')">返回全部项目</button>
  </div>

  <template v-else-if="detail">
    <DetailHero :detail="detail" @back="goBack" />

    <div class="container body-grid">
      <MetaCard :detail="detail" />

      <!-- 正文逐块浮现。注意：这里用外层 div 承载 v-reveal——
           ① BlockRenderer 根部是 v-if/v-else 链，指令挂组件上可能命中非元素根；
           ② 也绝不能把 v-reveal 加到 .body-grid 或 MetaCard 的祖先上——
              transform 会让 position: sticky 失效。 -->
      <article class="article">
        <template v-if="detail.blocks?.length">
          <div v-for="(b, i) in detail.blocks" :key="i" v-reveal="Math.min(i, 4) * 60">
            <BlockRenderer :block="b" />
          </div>
        </template>
        <p v-else class="no-body">暂无正文</p>
      </article>
    </div>

    <div class="container detail-foot">
      <button class="back" type="button" @click="goBack">← 返回全部项目</button>
    </div>
  </template>
</template>

<style scoped>
.dhero-sk {
  height: var(--hero-h);
  display: flex;
  align-items: flex-end;
  padding-bottom: 42px;
  background: linear-gradient(135deg, #16202e 0%, #22303f 100%);
}
.dhero-sk .sk {
  border-radius: var(--r-sm);
  background: rgba(255, 255, 255, 0.16);
}
.sk-crumb {
  width: 140px;
  height: 14px;
  margin-bottom: 12px;
}
.sk-title {
  width: min(420px, 60%);
  height: 44px;
  margin-bottom: 14px;
}
.sk-line {
  width: min(560px, 80%);
  height: 15px;
}

/* 两栏：资料卡吸顶 + 正文列（正文限宽 720px 保证阅读行宽） */
.body-grid {
  display: grid;
  grid-template-columns: 260px minmax(0, 1fr);
  gap: 48px;
  max-width: 1060px;
  padding-top: 36px;
}
.article {
  max-width: 720px;
}
.no-body {
  margin: 0;
  color: var(--muted);
  font-size: 14px;
}

.detail-foot {
  max-width: 1060px;
  margin-top: 48px;
  padding-top: 20px;
  border-top: 1px solid var(--hairline);
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

@media (max-width: 1023px) {
  .body-grid {
    grid-template-columns: 1fr;
    gap: 22px;
    padding-top: 22px;
  }
  .detail-foot {
    margin-top: 32px;
  }
}
</style>
