<script setup>
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { login } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const formRef = ref()
const loading = ref(false)

const form = reactive({ username: 'admin', password: '' })
const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

async function submit() {
  await formRef.value.validate().catch(() => Promise.reject())
  loading.value = true
  try {
    const data = await login(form.username, form.password)
    auth.setLogin(data)
    ElMessage.success('登录成功')
    router.push(route.query.redirect || '/dashboard')
  } catch (e) {
    // 错误提示已在拦截器处理
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <div class="card">
      <div class="title">
        <div class="logo-dot"></div>
        <h2>优秀项目展示系统</h2>
        <p>内容管理后台</p>
      </div>
      <el-form ref="formRef" :model="form" :rules="rules" size="large" @keyup.enter="submit">
        <el-form-item prop="username">
          <el-input v-model="form.username" placeholder="用户名" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="form.password" type="password" placeholder="密码" show-password />
        </el-form-item>
        <el-button type="primary" class="submit" :loading="loading" @click="submit">
          登 录
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #08345f 0%, #0d4a8e 100%);
}
.card {
  width: 380px;
  background: #fff;
  border-radius: 10px;
  padding: 40px 36px 32px;
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.18);
}
.title {
  text-align: center;
  margin-bottom: 24px;
}
.logo-dot {
  width: 40px;
  height: 40px;
  margin: 0 auto 10px;
  border-radius: 10px;
  background: #0d4a8e;
}
.title h2 {
  margin: 0;
  font-size: 20px;
  color: #1e2433;
}
.title p {
  margin: 4px 0 0;
  color: #a7b0c0;
  font-size: 13px;
}
.submit {
  width: 100%;
  margin-top: 4px;
}
</style>
