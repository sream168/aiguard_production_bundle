import type {
  CacheClearResult,
  HistoryItem,
  LogState,
  PrepareRepositoryResponse,
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

  on(name: string, callback: (payload: any) => void) {
    if (window.runtime?.EventsOn) {
      window.runtime.EventsOn(name, callback)
    }
  }
}
