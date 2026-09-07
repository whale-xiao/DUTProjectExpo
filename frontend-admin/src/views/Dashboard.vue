<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getOverview } from '@/api/admin'
import { STATUS } from '@/utils/constants'

const router = useRouter()
const data = ref({ projectTotal: 0, published: 0, draft: 0, recommended: 0, viewTotal: 0, recent: [] })
const loading = ref(true)

const stats = ref([])

function renderStats() {
  const d = data.value
  stats.value = [
    { label: '项目总数', value: d.projectTotal },
    { label: '已上架', value: d.published },
    { label: '草稿', value: d.draft },
    { label: '推荐项目', value: d.recommended },
    { label: '累计浏览量', value: d.viewTotal.toLocaleString() },
  ]
}

onMounted(async () => {
  try {
    data.value = await getOverview()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
    renderStats()
  }
})
</script>

<template>
  <div>
    <div class="stat-row" v-loading="loading">
      <div v-for="s in stats" :key="s.label" class="stat-card">
        <div class="stat-value">{{ s.value }}</div>
        <div class="stat-label">{{ s.label }}</div>
      </div>
    </div>

    <div class="quick">
      <el-button type="primary" @click="router.push('/projects/create')">+ 新建项目</el-button>
      <el-button @click="router.push('/projects')">项目管理</el-button>
    </div>

    <el-card shadow="never" class="recent-card">
      <template #header>
        <span>最近编辑</span>
      </template>
      <el-table :data="data.recent" v-loading="loading" size="small">
        <el-table-column prop="name" label="项目名称" min-width="180" />
        <el-table-column prop="categoryName" label="分类" width="110" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="STATUS[row.status]?.type" size="small">{{ STATUS[row.status]?.label }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="viewCount" label="浏览" width="80" />
        <el-table-column prop="updatedAt" label="更新时间" width="170">
          <template #default="{ row }">{{ row.updatedAt?.replace('T', ' ')?.slice(0, 16) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="80">
          <template #default="{ row }">
            <el-button link type="primary" @click="router.push(`/projects/${row.id}/edit`)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<style scoped>
.stat-row {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 16px;
}
.stat-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  text-align: center;
}
.stat-value {
  font-size: 26px;
  font-weight: 700;
  color: #0d4a8e;
}
.stat-label {
  margin-top: 6px;
  color: #5a6478;
  font-size: 13px;
}
.quick {
  margin: 16px 0;
}
.recent-card {
  border-radius: 8px;
}
@media (max-width: 1000px) {
  .stat-row {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
