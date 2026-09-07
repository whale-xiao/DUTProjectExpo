<script setup>
// 单图上传组件：v-model 绑定 url；封装 uploadImage 请求。
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Delete, Picture } from '@element-plus/icons-vue'
import { uploadImage } from '@/api/admin'

const props = defineProps({
  modelValue: { type: String, default: '' },
})
const emit = defineEmits(['update:modelValue'])

const uploading = ref(false)

async function doUpload(option) {
  uploading.value = true
  try {
    const data = await uploadImage(option.file)
    emit('update:modelValue', data.url)
    ElMessage.success('上传成功')
  } catch (e) {
    ElMessage.error(e.message)
  } finally {
    uploading.value = false
  }
}

function remove() {
  emit('update:modelValue', '')
}
</script>

<template>
  <div class="img-upload">
    <template v-if="modelValue">
      <div class="preview">
        <img :src="modelValue" alt="图片预览" />
        <div class="ops">
          <el-button link type="primary" :loading="uploading" @click="$refs.fileInput?.click()">
            <el-icon><Picture /></el-icon> 更换
          </el-button>
          <el-button link type="danger" @click="remove">
            <el-icon><Delete /></el-icon> 移除
          </el-button>
        </div>
      </div>
    </template>
    <div v-else class="empty" @click="$refs.fileInput?.click()">
      <el-icon><Plus /></el-icon>
      <span>点击上传（jpg/png/webp ≤5MB）</span>
    </div>
    <input
      ref="fileInput"
      type="file"
      accept=".jpg,.jpeg,.png,.webp,.gif"
      style="display: none"
      @change="(e) => e.target.files?.[0] && doUpload({ file: e.target.files[0] })"
    />
  </div>
</template>

<style scoped>
.img-upload {
  display: inline-block;
}
.preview img {
  width: 180px;
  height: 120px;
  object-fit: cover;
  border-radius: 6px;
  display: block;
  border: 1px solid #e6e9f0;
}
.ops {
  margin-top: 4px;
}
.empty {
  width: 180px;
  height: 120px;
  border: 1px dashed #c9d1dc;
  border-radius: 6px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #a7b0c0;
  cursor: pointer;
  font-size: 12px;
  gap: 4px;
}
.empty:hover {
  border-color: #0d4a8e;
  color: #0d4a8e;
}
</style>
