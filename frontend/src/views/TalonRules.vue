<template>
  <v-container fluid class="pa-4 pt-10">
    <v-card class="elevation-1 border mb-4">
      <v-card-title class="text-h6 font-weight-medium bg-grey-lighten-4 py-3 d-flex align-center">
        <v-icon color="primary" class="mr-2">mdi-shield-edit</v-icon>
        Automated Responses Maps (Talon Rules)
        <v-spacer></v-spacer>
        <v-btn color="primary" variant="outlined" size="small" class="text-none">
          Sync Rules from Cluster
        </v-btn>
      </v-card-title>
      <v-card-text class="pt-4 text-body-1 text-grey-darken-2">
        Aici se mapează fiecare alertă (Regulă Falco) pe care o generează securitatea clusterului tău către o acțiune automată de remediere SOAR prin Falco Talon. Modul "Audit" va simula acțiunea fără s-o execute.
      </v-card-text>
    </v-card>

    <v-data-table
      :headers="headers"
      :items="mockRules"
      hide-default-footer
      class="elevation-1"
    >
      <template v-slot:item.falco_rule="{ item }">
        <span class="font-weight-medium">{{ item.falco_rule }}</span>
      </template>

      <template v-slot:item.action="{ item }">
        <v-chip color="primary" variant="outlined" size="small">
          <v-icon start size="x-small">mdi-flash</v-icon> {{ item.action }}
        </v-chip>
      </template>

      <template v-slot:item.status="{ item }">
        <v-switch
          v-model="item.active"
          color="success"
          hide-details
          inset
          class="d-inline-flex ml-n2"
          :label="item.active ? 'Enforce' : 'Audit Mode'"
        ></v-switch>
      </template>

      <template v-slot:item.edit>
        <v-btn icon="mdi-pencil-outline" size="small" variant="text"></v-btn>
      </template>
    </v-data-table>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const headers = [
  { title: 'Falco Rule (Trigger)', key: 'falco_rule' },
  { title: 'Target Workload', key: 'match' },
  { title: 'Talon Parameter (Action)', key: 'action' },
  { title: 'Status / Mode', key: 'status' },
  { title: 'Configure', key: 'edit', sortable: false, align: 'end' }
]

const mockRules = ref([
  { id: 1, falco_rule: 'Terminal shell in container', match: 'k8s.ns.name=*', action: 'Terminate Pod', active: true },
  { id: 2, falco_rule: 'Write below etc', match: 'k8s.pod.labels.app=api', action: 'Label: quarantine=true', active: false },
  { id: 3, falco_rule: 'Drop and execute new binary', match: 'k8s.ns.name=production', action: 'Network Isolate', active: true },
  { id: 4, falco_rule: 'Run as Root User', match: 'k8s.ns.name!=kube-system', action: 'Network Isolate', active: false },
])
</script>
