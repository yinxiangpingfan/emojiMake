<template>
  <div class="home-container">
    <div class="hero-section">
      <div class="content">
        <h1 class="title">欢迎使用趣味表情生成</h1>
        <p class="subtitle">创建属于你的专属表情包</p>
        <div class="button-container">
          <el-button type="primary" size="large" class="generate-button" @click="handleGenerateMeme">
            立即制作表情包
          </el-button>
        </div>
      </div>
      <div class="meme-examples">
        <img v-for="meme in memeExamples" :key="meme.id" :src="meme.url" :alt="meme.name" class="meme-example">
      </div>
    </div>

    <div class="intro-section">
      <div class="section-container">
        <h2 class="section-title">为什么选择我们的趣味表情生成？</h2>
        <div class="features-grid">
          <div class="feature-card">
            <el-icon class="feature-icon">
              <component :is="Star" />
            </el-icon>
            <h3>简单易用</h3>
            <p>无需设计经验，三步即可制作专属表情包</p>
          </div>
          <div class="feature-card">
            <el-icon class="feature-icon">
              <component :is="Picture" />
            </el-icon>
            <h3>多种模板</h3>
            <p>海量热门模板供你选择，紧跟潮流</p>
          </div>
          <div class="feature-card">
            <el-icon class="feature-icon">
              <component :is="Download" />
            </el-icon>
            <h3>一键下载</h3>
            <p>高清表情包一键生成，直接保存到设备</p>
          </div>
        </div>
      </div>
    </div>

    <div class="usage-section">
      <div class="section-container">
        <h2 class="section-title">应用场景</h2>
        <div class="usage-grid">
          <div class="usage-card">
            <h3>社交聊天</h3>
            <p>在微信、QQ等社交平台中用个性化表情包表达情感</p>
          </div>
          <div class="usage-card">
            <h3>工作沟通</h3>
            <p>在团队协作中用轻松幽默的方式传达信息</p>
          </div>
          <div class="usage-card">
            <h3>内容创作</h3>
            <p>为博客、社交媒体帖子添加趣味性配图</p>
          </div>
          <div class="usage-card">
            <h3>个人娱乐</h3>
            <p>制作专属表情包，与朋友分享快乐时光</p>
          </div>
        </div>
      </div>
    </div>

    <div class="faq-section">
      <div class="section-container">
        <h2 class="section-title">常见问题</h2>
        <div class="faq-list">
          <el-collapse v-model="activeNames">
            <el-collapse-item title="如何制作表情包？" name="1">
              <p>点击首页的"立即制作表情包"按钮，选择上传图片、网络图片或模板，添加文字即可生成。</p>
            </el-collapse-item>
            <el-collapse-item title="支持哪些图片格式？" name="2">
              <p>目前支持 JPG、PNG 格式的图片，文件大小不超过 5MB。</p>
            </el-collapse-item>
            <el-collapse-item title="生成的表情包有水印吗？" name="3">
              <p>我们完全免费且无水印，您可以自由使用生成的表情包。</p>
            </el-collapse-item>
            <el-collapse-item title="表情包可以商用吗？" name="4">
              <p>使用您自己拥有的图片制作的表情包可以用于商业用途，但使用模板制作的请遵守相关版权规定。</p>
            </el-collapse-item>
          </el-collapse>
        </div>
      </div>
    </div>

    <div class="cta-section">
      <div class="section-container">
        <h2 class="section-title">开始制作你的专属表情包</h2>
        <p class="cta-subtitle">数百万用户的选择，立即体验</p>
        <el-button type="primary" size="large" class="cta-button" @click="handleGenerateMeme">
          免费制作表情包
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { Star, Picture, Download } from '@element-plus/icons-vue'

// 模拟一些表情包示例
const memeExamples = ref([
  { id: 1, name: 'Example 1', url: new URL('../assets/images/example1.gif', import.meta.url).href },
  { id: 2, name: 'Example 2', url: new URL('../assets/images/example2.gif', import.meta.url).href },
  { id: 3, name: 'Example 3', url: new URL('../assets/images/example3.gif', import.meta.url).href },
  { id: 4, name: 'Example 4', url: new URL('../assets/images/example4.gif', import.meta.url).href }
])

const activeNames = ref(['1'])

const router = useRouter()
const userStore = useUserStore()

const handleGenerateMeme = () => {
  if (userStore.isAuthenticated) {
    // 用户已登录，跳转到生成表情包页面
    router.push('/generate')
  } else {
    // 用户未登录，跳转到登录页面
    router.push('/login')
  }
}

</script>

<style scoped>
.home-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  scroll-behavior: smooth;
  width: 100%;
  position: relative;
  overflow-x: hidden;
}

.hero-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 80vh;
  padding: 20px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.8) 0%, rgba(118, 75, 162, 0.8) 100%);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  width: 100%;
  color: white;
  position: relative;
}

.content {
  text-align: center;
  margin-bottom: 40px;
  max-width: 800px;
}

.title {
  font-size: 3rem;
  font-weight: 700;
  margin-bottom: 1rem;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
}

.subtitle {
  font-size: 1.5rem;
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 2rem;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.3);
}

.button-container {
  margin-top: 2rem;
}

.generate-button {
  font-size: 1.5rem;
  padding: 20px 40px;
  border-radius: 50px;
  background: linear-gradient(45deg, #ff6b6b, #ffa502);
  border: none;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2);
  transition: all 0.3s ease;
}

.generate-button:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.3);
}

.generate-button:active {
  transform: translateY(1px);
}

.meme-examples {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 20px;
  max-width: 1000px;
}

.meme-example {
  width: 150px;
  height: 150px;
  object-fit: cover;
  border-radius: 10px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.2);
  transition: transform 0.3s ease;
}

.meme-example:hover {
  transform: scale(1.05);
}

.section-container {
  max-width: 1200px;
  padding: 60px 20px;
  width: 100%;
  margin: 0 auto;
}

.section-title {
  text-align: center;
  font-size: 2.5rem;
  margin-bottom: 50px;
  color: #333;
}

.intro-section {
  background-color: rgba(248, 249, 250, 0.8);
  width: 100%;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 30px;
}

.feature-card {
  background: white;
  padding: 30px;
  border-radius: 10px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
  text-align: center;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
}

.feature-icon {
  font-size: 3rem;
  margin-bottom: 20px;
  color: #667eea;
}

.feature-card h3 {
  font-size: 1.5rem;
  margin-bottom: 15px;
  color: #333;
}

.feature-card p {
  color: #666;
  font-size: 1rem;
  line-height: 1.6;
}

.usage-section {
  width: 100%;
  background-color: rgba(248, 249, 250, 0.8);
}

.usage-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 30px;
}

.usage-card {
  background: #f8f9fa;
  padding: 30px;
  border-radius: 10px;
  text-align: center;
  transition: background-color 0.3s ease;
}

.usage-card:hover {
  background: #e9ecef;
}

.usage-card h3 {
  font-size: 1.3rem;
  margin-bottom: 15px;
  color: #333;
}

.usage-card p {
  color: #666;
  line-height: 1.6;
}

.faq-section {
  background-color: rgba(248, 249, 250, 0.8);
  width: 100%;
}

.faq-list {
  max-width: 800px;
  margin: 0 auto;
}

.faq-list :deep(.el-collapse-item__header) {
  font-size: 1.1rem;
  padding: 20px 30px;
}

.faq-list :deep(.el-collapse-item__content) {
  padding: 20px 30px;
  font-size: 1rem;
  color: #666;
  line-height: 1.6;
  background: white;
}

.cta-section {
  width: 100%;
  text-align: center;
  padding: 80px 20px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.8) 0%, rgba(118, 75, 162, 0.8) 100%);
  color: white;
}

.cta-section .section-title {
  color: white;
  margin-bottom: 20px;
}

.cta-subtitle {
  font-size: 1.2rem;
  margin-bottom: 30px;
  color: rgba(255, 255, 255, 0.9);
}

.cta-button {
  font-size: 1.2rem;
  padding: 20px 40px;
  border-radius: 50px;
  background: white;
  color: #667eea;
  border: none;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2);
  transition: all 0.3s ease;
}

.cta-button:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.3);
  background: #f8f9fa;
}

@media (max-width: 768px) {
  .title {
    font-size: 2rem;
  }

  .subtitle {
    font-size: 1.2rem;
  }

  .section-title {
    font-size: 2rem;
  }

  .meme-examples {
    gap: 10px;
  }

  .meme-example {
    width: 100px;
    height: 100px;
  }

  .generate-button {
    font-size: 1.2rem;
    padding: 15px 30px;
  }

  .section-container {
    padding: 40px 15px;
  }

  .features-grid,
  .usage-grid {
    grid-template-columns: 1fr;
  }
}
</style>