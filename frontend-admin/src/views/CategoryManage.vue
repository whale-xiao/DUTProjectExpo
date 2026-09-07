<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import {
  listCategories,
  createCategory,
  updateCategory,
  deleteCategory,
} from '@/api/admin'

const loading = ref(false)
const list = ref([])
const keyword = ref('')

const dialog = ref(false)
const saving = ref(false)
const editingId = ref(null)
const form = ref({ name: '', description: '', sortOrder: 1 })
const rules = {
  name: [{ required: true, message: '请输入分类名', trigger: 'blur' }],
}

const filtered = () =>
  keyword.value
    ? list.value.filter((c) => c.name.includes(keyword.value) || (c.description || '').includes(keyword.value))
    : list.value

async function fetchList() {
  loading.value = true
  try {
    list.value = await listCategories()
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    loading.value = false
  }
}
onMounted(fetchList)

function openCreate() {
  editingId.value = null
  form.value = { name: '', description: '', sortOrder: list.value.length + 1 }
  dialog.value = true
}
function openEdit(row) {
  editingId.value = row.id
  form.value = { name: row.name, description: row.description, sortOrder: row.sortOrder }
  dialog.value = true
}

async function save() {
  saving.value = true
  try {
    if (editingId.value) {
      await updateCategory(editingId.value, form.value)
      ElMessage.success('已更新')
    } else {
      await createCategory(form.value)
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
    await ElMessageBox.confirm(`删除分类「${row.name}」？`, '确认删除', { type: 'warning' })
  } catch {
    return
  }
  try {
    await deleteCategory(row.id)
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
      <el-input
        v-model="keyword"
        placeholder="筛选分类"
        clearable
        style="width: 220px"
      >
        <template #prefix><el-icon><Search /></el-icon></template>
      </el-input>
      <span class="spacer"></span>
      <el-button type="primary" :icon="Plus" @click="openCreate">新增分类</el-button>
    </div>

    <el-table :data="filtered()" v-loading="loading" stripe>
      <el-table-column prop="name" label="分类名" width="160" />
      <el-table-column prop="description" label="描述" min-width="220" />
      <el-table-column prop="sortOrder" label="排序" width="90" />
      <el-table-column prop="createdAt" label="创建时间" width="180">
        <template #default="{ row }">{{ row.createdAt?.replace('T', ' ')?.slice(0, 16) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="140">
        <template #default="{ row }">
          <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button link type="danger" @click="onDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialog" :title="editingId ? '编辑分类' : '新增分类'" width="480px">
      <el-form :model="form" :rules="rules" label-width="80px">
        <el-form-item label="分类名" prop="name">
          <el-input v-model="form.name" maxlength="50" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" maxlength="255" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sortOrder" :min="1" />
          <span class="tip">越小越靠前</span>
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
.tip {
  color: #a7b0c0;
  font-size: 12px;
  margin-left: 8px;
}
</style>
