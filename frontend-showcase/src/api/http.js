import axios from 'axios'

// 开发走 Vite 代理(/api)，生产由 Nginx 同源反代，因此默认用相对路径；
// 若部署到独立域名后端，可在 .env 中设置 VITE_API_BASE_URL。
const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '',
  timeout: 8000,
})

// 拦截器统一解包 {code,message,data}：业务代码只见 data。
http.interceptors.response.use(
  (res) => {
    const body = res.data
    if (body && typeof body === 'object' && 'code' in body) {
      if (body.code === 0) return body.data
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
    const e = new Error(msg)
    e.code = error.response?.status || -1
    return Promise.reject(e)
  },
)

export default http
