import http from './http'

export function login(username, password) {
  return http.post('/api/auth/login', { username, password })
}

export function fetchMe() {
  return http.get('/api/auth/me')
}
