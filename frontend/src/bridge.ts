import type {
  CacheClearResult,
  HistoryItem,
  LogState,
  OpenPathResult,
  PrepareRepositoryResponse,
  RepositorySuggestion,
  RuntimeContextRequest,
  StartReviewRequest
} from './types'

declare global {
  interface Window {
    go?: any
    runtime?: {
      EventsOn?: (name: string, callback: (payload: any) => void) => void
    }
  }
}

export const backend = {
  async suggestRepository(payload: StartReviewRequest): Promise<RepositorySuggestion> {
    if (window.go?.main?.App?.SuggestRepository) {
      return await window.go.main.App.SuggestRepository(payload)
    }
    return {
      repoUrl: payload.repoUrl,
      candidates: payload.repoUrl ? [payload.repoUrl] : [],
      resolvedByMr: false,
      manualRepoUrl: payload.repoUrl,
      message: ''
    }
  },

  async pullCode(payload: StartReviewRequest): Promise<PrepareRepositoryResponse> {
    if (window.go?.main?.App?.PullCode) {
      return await window.go.main.App.PullCode(payload)
    }
    return {
      repoUrl: payload.repoUrl,
      sourceBranch: payload.sourceBranch,
      targetBranch: payload.targetBranch,
      availableBranches: [],
      logPath: `${payload.workspaceDir || './workspace'}/logs/aiguard.log`,
      message: 'Mock pull completed.'
    }
  },

  async startReview(payload: StartReviewRequest): Promise<string> {
    if (window.go?.main?.App?.StartReview) {
      return await window.go.main.App.StartReview(payload)
    }
    return crypto?.randomUUID?.() ?? `mock-${Date.now()}`
  },

  async cancelReview(taskId: string): Promise<void> {
    if (window.go?.main?.App?.CancelReview) {
      await window.go.main.App.CancelReview(taskId)
    }
  },

  async listHistory(payload: RuntimeContextRequest): Promise<HistoryItem[]> {
    if (window.go?.main?.App?.ListHistory) {
      return await window.go.main.App.ListHistory(payload)
    }
    return []
  },

  async getLogState(payload: RuntimeContextRequest): Promise<LogState> {
    if (window.go?.main?.App?.GetLogState) {
      return await window.go.main.App.GetLogState(payload)
    }
    return {
      logPath: `${payload.workspaceDir || './workspace'}/logs/aiguard.log`,
      content: '',
      updatedAt: ''
    }
  },

  async clearCache(payload: RuntimeContextRequest): Promise<CacheClearResult> {
    if (window.go?.main?.App?.ClearCache) {
      return await window.go.main.App.ClearCache(payload)
    }
    return {
      workspaceDir: payload.workspaceDir || './workspace',
      logPath: `${payload.workspaceDir || './workspace'}/logs/aiguard.log`,
      removed: [],
      message: 'Mock cache clear completed.'
    }
  },

  async openReport(payload: { htmlPath: string; reportDir: string }): Promise<OpenPathResult> {
    if (window.go?.main?.App?.OpenReport) {
      return await window.go.main.App.OpenReport(payload)
    }
    return {
      path: payload.htmlPath || payload.reportDir,
      mode: payload.htmlPath ? 'html' : 'directory',
      message: 'Mock report open completed.'
    }
  },

  async openReportDirectory(payload: { reportDir: string }): Promise<OpenPathResult> {
    if (window.go?.main?.App?.OpenReportDirectory) {
      return await window.go.main.App.OpenReportDirectory(payload)
    }
    return {
      path: payload.reportDir,
      mode: 'directory',
      message: 'Mock directory open completed.'
    }
  },

  on(name: string, callback: (payload: any) => void) {
    if (window.runtime?.EventsOn) {
      window.runtime.EventsOn(name, callback)
    }
  }
}
