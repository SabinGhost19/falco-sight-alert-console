<template>
  <v-container fluid class="pa-0 h-100 fill-height bg-background">
    <!-- Split Screen Layout -->
    <v-row no-gutters class="fill-height pb-0 mb-0">
      
      <!-- Left Panel: Context & Falco Alert Info -->
      <v-col cols="12" md="4" class="border-e bg-white pa-4 overflow-auto fill-height">
        <v-btn variant="text" prepend-icon="mdi-arrow-left" class="mb-4 text-none text-primary pa-0" @click="router.push('/')">
          Back to Incidents
        </v-btn>

        <template v-if="alert">
          <div class="d-flex align-center mb-2">
            <v-chip :color="priorityColor" label size="small" class="mr-2 font-weight-bold">{{ alert.priority }}</v-chip>
            <span class="text-caption text-grey">{{ new Date(alert.created_at).toLocaleString() }}</span>
          </div>

          <h2 class="text-h6 font-weight-medium text-grey-darken-3 mb-1">{{ alert.rule }}</h2>
          <p class="text-body-2 text-grey-darken-2 mb-4">{{ alert.message }}</p>

          <v-divider class="mb-4"></v-divider>

          <!-- Corelate K8s Workload Details -->
          <h3 class="text-subtitle-2 font-weight-bold text-grey-darken-3 mb-2 text-uppercase">Correlated Workload</h3>
          <v-list density="compact" class="bg-transparent pa-0">
            <v-list-item>
              <template v-slot:prepend>
                <v-icon color="#326CE5" size="small">mdi-kubernetes</v-icon>
              </template>
              <v-list-item-title class="font-weight-medium">{{ alert.pod_name }}</v-list-item-title>
              <v-list-item-subtitle>Namespace: {{ alert.namespace }}</v-list-item-subtitle>
            </v-list-item>
            <v-list-item v-if="alert.container_name">
               <template v-slot:prepend>
                <v-icon color="grey" size="small">mdi-docker</v-icon>
              </template>
              <v-list-item-title>{{ alert.container_name }}</v-list-item-title>
              <v-list-item-subtitle>Container</v-list-item-subtitle>
            </v-list-item>
          </v-list>

          <v-divider class="my-4"></v-divider>

          <!-- Falco Talon Action Center -->
          <h3 class="text-subtitle-2 font-weight-bold text-grey-darken-3 mb-3 text-uppercase">Talon Automated Response</h3>
          
          <v-card variant="outlined" :color="alert.talon_status ? 'success' : 'grey-lighten-2'" class="mb-4">
            <v-card-text v-if="alert.talon_status">
               <div class="d-flex align-center text-success font-weight-bold mb-1">
                 <v-icon size="small" start>mdi-check-circle</v-icon> Acted automatically by Talon
               </div>
               <div class="text-body-2">Action: <strong>{{ alert.talon_action }}</strong> via Talon Rules.</div>
               <div class="text-caption mt-1">Status: {{ alert.talon_status }}</div>
            </v-card-text>
            <v-card-text v-else>
               <div class="d-flex align-center text-grey-darken-1 font-weight-bold mb-1">
                 <v-icon size="small" start>mdi-alert-circle-outline</v-icon> No Automatic Action Taken
               </div>
               <div class="text-body-2 mb-3">Talon rules did not match or were running in audit mode.</div>
               
               <!-- Manual Overrides -->
               <v-btn color="primary" class="mb-2 text-none" block variant="flat" size="small" @click="triggerAction('network_isolate')" :disabled="!alert">
                  Execute Network Isolate
               </v-btn>
               <v-btn color="error" variant="outlined" block class="text-none" size="small" @click="triggerAction('terminate')" :disabled="!alert">
                  Terminate Pod
               </v-btn>
            </v-card-text>
          </v-card>
        </template>
      </v-col>

      <!-- Right Panel: Code Viewer (Monaco Editor) & Fixes -->
      <v-col cols="12" md="8" class="bg-grey-darken-4 d-flex flex-column fill-height">
        <v-toolbar color="#202124" flat dense class="border-b" style="border-bottom-color: #3C4043 !important;">
          <v-tabs v-model="tab" color="primary" bg-color="transparent" slider-color="primary">
             <v-tab value="raw" class="text-none">Live YAML Manifest</v-tab>
             <v-tab value="diff" class="text-none ml-4"><v-icon start size="small">mdi-auto-fix</v-icon> Proposed Security Fix</v-tab>
          </v-tabs>
        </v-toolbar>

        <div class="flex-grow-1 position-relative">
            <vue-monaco-editor
              v-if="tab === 'raw'"
              v-model:value="code"
              theme="vs-dark"
              language="yaml"
              :options="editorOptions"
              @mount="handleRawEditorMount"
            ></vue-monaco-editor>

             <!-- Diff Editor would go here, currently using regular editor disabled for simplicity -->
            <vue-monaco-editor
              v-if="tab === 'diff'"
              v-model:value="fixedCode"
              theme="vs-dark"
              language="yaml"
              :options="{...editorOptions, readOnly: true}"
            ></vue-monaco-editor>
        </div>
      </v-col>

    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAlertsStore } from '../store/alerts'
import { VueMonacoEditor } from '@guolao/vue-monaco-editor'

const route = useRoute()
const router = useRouter()
const store = useAlertsStore()

const id = parseInt(route.params.id as string)

const alert = computed(() => store.alerts.find(a => a.id === id))
const priorityColor = computed(() => {
  const p = alert.value?.priority?.toLowerCase()
  if (p === 'critical' || p === 'error') return 'error'
  if (p === 'warning') return 'warning'
  return 'primary'
})

const tab = ref('raw')
const code = ref('')
const fixedCode = ref('')
const monacoEditorInstance = ref<any>(null)

const editorOptions = {
  automaticLayout: true,
  minimap: { enabled: false },
  readOnly: true,
  scrollBeyondLastLine: false,
  wordWrap: 'on',
  fontSize: 13,
  padding: { top: 16 }
}

const handleRawEditorMount = (editor: any, monaco: any) => {
  monacoEditorInstance.value = editor
  highlightVulnerableLines(monaco)
}

const highlightVulnerableLines = (monaco: any) => {
  if (!alert.value || !alert.value.vulnerable_lines) return;
  
  const lines = alert.value.vulnerable_lines.split(',')
  
  const decorations = lines.map(lineStr => {
    const lineNum = parseInt(lineStr)
    return {
      range: new monaco.Range(lineNum, 1, lineNum, 100),
      options: {
        isWholeLine: true,
        className: 'vulnerable-line-bg',
        glyphMarginClassName: 'vulnerable-line-glyph'
      }
    }
  })

  monacoEditorInstance.value?.createDecorationsCollection(decorations)
}

const triggerAction = async (action: string) => {
  if(!alert.value) return;
  await store.triggerTalon(alert.value.id, action)
}

onMounted(async () => {
    if (store.alerts.length === 0) {
        await store.fetchAlerts()
    }
    if(alert.value) {
        code.value = alert.value.manifest_yaml || '# No manifest captured.'
        // Simulate a fix for display
        fixedCode.value = code.value.replace(/privileged: true/g, 'privileged: false\n          runAsNonRoot: true')
    }
})

// Update code when store loads
watch(alert, (newVal) => {
  if (newVal) {
    code.value = newVal.manifest_yaml || '# No manifest captured.'
    fixedCode.value = code.value.replace(/privileged: true/g, 'privileged: false\n          runAsNonRoot: true')
  }
})
</script>

<style>
/* CSS added to global scope for Monaco decorations */
.vulnerable-line-bg {
  background-color: rgba(234, 67, 53, 0.2) !important;
}
.vulnerable-line-glyph {
  background-color: #EA4335 !important;
}
</style>
