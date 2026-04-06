<template>
  <v-container fluid class="pa-4 pt-10">
    <v-card class="elevation-1 border mb-4">
      <v-card-title class="text-h6 font-weight-medium bg-grey-lighten-4 py-3 d-flex align-center">
        <v-icon color="primary" class="mr-2">mdi-shield-edit</v-icon>
        Automated Responses Maps (Talon Rules)
        <v-spacer></v-spacer>
        <v-btn color="primary" variant="outlined" size="small" class="text-none" @click="store.fetchRules()">
            <v-icon start icon="mdi-refresh" /> Sync Rules
        </v-btn>
        <v-btn color="primary" variant="flat" size="small" class="text-none ml-2" @click="openDialog">
            <v-icon start icon="mdi-plus" /> New Rule
        </v-btn>
      </v-card-title>
      <v-card-text class="pt-4 text-body-1 text-grey-darken-2">
        Aici se mapează fiecare alertă (Regulă Falco) pe care o generează securitatea clusterului tău către o acțiune automată de remediere SOAR prin Falco Talon. Modul "Audit" va simula acțiunea fără s-o execute.
      </v-card-text>
    </v-card>

    <v-data-table
      :headers="headers"
      :items="store.rules"
      :loading="store.loading"
      hide-default-footer
      class="elevation-1"
    >
      <template #item.falcoRule="{ item }">
        <span class="font-weight-medium">{{ item.falcoRule }}</span>
      </template>

      <template #item.action="{ item }">
        <v-chip color="primary" variant="outlined" size="small">
          <v-icon start size="x-small">mdi-flash</v-icon> {{ item.action }} {{ item.actionDetails ? `(${item.actionDetails})` : '' }}
        </v-chip>
      </template>

      <template #item.status="{ item }">
        <v-switch
          :model-value="item.enabled"
          @update:model-value="store.toggleRule(item.ID)"
          color="success"
          hide-details
          inset
          class="d-inline-flex ml-n2"
          :label="item.enabled ? 'Enforce' : 'Audit Mode'"
        ></v-switch>
      </template>

      <template #item.actions="{ item }">
        <v-btn icon="mdi-delete-outline" color="error" size="small" variant="text" @click="deleteRule(item.ID)"></v-btn>
      </template>
    </v-data-table>

    <-> Dialog for new rule -->
    <v-dialog v-model="dialog" max-width="500">
      <v-card>
        <v-card-title>Create Talon Rule</v-card-title>
        <v-card-text>
          <v-text-field v-model="newRule.name" label="Rule Name" required></v-text-field>
          <v-text-field v-model="newRule.falcoRule" label="Falco Rule Trigger" required></v-text-field>
          <v-select v-model="newRule.action" :items="['terminate_pod', 'network_isolate', 'label_pod']" label="Action"></v-select>
          <v-text-field v-if="newRule.action === 'label_pod'" v-model="newRule.actionDetails" label="Action Details (e.g. key=val)"></v-text-field>
          <v-switch v-model="newRule.enabled" label="Enabled"></v-switch>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text="Close" @click="dialog = false"></v-btn>
          <v-btn color="primary" @click="saveRule" variant="flat">Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRulesStore } from '../store/rules'

const store = useRulesStore()
const dialog = ref(false)

const newRule = ref({
    name: '',
    falcoRule: '',
    action: 'terminate_pod',
    actionDetails: '',
    enabled: true
})

const headers = [
  { title: 'Rule Name', key: 'name' },
  { title: 'Falco Rule (Trigger)', key: 'falcoRule' },
  { title: 'Talon Parameter (Action)', key: 'action' },
  { title: 'Status / Mode', key: 'status' },
  { title: 'Actions', key: 'actions', sortable: false, align: 'end' as const }
]

onMounted(() => {
    store.fetchRules()
})

const openDialog = () => {
    newRule.value = { name: '', falcoRule: '', action: 'terminate_pod', actionDetails: '', enabled: true }
    dialog.value = true
}

const saveRule = async () => {
    await store.createRule(newRule.value)
    dialog.value = false
}

const deleteRule = async (id: number) => {
    await store.deleteRule(id)
}
</script>
