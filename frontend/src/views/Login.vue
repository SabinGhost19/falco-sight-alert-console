<template>
  <v-container class="fill-height bg-background d-flex align-center justify-center">
    <v-card class="elevation-4 pb-4" max-width="400" width="100%" rounded="lg">
      <v-card-title class="d-flex justify-center pt-6 pb-2">
        <v-icon color="primary" size="large" class="mr-2">mdi-shield-check</v-icon>
        <span class="text-h5 font-weight-medium text-grey-darken-3">FalcoSight</span>
      </v-card-title>
      
      <v-card-text class="text-center text-body-2 text-grey-darken-1 mb-4">
        Sign in to view Kubernetes Security Alerts
      </v-card-text>

      <v-card-text class="px-6">
        <v-form @submit.prevent="handleLogin" ref="form">
          <v-text-field
            v-model="username"
            label="Username"
            variant="outlined"
            density="comfortable"
            prepend-inner-icon="mdi-account"
            color="primary"
            class="mb-2"
          ></v-text-field>

          <v-text-field
            v-model="password"
            label="Password"
            type="password"
            variant="outlined"
            density="comfortable"
            prepend-inner-icon="mdi-lock"
            color="primary"
            @keyup.enter="handleLogin"
          ></v-text-field>

          <v-alert v-if="error" type="error" variant="tonal" density="compact" class="mt-2 mb-4">
            {{ error }}
          </v-alert>

          <v-btn
            type="submit"
            color="primary"
            size="large"
            block
            class="mt-4 text-none font-weight-medium"
            :loading="loading"
            elevation="1"
          >
            Access Dashboard
          </v-btn>
        </v-form>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'

const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const router = useRouter()
const authStore = useAuthStore()

const handleLogin = async () => {
  if (!username.value || !password.value) {
    error.value = "Please enter both fields"
    return
  }

  error.value = ''
  loading.value = true

  const success = await authStore.login(username.value, password.value)
  loading.value = false

  if (success) {
    router.push('/')
  } else {
    error.value = "Invalid username or password"
  }
}
</script>
