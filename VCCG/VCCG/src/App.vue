<template>
  <div id="app">
    <!-- 泡泡背景容器 -->
    <div class="bubbles-container">
      <div v-for="bubble in bubbles" :key="bubble.id" class="bubble" :style="{
        left: bubble.left + 'px',
        width: bubble.size + 'px',
        height: bubble.size + 'px',
        backgroundColor: bubble.color,
        bottom: bubble.bottom + 'px',
        opacity: bubble.opacity
      }"></div>
    </div>

    <header>
      <nav class="navbar">
        <div class="nav-brand">
          <h2 class="logo-text">趣味表情生成</h2>
        </div>
        <div class="nav-items">
          <el-button type="primary" link @click="goToHome" class="nav-link">
            首页
          </el-button>
          <el-button type="primary" link @click="goToGenerator" class="nav-link">
            生成表情包
          </el-button>

          <!-- 用户头像/登录按钮 -->
          <div v-if="!userStore.isLoggedIn" class="user-actions">
            <el-button type="primary" link @click="goToLogin" class="nav-link">
              登录
            </el-button>
          </div>

          <div v-else class="user-avatar-dropdown">
            <el-dropdown @command="handleUserCommand">
              <div class="user-avatar">
                <span class="avatar-text">{{ getAvatarText() }}</span>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="changePassword">
                    修改密码
                  </el-dropdown-item>
                  <el-dropdown-item command="logout">
                    退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </nav>
    </header>
    <main>
      <router-view />
    </main>

    <!-- 修改密码对话框 -->
    <el-dialog v-model="changePasswordDialogVisible" title="修改密码" width="450px" center class="change-password-dialog">
      <div class="dialog-description">
        请输入您的新密码，密码长度至少8位
      </div>
      <el-form :model="passwordForm" ref="passwordFormRef" :rules="passwordRules" label-position="top"
        class="password-form">
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="passwordForm.newPassword" type="password" show-password placeholder="请输入至少8位新密码"
            class="password-input" size="large" />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" type="password" show-password placeholder="请再次输入新密码"
            class="password-input" size="large" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="changePasswordDialogVisible = false" size="large" class="cancel-button">
            取消
          </el-button>
          <el-button type="primary" @click="submitPasswordChange" size="large" class="confirm-button">
            确定修改
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import axios from 'axios'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// 修改密码表单相关
const changePasswordDialogVisible = ref(false)
const passwordForm = ref({
  newPassword: '',
  confirmPassword: ''
})

const passwordFormRef = ref<FormInstance>()

// 修改密码表单验证规则
const passwordRules = ref<FormRules>({
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 8, message: '密码长度至少8位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: any) => {
        if (value !== passwordForm.value.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
})

// 泡泡数据
const bubbles = ref<Array<{
  id: number,
  left: number,
  size: number,
  initialSize: number,
  color: string,
  bottom: number,
  opacity: number
}>>([])

let bubbleId = 0
let animationFrameId: number | null = null

// 生成随机颜色
const getRandomColor = () => {
  const colors = [
    'rgba(255, 107, 107, 0.7)',   // 红色
    'rgba(107, 185, 255, 0.7)',   // 蓝色
    'rgba(107, 255, 186, 0.7)',   // 青色
    'rgba(255, 217, 107, 0.7)',   // 黄色
    'rgba(204, 107, 255, 0.7)',   // 紫色
    'rgba(255, 159, 107, 0.7)'    // 橙色
  ]
  return colors[Math.floor(Math.random() * colors.length)]
}

// 创建新泡泡
const createBubble = () => {
  const initialSize = Math.random() * 20 + 10; // 初始大小 10px 到 30px
  const newBubble = {
    id: bubbleId++,
    left: Math.random() * window.innerWidth,
    size: initialSize,
    initialSize: initialSize,
    color: getRandomColor(),
    bottom: -50, // 从屏幕底部外开始
    opacity: Math.random() * 0.5 + 0.3 // 0.3 到 0.8 的透明度
  }

  bubbles.value.push(newBubble)
}

// 更新泡泡位置和大小
const updateBubbles = () => {
  bubbles.value = bubbles.value.filter(bubble => {
    // 泡泡向上移动
    bubble.bottom += 1.5

    // 泡泡大小随位置增加而变大（越靠近顶部越大）
    // 根据泡泡的位置计算大小，顶部的泡泡比初始大小大3倍
    const progress = Math.min(bubble.bottom / window.innerHeight, 1);
    bubble.size = bubble.initialSize * (1 + progress * 2);

    // 如果泡泡移出屏幕顶部，则移除
    return bubble.bottom < window.innerHeight + 100
  })

  // 随机生成新泡泡
  if (Math.random() > 0.7) {
    createBubble()
  }

  animationFrameId = requestAnimationFrame(animateBubbles)
}

// 动画循环
const animateBubbles = () => {
  updateBubbles()
}

// 开始动画
const startAnimation = () => {
  if (animationFrameId === null) {
    animationFrameId = requestAnimationFrame(animateBubbles)
  }
}

const goToHome = () => {
  router.push('/')
}

const goToGenerator = () => {
  router.push('/generate')
}

const goToLogin = () => {
  router.push('/login')
}

// 获取头像文本（手机号缩写）
const getAvatarText = () => {
  if (!userStore.userInfo.phone) return 'U'
  // 取手机号的后四位
  return userStore.userInfo.phone.slice(-4)
}

// 处理用户下拉菜单命令
const handleUserCommand = (command: string) => {
  switch (command) {
    case 'changePassword':
      openChangePasswordDialog()
      break
    case 'logout':
      logout()
      break
  }
}

// 打开修改密码对话框
const openChangePasswordDialog = () => {
  // 重置表单
  passwordForm.value = {
    newPassword: '',
    confirmPassword: ''
  }
  changePasswordDialogVisible.value = true
}

// 提交密码修改
const submitPasswordChange = async () => {
  if (!passwordFormRef.value) return

  try {
    await passwordFormRef.value.validate()

    // 调用修改密码接口
    const formData = new FormData()
    formData.append('newPassword', passwordForm.value.newPassword)

    const response = await axios.post('/api/v1/users/change-password', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Authorization': `Bearer ${userStore.token}`
      }
    })

    if (response.data.code === 0) {
      ElMessage.success('密码修改成功')
      changePasswordDialogVisible.value = false
    } else {
      ElMessage.error(response.data.message || '密码修改失败')
    }
  } catch (error: any) {
    console.error('修改密码失败:', error)
    ElMessage.error(error.response?.data?.message || '密码修改失败')
  }
}

// 退出登录
const logout = () => {
  ElMessageBox.confirm('确定要退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    userStore.logout()
    ElMessage.success('已退出登录')
    // 如果当前在需要登录的页面，跳转到首页
    if (route.path === '/generate') {
      router.push('/')
    }
  }).catch(() => {
    // 取消退出
  })
}

// 监听路由变化，根据当前页面调整导航栏样式
watch(() => route.path, (newPath) => {
  const navbar = document.querySelector('.navbar')
  if (navbar) {
    if (newPath === '/generate') {
      navbar.classList.add('generator-page')
    } else {
      navbar.classList.remove('generator-page')
    }
  }
})

// 组件挂载时启动动画
onMounted(() => {
  startAnimation()

  // 检查登录状态
  userStore.checkLoginStatus()

  // 初始创建一些泡泡
  for (let i = 0; i < 10; i++) {
    setTimeout(() => {
      createBubble()
    }, i * 300)
  }

  // 监听窗口大小变化，更新窗口高度
  window.addEventListener('resize', () => {
    // 窗口大小变化时不需要特殊处理，泡泡会自动适应
  });

  // 初始化导航栏样式
  if (route.path === '/generate') {
    const navbar = document.querySelector('.navbar')
    if (navbar) {
      navbar.classList.add('generator-page')
    }
  }
})

// 组件卸载前清理动画
onBeforeUnmount(() => {
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
  }
})
</script>

<style scoped>
#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  width: 100%;
  position: relative;
  overflow-x: hidden;
}

.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.1);
  width: 100%;
  position: relative;
  z-index: 10;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.navbar.generator-page {
  background: rgba(48, 43, 99, 0.8);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.nav-brand {
  display: flex;
  align-items: center;
}

.logo-text {
  color: #667eea;
  margin: 0;
  font-size: 1.8rem;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 1px;
}

.navbar.generator-page .logo-text {
  background: linear-gradient(135deg, #ffffff 0%, #e0e0ff 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.nav-items {
  display: flex;
  gap: 1.5rem;
  align-items: center;
}

.nav-link {
  font-size: 1.1rem;
  font-weight: 500;
  color: #666;
  position: relative;
  transition: all 0.3s ease;
  padding: 0.5rem 1rem;
  border-radius: 20px;
}

.navbar.generator-page .nav-link {
  color: rgba(255, 255, 255, 0.8);
}

.nav-link:hover {
  color: #667eea;
  background: rgba(102, 126, 234, 0.1);
  transform: translateY(-2px);
}

.navbar.generator-page .nav-link:hover {
  color: white;
  background: rgba(102, 126, 234, 0.3);
}

.nav-link::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 2px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  transition: width 0.3s ease;
}

.navbar.generator-page .nav-link::after {
  background: linear-gradient(135deg, #ffffff 0%, #e0e0ff 100%);
}

.nav-link:hover::after {
  width: 70%;
}

/* 用户头像样式 */
.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: white;
  font-weight: bold;
  font-size: 14px;
  transition: all 0.3s ease;
}

.user-avatar:hover {
  transform: scale(1.05);
  box-shadow: 0 0 10px rgba(102, 126, 234, 0.5);
}

/* 用户下拉菜单样式 */
.user-avatar-dropdown {
  position: relative;
}

.footer {
  margin-top: auto;
  text-align: center;
  padding: 1.5rem;
  background: rgba(245, 245, 245, 0.8);
  color: #666;
  width: 100%;
  position: relative;
  z-index: 10;
  backdrop-filter: blur(10px);
  border-top: 1px solid rgba(0, 0, 0, 0.05);
}

/* 泡泡背景样式 */
.bubbles-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: -1;
}

.bubble {
  position: absolute;
  border-radius: 50%;
  bottom: -100px;
  animation-timing-function: linear;
  animation-iteration-count: infinite;
  transform: translateX(-50%);
}

/* 修改密码对话框样式 */
.change-password-dialog {
  border-radius: 12px;
  overflow: hidden;
}

.change-password-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
  text-align: center;
}

.change-password-dialog :deep(.el-dialog__title) {
  color: white;
  font-size: 1.5rem;
  font-weight: 600;
}

.change-password-dialog :deep(.el-dialog__headerbtn) {
  top: 20px;
}

.change-password-dialog :deep(.el-dialog__headerbtn .el-dialog__close) {
  color: white;
  font-size: 24px;
}

.change-password-dialog :deep(.el-dialog__body) {
  padding: 30px;
}

.dialog-description {
  color: #666;
  font-size: 14px;
  margin-bottom: 25px;
  text-align: center;
}

.password-form :deep(.el-form-item__label) {
  font-weight: 500;
  color: #333;
  font-size: 15px;
}

.password-input :deep(.el-input__wrapper) {
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(102, 126, 234, 0.2);
  transition: all 0.3s ease;
}

.password-input :deep(.el-input__wrapper:hover) {
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.2);
  border-color: rgba(102, 126, 234, 0.5);
}

.password-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
  border-color: #667eea;
}

.dialog-footer {
  display: flex;
  justify-content: center;
  gap: 20px;
  padding: 10px 0;
}

.cancel-button {
  border-radius: 25px;
  padding: 12px 30px;
  font-size: 16px;
  border: 1px solid #ddd;
  color: #666;
  transition: all 0.3s ease;
}

.cancel-button:hover {
  background: #f5f5f5;
  border-color: #ccc;
  transform: translateY(-2px);
}

.confirm-button {
  border-radius: 25px;
  padding: 12px 30px;
  font-size: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  transition: all 0.3s ease;
}

.confirm-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}
</style>
