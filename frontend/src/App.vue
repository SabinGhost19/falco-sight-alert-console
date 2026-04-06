<template>
  <v-app>
    <-> Top Navigation Bar -->
    <v-app-bar color="#FFFFFF" elevation="0" border="b" v-if="authStore.isAuthenticated" class="gc-border" density="compact">
      <div class="px-4 d-flex align-center w-100">
        <v-btn icon="mdi-menu" variant="text" size="small" color="#5F6368" @click="drawer = !drawer"></v-btn>

        <-> Brand -->
        <img src="@/assets/logo.png" alt="Falco Logo" class="ml-2 mr-2" style="height: 26px; cursor: pointer" @click="router.push('/')" />
        <span class="font-weight-medium text-grey-darken-3 text-subtitle-1" style="cursor: pointer" @click="router.push('/')">
          FalcoSight
        </span>
        <span class="ml-2 text-grey text-subtitle-2 font-italic border-l pl-2 border-opacity-50">Command Center</span>

        <v-spacer></v-spacer>

        <-> Project/Cluster Selector (GCP Style) -->
        <v-menu>
          <template v-slot:activator="{ props }">
            <v-btn v-bind="props" variant="text" class="text-none ml-4 text-primary font-weight-medium" append-icon="mdi-menu-down" prepend-icon="mdi-cloud">
              production-cluster-01
            </v-btn>
          </template>
          <v-list density="compact" class="gc-border">
            <v-list-item><v-list-item-title class="text-caption">sandbox-dev</v-list-item-title></v-list-item>
            <v-list-item><v-list-item-title class="text-caption font-weight-bold">production-cluster-01</v-list-item-title></v-list-item>
          </v-list>
        </v-menu>

        <-> Helper Icons (search, help, notifications) -->
        <v-btn icon="mdi-magnify" variant="text" size="small" color="#5F6368" class="ml-2"></v-btn>
        <v-btn icon="mdi-help-circle-outline" variant="text" size="small" color="#5F6368"></v-btn>
        
        <-> User Profile Avatar & Logout -->
        <v-menu>
          <template v-slot:activator="{ props }">
            <div class="ml-4 pl-1 align-center d-flex cursor-pointer" v-bind="props">
               <v-avatar size="32" color="blue-lighten-4">
                  <span class="text-primary text-subtitle-2 font-weight-bold">{{ (authStore.user || 'Admin').charAt(0) }}</span>
               </v-avatar>
            </div>
          </template>
          <v-list density="compact" class="gc-border">
            <v-list-item @click="handleLogout">
              <template v-slot:prepend><v-icon color="error" size="small">mdi-logout</v-icon></template>
              <v-list-item-title class="text-error font-weight-medium">Sign out</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </div>
    </v-app-bar>

    <-> Sidebar navigation (Navigation Drawer) -->
    <v-navigation-drawer v-model="drawer" v-if="authStore.isAuthenticated" permanent elevation="0" border="r" color="#F8F9FA" class="gc-border">
      <v-list density="compact" nav class="px-2 pt-4">
        <v-list-subheader class="text-uppercase text-caption font-weight-bold text-grey-darken-1 mb-2">Security Posture</v-list-subheader>
        
        <v-list-item to="/" prepend-icon="mdi-view-dashboard-outline" color="primary" rounded="sm" exact>
          <v-list-item-title class="font-weight-medium text-body-2">Overview</v-list-item-title>
        </v-list-item>

        <v-list-item to="/logs" prepend-icon="mdi-alert-octagon-outline" color="primary" rounded="sm">
          <v-list-item-title class="font-weight-medium text-body-2">Alerts Explorer</v-list-item-title>
        </v-list-item>

        <v-divider class="my-4"></v-divider>
        <v-list-subheader class="text-uppercase text-caption font-weight-bold text-grey-darken-1 mb-2">Response Automation</v-list-subheader>

        <v-list-item to="/rules" prepend-icon="mdi-robot-outline" color="primary" rounded="sm">
          <v-list-item-title class="font-weight-medium text-body-2">Talon Actions (SOAR)</v-list-item-title>
        </v-list-item>
        
        <v-list-item prepend-icon="mdi-shield-check-outline" color="primary" rounded="sm" href="#" target="_blank">
          <v-list-item-title class="font-weight-medium text-body-2">Policies</v-list-item-title>
        </v-list-item>

      </v-list>
    </v-navigation-drawer>

    <-> Main Content Area -->
    <!-- Notificări globale -->
    <AppSnackbar />

    <v-main class="bg-background">
      <div class="main-container h-100">
         <router-view></router-view>
      </div>
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import { useAuthStore } from './store/auth'
import { useRouter } from 'vue-router'
import { onMounted, ref } from 'vue'
import axios from 'axios'
import AppSnackbar from './components/AppSnackbar.vue'

const authStore = useAuthStore()
const router = useRouter()
const drawer = ref(true)

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
@import url('https://fonts.googleapis.com/css2?family=Roboto:wght@300;400;500;700&display=swap');

body {
  font-family: 'Roboto', "Helvetica Neue", sans-serif;
  margin: 0;
  -webkit-font-smoothing: antialiased;
}
.bg-background {
  background-color: #F8F9FA !important;
}
/* Google Cloud border standard */
.gc-border {
  border-color: #DADCE0 !important;
}

.main-container {
  max-width: 1600px;
  margin: 0 auto;
}

/* Material Design specific link selections */
.v-list-item--active {
  background-color: #E8F0FE;
  color: #1A73E8;
}

.v-list-item--active .v-icon {
  color: #1A73E8 !important;
}
</style>
