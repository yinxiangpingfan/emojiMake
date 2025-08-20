<template>
  <div class="meme-generator-page">
    <!-- 星空背景装饰 -->
    <div class="star-background">
      <!-- 静态星星 -->
      <div v-for="star in stars" :key="star.id" class="star" :style="{
        left: star.left + 'px',
        top: star.top + 'px',
        width: star.size + 'px',
        height: star.size + 'px',
        opacity: star.opacity
      }"></div>

      <!-- 流星 -->
      <div v-for="meteor in meteors" :key="meteor.id" class="meteor" :style="{
        top: meteor.top + 'px',
        left: meteor.left + 'px',
        width: meteor.length + 'px',
        height: meteor.height + 'px',
        opacity: meteor.opacity,
        animationDelay: meteor.delay + 's',
        animationDuration: meteor.duration + 's'
      }">
        <div class="meteor-effect"></div>
      </div>
    </div>

    <div class="container">
      <h1 class="page-title">AI趣味表情生成</h1>

      <div class="generator-section">
        <div class="generator-card">
          <h2 class="section-title">生成表情包</h2>

          <div class="custom-alert glass-alert">
            <div class="alert-icon">
              <el-icon>
                <Warning />
              </el-icon>
            </div>
            <div class="alert-content">
              <p>目前联网生成是根据角色，通过联网，用文本描述样貌之后进行生成。效果远不及图生表情。建议通过图生视频</p>
            </div>
          </div>

          <el-tabs v-model="activeTab" class="generator-tabs">
            <el-tab-pane label="文生视频" name="textToVideo">
              <TextToVideo 
                :prompt="prompt" 
                :size="size"
                @update:prompt="prompt = $event"
                @update:size="size = $event"
              />
            </el-tab-pane>

            <el-tab-pane label="图生视频" name="imageToVideo">
              <ImageToVideo 
                :prompt="prompt"
                :resolution="resolution"
                :selectedFilePreview="selectedFilePreview"
                @update:prompt="prompt = $event"
                @update:resolution="resolution = $event"
                @update:selectedFilePreview="selectedFilePreview = $event"
                @fileChange="handleFileChange"
                @removeImage="removeImage"
                @selectExample="selectExampleImage"
              />
            </el-tab-pane>

            <el-tab-pane label="角色生成(高级)" name="roleBased">
              <RoleBasedGeneration 
                :roleForm="roleForm"
                @update:roleForm="roleForm = $event"
              />
            </el-tab-pane>
          </el-tabs>

          <el-button type="primary" @click="generateVideoMeme" :loading="generating" class="generate-button"
            size="large" :disabled="!isGenerateEnabled">
            生成表情包
          </el-button>
        </div>

        <VideoPreview 
          :generatedVideo="generatedVideo"
          :jobId="jobId"
          @download="downloadVideo"
        />
      </div>

      <PopularMemesCarousel />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { ElMessage, ElNotification } from 'element-plus'
import axios from 'axios'
import { useUserStore } from '@/stores/user'
import {
  Warning
} from '@element-plus/icons-vue'

// 导入子组件
import TextToVideo from '@/components/meme/TextToVideo.vue'
import ImageToVideo from '@/components/meme/ImageToVideo.vue'
import RoleBasedGeneration from '@/components/meme/RoleBasedGeneration.vue'
import VideoPreview from '@/components/meme/VideoPreview.vue'
import PopularMemesCarousel from '@/components/meme/PopularMemesCarousel.vue'

// 使用user store
const userStore = useUserStore()

// 活动标签页
const activeTab = ref('textToVideo')

// 文生视频相关
const prompt = ref('')
const size = ref('624*624')

// 图生视频相关
const selectedFile = ref<File | null>(null)
const selectedFilePreview = ref<string | null>(null)
const resolution = ref('480P')

// 角色生成相关
const roleForm = ref({
  role: '',
  source: '',
  action: '',
  size: '624*626'
})

// 生成相关
const generating = ref(false)
const jobId = ref('')
const generatedVideo = ref<string | null>(null)
const pollingInterval = ref<number | null>(null)

// 星星数据
const stars = ref<Array<{
  id: number,
  left: number,
  top: number,
  size: number,
  opacity: number
}>>([])

// 流星数据
const meteors = ref<Array<{
  id: number,
  top: number,
  left: number,
  length: number,
  height: number,
  opacity: number,
  delay: number,
  duration: number
}>>([])

let meteorId = 0
let starId = 0

// 计算是否可以生成
const isGenerateEnabled = computed(() => {
  if (activeTab.value === 'textToVideo') {
    return !!prompt.value.trim()
  } else if (activeTab.value === 'imageToVideo') {
    return !!selectedFile.value && !!prompt.value.trim()
  } else if (activeTab.value === 'roleBased') {
    return !!roleForm.value.role && !!roleForm.value.action
  }
  return false
})

// 创建星星
const createStars = () => {
  const starsArray = []
  // 创建100个随机分布的星星
  for (let i = 0; i < 100; i++) {
    starsArray.push({
      id: starId++,
      left: Math.random() * window.innerWidth,
      top: Math.random() * window.innerHeight,
      size: Math.random() * 3 + 1, // 1px 到 4px
      opacity: Math.random() * 0.8 + 0.2 // 0.2 到 1.0
    })
  }
  stars.value = starsArray
}

// 创建流星
const createMeteors = () => {
  const meteorsArray = []
  // 创建10个流星
  for (let i = 0; i < 10; i++) {
    meteorsArray.push({
      id: meteorId++,
      top: Math.random() * window.innerHeight * 0.5, // 只在画面上半部分出现
      left: Math.random() * window.innerWidth,
      length: Math.random() * 100 + 50, // 50px 到 150px
      height: Math.random() * 2 + 1, // 1px 到 3px
      opacity: Math.random() * 0.7 + 0.3, // 0.3 到 1.0
      delay: Math.random() * 10, // 0s 到 10s 的延迟
      duration: Math.random() * 3 + 2 // 2s 到 5s 的持续时间
    })
  }
  meteors.value = meteorsArray
}

// 处理文件上传
const handleFileChange = (file: File) => {
  selectedFile.value = file
}

// 压缩图片函数
const compressImage = (file: File): Promise<Blob> => {
  return new Promise((resolve, reject) => {
    const img = new Image();
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');

    img.onload = () => {
      // 计算压缩比例，确保最长边不超过1080像素
      const maxWidth = 1080;
      const maxHeight = 1080;
      let { width, height } = img;

      if (width > height) {
        if (width > maxWidth) {
          height = (height * maxWidth) / width;
          width = maxWidth;
        }
      } else {
        if (height > maxHeight) {
          width = (width * maxHeight) / height;
          height = maxHeight;
        }
      }

      // 设置canvas尺寸
      canvas.width = width;
      canvas.height = height;

      // 在canvas上绘制图片
      ctx?.drawImage(img, 0, 0, width, height);

      // 转换为Blob
      canvas.toBlob(
        (blob) => {
          if (blob) {
            resolve(blob);
          } else {
            reject(new Error('Canvas toBlob failed'));
          }
        },
        'image/jpeg',
        0.9
      );
    };

    img.onerror = () => {
      reject(new Error('Image loading failed'));
    };

    img.src = URL.createObjectURL(file);
  });
};

// 移除图片
const removeImage = () => {
  selectedFile.value = null
}

// 选择示例图片
const selectExampleImage = (example: { id: number, name: string, url: string }) => {
  console.log('Selected example image:', example.name)
}

// 生成视频表情包
const generateVideoMeme = async () => {
  if (!isGenerateEnabled.value) {
    ElMessage.warning('请填写必要的生成信息')
    return
  }

  generating.value = true
  jobId.value = ''
  generatedVideo.value = null

  try {
    let response: any

    if (activeTab.value === 'textToVideo') {
      // 文生视频
      const formData = new FormData()
      formData.append('type', 'text_to_video')
      formData.append('prompt', prompt.value)
      formData.append('size', size.value)

      response = await axios.post('/api/v1/video/create', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
          'Authorization': `Bearer ${userStore.token}` // 使用store中的token
        }
      })
    } else if (activeTab.value === 'imageToVideo' && selectedFile.value) {
      // 图生视频
      try {
        // 压缩图片
        const compressedImageBlob = await compressImage(selectedFile.value);

        // 转换为base64
        const base64 = await new Promise<string>((resolve, reject) => {
          const reader = new FileReader();
          reader.onload = (e) => {
            resolve(e.target?.result as string);
          };
          reader.onerror = () => {
            reject(new Error('Failed to read compressed image'));
          };
          reader.readAsDataURL(compressedImageBlob);
        });

        const formData = new FormData()
        formData.append('type', 'image_to_video')
        formData.append('prompt', prompt.value)
        formData.append('resolution', resolution.value)
        formData.append('img_base64', base64)

        response = await axios.post('/api/v1/video/create', formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
            'Authorization': `Bearer ${userStore.token}` // 使用store中的token
          }
        })
      } catch (error) {
        console.error('图片压缩失败:', error);
        ElMessage.error('图片处理失败，请重试');
        generating.value = false;
        return;
      }
    } else if (activeTab.value === 'roleBased') {
      // 角色生成
      const formData = new FormData()
      formData.append('role', roleForm.value.role)
      formData.append('source', roleForm.value.source)
      formData.append('action', roleForm.value.action)
      formData.append('size', roleForm.value.size)

      response = await axios.post('/api/v1/video/create_with_prompt', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
          'Authorization': `Bearer ${userStore.token}` // 使用store中的token
        }
      })
    } else {
      ElMessage.warning('请选择生成方式并填写必要信息')
      generating.value = false
      return
    }

    if (response.data.code === 200) {
      jobId.value = response.data.data.job_id
      ElMessage.success('任务创建成功，正在生成中...')

      // 开始轮询任务状态
      startPolling(jobId.value)
    } else {
      ElMessage.error(response.data.message || '任务创建失败')
      generating.value = false
    }
  } catch (error: any) {
    console.error('生成视频表情包失败:', error)
    ElMessage.error(error.response?.data?.message || '生成视频表情包失败，请稍后重试')
    generating.value = false
  }
}

// 开始轮询任务状态
const startPolling = (jobId: string) => {
  // 清除之前的轮询
  if (pollingInterval.value) {
    clearInterval(pollingInterval.value)
  }

  // 每3秒查询一次任务状态
  pollingInterval.value = window.setInterval(async () => {
    try {
      const response = await axios.get(`/api/v1/video/query/${jobId}`, {
        headers: {
          'Authorization': `Bearer ${userStore.token}` // 使用store中的token
        }
      })

      const { status, video_url } = response.data.data

      if (status === 'SUCCEEDED') {
        // 处理视频URL，通过代理访问避免证书问题
        generatedVideo.value = proxyVideoUrl(video_url)
        stopPolling()
        generating.value = false
        ElNotification({
          title: '生成完成',
          message: '视频表情包已生成完成！',
          type: 'success'
        })
      } else if (status === 'FAILED') {
        stopPolling()
        generating.value = false
        ElMessage.error(`生成失败: ${response.data.data.error_message}`)
      } else if (status === 'UNKNOWN') {
        stopPolling()
        generating.value = false
        ElMessage.error('任务不存在')
      }
      // 如果是 PENDING 或 RUNNING 状态，则继续轮询
    } catch (error) {
      console.error('查询任务状态失败:', error)
      stopPolling()
      generating.value = false
      ElMessage.error('查询任务状态失败')
    }
  }, 3000)
}

// 通过代理处理视频URL以避免SSL证书问题
const proxyVideoUrl = (url: string) => {
  if (!url) return url

  // 如果URL包含目标主机地址，则通过代理访问
  if (url.startsWith('https://82.156.59.17:8000/')) {
    // 将原始URL作为查询参数传递给代理端点
    return `/api/proxy?url=${encodeURIComponent(url)}`
  }

  return url
}

// 停止轮询
const stopPolling = () => {
  if (pollingInterval.value) {
    clearInterval(pollingInterval.value)
    pollingInterval.value = null
  }
}

// 下载视频
const downloadVideo = () => {
  if (!generatedVideo.value) return

  // 直接使用当前URL下载（可能是代理URL或原始URL）
  const downloadUrl = generatedVideo.value

  // 创建一个隐藏的iframe来处理下载
  const iframe = document.createElement('iframe')
  iframe.style.display = 'none'
  document.body.appendChild(iframe)

  // 在iframe中创建一个链接并点击它
  if (iframe.contentDocument) {
    const a = iframe.contentDocument.createElement('a')
    a.href = downloadUrl
    a.download = downloadUrl.includes('.gif') ? 'meme.gif' : 'meme.mp4'
    iframe.contentDocument.body.appendChild(a)
    a.click()
  }

  // 清理
  setTimeout(() => {
    document.body.removeChild(iframe)
  }, 100)

  ElMessage.success('下载成功！')
}

// 组件挂载时初始化
onMounted(() => {
  createStars()
  createMeteors()
})

// 组件卸载前清理
onBeforeUnmount(() => {
  stopPolling()
})
</script>

<style scoped>
.meme-generator-page {
  min-height: 100vh;
  padding: 20px;
  background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%);
  position: relative;
  overflow: hidden;
}

/* 星空背景 */
.star-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
}

.star {
  position: absolute;
  background-color: white;
  border-radius: 50%;
  animation: twinkle var(--duration, 5s) infinite ease-in-out;
  animation-delay: var(--delay, 0s);
}

@keyframes twinkle {

  0%,
  100% {
    opacity: var(--min-opacity, 0.2);
  }

  50% {
    opacity: var(--max-opacity, 1);
  }
}

/* 为星星设置不同的动画变量 */
.star:nth-child(2n) {
  --duration: 3s;
  --min-opacity: 0.3;
  --max-opacity: 0.9;
}

.star:nth-child(3n) {
  --duration: 4s;
  --min-opacity: 0.2;
  --max-opacity: 0.8;
}

.star:nth-child(4n) {
  --duration: 6s;
  --min-opacity: 0.4;
  --max-opacity: 1;
}

.meteor {
  position: absolute;
  transform: rotate(30deg);
  overflow: hidden;
}

.meteor-effect {
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.8), transparent);
  animation: meteorMove linear infinite;
}

@keyframes meteorMove {
  0% {
    left: -100%;
  }

  100% {
    left: 100%;
  }
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  position: relative;
  z-index: 1;
}

.page-title {
  text-align: center;
  font-size: 2.5rem;
  color: #e0e0ff;
  margin-bottom: 30px;
  font-weight: 700;
  text-shadow: 0 0 10px rgba(224, 224, 255, 0.3);
}

.generator-section {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 30px;
  margin-bottom: 50px;
}

.generator-card {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 15px;
  padding: 30px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.section-title {
  font-size: 1.5rem;
  color: #e0e0ff;
  margin-bottom: 20px;
  font-weight: 600;
  text-shadow: 0 0 5px rgba(224, 224, 255, 0.2);
}

.generator-tabs {
  margin-bottom: 30px;
}

.generator-tabs :deep(.el-tabs__item) {
  font-size: 1.1rem;
  color: rgba(224, 224, 255, 0.8);
}

.generator-tabs :deep(.el-tabs__item.is-active) {
  color: #ffffff;
  font-weight: 600;
}

.generator-tabs :deep(.el-tabs__active-bar) {
  background-color: #667eea;
}

.generator-tabs :deep(.el-tabs__nav-wrap::after) {
  background-color: rgba(255, 255, 255, 0.1);
}

.generate-button {
  width: 100%;
  font-size: 1.1rem;
  padding: 15px;
  border-radius: 25px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  transition: all 0.3s ease;
  margin-top: 10px;
  color: white;
}

.generate-button:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.6);
}

.generate-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.glass-alert {
  display: flex;
  align-items: center;
  padding: 20px;
  border-radius: 15px;
  background: rgba(40, 40, 80, 0.4);
  border: 1px solid rgba(102, 126, 234, 0.3);
  backdrop-filter: blur(20px);
  margin-bottom: 25px;
  color: #e0e0ff;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.2);
  position: relative;
  overflow: hidden;
}

.glass-alert::before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(102, 126, 234, 0.5), transparent);
}

.glass-alert .alert-icon {
  font-size: 24px;
  margin-right: 15px;
  color: #ffd54f;
  flex-shrink: 0;
}

.glass-alert .alert-content p {
  margin: 0;
  font-size: 15px;
  line-height: 1.6;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}

@media (max-width: 992px) {
  .generator-section {
    grid-template-columns: 1fr;
  }

  .page-title {
    font-size: 2rem;
  }

  .generator-card,
  .preview-section,
  .popular-section {
    padding: 20px;
  }
}

@media (max-width: 768px) {
  .meme-generator-page {
    padding: 15px;
  }

  .generator-card,
  .preview-section,
  .popular-section {
    padding: 15px;
  }

  .page-title {
    font-size: 1.8rem;
  }
}
</style>