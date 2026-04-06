# FalcoSight - Falco Alert Console

FalcoSight este o consolă avansată de management și răspuns pentru alertele de securitate generate de [Falco](https://falco.org/) în medii Kubernetes. Soluția colectează, corelează și îmbogățește alertele standard Falco cu informații din manifestele Kubernetes on-the-fly, calculează Blast Radius-ul și se integrează fluent cu **Falco Talon** pentru remediere automată și manuală.

## Arhitectură și Componente

Proiectul este compus din două părți principale:
- **Backend (Go/Fiber):** API ultra-rapid de ingestie webhook-uri din Falco, ce rulează o analiză statică a Pod-urilor expuse (K8s API Correlator, Mitre mapping, Process tree generation) și stocare în PostgreSQL. Forwardează deciziile către Talon în mod asincron.
- **Frontend (Vue 3 / Vite):** Interfață modernă ce vizualizează Dashboard-ul de Alerte, metricile de securitate (Trendlines) și pune la dispoziție Explorer-ul Log-urilor pe baza framework-ului de design Google.

## Tehnologii Core
- Miez API: **Golang** cu framework-ul [Fiber](https://gofiber.io/)
- UI: **Vue.js 3** + **Vuetify** + **ApexCharts** + **Monaco Editor**
- DB: **PostgreSQL** + **GORM**
- Integrare: Kubernetes Client-Go, Falco Webhooks, Falco Talon (pentru remediere RBAC, pod deletion, etc.)

## Pornirea și Instalarea Componentelor (Docker / Helm)

### Rulare locală folosind Docker Compose
Dacă se dorește rularea locală pentru testare în afara clusterului (Kubeconfig implicit se va încărca, sau în funcție de environment variables):
```bash
docker-compose up --build
```
*Notă: Baza de date este instanțiată local printr-un container PostgreSQL separat.*

### Deploy pe Kubernetes (Helm)
Chart-ul principal de Helm se află în `helm/falcosight/`.
```bash
helm upgrade --install falcosight ./helm/falcosight -n falco --create-namespace
```

## Cum rulează fluxul decizional?
1. **Detectare:** Falco Kernel Module detectează activitate malițioasă la nivel de syscall.
2. **Ingestie FalcoSight:** Payload-ul webhook este preluat de `/api/falco/webhook`.
3. **Corelare K8s:** Se interoghează automat K8s API pentru a prelua manifestul `.yaml` al workload-ului atacat direct. Dacă alertă vine de la Host Level, se adaptează.
4. **Trimitere Alertă UI & Talon:** Datele procesate apar în Dashboard în timp real, iar Falco Talon poate lua măsuri corective pe baza regulilor hardcodate (#Talon).
5. **Human-in-the-Loop:** Un operator securitate poate iniția o comandă de izolare direct din Dashboard apasand butoanele specifice din `Incident Context`.

## Testare Automată
Folderul `test-pod-alert-triggers/` furnizează `.yaml`-uri folositoare (blast radius expus, repeat offenders) care creează intenționat comportamente non-compliant pentru a arunca excepții monitorizate de Falco și redate în Consola FalcoSight.

## Documentație suplimentară
Directorul `doc/` conține:
- **overview_1.md**: Abordarea inițială și design-ul arhitectural.
- **overview_2.md**: Prezentarea integrării grafice, a log-urilor Falco și particularitățile stivei.
- **overview_3.md**: Aspecte specifice despre Blast Radius, corelarea pe arbore de procese și etichete MITRE.
- **k8s_correlator_and_talon.md**: Adăugat recent pentru a detalia mecanismul in-memory de analiză on-the-fly a manifestelor de Pod și integrarea cu webhook-ul Talon.
