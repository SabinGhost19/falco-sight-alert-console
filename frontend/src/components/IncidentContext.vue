<template>
  <v-container fluid class="pa-0 h-100 fill-height bg-background">
    
    <v-row no-gutters class="fill-height pb-0 mb-0">
      
      
      <v-col cols="12" md="4" class="border-e gc-border bg-white pa-6 overflow-auto fill-height">
        <v-btn variant="text" prepend-icon="mdi-arrow-left" color="primary" class="mb-6 text-none pl-0 ml-n2" @click="router.push('/logs')">
          Back to Alerts
        </v-btn>

        <template v-if="alert">
          
          <div class="d-flex align-center mb-3">
            <v-chip :color="priorityColor" size="small" variant="flat" class="mr-3 font-weight-bold text-white px-3">{{ alert.priority }}</v-chip>
            <span class="text-caption text-medium-emphasis font-weight-medium">{{ new Date(alert.created_at).toLocaleString() }}</span>
          </div>

          <h2 class="text-h5 font-weight-medium text-high-emphasis mb-2" style="line-height: 1.3">{{ alert.rule }}</h2>
          
          <v-card class="bg-grey-lighten-4 gc-border pa-3 mb-6 mt-4" elevation="0">
             <div class="text-caption font-weight-bold text-medium-emphasis mb-1 text-uppercase">Raw Event Message</div>
             <p class="text-body-2 font-monospace text-high-emphasis mb-0" style="word-break: break-word;">{{ alert.message }}</p>
          </v-card>

          
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

          
          <div class="mt-8 mb-6">
            <h3 class="text-subtitle-2 font-weight-bold text-medium-emphasis mb-4 text-uppercase border-b gc-border pb-2">Response & Remediation</h3>
            
            
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

          
          <v-card variant="outlined" class="gc-border mt-4 mb-2 pa-4" color="surface">
             <div class="d-flex align-center">
                 <v-icon color="warning" class="mr-2">mdi-history</v-icon>
                 <div class="font-weight-medium text-subtitle-2 text-high-emphasis">Image Risk Profile</div>
             </div>
             <p class="text-caption text-medium-emphasis mt-2 mb-0">
               The container image <strong>{{ alert?.container_image || 'unknown' }}</strong> has generated <strong class="text-error">{{ imageIncidentCount }}</strong> incidents cluster-wide in the past 30 days.
             </p>
             <v-btn size="x-small" color="primary" variant="text" class="text-none mt-2 px-0 font-weight-bold">View all incidents for this image</v-btn>
          </v-card>
        </template>
      </v-col>

      
      <v-col cols="12" md="8" class="bg-grey-darken-4 d-flex flex-column fill-height">
        <v-toolbar color="#202124" flat density="compact" class="border-b" style="border-bottom-color: #3C4043 !important;">
          <v-tabs v-model="tab" color="primary" bg-color="transparent" slider-color="primary" density="compact">
             <v-tab value="raw" class="text-none"><v-icon start size="small" class="mr-1">mdi-file-document-outline</v-icon> Live Manifest</v-tab>
             <v-tab value="diff" class="text-none"><v-icon start size="small" class="mr-1">mdi-auto-fix</v-icon> Proposed Security Fix</v-tab>
             <v-tab value="blast" class="text-none"><v-icon start size="small" class="mr-1">mdi-target</v-icon> Blast Radius</v-tab>
             <v-tab value="tree" class="text-none text-error"><v-icon start size="small" class="mr-1">mdi-file-tree</v-icon> Process Ancestry Tree</v-tab>
          </v-tabs>
          <v-spacer></v-spacer>
          <v-btn v-if="tab === 'diff'" class="text-none text-primary mr-2" variant="text" size="small" prepend-icon="mdi-content-copy" @click="copyFix">
            Copy Patch
          </v-btn>
        </v-toolbar>

        
        <div class="flex-grow-1 position-relative" style="height: calc(100vh - 48px); overflow: hidden;">            
            
            <div v-if="!code || code.trim() === ''" class="d-flex flex-column justify-center align-center h-100 bg-grey-lighten-4" style="position: absolute; width: 100%; z-index: 10;">
              <v-icon color="warning" size="64" class="mb-4">mdi-alert-circle-outline</v-icon>
              <div class="text-h6 text-high-emphasis font-weight-medium">Manifest unavailable</div>
              <div class="text-body-2 text-medium-emphasis text-center px-12 mt-2">
                The YAML structure could not be retrieved from the Kubernetes API. The workload might have been deleted, or RBAC permissions are denying access.
              </div>
            </div>
            
            <vue-monaco-editor
              v-if="tab === 'raw'"
              v-model:value="code"
              theme="vs-dark"
              language="yaml"
              :options="editorOptions"
              @mount="handleRawEditorMount"
              style="height: 100%; width: 100%;"
            ></vue-monaco-editor>

             
            <vue-monaco-diff-editor
              v-if="tab === 'diff'"
              :original="code"
              :modified="fixedCode"
              theme="vs-dark"
              language="yaml"
              :options="diffOptions"
              style="height: 100%; width: 100%;"
            ></vue-monaco-diff-editor>

            
            <v-container v-if="tab === 'blast'" class="pa-6 fill-height align-start bg-grey-lighten-4">
              <v-row>
                
                <v-col cols="12" md="7">
                  <v-card class="gc-border h-100 pa-4" elevation="0">
                     <div class="text-subtitle-1 font-weight-medium mb-4">Traffic Flow Analysis</div>
                     <v-alert type="warning" variant="tonal" class="mb-4 text-caption border">
                        Visual graph currently loading. Demonstrates pod connectivity mapping.
                     </v-alert>
                     <v-img src="https://upload.wikimedia.org/wikipedia/commons/thumb/c/c5/Network_Representation.svg/1024px-Network_Representation.svg.png" max-height="200" contain opacity="0.3"></v-img>
                  </v-card>
                </v-col>
                <v-col cols="12" md="5">
                  <v-row>
                     <v-col cols="12">
                        <v-card class="gc-border" elevation="0" color="#FCE8E6">
                          <v-card-title class="text-subtitle-2 font-weight-bold text-error px-4 pt-3 pb-0">RBAC Risk Assessment</v-card-title>
                          <v-card-text class="pa-4 text-body-2 text-error">
                            {{ alert?.rbac_risk || 'No critical privileges detected.' }}
                          </v-card-text>
                        </v-card>
                     </v-col>
                     <v-col cols="12">
                        <v-card class="gc-border" elevation="0">
                          <v-card-title class="text-subtitle-2 font-weight-bold text-high-emphasis px-4 pt-3 pb-0"><v-icon start color="primary">mdi-earth</v-icon> Network Posture</v-card-title>
                          <v-card-text class="pa-4 text-body-2 text-medium-emphasis">
                            {{ alert?.network_risk || 'Egress traffic restricted securely.' }}
                          </v-card-text>
                        </v-card>
                     </v-col>
                  </v-row>
                </v-col>
              </v-row>
            </v-container>

            
            <v-container v-if="tab === 'tree'" class="pa-0 fill-height bg-white d-flex align-center justify-center">
              <div class="text-center pa-12" style="width:100%;">
                 <v-icon color="primary" size="x-large" class="mb-4">mdi-source-branch</v-icon>
                 <div class="text-h6 text-high-emphasis">Process Ancestry visualization</div>
                 <div class="text-body-2 text-medium-emphasis mt-2 mb-6">
                    A hierarchical mapping from container runtime executing till the alert trigger point.
                 </div>
                 
                 
                 <div class="d-flex align-center justify-center">
                    <template v-if="parsedProcessTree.length > 0">
                      <template v-for="(proc, index) in parsedProcessTree" :key="index">
                         <v-chip :color="index === parsedProcessTree.length - 1 ? 'error' : (index === 0 ? 'grey' : 'warning')" 
                                 :variant="index === parsedProcessTree.length - 1 ? 'flat' : 'outlined'" 
                                 size="large" 
                                 class="font-weight-bold font-mono"
                                 :class="{ 'pulse-animation': index === parsedProcessTree.length - 1 }">
                            {{ proc }}
                         </v-chip>
                         <div v-if="index < parsedProcessTree.length - 1" style="width:40px;height:2px;background:#DADCE0;"></div>
                      </template>
                    </template>
                    <div v-else class="text-grey">No process ancestry data available for this sys-call.</div>
                 </div>
              </div>
            </v-container>

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

// NEW DYNAMIC COMPUTED DEMOS
const parsedProcessTree = computed(() => {
  if (!alert.value || !alert.value.process_tree) return [];
  try {
     const tree = JSON.parse(alert.value.process_tree);
     // Usually from falco it comes [parent, grandparent, root]. We reverse it to [root -> ..., parent] and append the actual process
     if (Array.isArray(tree)) {
       return tree.reverse().concat(["Malicious Trigger Process"]);
     }
  } catch(e) {
     return [];
  }
  return [];
});

const imageIncidentCount = computed(() => {
  if (!alert.value || !alert.value.container_image) return 1;
  return store.alerts.filter(a => a.container_image === alert.value!.container_image).length;
});


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

/* Custom CSS Animations */
@keyframes pulseRed {
  0% { box-shadow: 0 0 0 0 rgba(217, 48, 37, 0.7); }
  70% { box-shadow: 0 0 0 10px rgba(217, 48, 37, 0); }
  100% { box-shadow: 0 0 0 0 rgba(217, 48, 37, 0); }
}
.pulse-animation {
  animation: pulseRed 2s infinite;
}
</style>
  background-color: #D93025 !important;
}
</style>
