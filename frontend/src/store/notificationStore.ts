import { defineStore } from 'pinia'

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    visible: false,
    message: '',
    details: '',
    color: 'grey-darken-3',
    icon: 'mdi-information',
  }),
  actions: {
    showError(msg: string, detailsInfo: string = '') {
      this.message = msg
      this.details = detailsInfo
      this.color = '#D93025' // Google Error Red
      this.icon = 'mdi-alert-circle'
      this.visible = true
    },
    showWarning(msg: string, detailsInfo: string = '') {
      this.message = msg
      this.details = detailsInfo
      this.color = '#F29900' // Google Warning Yellow
      this.icon = 'mdi-alert'
      this.visible = true
    },
    showSuccess(msg: string) {
      this.message = msg
      this.details = ''
      this.color = '#1E8E3E' // Google Success Green
      this.icon = 'mdi-check-circle'
      this.visible = true
    }
  }
})
