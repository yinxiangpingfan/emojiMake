<template>
  <div class="register-page">
    <div class="register-container">
      <div class="register-box">
        <div class="register-header">
          <div class="logo">
            <el-icon class="logo-icon"><UserFilled /></el-icon>
          </div>
          <h2 class="register-title">创建账户</h2>
          <p class="register-subtitle">加入AI视频趣味表情生成社区</p>
        </div>
        
        <el-alert
          v-if="showApiError"
          title="无法连接到服务器"
          type="error"
          description="请确保后端服务正在运行，或联系系统管理员。"
          show-icon
          closable
          @close="showApiError = false"
          class="error-alert"
        />
        
        <el-form 
          :model="registerForm" 
          :rules="registerRules" 
          ref="registerFormRef"
          class="register-form"
        >
          <el-form-item prop="phone">
            <div class="input-wrapper">
              <el-icon class="input-icon"><User /></el-icon>
              <el-input 
                v-model="registerForm.phone" 
                placeholder="请输入手机号"
                class="register-input"
                maxlength="11"
              />
            </div>
          </el-form-item>
          
          <el-form-item prop="password">
            <div class="input-wrapper">
              <el-icon class="input-icon"><Lock /></el-icon>
              <el-input 
                v-model="registerForm.password" 
                type="password"
                placeholder="请输入密码"
                show-password
                class="register-input"
                maxlength="20"
              />
            </div>
          </el-form-item>
          
          <el-form-item prop="confirmPassword">
            <div class="input-wrapper">
              <el-icon class="input-icon"><Lock /></el-icon>
              <el-input 
                v-model="registerForm.confirmPassword" 
                type="password"
                placeholder="请确认密码"
                show-password
                class="register-input"
                maxlength="20"
              />
            </div>
          </el-form-item>
          
          <el-form-item>
            <el-button 
              type="primary" 
              size="large" 
              class="register-button"
              @click="handleRegister"
              :loading="loading"
              round
            >
              <el-icon v-if="!loading"><Right /></el-icon>
              注册
            </el-button>
          </el-form-item>
        </el-form>
        
        <div class="register-footer">
          <div class="login-link">
            <span>已有账号？</span>
            <el-button type="primary" link @click="goToLogin">
              <el-icon><Position /></el-icon>
              立即登录
            </el-button>
          </div>
          
          <el-button type="primary" link @click="goToHome" class="home-link">
            <el-icon><House /></el-icon>
            返回首页
          </el-button>
        </div>
      </div>
      
      <div class="register-decoration">
        <div class="decoration-content">
          <div class="decoration-title">
            <h3>加入我们的创意社区</h3>
            <p>释放你的创造力，制作独一无二的视频表情包</p>
          </div>
          <div class="features">
            <div class="feature-item">
              <el-icon><Picture /></el-icon>
              <span>AI视频生成</span>
            </div>
            <div class="feature-item">
              <el-icon><MagicStick /></el-icon>
              <span>智能生成</span>
            </div>
            <div class="feature-item">
              <el-icon><Download /></el-icon>
              <span>一键下载</span>
            </div>
          </div>
          <div class="bubble" v-for="i in 8" :key="i"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElNotification } from 'element-plus'
import type { FormInstance } from 'element-plus'
import axios from 'axios'
// 导入Element Plus图标
import {
  User, 
  Lock, 
  Right, 
  Position, 
  House, 
  Picture, 
  MagicStick, 
  Download,
  UserFilled
} from '@element-plus/icons-vue'

const router = useRouter()
const registerFormRef = ref<FormInstance>()
const loading = ref(false)
const showApiError = ref(false)

const registerForm = reactive({
  phone: '',
  password: '',
  confirmPassword: ''
})

const registerRules = {
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 8, message: '密码长度不能少于8位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: any) => {
        if (value !== registerForm.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const handleRegister = async () => {
  if (!registerFormRef.value) return
  
  await registerFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      showApiError.value = false
      
      try {
        // 调用注册API
        const formData = new FormData()
        formData.append('phone', registerForm.phone)
        formData.append('password', registerForm.password)
        
        const response = await axios.post('/api/v1/users/register', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        
        if (response.data.code === 0) {
          // 注册成功
          ElMessage.success('注册成功')
          router.push('/login')
        } else {
          // 注册失败
          ElMessage.error(response.data.message || '注册失败')
        }
      } catch (error: any) {
        console.error('Register error:', error)
        showApiError.value = true
        
        if (error.code === 'ERR_NETWORK') {
          ElNotification({
            title: '网络错误',
            message: '无法连接到服务器，请确保后端服务正在运行',
            type: 'error',
            duration: 0
          })
        } else {
          ElMessage.error(error.response?.data?.message || '注册请求失败，请稍后重试')
        }
      } finally {
        loading.value = false
      }
    }
  })
}

const goToLogin = () => {
  router.push('/login')
}

const goToHome = () => {
  router.push('/')
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  position: relative;
  overflow: hidden;
}

.register-container {
  position: relative;
  width: 100%;
  max-width: 1200px;
  min-height: 80vh;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 20px;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.2);
  display: flex;
  overflow: hidden;
}

.register-box {
  flex: 1;
  padding: 60px 80px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  z-index: 2;
}

.register-header {
  text-align: center;
  margin-bottom: 30px;
}

.logo {
  margin-bottom: 20px;
  display: flex;
  justify-content: center;
}

.logo-icon {
  font-size: 3rem;
  color: #667eea;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.register-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: #333;
  margin-bottom: 10px;
}

.register-subtitle {
  font-size: 1.1rem;
  color: #666;
}

.error-alert {
  margin-bottom: 20px;
}

.register-form {
  width: 100%;
}

.input-wrapper {
  position: relative;
  width: 100%;
}

.input-icon {
  position: absolute;
  left: 15px;
  top: 50%;
  transform: translateY(-50%);
  z-index: 1;
  color: #909399;
  font-size: 1.2rem;
}

.register-input {
  width: 100%;
}

.register-input :deep(.el-input__wrapper) {
  border-radius: 25px;
  padding: 0 20px 0 45px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  border: 1px solid #dcdfe6;
  height: 50px;
  width: 100%;
}

.register-input :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 1px #667eea;
}

.register-input :deep(.el-input__inner) {
  color: #333;
  height: 100%;
  width: 100%;
}

.register-input :deep(.el-input__inner::placeholder) {
  color: #a8abb2;
  opacity: 0.8;
}

.register-button {
  width: 100%;
  height: 50px;
  font-size: 1.1rem;
  margin-top: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.register-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4);
}

.register-button :deep(.el-icon) {
  margin-right: 8px;
}

.register-footer {
  margin-top: 30px;
  text-align: center;
}

.login-link {
  margin-bottom: 20px;
}

.login-link span {
  color: #666;
  margin-right: 5px;
}

.login-link :deep(.el-button) {
  color: #667eea;
  font-weight: 500;
}

.home-link {
  font-size: 0.9rem;
  color: #667eea;
}

.home-link :deep(.el-icon) {
  margin-right: 5px;
}

.register-decoration {
  flex: 1;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  display: none;
  color: white;
  padding: 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.decoration-content {
  max-width: 500px;
  margin: 0 auto;
  text-align: center;
}

.decoration-title h3 {
  font-size: 2rem;
  margin-bottom: 15px;
}

.decoration-title p {
  font-size: 1.2rem;
  opacity: 0.9;
  margin-bottom: 40px;
}

.features {
  display: flex;
  justify-content: space-around;
  flex-wrap: wrap;
  gap: 20px;
  margin-bottom: 40px;
}

.feature-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
  min-width: 100px;
}

.feature-item .el-icon {
  font-size: 2.5rem;
  margin-bottom: 15px;
}

.feature-item span {
  font-size: 1.1rem;
  font-weight: 500;
}

.bubble {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 8s infinite ease-in-out;
}

.bubble:nth-child(1) {
  width: 80px;
  height: 80px;
  top: 10%;
  left: 20%;
  animation-delay: 0s;
}

.bubble:nth-child(2) {
  width: 60px;
  height: 60px;
  top: 20%;
  right: 30%;
  animation-delay: 1s;
}

.bubble:nth-child(3) {
  width: 100px;
  height: 100px;
  bottom: 15%;
  left: 15%;
  animation-delay: 2s;
}

.bubble:nth-child(4) {
  width: 70px;
  height: 70px;
  bottom: 30%;
  right: 25%;
  animation-delay: 3s;
}

.bubble:nth-child(5) {
  width: 50px;
  height: 50px;
  top: 40%;
  left: 40%;
  animation-delay: 4s;
}

.bubble:nth-child(6) {
  width: 90px;
  height: 90px;
  top: 60%;
  right: 20%;
  animation-delay: 5s;
}

.bubble:nth-child(7) {
  width: 60px;
  height: 60px;
  bottom: 20%;
  right: 40%;
  animation-delay: 6s;
}

.bubble:nth-child(8) {
  width: 80px;
  height: 80px;
  top: 30%;
  right: 10%;
  animation-delay: 7s;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0);
  }
  25% {
    transform: translate(20px, 20px);
  }
  50% {
    transform: translate(-20px, 20px);
  }
  75% {
    transform: translate(-20px, -20px);
  }
}

@media (min-width: 992px) {
  .register-decoration {
    display: flex;
  }
}

@media (max-width: 768px) {
  .register-box {
    padding: 40px 30px;
  }
  
  .register-title {
    font-size: 2rem;
  }
  
  .register-container {
    min-height: auto;
  }
  
  .features {
    flex-direction: column;
    align-items: center;
  }
}
</style>