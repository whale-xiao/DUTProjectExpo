import http from './http'

// GET /api/categories 分类列表
export function getCategories() {
  return http.get('/api/categories')
}

// GET /api/tags 标签列表
export function getTags() {
  return http.get('/api/tags')
}
