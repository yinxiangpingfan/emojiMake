import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    host: '0.0.0.0',
    allowedHosts: ['emoji.easyimpr.com'],
    proxy: {
      // 代理规则按照具体到一般的顺序排列
      '/api/proxy': {
        target: 'https://82.156.59.17:8000',
        changeOrigin: true,
        secure: false,
        rewrite: (path) => {
          // 提取查询参数中的url并将其作为目标路径
          const urlParams = new URLSearchParams(path.replace('/api/proxy?', ''));
          const targetUrl = urlParams.get('url');
          if (targetUrl) {
            try {
              const urlObj = new URL(targetUrl);
              return urlObj.pathname + urlObj.search;
            } catch (e) {
              // 如果URL解析失败，返回原始路径
              return path;
            }
          }
          return path;
        }
      },
      '/api': {
        target: 'https://82.156.59.17:8000',
        changeOrigin: true,
        secure: false,
        rewrite: (path) => path.replace(/^\/api/, '/api')
      }
    }
  }
})