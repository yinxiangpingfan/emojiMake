<template>
  <div class="role-based-section">
    <el-form :model="localRoleForm" label-position="top">
      <el-form-item label="角色名称">
        <el-input v-model="localRoleForm.role" placeholder="例如：孙悟空、马里奥" class="glass-input" />
      </el-form-item>
      <el-form-item label="角色来源（可选）">
        <el-input v-model="localRoleForm.source" placeholder="例如：七龙珠、任天堂游戏" class="glass-input" />
      </el-form-item>
      <el-form-item label="执行动作">
        <el-input v-model="localRoleForm.action" placeholder="例如：正在跳舞、正在奔跑" class="glass-input" />
      </el-form-item>
      <el-form-item label="分辨率">
        <el-select v-model="localRoleForm.size" placeholder="选择分辨率" class="glass-select">
          <el-option label="624*624 (推荐)" value="624*624" />
          <el-option label="832*480" value="832*480" />
          <el-option label="480*832" value="480*832" />
          <el-option label="960*960" value="960*960" />
          <el-option label="1280*720" value="1280*720" />
          <el-option label="720*1280" value="720*1280" />
        </el-select>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

// 定义组件属性
const props = defineProps<{
  roleForm: {
    role: string,
    source: string,
    action: string,
    size: string
  }
}>()

// 定义事件
const emit = defineEmits<{
  (e: 'update:roleForm', value: { role: string, source: string, action: string, size: string }): void
}>()

// 本地响应式数据
const localRoleForm = ref({ ...props.roleForm })

// 监听属性变化
watch(() => props.roleForm, (newVal) => {
  localRoleForm.value = { ...newVal }
}, { deep: true })

// 监听本地数据变化并发出事件
watch(localRoleForm, (newVal) => {
  emit('update:roleForm', { ...newVal })
}, { deep: true })
</script>

<style scoped>
.glass-input :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.glass-input :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 15px rgba(102, 126, 234, 0.4);
}

.glass-input :deep(.el-input.is-focus .el-input__wrapper) {
  border-color: #667eea;
  box-shadow: 0 0 15px rgba(102, 126, 234, 0.4);
}

.glass-input :deep(.el-input__inner) {
  color: #fff;
  background: transparent;
}

.glass-input :deep(.el-input__inner::placeholder) {
  color: rgba(224, 224, 255, 0.6);
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
</style>