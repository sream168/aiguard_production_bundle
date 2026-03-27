<template>
  <div v-if="progress" class="progress-panel">
    <div class="progress-header">
      <div>
        <div class="stage">{{ progress.stage }}</div>
        <p class="progress-message">{{ progress.message }}</p>
      </div>
      <div class="progress-right">
        <span class="percent">{{ progress.percent }}%</span>
        <div class="risk-pills">
          <span class="pill high">高 {{ progress.high }}</span>
          <span class="pill medium">中 {{ progress.medium }}</span>
          <span class="pill low">低 {{ progress.low }}</span>
        </div>
      </div>
    </div>
    <div class="progress-bar">
      <div class="progress-fill" :style="{ width: `${progress.percent}%` }"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ProgressEvent } from '../types'

defineProps<{
  progress: ProgressEvent | null
}>()
</script>

<style scoped>
.progress-panel {
  background: rgba(10, 18, 35, 0.82);
  border: 1px solid rgba(99, 102, 241, 0.18);
  padding: 1rem 1.1rem;
  border-radius: 1.2rem;
  box-shadow: 0 14px 30px rgba(0, 0, 0, 0.18);
}

.progress-header {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 0.8rem;
}

.stage {
  font-weight: 700;
  color: #93c5fd;
  font-size: 1rem;
}

.progress-message {
  margin: 0.35rem 0 0;
  color: #cbd5e1;
  font-size: 0.9rem;
}

.progress-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.55rem;
}

.percent {
  color: #e2e8f0;
  font-size: 1.15rem;
  font-weight: 700;
}

.risk-pills {
  display: flex;
  gap: 0.45rem;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.pill {
  display: inline-flex;
  align-items: center;
  padding: 0.28rem 0.6rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 700;
}

.pill.high {
  background: rgba(190, 24, 93, 0.14);
  color: #fecdd3;
}

.pill.medium {
  background: rgba(217, 119, 6, 0.14);
  color: #fde68a;
}

.pill.low {
  background: rgba(13, 148, 136, 0.14);
  color: #ccfbf1;
}

.progress-bar {
  height: 12px;
  background: rgba(15, 23, 42, 0.92);
  border-radius: 999px;
  overflow: hidden;
  border: 1px solid rgba(96, 165, 250, 0.14);
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #3b82f6, #22c55e, #8b5cf6);
  transition: width 0.3s ease;
}

@media (max-width: 860px) {
  .progress-header {
    flex-direction: column;
  }

  .progress-right {
    align-items: flex-start;
  }

  .risk-pills {
    justify-content: flex-start;
  }
}
</style>
