<template>
  <v-container fluid class="pa-0 h-100 fill-height bg-background">
    <-> Split Screen Layout -->
    <v-row no-gutters class="fill-height pb-0 mb-0">
      
      <-> Left Panel: Context & Falco Alert Info -->
      <v-col cols="12" md="4" class="border-e gc-border bg-white pa-6 overflow-auto fill-height">
        <v-btn variant="text" prepend-icon="mdi-arrow-left" color="primary" class="mb-6 text-none pl-0 ml-n2" @click="router.push('/logs')">
          Back to Alerts
        </v-btn>

        <template v-if="alert">
          <-> Alert Header -->
          <div class="d-flex align-center mb-3">
            <v-chip :color="priorityColor" size="small" variant="flat" class="mr-3 font-weight-bold text-white px-3">{{ alert.priority }}</v-chip>
            <span class="text-caption text-medium-emphasis font-weight-medium">{{ new Date(alert.created_at).toLocaleString() }}</span>
          </div>

          <h2 class="text-h5 font-weight-medium text-high-emphasis mb-2" style="line-height: 1.3">{{ alert.rule }}</h2>
          
          <v-card class="bg-grey-lighten-4 gc-border pa-3 mb-6 mt-4" elevation="0">
             <div class="text-caption font-weight-bold text-medium-emphasis mb-1 text-uppercase">Raw Event Message</div>
             <p class="text-body-2 font-monospace text-high-emphasis mb-0" style="word-break: break-word;">{{ alert.message }}</p>
          </v-card>

          <-> Corelate K8s Workload Details (GCP style key-value) -->
          <div class="mb-6">
            <h3 class="text-subtitle-2 font-weight-bold text-medium-emphasis mb-3 text-uppercase border-b gc-border pb-2">Target Workload</h3>
            
            <v-row no-gutters class="mb-2">
               <v-col cols="4"><span class="text-body-2 text-medium-emphasis">Cluster</span></v-col>
               <v-col cols="8"><span class="text-body-2 font-weight-medium text-high-emphasis">production-cluster-01</span></v-col>
            </v-row>
            <v-row no-gutters class="mb-2">
               <v-col cols="4"><span class="text-body-2 text-medium-emphasis">Namespace</span></v-col>
               <v-col cols="8"><span class="text-body-2 font-weight-medium text-primary cursor-pointer">{{ alert.namespace }}</span></v-col>
            </v-row>
            <v-row no-gutters class="mb-2">
               <v-col cols="4"><span class="text-body-2 text-medium-emphasis">Pod</span></v-col>
               <v-col cols="8"><span class="text-body-2 font-weight-medium text-primary cursor-pointer">{{ alert.pod_name }}</span></v-col>
            </v-row>
            <v-row no-gutters v-if="alert.container_name">
               <v-col cols="4"><span class="text-body-2 text-medium-emphasis">Container</span></v-col>
               <v-col cols="8"><span class="text-body-2 font-weight-medium text-high-emphasis">{{ alert.container_name }}</span></v-col>
            </v-row>
          </div>

          <-> Falco Talon Action Center (SOAR) -->
          <div class="mt-8">
            <h3 class="text-subtitle-2 font-weight-bold text-medium-emphasis mb-4 text-uppercase border-b gc-border pb-2">Response & Remediation</h3>
            
            <-> Timeline / Stepper visualization -->
            <div class="d-flex align-start mb-4">
              <div class="d-flex flex-column align-center mr-4">
                 <v-icon color="error" size="small">mdi-shield-alert</v-icon>
                 <div style="width: 2px; height: 30px; background-color: #DADCE0;" class="my-1"></div>
                 <v-icon :color="alert.talon_status ? 'success' : 'grey-lighten-1'" size="small">
                    {{ alert.talon_status ? 'mdi-check-circle' : 'mdi-circle-outline' }}
                 </v-icon>
              </div>
              <div>
                 <div class="mb-4">
                    <div class="text-body-2 font-weight-bold text-high-emphasis">Threat Detected</div>
                    <div class="text-caption text-grey">Falco engine intercepted system call.</div>
                 </div>
                 
                 <div>
                    <div class="text-body-2 font-weight-bold" :class="alert.talon_status ? 'text-success' : 'text-high-emphasis'">
                       SOAR Automation (Talon)
                    </div>
                    <div class="text-caption text-grey" v-if="alert.talon_status">
                       Action applied: <strong>{{ alert.talon_action }}</strong> ({{ alert.talon_status }})
                    </div>
                    <div class="text-caption text-grey mb-3" v-else>
                       No automatic policies matched or running in Audit Mode.
                    </div>

                    <-> Manual Overrides if unchecked -->
                    <div v-if="!alert.talon_status" class="mt-2 d-flex flex-column gap-2">
                       <v-btn color="primary" class="text-none font-weight-medium" variant="flat" size="small" @click="triggerAction('network_isolate')" :disabled="!alert" prepend-icon="mdi-lan-disconnect">
                          Manual Override: Isolate Pod
                       </v-btn>
                       <v-btn color="error" class="text-none mt-2 font-weight-medium" variant="outlined" size="small" @click="triggerAction('terminate')" :disabled="!alert" prepend-icon="mdi-close-octagon">
                          Terminate Workload
                       </v-btn>
                    </div>
                 </div>
              </div>
            </div>
          </div>
        </template>
      </v-col>

      <-> Right Panel: Code Viewer (Monaco Editor) & Fixes -->
      <v-col cols="12" md="8" class="bg-grey-darken-4 d-flex flex-column fill-height">
        <v-toolbar color="#202124" flat density="compact" class="border-b" style="border-bottom-color: #3C4043 !important;">
          <v-tabs v-model="tab" color="primary" bg-color="transparent" slider-color="primary" density="compact">
             <v-tab value="raw" class="text-none"><v-icon start size="small" class="mr-1">mdi-file-document-outline</v-icon> Live Manifest</v-tab>
             <v-tab value="diff" class="text-none ml-2"><v-icon start size="small" class="mr-1">mdi-auto-fix</v-icon> Proposed Security Fix</v-tab>
          </v-tabs>
          <v-spacer></v-spacer>
          <v-btn v-if="tab === 'diff'" class="text-none text-primary mr-2" variant="text" size="small" prepend-icon="mdi-content-copy" @click="copyFix">
            Copy Patch
          </v-btn>
        </v-toolbar>

        <-> Container for Monaco Editor -->
        <div class="flex-grow-1 position-relative" style="height: calc(100vh - 48px); overflow: hidden;">            
            <!-- Empty State (Graceful Degradation) -->
            <div v-if="!code || code.trim() === ''" class="d-flex flex-column justify-center align-center h-100 bg-grey-lighten-4" style="position: absolute; width: 100%; z-index: 10;">
              <v-icon color="warning" size="64" class="mb-4">mdi-alert-circle-outline</v-icon>
              <div class="text-h6 text-high-emphasis font-weight-medium">Manifest unavailable</div>
              <div class="text-body-2 text-medium-emphasis text-center px-12 mt-2">
                The YAML structure could not be retrieved from the Kubernetes API. The workload might have been deleted, or RBAC permissions are denying access.
              </div>
            </div>
            <-> RAW YAML Viewer with Custom Decorations -->
            <vue-monaco-editor
              v-if="tab === 'raw'"
              v-model:value="code"
              theme="vs-dark"
              language="yaml"
              :options="editorOptions"
              @mount="handleRawEditorMount"
              style="height: 100%; width: 100%;"
            ></vue-monaco-editor>

             <-> DIFF Editor for seeing Fixes -->
            <vue-monaco-diff-editor
              v-if="tab === 'diff'"
              :original="code"
              :modified="fixedCode"
              theme="vs-dark"
              language="yaml"
              :options="diffOptions"
              style="height: 100%; width: 100%;"
            ></vue-monaco-diff-editor>
        </div>
      </v-col>

    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAlertsStore } from '../store/alerts'
import { VueMonacoEditor, VueMonacoDiffEditor } from '@guolao/vue-monaco-editor'
import { useClipboard } from '@vueuse/core'

const route = useRoute()
const router = useRouter()
const store = useAlertsStore()
const { copy } = useClipboard()

const id = parseInt(route.params.id as string)

const alert = computed(() => store.alerts.find(a => a.id === id))
const priorityColor = computed(() => {
  const p = alert.value?.priority?.toLowerCase()
  if (p === 'emergency') return '#B31412'
  if (p === 'critical' || p === 'error') return '#D93025'
  if (p === 'warning') return '#F29900'
  return '#1A73E8'
})

const tab = ref('raw')
const code = ref('')
const fixedCode = ref('')
const monacoEditorInstance = ref<any>(null)

// Monaco Editor Config matching Google Cloud Shell
const editorOptions = {
  automaticLayout: true,
  minimap: { enabled: false },
  readOnly: true,
  scrollBeyondLastLine: false,
  wordWrap: 'on',
  fontSize: 14,
  fontFamily: "'JetBrains Mono', 'Fira Code', 'Roboto Mono', monospace",
  padding: { top: 16 },
  renderLineHighlight: 'all',
}

const diffOptions = {
  ...editorOptions,
  enableSplitViewResizing: true,
  renderSideBySide: true,
  readOnly: true,
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

const copyFix = () => {
   copy(fixedCode.value)
   // Provide transient notification via Vuetify snackbar if available natively
}

const triggerAction = async (action: string) => {
  if(!alert.value) return;
  await store.triggerTalon(alert.value.id, action)
}

onMounted(async () => {
    if (store.alerts.length === 0) {
        await store.fetchAlerts()
    }
    setupCode()
})

watch(alert, () => setupCode())

function setupCode() {
  if (alert.value) {
    code.value = alert.value.manifest_yaml || '# No live manifest captured for this workload.'
    // Simulate a basic remediation config replace 
    fixedCode.value = code.value.replace(/privileged:\s*true/g, 'privileged: false\n          allowPrivilegeEscalation: false\n          runAsNonRoot: true')
  }
}
</script>

<style>
/* Monaco Decorations injected into DOM securely */
.vulnerable-line-bg {
  background-color: rgba(217, 48, 37, 0.2) !important; /* Google Red transparent */
}
.vulnerable-line-glyph {
  background-color: #D93025 !important;
}
</style>
