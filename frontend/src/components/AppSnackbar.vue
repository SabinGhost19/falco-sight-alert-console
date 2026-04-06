<template>
  <v-snackbar
    v-model="notify.visible"
    :color="notify.color"
    location="top right"
    timeout="6000"
    variant="flat"
    class="mt-12"
  >
    <div class="d-flex align-start">
      <v-icon class="mr-3 mt-1">{{ notify.icon }}</v-icon>
      <div>
        <div class="text-subtitle-2 font-weight-bold">{{ notify.message }}</div>
        <div v-if="notify.details" class="text-caption mt-1 opacity-80" style="font-family: monospace;">
            {{ notify.details }}
        </div>
      </div>
    </div>
    <template v-slot:actions>
      <v-btn icon="mdi-close" variant="text" size="small" @click="notify.visible = false"></v-btn>
    </template>
  </v-snackbar>
</template>

<script setup lang="ts">
import { useNotificationStore } from '../store/notificationStore'
const notify = useNotificationStore()
</script>
