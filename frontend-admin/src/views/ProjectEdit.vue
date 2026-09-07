<script setup>
// 项目编辑页：上段基础信息 + 下段内容块编排；「保存」整单提交（事务内全量替换）。
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Plus, Delete, View } from '@element-plus/icons-vue'
import ImageUpload from '@/components/ImageUpload.vue'
import BlockEditor from '@/components/BlockEditor.vue'
import {
  listCategories,
  listTags,
  getProjectDetail,
  createProject,
  updateProject,
} from '@/api/admin'
import { STATUS_OPTIONS } from '@/utils/constants'

const route = useRoute()
const router = useRouter()

const id = computed(() => {
  const v = Number(route.params.id)
  return Number.isInteger(v) && v > 0 ? v : null
})

const formRef = ref()
const loading = ref(false)
const saving = ref(false)
const categories = ref([])
const tagOptions = ref([])

const form = reactive({
  name: '',
  summary: '',
  categoryId: null,
  coverUrl: '',
  tags: [],
  status: 'draft',
  isRecommended: false,
  recommendOrder: 0,
  members: [],
  blocks: [],
})

const rules = {
  name: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
  summary: [{ required: true, message: '请输入一句话介绍', trigger: 'blur' }],
}

onMounted(async () => {
  loading.value = true
  try {
    const [cats, tags] = await Promise.all([listCategories(), listTags()])
    categories.value = cats
    tagOptions.value = tags.map((t) => t.name)
    if (id.value) {
      const d = await getProjectDetail(id.value)
      form.name = d.name
      form.summary = d.summary
      form.categoryId = d.category?.id ?? null
      form.coverUrl = d.coverUrl
      form.tags = d.tags || []
      form.status = d.status || 'draft'
      form.isRecommended = d.isRecommended
      form.recommendOrder = d.recommendOrder
      form.members = (d.members || []).map((m) => ({ name: m.name, role: m.role }))
      form.blocks = (d.blocks || []).map((b) => ({ ...b }))
    } else {
      form.members = [{ name: '', role: '' }]
    }
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
})

function addMember() {
  form.members.push({ name: '', role: '' })
}
function removeMember(i) {
  form.members.splice(i, 1)
}

async function save() {
  await formRef.value.validate().catch(() => Promise.reject())
  if (!form.categoryId) {
    ElMessage.warning('请选择项目分类')
    return
  }
  const payload = {
    name: form.name.trim(),
    summary: form.summary.trim(),
    categoryId: form.categoryId,
    coverUrl: form.coverUrl,
    status: form.status,
    isRecommended: form.isRecommended,
    recommendOrder: form.isRecommended ? form.recommendOrder : 0,
    tags: form.tags,
    members: form.members.filter((m) => m.name?.trim()),
    blocks: form.blocks.filter((b) => {
      if (b.type === 'image') return true
      return String(b.text || '').trim() !== ''
    }),
  }
  saving.value = true
  try {
    if (id.value) {
      await updateProject(id.value, payload)
      ElMessage.success('保存成功')
    } else {
      await createProject(payload)
      ElMessage.success('创建成功')
    }
    router.push('/projects')
  } catch (e) {
    // 拦截器已提示
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div v-loading="loading">
    <div class="page" v-if="!loading">
      <div class="head">
        <el-button :icon="ArrowLeft" text @click="router.push('/projects')">返回</el-button>
        <h2 class="title">{{ id ? '编辑项目' : '新建项目' }}</h2>
        <span class="spacer"></span>
        <el-button :icon="View" @click="router.push('/')" v-if="false" />
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </div>

      <el-card shadow="never" class="card">
        <template #header><span>基础信息</span></template>
        <el-form ref="formRef" :model="form" :rules="rules" label-width="110px">
          <el-form-item label="项目名称" prop="name">
            <el-input v-model="form.name" maxlength="100" show-word-limit placeholder="如 MyGoVoice" />
          </el-form-item>
          <el-form-item label="一句话介绍" prop="summary">
            <el-input
              v-model="form.summary"
              type="textarea"
              :rows="2"
              maxlength="255"
              show-word-limit
              placeholder="卡片上展示的一句话，说明项目是什么"
            />
          </el-form-item>
          <el-form-item label="封面图片">
            <ImageUpload v-model="form.coverUrl" />
          </el-form-item>
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="分类" required>
                <el-select v-model="form.categoryId" placeholder="选择分类" style="width: 100%">
                  <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="技术标签">
                <el-select
                  v-model="form.tags"
                  multiple
                  filterable
                  allow-create
                  default-first-option
                  :reserve-keyword="false"
                  placeholder="输入或选择标签"
                  style="width: 100%"
                >
                  <el-option v-for="t in tagOptions" :key="t" :label="t" :value="t" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="状态">
                <el-radio-group v-model="form.status">
                  <el-radio-button v-for="s in STATUS_OPTIONS" :key="s.value" :value="s.value">
                    {{ s.label }}
                  </el-radio-button>
                </el-radio-group>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="推荐位">
                <el-switch v-model="form.isRecommended" />
                <el-input-number
                  v-model="form.recommendOrder"
                  :min="1"
                  :max="99"
                  :disabled="!form.isRecommended"
                  style="margin-left: 12px"
                  size="small"
                />
                <span class="tip">数值越小越靠前</span>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="项目成员">
            <div class="members">
              <div v-for="(m, i) in form.members" :key="i" class="member-row">
                <el-input v-model="m.name" placeholder="姓名" style="width: 200px" />
                <el-input v-model="m.role" placeholder="角色 / 职责" style="width: 240px" />
                <el-button :icon="Delete" circle text type="danger" :disabled="form.members.length <= 1" @click="removeMember(i)" />
              </div>
              <el-button size="small" :icon="Plus" @click="addMember">添加成员</el-button>
            </div>
          </el-form-item>
        </el-form>
      </el-card>

      <el-card shadow="never" class="card">
        <template #header><span>正文编排（作品文章）</span></template>
        <BlockEditor v-model="form.blocks" />
      </el-card>

      <div class="foot">
        <el-button @click="router.push('/projects')">取消</el-button>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.head {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}
.title {
  margin: 0;
  font-size: 18px;
}
.spacer {
  flex: 1;
}
.card {
  margin-bottom: 16px;
  border-radius: 8px;
}
.members {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.member-row {
  display: flex;
  gap: 8px;
  align-items: center;
}
.tip {
  color: #a7b0c0;
  font-size: 12px;
  margin-left: 8px;
}
.foot {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>
