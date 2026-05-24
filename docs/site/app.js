const params = new URLSearchParams(window.location.search);
const apiBase = params.get('api')?.replace(/\/+$/, '') || window.__AIHYPERVISOR_API__ || '';

const apiStatus = document.getElementById('api-status');
const apiDetail = document.getElementById('api-detail');
const readyStatus = document.getElementById('ready-status');
const versionStatus = document.getElementById('version-status');
const statusSource = document.getElementById('status-source');
const apiBaseLabel = 'static demo';
const metricsGrid = document.getElementById('metrics-grid');
const workflowGrid = document.getElementById('workflow-grid');
const activityList = document.getElementById('activity-list');
const metricsSource = document.getElementById('metrics-source');
const activitySource = document.getElementById('activity-source');
const demoRefresh = document.getElementById('demo-refresh');

const demoMetrics = [
  {
    label: 'GPU allocation',
    value: '92%',
    note: '12 GPUs reserved across 3 nodes',
    fill: 92,
    tone: 'good',
  },
  {
    label: 'Active VMs',
    value: '38',
    note: '26 running, 8 paused, 4 queued',
    fill: 76,
    tone: 'live',
  },
  {
    label: 'Queue depth',
    value: '4',
    note: 'Provisioning jobs awaiting execution',
    fill: 34,
    tone: 'warn',
  },
  {
    label: 'Policy pass rate',
    value: '99.7%',
    note: 'Last 200 requests passed checks',
    fill: 99,
    tone: 'good',
  },
  {
    label: 'Median startup',
    value: '1m 24s',
    note: 'From request accepted to VM ready',
    fill: 61,
    tone: 'live',
  },
  {
    label: 'Health score',
    value: 'A-',
    note: 'All control-plane dependencies online',
    fill: 88,
    tone: 'good',
  },
];

const demoWorkflows = [
  {
    title: 'Release readiness',
    detail: 'Create a checklist, open a release PR, and notify #releases.',
    tags: ['approval gate', 'devops'],
    status: 'running',
  },
  {
    title: 'Customer escalation',
    detail: 'Triage a case, open a ticket, and email the account owner.',
    tags: ['policy check', 'urgent'],
    status: 'queued',
  },
  {
    title: 'Capacity expansion',
    detail: 'Score hosts, reserve GPUs, and schedule additional VMs.',
    tags: ['scheduler', 'gpu'],
    status: 'ready',
  },
  {
    title: 'Audit export',
    detail: 'Package events, sign the bundle, and publish the archive.',
    tags: ['compliance', 'export'],
    status: 'done',
  },
];

const demoActivity = [
  {
    title: 'GPU pool expanded',
    detail: 'Two nodes were added to the scheduling pool after health checks passed.',
    time: '2m ago',
    tone: 'success',
  },
  {
    title: 'Approval requested',
    detail: 'A high-risk workflow is waiting on a human approval gate.',
    time: '7m ago',
    tone: 'warning',
  },
  {
    title: 'VM lifecycle completed',
    detail: 'Provisioning, network attach, and health validation all finished.',
    time: '13m ago',
    tone: 'info',
  },
  {
    title: 'Telemetry exported',
    detail: 'Prometheus and trace samples were shipped to the observability stack.',
    time: '19m ago',
    tone: 'success',
  },
];

function setStatus(state, detail, sourceLabel) {
  apiStatus.textContent = state;
  apiDetail.textContent = detail;
  statusSource.textContent = sourceLabel;
}

function renderMetricCards() {
  if (!metricsGrid) {
    return;
  }

  metricsGrid.innerHTML = demoMetrics
    .map(
      (metric) => `
        <article class="metric-card">
          <div class="metric-top">
            <p class="metric-label">${metric.label}</p>
            <span class="pill ${metric.tone}">${metric.tone}</span>
          </div>
          <p class="metric-value">${metric.value}</p>
          <p class="metric-note">${metric.note}</p>
          <div class="meter" aria-hidden="true">
            <div class="meter-fill" style="width: ${metric.fill}%"></div>
          </div>
        </article>
      `
    )
    .join('');
}

function renderWorkflowCards() {
  if (!workflowGrid) {
    return;
  }

  workflowGrid.innerHTML = demoWorkflows
    .map(
      (workflow) => `
        <article class="workflow-card">
          <div class="workflow-top">
            <p class="workflow-title">${workflow.title}</p>
            <span class="pill ${workflow.status === 'done' ? 'good' : workflow.status === 'running' ? 'live' : 'warn'}">${workflow.status}</span>
          </div>
          <p class="workflow-detail">${workflow.detail}</p>
          <div class="workflow-meta">
            ${workflow.tags.map((tag) => `<span class="pill">${tag}</span>`).join('')}
          </div>
        </article>
      `
    )
    .join('');
}

function renderActivityFeed() {
  if (!activityList) {
    return;
  }

  activityList.innerHTML = demoActivity
    .map(
      (entry) => `
        <article class="activity-item">
          <div class="activity-top">
            <p class="activity-title">${entry.title}</p>
            <span class="pill ${entry.tone === 'success' ? 'good' : entry.tone === 'warning' ? 'warn' : 'live'}">${entry.tone}</span>
          </div>
          <p class="activity-detail">${entry.detail}</p>
          <span class="activity-time">${entry.time}</span>
        </article>
      `
    )
    .join('');
}

function renderDemoConsole() {
  renderMetricCards();
  renderWorkflowCards();
  renderActivityFeed();

  if (metricsSource) {
    metricsSource.textContent = apiBase ? 'live shell + demo data' : 'static demo';
  }

  if (activitySource) {
    activitySource.textContent = apiBase ? 'live shell + demo data' : 'demo feed';
  }
}

async function fetchHealth(path) {
  const response = await fetch(`${apiBase}${path}`, { cache: 'no-store' });
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`);
  }

  const body = await response.json().catch(() => ({}));
  return { status: response.status, body };
}

async function hydrateStatus() {
  renderDemoConsole();

  if (!apiBase) {
    setStatus(
      'demo',
      'Static demo mode is active. Add ?api=https://your-api.example.com to query a live deployment.',
      apiBaseLabel
    );
    readyStatus.textContent = 'demo-ready';
    versionStatus.textContent = 'github pages';
    return;
  }

  statusSource.textContent = apiBase;
  try {
    const health = await fetchHealth('/health');
    const ready = await fetchHealth('/ready').catch(() => null);

    const healthState = health.body?.status || health.body?.state || 'healthy';
    setStatus('online', `Health endpoint returned ${healthState}.`, apiBase);
    readyStatus.textContent = ready?.body?.status || ready?.body?.state || 'ready';
    versionStatus.textContent = health.body?.version || 'live';
  } catch (error) {
    setStatus('degraded', `Unable to reach ${apiBase}. ${error.message}`, apiBase);
    readyStatus.textContent = 'unknown';
    versionStatus.textContent = 'offline';
  }
}

if (demoRefresh) {
  demoRefresh.addEventListener('click', () => {
    renderDemoConsole();
  });
}

hydrateStatus();