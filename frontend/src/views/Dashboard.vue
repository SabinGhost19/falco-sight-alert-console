<template>
  <v-container fluid class="pa-6">
    <div class="text-h5 font-weight-medium mb-1 text-high-emphasis">Command Center</div>
    <div class="text-body-2 text-medium-emphasis mb-6">Overview of K8s workload security posture across your environments over the last 24h.</div>
    
    <!-- Summary Cards -->
    <v-row>
      <!-- Total Alerts -->
      <v-col cols="12" md="3">
        <v-card class="gc-border pb-2 d-flex flex-column justify-center px-4" height="120" color="surface">
          <div class="text-overline text-medium-emphasis font-weight-medium d-flex align-center">
             TOTAL PREVENTIONS <v-icon size="small" class="ml-auto" color="grey-lighten-1">mdi-bell-outline</v-icon>
          </div>
          <div class="text-h3 font-weight-regular mt-1">{{ store.alerts.length }}</div>
        </v-card>
      </v-col>
      <!-- Critical Threats -->
      <v-col cols="12" md="3">
        <v-card class="gc-border pb-2 d-flex flex-column justify-center px-4" height="120" color="surface">
          <div class="text-overline text-medium-emphasis font-weight-medium d-flex align-center">
             CRITICAL THREATS <v-icon size="small" class="ml-auto" color="error">mdi-alert-circle-outline</v-icon>
          </div>
          <div class="text-h3 font-weight-regular mt-1 text-error">{{ highPriority.length }}</div>
        </v-card>
      </v-col>
      <!-- Talon Status -->
      <v-col cols="12" md="3">
        <v-card class="gc-border pb-2 d-flex flex-column justify-center px-4" height="120" color="surface">
          <div class="text-overline text-medium-emphasis font-weight-medium d-flex align-center">
             TALON REMEDIATIONS <v-icon size="small" class="ml-auto" color="success">mdi-robot-outline</v-icon>
          </div>
          <div class="text-h3 font-weight-regular mt-1 text-success">{{ talonRemediated }}</div>
        </v-card>
      </v-col>
      <!-- Workloads Affected -->
      <v-col cols="12" md="3">
        <v-card class="gc-border pb-2 d-flex flex-column justify-center px-4" height="120" color="surface">
          <div class="text-overline text-medium-emphasis font-weight-medium d-flex align-center">
             WORKLOADS AFFECTED <v-icon size="small" class="ml-auto" color="primary">mdi-kubernetes</v-icon>
          </div>
          <div class="text-h3 font-weight-regular mt-1 text-primary">{{ uniqueWorkloads }}</div>
        </v-card>
      </v-col>
    </v-row>

    <!-- Charts & Analytics -->
    <v-row class="mt-4">
       <!-- Trend Line Chart -->
      <v-col cols="12" md="8">
        <v-card class="gc-border h-100" color="surface">
          <v-card-title class="text-subtitle-1 font-weight-medium text-high-emphasis border-b gc-border py-3">Alerts Trendline (24h)</v-card-title>
          <v-card-text class="pt-4 px-2 pb-0">
             <!-- Using ApexCharts -->
             <apexchart type="area" height="280" :options="chartData.options" :series="chartData.series"></apexchart>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Top Offenders Table -->
      <v-col cols="12" md="4">
        <v-card class="gc-border h-100" color="surface">
          <v-card-title class="text-subtitle-1 font-weight-medium text-high-emphasis border-b gc-border py-3">Top Vulnerable Workloads</v-card-title>
          <v-card-text class="pa-0">
             <v-list density="compact" lines="two">
               <v-list-item v-for="(item, i) in topOffenders" :key="i" class="border-b gc-border" @click="goToLogs(item.pod_name)">
                 <template v-slot:prepend>
                    <v-list-item-action start>
                        <v-icon color="primary" class="mr-3">mdi-kubernetes</v-icon>
                    </v-list-item-action>
                 </template>
                 <v-list-item-title class="font-weight-medium text-body-2">{{ item.pod_name }}</v-list-item-title>
                 <v-list-item-subtitle class="text-caption">Namespace: {{ item.namespace }}</v-list-item-subtitle>
                 
                 <template v-slot:append>
                    <v-chip size="small" color="error" variant="flat" class="font-weight-bold">{{ item.count }} alerts</v-chip>
                 </template>
               </v-list-item>

               <div v-if="topOffenders.length === 0" class="pa-6 text-center text-grey">All workloads are secure. No offenses detected.</div>
             </v-list>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAlertsStore } from '../store/alerts'
import VueApexCharts from 'vue3-apexcharts'

// Configurare globala vue component
const apexchart = VueApexCharts

const store = useAlertsStore()
const router = useRouter()

const highPriority = computed(() => store.alerts.filter(a => ['Critical', 'Emergency', 'error'].includes(a.priority) || a.priority === 'High'))
const talonRemediated = computed(() => store.alerts.filter(a => a.talon_status === 'Success' || a.talon_status === 'Requested').length)

// Numarul unic de Workload-uri atacate
const uniqueWorkloads = computed(() => {
  const workloads = new Set<string>()
  store.alerts.forEach(a => { if (a.pod_name) workloads.add(`${a.namespace}/${a.pod_name}`) })
  return workloads.size
})

// Topul vulnerabilitatilor
const topOffenders = computed(() => {
  const map: Record<string, {pod_name: string, namespace: string, count: number}> = {}
  store.alerts.forEach(a => {
      if(!a.pod_name) return
      const k = `${a.namespace}/${a.pod_name}`
      if(!map[k]) {
          map[k] = { pod_name: a.pod_name, namespace: a.namespace, count: 0 }
      }
      map[k].count++
  })
  
  return Object.values(map)
     .sort((a,b) => b.count - a.count)
     .slice(0, 5) // Extrage primele 5 locuri
})

const goToLogs = (podQuery: string) => {
   // Așa interconectăm Top Offenders cu Alerts Explorer
   router.push({ path: '/logs', query: { q: podQuery } })
}

// Date complet dinamice pentru grafic bazat pe creat_at
const chartData = computed(() => {
  // Vom construi un grafic pe ultimele 12 ore
  const hours = 12;
  const now = new Date();
  
  // Inițializăm bucket-uri pentru ultimele 12 ore
  const buckets = Array.from({length: hours}, (_, i) => {
    const d = new Date(now.getTime() - (hours - 1 - i) * 60 * 60 * 1000);
    return {
      hourLabel: `${d.getHours().toString().padStart(2, '0')}:00`,
      count: 0,
      timestamp: d.getTime()
    };
  });

  // Distribuim altele din store
  const twelveHoursAgo = now.getTime() - hours * 60 * 60 * 1000;
  
  store.alerts.forEach(alert => {
    if (!alert.created_at) return;
    const alertTime = new Date(alert.created_at).getTime();
    if (alertTime >= twelveHoursAgo) {
      // Găsim bucket-ul corect
      for (let i = 0; i < hours; i++) {
        const bucketStart = buckets[i].timestamp - (30 * 60 * 1000); // approx centering
        const bucketEnd = buckets[i].timestamp + (30 * 60 * 1000);
        if (alertTime >= bucketStart && alertTime < bucketEnd) {
          buckets[i].count++;
          break;
        }
      }
    }
  });

  return {
    series: [{
      name: 'Security Incidents',
      data: buckets.map(b => b.count)
    }],
    options: {
      chart: {
        fontFamily: 'Roboto, sans-serif',
        toolbar: { show: false },
        sparkline: { enabled: false },
        animations: { enabled: true }
      },
      colors: ['#D93025'], // Google Red
      stroke: { curve: 'smooth' as const, width: 2 },
      fill: {
        type: 'gradient',
        gradient: { shadeIntensity: 1, opacityFrom: 0.4, opacityTo: 0.05, stops: [0, 90, 100] }
      },
      dataLabels: { enabled: false },
      xaxis: {
        categories: buckets.map(b => b.hourLabel),
        axisBorder: { show: false },
        axisTicks: { show: false },
        labels: { style: { colors: '#9AA0A6' } }
      },
      yaxis: {
        labels: { style: { colors: '#9AA0A6' }, formatter: (val: number) => Math.round(val).toString() },
        min: 0,
        forceNiceScale: true
      },
      grid: {
        borderColor: '#F1F3F4',
        strokeDashArray: 4,
        xaxis: { lines: { show: true } },
        yaxis: { lines: { show: true } }
      }
    }
  };
});

onMounted(() => {
  store.fetchAlerts()
})
</script>
