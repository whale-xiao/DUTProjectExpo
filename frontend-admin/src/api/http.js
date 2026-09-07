import axios from 'axios'
import { ElMessage } from 'element-plus'

export const TOKEN_KEY = 'admin_token'

const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '',
  timeout: 10000,
})

// 请求：自动附 token
http.interceptors.request.use((config) => {
  const token = localStorage.getItem(TOKEN_KEY)
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

function logoutToLogin() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem('admin_user')
  if (!window.location.pathname.startsWith('/login')) {
    window.location.href = '/login'
  }
}

// 响应：解包 {code,message,data}；40100 清态回登录
http.interceptors.response.use(
  (res) => {
    const body = res.data
    if (body && typeof body === 'object' && 'code' in body) {
      if (body.code === 0) return body.data
      if (body.code === 40100) {
        ElMessage.error(body.message || '登录已过期')
        logoutToLogin()
        return Promise.reject(new Error('unauthorized'))
      }
      const err = new Error(body.message || '请求失败')
      err.code = body.code
      return Promise.reject(err)
    }
    return body
  },
  (error) => {
    let msg = '网络错误，请检查后端服务'
    if (error.response?.data?.message) msg = error.response.data.message
    else if (error.code === 'ECONNABORTED') msg = '请求超时'
    ElMessage.error(msg)
    return Promise.reject(new Error(msg))
  },
)

export default http
