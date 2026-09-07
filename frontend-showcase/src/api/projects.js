import http from './http'

/**
 * 项目列表（全部 / 搜索 / 分类标签筛选 / 排序 / 分页）。
 * @param {{ keyword?: string, categoryId?: number, tagId?: number, sort?: string, page?: number, pageSize?: number }} params
 */
export function listProjects(params = {}) {
  return http.get('/api/projects', {
    params: {
      page: 1,
      page_size: 20,
      sort: 'recommended',
      ...params,
    },
  })
}

/** 搜索联想（回车可到搜索结果页）。signal 用于中断过期请求。 */
export function suggestProjects(keyword, limit = 6, signal) {
  return http.get('/api/search/suggest', { params: { keyword, limit }, signal })
}

/** 项目详情（作品文章，含 blocks/members/tags/category）。 */
export function getProject(id) {
  return http.get(`/api/projects/${id}`)
}
