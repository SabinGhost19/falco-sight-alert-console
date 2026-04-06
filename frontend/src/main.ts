import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import vuetify from './plugins/vuetify'
import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import axios from 'axios'
import { useNotificationStore } from './store/notificationStore'

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
app.use(router)
app.use(vuetify)

// Inregistram interceptorul global pt Axios DOAR DUPA initializarea Pinia
const notify = useNotificationStore(pinia)

axios.interceptors.response.use(
  (response) => response,
  (error) => {
    // 1. Fara retea / Backend off
    if (!error.response) {
      notify.showError("Network Error: Nu se poate contacta serverul FalcoSight. Verificați conexiunea VPN sau starea K8s.", "ERR_CONNECTION_REFUSED")
      return Promise.reject(error)
    }

    const status = error.response.status
    const errorData = error.response.data?.error

    // 2. Structura noastra Enterprise-Grade ErrorHandling
    if (errorData) {
      if (status === 403 || status === 401) {
         notify.showWarning(`Acces Respins: ${errorData.message} (${errorData.code})`)
      } else if (status >= 500) {
         notify.showError(`Eroare Platformă: ${errorData.message}`, errorData.details)
      } else {
         notify.showWarning(errorData.message, errorData.details)
      }
    } 
    // 3. Erori venite din Proxy (Nginx) sau timeout
    else {
      if (status === 502 || status === 503 || status === 504) {
         notify.showError("Serviciul de Backend K8s este temporar indisponibil.", `50x Proxy Failure`)
      } else {
         notify.showError(`Eroare HTTP Neașteptată (${status})`, error.message)
      }
    }

    return Promise.reject(error)
  }
)

// Timeout global siguranta pt Axios
axios.defaults.timeout = 10000

app.mount('#app')
