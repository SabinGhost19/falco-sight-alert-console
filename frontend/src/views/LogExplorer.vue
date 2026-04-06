<template>
  <v-container fluid class="pa-4 pt-10">
    <!-- Log Explorer Advanced Search -->
    <v-card class="mb-4">
      <v-card-text>
         <v-row align="center">
           <v-col cols="12" md="10">
              <v-text-field
                v-model="searchQuery"
                prepend-inner-icon="mdi-magnify"
                label="Advanced Search (e.g. rule='Drop and execute new binary' AND namespace!='kube-system')"
                variant="outlined"
                density="compact"
                hide-details
                color="primary"
                clearable
              ></v-text-field>
           </v-col>
           <v-col cols="12" md="2">
              <v-btn color="primary" block class="text-none" prepend-icon="mdi-filter-variant">Search Logs</v-btn>
           </v-col>
         </v-row>
      </v-card-text>
    </v-card>

    <!-- Expanding Data Table -->
    <v-data-table
      :headers="headers"
      :items="filteredAlerts"
      :search="searchQuery"
      :loading="store.loading"
      item-value="id"
      show-expand
      class="elevation-1"
    >
      <template v-slot:item.priority="{ item }">
        <v-chip :color="item.priority === 'Critical' ? 'error' : 'warning'" size="small" label>
          {{ item.priority }}
        </v-chip>
      </template>

      <!-- Expanded Row with Raw JSON -->
      <template v-slot:expanded-row="{ columns, item }">
        <tr>
          <td :colspan="columns.length" class="bg-grey-darken-4 pa-4">
            <div class="d-flex justify-space-between mb-2 text-white">
              <strong>Raw Event JSON:</strong>
              <v-btn size="x-small" color="primary" variant="text" prepend-icon="mdi-content-copy">Copy Payload</v-btn>
            </div>
            <pre class="text-caption text-light-green-accent-2" style="white-space: pre-wrap; word-break: break-all;">
{{ formatRaw(item) }}
            </pre>
          </td>
        </tr>
      </template>
    </v-data-table>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAlertsStore } from '../store/alerts'

const store = useAlertsStore()
const searchQuery = ref('')

const headers = [
  { title: 'Time', key: 'created_at' },
  { title: 'Priority', key: 'priority' },
  { title: 'Namespace', key: 'namespace' },
  { title: 'Pod', key: 'pod_name' },
  { title: 'Rule', key: 'rule' },
  { title: 'Message', key: 'message' },
]

const filteredAlerts = computed(() => store.alerts)

const formatRaw = (item: any) => {
  return JSON.stringify(item, null, 2)
}

onMounted(() => {
  store.fetchAlerts()
})
</script>
