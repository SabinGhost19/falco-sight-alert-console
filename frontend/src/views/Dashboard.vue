<template>
  <v-container fluid class="pa-4 pt-10">
    <v-row>
      <!-- Summary Cards -->
      <v-col cols="12" md="3">
        <v-card class="elevation-1 rounded-sm">
          <v-card-text>
            <div class="text-overline text-grey-darken-1 mb-1">TOTAL ALERTS</div>
            <div class="text-h3 font-weight-medium">{{ store.alerts.length }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card class="elevation-1 rounded-sm">
          <v-card-text>
            <div class="text-overline text-grey-darken-1 mb-1">CRITICAL THREATS</div>
            <div class="text-h3 font-weight-medium text-error">{{ highPriority.length }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card class="elevation-1 rounded-sm">
          <v-card-text>
            <div class="text-overline text-grey-darken-1 mb-1">REMEDIATED BY FALCO TALON</div>
            <div class="text-h3 font-weight-medium text-success">{{ talonRemediated }}</div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Workload Security Posture Table -->
    <v-row class="mt-6">
      <v-col cols="12">
         <v-sheet class="elevation-1 rounded-sm pa-0">
            <v-toolbar color="white" flat class="border-b px-2">
              <v-toolbar-title class="text-subtitle-1 font-weight-medium text-grey-darken-3">
                Security Alerts Explorer
              </v-toolbar-title>
              <v-spacer></v-spacer>
              <v-text-field
                v-model="search"
                prepend-inner-icon="mdi-magnify"
                label="Filter by Rule, Namespace..."
                single-line
                hide-details
                density="compact"
                variant="outlined"
                class="narrow-search"
                color="primary"
              ></v-text-field>
            </v-toolbar>

            <v-data-table
              :headers="headers"
              :items="store.alerts"
              :search="search"
              :loading="store.loading"
              items-per-page="10"
              hover
              @click:row="openIncidentContext"
            >
              <!-- Priority Column -->
              <template v-slot:item.priority="{ item }">
                <v-chip
                  :color="getPriorityColor(item.priority)"
                  size="small"
                  label
                  class="font-weight-bold"
                >
                  {{ item.priority }}
                </v-chip>
              </template>

              <!-- Workload Column -->
              <template v-slot:item.pod_name="{ item }">
                <div class="d-flex align-center">
                  <v-icon color="#326CE5" size="small" class="mr-2">mdi-kubernetes</v-icon>
                  <span class="font-weight-medium">{{ item.namespace }}/</span>
                  {{ item.pod_name }}
                </div>
              </template>

              <!-- Talon Action -->
              <template v-slot:item.talon_status="{ item }">
                <v-chip v-if="item.talon_status" size="small" variant="outlined" color="success">
                  <v-icon start size="x-small">mdi-robot-outline</v-icon>
                  {{ item.talon_action }} ({{ item.talon_status }})
                </v-chip>
                <span v-else class="text-caption text-grey">Unmanaged</span>
              </template>

              <!-- Actions -->
              <template v-slot:item.actions="{ item }">
                <v-btn icon="mdi-open-in-new" variant="text" size="small" color="primary"></v-btn>
              </template>
            </v-data-table>
        </v-sheet>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAlertsStore } from '../store/alerts'

const store = useAlertsStore()
const router = useRouter()
const search = ref('')

const headers = [
  { title: 'Severity', key: 'priority', width: '100px' },
  { title: 'Time', key: 'created_at' },
  { title: 'Rule Triggered', key: 'rule' },
  { title: 'Workload (NS/Pod)', key: 'pod_name' },
  { title: 'Talon Remediations', key: 'talon_status' },
  { title: '', key: 'actions', sortable: false, align: 'end' }
]

const highPriority = computed(() => store.alerts.filter(a => ['Critical', 'Emergency'].includes(a.priority) || a.priority === 'High'))
const talonRemediated = computed(() => store.alerts.filter(a => a.talon_status === 'Success' || a.talon_status === 'Requested').length)

const getPriorityColor = (priority: string) => {
  switch (priority?.toLowerCase()) {
    case 'critical': return 'error'
    case 'error': return 'error'
    case 'warning': return 'warning'
    case 'notice': return 'info'
    default: return 'primary'
  }
}

const openIncidentContext = (event: Event, { item }: { item: any }) => {
  router.push({ name: 'IncidentContext', params: { id: item.id } })
}

onMounted(() => {
  store.fetchAlerts()
})
</script>

<style scoped>
.narrow-search {
  max-width: 300px;
}
:deep(.v-data-table) {
  font-size: 14px;
}
:deep(.v-data-table-header__content) {
  font-weight: 600 !important;
  color: #5F6368 !important;
}
</style>
