<template>
  <v-container fluid class="pa-6">
    <div class="text-h5 font-weight-medium mb-1 text-high-emphasis">Alerts Explorer</div>
    <div class="text-body-2 text-medium-emphasis mb-6">Perform Threat Hunting across your generated Kubernetes security incidents.</div>

    <v-card class="gc-border pb-0" color="surface">
      <v-toolbar color="surface" flat class="border-b gc-border px-2">
         <v-toolbar-title class="text-subtitle-1 font-weight-medium text-high-emphasis">
            Security Alerts
         </v-toolbar-title>
         <v-spacer></v-spacer>
         <v-text-field
           v-model="search"
           prepend-inner-icon="mdi-magnify"
           label="Filter by rule, namespace, or pod..."
           single-line
           hide-details
           density="compact"
           variant="outlined"
           style="max-width: 350px;"
           color="primary"
         ></v-text-field>
      </v-toolbar>

      <v-data-table
        :headers="headers"
        :items="store.alerts"
        :search="search"
        :loading="store.loading"
        items-per-page="15"
        hover
        density="comfortable"
        @click:row="openIncidentContext"
      >
        <-> Severity -->
        <template v-slot:item.priority="{ item }">
           <v-chip :color="getPriorityColor(item.priority)" size="small" variant="flat" class="font-weight-bold text-white">
               {{ item.priority }}
           </v-chip>
        </template>
        
        <-> Timestamp Formatted -->
        <template v-slot:item.created_at="{ item }">
           <span class="text-body-2 text-medium-emphasis">{{ formatDate(item.created_at) }}</span>
        </template>

        <-> Rule Name -->
        <template v-slot:item.rule="{ item }">
           <span class="font-weight-medium text-high-emphasis">{{ item.rule }}</span>
        </template>

        <-> Workload Details -->
        <template v-slot:item.pod_name="{ item }">
           <div class="d-flex align-center">
              <v-icon color="#326CE5" size="small" class="mr-2">mdi-kubernetes</v-icon>
              <span class="font-weight-bold text-primary">{{ item.namespace }}</span>
              <span class="text-medium-emphasis mx-1">/</span>
              <span class="text-body-2">{{ item.pod_name }}</span>
           </div>
        </template>

        <-> Talon Action Status -->
        <template v-slot:item.talon_status="{ item }">
           <div v-if="item.talon_status" class="d-flex align-center">
              <v-icon size="small" :color="item.talon_status === 'Success' ? 'success' : 'warning'" class="mr-1">mdi-robot-outline</v-icon>
              <span class="text-caption font-weight-medium">{{ item.talon_action }} ({{ item.talon_status }})</span>
           </div>
           <span v-else class="text-caption text-grey">Unmanaged</span>
        </template>
      </v-data-table>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAlertsStore } from '../store/alerts'
import { formatDistanceToNow, parseISO } from 'date-fns'

const store = useAlertsStore()
const router = useRouter()
const route = useRoute()
const search = ref('')

const headers = [
  { title: 'Severity', key: 'priority', width: '110px' },
  { title: 'Time Detected', key: 'created_at', width: '180px' },
  { title: 'Rule Name', key: 'rule' },
  { title: 'Workload (NS / Pod)', key: 'pod_name' },
  { title: 'Talon Remediations', key: 'talon_status' }
]

// Priority Color Mapper for Google Theme
const getPriorityColor = (priority: string) => {
  switch (priority?.toLowerCase()) {
    case 'emergency': return '#B31412' // Darker red
    case 'critical': return '#D93025' // Google Red
    case 'error': return '#D93025'
    case 'warning': return '#F29900' // Google Yellow
    case 'notice': return '#1A73E8'   // Google Blue
    case 'info': return '#5F6368'     // Google Grey
    default: return '#1A73E8'
  }
}

// Date formatter human readable
const formatDate = (isoString: string) => {
   if(!isoString) return 'Unknown'
   try {
      return formatDistanceToNow(parseISO(isoString), { addSuffix: true })
   } catch(e) {
      return isoString
   }
}

const openIncidentContext = (_event: Event, { item }: { item: any }) => {
  router.push({ name: 'IncidentContext', params: { id: item.id } })
}

// Prefill search if query exists
watch(() => route.query.q, (newQ) => {
   if(newQ) search.value = newQ as string
}, { immediate: true })

onMounted(() => {
  store.fetchAlerts()
  if (route.query.q) search.value = route.query.q as string
})
</script>

<style scoped>
:deep(.v-data-table) {
  font-size: 14px;
}
:deep(.v-data-table-header__content) {
  font-weight: 600 !important;
  color: #5F6368 !important;
  text-transform: uppercase;
  font-size: 11px;
  letter-spacing: 0.5px;
}
</style>
