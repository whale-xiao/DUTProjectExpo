import http from './http'

// —— 仪表盘 ——
export function getOverview() {
  return http.get('/api/admin/overview')
}

// —— 项目管理 ——
export function listProjects(params = {}) {
  return http.get('/api/admin/projects', {
    params: { page: 1, page_size: 20, ...params },
  })
}
export function getProjectDetail(id) {
  return http.get(`/api/admin/projects/${id}`)
}
export function createProject(data) {
  return http.post('/api/admin/projects', data)
}
export function updateProject(id, data) {
  return http.put(`/api/admin/projects/${id}`, data)
}
export function deleteProject(id) {
  return http.delete(`/api/admin/projects/${id}`)
}
export function setProjectStatus(id, status) {
  return http.put(`/api/admin/projects/${id}/status`, { status })
}
export function setProjectRecommend(id, isRecommended, recommendOrder) {
  return http.put(`/api/admin/projects/${id}/recommend`, { isRecommended, recommendOrder })
}

// —— 分类管理 ——
export const listCategories = () => http.get('/api/admin/categories')
export const createCategory = (data) => http.post('/api/admin/categories', data)
export const updateCategory = (id, data) => http.put(`/api/admin/categories/${id}`, data)
export const deleteCategory = (id) => http.delete(`/api/admin/categories/${id}`)

// —— 标签管理 ——
export const listTags = () => http.get('/api/admin/tags')
export const createTag = (data) => http.post('/api/admin/tags', data)
export const updateTag = (id, data) => http.put(`/api/admin/tags/${id}`, data)
export const deleteTag = (id) => http.delete(`/api/admin/tags/${id}`)

// —— 上传 ——
export function uploadImage(file) {
  const form = new FormData()
  form.append('file', file)
  return http.post('/api/admin/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}
