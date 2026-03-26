<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { backend } from './bridge'
import type { HistoryItem, LogState, ProgressEvent, ReviewDoneEvent, StartReviewRequest } from './types'

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
const logState = ref<LogState>({
  logPath: '',
  content: '',
  updatedAt: ''
})
const logViewerRef = ref<HTMLElement | null>(null)

let logTimer: number | null = null

const runtimeRequest = computed(() => ({
  configPath: form.value.configPath,
  workspaceDir: form.value.workspaceDir
}))

const findings = computed(() => doneEvent.value?.report.findings ?? [])
const report = computed(() => doneEvent.value?.report)

function normalizeError(err: unknown): string {
  if (err instanceof Error) {
    return err.message
  }
  if (typeof err === 'string') {
    return err
  }
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
    if (logViewerRef.value) {
      logViewerRef.value.scrollTop = logViewerRef.value.scrollHeight
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
  logTimer = window.setInterval(() => {
    void loadLogState(true)
  }, 1000)
}

async function pullCode() {
  clearMessages()
  pullingCode.value = true
  try {
    const result = await backend.pullCode(form.value)
    if (result.repoUrl) {
      form.value.repoUrl = result.repoUrl
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
    if (showLogs.value) {
      await loadLogState(true)
    }
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
    currentTaskId.value = ''
    progress.value = null
    doneEvent.value = null
    history.value = []
    availableBranches.value = []
    form.value.sourceBranch = ''
    form.value.targetBranch = ''
    setInfo(result.message || '缓存清理完成。')
    await Promise.all([loadHistory(), loadLogState(true)])
  } catch (err) {
    setError(normalizeError(err))
  } finally {
    clearingCache.value = false
  }
}

function toggleLogs() {
  showLogs.value = !showLogs.value
}

watch(
  () => [form.value.configPath, form.value.workspaceDir],
  async () => {
    try {
      await Promise.all([loadHistory(), loadLogState()])
    } catch (err) {
      setError(normalizeError(err))
    }
  },
  { immediate: true }
)

watch(showLogs, async (value) => {
  if (value) {
    await loadLogState(true)
    startLogPolling()
  } else {
    stopLogPolling()
  }
})

onMounted(() => {
  backend.on('review:progress', (payload: ProgressEvent) => {
    if (!currentTaskId.value || payload.taskId === currentTaskId.value) {
      progress.value = payload
    }
  })

  backend.on('review:done', async (payload: ReviewDoneEvent) => {
    if (!currentTaskId.value || payload.taskId === currentTaskId.value) {
      doneEvent.value = payload
      running.value = false
      currentTaskId.value = ''
      setInfo('监视完成，报告已生成。')
      await Promise.all([loadHistory(), loadLogState(true)])
    }
  })

  backend.on('review:error', async (payload: { taskId: string; message: string }) => {
    if (!currentTaskId.value || payload.taskId === currentTaskId.value) {
      setError(payload.message)
      running.value = false
      currentTaskId.value = ''
      await loadLogState(true)
    }
  })
})

onUnmounted(() => {
  stopLogPolling()
})
</script>

<template>
  <div class="app-shell">
    <div class="hero">
      <div>
        <h1>AI代码监视</h1>
        <p>面向 GitHub PR、GitLab MR 与本地仓库改动的桌面审计工具，现已支持代理配置、GitLab 企业仓库拉取策略、实时日志与缓存清理。</p>
      </div>
      <div class="hero-tip small">建议流程：先拉取代码 → 自动识别分支 → 开始监视</div>
    </div>

    <div class="layout">
      <div class="left-column">
        <div class="card">
          <div class="panel-header">
            <div>
              <h2>仓库与监视配置</h2>
              <div class="small">远程仓库模式下，开始监视前请先点击“拉取代码”。</div>
            </div>
          </div>

          <div class="field">
            <label>MR / PR 链接</label>
            <input v-model="form.mrUrl" placeholder="https://gitlab.example.com/group/project/-/merge_requests/123" />
          </div>

          <div class="field">
            <label>仓库地址（可选，自动识别失败时填写）</label>
            <input v-model="form.repoUrl" placeholder="git@gitlab.example.com:group/project.git 或 https://host/group/project.git" />
          </div>

          <div class="field">
            <label>本地仓库路径（可选，本地模式）</label>
            <input v-model="form.localRepoPath" placeholder="D:/code/project 或 /Users/me/code/project" />
          </div>

          <datalist id="branch-options">
            <option v-for="branch in availableBranches" :key="branch" :value="branch" />
          </datalist>

          <div class="field-row">
            <div class="field">
              <label>源分支</label>
              <input v-model="form.sourceBranch" list="branch-options" placeholder="拉取代码后自动填充，可手动修改" />
            </div>
            <div class="field">
              <label>目标分支</label>
              <input v-model="form.targetBranch" list="branch-options" placeholder="默认优先 master / main / develop" />
            </div>
          </div>

          <div class="field-row">
            <div class="field">
              <label>配置文件路径</label>
              <input v-model="form.configPath" placeholder="./config.yaml" />
            </div>
            <div class="field">
              <label>工作区路径</label>
              <input v-model="form.workspaceDir" placeholder="./workspace" />
            </div>
          </div>

          <div class="action-grid">
            <button @click="pullCode" :disabled="running || pullingCode || clearingCache">
              {{ pullingCode ? '拉取中...' : '拉取代码' }}
            </button>
            <button @click="startReview" :disabled="running || pullingCode || clearingCache">
              {{ running ? '监视中...' : '开始监视' }}
            </button>
            <button class="secondary" @click="cancelReview" :disabled="!running">取消任务</button>
            <button class="secondary" @click="toggleLogs">{{ showLogs ? '隐藏日志' : '查看日志' }}</button>
            <button class="secondary danger" @click="clearCache" :disabled="running || pullingCode || clearingCache">
              {{ clearingCache ? '清理中...' : '清理缓存' }}
            </button>
            <button class="secondary" @click="loadHistory">刷新历史</button>
          </div>

          <div class="inline-meta">
            <div class="small">默认日志文件：{{ logState.logPath || '未生成' }}</div>
            <div class="small" v-if="availableBranches.length > 0">已识别分支：{{ availableBranches.length }} 个</div>
          </div>

          <div v-if="progress" class="status-panel">
            <div class="status-top">
              <div class="small">阶段：{{ progress.stage }}</div>
              <div class="small">进度：{{ progress.percent }}%</div>
            </div>
            <div class="progress"><div :style="{ width: `${progress.percent}%` }" /></div>
            <div class="small">{{ progress.message }}</div>
          </div>

          <div v-if="infoMessage" class="notice success">
            {{ infoMessage }}
          </div>

          <div v-if="errorMessage" class="notice error">
            {{ errorMessage }}
          </div>
        </div>

        <div class="card" v-if="showLogs">
          <div class="panel-header compact">
            <div>
              <h2>实时日志</h2>
              <div class="small">日志文件：{{ logState.logPath || '未生成' }}</div>
            </div>
            <div class="small">更新时间：{{ logState.updatedAt || '—' }}</div>
          </div>
          <pre ref="logViewerRef" class="log-viewer">{{ logState.content || '当前暂无日志内容。' }}</pre>
        </div>

        <div class="card">
          <div class="panel-header compact">
            <div>
              <h2>历史记录</h2>
              <div class="small">展示当前工作区下的最近报告</div>
            </div>
            <div class="small">{{ history.length }} 条</div>
          </div>
          <ul class="history-list">
            <li v-for="item in history" :key="item.taskId" class="history-item">
              <strong>{{ item.title }}</strong>
              <div class="small">{{ item.createdAt }}</div>
              <div class="small">{{ item.sourceRef }} → {{ item.targetRef }}</div>
              <div class="small">总问题 {{ item.totalIssues }} · 报告目录 {{ item.reportDir }}</div>
            </li>
            <li v-if="history.length === 0" class="placeholder">暂无历史记录</li>
          </ul>
        </div>
      </div>

      <div class="right-column">
        <div class="card">
          <div class="panel-header compact">
            <div>
              <h2>审计概览</h2>
              <div class="small">审计过程中会持续推送阶段进度，结束后展示最新报告。</div>
            </div>
            <div v-if="doneEvent" class="small">
              报告目录：{{ doneEvent.reportDir }}
            </div>
          </div>

          <div class="metrics">
            <div class="metric">
              <div class="label">高（高危+严重）</div>
              <div class="value">{{ report?.summary.high ?? 0 }}</div>
            </div>
            <div class="metric">
              <div class="label">中（一般）</div>
              <div class="value">{{ report?.summary.medium ?? 0 }}</div>
            </div>
            <div class="metric">
              <div class="label">低（建议）</div>
              <div class="value">{{ report?.summary.low ?? 0 }}</div>
            </div>
            <div class="metric">
              <div class="label">总计</div>
              <div class="value">{{ report?.summary.total ?? 0 }}</div>
            </div>
          </div>

          <div v-if="report" class="artifact-list small">
            <div>HTML：{{ doneEvent?.htmlPath }}</div>
            <div>Markdown：{{ doneEvent?.markdownPath }}</div>
            <div>JSON：{{ doneEvent?.jsonPath }}</div>
          </div>
        </div>

        <div class="card" v-if="report">
          <div class="panel-header compact">
            <div>
              <h2>质量与健康度</h2>
              <div class="small">报告会保留规则预扫与 AI 审计的综合结论。</div>
            </div>
          </div>

          <div class="metrics metrics-5">
            <div class="metric">
              <div class="label">安全性</div>
              <div class="value">{{ report.health.security }}</div>
            </div>
            <div class="metric">
              <div class="label">性能</div>
              <div class="value">{{ report.health.performance }}</div>
            </div>
            <div class="metric">
              <div class="label">健壮性</div>
              <div class="value">{{ report.health.robustness }}</div>
            </div>
            <div class="metric">
              <div class="label">可维护性</div>
              <div class="value">{{ report.health.maintainability }}</div>
            </div>
            <div class="metric">
              <div class="label">框架最佳实践</div>
              <div class="value">{{ report.health.frameworkPractice }}</div>
            </div>
          </div>

          <h3 class="sub-title">其他说明</h3>
          <ul class="note-list">
            <li v-for="note in report.notes" :key="note">{{ note }}</li>
          </ul>
        </div>

        <div class="card">
          <div class="panel-header compact">
            <div>
              <h2>问题清单</h2>
              <div class="small">输出最新一次监视结果。</div>
            </div>
            <div class="small">{{ findings.length }} 条</div>
          </div>

          <div v-if="findings.length === 0" class="placeholder">暂无问题结果，等待监视完成。</div>
          <div v-for="item in findings" :key="item.id || `${item.file}-${item.lineStart}-${item.title}`" class="finding">
            <h3>{{ item.id || '未编号' }}：{{ item.title }}（{{ item.severity }}）</h3>
            <div class="chips">
              <span class="chip">{{ item.category }}</span>
              <span class="chip">{{ item.confidence }}</span>
              <span class="chip">{{ item.file }}:{{ item.lineStart }}-{{ item.lineEnd }}</span>
            </div>
            <p><strong>详细描述：</strong>{{ item.description }}</p>
            <p><strong>影响分析：</strong>{{ item.impact }}</p>
            <p><strong>证据：</strong>{{ item.evidence }}</p>
            <p><strong>修复建议：</strong>{{ item.recommendation }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
