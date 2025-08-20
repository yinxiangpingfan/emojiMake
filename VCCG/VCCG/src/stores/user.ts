import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    isLoggedIn: false,
    userInfo: {
      phone: ''
    },
    token: ''
  }),
  
  actions: {
    login(token: string, userInfo: any) {
      this.isLoggedIn = true
      this.token = token
      this.userInfo = userInfo
      localStorage.setItem('userToken', token)
      // 保存用户信息到localStorage
      localStorage.setItem('userInfo', JSON.stringify(userInfo))
    },
    
    logout() {
      this.isLoggedIn = false
      this.token = ''
      this.userInfo = {
        phone: ''
      }
      localStorage.removeItem('userToken')
      localStorage.removeItem('userInfo')
    },
    
    checkLoginStatus() {
      const token = localStorage.getItem('userToken')
      const userInfo = localStorage.getItem('userInfo')
      if (token) {
        this.isLoggedIn = true
        this.token = token
        // 恢复用户信息
        if (userInfo) {
          try {
            this.userInfo = JSON.parse(userInfo)
          } catch (e) {
            console.error('解析用户信息失败:', e)
            // 如果解析失败，使用默认值
            this.userInfo = {
              phone: ''
            }
          }
        }
      }
    }
  },
  
  getters: {
    isAuthenticated: (state) => state.isLoggedIn && !!state.token
  }
})