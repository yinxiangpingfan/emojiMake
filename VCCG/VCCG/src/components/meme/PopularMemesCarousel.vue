<template>
  <div class="popular-section">
    <h2 class="section-title">热门网络表情包</h2>
    <el-carousel :interval="4000" type="card" height="300px" class="popular-carousel soft-edges">
      <el-carousel-item v-for="meme in popularMemes" :key="meme.id">
        <img :src="meme.url" :alt="meme.name" class="carousel-image">
      </el-carousel-item>
    </el-carousel>
  </div>
</template>

<script setup lang="ts">
// 热门表情包轮播数据
const popularMemes = [
  { id: 1, name: 'Doge', url: new URL('../../assets/images/popular1.gif', import.meta.url).href },
  { id: 2, name: 'Woman Yelling At Cat', url: new URL('../../assets/images/popular2.gif', import.meta.url).href },
  { id: 3, name: 'This Is Fine', url: new URL('../../assets/images/popular3.gif', import.meta.url).href },
  { id: 4, name: 'Drake Hotline Bling', url: new URL('../../assets/images/popular4.gif', import.meta.url).href },
  { id: 5, name: 'Two Buttons', url: new URL('../../assets/images/popular5.gif', import.meta.url).href },
  { id: 6, name: 'Running Away Balloon', url: new URL('../../assets/images/popular6.gif', import.meta.url).href }
]
</script>

<style scoped>
.popular-section {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 15px;
  padding: 30px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  margin-bottom: 30px;
}

.popular-section .section-title {
  text-align: center;
  color: #e0e0ff;
}

.soft-edges {
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.2);
}

.soft-edges :deep(.el-carousel__item) {
  border-radius: 15px;
  overflow: hidden;
  border: none;
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.2);
}

.carousel-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 15px;
}
</style>