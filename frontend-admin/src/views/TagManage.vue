<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import { listTags, createTag, updateTag, deleteTag } from '@/api/admin'

const loading = ref(false)
const list = ref([])
const keyword = ref('')

const dialog = ref(false)
const saving = ref(false)
const editingId = ref(null)
const form = ref({ name: '', color: '#0d4a8e' })
const rules = {
  name: [{ required: true, message: '请输入标签名', trigger: 'blur' }],
}

const filtered = () => (keyword.value ? list.value.filter((t) => t.name.includes(keyword.value)) : list.value)

async function fetchList() {
  loading.value = true
  try {
    list.value = await listTags()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}
onMounted(fetchList)

function openCreate() {
  editingId.value = null
  form.value = { name: '', color: '#0d4a8e' }
  dialog.value = true
}
function openEdit(row) {
  editingId.value = row.id
  form.value = { name: row.name, color: row.color || '#0d4a8e' }
  dialog.value = true
}

async function save() {
  saving.value = true
  try {
    if (editingId.value) {
      await updateTag(editingId.value, form.value)
      ElMessage.success('已更新')
    } else {
      await createTag(form.value)
      ElMessage.success('创建成功')
    }
    dialog.value = false
    fetchList()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    saving.value = false
  }
}

async function onDelete(row) {
  try {
    await ElMessageBox.confirm(`删除标签「${row.name}」？若仍被项目引用将无法删除。`, '确认删除', { type: 'warning' })
  } catch {
    return
  }
  try {
    await deleteTag(row.id)
    ElMessage.success('已删除')
    fetchList()
  } catch (e) {
    ElMessage.error(e.message)
  }
}
</script>

<template>
  <div class="page">
    <div class="page-toolbar">
      <el-input v-model="keyword" placeholder="筛选标签" clearable style="width: 220px">
        <template #prefix><el-icon><Search /></el-icon></template>
      </el-input>
      <span class="spacer"></span>
      <el-button type="primary" :icon="Plus" @click="openCreate">新增标签</el-button>
    </div>

    <el-table :data="filtered()" v-loading="loading" stripe>
      <el-table-column prop="name" label="标签名" width="200">
        <template #default="{ row }">
          <el-tag :color="row.color || '#0d4a8e'" effect="dark">{{ row.name }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="color" label="颜色" width="140">
        <template #default="{ row }"><span :style="{ color: row.color }">{{ row.color }}</span></template>
      </el-table-column>
      <el-table-column prop="createdAt" label="创建时间" width="200">
        <template #default="{ row }">{{ row.createdAt?.replace('T', ' ')?.slice(0, 16) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="140">
        <template #default="{ row }">
          <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button link type="danger" @click="onDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialog" :title="editingId ? '编辑标签' : '新增标签'" width="440px">
      <el-form :model="form" :rules="rules" label-width="70px">
        <el-form-item label="标签名" prop="name">
          <el-input v-model="form.name" maxlength="30" />
        </el-form-item>
        <el-form-item label="颜色">
          <el-color-picker v-model="form.color" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.spacer {
  flex: 1;
}
</style>
