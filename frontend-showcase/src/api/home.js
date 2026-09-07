import http from './http'

// GET /api/home 首页聚合：推荐 + 分类 + 最新
export function getHome() {
  return http.get('/api/home')
}
