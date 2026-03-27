<template>
  <div class="history-wrapper">
    <div class="panel-header compact">
      <div>
        <h2>历史记录</h2>
        <p class="small">修复了历史列表字段映射问题，现在会稳定显示总问题数及高/中/低分布。</p>
      </div>
      <div class="small">共 {{ history.length }} 条</div>
    </div>

    <div v-if="history.length === 0" class="empty">暂无历史记录</div>
    <div v-else class="history-items">
      <div v-for="item in history" :key="item.taskId" class="history-item-card">
        <div class="history-header">
          <div>
            <div class="task-id">{{ item.taskId.slice(0, 8) }}</div>
            <div class="history-title">{{ item.title || 'AI代码监视报告' }}</div>
          </div>
          <span class="timestamp">{{ item.createdAt }}</span>
        </div>

        <div class="history-meta">
          <span>{{ item.sourceRef }} → {{ item.targetRef }}</span>
          <span class="repo-url" :title="item.repoUrl">{{ item.repoUrl || '未记录仓库地址' }}</span>
        </div>

        <div class="history-summary-grid">
          <div class="summary-chip total">总计 {{ item.totalIssues }}</div>
          <div class="summary-chip high">高 {{ item.summary.high }}</div>
          <div class="summary-chip medium">中 {{ item.summary.medium }}</div>
          <div class="summary-chip low">低 {{ item.summary.low }}</div>
        </div>

        <div class="history-actions">
          <button class="mini primary" @click="$emit('open-report', item)">查看报告</button>
          <button class="mini secondary" @click="$emit('open-dir', item)">打开目录</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { HistoryItem } from '../types'

defineProps<{
  history: HistoryItem[]
}>()

defineEmits<{
  'open-report': [item: HistoryItem]
  'open-dir': [item: HistoryItem]
}>()
</script>

<style scoped>
.history-wrapper {
  display: flex;
  flex-direction: column;
  gap: 0.9rem;
}

.empty {
  color: #94a3b8;
  font-size: 0.95rem;
  padding: 1.25rem;
  text-align: center;
  background: rgba(15, 23, 42, 0.78);
  border: 1px solid rgba(99, 102, 241, 0.16);
  border-radius: 1rem;
}

.history-items {
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
}

.history-item-card {
  background: rgba(10, 18, 35, 0.82);
  border: 1px solid rgba(99, 102, 241, 0.18);
  border-radius: 1.2rem;
  padding: 1rem 1rem 0.95rem;
  box-shadow: 0 14px 30px rgba(0, 0, 0, 0.18);
}

.history-header,
.history-meta,
.history-actions {
  display: flex;
  justify-content: space-between;
  gap: 0.75rem;
}

.history-header {
  align-items: flex-start;
  margin-bottom: 0.55rem;
}

.history-title {
  color: #e2e8f0;
  font-size: 0.93rem;
  margin-top: 0.2rem;
}

.task-id {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  color: #93c5fd;
  font-weight: 700;
}

.timestamp {
  color: #94a3b8;
  font-size: 0.82rem;
  white-space: nowrap;
}

.history-meta {
  color: #cbd5e1;
  font-size: 0.88rem;
  margin-bottom: 0.8rem;
  flex-wrap: wrap;
}

.repo-url {
  max-width: 52%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #a5b4fc;
}

.history-summary-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 0.65rem;
  margin-bottom: 0.85rem;
}

.summary-chip {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 42px;
  border-radius: 0.9rem;
  font-size: 0.88rem;
  font-weight: 700;
  border: 1px solid rgba(99, 102, 241, 0.12);
}

.summary-chip.total {
  color: #dbeafe;
  background: rgba(30, 64, 175, 0.14);
}

.summary-chip.high {
  color: #fecdd3;
  background: rgba(190, 24, 93, 0.14);
}

.summary-chip.medium {
  color: #fde68a;
  background: rgba(217, 119, 6, 0.14);
}

.summary-chip.low {
  color: #ccfbf1;
  background: rgba(13, 148, 136, 0.14);
}

.history-actions {
  justify-content: flex-end;
}

button.mini {
  border: none;
  border-radius: 0.9rem;
  padding: 0.55rem 0.85rem;
  color: white;
  cursor: pointer;
}

button.primary {
  background: linear-gradient(135deg, #3b82f6, #7c3aed);
}

button.secondary {
  background: rgba(51, 65, 85, 0.9);
}

@media (max-width: 860px) {
  .history-summary-grid {
    grid-template-columns: 1fr 1fr;
  }

  .history-header,
  .history-meta,
  .history-actions {
    flex-direction: column;
  }

  .repo-url {
    max-width: 100%;
  }
}
</style>
