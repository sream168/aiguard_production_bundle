<template>
  <div class="review-form">
    <div class="panel-header compact">
      <div>
        <h2>监视配置</h2>
        <p class="small">支持 GitLab MR、GitHub PR 和本地仓库模式。仓库地址会根据 MR/PR 链接自动推导，仍可手动改写。</p>
      </div>
    </div>

    <div class="field">
      <label>MR/PR 链接</label>
      <input v-model="form.mrUrl" placeholder="https://gitlab.example.com/group/project/-/merge_requests/123" />
    </div>

    <div class="field">
      <label>仓库地址（自动推导，可编辑）</label>
      <input
        v-model="form.repoUrl"
        placeholder="git@gitlab.example.com:group/project.git"
        @input="$emit('repo-url-edited', form.repoUrl)"
      />
      <div class="field-hint" v-if="repoSuggestionMessage || repoSuggestionCandidates.length > 0">
        <span class="hint-pill" :class="repoSuggestionResolvedByMr ? 'primary' : 'muted'">
          {{ repoSuggestionResolvedByMr ? '已根据 MR 自动推导' : '当前为手动或未识别' }}
        </span>
        <span class="small">{{ repoSuggestionMessage }}</span>
      </div>
      <div class="candidate-list" v-if="repoSuggestionCandidates.length > 1">
        <span class="small">候选地址：</span>
        <code v-for="candidate in repoSuggestionCandidates" :key="candidate">{{ candidate }}</code>
      </div>
    </div>

    <div class="field">
      <label>本地仓库路径（可选）</label>
      <input v-model="form.localRepoPath" placeholder="/path/to/repo" />
    </div>

    <datalist id="branch-options">
      <option v-for="branch in availableBranches" :key="branch" :value="branch" />
    </datalist>

    <div class="field-row">
      <div class="field">
        <label>源分支</label>
        <input v-model="form.sourceBranch" list="branch-options" placeholder="feature/xxx" />
      </div>
      <div class="field">
        <label>目标分支</label>
        <input v-model="form.targetBranch" list="branch-options" placeholder="main / master / develop" />
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
      <button @click="$emit('pull-code')" :disabled="disabled">
        {{ pullingCode ? '拉取中...' : '拉取代码' }}
      </button>
      <button class="gradient" @click="$emit('start-review')" :disabled="disabled">
        {{ running ? '监视中...' : '开始监视' }}
      </button>
      <button class="secondary" @click="$emit('cancel-review')" :disabled="!running">取消任务</button>
      <button class="secondary" @click="$emit('toggle-logs')">{{ showLogs ? '隐藏日志' : '查看日志' }}</button>
      <button class="secondary danger" @click="$emit('clear-cache')" :disabled="disabled">
        {{ clearingCache ? '清理中...' : '清理缓存' }}
      </button>
      <button class="secondary" @click="$emit('load-history')">刷新历史</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { StartReviewRequest } from '../types'

defineProps<{
  form: StartReviewRequest
  availableBranches: string[]
  running: boolean
  pullingCode: boolean
  clearingCache: boolean
  showLogs: boolean
  disabled: boolean
  repoSuggestionMessage: string
  repoSuggestionCandidates: string[]
  repoSuggestionResolvedByMr: boolean
}>()

defineEmits<{
  'pull-code': []
  'start-review': []
  'cancel-review': []
  'toggle-logs': []
  'clear-cache': []
  'load-history': []
  'repo-url-edited': [value: string]
}>()
</script>

<style scoped>
.review-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.field label {
  font-size: 0.9rem;
  font-weight: 600;
  color: #dbe7ff;
}

.field input {
  width: 100%;
  padding: 0.95rem 1rem;
  background: rgba(7, 14, 30, 0.88);
  border: 1px solid rgba(99, 102, 241, 0.24);
  border-radius: 1rem;
  color: #eef2ff;
  font-size: 0.95rem;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
}

.field input:focus {
  outline: none;
  border-color: rgba(96, 165, 250, 0.72);
  box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.14);
}

.field-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.field-hint {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.hint-pill {
  display: inline-flex;
  align-items: center;
  padding: 0.3rem 0.7rem;
  border-radius: 999px;
  font-size: 0.78rem;
  border: 1px solid rgba(148, 163, 184, 0.24);
  background: rgba(30, 41, 59, 0.78);
}

.hint-pill.primary {
  color: #bfdbfe;
  border-color: rgba(96, 165, 250, 0.28);
  background: rgba(30, 64, 175, 0.16);
}

.hint-pill.muted {
  color: #cbd5e1;
}

.candidate-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem;
  align-items: center;
}

.candidate-list code {
  background: rgba(15, 23, 42, 0.92);
  border: 1px solid rgba(99, 102, 241, 0.14);
  color: #c4b5fd;
  padding: 0.2rem 0.45rem;
  border-radius: 0.55rem;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0.8rem;
  margin-top: 0.2rem;
}

button {
  padding: 0.82rem 1rem;
  border: none;
  border-radius: 1rem;
  font-size: 0.92rem;
  color: #fff;
  cursor: pointer;
  transition: transform 0.18s ease, box-shadow 0.18s ease, opacity 0.18s ease;
  background: linear-gradient(135deg, #2563eb, #7c3aed);
  box-shadow: 0 10px 22px rgba(37, 99, 235, 0.24);
}

button.gradient {
  background: linear-gradient(135deg, #22c55e, #06b6d4 45%, #6366f1 100%);
  box-shadow: 0 12px 24px rgba(34, 197, 94, 0.18);
}

button:hover:not(:disabled) {
  transform: translateY(-1px);
}

button:disabled {
  opacity: 0.55;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

button.secondary {
  background: rgba(51, 65, 85, 0.88);
  box-shadow: none;
}

button.secondary:hover:not(:disabled) {
  background: rgba(71, 85, 105, 0.95);
}

button.danger {
  background: linear-gradient(135deg, rgba(220, 38, 38, 0.9), rgba(190, 24, 93, 0.9));
}

@media (max-width: 860px) {
  .field-row,
  .action-grid {
    grid-template-columns: 1fr;
  }
}
</style>
