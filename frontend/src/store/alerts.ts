import { defineStore } from 'pinia'
import axios from 'axios'

export interface Alert {
  id: number
  created_at: string
  priority: string
  rule: string
  message: string
  namespace: string
  pod_name: string
  container_name: string
  manifest_yaml: string
  vulnerable_lines: string
  talon_action: string
  talon_status: string
}

export const useAlertsStore = defineStore('alerts', {
  state: () => ({
    alerts: [] as Alert[],
    loading: false
  }),
  
  actions: {
    async fetchAlerts() {
      this.loading = true
      try {
        const response = await axios.get('/api/v1/alerts')
        this.alerts = response.data
      } catch (error) {
        console.error('Error fetching alerts:', error)
      } finally {
        this.loading = false
      }
    },
    
    async triggerTalon(alertId: number, actionName: string) {
      try {
        await axios.post('/api/v1/talon/trigger', {
          alert_id: alertId,
          action: actionName
        })
        const alert = this.alerts.find(a => a.id === alertId)
        if (alert) {
          alert.talon_action = actionName
          alert.talon_status = 'Requested'
        }
      } catch (error) {
        console.error('Error triggering talon:', error)
      }
    }
  }
})
