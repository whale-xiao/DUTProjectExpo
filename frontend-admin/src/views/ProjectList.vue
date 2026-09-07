<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Star } from '@element-plus/icons-vue'
import {
  listProjects,
  listCategories,
  deleteProject,
  setProjectStatus,
  setProjectRecommend,
} from '@/api/admin'
import { STATUS, STATUS_OPTIONS } from '@/utils/constants'

const router = useRouter()

const loading = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const query = ref({ keyword: '', categoryId: '', status: '' })
const categories = ref([])

async function fetchList() {
  loading.value = true
  try {
    const data = await listProjects({
      keyword: query.value.keyword || undefined,
      category_id: query.value.categoryId || undefined,
      status: query.value.status || undefined,
      page: page.value,
      page_size: pageSize.value,
    })
    list.value = data.list
    total.value = data.total
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    categories.value = await listCategories()
  } catch { /* 可选 */ }
  fetchList()
})

function search() {
  page.value = 1
  fetchList()
}

async function toggleStatus(row) {
  const next = row.status === 'published' ? 'archived' : 'published'
  await setProjectStatus(row.id, next)
  ElMessage.success(next === 'published' ? '已上架' : '已下架')
  fetchList()
}

async function toggleRecommend(row) {
  const next = !row.isRecommended
  await setProjectRecommend(row.id, next, next ? 0 : 0)
  ElMessage.success(next ? '已设为推荐' : '已取消推荐')
  row.isRecommended = next
}

async function onDelete(row) {
  await ElMessageBox.confirm(`删除项目「${row.name}」？其正文内容将一并删除且不可恢复。`, '确认删除', {
    type: 'warning',
  }).catch(() => Promise.reject())
  await deleteProject(row.id)
  ElMessage.success('已删除')
  if (list.value.length === 1 && page.value > 1) page.value -= 1
  fetchList()
}
</script>

<template>
  <div class="page">
    <div class="page-toolbar">
      <el-input
        v-model="query.keyword"
        placeholder="搜索项目名称 / 简介 / 标签"
        clearable
        style="width: 260px"
        @keyup.enter="search"
        @clear="search"
      >
        <template #prefix><el-icon><Search /></el-icon></template>
      </el-input>
      <el-select v-model="query.categoryId" placeholder="分类" clearable style="width: 140px" @change="search">
        <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
      </el-select>
      <el-select v-model="query.status" placeholder="状态" clearable style="width: 120px" @change="search">
        <el-option v-for="s in STATUS_OPTIONS" :key="s.value" :label="s.label" :value="s.value" />
      </el-select>
      <span class="spacer"></span>
      <el-button type="primary" :icon="Plus" @click="router.push('/projects/create')">新建项目</el-button>
    </div>

    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column label="项目" min-width="240">
        <template #default="{ row }">
          <div class="proj">
            <el-image v-if="row.coverUrl" :src="row.coverUrl" fit="cover" class="thumb" />
            <div v-else class="thumb thumb-ph">{{ row.name.slice(0, 1) }}</div>
            <div class="meta">
              <div class="name">{{ row.name }}</div>
              <div class="summary">{{ row.summary }}</div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="categoryName" label="分类" width="110" />
      <el-table-column label="技术标签" min-width="150">
        <template #default="{ row }">
          <el-tag v-for="t in (row.tags || []).slice(0, 3)" :key="t" size="small" effect="plain" class="tg">
            {{ t }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="STATUS[row.status]?.type" size="small">{{ STATUS[row.status]?.label }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="推荐" width="60">
        <template #default="{ row }">
          <el-switch
            :model-value="row.isRecommended"
            :active-action-icon="Star"
            :inactive-action-icon="Star"
            inline-prompt
            style="--el-switch-on-color: #e0700a"
            @change="toggleRecommend(row)"
          />
        </template>
      </el-table-column>
      <el-table-column prop="viewCount" label="浏览" width="80" />
      <el-table-column label="更新时间" width="160">
        <template #default="{ row }">{{ row.updatedAt?.replace('T', ' ')?.slice(0, 16) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="router.push(`/projects/${row.id}/edit`)">编辑</el-button>
          <el-button link :type="row.status === 'published' ? 'warning' : 'success'" @click="toggleStatus(row)">
            {{ row.status === 'published' ? '下架' : '上架' }}
          </el-button>
          <el-button link type="danger" @click="onDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pager">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="fetchList"
      />
    </div>
  </div>
</template>

<style scoped>
.spacer {
  flex: 1;
}
.proj {
  display: flex;
  align-items: center;
  gap: 10px;
}
.thumb {
  width: 56px;
  height: 40px;
  border-radius: 4px;
  flex: none;
}
.thumb-ph {
  background: #f0f2f5;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #a7b0c0;
  font-weight: 700;
}
.name {
  font-weight: 600;
}
.summary {
  color: #a7b0c0;
  font-size: 12px;
  max-width: 320px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.tg {
  margin-right: 4px;
}
</style>
