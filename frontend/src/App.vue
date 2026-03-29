<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { backend } from './bridge'
import type {
  HistoryItem,
  LogState,
  ProgressEvent,
  RepositorySuggestion,
  ReviewDoneEvent,
  StartReviewRequest
} from './types'
import ReviewForm from './components/ReviewForm.vue'
import ProgressPanel from './components/ProgressPanel.vue'
import LogViewer from './components/LogViewer.vue'
import HistoryList from './components/HistoryList.vue'
import ReportViewer from './components/ReportViewer.vue'

const form = ref<StartReviewRequest>({
  mrUrl: '',
  repoUrl: '',
  localRepoPath: '',
  sourceBranch: '',
  targetBranch: '',
  configPath: 'config.yaml',
  workspaceDir: './workspace'
})

const currentTaskId = ref('')
const progress = ref<ProgressEvent | null>(null)
const doneEvent = ref<ReviewDoneEvent | null>(null)
const history = ref<HistoryItem[]>([])
const running = ref(false)
const pullingCode = ref(false)
const clearingCache = ref(false)
const errorMessage = ref('')
const infoMessage = ref('')
const availableBranches = ref<string[]>([])
const showLogs = ref(false)
const logState = ref<LogState>({ logPath: '', content: '', updatedAt: '' })
const logViewerRef = ref<any>(null)
const repoSuggestion = ref<RepositorySuggestion>({
  repoUrl: '',
  candidates: [],
  resolvedByMr: false,
  manualRepoUrl: '',
  message: ''
})

let logTimer: number | null = null
let repoSuggestTimer: number | null = null
const lastSuggestedRepoUrl = ref('')
const manualRepoOverride = ref(false)

const runtimeRequest = computed(() => ({
  configPath: form.value.configPath,
  workspaceDir: form.value.workspaceDir
}))

const report = computed(() => doneEvent.value?.report ?? null)
const disabled = computed(() => running.value || pullingCode.value || clearingCache.value)
const statusText = computed(() => {
  if (running.value) return '监视进行中'
  if (doneEvent.value) return '报告已生成'
  return '等待任务启动'
})

function normalizeError(err: unknown): string {
  if (err instanceof Error) return err.message
  if (typeof err === 'string') return err
  return '操作失败，请查看日志了解详情。'
}

function clearMessages() {
  errorMessage.value = ''
  infoMessage.value = ''
}

function setInfo(message: string) {
  infoMessage.value = message
  errorMessage.value = ''
}

function setError(message: string) {
  errorMessage.value = message
}

async function loadHistory() {
  history.value = await backend.listHistory(runtimeRequest.value)
}

async function loadLogState(forceScroll = false) {
  logState.value = await backend.getLogState(runtimeRequest.value)
  if (showLogs.value || forceScroll) {
    await nextTick()
    const logEl = logViewerRef.value?.logViewerRef
    if (logEl) {
      logEl.scrollTop = logEl.scrollHeight
    }
  }
}

function stopLogPolling() {
  if (logTimer !== null) {
    window.clearInterval(logTimer)
    logTimer = null
  }
}

function startLogPolling() {
  stopLogPolling()
  logTimer = window.setInterval(() => void loadLogState(true), 1000)
}

function handleRepoUrlEdited(value: string) {
  const normalized = value.trim()
  manualRepoOverride.value = normalized !== '' && normalized !== lastSuggestedRepoUrl.value
  if (normalized === '' || normalized === lastSuggestedRepoUrl.value) {
    manualRepoOverride.value = false
  }
}

async function suggestRepository() {
  if (!form.value.mrUrl.trim() && !form.value.repoUrl.trim()) {
    repoSuggestion.value = {
      repoUrl: '',
      candidates: [],
      resolvedByMr: false,
      manualRepoUrl: '',
      message: ''
    }
    lastSuggestedRepoUrl.value = ''
    return
  }

  try {
    const suggestion = await backend.suggestRepository(form.value)
    repoSuggestion.value = suggestion
    if (
      suggestion.repoUrl &&
      (!manualRepoOverride.value || !form.value.repoUrl.trim() || form.value.repoUrl.trim() === lastSuggestedRepoUrl.value)
    ) {
      form.value.repoUrl = suggestion.repoUrl
    }
    lastSuggestedRepoUrl.value = suggestion.repoUrl
  } catch {
    // 自动提示不阻断主流程
  }
}

async function pullCode() {
  clearMessages()
  pullingCode.value = true
  try {
    const result = await backend.pullCode(form.value)
    if (result.repoUrl) {
      form.value.repoUrl = result.repoUrl
      lastSuggestedRepoUrl.value = result.repoUrl
    }
    form.value.sourceBranch = result.sourceBranch || form.value.sourceBranch
    form.value.targetBranch = result.targetBranch || form.value.targetBranch
    availableBranches.value = result.availableBranches ?? []
    setInfo(result.message || '代码拉取完成。')
    await Promise.all([loadHistory(), loadLogState(true)])
  } catch (err) {
    setError(normalizeError(err))
  } finally {
    pullingCode.value = false
  }
}

async function startReview() {
  clearMessages()
  doneEvent.value = null
  progress.value = null
  running.value = true
  try {
    currentTaskId.value = await backend.startReview(form.value)
    setInfo('监视已启动，正在执行前置校验与审计流程。')
    if (showLogs.value) await loadLogState(true)
  } catch (err) {
    running.value = false
    setError(normalizeError(err))
  }
}

async function cancelReview() {
  if (!currentTaskId.value) return
  clearMessages()
  try {
    await backend.cancelReview(currentTaskId.value)
    running.value = false
    currentTaskId.value = ''
    setInfo('任务已取消。')
    await loadLogState(true)
  } catch (err) {
    setError(normalizeError(err))
  }
}

async function clearCache() {
  clearMessages()
  clearingCache.value = true
  try {
    const result = await backend.clearCache(runtimeRequest.value)
    setInfo(result.message || '缓存已清理。')
    history.value = []
    doneEvent.value = null
    await loadLogState(true)
  } catch (err) {
    setError(normalizeError(err))
  } finally {
    clearingCache.value = false
  }
}

function toggleLogs() {
  showLogs.value = !showLogs.value
  if (showLogs.value) {
    void loadLogState(true)
    startLogPolling()
  } else {
    stopLogPolling()
  }
}

async function openReport(payload?: { htmlPath: string; reportDir: string }) {
  clearMessages()
  const target = payload ?? {
    htmlPath: doneEvent.value?.htmlPath || '',
    reportDir: doneEvent.value?.reportDir || ''
  }
  try {
    const result = await backend.openReport(target)
    setInfo(result.message || '已打开报告。')
  } catch (err) {
    setError(normalizeError(err))
  }
}

function openHistoryReport(item: HistoryItem) {
  void openReport({ htmlPath: item.htmlPath, reportDir: item.reportDir })
}

function openHistoryReportDirectory(item: HistoryItem) {
  void openReportDirectory({ reportDir: item.reportDir })
}

async function openReportDirectory(payload?: { reportDir: string }) {
  clearMessages()
  const target = payload ?? { reportDir: doneEvent.value?.reportDir || '' }
  try {
    const result = await backend.openReportDirectory(target)
    setInfo(result.message || '已打开报告目录。')
  } catch (err) {
    setError(normalizeError(err))
  }
}

watch(
  () => [form.value.mrUrl, form.value.configPath],
  () => {
    if (repoSuggestTimer !== null) {
      window.clearTimeout(repoSuggestTimer)
    }
    repoSuggestTimer = window.setTimeout(() => void suggestRepository(), 280)
  }
)

onMounted(() => {
  void loadHistory()
  void loadLogState()
  void suggestRepository()

  window.runtime?.EventsOn?.('review:progress', (event: ProgressEvent) => {
    progress.value = event
  })
  window.runtime?.EventsOn?.('review:error', (event: { taskId: string; message: string }) => {
    running.value = false
    setError(event.message)
    stopLogPolling()
  })
  window.runtime?.EventsOn?.('review:done', (event: ReviewDoneEvent) => {
    running.value = false
    doneEvent.value = event
    stopLogPolling()
    setInfo('监视完成，报告已生成。你可以直接点击“查看报告”打开 HTML 报告。')
    void loadHistory()
  })
})

onUnmounted(() => {
  stopLogPolling()
  if (repoSuggestTimer !== null) {
    window.clearTimeout(repoSuggestTimer)
    repoSuggestTimer = null
  }
})
</script>

<template>
  <div class="app-shell">
    <header class="hero hero-panel">
      <div>
        <h1>AI代码监视</h1>
        <p>
          面向生产交付的桌面审计界面：支持 MR/PR 自动识别仓库地址、历史报告回溯、HTML 报告快捷打开，
          并对结果卡片做了默认折叠与可展开详情优化。
        </p>
      </div>
      <div class="hero-status">
        <div class="status-badge">{{ statusText }}</div>
        <div class="small">工作区：{{ form.workspaceDir || './workspace' }}</div>
      </div>
    </header>

    <div v-if="errorMessage" class="notice error">{{ errorMessage }}</div>
    <div v-if="infoMessage" class="notice success">{{ infoMessage }}</div>

    <main class="layout">
      <section class="left-column">
        <div class="card glow-card">
          <ReviewForm
            :form="form"
            :available-branches="availableBranches"
            :running="running"
            :pulling-code="pullingCode"
            :clearing-cache="clearingCache"
            :show-logs="showLogs"
            :disabled="disabled"
            :repo-suggestion-message="repoSuggestion.message"
            :repo-suggestion-candidates="repoSuggestion.candidates"
            :repo-suggestion-resolved-by-mr="repoSuggestion.resolvedByMr"
            @pull-code="pullCode"
            @start-review="startReview"
            @cancel-review="cancelReview"
            @toggle-logs="toggleLogs"
            @clear-cache="clearCache"
            @load-history="loadHistory"
            @repo-url-edited="handleRepoUrlEdited"
          />
        </div>

        <div class="card glow-card">
          <HistoryList
            :history="history"
            @open-report="openHistoryReport"
            @open-dir="openHistoryReportDirectory"
          />
        </div>
      </section>

      <section class="right-column">
        <div class="card glow-card" v-if="progress || running">
          <div class="panel-header compact">
            <div>
              <h2>执行状态</h2>
              <p class="small">监视完成后会自动生成 HTML / Markdown / JSON 报告。</p>
            </div>
            <div class="report-actions-inline" v-if="doneEvent">
              <button class="small-action primary" @click="openReport()">查看报告</button>
              <button class="small-action secondary" @click="openReportDirectory()">打开目录</button>
            </div>
          </div>
          <ProgressPanel :progress="progress" />
        </div>

        <div class="card glow-card" v-if="showLogs">
          <div class="panel-header compact">
            <div>
              <h2>运行日志</h2>
              <p class="small">{{ logState.updatedAt ? `最后更新：${logState.updatedAt}` : '等待日志输出' }}</p>
            </div>
          </div>
          <LogViewer ref="logViewerRef" :show="showLogs" :log-path="logState.logPath" :content="logState.content" />
        </div>

        <div class="card glow-card" v-if="report && doneEvent">
          <ReportViewer :report="report" />
        </div>
      </section>
    </main>
  </div>
</template>

<style scoped>
.hero-panel {
  background: linear-gradient(135deg, rgba(14, 23, 42, 0.92), rgba(30, 41, 59, 0.82));
  border: 1px solid rgba(99, 102, 241, 0.18);
  border-radius: 1.5rem;
  padding: 1.5rem 1.6rem;
  box-shadow: 0 20px 42px rgba(0, 0, 0, 0.22);
}

.hero-panel h1 {
  margin: 0 0 0.5rem;
  font-size: 2.1rem;
  color: #eef2ff;
}

.hero-status {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.6rem;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 0.45rem 0.85rem;
  border-radius: 999px;
  background: linear-gradient(135deg, rgba(37, 99, 235, 0.22), rgba(124, 58, 237, 0.22));
  color: #dbeafe;
  border: 1px solid rgba(96, 165, 250, 0.24);
  font-weight: 700;
}

.glow-card {
  position: relative;
  overflow: hidden;
}

.glow-card::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background: radial-gradient(circle at top right, rgba(96, 165, 250, 0.08), transparent 28%),
    radial-gradient(circle at bottom left, rgba(34, 197, 94, 0.06), transparent 26%);
}

.report-actions-inline,
.report-actions-inline + * {
  position: relative;
  z-index: 1;
}

.small-action {
  border: none;
  border-radius: 0.9rem;
  padding: 0.6rem 0.85rem;
  color: white;
  cursor: pointer;
}

.small-action.primary {
  background: linear-gradient(135deg, #2563eb, #7c3aed);
}

.small-action.secondary {
  background: rgba(51, 65, 85, 0.92);
}

@media (max-width: 1280px) {
  .hero-status {
    align-items: flex-start;
  }
}
</style>
