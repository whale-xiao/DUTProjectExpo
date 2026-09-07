<script setup>
// 内容块编辑器：v-model 绑定 blocks 数组；按 heading/paragraph/image 分发编辑；
// 支持增删、上移下移（顺序即保存时的 sort_order）。
import { ArrowUp, ArrowDown, Delete, DocumentAdd } from '@element-plus/icons-vue'
import ImageUpload from './ImageUpload.vue'
import { BLOCK_TYPES } from '@/utils/constants'

const props = defineProps({
  modelValue: { type: Array, default: () => [] },
})
const emit = defineEmits(['update:modelValue'])

function update(list) {
  emit('update:modelValue', list)
}

function addBlock(type) {
  const block = type === 'image' ? { type, url: '', caption: '' } : { type, text: '' }
  update([...props.modelValue, block])
}

function removeAt(i) {
  const list = props.modelValue.slice()
  list.splice(i, 1)
  update(list)
}

function move(i, dir) {
  const list = props.modelValue.slice()
  const j = i + dir
  if (j < 0 || j >= list.length) return
  ;[list[i], list[j]] = [list[j], list[i]]
  update(list)
}

function setField(i, key, val) {
  const list = props.modelValue.slice()
  list[i] = { ...list[i], [key]: val }
  update(list)
}
</script>

<template>
  <div class="block-editor">
    <div class="toolbar">
      <el-button size="small" type="primary" plain :icon="DocumentAdd" @click="addBlock('heading')">
        小节标题
      </el-button>
      <el-button size="small" type="primary" plain :icon="DocumentAdd" @click="addBlock('paragraph')">
        段落
      </el-button>
      <el-button size="small" type="primary" plain :icon="DocumentAdd" @click="addBlock('image')">
        图片
      </el-button>
      <span class="hint">按顺序即页面展示顺序，可用上移/下移调整</span>
    </div>

    <el-empty v-if="!modelValue.length" description="还没有内容，先添加一个小节/段落/图片" :image-size="80" />

    <div v-for="(b, i) in modelValue" :key="i" class="block-row">
      <div class="block-head">
        <el-tag size="small" effect="plain">{{ BLOCK_TYPES[b.type] || b.type }}</el-tag>
        <span class="idx">#{{ i + 1 }}</span>
        <span class="spacer"></span>
        <el-button-group size="small">
          <el-button :icon="ArrowUp" :disabled="i === 0" @click="move(i, -1)" />
          <el-button :icon="ArrowDown" :disabled="i === modelValue.length - 1" @click="move(i, 1)" />
        </el-button-group>
        <el-button size="small" type="danger" plain :icon="Delete" @click="removeAt(i)" />
      </div>

      <div class="block-body">
        <template v-if="b.type === 'heading'">
          <el-input
            :model-value="b.text"
            placeholder="小节标题，如「项目介绍」"
            @update:model-value="(v) => setField(i, 'text', v)"
          />
        </template>
        <template v-else-if="b.type === 'paragraph'">
          <el-input
            :model-value="b.text"
            type="textarea"
            :rows="3"
            placeholder="段落正文"
            @update:model-value="(v) => setField(i, 'text', v)"
          />
        </template>
        <template v-else>
          <div class="img-block">
            <ImageUpload :model-value="b.url" @update:model-value="(v) => setField(i, 'url', v)" />
            <el-input
              class="cap"
              :model-value="b.caption"
              placeholder="图片说明（可选）"
              @update:model-value="(v) => setField(i, 'caption', v)"
            />
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<style scoped>
.block-editor {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
}
.hint {
  color: #a7b0c0;
  font-size: 12px;
  margin-left: 8px;
}
.block-row {
  border: 1px solid #e6e9f0;
  border-radius: 8px;
  overflow: hidden;
}
.block-head {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  background: #f7f8fb;
  border-bottom: 1px solid #e6e9f0;
}
.idx {
  color: #a7b0c0;
  font-size: 12px;
}
.spacer {
  flex: 1;
}
.block-body {
  padding: 12px;
}
.img-block {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}
.cap {
  max-width: 260px;
}
</style>
