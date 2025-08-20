<template>
  <div class="text-to-video-section">
    <el-input v-model="localPrompt" type="textarea" placeholder="请输入视频描述，例如：一只可爱的猫咪在跳舞" :rows="4"
      class="prompt-input glass-input" />
    <div class="advanced-settings">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-select v-model="localSize" placeholder="选择分辨率" class="glass-select">
            <el-option label="624*624 (推荐)" value="624*624" />
            <el-option label="832*480" value="832*480" />
            <el-option label="480*832" value="480*832" />
          </el-select>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

// 定义组件属性
const props = defineProps<{
  prompt: string,
  size: string
}>()

// 定义事件
const emit = defineEmits<{
  (e: 'update:prompt', value: string): void,
  (e: 'update:size', value: string): void
}>()

// 本地响应式数据
const localPrompt = ref(props.prompt)
const localSize = ref(props.size)

// 监听属性变化
watch(() => props.prompt, (newVal) => {
  localPrompt.value = newVal
})

watch(() => props.size, (newVal) => {
  localSize.value = newVal
})

// 监听本地数据变化并发出事件
watch(localPrompt, (newVal) => {
  emit('update:prompt', newVal)
})

watch(localSize, (newVal) => {
  emit('update:size', newVal)
})
</script>

<style scoped>
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

.prompt-input {
  margin-bottom: 20px;
}

.advanced-settings {
  margin: 20px 0;
}
</style>