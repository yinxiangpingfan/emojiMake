<template>
  <div class="image-to-video-section">
    <div class="example-images-section">
      <h3 class="example-title">示例图片（点击选择）</h3>
      <div class="example-images-container">
        <div 
          v-for="(example, index) in exampleImages" 
          :key="index"
          class="example-image-item"
          @click="selectExampleImage(example)"
        >
          <img 
            :src="example.url" 
            :alt="example.name"
            class="example-image"
          />
        </div>
      </div>
    </div>

    <div class="upload-area" v-if="!selectedFilePreview">
      <el-upload 
        ref="uploadRef"
        class="upload-demo glass-upload" 
        drag 
        :auto-upload="false" 
        :on-change="handleFileChange"
        :show-file-list="false" 
        accept="image/*">
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          拖拽图片到此处或 <em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            支持 JPG/PNG 格式，文件大小不超过 5MB
          </div>
        </template>
      </el-upload>
    </div>

    <div class="image-preview-area glass-preview" v-else>
      <img :src="selectedFilePreview" alt="Preview" class="preview-image" @click="triggerUpload" />
      <p class="preview-tip">点击图片重新上传</p>
      <el-button type="primary" @click="removeImage" size="small" class="remove-button">
        删除图片
      </el-button>
    </div>

    <el-input v-model="localPrompt" type="textarea" placeholder="请输入视频描述，例如：让图片中的人物开始跳舞" :rows="3"
      class="prompt-input glass-input styled-prompt" />

    <div class="advanced-settings" style="margin-top: 20px;">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-select v-model="localResolution" placeholder="选择分辨率" class="glass-select">
            <el-option label="480P (推荐)" value="480P" />
            <el-option label="720P" value="720P" />
          </el-select>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'

// 定义组件属性
const props = defineProps<{
  prompt: string,
  resolution: string,
  selectedFilePreview: string | null
}>()

// 定义事件
const emit = defineEmits<{
  (e: 'update:prompt', value: string): void,
  (e: 'update:resolution', value: string): void,
  (e: 'update:selectedFilePreview', value: string | null): void,
  (e: 'fileChange', file: File): void,
  (e: 'removeImage'): void,
  (e: 'selectExample', example: { id: number, name: string, url: string }): void
}>()

// 本地响应式数据
const localPrompt = ref(props.prompt)
const localResolution = ref(props.resolution)

// 监听属性变化
watch(() => props.prompt, (newVal) => {
  localPrompt.value = newVal
})

watch(() => props.resolution, (newVal) => {
  localResolution.value = newVal
})

// 监听本地数据变化并发出事件
watch(localPrompt, (newVal) => {
  emit('update:prompt', newVal)
})

watch(localResolution, (newVal) => {
  emit('update:resolution', newVal)
})

// 示例图片数据
const exampleImages = ref([
  { id: 1, name: '示例图片1', url: new URL('../../assets/images/ps1.jpg', import.meta.url).href },
  { id: 2, name: '示例图片2', url: new URL('../../assets/images/ps2.jpg', import.meta.url).href },
  { id: 3, name: '示例图片3', url: new URL('../../assets/images/ps3.jpg', import.meta.url).href },
  { id: 4, name: '示例图片4', url: new URL('../../assets/images/ps4.jpg', import.meta.url).href }
])

const uploadRef = ref()

// 处理文件上传
const handleFileChange = (file: any) => {
  emit('update:selectedFilePreview', URL.createObjectURL(file.raw))
  emit('fileChange', file.raw)
  
  // 确保DOM更新
  nextTick(() => {
    // 什么都不需要做，只是确保DOM更新完成
  })
}

// 触发上传
const triggerUpload = () => {
  // 直接触发上传组件的点击事件
  if (uploadRef.value) {
    uploadRef.value.$el.querySelector('input[type="file"]')?.click()
  }
}

// 移除图片
const removeImage = () => {
  emit('removeImage')
  emit('update:selectedFilePreview', null)
  // 清空上传组件的文件列表
  uploadRef.value?.clearFiles()
}

// 选择示例图片
const selectExampleImage = async (example: { id: number, name: string, url: string }) => {
  try {
    // 获取图片数据
    const response = await fetch(example.url);
    const blob = await response.blob();
    
    // 创建File对象
    const file = new File([blob], `example${example.id}.jpg`, { type: blob.type });
    
    // 创建文件对象以匹配Element Plus Upload的格式
    const uploadFile = {
      name: `example${example.id}.jpg`,
      raw: file,
      uid: Date.now()
    };
    
    // 处理文件
    handleFileChange(uploadFile);
    
    emit('selectExample', example)
    ElMessage.success(`已选择示例图片: ${example.name}`);
  } catch (error) {
    console.error('选择示例图片失败:', error);
    ElMessage.error('选择示例图片失败，请重试');
  }
}
</script>

<style scoped>
.glass-upload :deep(.el-upload) {
  display: block;
}

.glass-upload :deep(.el-upload-dragger) {
  background: rgba(255, 255, 255, 0.1);
  border: 1px dashed rgba(255, 255, 255, 0.3);
  border-radius: 15px;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
  width: 100%;
  padding: 30px 0;
}

.glass-upload :deep(.el-upload-dragger:hover) {
  border-color: #667eea;
  background: rgba(255, 255, 255, 0.15);
  box-shadow: 0 0 20px rgba(102, 126, 234, 0.3);
}

.glass-upload :deep(.el-upload-dragger .el-icon--upload) {
  font-size: 60px;
  color: rgba(224, 224, 255, 0.6);
  margin-bottom: 15px;
}

.glass-upload :deep(.el-upload__text) {
  color: rgba(224, 224, 255, 0.8);
  font-size: 16px;
  margin-bottom: 10px;
}

.glass-upload :deep(.el-upload__text em) {
  color: #667eea;
  font-style: normal;
  font-weight: 500;
}

.glass-upload :deep(.el-upload__tip) {
  color: rgba(224, 224, 255, 0.6);
  font-size: 14px;
}

.glass-preview {
  text-align: center;
  padding: 25px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 15px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
  margin-bottom: 20px;
}

.glass-preview:hover {
  background: rgba(255, 255, 255, 0.15);
  box-shadow: 0 0 20px rgba(102, 126, 234, 0.3);
}

.glass-select :deep(.el-select .el-input__wrapper) {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.glass-select :deep(.el-select .el-input__wrapper:hover) {
  box-shadow: 0 0 15px rgba(102, 126, 234, 0.4);
}

.glass-select :deep(.el-select .el-input.is-focus .el-input__wrapper) {
  border-color: #667eea;
  box-shadow: 0 0 15px rgba(102, 126, 234, 0.4);
}

.glass-select :deep(.el-select .el-input__inner) {
  color: #fff;
  background: transparent;
}

.glass-select :deep(.el-select .el-input__inner::placeholder) {
  color: rgba(224, 224, 255, 0.6);
}

.glass-select :deep(.el-select-dropdown) {
  background: rgba(30, 30, 60, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  backdrop-filter: blur(20px);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
}

.glass-select :deep(.el-select-dropdown__item) {
  color: rgba(224, 224, 255, 0.8);
}

.glass-select :deep(.el-select-dropdown__item:hover) {
  background: rgba(102, 126, 234, 0.2);
}

.glass-select :deep(.el-select-dropdown__item.selected) {
  color: #667eea;
  font-weight: 500;
}

.example-images-section {
  margin-bottom: 20px;
}

.example-title {
  font-size: 1.1rem;
  color: #e0e0ff;
  margin-bottom: 15px;
  font-weight: 500;
  text-shadow: 0 0 5px rgba(224, 224, 255, 0.2);
}

.example-images-container {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.example-image-item {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.1);
}

.example-image-item:hover {
  transform: translateY(-5px);
  border-color: #667eea;
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}

.example-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.upload-demo {
  width: 100%;
}

.image-preview-area {
  text-align: center;
  padding: 20px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  border: 1px dashed rgba(255, 255, 255, 0.3);
  cursor: pointer;
  transition: all 0.3s ease;
}

.image-preview-area:hover {
  background: rgba(255, 255, 255, 0.2);
}

.preview-image {
  max-width: 100%;
  max-height: 300px;
  border-radius: 8px;
  margin-bottom: 15px;
}

.preview-tip {
  margin: 10px 0;
  color: #ccc;
  font-size: 14px;
}

.remove-button {
  margin-top: 10px;
}

.glass-input :deep(.el-textarea__inner),
.glass-input :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  color: #fff;
  font-size: 1rem;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.glass-input :deep(.el-textarea__inner) {
  padding: 15px;
}

.glass-input :deep(.el-input__wrapper) {
  padding: 10px 15px;
}

.glass-input :deep(.el-textarea__inner):focus,
.glass-input :deep(.el-input.is-focus .el-input__wrapper) {
  border-color: #667eea;
  box-shadow: 0 0 15px rgba(102, 126, 234, 0.4);
  outline: none;
}

.glass-input :deep(.el-textarea__inner)::placeholder,
.glass-input :deep(.el-input__inner)::placeholder {
  color: rgba(224, 224, 255, 0.6);
}

.glass-input :deep(.el-input__inner) {
  background: transparent;
  color: #fff;
}

.styled-prompt :deep(.el-textarea__inner) {
  background: rgba(255, 255, 255, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  color: #fff;
  font-size: 1rem;
  padding: 15px;
  transition: all 0.3s ease;
  backdrop-filter: blur(5px);
}

.styled-prompt :deep(.el-textarea__inner):focus {
  border-color: #667eea;
  box-shadow: 0 0 15px rgba(102, 126, 234, 0.4);
  outline: none;
}

.styled-prompt :deep(.el-textarea__inner)::placeholder {
  color: rgba(224, 224, 255, 0.6);
}
</style>