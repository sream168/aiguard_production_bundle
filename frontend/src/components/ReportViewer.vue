<template>
  <div class="report-shell" v-if="report">
    <div class="panel-header">
      <div>
        <h2>审计报告</h2>
        <p class="small">默认折叠显示 3 行摘要，点击展开后可查看证据、影响、建议和可选修正代码片段。</p>
      </div>
      <div class="report-actions">
        <button class="primary" @click="$emit('open-report')">查看报告</button>
        <button class="secondary" @click="$emit('open-report-dir')">打开目录</button>
      </div>
    </div>

    <div class="summary-grid">
      <div class="metric-card total">
        <span class="label">总问题数</span>
        <span class="value">{{ report.summary.total }}</span>
      </div>
      <div class="metric-card danger">
        <span class="label">高（高危+严重）</span>
        <span class="value">{{ report.summary.high }}</span>
      </div>
      <div class="metric-card warn">
        <span class="label">中（一般）</span>
        <span class="value">{{ report.summary.medium }}</span>
      </div>
      <div class="metric-card low">
        <span class="label">低（建议）</span>
        <span class="value">{{ report.summary.low }}</span>
      </div>
    </div>

    <div class="health-grid">
      <div class="health-item"><span>安全性</span><strong>{{ report.health.security }}</strong></div>
      <div class="health-item"><span>性能</span><strong>{{ report.health.performance }}</strong></div>
      <div class="health-item"><span>健壮性</span><strong>{{ report.health.robustness }}</strong></div>
      <div class="health-item"><span>可维护性</span><strong>{{ report.health.maintainability }}</strong></div>
      <div class="health-item"><span>框架实践</span><strong>{{ report.health.frameworkPractice }}</strong></div>
    </div>

    <div v-if="report.findings.length === 0" class="empty-report">本次未发现可定位问题。</div>

    <div v-else class="finding-list">
      <article v-for="finding in report.findings" :key="finding.id || `${finding.file}-${finding.lineStart}-${finding.title}`" class="finding-card">
        <div class="finding-header">
          <div>
            <div class="tag-row">
              <span class="tag" :class="severityClass(finding.severity)">{{ finding.severity }}</span>
              <span class="tag neutral">{{ finding.category }}</span>
              <span class="tag neutral">{{ finding.confidence }}</span>
            </div>
            <h3>{{ finding.title }}</h3>
          </div>
          <div class="location">{{ finding.file }}:{{ finding.lineStart }}<template v-if="finding.lineEnd !== finding.lineStart">-{{ finding.lineEnd }}</template></div>
        </div>

        <p class="description" :class="{ collapsed: !isExpanded(finding) }">{{ finding.description }}</p>

        <button class="toggle-button" @click="toggle(finding)">
          {{ isExpanded(finding) ? '收起详情' : '展开详情' }}
        </button>

        <div v-if="isExpanded(finding)" class="detail-block">
          <div class="detail-item">
            <span class="detail-label">影响分析</span>
            <p>{{ finding.impact || '未提供' }}</p>
          </div>
          <div class="detail-item">
            <span class="detail-label">证据</span>
            <p>{{ finding.evidence || '未提供' }}</p>
          </div>
          <div class="detail-item">
            <span class="detail-label">修复建议</span>
            <p>{{ finding.recommendation || '未提供' }}</p>
          </div>
          <pre v-if="finding.recommendationCode" class="code-snippet"><code>{{ finding.recommendationCode }}</code></pre>
        </div>
      </article>
    </div>

    <div v-if="report.notes?.length" class="note-panel">
      <h3>补充说明</h3>
      <ul>
        <li v-for="(note, idx) in report.notes" :key="idx">{{ note }}</li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Finding, Report } from '../types'

const expanded = ref<Record<string, boolean>>({})

function findingKey(finding: Finding) {
  return finding.id || `${finding.file}-${finding.lineStart}-${finding.lineEnd}-${finding.title}`
}

function toggle(finding: Finding) {
  const key = findingKey(finding)
  expanded.value[key] = !expanded.value[key]
}

function isExpanded(finding: Finding) {
  return !!expanded.value[findingKey(finding)]
}

function severityClass(severity: string) {
  switch (severity) {
    case '高危':
      return 'high'
    case '严重':
      return 'severe'
    case '一般':
      return 'medium'
    default:
      return 'low'
  }
}

defineProps<{
  report: Report | null
  htmlPath: string
  reportDir: string
}>()

defineEmits<{
  'open-report': []
  'open-report-dir': []
}>()
</script>

<style scoped>
.report-shell {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.report-actions {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

button {
  border: none;
  border-radius: 1rem;
  padding: 0.72rem 1rem;
  color: white;
  cursor: pointer;
}

button.primary {
  background: linear-gradient(135deg, #2563eb, #7c3aed);
}

button.secondary {
  background: rgba(51, 65, 85, 0.92);
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 0.85rem;
}

.metric-card,
.health-item,
.finding-card,
.note-panel,
.empty-report {
  background: rgba(10, 18, 35, 0.82);
  border: 1px solid rgba(99, 102, 241, 0.18);
  box-shadow: 0 14px 30px rgba(0, 0, 0, 0.18);
}

.metric-card {
  border-radius: 1.2rem;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.45rem;
}

.metric-card .label {
  font-size: 0.85rem;
  color: #94a3b8;
}

.metric-card .value {
  font-size: 1.75rem;
  font-weight: 800;
}

.metric-card.total .value { color: #dbeafe; }
.metric-card.danger .value { color: #fda4af; }
.metric-card.warn .value { color: #fde68a; }
.metric-card.low .value { color: #ccfbf1; }

.health-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 0.85rem;
}

.health-item {
  border-radius: 1rem;
  padding: 0.9rem 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.45rem;
}

.health-item span {
  color: #94a3b8;
  font-size: 0.85rem;
}

.health-item strong {
  font-size: 1.3rem;
  color: #dbeafe;
}

.finding-list {
  display: flex;
  flex-direction: column;
  gap: 0.95rem;
}

.finding-card {
  border-radius: 1.25rem;
  padding: 1rem 1.05rem;
}

.finding-header {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  align-items: flex-start;
}

.finding-header h3 {
  margin: 0;
  color: #f8fafc;
  line-height: 1.45;
}

.tag-row {
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem;
  margin-bottom: 0.65rem;
}

.tag {
  display: inline-flex;
  align-items: center;
  padding: 0.28rem 0.62rem;
  border-radius: 999px;
  font-size: 0.78rem;
  font-weight: 700;
}

.tag.high {
  background: rgba(217, 70, 239, 0.18);
  color: #f5d0fe;
}

.tag.severe {
  background: rgba(244, 63, 94, 0.18);
  color: #fecdd3;
}

.tag.medium {
  background: rgba(250, 204, 21, 0.16);
  color: #fde68a;
}

.tag.low {
  background: rgba(34, 197, 94, 0.16);
  color: #bbf7d0;
}

.tag.neutral {
  background: rgba(59, 130, 246, 0.14);
  color: #bfdbfe;
}

.location {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 0.82rem;
  color: #94a3b8;
}

.description {
  margin: 0.8rem 0 0;
  color: #dbe7ff;
  line-height: 1.7;
}

.description.collapsed {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.toggle-button {
  margin-top: 0.85rem;
  padding: 0;
  background: transparent;
  color: #93c5fd;
  border-radius: 0;
}

.detail-block {
  margin-top: 0.9rem;
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
  border-top: 1px dashed rgba(148, 163, 184, 0.18);
  padding-top: 0.9rem;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.detail-item p {
  margin: 0;
  color: #dbe7ff;
  line-height: 1.7;
}

.detail-label {
  color: #94a3b8;
  font-size: 0.85rem;
  font-weight: 700;
}

.code-snippet {
  margin: 0;
  border-radius: 1rem;
  background: rgba(2, 6, 23, 0.96);
  border: 1px solid rgba(96, 165, 250, 0.14);
  padding: 0.95rem;
  color: #dbeafe;
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-word;
}

.note-panel,
.empty-report {
  border-radius: 1.2rem;
  padding: 1rem 1.1rem;
}

.note-panel h3 {
  margin-top: 0;
}

.note-panel ul {
  margin: 0;
  padding-left: 1.1rem;
  color: #dbe7ff;
}

.note-panel li + li {
  margin-top: 0.4rem;
}

@media (max-width: 1024px) {
  .summary-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .health-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 760px) {
  .summary-grid,
  .health-grid {
    grid-template-columns: 1fr;
  }

  .finding-header {
    flex-direction: column;
  }

  .report-actions {
    width: 100%;
  }
}
</style>
