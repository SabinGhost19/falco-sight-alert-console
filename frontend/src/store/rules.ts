import { defineStore } from 'pinia'
import axios from 'axios'

export interface TalonRule {
  ID: number
  name: string
  falcoRule: string
  action: string
  actionDetails: string
  enabled: boolean
}

export const useRulesStore = defineStore('rules', {
  state: () => ({
    rules: [] as TalonRule[],
    loading: false
  }),
  actions: {
    async fetchRules() {
      this.loading = true
      try {
        const res = await axios.get('/api/v1/rules')
        this.rules = res.data
      } catch (err) {
        console.error('Error fetching rules:', err)
      } finally {
        this.loading = false
      }
    },
    async toggleRule(id: number) {
      try {
        const res = await axios.patch(`/api/v1/rules/${id}/toggle`)
        const idx = this.rules.findIndex(r => r.ID === id)
        if (idx !== -1) {
          this.rules[idx].enabled = res.data.enabled
        }
      } catch (err) {
        console.error('Error toggling rule:', err)
      }
    },
    async deleteRule(id: number) {
        try {
          await axios.delete(`/api/v1/rules/${id}`)
          this.rules = this.rules.filter(r => r.ID !== id)
        } catch (err) {
          console.error('Error deleting rule:', err)
        }
    },
    async createRule(payload: Partial<TalonRule>) {
      try {
        const res = await axios.post('/api/v1/rules', payload)
        this.rules.unshift(res.data)
      } catch (err) {
        console.error('Error creating rule:', err)
      }
    }
  }
})
