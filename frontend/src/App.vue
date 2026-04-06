<template>
  <v-app>
    <-> Top Navigation Bar is visible only if authenticated -->
    <v-app-bar color="#202124" density="compact" elevation="2" v-if="authStore.isAuthenticated">
      <v-container class="d-flex align-center py-0 px-2 mx-auto" fluid>
        
        <!-- Custom FalcoSight Logo (SVG inline) -->
        <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#4285F4" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-3">
          <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
          <path d="m9 12 2 2 4-4"/>
        </svg>

        <v-app-bar-title class="text-white text-h6 font-weight-regular" style="cursor: pointer" @click="router.push('/')">
          FalcoSight <span class="text-caption text-grey-lighten-1">Security Observability</span>
        </v-app-bar-title>

        <v-spacer></v-spacer>

        <!-- Main Navigation Links -->
        <v-btn variant="text" class="text-none ml-2 text-white" to="/" exact>Dashboard</v-btn>
        <v-btn variant="text" class="text-none ml-2 text-white" to="/logs">Log Explorer</v-btn>
        <v-btn variant="text" class="text-none ml-2 text-white" to="/rules">Talon Rules</v-btn>

        <!-- Project/Cluster Selector -->
        <v-btn variant="outlined" prepend-icon="mdi-kubernetes" size="small" color="primary" class="text-none ml-6 text-white border-opacity-50">
          Cluster: microk8s-prod
        </v-btn>

        <-> User Profile Avatar & Logout -->
        <v-menu>
          <template v-slot:activator="{ props }">
            <v-chip
              v-bind="props"
              color="#4285F4"
              class="ml-4 pl-1 cursor-pointer"
              variant="flat"
              prepend-icon="mdi-account-circle"
            >
              {{ authStore.user || 'Admin' }}
            </v-chip>
          </template>
          <v-list density="compact" elevation="2">
            <v-list-item @click="handleLogout">
              <template v-slot:prepend><v-icon color="error" size="small">mdi-logout</v-icon></template>
              <v-list-item-title class="text-error font-weight-medium">Logout</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
        
      </v-container>
    </v-app-bar>

    <v-main class="bg-background">
      <router-view></router-view>
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import { useAuthStore } from './store/auth'
import { useRouter } from 'vue-router'
import { onMounted } from 'vue'
import axios from 'axios'

const authStore = useAuthStore()
const router = useRouter()

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

// Ensure the token is attached across page reloads
onMounted(() => {
  if (authStore.token) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${authStore.token}`
  }
})
</script>

<style>
/* Clean typography baseline similar to Roboto used in Google Cloud Console */
body {
  font-family: Roboto, "Helvetica Neue", sans-serif;
  margin: 0;
  background-color: #F8F9FA;
}
.bg-background {
  background-color: #F8F9FA !important;
}
</style>
