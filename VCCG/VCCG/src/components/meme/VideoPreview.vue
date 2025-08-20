<template>
  <div class="preview-section">
    <h2 class="section-title">预览</h2>
    <div class="preview-container">
      <div v-if="generatedVideo" class="video-preview">
        <img v-if="generatedVideo.endsWith('.gif')" :src="generatedVideo" alt="生成的视频表情包" class="preview-image" />
        <video v-else :src="generatedVideo" controls class="preview-video" loop muted />
        <el-button type="success" @click="downloadVideo" class="download-button">
          <el-icon>
            <download />
          </el-icon>
          下载表情包
        </el-button>
      </div>
      <div v-else-if="jobId && !generatedVideo" class="generating-status">
        <div class="status-info">
          <el-icon class="is-loading" color="#667eea">
            <loading />
          </el-icon>
          <p>表情包生成中，请稍候...</p>
        </div>
      </div>
      <div v-else class="no-preview">
        <el-icon>
          <picture />
        </el-icon>
        <p>生成的表情包将在这里显示</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Download, Picture, Loading } from '@element-plus/icons-vue'

// 定义组件属性
const props = defineProps<{
  generatedVideo: string | null,
  jobId: string
}>()

// 定义事件
const emit = defineEmits<{
  (e: 'download'): void
}>()

// 下载视频
const downloadVideo = () => {
  emit('download')
}
</script>

<style scoped>
.preview-section {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 15px;
  padding: 30px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  flex-direction: column;
}

.section-title {
  font-size: 1.5rem;
  color: #e0e0ff;
  margin-bottom: 20px;
  font-weight: 600;
  text-shadow: 0 0 5px rgba(224, 224, 255, 0.2);
}

.preview-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  flex: 1;
}

.video-preview {
  width: 100%;
  text-align: center;
}

.preview-image,
.preview-video {
  max-width: 100%;
  max-height: 300px;
  border-radius: 10px;
  margin-bottom: 20px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.3);
}

.download-button {
  font-size: 1.1rem;
  padding: 15px 30px;
  border-radius: 25px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  transition: all 0.3s ease;
  color: white;
}

.download-button:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.6);
}

.generating-status {
  text-align: center;
  color: rgba(224, 224, 255, 0.7);
  padding: 50px 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.status-info .el-icon {
  font-size: 3rem;
  margin-bottom: 20px;
  color: #667eea;
}

.status-info p {
  margin: 10px 0;
  font-size: 1.1rem;
  color: #e0e0ff;
}

.job-id {
  font-size: 0.9rem;
  color: rgba(224, 224, 255, 0.7);
  font-family: monospace;
}

.no-preview {
  text-align: center;
  color: rgba(224, 224, 255, 0.7);
  padding: 50px 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.no-preview .el-icon {
  font-size: 3rem;
  margin-bottom: 20px;
  color: rgba(224, 224, 255, 0.4);
}
</style>