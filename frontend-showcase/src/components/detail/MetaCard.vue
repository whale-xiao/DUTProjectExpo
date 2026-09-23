<script setup>
// 详情页资料卡：封面 + 演职员表（members[]）+ 作品信息（见 docs/12 §5.2/§5.4）。
// 桌面端吸顶跟随长正文；移动端降级为横向信息条（由下方 media query 处理）。
import { computed } from 'vue'

const props = defineProps({
  detail: { type: Object, required: true },
})

const members = computed(() => props.detail?.members || [])
const tags = computed(() => props.detail?.tags || [])

function fmtViews(n) {
  return n >= 1000 ? `${(n / 1000).toFixed(1)}k` : String(n)
}
</script>

<template>
  <aside class="mcard">
    <div v-if="detail.coverUrl" class="poster">
      <img :src="detail.coverUrl" :alt="detail.name" loading="lazy" decoding="async" />
    </div>

    <section v-if="members.length" class="block">
      <h2 class="bt">项目成员</h2>
      <ul class="cast">
        <li v-for="m in members" :key="m.name">
          <span class="cast-name">{{ m.name }}</span>
          <span v-if="m.role" class="cast-role">{{ m.role }}</span>
        </li>
      </ul>
    </section>

    <section class="block">
      <h2 class="bt">作品信息</h2>
      <dl class="facts">
        <template v-if="detail.category">
          <dt>分类</dt>
          <dd>{{ detail.category.name }}</dd>
        </template>
        <template v-if="tags.length">
          <dt>标签</dt>
          <dd class="facts-tags">
            <i v-for="t in tags" :key="t">{{ t }}</i>
          </dd>
        </template>
        <dt>浏览</dt>
        <dd>{{ fmtViews(detail.viewCount) }}</dd>
        <template v-if="detail.publishedAt">
          <dt>收录</dt>
          <dd>{{ detail.publishedAt }}</dd>
        </template>
      </dl>
    </section>
  </aside>
</template>

<style scoped>
.mcard {
  position: sticky;
  top: 24px;
  align-self: start;
}
.poster {
  border: 1px solid var(--hairline);
  border-radius: var(--r-md);
  overflow: hidden;
  background: var(--bg);
  margin-bottom: 22px;
}
.poster img {
  width: 100%;
  aspect-ratio: 3 / 2;
  object-fit: cover;
  display: block;
}
.block + .block {
  margin-top: 22px;
}
.bt {
  margin: 0 0 10px;
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.14em;
  color: var(--muted);
}

/* 演职员表 */
.cast {
  list-style: none;
  margin: 0;
  padding: 0;
}
.cast li {
  display: grid;
  grid-template-columns: 1fr auto;
  align-items: baseline;
  gap: 10px;
  padding: 6px 0;
  border-bottom: 1px solid var(--hairline);
}
.cast li:last-child {
  border-bottom: none;
}
.cast-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--ink);
}
.cast-role {
  font-size: 12px;
  color: var(--muted);
}

/* 作品信息 */
.facts {
  margin: 0;
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 8px 14px;
  align-items: baseline;
}
.facts dt {
  font-size: 13px;
  color: var(--muted);
  white-space: nowrap;
}
.facts dd {
  margin: 0;
  font-size: 13px;
  color: var(--sub);
  min-width: 0;
}
.facts-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}
.facts-tags i {
  font-style: normal;
  font-size: 11px;
  color: var(--sub);
  background: var(--bg);
  border: 1px solid var(--hairline);
  padding: 1px 6px;
  border-radius: 6px;
}

/* —— 移动端：取消吸顶，降级为横向信息条 —— */
@media (max-width: 1023px) {
  .mcard {
    position: static;
  }
  .poster {
    display: none;
  }
  .mcard {
    border-top: 1px solid var(--hairline);
    border-bottom: 1px solid var(--hairline);
    padding: 18px 0;
    margin-bottom: 8px;
  }
  .block + .block {
    margin-top: 16px;
  }
  .cast {
    display: flex;
    flex-wrap: wrap;
    gap: 6px 20px;
  }
  .cast li {
    display: inline-flex;
    gap: 6px;
    padding: 0;
    border-bottom: none;
  }
}
</style>
