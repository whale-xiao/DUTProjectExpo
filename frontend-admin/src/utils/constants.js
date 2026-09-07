// 状态 → 展示标签映射（Element Plus el-tag type）
export const STATUS = {
  draft: { label: '草稿', type: 'info' },
  published: { label: '已上架', type: 'success' },
  archived: { label: '已下架', type: 'warning' },
}
export const STATUS_OPTIONS = [
  { value: 'draft', label: '草稿' },
  { value: 'published', label: '上架' },
  { value: 'archived', label: '下架' },
]

// 内容块类型
export const BLOCK_TYPES = {
  heading: '小节标题',
  paragraph: '段落',
  image: '图片',
}
