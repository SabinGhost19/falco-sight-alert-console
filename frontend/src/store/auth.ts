import { defineStore } from 'pinia'
import axios from 'axios'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('falcosight_token') || null,
    user: localStorage.getItem('falcosight_user') || null,
  }),
  
  getters: {
    isAuthenticated: (state) => !!state.token
  },

  actions: {
    async login(username: string, password: string) {
      try {
        const response = await axios.post('http://localhost:3000/api/login', {
          username,
          password
        })

        const { token, user } = response.data
        
        // Save to state & localstorage
        this.token = token
        this.user = user
        localStorage.setItem('falcosight_token', token)
        localStorage.setItem('falcosight_user', user)

        // Set global Axios Authorization header
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        
        return true
      } catch (error) {
        console.error('Login failed', error)
        return false
      }
    },
    
    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem('falcosight_token')
      localStorage.removeItem('falcosight_user')
      delete axios.defaults.headers.common['Authorization']
    }
  }
})
